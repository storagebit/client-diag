package main

import (
	"fmt"
	"strings"
)

func parseOSinfo (o string, k string) {

	strKernel := k
	if strings.Contains(o, "CentOS") {
		redhat_release ,_ := runCommand(strings.Fields("cat /etc/redhat-release"))
		fmt.Println("\t", "Linux Distribution:", strings.TrimSpace(redhat_release))
	} else {
		slcOS := strings.Split(o,"\n")
		for _, line := range slcOS {
			if strings.Contains(line, "Description:") {
				fmt.Println("\t", "Linux Distribution:", strings.TrimSpace(strings.Split(line, ":")[1]))
			}
			if strings.Contains(line, "PRETTY_NAME=") {
				fmt.Println("\t", "Linux Distribution:", strings.Replace(strings.TrimSpace(strings.Split(line, "=")[1]), "\"", "", -1))	
			}
		}
	}
	if len(strKernel) > 1 {
		fmt.Print("\t", " Kernel:", strKernel)
	}
}
