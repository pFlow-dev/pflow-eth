package config

import (
	"os"
)


// Global variables, set at build time with ldflags or dynamically via environment variables.
var (
	JsBuild         = "267f1a3e"                                                       // Update to match the path ./public/p/static/js/main.<JsBuild>.js
	CssBuild        = "223157f6"                                                       // Update to match the path ./public/p/static/css/main.<CssBuild>.css
	Endpoint        = "http://hardhat:8545"                                            // Default endpoint, can be overridden by ENDPOINT env var
	Port            = "8083"                                                           // Default port, can be overridden by PORT env var
	Host            = "127.0.0.1"                                                      // Default host, can be overridden by HOST env var
	NewRelicLicense string                                                             // New Relic License Key for monitoring, set via ldflags
	NewRelicApp     string                                                             // New Relic Application Name, set via ldflags
	Address         string                                                             // Address for the faucet, set via ldflags
	SessionSalt     = "pflow"                                                          // Salt for session cookies
	Network         = "hardhat"                                                        // Network to connect to
)

func init() {
	// Override Endpoint if ENDPOINT environment variable is set.
	if endpoint := os.Getenv("ENDPOINT"); endpoint != "" {
		Endpoint = endpoint
		// fmt.Printf("Endpoint set to %s\n", endpoint)
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

	// Override Network if NETWORK environment variable is set.
	if network := os.Getenv("NETWORK"); network != "" {
		Network = network
	}

	if Address == "" {
		if Network == "hardhat" {
			Address = "0x5fbdb2315678afecb367f032d93f642f64180aa3" // hardhat default address
		}
		if Network == "sepolia_optimism" {
			Address = "0x9265fd7b41b3f96c3123319b713a5c5a761981f1" // sepolia_optimism
		}
	}

}
