package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

var inputs = map[string]string{
	"input1": "value 1",
	"input2": "value 2",
}

func TestTaskExec(t *testing.T) {

	task := Task{Command: "echo 'output1:{input1}\noutput2:{input2}'"}
	execOutputs, console, err := task.exec(mainPipeline.ID, inputs)

	require.NoError(t, err)
	require.Equal(t, map[string]string{"output1": "value 1", "output2": "value 2"}, execOutputs)
	require.Equal(t, "output1:value 1\noutput2:value 2\n", console)
}

func TestTaskPrepareCommand(t *testing.T) {
	task := Task{
		Command: "echo 'output1:{input1}\noutput2:{input2}'",
		Overwrite: map[string]string{
			"input1": "overwrited",
		},
	}

	task.prepareCommand("1", inputs)

	require.Equal(t, "cd ./data/pipelines/1/workdir && echo 'output1:overwrited\noutput2:value 2'", task.Command)
}

func TestTaskPrepareCommandWithSetting(t *testing.T) {
	task := Task{
		Command: "echo 'output1:{input1}\noutput2:{input2}'",
		Setting: map[string]map[string]string{
			"abc": {"Value": "101"},
			"def": {"Value": "102"},
			"xyz": {"Value": ""},
		},
	}

	task.prepareCommand("1", nil)

	// check the setting flags separately, because map's order not guaranteed
	require.Contains(t, task.Command, "-abc '101'")
	require.Contains(t, task.Command, "-def '102'")
}

func TestTaskCommandOutputs(t *testing.T) {
	outputs := commandOutputs("output1:value 1\noutput2:value 2\noutput3")

	require.Equal(t, map[string]string{"output1": "value 1", "output2": "value 2"}, outputs)
}
