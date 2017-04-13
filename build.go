package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

// Build type represents a single build information
// and the required methods for running a build
type Build struct {
	PipelineID string
	Success    bool
	StartAt    time.Time
	EndAt      time.Time
	Logs       []BuildLog
	tasks      map[string]Task
	channel    chan interface{}
}

// BuildLog type represents a single build's log
type BuildLog struct {
	TaskID   string
	TaskName string
	Command  string
	Console  string
	Level    int
}

// start the build process
func (build *Build) run(pipeline Pipeline) error {
	// increase the pipeline last build number and save the pipeline
	// to avid build id collision if two builds run in the same time
	pipeline.LastBuild++
	pipelineSave(pipeline)

	build.StartAt = time.Now()

	// run the first task "start"
	err := build.task("start", nil, 0)
	if err != nil {
		log.Println("Build Error:", err.Error())
	}

	build.EndAt = time.Now()
	build.Success = err == nil

	// save the build logs
	err = build.save(pipeline)
	if err != nil {
		if build.channel != nil {
			build.channel <- false
		}

		return fmt.Errorf("Save Build Error: " + err.Error())
	}

	// send true to the message to close the websocket
	if build.channel != nil {
		build.channel <- true
	}

	return nil
}

// run a task
func (build *Build) task(id string, inputsValues map[string]string, level int) error {
	// check is the task exist first
	task, ok := build.tasks[id]
	if !ok {
		return fmt.Errorf("Task [%s] not exist", id)
	}

	var err error
	var console string
	var execOutputs map[string]string

	// if the task have command execute it
	if task.Command != "" {
		execOutputs, console, err = task.exec(build.PipelineID, inputsValues)
	}

	build.log(BuildLog{TaskID: id, TaskName: task.Name, Command: task.Command, Console: console, Level: level})
	if err != nil {
		return err
	}

	err = build.nextTasks(task.Outputs, execOutputs, level)
	if err != nil {
		return err
	}

	return nil
}

// save pipeline json to file
func (build *Build) save(pipeline Pipeline) error {

	// encode the build details to json
	buildJSON, err := json.Marshal(build)
	if err != nil {
		return err
	}

	// save the build log file inside the pipeline build folder
	path := fmt.Sprintf("%s%s/builds/%d.json", pipelinesDir, pipeline.ID, pipeline.LastBuild)
	err = ioutil.WriteFile(path, buildJSON, 0644)
	if err != nil {
		return err
	}

	// add the build id to the pipeline's builds and save it
	pipeline.Builds = append(pipeline.Builds, pipeline.LastBuild)
	pipelineSave(pipeline)

	// cleanup working directory
	os.RemoveAll(pipelinesDir + pipeline.ID + "/workdir")
	os.Mkdir(pipelinesDir+pipeline.ID+"/workdir", 0700)

	return nil
}

// collect the next tasks and its inputs values
func (build *Build) nextTasks(outputs TaskOutputs, outputsValues map[string]string, level int) error {
	// prepare the outputs and categorize them by the destination task
	nextTasks := map[string]map[string]string{}
	for outputID, output := range outputs {
		for _, destination := range output.Destination {
			// initial the destination map if it not exist
			if _, ok := nextTasks[destination.Task]; !ok {
				nextTasks[destination.Task] = make(map[string]string)
			}

			// add the output to the destination map
			nextTasks[destination.Task][destination.Input] = outputsValues[outputID]
		}
	}

	// run the next tasks
	level++
	for destinationTask, destinationInputs := range nextTasks {
		err := build.task(destinationTask, destinationInputs, level)
		if err != nil {
			return err
		}
	}

	return nil
}

// log build console output
func (build *Build) log(l BuildLog) {
	if build.channel != nil {
		build.channel <- l
	}

	build.Logs = append(build.Logs, l)
}

// delete a build
func buildDelete(pipelineID string, buildID string) error {

	// load the pipeline
	pipeline, err := pipelineLoad(pipelineID)
	if err != nil {
		return err
	}

	// delete the build id from the pipeline
	pipeline.Builds = removeBuildID(pipeline.Builds, buildID)

	// update the pipeline file
	err = pipelineSave(pipeline)
	if err != nil {
		return err
	}

	// delete the build file
	err = os.Remove(pipelinesDir + pipelineID + "/builds/" + buildID + ".json")
	if err != nil {
		return err
	}

	return nil
}

// remove a build id from the pipeline builds array
func removeBuildID(list []int, item string) []int {

	// convert id to int
	id, _ := strconv.Atoi(item)

	// build the new list
	var newList []int
	for _, i := range list {
		if i == id {
			continue
		}
		newList = append(newList, i)
	}

	return newList
}
