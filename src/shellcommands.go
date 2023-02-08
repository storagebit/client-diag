/*The MIT License (MIT)
Copyright © 2020 StorageBIT.ch
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
documentation files (the “Software”), to deal in the Software without restriction, including without limitation the
rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit
persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
Software.
THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO
THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.IN NO EVENT SHALL THE AUTHORS
OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"bytes"
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
		writeOutLn(formatYellow("\t" + cmd.String() + " " + err.Error()))
		for _, line := range strings.Split(errStr, "\n") {
			writeOutLn("\t" + line)
		}
	}
	return outStr, errStr
}

func getCommandReturnCode(commandParts []string) int {

	cmd := exec.Command(commandParts[0], commandParts[1:]...)
	var returnCode int
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			returnCode = exitError.ExitCode()
		} else {
			returnCode = 0
		}
	}
	return returnCode
}
