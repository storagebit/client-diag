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
	"strconv"
	"strings"
)

func parseLsCpu(s string) {

	slcLsCpu := strings.Split(s, "\n")

	cpuMax := 0
	cpuSpeed := 0

	for _, line := range slcLsCpu {
		if strings.Contains(line, "Architecture") {
			writeOutLn("\t", line)
		}
		if strings.Contains(line, "CPU(s)") {
			writeOutLn("\t", line)
			if strings.Contains(line, "On-line") {
				cpuList = strings.Fields(line)[3]
			}
		}
		if strings.Contains(line, "Thread(s)") {
			writeOutLn("\t", line)
		}
		if strings.Contains(line, "Core(s)") {
			writeOutLn("\t", line)
		}
		if strings.Contains(line, "Socket(s)") {
			writeOutLn("\t", line)
		}
		if strings.Contains(line, "Model name") {
			writeOutLn("\t", line)
		}
		if strings.Split(line, ":")[0] == "CPU MHz" {
			writeOutLn("\t", line)
			strCPUspeed := strings.TrimSpace(strings.Split(strings.Split(line, ":")[1], ".")[0])
			cpuSpeed, _ = strconv.Atoi(strCPUspeed)
		}
		if strings.Split(line, ":")[0] == "CPU max MHz" {
			writeOutLn("\t", line)
			strCPUMaxSpeed := strings.TrimSpace(strings.Split(strings.Split(line, ":")[1], ".")[0])
			cpuMax, _ = strconv.Atoi(strCPUMaxSpeed)

			if cpuMax-cpuSpeed > 100 {
				sWarning := "Warning! The CPU runs on speeds below its capabilities. Please check tuned, c-state and power saving settings."
				writeOutLn(formatYellow("\t " + sWarning))
				troubleReport = append(troubleReport, "CPU: "+sWarning)
			}
		}

	}
}
