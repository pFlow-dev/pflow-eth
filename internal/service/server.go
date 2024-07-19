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
	"github.com/pflow-dev/pflow-eth/internal/config"
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
}

type Server interface {
	Serve(box *rice.Box)
}

func New() Server {
	s := &Service{
		Router: mux.NewRouter(),
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

func (s Service) apiRoutes() {
	s.Router.HandleFunc("/v0/ping", s.SyncStatsHandler)
	s.Router.HandleFunc("/v0/control", s.ControlPanelHandler)
	s.Router.HandleFunc("/v0/transactions", s.TransactionsHandler)
	s.Router.HandleFunc("/v0/model", s.ContextHandler)
	s.Router.HandleFunc("/v0/signal", s.SignalHandler)
	s.Router.HandleFunc("/v0/faucet", s.FaucetHandler)
}

func (s *Service) getHightestSequence() (highestSequence int64) {
	err := s.NodeDb.QueryRow("SELECT MAX(sequence) FROM transaction_states").Scan(&highestSequence)
	if err != nil {
		log.Printf("Error finding highest sequence: %v", err)
		return -1
	}
	return highestSequence
}

func (s *Service) insertSequence(nextSeq int64) (int64, error) {
	sqlStatement := `
    INSERT INTO transaction_states (sequence, transaction_hash, state)
    SELECT sequence, transaction_hash, NULL
    FROM (
        SELECT sequence, transaction_hash
        FROM transaction_log_view
        WHERE sequence = {{}} 
    ) AS next_sequence
    `

	result, err := s.NodeDb.Exec(sqlStatement, nextSeq)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (s *Service) syncNodeWithBlockchain(schema string) bool {
	s.setSearchPath(schema)
	_, err := s.NodeDb.Exec("REFRESH MATERIALIZED VIEW node_sync_data_view")
	if err != nil {
		log.Printf("Error refreshing materialized view node_sync_data_view: %v", err)
		return false
	} else {
		log.Println("Materialized view node_sync_data_view refreshed successfully.")
		return true
	}
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

func (s *Service) importNextTransaction(schema string) (ok bool) {
	s.setSearchPath(schema)
	highestSequence := s.getHightestSequence()
	if highestSequence < 0 {
		return false
	}
	nextSeq := highestSequence + 1
	rows, err := s.insertSequence(nextSeq)
	if err != nil {
		log.Printf("Error inserting sequence %d: %v", nextSeq, err)
		return false
	}
	if rows > 0 {
		s.Event("sequence", map[string]interface{}{
			"sequence": nextSeq,
		})
		log.Printf("Inserted sequence %d", nextSeq)
		return true
	}
	return false
}

func (s *Service) computeTransactionStates(schema string) (ok bool) {
	s.setSearchPath(schema)
	// REVIEW: test this function
	sqlStatement := `SELECT sequence, transaction_hash FROM transaction_states WHERE state IS NULL order by sequence asc`

	// Execute the query
	rows, err := s.NodeDb.Query(sqlStatement)
	if err != nil {
		log.Printf("Error querying null state transactions: %v", err)
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var sequence int64
		var transactionHash string

		err := rows.Scan(&sequence, &transactionHash)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return false
		}
		log.Printf("Found null state transaction: sequence %d, hash %s", sequence, transactionHash)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return false
	}

	return true
}

func (s Service) getNetwork(host string) string {
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

func (s Service) setSearchPathForRequest(r *http.Request) {
	host := r.Host
	schema := s.getNetwork(host)
	// fmt.Printf("Setting search path to %s => %s\n", host, schema)
	s.setSearchPath(schema)
}

func (s Service) setSearchPath(schema string) {
	_, err := s.NodeDb.Exec("SET search_path TO " + schema)
	if err != nil {
		panic(err)
	}
}
