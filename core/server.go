package core

import (
	"github.com/gorilla/mux"
	"os"
	"net/http"
	"time"
	"log"
)

func ServerRun() {


	fail := pipelineRun("1")

	if fail == nil {
		println("[Build Success]")
	} else {
		println("[Build Fail]", fail.Error())
	}

	return
	r := mux.NewRouter()

	// get the server port from env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

    srv := &http.Server{
        Handler:      r,
        Addr:         ":" + port,
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }

    log.Fatal(srv.ListenAndServe())
}
