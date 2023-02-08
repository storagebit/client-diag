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
	"regexp"
	"strconv"
	"strings"
)

func parseIBDEVInfo() {

	if checkExecutableExists("ibv_devinfo") {

		strIBDEVInfoOut, _ := runCommand(strings.Fields("ibv_devinfo"))
		slcHCA := strings.Split(strIBDEVInfoOut, "\n\n")
		for _, hca := range slcHCA {
			if len(hca) > 0 {
				r := regexp.MustCompile(`(mlx\d*\_\d*)`)
				hcaId := r.FindStringSubmatch(hca)[0]
				writeOutLn("\tHCA Id:", hcaId)

				r = regexp.MustCompile(`(\d*\.\d*\.\d*)`)
				firmWareLevel := r.FindStringSubmatch(hca)[0]
				writeOutLn("\t\tFirmware level:", firmWareLevel)

				r = regexp.MustCompile(`([A-Za-z0-9]{4}\:[A-Za-z0-9]{4}\:[A-Za-z0-9]{4}\:[A-Za-z0-9]{4})`)
				strGUID := r.FindStringSubmatch(hca)[0]
				writeOutLn("\t\tGUID:", strGUID)
				if strings.Split(strGUID, ":")[0] == "0000" {
					sWarning := "Warning! GUID seems invalid. Please double-check and verify."
					writeOutLn(formatYellow("\t\t" + sWarning))
					troubleReport = append(troubleReport, hcaId+": "+sWarning)
				}

				slcPort := regexp.MustCompile(`(?m)^\s*port:\s*\d*`).Split(hca, -1)[1:]

				for i, port := range slcPort {

					portNumber := i + 1
					writeOutLn("\t\t\tPort:", strconv.Itoa(portNumber))

					r = regexp.MustCompile(`\s*link_layer\:\s*([A-Za-z]*)`)
					linkLayer := r.FindStringSubmatch(port)[1]
					writeOutLn("\t\t\t\tLink layer:", linkLayer)

					r = regexp.MustCompile(`state\:\s*PORT_([A-Za-z]*)`)
					portStatus := r.FindStringSubmatch(port)[1]
					writeOutLn("\t\t\t\tStatus:", portStatus)

					r = regexp.MustCompile(`max_mtu\:\s*(\d*)`)
					maxMtu := r.FindStringSubmatch(port)[1]
					writeOutLn("\t\t\t\tMax MTU:", maxMtu)

					r = regexp.MustCompile(`active_mtu\:\s*(\d*)`)
					activeMtu := r.FindStringSubmatch(port)[1]
					writeOutLn("\t\t\t\tActive MTU:", activeMtu)

					intMaxMTU, _ := strconv.Atoi(maxMtu)
					intActiveMTU, _ := strconv.Atoi(activeMtu)

					if intMaxMTU != intActiveMTU {
						sWarning := "Warning! MTU Mismatch!"
						writeOutLn(formatYellow("\t\t\t\t" + sWarning))
					}
				}
				fmt.Print("\n")
			}
		}
	}
	return
}
