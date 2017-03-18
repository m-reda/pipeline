package core

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
)

type Pipeline struct {
	ID string
	Name string
	Tasks map[string]Task
}

// run a pipeline
func pipelineRun(id string) error {

	pipeline, err := pipelineLoad(id)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Pipeline >>", pipeline.Name)

	return taskRun("start", nil, pipeline.Tasks)
}

// load pipeline json file
func pipelineLoad(id string) (Pipeline, error) {
	var pipeline Pipeline

	file, err := ioutil.ReadFile("./.data/pipelines/" + id + "/pipeline.json")

	if err != nil {
		return pipeline, err
	}

	err = json.Unmarshal(file, &pipeline)
	if err != nil {
		return pipeline, err
	}

	return pipeline, nil
}

// save pipeline json to file
func pipelineSave(id string, pipelineJson []byte) error{
	var pipeline Pipeline

	err := json.Unmarshal(pipelineJson, &pipeline)

	if err != nil {
		return err
	}

	pipelineJson, err = json.Marshal(pipeline)

	if err != nil {
		return err
	}

    err = ioutil.WriteFile("./.data/pipelines/" + id + "/pipeline.json", pipelineJson, 0644)

	if err != nil {
		return err
	}

	return nil
}

// TODO: execute command timeout
