package main

import (
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	cmdName := cmd[0]
	command := exec.Command(cmdName, cmd[1:]...)

	for envName, envValue := range env {
		if err := processEnv(envName, envValue); err != nil {
			log.Fatal(err)
		}
	}
	command.Env = os.Environ()

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		//nolint:errorlint
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode()
		}
	}

	return
}

func processEnv(envName string, value EnvValue) (err error) {
	if value.NeedRemove {
		err = os.Unsetenv(envName)
		return
	}

	err = os.Setenv(envName, value.Value)
	return
}
