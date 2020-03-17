package main

import (
	"context"
	"github.com/ottotech/riskmanagement/pkg/adding"
	"github.com/ottotech/riskmanagement/pkg/config"
	"github.com/ottotech/riskmanagement/pkg/deleting"
	"github.com/ottotech/riskmanagement/pkg/http/rest"
	"github.com/ottotech/riskmanagement/pkg/listing"
	mw "github.com/ottotech/riskmanagement/pkg/middlewares"
	"github.com/ottotech/riskmanagement/pkg/storage/json"
	"github.com/ottotech/riskmanagement/pkg/storage/memory"
	"github.com/ottotech/riskmanagement/pkg/updating"
	"log"
	"net/http"
	"os"
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
	// configure logger
	logFile, err := os.OpenFile("error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	config.Logger = log.New(logFile, "Error Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
	defer func() {
		err := logFile.Close()
		if err != nil {
			config.Logger.Println(err)
		}
	}()
	// set up storage
	var adder adding.Service
	var lister listing.Service
	var deleter deleting.Service
	var updater updating.Service

	switch config.StorageType {
	case config.MEMORY:
		s := new(memory.Storage)
		adder = adding.NewService(s)
		lister = listing.NewService(s)
		deleter = deleting.NewService(s)
		updater = updating.NewService(s)
	case config.JSON:
		s, err := json.NewStorage()
		if err != nil {
			config.Logger.Println(err)
			log.Fatal(err)
		}
		adder = adding.NewService(s)
		lister = listing.NewService(s)
		deleter = deleting.NewService(s)
		updater = updating.NewService(s)
	default:
		config.Logger.Fatalln("No valid data store specified.")
	}
	idleConnsClosed := make(chan struct{})
	shutDownSignal := make(chan bool, 1)
	app := new(rest.App)
	mux := http.NewServeMux()
	mux.Handle("/", mw.Chain(app.ListMatrix.Handler(lister), mw.MediaPathRequired(lister)))
	mux.Handle("/add", mw.Chain(app.AddMatrix.Handler(adder, lister), mw.MediaPathRequired(lister)))
	mux.Handle("/get/", mw.Chain(app.GetMatrix.Handler(lister), mw.MediaPathRequired(lister)))
	mux.Handle("/add-risks", mw.Chain(app.AddRisk.Handler(adder, lister, updater), mw.MediaPathRequired(lister)))
	mux.Handle("/delete-risks", mw.Chain(app.DeleteRisk.Handler(deleter, lister), mw.MediaPathRequired(lister)))
	mux.Handle("/delete-risk-matrix", mw.Chain(app.DeleteRiskMatrix.Handler(deleter, lister), mw.MediaPathRequired(lister)))
	mux.Handle("/set-media-path", app.AddMediaPath.Handler(adder))
	mux.Handle("/media/", app.Media.Handler())
	mux.Handle("/shutdown", shutDownHandler(shutDownSignal))
	mux.Handle("/favicon.ico", http.NotFoundHandler())
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	go func() {
		<-shutDownSignal
		if err := server.Shutdown(context.Background()); err != nil {
			config.Logger.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		config.Logger.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	<-idleConnsClosed
}
