package service

import (
	"context"
	"database/sql"
	"encoding/json"
	rice "github.com/GeertJohan/go.rice"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/logWriter"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pflow-dev/pflow-eth/config"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Options struct {
	Port            string
	Host            string
	Url             string
	DbPath          string
	NewRelicLicense string
	NewRelicApp     string
}

type Server struct {
	NodeDb     *sql.DB
	Client     *ethclient.Client
	Logger     *log.Logger
	Options    Options
	Router     *mux.Router
	indexPage  *template.Template
	apm        *newrelic.Application
	cancelFunc context.CancelFunc
}

func New(options Options) *Server {
	s := &Server{
		Options: options,
		Router:  mux.NewRouter(),
	}
	s.NodeDb, _ = sql.Open("postgres", config.DbConn)
	s.Client, _ = ethclient.Dial(config.Endpoint)

	if s.Options.NewRelicLicense != "" {
		s.apm, _ = newrelic.NewApplication(
			newrelic.ConfigAppName(s.Options.NewRelicApp),
			newrelic.ConfigLicense(s.Options.NewRelicLicense),
			newrelic.ConfigAppLogForwardingEnabled(true),
		)
		writer := logWriter.New(os.Stdout, s.apm)
		s.Logger = log.New(writer, "", log.Default().Flags())
		s.Logger.Printf("NewRelic license set, APM enabled %s\n", s.Options.NewRelicApp)
	} else {
		s.Logger = log.Default()
		s.Logger.Print("NewRelic license not set, skipping APM, disable browser tracking\n")
	}
	indexSource := s.IndexTemplateSource()
	s.indexPage = template.Must(template.New("index.html").Parse(indexSource))
	s.Logger.Printf("DBPath: %s\n", s.Options.DbPath)
	s.Logger.Printf("Listening on %s:%s\n", s.Options.Host, s.Options.Port)
	return s
}

type HandlerWithVars = func(vars map[string]string, w http.ResponseWriter, r *http.Request)

type VarsFactory = func(r *http.Request) map[string]string

func WithVars(handler HandlerWithVars, getVarsFunc VarsFactory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(getVarsFunc(r), w, r)
	}
}

func (s *Server) WrapHandler(pattern string, handler HandlerWithVars) {
	if s.apm != nil {
		s.Router.HandleFunc(newrelic.WrapHandleFunc(s.apm, pattern, func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			handler(vars, w, r)
		}))
	} else {
		s.Router.HandleFunc(
			pattern,
			func(w http.ResponseWriter, r *http.Request) {
				vars := mux.Vars(r)
				handler(vars, w, r)
			})
	}
}
func (s *Server) Event(eventType string, params map[string]interface{}) {
	if s.apm != nil {
		s.apm.RecordCustomEvent(eventType, params)
	}
	data, _ := json.Marshal(params)
	s.Logger.Printf("event %s %s\n", eventType, data)
}

func (*Server) GetState(r *http.Request) (state []int64, ok bool) {
	// FIXME
	return nil, false
}

func (s *Server) Serve(box *rice.Box) {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFunc = cancel

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		cancel()
	}()

	go s.runSequencer(ctx)
	s.applicationRoutes(box)
	s.ApiRoutes()
	err := http.ListenAndServe(s.Options.Host+":"+s.Options.Port, s.Router)
	if err != nil {
		panic(err)
	}
}

// add pages here if using react-router
var pageRoutes = []string{
	"/",
}

func (s *Server) AppPage(vars map[string]string, w http.ResponseWriter, r *http.Request) {
	if err := s.indexPage.Execute(w, nil); err != nil {
		s.Logger.Printf("Error rendering index page: %v\n", err)
	}
}

func (s *Server) applicationRoutes(box *rice.Box) {
	for _, route := range pageRoutes {
		s.WrapHandler(route, s.AppPage)
	}

	s.Router.HandleFunc("/{file}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		http.StripPrefix("/", http.FileServer(box.HTTPBox())).ServeHTTP(w, r)
	})

	s.Router.HandleFunc("/static/js/{jsBuild}",
		func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			f, boxErr := box.Open("static/js/" + vars["jsBuild"])
			if boxErr != nil {
				http.Error(w, boxErr.Error(), http.StatusNotFound)
				return
			}

			fileInfo, fileErr := f.Stat()
			if fileErr != nil {
				http.Error(w, fileErr.Error(), http.StatusNotFound)
				return
			}
			w.Header().Set("Cache-Control", "public, max-age=31536000")
			http.ServeContent(w, r, "main."+vars["jsBuild"]+".js", fileInfo.ModTime(), f)
		})

	s.Router.HandleFunc("/static/css/{cssBuild}",
		func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			f, boxErr := box.Open("static/css/" + vars["cssBuild"])
			if boxErr != nil {
				http.Error(w, boxErr.Error(), http.StatusNotFound)
				return
			}
			fileInfo, fileErr := f.Stat()
			if fileErr != nil {
				http.Error(w, fileErr.Error(), http.StatusNotFound)
				return
			}
			w.Header().Set("Cache-Control", "public, max-age=31536000")
			http.ServeContent(w, r, "main."+vars["cssBuild"]+".css", fileInfo.ModTime(), f)
		})

}
