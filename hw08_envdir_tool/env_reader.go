package main

import (
	"bufio"
	"errors"
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

var (
	ErrNotDirectory   = errors.New("not a directory")
	ErrForbiddenChars = errors.New("forbidden character `=` found in filename")
)

func modifyLine(b []byte) string {
	s := string(b)
	s = strings.ReplaceAll(s, "\x00", "\n")
	s = strings.TrimRight(s, " \t")
	return s
}

func checkName(name string) bool {
	for _, r := range name {
		if r == '=' {
			return false
		}
	}
	return true
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, ErrNotDirectory
	}
	readDir, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	env := make(Environment)
	for _, f := range readDir {
		if f.IsDir() {
			continue
		}
		if ok := checkName(f.Name()); !ok {
			return nil, ErrForbiddenChars
		}
		fPath := path.Join(dir, f.Name())
		file, err := os.Open(fPath)
		if err != nil {
			return nil, err
		}
		stat, err := file.Stat()
		if err != nil {
			return nil, err
		}

		if stat.Size() == 0 {
			env[f.Name()] = EnvValue{NeedRemove: true}
			continue
		}

		reader := bufio.NewReader(file)
		line, _, err := reader.ReadLine()
		if err != nil {
			return nil, err
		}
		env[f.Name()] = EnvValue{Value: modifyLine(line)}
	}
	return env, nil
}
