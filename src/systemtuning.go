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
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func checkSystemTuning() {
	fmt.Println(formatBoldWhite("\nServer/Client System Tuning Settings:"))
	checkTuneD()
	checkIRQBalance()
	checkCstate()
	checkCPUScalingGovernor()
}

func checkCstate() {
	if checkIfFileExists("/sys/module/intel_idle/parameters/max_cstate") {
		c, err := ioutil.ReadFile("/sys/module/intel_idle/parameters/max_cstate")
		if err != nil {
			log.Fatal(err)
		}
		sCstateSetting := strings.TrimRight(string(c), "\n\r")
		iCstateSetting, err := strconv.Atoi(sCstateSetting)

		if iCstateSetting > 0 {
			sWarning := "The CPU maximum cstate setting is at " + sCstateSetting + "! For best performance it is recommended to configure it to 0."
			fmt.Println("\tCPU maximum cstate setting: ", formatYellow(sWarning))
			troubleReport = append(troubleReport, "CPU cstate: "+sWarning)
			return
		} else {
			fmt.Println("\tCPU maximum cstate setting: " + sCstateSetting + " - " + formatGreen("OK"))
			return
		}
	} else {
		fmt.Println("\tCannot read the systems cstate setting. Cannot read the file '/sys/module/intel_idle/parameters/max_cstate'")
	}
}

func checkCPUScalingGovernor() {
	fmt.Println("\n\tReading additional CPU information. ")
	if checkExecutableExists("cpupower") {
		cpuPowerOutput, _ := runCommand(strings.Fields("cpupower -c " + cpuList + " frequency-info"))
		for _, line := range strings.Split(cpuPowerOutput, "\n") {
			fmt.Println("\t", line)
		}
	}
}

func checkTuneD() {
	tunedReturnCode := getCommandReturnCode(strings.Fields("systemctl status tuned"))
	if tunedReturnCode == 4 {
		sWarning := "Not installed! To achieve best performance you should run tuned in the 'latency-performance' profile."
		fmt.Println("\tTuneD service: ", formatYellow(sWarning))
		troubleReport = append(troubleReport, "TuneD service: "+sWarning)
		return
	}
	if tunedReturnCode != 0 {
		sWarning := "Service failed! Please check the tuned service."
		fmt.Println("\tTuneD service: ", formatYellow(sWarning))
		troubleReport = append(troubleReport, "TuneD service: "+sWarning)
		return
	}
	if tunedReturnCode == 0 {
		tunedProfile, _ := runCommand(strings.Fields("tuned-adm active"))
		{
			if strings.Contains(tunedProfile, "latency-performance") {
				fmt.Println("\tTuneD service: ", formatGreen("OK"))
				return
			} else {
				sWarning := "TuneD is not running. To achieve best performance you should run tuned with the 'latency-performance' profile."
				fmt.Println("\tTuneD service: ", formatYellow(sWarning))
				troubleReport = append(troubleReport, "TuneD service: "+sWarning)
				return
			}
		}
	}
}

func checkIRQBalance() {
	irqBalance := getCommandReturnCode(strings.Fields("systemctl status irqbalance"))
	if irqBalance == 4 {
		sWarning := "Not installed! To achieve best performance and user experience you should run IRQ Balance."
		fmt.Println("\tIRQ Balance: ", formatYellow(sWarning))
		troubleReport = append(troubleReport, "IRQ Balance: "+sWarning)
		return
	}
	if irqBalance != 0 {
		sWarning := "IRQBalance service is failed! Please check the IRQ Balance service."
		fmt.Println("\tIRQ Balance: ", formatYellow(sWarning))
		troubleReport = append(troubleReport, "IRQ Balance: "+sWarning)
		return
	}
	if irqBalance == 0 {
		fmt.Println("\tIRQBalance: ", formatGreen("OK"))
		return
	}
}
