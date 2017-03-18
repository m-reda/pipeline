package core

import (
	"net/http"
	"io/ioutil"
)

func singlePipelineHandler(w http.ResponseWriter, r *http.Request)  {

	file, _ := ioutil.ReadFile("./.data/pipelines/1/pipeline.json")
	w.Write(file)
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
