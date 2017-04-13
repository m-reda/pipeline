package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	schedulerStart()
	log.Fatal(serverStart())
}

func handlers() *mux.Router {
	router := mux.NewRouter()

	// pipelines routes
	router.HandleFunc("/api/pipelines", allPipelinesHandler).Methods("GET")
	router.HandleFunc("/api/pipelines", storePipelineHandler).Methods("POST", "PUT")
	router.HandleFunc("/api/pipelines/{id}", singlePipelineHandler).Methods("GET")
	router.HandleFunc("/api/pipelines/{id}", deletePipelineHandler).Methods("DELETE")

	// builds routes
	router.HandleFunc("/api/pipelines/{id}/build", buildPipelineHandler)
	router.HandleFunc("/api/pipelines/{id}/build/{bid}", buildDetailsHandler).Methods("GET")
	router.HandleFunc("/api/pipelines/{id}/build/{bid}", deleteBuildHandler).Methods("DELETE")

	// units routes
	router.HandleFunc("/api/units", allUnitsHandler).Methods("GET")

	// static files server
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./data/ui")))

	return router
}

func serverStart() error {
	// server port
	port := "8080"
	if env := os.Getenv("PORT"); env != "" {
		port = env
	}

	fmt.Printf("Start pipeline server on :%s.\n", port)

	// setting the server
	server := &http.Server{
		Handler:      handlers(),
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return server.ListenAndServe()
}
