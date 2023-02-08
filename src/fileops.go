package main

import (
	"bufio"
	"fmt"
	"os"
)

func readFile(filePath string) string {

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(formatRed("Error opening file: " + err.Error()))
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(formatRed("Error closing file: " + err.Error()))
		}
	}(file)

	scanner := bufio.NewScanner(file)

	if err := scanner.Err(); err != nil {
		fmt.Println(formatRed("Error reading file: " + err.Error()))
	}
	return scanner.Text()
}
