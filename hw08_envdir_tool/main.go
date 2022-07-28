package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 3 {
		log.Fatal("Command to run was not set. Usage: go-envdir dirName cmd")
	}

	env, err := ReadDir(args[1])
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(RunCmd(args[2:], env))
}
