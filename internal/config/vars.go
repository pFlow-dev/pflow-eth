package config

import (
	"fmt"
	"os"
)

// Global variables, set at build time with ldflags or dynamically via environment variables.
var (
	JsBuild         = "51c9bf03"                                                       // Update to match the path ./public/p/static/js/main.<JsBuild>.js
	CssBuild        = "16b9238e"                                                       // Update to match the path ./public/p/static/css/main.<CssBuild>.css
	Endpoint        = "http://hardhat:8545"                                            // Default endpoint, can be overridden by ENDPOINT env var
	DbConn          = "dbname=pflow user=pflow password=pflow sslmode=disable host=db" // Default DB connection string, can be overridden by DB_HOST env var
	Port            = "8083"                                                           // Default port, can be overridden by PORT env var
	Host            = "127.0.0.1"                                                      // Default host, can be overridden by HOST env var
	NewRelicLicense string                                                             // New Relic License Key for monitoring, set via ldflags
	NewRelicApp     string                                                             // New Relic Application Name, set via ldflags
	Address         string                                                             // Address for the faucet, set via ldflags
)

func init() {
	// Override DbConn if DB_HOST environment variable is set.
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		DbConn = fmt.Sprintf("dbname=pflow user=pflow password=pflow sslmode=disable host=%s", dbHost)
		fmt.Printf("DB connection string set to %s\n", dbHost)
	}

	// Override Endpoint if ENDPOINT environment variable is set.
	if endpoint := os.Getenv("ENDPOINT"); endpoint != "" {
		Endpoint = endpoint
	}

	// Override Port if PORT environment variable is set.
	if listenPort, portSet := os.LookupEnv("PORT"); portSet {
		Port = listenPort
	}

	// Override Host if HOST environment variable is set.
	if listenHost, hostSet := os.LookupEnv("HOST"); hostSet {
		Host = listenHost
	}
	Address = os.Getenv("ADDRESS")
	if Address == "" {
		Address = "0x5fbdb2315678afecb367f032d93f642f64180aa3"
	}
}
