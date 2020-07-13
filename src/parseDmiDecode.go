package main

import (
	"fmt"
	"strings"
)

func parseDmiDecode(s string){

	slcDmiInformation := strings.Split(s, "\n\n")

	for _, line := range strings.Split(slcDmiInformation[0], "\n"){
		if strings.Contains(line, "Product Name"){
			fmt.Println(formatBoldWhite("Baseboard type:"), strings.TrimSpace(strings.Split(line, ":")[1]))
		}
		if strings.Contains(line, "Version"){
			fmt.Println(formatBoldWhite("Baseboard version:"), strings.TrimSpace(strings.Split(line, ":")[1]))
		}
		if strings.Contains(line, "Serial Number"){
			fmt.Println(formatBoldWhite("Baseboard serial #:"), strings.TrimSpace(strings.Split(line, ":")[1]))
		}
	}
}