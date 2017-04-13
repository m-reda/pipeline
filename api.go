package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
)

// request all pipelines list
func allPipelinesHandler(w http.ResponseWriter, _ *http.Request) {

	// encode the pipelines list
	unitsJSON, err := json.Marshal(loadAllPipelines())
	if err != nil {
		w.Write([]byte(`{"success": false, "message": "Pipeline not exist"}`))
		return
	}

	w.Write(unitsJSON)
}

// request a pipeline details
func singlePipelineHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	file, err := ioutil.ReadFile(pipelinesDir + vars["id"] + "/pipeline.json")
	if err != nil {
		w.Write([]byte(`{"success": false, "message": "Pipeline not exist"}`))
		return
	}

	w.Write(file)
}

// save a pipeline details
func storePipelineHandler(w http.ResponseWriter, r *http.Request) {
	pipelineJSON := r.FormValue("pipeline")

	isNew := r.Method == "POST"

	// check the pipeline not empty
	if pipelineJSON == "" {
		w.Write([]byte(`{"success": false, "message": "pipeline field empty"}`))
		return
	}

	var pipeline Pipeline
	var err error

	if isNew {
		pipeline.ID, err = pipelineNew()
	}

	if err == nil {
		// check that the json is valid
		err = json.Unmarshal([]byte(pipelineJSON), &pipeline)
	}

	if err == nil {
		err = pipelineSave(pipeline)
	}

	if err != nil {
		if isNew && pipeline.ID != "" {
			pipelineDelete(pipeline.ID)
		}

		w.Write([]byte(`{"success": false, "message": "` + err.Error() + `"}`))
		return
	}

	w.Write([]byte(`{"success": true, "id": "` + pipeline.ID + `"}`))
}

// delete a pipeline
func deletePipelineHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := pipelineDelete(vars["id"])
	if err != nil {
		w.Write([]byte(`{"success": false, "message": ` + err.Error() + `}`))
		return
	}

	w.Write([]byte(`{"success": true}`))
}

// start building a pipeline
func buildPipelineHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	isWebSocket := websocket.IsWebSocketUpgrade(r)

	var ch chan interface{}

	// initial only if request is websocket
	if isWebSocket {
		ch = make(chan interface{})
	}

	go pipelineRun(vars["id"], ch)

	if isWebSocket {
		websocketServer(w, r, ch)
	} else {
		w.Write([]byte(`{"success": true}`))
	}
}

// request a build details
func buildDetailsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	file, err := ioutil.ReadFile(pipelinesDir + vars["id"] + "/builds/" + vars["bid"] + ".json")
	if err != nil {
		w.Write([]byte(`{"success": false, "message": "` + err.Error() + `"}`))
		return
	}

	w.Write(file)
}

// delete a build
func deleteBuildHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := buildDelete(vars["id"], vars["bid"])
	if err != nil {
		w.Write([]byte(`{"success": false, "message": "` + err.Error() + `"}`))
		return
	}

	w.Write([]byte(`{"success": true}`))
}

// request all units list
func allUnitsHandler(w http.ResponseWriter, _ *http.Request) {

	// encode the units list
	unitsJSON, _ := json.Marshal(loadAllUnits())
	w.Write(unitsJSON)
}

func websocketServer(w http.ResponseWriter, r *http.Request, ch chan interface{}) {

	upgrader := new(websocket.Upgrader)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	defer conn.Close()

	for buildLog := range ch {
		switch buildLog.(type) {
		// close channel if bool
		case bool:
			msg := "fail"
			if buildLog.(bool) {
				msg = "done"
			}

			conn.WriteMessage(websocket.TextMessage, []byte(msg))
			return

		// write log to channel
		default:
			err = conn.WriteJSON(buildLog)
			if err != nil {
				return
			}
		}
	}
}
