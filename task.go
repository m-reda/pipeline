package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// Task type represents a single pipeline's task
// with the required properties for the UI like position
type Task struct {
	ID        string
	Name      string
	Unit      string
	X         int
	Y         int
	Command   string
	Setting   map[string]map[string]string
	Overwrite map[string]string
	Inputs    map[string]string
	Outputs   TaskOutputs
}

// TaskOutputs type represents a map of the task outputs
// and each output destinations
type TaskOutputs map[string]struct {
	Name        string
	Destination []struct {
		Task  string
		Input string
	}
}

func (task *Task) exec(pipelineID string, inputs map[string]string) (map[string]string, string, error) {
	task.prepareCommand(pipelineID, inputs)

	// execute the task command
	consoleBytes, err := exec.Command("/bin/sh", "-c", task.Command).Output()
	console := string(consoleBytes)

	if err != nil {
		// append the error to the console output
		console = "Unit command error: " + err.Error() + ". " + console
		return nil, console, err
	}

	execOutputs := commandOutputs(console)

	return execOutputs, console, nil
}

func (task *Task) prepareCommand(pipelineID string, inputs map[string]string) {
	// set the inputs overwrites
	for key, value := range task.Overwrite {
		if value != "" {
			inputs[key] = value
		}
	}

	// prepare the inputs
	for inputID, inputValue := range inputs {
		task.Command = strings.Replace(task.Command, "{"+inputID+"}", inputValue, -1)
	}

	task.Command += task.settingFlags()

	// change the working directory to the pipeline workdir
	task.Command = fmt.Sprintf("cd %s%s/workdir && %s", pipelinesDir, pipelineID, task.Command)
}

func (task *Task) settingFlags() string {

	flags := ""

	// add setting to command
	for flagName, flagSetting := range task.Setting {
		if flagSetting["Value"] == "" {
			continue
		}

		flags += " -" + flagName

		// not add value for boolean setting
		if flagSetting["Type"] != "checkbox" {
			flags += " '" + flagSetting["Value"] + "'"
		}
	}

	return flags
}

func commandOutputs(console string) map[string]string {
	outputs := make(map[string]string)

	// set outputs values from the console
	for _, line := range strings.Split(console, "\n") {
		if line == "" {
			continue
		}

		if o := strings.SplitN(line, ":", 2); len(o) > 1 {
			outputs[o[0]] = o[1]
		}
	}

	return outputs
}
