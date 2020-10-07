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
	"flag"
	"fmt"
	"strings"
)

var (
	troubleReport             []string
	bPlainOutput              = false
	bLustreInstalled          = false
	iInstalledLustrePackages  = 0
	iLoadedLustreModules      = 0
	bLnetLoaded               = false
	bParseMellanoxLspciOutput = false
	bLustreLoaded             = false
	cpuList                   string
)

func init() {
}

func main() {

	flag.BoolVar(&bPlainOutput, "plain", false, "No color and no formatted output. The console output will be in plain text.")
	flag.Parse()

	checkUser()

	fmt.Println(formatBoldWhite("Starting System Scan. Please wait..."))

	if checkExecutableExists("lshw") {
		strLshwOutput, _ := runCommand(strings.Fields("lshw -c system -quiet"))
		parseLshw(strLshwOutput)
	} else {
		sHostname, _ := runCommand(strings.Fields("hostname"))
		sHostname = strings.TrimRight(sHostname, "\n\r")
		fmt.Println(formatBoldWhite("Server/Hostname:"), sHostname)
	}

	if rootUser() {
		strDmiDecodeOutput, _ := runCommand(strings.Fields("dmidecode --no-sysfs -t baseboard -q"))
		parseDmiDecode(strDmiDecodeOutput)
	}
	if checkExecutableExists("free") {
		strFree, _ := runCommand(strings.Fields("free --si -h"))
		slcInstMem := strings.Split(strFree, "\n")
		fmt.Println(formatBoldWhite("Installed Memory:"), strings.Fields(slcInstMem[1])[1])
	}
	fmt.Println(formatBoldWhite("\nOperating System:"))
	var strOSinfo string
	if checkExecutableExists("lsb_release") {
		strOSinfo, _ = runCommand(strings.Fields("lsb_release -a"))
	} else if checkExecutableExists("hostnamectl") {
		strOSinfo, _ = runCommand(strings.Fields("hostnamectl"))
	} else {
		strOSinfo, _ = runCommand(strings.Fields("cat /etc/os-release"))
	}
	strKernel, _ := runCommand(strings.Fields("uname -r"))
	parseOSinfo(strOSinfo, strKernel)

	fmt.Println(formatBoldWhite("\nCPU Information:"))
	strLsCpuOutput, _ := runCommand(strings.Fields("lscpu"))
	parseLsCpu(strLsCpuOutput)

	checkSystemTuning()

	if !checkFirewall() == false {
		sWarning := "Firewall is running! Make sure that all necessary TCP/IP ports for your environment are open."
		fmt.Println(formatBoldWhite("\nFirewall:"), formatYellow(sWarning))
		troubleReport = append(troubleReport, "Firewall: "+sWarning)

	} else {
		fmt.Println(formatBoldWhite("\nFirewall:"), "No firewall found.")
	}

	fmt.Println(formatBoldWhite("\nInstalled Lustre packages:"))
	parseLustrePackages()

	if bLustreInstalled {
		fmt.Println(formatBoldWhite("\nLoaded Lustre Kernel modules:"))
		parseLoadedLustreKernelModules()

		fmt.Println(formatBoldWhite("\nLustre Kernel module configuration (\"/etc/modprobe.d/lustre.conf\"):"))
		parseLustreKernelModuleConfig()

		fmt.Println(formatBoldWhite("\nLustre Filesystem OST and MST Information:"))
		if bLustreLoaded {
			parseLfsDf()
		} else {
			sWarning := "The Lustre kernel module is not loaded, cannot get OST and MST usage details."
			fmt.Println(formatYellow("\tWarning: " + sWarning))
			troubleReport = append(troubleReport, "Lustre OST and MST info: "+sWarning)
		}
		if rootUser() {
			fmt.Println(formatBoldWhite("\nLustre LNET Information:"))
			if bLnetLoaded {
				parseLnetInfo()
			} else {
				sWarning := "The lustre LNET kernel module is not loaded, cannot get LNET details.."
				fmt.Println(formatYellow("\tWarning: " + sWarning))
				troubleReport = append(troubleReport, "LNET info: "+sWarning)
			}
		}
	}

	if checkExecutableExists("ofed_info") {
		ofedVersion, _ := runCommand(strings.Fields("ofed_info -n"))
		fmt.Println(formatBoldWhite("\nMellanox OFED:"), strings.TrimRight(ofedVersion, "\n\r"))
	} else {
		fmt.Println(formatBoldWhite("\nMellanox OFED:"), "No OFED found.")
	}

	if checkExecutableExists("ibv_devinfo") {
		parseIBDEVInfo()
	} else {
		bParseMellanoxLspciOutput = true
		fmt.Println("\nCannot find \"ibv_devinfo\" in $PATH. Will parse \"lspci -vvv\" output for Mellanox HCA information instead. Is OFED installed?")
	}

	strLspciOutput, _ := runCommand(strings.Fields("lspci -vvv"))
	parseLSPCI(strLspciOutput)

	if checkExecutableExists("ibnetdiscover") {
		fmt.Println(formatBoldWhite("\nInfiniband fabric peers information (\"ibnetdiscover\" output):"))
		sIBNetDiscover, _ := runCommand(strings.Fields("ibnetdiscover -H"))
		slcIBNetDiscover := strings.Split(sIBNetDiscover, "\n")
		for _, line := range slcIBNetDiscover {
			fmt.Println("\t", line)
		}
	} else {
		bParseMellanoxLspciOutput = true
		fmt.Println("\nCannot find \"ibnetdiscover\" in $PATH and therefore no IB fabric peer information will be available. Is OFED installed?")
	}

	fmt.Println(formatBoldWhite("\nIP Network Interface Information:"))

	for _, ipCommand := range []string{"ip -4 a s", "ip -s link"} {

		strIpOutput, _ := runCommand(strings.Fields(ipCommand))
		slcIpOutput := strings.Split(strIpOutput, "\n")
		for _, line := range slcIpOutput {
			fmt.Println("\t", line)
		}
	}

	fmt.Println(formatBoldWhite("\nClient lustre filesystem capacity information:"))

	strDfOutput, err := runCommand(strings.Fields("df -t lustre -H"))
	if len(err) > 0 {
		sWarning := "Cannot find any active lustre filesystem. Are all lustre resources mounted?"
		fmt.Println(formatYellow("\tWarning: " + sWarning))
		troubleReport = append(troubleReport, "Lustre filesystem: "+sWarning)
	} else {
		slcDfOutput := strings.Split(strDfOutput, "\n")
		if len(slcDfOutput) > 0 {
			for _, line := range slcDfOutput {
				fmt.Println("\t", line)
			}
		} else {
			sWarning := "Cannot find any active lustre filesystem. Are all lustre resources mounted?"
			fmt.Println(formatYellow("\tWarning: " + sWarning))
			troubleReport = append(troubleReport, "Lustre filesystem: "+sWarning)
		}
	}

	fmt.Println(formatBoldWhite("\nClient lustre mount information:"))

	strMountOutput, _ := runCommand(strings.Fields("mount -t lustre -l"))
	slcMountOutput := strings.Split(strMountOutput, "\n")
	if len(slcMountOutput) > 1 {
		for _, line := range slcMountOutput {
			fmt.Println("\t", line)
		}
	} else {
		sWarning := "Cannot find any active lustre device mounts. Are all lustre resources mounted?"
		fmt.Println(formatYellow("\tWarning" + sWarning))
		troubleReport = append(troubleReport, "Lustre devices: "+sWarning)
	}

	fmt.Println(formatBoldWhite("\nTrouble Summary:"))
	if len(troubleReport) < 1 {
		fmt.Println(formatGreen("\tNo troubles found."))
	} else {
		for _, line := range troubleReport {
			fmt.Println(formatYellow("\t" + line))
		}
	}
}
