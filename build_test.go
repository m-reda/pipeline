package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"strconv"
	"testing"
)

var mainBuild = Build{
	PipelineID: mainPipeline.ID,
	tasks:      mainPipeline.Tasks,
}

func TestBuildRun(t *testing.T) {
	err := mainBuild.run(mainPipeline)
	require.NoError(t, err)
}

func TestBuildTask(t *testing.T) {
	err := mainBuild.task("task2", nil, 0)
	require.NoError(t, err)

	err = mainBuild.task("----", nil, 0)
	require.Error(t, err)
}

func TestBuildTaskError(t *testing.T) {
	err := mainBuild.task("task3", nil, 0)
	require.Error(t, err)
}

func TestBuildSave(t *testing.T) {
	err := mainBuild.save(mainPipeline)
	require.NoError(t, err)
}

func TestBuildLog(t *testing.T) {
	logsCount := len(mainBuild.Logs)
	mainBuild.log(BuildLog{
		TaskID:   "id",
		TaskName: "name",
		Command:  "command",
		Console:  "console",
		Level:    0,
	})

	if len(mainBuild.Logs) <= logsCount {
		t.Errorf("Build log error expected more than %d logs, got %d", logsCount, len(mainBuild.Logs))
	}
}

func TestBuildDelete(t *testing.T) {
	// create new build
	err := mainBuild.run(mainPipeline)
	require.NoError(t, err)

	mainPipeline.LastBuild++
	bid := strconv.Itoa(mainPipeline.LastBuild)

	// check that the build file exist before delete it
	require.True(t, isBuildFileExist(mainPipeline.ID, bid))

	// delete the build
	err = buildDelete(mainPipeline.ID, bid)
	require.NoError(t, err)

	// check that the build file not exist after delete it
	require.False(t, isBuildFileExist(mainPipeline.ID, bid))
}

func TestBuildRemoveFrom(t *testing.T) {
	list := []int{1, 2, 3, 4, 5}
	newList := removeBuildID(list, "2")
	require.Equal(t, []int{1, 3, 4, 5}, newList)
}

func isBuildFileExist(pipelineID string, buildID string) bool {
	path := fmt.Sprintf("%s%s/builds/%s.json", pipelinesDir, pipelineID, buildID)
	if _, err := os.Stat(path); err != nil {
		return false
	}

	return true
}
