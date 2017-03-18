package core

import (
	"os/exec"
	"strings"
)

type Task struct {
	Name string
	X int
	Y int
	Command string
	Setting map[string]string
	Inputs map[string]string
	Outputs map[string] struct{
		Name string
		Destination []struct{
			Task string
			Input string
		}
	}
}

// run a task
func taskRun(id string, inputs map[string]string, tasks map[string]Task) error {
	// check is the task exist first
	if task, ok := tasks[id]; ok {
		println("taskRun >>", task.Name)

		outputsValues := map[string]string{}

		// if the task have command
		if task.Command != "" {
			execCommand(task.Command, inputs, task.Name)
		}

		// prepare the outputs and categorize them by the destination task
		nextTasks := map[string] map[string]string{}
		for outputID, output := range task.Outputs {
			for _, destination := range output.Destination {
				// set the destination map if it not exist
				if _, ok := nextTasks[destination.Task]; !ok {
					nextTasks[destination.Task] = make(map[string]string)
				}

				// add the output to the destination map
				nextTasks[destination.Task][destination.Input] = outputsValues[outputID]
			}
		}

		// run the next tasks
		for destinationTask, destinationInputs := range nextTasks {
			err := taskRun(destinationTask, destinationInputs, tasks)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func execCommand(command string, inputs map[string]string, taskName string) (map[string]string, error) {
	// receive the inputs
	for inputID, inputValue := range inputs {
		command = strings.Replace(command, "{"+inputID+"}", inputValue, -1)
	}

	outputsValues := map[string]string{}

	// execute the task command
	cmd := strings.Split(command, "|")
	console, err := exec.Command(cmd[0], cmd[1:]...).Output()

	// logConsole(taskName, command, console)

	if err != nil {
		return outputsValues, err
	}

	// set the outputs values from the console
	for _, line := range strings.Split(string(console), "\n") {
		args := strings.Split(line, ":")

		if len(args) < 2 {
			args = append(args, "")
		}

		outputsValues[args[0]] = args[1]
	}

	return outputsValues, nil
}

// log build console output
func logConsole(task string, command string, output []byte) {
	println("==", task, "===============")
	println("COMAND:", command)
	println("+++++++++++++++++++++++++")
	println(string(output))
	// TODO
}

// TODO: execute command timeout
