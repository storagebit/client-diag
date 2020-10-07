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

func parseLSPCI(s string) {

	slcPCIDevices := strings.Split(s, "\n\n")
	var slcNetworkAdapters []string
	var slcMellanoxHCAs []string

	for _, device := range slcPCIDevices {

		if strings.Contains(strings.ToLower(string(device)), "ethernet controller") {
			slcNetworkAdapters = append(slcNetworkAdapters, device)
		}
		if strings.Contains(strings.ToLower(string(device)), "mellanox") {
			slcMellanoxHCAs = append(slcMellanoxHCAs, device)
		}
	}

	if bParseMellanoxLspciOutput {
		fmt.Println(formatBoldWhite("\nMellanox HCAs: " + strconv.Itoa(len(slcMellanoxHCAs)) + " HCAs found in lspci output."))
		for _, n := range slcMellanoxHCAs {
			parsePCIDeviceDetail(strings.Split(n, "\n"))
		}
	}

	fmt.Println(formatBoldWhite("\nNetwork Ports/Interfaces Found: " + strconv.Itoa(len(slcNetworkAdapters)) + " NICs found."))
	for _, n := range slcNetworkAdapters {
		parsePCIDeviceDetail(strings.Split(n, "\n"))
	}
}

func parsePCIDeviceDetail(d []string) {

	fmt.Println("\t", string(d[0][8:]))
	fmt.Println("\t\t", "PCI Address:\t", d[0][:8])

	strLinkCap := ""
	strLinkSta := ""

	for _, detail := range d {
		if strings.Contains(string(detail), "Subsystem:") {
			fmt.Println("\t\t", "Vendor/OEM:\t", strings.TrimSpace(strings.Split(detail, ":")[1]))
		}
		if strings.Contains(string(detail), "NUMA node:") {
			fmt.Println("\t\t", "NUMA node:\t", strings.TrimSpace(strings.Split(detail, ":")[1]))
		}
		if strings.Contains(string(detail), "LnkCap:") {
			r := regexp.MustCompile(`Speed\s(\d*).*Width\s\d*x(\d*)`)
			fmt.Println("\t\t", "Capabilities:\t Speed", r.FindStringSubmatch(detail)[1], "GT/s, Width x", r.FindStringSubmatch(detail)[2])
			strLinkCap = r.FindStringSubmatch(detail)[1] + r.FindStringSubmatch(detail)[2]
		}
		if strings.Contains(string(detail), "LnkSta:") {
			r := regexp.MustCompile(`LnkSta:\sSpeed\s(\d*).*Width\s\d*x(\d*)`)
			fmt.Println("\t\t", "Status:\t Speed", r.FindStringSubmatch(detail)[1], "GT/s, Width x", r.FindStringSubmatch(detail)[2])
			strLinkSta = r.FindStringSubmatch(detail)[1] + r.FindStringSubmatch(detail)[2]

			if strLinkCap != strLinkSta {
				sWarning := "WARNING! PCI link capabilities <-> PCI link status mismatch!"
				fmt.Println(formatYellow("\t\t " + sWarning))
				troubleReport = append(troubleReport, "PCI Device "+d[0][:8]+": "+d[0][8:]+" "+sWarning)
			}
		}
		if strings.Contains(string(detail), "Part number") {
			r := regexp.MustCompile(`^\s*\[\w*\]\s\w*\s\w*:\s([\w\W]*)`)
			fmt.Println("\t\t", "Part Number:\t", r.FindStringSubmatch(detail)[1])
		}
		if strings.Contains(string(detail), "Serial number") {
			r := regexp.MustCompile(`^\s*\[\w*\]\s\w*\s\w*:\s([\w\W]*)`)
			fmt.Println("\t\t", "Serial Number:\t", r.FindStringSubmatch(detail)[1])
		}
		if strings.Contains(string(detail), "Kernel driver in use:") {
			r := regexp.MustCompile(`^\s*[A-Za-z\s]*\:\s([A-Za-z\_\d]*)`)
			fmt.Println("\t\t", "Kernel Driver:\t", r.FindStringSubmatch(detail)[1])
		}
		if strings.Contains(string(detail), "Kernel modules:") {
			r := regexp.MustCompile(`^\s*[A-Za-z\s]*\:\s([A-Za-z\_\d]*)`)
			fmt.Println("\t\t", "Kernel Module:\t", r.FindStringSubmatch(detail)[1])
		}
	}
}
