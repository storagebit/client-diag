package main

import (
	"fmt"
	"strconv"
	"strings"
)

func parseLsCpu (s string){

	slcLsCpu := strings.Split(s, "\n")

	cpuMax := 0
	cpuSpeed := 0

	for _, line := range slcLsCpu{
		if strings.Contains(line, "Architecture"){
			fmt.Println("\t", line)
		}
		if strings.Contains(line, "CPU(s)"){
			fmt.Println("\t", line)
		}
		if strings.Contains(line, "Thread(s)"){
			fmt.Println("\t", line)
		}
		if strings.Contains(line, "Core(s)"){
			fmt.Println("\t", line)
		}
		if strings.Contains(line, "Socket(s)"){
			fmt.Println("\t", line)
		}
		if strings.Contains(line, "Model name"){
			fmt.Println("\t", line)
		}
		if strings.Split(line, ":")[0] == "CPU MHz"{
			fmt.Println("\t", line)
			strCPUspeed := strings.TrimSpace(strings.Split(strings.Split(line, ":")[1], ".")[0])
			cpuSpeed, _ = strconv.Atoi(strCPUspeed)
		}
		if strings.Split(line, ":")[0] == "CPU max MHz"{
			fmt.Println("\t", line)
			strCPUMaxSpeed := strings.TrimSpace(strings.Split(strings.Split(line, ":")[1], ".")[0])
			cpuMax, _ = strconv.Atoi(strCPUMaxSpeed)

			if cpuMax - cpuSpeed > 100{
				sWarning := "Warning! The CPU runs on speeds below its capabilities. Please check tuned, c-state and power saving settings."
				fmt.Println(formatYellow("\t" + sWarning))
				mReport["CPU"] = sWarning
			}
		}

	}
}
