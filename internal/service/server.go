package service

import (
	"context"
	"database/sql"
	"encoding/json"
	rice "github.com/GeertJohan/go.rice"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/logWriter"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pflow-dev/pflow-eth/internal/config"
	"github.com/pflow-dev/pflow-eth/metamodel"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Service struct {
	NodeDb          *sql.DB
	Client          *ethclient.Client
	Logger          *log.Logger
	Router          *mux.Router
	applicationPage *template.Template
	apm             *newrelic.Application
	cancelFunc      context.CancelFunc
	addressHeights  map[common.Address]uint64
}

type Server interface {
	Serve(box *rice.Box)
}

func New() Server {
	s := &Service{
		Router:         mux.NewRouter(),
		addressHeights: map[common.Address]uint64{},
	}
	var err error
	s.NodeDb, err = sql.Open("postgres", config.DbConn)
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
	s.Logger.Printf("event %s %s\n", eventType, data)
}

func (s *Service) Serve(box *rice.Box) {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFunc = cancel

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		cancel()
	}()

	go s.runBlockIndexer(ctx, map[string]bool{
		"hardhat": true,
	})

	s.applicationRoutes(box)
	s.apiRoutes()
	err := http.ListenAndServe(config.Host+":"+config.Port, s.Router)
	if err != nil {
		panic(err)
	}
}

func (s *Service) apiRoutes() {
	s.Router.HandleFunc("/v0/ping", s.SyncStatsHandler)
	s.Router.HandleFunc("/v0/control", s.ControlPanelHandler)
	s.Router.HandleFunc("/v0/transactions", s.TransactionsHandler)
	s.Router.HandleFunc("/v0/model", s.ContextHandler)
	s.Router.HandleFunc("/v0/signal", s.SignalHandler)
	s.Router.HandleFunc("/v0/faucet", s.FaucetHandler)
	s.Router.HandleFunc("/v0/raw_tx", s.RawTransactionHandler)
}

func (s *Service) getHightestSequence() (highestSequence int64) {
	err := s.NodeDb.QueryRow("SELECT MAX(sequence) FROM transaction_logs_view").Scan(&highestSequence)
	if err != nil {
		log.Printf("Error finding highest sequence: %v", err)
		return -1
	}
	return highestSequence
}
func (s *Service) refreshDataView(schema string) bool {
	s.setSearchPath(schema)
	_, err := s.NodeDb.Exec("REFRESH MATERIALIZED VIEW node_sync_data_view")
	ok := err == nil
	return ok
}

func (s *Service) syncNodeWithBlockchain(schema string) bool {
	ok := s.refreshDataView(schema)
	s.Event("refresh_materialized_view", map[string]interface{}{
		"schema": schema,
		"ok":     ok,
	})

	if !ok {
		log.Println("Error refreshing materialized view node_sync_data_view")
	}
	return ok
}

func (s *Service) confirmSentTransactions(schema string) (ok bool) {
	var jsonString string
	s.setSearchPath(schema) // REVIEW: test this function
	err := s.NodeDb.QueryRow("SELECT confirm_transactions();").Scan(&jsonString)
	if err != nil {
		log.Printf("Error calling remove_confirmed_transactions: %v", err)
		return
	}

	// Unmarshal the JSON string into a map
	var result map[string]interface{}
	err = json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		log.Printf("Error unmarshaling JSON from confirm_transactions: %v", err)
		return false
	}

	s.Event("confirm_transactions", result)
	return true
}

func (s *Service) getNetwork(host string) string {
	if strings.HasPrefix(host, "localhost") || strings.HasPrefix(host, "127.0.0.1") {
		return "hardhat"
	}
	if strings.HasPrefix(host, "optimism-sepolia") {
		return "optimism_sepolia"
	}
	if strings.HasPrefix(host, "optimism") {
		return "optimism"
	}
	return "public"
}

func (s *Service) setSearchPathForRequest(r *http.Request) {
	host := r.Host
	schema := s.getNetwork(host)
	// fmt.Printf("Setting search path to %s => %s\n", host, schema)
	s.setSearchPath(schema)
}

func (s *Service) setSearchPath(schema string) {
	_, err := s.NodeDb.Exec("SET search_path TO " + schema)
	if err != nil {
		panic(err)
	}
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
