package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func checkExecutableExists(c string) bool {

	_, err := exec.LookPath(c)
	if err != nil {
		return false
	} else {
		return true
	}
}

func checkIfFileExists(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func runCommand(commandParts []string) (string, string) {

	cmd := exec.Command(commandParts[0], commandParts[1:]...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(formatYellow(err.Error()))
	}
	outStr := string(stdout.Bytes())
	errStr := string(stderr.Bytes())
	return outStr, errStr
}

func getCommandReturnCode(commandParts []string) int{

	cmd := exec.Command(commandParts[0], commandParts[1:]...)
	var returnCode int
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			returnCode = exitError.ExitCode()
		} else{
			returnCode = 0
		}
	}
	return returnCode
}