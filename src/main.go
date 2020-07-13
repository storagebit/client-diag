package main

import (
	"flag"
	"fmt"
	"strings"
)

var (
	mReport                  = make(map[string]string)
	plainOutput              = false
	lustreInstalled          = false
	iInstalledLustrePackages = 0
	iLoadedLustreModules     = 0
	bLnetLoaded              = false
)

func init() {
}

func main() {

	flag.BoolVar(&plainOutput, "plain", false, "No color and no formatted console output in plain text.")
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

	parseLustrePackages()

	if lustreInstalled {
		println(formatBoldWhite("\nLoaded Lustre kernel modules:"))
		parseLoadedLustreKernelModules()
		parseLustreKernelModuleConfig()
		if rootUser() {
			fmt.Println(formatBoldWhite("\nLustre LNET Information:"))
			if bLnetLoaded {
				parseLnetInfo()
			} else {
				fmt.Println("\tThe Lustre LNET kernel module is not loaded, cannot get LNET details.")
			}
		}
	}

	if checkExecutableExists("ofed_info") {
		ofedVersion, _ := runCommand(strings.Fields("ofed_info -n"))
		fmt.Println(formatBoldWhite("\nMellanox OFED:"), strings.TrimRight(ofedVersion, "\n\r"))
	} else {
		fmt.Println(formatBoldWhite("\nMellanox OFED:"), "No OFED found.")
	}

	if !checkFirewall() == false {
		sWarning := "Firewall is running! Make sure that all necessary TCP/IP ports for your environment are open."
		fmt.Println(formatBoldWhite("\nFirewall:"), formatYellow(sWarning))
		mReport["Firewall"] = sWarning
	} else {
		fmt.Println(formatBoldWhite("\nFirewall:"), "No firewall found.")
	}

	strLspciOutput, _ := runCommand(strings.Fields("lspci -vvv"))
	parseLSPCI(strLspciOutput)

	fmt.Println(formatBoldWhite("\nIP Network Interface Information:"))

	for _, ipCommand := range []string{"ip -4 a s", "ip -s link"} {

		strIpOutput, _ := runCommand(strings.Fields(ipCommand))
		slcIpOutput := strings.Split(strIpOutput, "\n")
		for _, line := range slcIpOutput {
			fmt.Println("\t", line)
		}
	}

	fmt.Println(formatBoldWhite("\nClient filesystem consumption information:"))

	strDfOutput, _ := runCommand(strings.Fields("df -H"))
	slcDfOutput := strings.Split(strDfOutput, "\n")
	for _, line := range slcDfOutput {
		fmt.Println("\t", line)
	}

	fmt.Println(formatBoldWhite("\nClient mount information:"))

	strMountOutput, _ := runCommand(strings.Fields("mount -l"))
	slcMountOutput := strings.Split(strMountOutput, "\n")
	for _, line := range slcMountOutput {
		fmt.Println("\t", line)
	}

	parseIBDEVInfo()

	fmt.Println(formatBoldWhite("\nSummary:"))
	if len(mReport) < 1 {
		fmt.Println("\t No troubles found.")
	} else {
		for k, v := range mReport {
			fmt.Println("\t", k, ":", v)
		}
	}
}
