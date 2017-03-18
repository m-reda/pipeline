package core

import (
	"github.com/gorilla/mux"
	"os"
	"net/http"
	"time"
	"log"
)

func ServerRun() {

	r := mux.NewRouter()
	r.HandleFunc("/api/pipelines/{key}", singlePipelineHandler).Methods("GET")
	r.HandleFunc("/api/pipelines/{key}", storePipelineHandler).Methods("PUT")
	r.HandleFunc("/api/pipelines/{key}/build", buildPipelineHandler)

    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	// get the server port from env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}


    srv := &http.Server{
        Handler:      r,
        Addr:         ":" + port,
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }

    log.Fatal(srv.ListenAndServe())
}
