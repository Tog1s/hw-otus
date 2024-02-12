package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	result := Environment{}
	for _, item := range files {
		if strings.Contains(item.Name(), "=") {
			fmt.Printf("item %s contains = in name", item.Name())
			continue
		}

		s, err := readFirstLine(path.Join(dir, item.Name()))
		needRemove := false
		if err != nil {
			return nil, err
		}

		if s == "" {
			needRemove = true
		}
		result[item.Name()] = EnvValue{s, needRemove}
	}

	return result, nil
}

func readFirstLine(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	line, err := reader.ReadBytes('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}

	line = bytes.ReplaceAll(line, []byte("\n"), []byte(""))
	line = bytes.ReplaceAll(line, []byte("\x00"), []byte("\n"))

	return strings.TrimRight(string(line), " \t"), nil
}
