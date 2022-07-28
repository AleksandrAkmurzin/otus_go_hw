package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
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
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	environment := make(Environment, len(dirEntries))
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}

		envValue, err := fileToEnvValue(dir, dirEntry)
		if err != nil {
			return nil, err
		}

		environment[dirEntry.Name()] = envValue
	}

	return environment, nil
}

func fileToEnvValue(dir string, fileEntry os.DirEntry) (envValue EnvValue, err error) {
	fileName := fileEntry.Name()
	if strings.Contains(fileName, "=") {
		err = ErrUnsupportedFileName
		return
	}

	fileInfo, err := fileEntry.Info()
	if err != nil {
		return
	}
	if fileInfo.Size() == 0 {
		return EnvValue{NeedRemove: true}, nil
	}

	file, err := os.Open(filepath.Join(dir, fileName))
	if err != nil {
		return
	}
	defer file.Close()

	b := bufio.NewReader(file)
	line, err := b.ReadBytes('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return
	}

	value := strings.TrimRight(
		string(bytes.ReplaceAll(
			bytes.TrimRight(line, "\n"),
			[]byte{'\x00'},
			[]byte{'\n'},
		)),
		" \t",
	)

	return EnvValue{Value: value}, nil
}
