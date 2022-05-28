package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var ErrUnsupportedFileName = errors.New("file name cannot contain '='")

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	environment := make(Environment, len(files))
	for _, fileInfo := range files {

		envValue, err := fileToEnvValue(dir, fileInfo)
		if err != nil {
			return nil, err
		}

		environment[fileInfo.Name()] = envValue
	}

	return environment, nil
}

func fileToEnvValue(dir string, info os.FileInfo) (envValue EnvValue, err error) {
	fileName := info.Name()
	if strings.Contains(fileName, "=") {
		err = ErrUnsupportedFileName
		return
	}

	if info.Size() == 0 {
		return EnvValue{NeedRemove: true}, nil
	}

	file, err := os.Open(dir + string(os.PathSeparator) + fileName)
	if err != nil {
		return
	}

	b := bufio.NewReader(file)
	line, err := b.ReadBytes('\n')
	if err != nil {
		if err != io.EOF {
			return
		}
	} else {
		// Trim \n at the end of line.
		line = line[:len(line)-1]
	}

	value := strings.TrimRight(
		string(bytes.ReplaceAll(
			line,
			[]byte{'\x00'},
			[]byte{'\n'},
		)),
		" \t",
	)

	return EnvValue{Value: value}, nil
}
