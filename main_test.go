package main

import (
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"os"
	"testing"
)

var server = httptest.NewServer(handlers())

var mainPipeline = Pipeline{
	Name:     "Main Pipeline",
	Schedule: []string{"0 0 * * *"},
	Tasks: map[string]Task{
		"start": {ID: "start", Name: "Start", Outputs: TaskOutputs{
			"run": {
				Destination: []struct {
					Task  string
					Input string
				}{
					{Task: "task2"},
				},
			},
		}},
		"task2": {ID: "task2", Name: "Task 2"},
		"task3": {ID: "task3", Name: "Task 3", Command: "not_program"},
	},
}

func TestMain(m *testing.M) {
	var err error

	// create new pipeline
	mainPipeline.ID, err = pipelineNew()
	if err != nil {
		panic(err.Error())
	}

	// save the pipeline json file
	err = pipelineSave(mainPipeline)
	if err != nil {
		panic(err.Error())
	}

	// run the test cases
	exit := m.Run()

	// delete the main pipeline
	err = pipelineDelete(mainPipeline.ID)
	if err != nil {
		panic(err.Error())
	}

	os.Exit(exit)
}

func TestServerStart(t *testing.T) {
	os.Setenv("PORT", "9999999")
	err := serverStart()

	require.Error(t, err)
}
