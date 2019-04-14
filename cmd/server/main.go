package main

import (
	"github.com/ottotech/riskmanagement/pkg/adding"
	"github.com/ottotech/riskmanagement/pkg/http/rest"
	"github.com/ottotech/riskmanagement/pkg/listing"
	"github.com/ottotech/riskmanagement/pkg/storage/memory"
	"log"
	"net/http"
)

const (
	Memory int = 1
)

func main() {
	// set up storage
	storageType := Memory // this could be a flag; hardcoded here for simplicity

	var adder adding.Service
	var lister listing.Service

	switch storageType {
	case Memory:
		s := new(memory.Storage)
		adder = adding.NewService(s)
		lister = listing.NewService(s)
	// more data stores can be supported
	}

	app := new(rest.App)
	mux := http.NewServeMux()
	mux.Handle("/", app.HomeHandler.Handler())
	mux.Handle("/add", app.AddHandler.Handler(adder))
	mux.Handle("/list", app.ListHandler.Handler(lister))
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}
