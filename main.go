package main

import (
	"context"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"log"
	"net/http"
	"source.golabs.io/engineering-platforms/lens/trace-app-golang/config"
	"source.golabs.io/engineering-platforms/lens/trace-app-golang/environment"
	"source.golabs.io/engineering-platforms/lens/trace-app-golang/handler"
)

func main() {
	cfg := config.Load()
	env := environment.Init(cfg)
	defer env.TraceExporter.Shutdown(context.Background())

	r := mux.NewRouter()

	r.HandleFunc("/v1/books", handler.CreateBook(env)).Methods("POST")
	r.HandleFunc("/v1/books/{title}", handler.ReadBook(env)).Methods("GET")
	r.Use(otelmux.Middleware(cfg.AppName))

	log.Printf("Listening on port 8080...")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Printf("Shutting down server")
	}
}
