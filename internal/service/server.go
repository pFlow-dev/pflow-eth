package service

import (
	"context"
	"encoding/json"
	"fmt"
	rice "github.com/GeertJohan/go.rice"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/logWriter"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pflow-dev/pflow-eth/internal/config"
	"github.com/pflow-dev/pflow-eth/metamodel"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Service struct {
	*Sessions
	Client          *ethclient.Client
	Logger          *log.Logger
	Router          *mux.Router
	applicationPage *template.Template
	apm             *newrelic.Application
	cancelFunc      context.CancelFunc
}

type Server interface {
	Serve(box *rice.Box)
}

func New() Server {
	s := &Service{
		Sessions: &Sessions{
			sessions: make(map[string]*SessionData),
		},
		Router: mux.NewRouter(),
	}
	var err error
	if err != nil {
		log.Fatal(err)
	}
	s.Client, err = ethclient.Dial(config.Endpoint)
	if err != nil {
		log.Fatal(err)
	}

	if config.NewRelicLicense != "" {
		s.apm, _ = newrelic.NewApplication(
			newrelic.ConfigAppName(config.NewRelicApp),
			newrelic.ConfigLicense(config.NewRelicLicense),
			newrelic.ConfigAppLogForwardingEnabled(true),
		)
		writer := logWriter.New(os.Stdout, s.apm)
		s.Logger = log.New(writer, "", log.Default().Flags())
		s.Logger.Printf("NewRelic license set, APM enabled %s\n", config.NewRelicApp)
	} else {
		s.Logger = log.Default()
		s.Logger.Print("NewRelic license not set, skipping APM, disable browser tracking\n")
	}
	s.applicationPage = template.Must(template.New("index.html").Parse(s.appSource()))
	s.Logger.Printf("Listening on %s:%s\n", config.Host, config.Port)
	return s
}

type HandlerWithVars = func(vars map[string]string, w http.ResponseWriter, r *http.Request)

type VarsFactory = func(r *http.Request) map[string]string

func (s *Service) WrapHandler(pattern string, handler HandlerWithVars) {
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

// Event record custom event in apm and log
func (s *Service) Event(eventType string, params map[string]interface{}) {
	if s.apm != nil {
		s.apm.RecordCustomEvent(eventType, params)
	}
	data, _ := json.Marshal(params)
	paddedEventType := fmt.Sprintf("%-25s", eventType)
	s.Logger.Printf("| %s | %s\n", paddedEventType, data)
}

func (s *Service) Serve(box *rice.Box) {

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
	}()

	go s.CleanupSessions()
	s.applicationRoutes(box)
	s.apiRoutes()
	err := http.ListenAndServe(config.Host+":"+config.Port, s.Router)
	if err != nil {
		panic(err)
	}
}

func (s *Service) apiRoutes() {
	s.Router.HandleFunc("/v0/ping", s.PingHandler)
	s.Router.HandleFunc("/v0/authenticate", s.loginHandler)
}

func (s *Service) getContractSequence() int64 {
	address := common.HexToAddress(config.Address)
	contract, _ := metamodel.NewMetamodel(address, s.Client)
	res, err := contract.Sequence(nil)
	if err != nil {
		return -1
	}
	return res.Int64()
}

func (s *Service) getLatestBlockNumber() int64 {
	header, err := s.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return -1
	}
	return header.Number.Int64()
}
