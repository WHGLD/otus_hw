package main

import "os"

func main() {
	envDir := os.Args[1]
	cmd := os.Args[2:]

	envs, err := ReadDir(envDir)
	if err != nil {
		return
	}

	os.Exit(RunCmd(cmd, envs))
}
