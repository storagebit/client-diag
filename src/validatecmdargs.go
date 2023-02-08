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
	"strings"
)

func validateCommandLineArgs() {

	if !checkIfSosreportExists() && bSOSreportRequested {
		fmt.Println("sos not found. Please install it before using the -s option. Will not create and collect a sos report.")
		fmt.Print("Do you want to continue? (y/N): ")
		if !bAnswerYes {
			var strAnswer string
			_, err := fmt.Scanln(&strAnswer)
			if err != nil {
				return
			}
			if strings.ToLower(strAnswer) != "y" {
				os.Exit(0)
			}
		} else {
			bSOSreportRequested = false
		}
	}

	if bKeepDiagBundle && !bCreateClientDiagBundle {
		fmt.Println("The -k option only works if the -c option is used. Will not keep the diagnostic bundle.")
		fmt.Print("Do you want to continue? (y/N): ")
		if !bAnswerYes {
			var strAnswer string
			_, err := fmt.Scanln(&strAnswer)
			if err != nil {
				return
			}
			if strings.ToLower(strAnswer) != "y" {
				os.Exit(0)
			}
		} else {
			bKeepDiagBundle = false
		}
	}

	if bQuietMode && !bCreateClientDiagBundle {
		fmt.Println("The -q option only works if the -c option is used. Will not run in quiet mode.")
		fmt.Print("Do you want to continue? (y/N): ")
		if !bAnswerYes {
			var strAnswer string
			_, err := fmt.Scanln(&strAnswer)
			if err != nil {
				return
			}
			if strings.ToLower(strAnswer) != "y" {
				os.Exit(0)
			}
		} else {
			bQuietMode = false
		}
	}

	if bSOSreportRequested && !bCreateClientDiagBundle {
		fmt.Println("The -s option only works if the -c option is used. Will not create a sos report.")
		fmt.Print("Do you want to continue? (y/N): ")
		if !bAnswerYes {
			var strAnswer string
			_, err := fmt.Scanln(&strAnswer)
			if err != nil {
				log.Println(formatRed("Error: " + err.Error()))
			}
			if strings.ToLower(strAnswer) != "y" {
				os.Exit(0)
			}
		} else {
			bSOSreportRequested = false
		}
	}
}
