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
	"strconv"
	"strings"
)

var (
	sysctlParameters = []string{
		"net.ipv4.conf.all.rp_filter",
		"net.ipv4.conf.all.arp_announce",
		"net.ipv4.conf.all.arp_ignore",
		"net.ipv4.conf.default.rp_filter",
		"net.ipv4.conf.default.arp_announce",
		"net.ipv4.conf.default.arp_ignore",
	}
)

func checkSystemTuning() {
	writeOutLn(formatBoldWhite("\nServer/Client OS and System Tuning Settings:"))
	checkTuneD()
	checkIRQBalance()
	checkCstate()
	checkOSTunables()
	checkCPUScalingGovernor()
}

func checkCstate() {
	if checkIfFileExists("/sys/module/intel_idle/parameters/max_cstate") {
		cState := readFile("/sys/module/intel_idle/parameters/max_cstate")

		sCstateSetting := strings.TrimRight(string(cState), "\n\r")
		iCstateSetting, _ := strconv.Atoi(sCstateSetting)

		if iCstateSetting > 0 {
			sWarning := "The CPU maximum cstate setting is at " + sCstateSetting + "! For best performance it is recommended to configure it to 0."
			writeOutLn("\tCPU maximum cstate setting: ", formatYellow(sWarning))
			troubleReport = append(troubleReport, "CPU cstate: "+sWarning)
			return
		} else {
			writeOutLn("\tCPU maximum cstate setting: " + sCstateSetting + " - " + formatGreen("OK"))
			return
		}
	} else {
		writeOutLn("\tCannot read the systems cstate setting. Cannot read the file '/sys/module/intel_idle/parameters/max_cstate'")
	}
}

func checkCPUScalingGovernor() {
	writeOutLn("\n\tReading additional CPU information. ")
	if checkExecutableExists("cpupower") {
		cpuPowerOutput, _ := runCommand(strings.Fields("cpupower -c " + cpuList + " frequency-info"))
		for _, line := range strings.Split(cpuPowerOutput, "\n") {
			writeOutLn("\t", line)
		}
	}
}

func checkTuneD() {
	tunedReturnCode := getCommandReturnCode(strings.Fields("systemctl status tuned"))

	if tunedReturnCode == 0 {
		tunedProfile, _ := runCommand(strings.Fields("tuned-adm active"))
		{
			if strings.Contains(tunedProfile, "latency-performance") {
				writeOutLn("\tTuneD service: ", formatGreen("OK"))
				return
			} else {
				sWarning := "To achieve best performance you should run tuned with the 'latency-performance' profile."
				writeOutLn("\tTuneD service: ", formatYellow(sWarning))
				troubleReport = append(troubleReport, "TuneD service: "+sWarning)
				return
			}
		}
	}

	if tunedReturnCode == 4 {
		sWarning := "Not installed! To achieve best performance you should run tuned in the 'latency-performance' profile."
		writeOutLn("\tTuneD service: ", formatYellow(sWarning))
		troubleReport = append(troubleReport, "TuneD service: "+sWarning)
		return
	}
	if tunedReturnCode != 0 {
		sWarning := "Service failed! Please check the tuned service."
		writeOutLn("\tTuneD service: ", formatYellow(sWarning))
		troubleReport = append(troubleReport, "TuneD service: "+sWarning)
		return
	}
}

func checkIRQBalance() {
	irqBalance := getCommandReturnCode(strings.Fields("systemctl status irqbalance"))
	if irqBalance == 4 {
		sWarning := "Not installed! To achieve best performance and user experience you should run IRQ Balance."
		writeOutLn("\tIRQ Balance: ", formatYellow(sWarning))
		troubleReport = append(troubleReport, "IRQ Balance: "+sWarning)
		return
	}
	if irqBalance != 0 {
		sWarning := "IRQBalance service is failed! Please check the IRQ Balance service."
		writeOutLn("\tIRQ Balance: ", formatYellow(sWarning))
		troubleReport = append(troubleReport, "IRQ Balance: "+sWarning)
		return
	}
	if irqBalance == 0 {
		writeOutLn("\tIRQBalance: ", formatGreen("OK"))
		return
	}
}

func checkOSTunables() {
	if checkExecutableExists("cat") {
		fmt.Print("\t/proc/sys/fs/aio-max-nr: ", strings.Trim(runCommand(
			strings.Fields("cat /proc/sys/fs/aio-max-nr"))))
		fmt.Print("\t/sys/kernel/mm/transparent_hugepage/enabled: ", strings.Trim(runCommand(
			strings.Fields("cat /sys/kernel/mm/transparent_hugepage/enabled"))))
	}
	if checkExecutableExists("sysctl") {
		for _, sysctlParameter := range sysctlParameters {
			fmt.Print("\t", strings.Trim(runCommand(strings.Fields("sysctl "+sysctlParameter))))
		}
	}
}
