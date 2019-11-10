package main

import (
	"context"
	"github.com/ottotech/riskmanagement/pkg/adding"
	"github.com/ottotech/riskmanagement/pkg/deleting"
	"github.com/ottotech/riskmanagement/pkg/http/rest"
	"github.com/ottotech/riskmanagement/pkg/listing"
	"github.com/ottotech/riskmanagement/pkg/storage/memory"
	"github.com/ottotech/riskmanagement/pkg/updating"
	"log"
	"net/http"
)

const (
	Memory int = 1
)

func shutDownHandler(signal chan bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
			return
		}
		signal <- true
		return
	})
}

func main() {
	// set up storage
	storageType := Memory // this could be a flag; hardcoded here for simplicity

	var adder adding.Service
	var lister listing.Service
	var deleter deleting.Service
	var updater updating.Service

	switch storageType {
	case Memory:
		s := new(memory.Storage)
		adder = adding.NewService(s)
		lister = listing.NewService(s)
		deleter = deleting.NewService(s)
		updater = updating.NewService(s)
	}
	idleConnsClosed := make(chan struct{})
	shutDownSignal := make(chan bool, 1)
	app := new(rest.App)
	mux := http.NewServeMux()
	mux.Handle("/", app.List.Handler(lister)) // home
	mux.Handle("/add", app.Add.Handler(adder, lister))
	mux.Handle("/get/", app.Get.Handler(lister))
	mux.Handle("/add-risks", app.AddRisk.Handler(adder, lister, updater))
	mux.Handle("/delete-risks", app.DeleteRisk.Handler(deleter, lister))
	mux.Handle("/delete-risk-matrix", app.DeleteRiskMatrix.Handler(deleter, lister))
	mux.Handle("/media/", app.Media.Handler())
	mux.Handle("/shutdown", shutDownHandler(shutDownSignal))
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	go func() {
		<-shutDownSignal
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	<-idleConnsClosed
}
