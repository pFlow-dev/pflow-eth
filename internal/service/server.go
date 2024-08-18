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
	"strings"
)

type Service struct {
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

func (s *Service) Metric(name string, value float64) {
	if s.apm != nil {
		s.apm.RecordCustomMetric(name, value)
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
	s.applicationRoutes(box)
	s.apiRoutes()
	err := http.ListenAndServe(config.Host+":"+config.Port, s.Router)
	s.cancelFunc()
	if err != nil {
		panic(err)
	}
}

func (s *Service) apiRoutes() {
	s.Router.HandleFunc("/v0/authenticate", s.loginHandler)
	s.Router.HandleFunc("/v0/faucet", s.FaucetHandler)
}

func (s *Service) getContractSequence(address common.Address) int64 {
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

func (s *Service) getNonce(address common.Address) string {
	nonce, err := s.Client.PendingNonceAt(context.Background(), address)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%d", nonce)
}

func (*Service) getNetwork(r *http.Request) string {
	host := r.Host
	hostname := strings.Split(host, ":")[0]

	if hostname == "127.0.0.1" || hostname == "localhost" {
		return "hardhat"
	}

	subdomain := strings.Split(hostname, ".")[0]
	if network, exists := config.SubdomainToNetwork[subdomain]; exists {
		return network
	}
	return "unknown"
}

// construct a uuid build string from the JS and CSS build numbers
func getBuild() string {
	hexString := config.JsBuild + "000000000001" + "0000" + config.CssBuild

	return fmt.Sprintf("%s-%s-%s-%s-%s",
		hexString[0:8], hexString[8:12], hexString[12:16], hexString[16:20], hexString[20:32])
}

var buildUuid = getBuild()

func (s *Service) getSessionData(sessionID string) (map[string]interface{}, error) {
	status := make(map[string]interface{})
	status["status"] = "ok"
	status["build"] = buildUuid
	return status, nil
}
