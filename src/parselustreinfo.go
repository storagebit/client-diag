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
	"bufio"
	"os"
	"strconv"
	"strings"
)

var (
	perfParameters, LUSTRE_CONF = []string{
		"osc.*.max_pages_per_rpc",
		"osc.*.max_rpcs_in_flight",
		"osc.*.max_dirty_mb",
		"osc.*.checksums",
		"llite.*.max_read_ahead_mb",
		"llite.*.max_read_ahead_per_file_mb",
		"ldlm.namespaces.*.lru_size",
		"ldlm.namespaces.*.lru_max_age",
		"mdc.*.max_rpcs_in_flight",
		"mdc.*.max_mod_rpcs_in_flight",
	}, []string{
		"/etc/modprobe.d/lustre.conf",
		"/etc/modprobe.d/lnet.conf",
	}
)

func parseLustrePackages() {
	if checkExecutableExists("rpm") {
		sRPMOutput, _ := runCommand(strings.Fields("rpm -qa"))
		slcRPMOutput := strings.Split(sRPMOutput, "\n")
		for _, line := range slcRPMOutput {
			if strings.Contains(line, "lustre") {
				writeOutLn("\t", line)
				iInstalledLustrePackages++
			}
		}
		writeOutLn("\t", "Found", strconv.Itoa(iInstalledLustrePackages), "installed packages.")
		if iInstalledLustrePackages < 1 {
			bLustreInstalled = false
		} else {
			bLustreInstalled = true
		}
		return

	} else if checkExecutableExists("apt") {
		sDPKGOutput, _ := runCommand(strings.Fields("apt list --installed"))
		slcDPKGOutput := strings.Split(sDPKGOutput, "\n")
		for _, line := range slcDPKGOutput {
			if strings.Contains(line, "lustre") {
				writeOutLn("\t", line)
				iInstalledLustrePackages++
			}
		}
		writeOutLn("\t", "Found", strconv.Itoa(iInstalledLustrePackages), "installed packages.")
		if iInstalledLustrePackages < 1 {
			bLustreInstalled = false
		} else {
			bLustreInstalled = true
		}
		return
	}
}

func parseLoadedLustreKernelModules() {
	if checkIfFileExists("/proc/modules") {
		file, _ := os.Open("/proc/modules")
		fScanner := bufio.NewScanner(file)
		fScanner.Split(bufio.ScanLines)
		for fScanner.Scan() {
			sModule := strings.Split(fScanner.Text(), " ")[0]
			sModinfoOutput, _ := runCommand(strings.Fields("modinfo " + sModule))
			if strings.Contains(sModule, "lnet") {
				bLnetLoaded = true
			}
			if strings.Contains(sModule, "lustre") {
				bLustreLoaded = true
			}
			slcModinfoOutput := strings.Split(sModinfoOutput, "\n")

			if strings.Contains(slcModinfoOutput[3], "Lustre") {
				writeOutLn("\tKernel module", sModule, "is loaded. Details as below:")
				for _, line := range slcModinfoOutput {
					writeOutLn("\t\t", line)
				}
				iLoadedLustreModules++
			}
			file.Close()
		}
	}
	if iLoadedLustreModules < 1 {
		sWarning := "No lustre kernel modules loaded."
		writeOutLn(formatYellow("\tWarning:"), formatYellow(sWarning))
		troubleReport = append(troubleReport, "Lustre kernel module check: "+sWarning)
	}
	return
}

func parseLnetInfo() {
	if checkExecutableExists("lnetctl") {
		sLnetOutput, _ := runCommand(strings.Fields("lnetctl export"))
		slcLnetOutput := strings.Split(sLnetOutput, "\n")

		for _, line := range slcLnetOutput {
			writeOutLn("\t", line)
		}
	}
	return
}

func parseLfsDf() {
	if checkExecutableExists("lfs") {
		sLfsDfOutput, _ := runCommand(strings.Fields("lfs df -h"))
		slcLfsDfOutput := strings.Split(sLfsDfOutput, "\n")

		if len(slcLfsDfOutput) < 2 {
			sWarning := "Cannot read Lustre filesystem information! Is a Lustre filesystem mounted?"
			writeOutLn(formatYellow("\tWarning: " + sWarning))
			troubleReport = append(troubleReport, "Lustre filesystem info: "+sWarning)
			return
		}
		for _, line := range slcLfsDfOutput {
			writeOutLn("\t", line)
		}
	}
	return
}

func parseLustreKernelModuleConfig() {

	for _, sPath := range LUSTRE_CONF {

		if checkIfFileExists(sPath) {
			writeOutLn("\tFound a Lustre kernel module config file at", sPath)
			for _, line := range strings.Split(readFile("etc/modprobe.d/lustre.conf"), "\n") {
				writeOutLn("\t" + line)
			}
		} else {
			sWarning := "No " + sPath + "defined or to be found."
			println(formatYellow("\tWarning: " + sWarning))
			troubleReport = append(troubleReport, "Lustre kernel module config: "+sWarning)
		}
	}
}

func parseLustreFilesystemTuning() {
	if checkExecutableExists("lctl") {
		for _, perfParameter := range perfParameters {
			lctlOutput, _ := runCommand(strings.Fields("lctl get_param " + perfParameter))
			slcLctlOutput := strings.Split(lctlOutput, "\n")
			for _, line := range slcLctlOutput {
				writeOutLn("\t" + strings.Trim(line, "\n"))
			}
		}

	}
}
