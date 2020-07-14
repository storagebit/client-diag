package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseLustrePackages() {
	if checkExecutableExists("rpm") {
		sRPMOutput, _ := runCommand(strings.Fields("rpm -qa"))
		slcRPMOutput := strings.Split(sRPMOutput, "\n")
		fmt.Println(formatBoldWhite("\nInstalled Lustre packages:"))
		for _, line := range slcRPMOutput {
			if strings.Contains(line, "lustre") {
				fmt.Println("\t", line)
				iInstalledLustrePackages++
			}
		}
		fmt.Println("\t", "Found", strconv.Itoa(iInstalledLustrePackages), "installed packages.")
		if iInstalledLustrePackages < 1 {
			bLustreInstalled = false
		} else {
			bLustreInstalled = true
		}
		return

	} else if checkExecutableExists("apt") {
		sDPKGOutput, _ := runCommand(strings.Fields("apt list --installed"))
		slcDPKGOutput := strings.Split(sDPKGOutput, "\n")
		fmt.Println(formatBoldWhite("\nInstalled Lustre packages:"))
		for _, line := range slcDPKGOutput {
			if strings.Contains(line, "lustre") {
				fmt.Println("\t", line)
				iInstalledLustrePackages++
			}
		}
		fmt.Println("\t", "Found", strconv.Itoa(iInstalledLustrePackages), "installed packages.")
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
		fmt.Println(formatBoldWhite("\nLoaded Lustre Kernel modules:"))
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
				fmt.Println("\tKernel module", sModule, "is loaded. Details as below:")
				for _, line := range slcModinfoOutput {
					fmt.Println("\t\t", line)
				}
				iLoadedLustreModules++
			}
			file.Close()
		}
	}
	if iLoadedLustreModules < 1 {
		fmt.Println("\tNo Lustre Kernel module loaded.")
	}
	return
}

func parseLnetInfo() {
	if checkExecutableExists("lnetctl") {
		sLnetOutput, _ := runCommand(strings.Fields("lnetctl export"))
		slcLnetOutput := strings.Split(sLnetOutput, "\n")

		for _, line := range slcLnetOutput {
			fmt.Println("\t", line)
		}
	}
	return
}

func parseLfsDf() {
	if checkExecutableExists("lfs") {
		sLfsDfOutput, _ := runCommand(strings.Fields("lfs df -h"))
		slcLfsDfOutput := strings.Split(sLfsDfOutput, "\n")

		if len(slcLfsDfOutput) < 2 {
			fmt.Println(formatYellow("\tWarning. Cannot read Lustre filesystem information! Is a Lustre filesystem mounted?"))
			return
		}
		for _, line := range slcLfsDfOutput {
			fmt.Println("\t", line)
		}
	}
	return
}

func parseLustreKernelModuleConfig() {
	fmt.Println(formatBoldWhite("\nLustre Kernel module configuration (\"/etc/modprobe.d/lustre.conf\"):"))
	sPath := "/etc/modprobe.d/lustre.conf"

	if checkIfFileExists(sPath) {
		file, _ := os.Open(sPath)
		fScanner := bufio.NewScanner(file)
		fScanner.Split(bufio.ScanLines)

		for fScanner.Scan() {
			fmt.Println("\t" + fScanner.Text())
		}
		file.Close()
	} else {
		println(formatYellow("\tWarning: No \"/etc/modprobe.d/lustre.conf\" defined or to be found."))
	}
}
