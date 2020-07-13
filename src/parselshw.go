package main

import (
	"fmt"
	"strings"
)

func parseLshw (s string){

	slcLshw := strings.Split(strings.Split(s, "*")[0], "\n")

	fmt.Println(formatBoldWhite("Server/Hostname:"), slcLshw[0])

	for _, line := range slcLshw{
		if strings.Contains(line, "product"){
			fmt.Println(formatBoldWhite("Platform:"), strings.TrimSpace(strings.Split(line, ":")[1]))
		}
		if strings.Contains(line, "vendor"){
			fmt.Println(formatBoldWhite("Manufacturer:"), strings.TrimSpace(strings.Split(line, ":")[1]))
		}
		if strings.Contains(line, "serial"){
			fmt.Println(formatBoldWhite("Serial #:"), strings.TrimSpace(strings.Split(line, ":")[1]))
		}
	}
}

