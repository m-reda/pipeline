package core

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

func singlePipelineHandler(w http.ResponseWriter, r *http.Request)  {

	file, _ := ioutil.ReadFile(pipelinesDir + "/1/pipeline.json")
	w.Write(file)
}


func buildPipelineHandler(w http.ResponseWriter, _ *http.Request)  {
	fail := pipelineRun("1")

	var msg string
	if fail == nil {
		msg = "[Build Success]"
	} else {
		msg = "[Build Fail]" + fail.Error()
	}

	w.Write([]byte(msg));
}

func storePipelineHandler(w http.ResponseWriter, r *http.Request)  {
	pipelineJson := r.FormValue("pipeline")

	if pipelineJson == "" {
		w.Write([]byte(`{"success": false, "message": "pipeline field empty"}`));
		return
	}

	if err := pipelineSave("1", []byte(pipelineJson)); err != nil {
		w.Write([]byte(`{"success": false, "message": "`+ err.Error() +`"}`));
		return
	}

	w.Write([]byte(`{"success": true}`));
}


func allUnitsHandler(w http.ResponseWriter, r *http.Request)  {

	var units []Unit

	files, _ := ioutil.ReadDir(unitsDir)
    for _, f := range files {

		if unit, err := unitLoad(f.Name()); err == nil {
			units = append(units, unit)
		}
    }

	unitsJson, _ := json.Marshal(units)

	w.Write(unitsJson)
}
