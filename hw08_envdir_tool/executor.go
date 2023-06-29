package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

var ee *exec.ExitError

// RunCmd - функция запуска процесса ОС с применением к нему окружения переменных.
func RunCmd(cmd []string, environment Environment) (returnCode int) {
	Apply(environment)

	if len(cmd) == 0 {
		return
	}
	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Run()
	if err != nil {
		if errors.As(err, &ee) {
			returnCode = ee.ExitCode()
		}
	}
	return
}

// RunCmdVariant2 - функция запуска процесса ОС с применением к нему окружения переменных.
// В данной реализации задействуется установление контекста посредством `command.Env`.
func RunCmdVariant2(cmd []string, environment Environment) (returnCode int) {
	UnsetOnlyNeeded(environment)

	if len(cmd) == 0 {
		return
	}
	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	command.Env = os.Environ()
	for envVarName := range environment {
		if !environment[envVarName].NeedRemove {
			command.Env = append(command.Env, fmt.Sprintf("%s=%s", envVarName, environment[envVarName].Value))
		}
	}

	err := command.Run()
	if err != nil {
		if errors.As(err, &ee) {
			returnCode = ee.ExitCode()
		}
	}
	return
}
