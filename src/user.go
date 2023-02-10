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
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
)

func rootUser() bool {
	currentUser, userError := user.Current()

	if userError != nil {
		log.Fatalf("Error while retrieving current user details: '%s'", userError)
	}

	if currentUser.Uid == "0" {
		return true
	} else {
		return false
	}
}

func checkUser() {

	currentUser, err := user.Current()

	if err != nil {
		log.Fatalf("Error while retrieving current user details: '%s'", err)
	}

	writeOutLn("Running client-diag as user: " + currentUser.Username)

	if !rootUser() {

		writeOutLn("Executing client-diag without root privileges or sudo will limit the dianostic/reporting capabilities.\n" +
			"Run as root or sudo if you want to see more.")
		if !bAnswerYes {
			fmt.Print("Do you want to continue? [y/N]: ")
			var input string
			_, err := fmt.Scanln(&input)
			if err != nil {
				return
			}
			if strings.ToLower(input) != "y" {
				writeOutLn("Exiting...")
				os.Exit(2)
			}
		}
	} else {

		writeOutLn("client-diag is being executed with elevated/root privileges.")
	}

}
