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
	"strings"
)

func parseLshw(s string) {

	slcLshw := strings.Split(strings.Split(s, "*")[0], "\n")

	writeOutLn(formatBoldWhite("Server/Hostname: " + slcLshw[0]))

	for _, line := range slcLshw {
		if strings.Contains(line, "product") {
			writeOutLn(formatBoldWhite("Platform: " + strings.TrimSpace(strings.Split(line, ":")[1])))
		}
		if strings.Contains(line, "vendor") {
			writeOutLn(formatBoldWhite("Manufacturer: " + strings.TrimSpace(strings.Split(line, ":")[1])))
		}
		if strings.Contains(line, "serial") {
			writeOutLn(formatBoldWhite("Serial #: "), strings.TrimSpace(strings.Split(line, ":")[1]))
		}
	}
}
