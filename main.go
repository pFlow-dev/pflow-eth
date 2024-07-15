package main

import (
	rice "github.com/GeertJohan/go.rice"
	"github.com/pflow-dev/pflow-eth/internal/service"
	"os"
)

var (
	options = service.Options{
		Host:            "127.0.0.1",
		Port:            "8083",
		DbPath:          "/tmp/pflow.db",
		NewRelicApp:     "pflow.dev",
		NewRelicLicense: os.Getenv("NEW_RELIC_LICENSE"),
	}
)

func main() {
	dbPath, pathSet := os.LookupEnv("DB_PATH")
	if pathSet {
		options.DbPath = dbPath
	}
	baseUrl, urlSet := os.LookupEnv("URL_BASE")
	if urlSet {
		options.Url = baseUrl
	}
	listenPort, portSet := os.LookupEnv("PORT")
	if portSet {
		options.Port = listenPort
	}
	listenHost, hostSet := os.LookupEnv("HOST")
	if hostSet {
		options.Host = listenHost
	}

	s := service.New(options)
	s.Serve(rice.MustFindBox("./public"))
}
