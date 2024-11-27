package main

import (
	"bufio"
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
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
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)

	for _, entry := range entries {
		if entry.IsDir() || !entry.Type().IsRegular() {
			continue
		}

		err = processFile(dir, entry, env)
		if err != nil {
			return nil, err
		}
	}

	return env, nil
}

func processFile(dir string, entry fs.DirEntry, env Environment) error {
	f, err := os.Open(filepath.Join(dir, entry.Name()))
	if err != nil {
		return err
	}
	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		return err
	}

	if fileInfo.Size() == 0 {
		env[entry.Name()] = EnvValue{NeedRemove: true}
		return nil
	}

	scanner := bufio.NewScanner(f)

	if scanner.Scan() {
		firstLine := bytes.ReplaceAll(scanner.Bytes(), []byte{0}, []byte{'\n'})

		envValue := EnvValue{
			Value: strings.TrimRight(string(firstLine), " \t"),
		}

		env[entry.Name()] = envValue
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	return nil
}
