package main

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"os"
)

var pipelinesDir = "./data/pipelines/"

// Pipeline type represents single pipeline information
type Pipeline struct {
	ID        string
	Name      string
	Builds    []int
	Schedule  []string
	LastBuild int
	Tasks     map[string]Task
}

// run a pipeline build process
func pipelineRun(id string, ch chan interface{}) {

	pipeline, err := pipelineLoad(id)
	if err != nil {
		if ch != nil {
			ch <- false
		}
		return
	}

	build := Build{
		PipelineID: pipeline.ID,
		tasks:      pipeline.Tasks,
		channel:    ch,
	}

	err = build.run(pipeline)
	if err != nil {
		log.Println("Pipeline Run Error:" + err.Error())
	}
}

// load pipeline json file
func pipelineLoad(id string) (Pipeline, error) {
	var pipeline Pipeline

	// load the json file
	file, err := ioutil.ReadFile(pipelinesDir + id + "/pipeline.json")
	if err != nil {
		return pipeline, err
	}

	// decode the json into pipeline struct
	err = json.Unmarshal(file, &pipeline)
	if err != nil {
		return pipeline, err
	}

	return pipeline, nil
}

// new pipeline
func pipelineNew() (string, error) {

	// generate new id
	id := uuid.NewV4().String()

	// make the pipeline dir
	err := os.Mkdir(pipelinesDir+id, 0700)
	if err != nil {
		return "", err
	}

	// make the pipeline builds dir
	err = os.Mkdir(pipelinesDir+id+"/builds", 0700)
	if err != nil {
		return "", err
	}

	// make the pipeline work dir
	// TODO: make workdir for each build to avoid collision if two builds run in the same time
	err = os.Mkdir(pipelinesDir+id+"/workdir", 0700)
	if err != nil {
		return "", err
	}

	return id, nil
}

// save pipeline json to file
func pipelineSave(pipeline Pipeline) error {

	// reset the scheduler for this pipeline
	schedulerSet(pipeline.ID, pipeline.Schedule)

	// encode the pipeline data to json
	pipelineJSON, err := json.Marshal(pipeline)
	if err != nil {
		return err
	}

	// save the pipeline file
	err = ioutil.WriteFile(pipelinesDir+pipeline.ID+"/pipeline.json", pipelineJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}

// delete pipeline data directory
func pipelineDelete(id string) error {

	err := os.RemoveAll(pipelinesDir + id)
	if err != nil {
		return err
	}

	return nil
}

// All pipelines
func loadAllPipelines() []Pipeline {
	var pipelines []Pipeline

	// load all the units folders
	files, _ := ioutil.ReadDir(pipelinesDir)

	for _, f := range files {
		// load the pipeline and add it to the units list
		pipeline, err := pipelineLoad(f.Name())
		if err != nil {
			continue
		}

		pipelines = append(pipelines, pipeline)
	}

	return pipelines
}
