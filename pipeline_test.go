package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strconv"
	"testing"
)

func TestPipelineRun(t *testing.T) {

	pipelineRun(mainPipeline.ID, nil)

	mainPipeline.LastBuild++
	buildID := strconv.Itoa(mainPipeline.LastBuild)

	require.True(t, isBuildFileExist(mainPipeline.ID, buildID))
}

func TestPipelineRunNotExist(t *testing.T) {

	ch := make(chan interface{})
	go pipelineRun("-", ch)

	done := <-ch

	switch done.(type) {
	case bool:
		require.False(t, done.(bool))
	default:
		t.Fatal("")
	}
}

func TestPipelineLoad(t *testing.T) {
	_, err := pipelineLoad(mainPipeline.ID)
	require.NoError(t, err)

	_, err = pipelineLoad("-")
	require.Error(t, err)
}

func TestPipelineNewAndDelete(t *testing.T) {
	for i := 0; i < 5; i++ {
		id, err := pipelineNew()

		require.NoError(t, err)
		require.NotEmpty(t, id)

		err = pipelineDelete(id)
		require.NoError(t, err)
	}
}

func TestPipelineSave(t *testing.T) {

	// create new pipeline
	id, err := pipelineNew()
	require.NoError(t, err)

	pipeline := Pipeline{ID: id}

	// save the pipeline json file
	err = pipelineSave(pipeline)
	assert.NoError(t, err)

	// check that the pipeline file exist
	if _, err := os.Stat(pipelinesDir + id + "/pipeline.json"); err != nil {
		t.Errorf("Pipeline %s not exist", id)
	}

	// delete the pipeline
	err = pipelineDelete(id)
	assert.NoError(t, err)
}

func TestLoadAllPipelines(t *testing.T) {

	// check if the main pipeline counts
	pipelinesCount := len(loadAllPipelines())
	if pipelinesCount < 1 {
		t.Error("Shoud be one or more.")
	}
}
