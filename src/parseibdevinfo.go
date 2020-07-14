package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func parseIBDEVInfo() {

	if checkExecutableExists("ibv_devinfo") {

		fmt.Println(formatBoldWhite("\nMellanox HCA Details/Information:"))

		strIBDEVInfoOut, _ := runCommand(strings.Fields("ibv_devinfo"))
			slcHCA := strings.Split(strIBDEVInfoOut, "\n\n")
			for _, hca := range slcHCA {
				if len(hca) > 0 {
					r := regexp.MustCompile(`(mlx\d*\_\d*)`)
					hcaId := r.FindStringSubmatch(hca)[0]
					fmt.Println("\tHCA Id:", hcaId)

					r = regexp.MustCompile(`(\d*\.\d*\.\d*)`)
					firmWareLevel := r.FindStringSubmatch(hca)[0]
					fmt.Println("\t\tFirmware level:", firmWareLevel)

					r = regexp.MustCompile(`([A-Za-z0-9]{4}\:[A-Za-z0-9]{4}\:[A-Za-z0-9]{4}\:[A-Za-z0-9]{4})`)
					strGUID := r.FindStringSubmatch(hca)[0]
					fmt.Println("\t\tGUID:", strGUID)
					if strings.Split(strGUID, ":")[0] == "0000" {
						sWarning := "Warning! GUID seems invalid. Please double-check and verify."
						fmt.Println(formatYellow("\t\t" + sWarning))
						mReport[hcaId] = sWarning
					}

					slcPort := regexp.MustCompile(`(?m)^\s*port:\s*\d*`).Split(hca, -1)[1:]

					for i, port := range slcPort {

						portNumber := i + 1
						fmt.Println("\t\t\tPort:", strconv.Itoa(portNumber))

						r = regexp.MustCompile(`\s*link_layer\:\s*([A-Za-z]*)`)
						linkLayer := r.FindStringSubmatch(port)[1]
						fmt.Println("\t\t\t\tLink layer:", linkLayer)

						r = regexp.MustCompile(`state\:\s*PORT_([A-Za-z]*)`)
						portStatus := r.FindStringSubmatch(port)[1]
						fmt.Println("\t\t\t\tStatus:", portStatus)

						r = regexp.MustCompile(`max_mtu\:\s*(\d*)`)
						maxMtu := r.FindStringSubmatch(port)[1]
						fmt.Println("\t\t\t\tMax MTU:", maxMtu)

						r = regexp.MustCompile(`active_mtu\:\s*(\d*)`)
						activeMtu := r.FindStringSubmatch(port)[1]
						fmt.Println("\t\t\t\tActive MTU:", activeMtu)

						intMaxMTU, _ := strconv.Atoi(maxMtu)
						intActiveMTU, _ := strconv.Atoi(activeMtu)

						if intMaxMTU < intActiveMTU {
							sWarning := "Warning! MTU Mismatch!"
							fmt.Println(formatYellow("\t\t\t\t" + sWarning))
						}
					}
					fmt.Print("\n")
				}
			}
		}
	return
}
