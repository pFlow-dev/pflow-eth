package config

import (
	"fmt"
	"os"
)

var ( // set at buildtime w/ ldflags
	JsBuild  = `d0b89f71` // update to match ./public/p/static/js/main.<JsBuild>.js
	CssBuild = `459ab37a` // update to match ./public/p/static/css/main.<CssBuild>.css

	Endpoint = `http://harhat:8545`                                             // assume docker usage
	DbConn   = "dbname=pflow user=pflow password=pflow sslmode=disable host=db" // docker
)

func init() {
	if DB_HOST := os.Getenv("DB_HOST"); DB_HOST != "" {
		DbConn = fmt.Sprintf("dbname=pflow user=pflow password=pflow sslmode=disable host=%s", DB_HOST)
	}
	if ENDPOINT := os.Getenv("ENDPOINT"); ENDPOINT != "" {
		Endpoint = ENDPOINT
	}
}
