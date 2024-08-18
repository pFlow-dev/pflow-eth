package config

import (
	"os"
)

// FIXME: Set Webui Version: main.6a6e9436.js main.88d73541.css

var (
	JsBuild         = "6a6e9436"                     // Update to match the path ./public/p/static/js/main.<JsBuild>.js
	CssBuild        = "88d73541"                     // Update to match the path ./public/p/static/css/main.<CssBuild>.css
	Endpoint        = "http://hardhat:8545"          // Default endpoint, can be overridden by ENDPOINT env var
	Port            = "8083"                         // Default port, can be overridden by PORT env var
	Host            = "127.0.0.1"                    // Default host, can be overridden by HOST env var
	NewRelicApp     = "explore.pflow.xyz"            // New Relic Application Name, set via ldflags
	NewRelicLicense = os.Getenv("NEW_RELIC_LICENSE") // New Relic License, set via environment variable
	SessionSalt     = "pflow"                        // Salt for session cookies
	Network         = "hardhat"                      // Network to connect to

	SubdomainToNetwork = map[string]string{
		"sepolia-optimism": "sepolia_optimism",
		"optimism":         "optimism",
	}
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

	// Override Network if NETWORK environment variable is set.
	if network := os.Getenv("NETWORK"); network != "" {
		Network = network
	}

}
