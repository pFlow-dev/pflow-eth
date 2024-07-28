package main

import (
	rice "github.com/GeertJohan/go.rice"
	"github.com/pflow-dev/pflow-eth/internal/service"
)

func main() {
	s := service.New()
	s.Serve(rice.MustFindBox("./public"))
}
