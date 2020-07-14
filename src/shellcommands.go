package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
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
	outStr := string(stdout.Bytes())
	errStr := string(stderr.Bytes())
	if err != nil {
		fmt.Println(formatYellow("\t" + cmd.String() + " " + err.Error()))
		for _, line := range strings.Split(errStr, "\n"){
			fmt.Println("\t" + line)
		}
	}
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