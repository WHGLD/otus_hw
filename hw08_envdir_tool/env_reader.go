package main

import (
	"bufio"
	"bytes"
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
	result := make(Environment)

	envs, err := os.ReadDir(dir)
	if err != nil {
		return result, err
	}

	for _, env := range envs {
		envName := env.Name()
		isNotCorrectEnvName := strings.Contains(envName, "=")
		if isNotCorrectEnvName {
			continue
		}

		envPath := path.Join(dir, envName)
		value, needRemove, errEnv := getEnvValue(envPath)
		if errEnv != nil {
			continue
		}

		envItem := EnvValue{
			Value:      value,
			NeedRemove: needRemove,
		}
		result[envName] = envItem
	}

	return result, nil
}

func getEnvValue(fileName string) (string, bool, error) {
	envFile, errOpen := os.Open(fileName)
	if errOpen != nil {
		return "", false, errOpen
	}
	defer envFile.Close()

	envFileStat, errEnvStat := envFile.Stat()
	if errEnvStat != nil {
		return "", false, errEnvStat
	}
	needRemove := envFileStat.Size() == 0

	firstLine, errLine := getFirstLineOfEnv(envFile)
	if errLine != nil {
		return "", false, errLine
	}

	return firstLine, needRemove, nil
}

func getFirstLineOfEnv(env *os.File) (string, error) {
	reader := bufio.NewReader(env)

	firstLine, err := reader.ReadBytes('\n')
	if err != nil && err != io.EOF {
		return "", err
	}

	firstLine = bytes.Replace(firstLine, []byte("\x00"), []byte("\n"), -1)
	firstLineStr := strings.TrimRight(string(firstLine), "\n ")

	return firstLineStr, nil
}
