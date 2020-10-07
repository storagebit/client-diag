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
	"strings"
)

func parseOSinfo(o string, k string) {

	strKernel := k
	if strings.Contains(o, "CentOS") {
		redhat_release, _ := runCommand(strings.Fields("cat /etc/redhat-release"))
		fmt.Println("\t", "Linux Distribution:", strings.TrimSpace(redhat_release))
	} else {
		slcOS := strings.Split(o, "\n")
		for _, line := range slcOS {
			if strings.Contains(line, "Description:") {
				fmt.Println("\t", "Linux Distribution:", strings.TrimSpace(strings.Split(line, ":")[1]))
			}
			if strings.Contains(line, "PRETTY_NAME=") {
				fmt.Println("\t", "Linux Distribution:", strings.Replace(strings.TrimSpace(strings.Split(line, "=")[1]), "\"", "", -1))
			}
		}
	}
	if len(strKernel) > 1 {
		fmt.Print("\t", " Kernel:", strKernel)
	}
}
