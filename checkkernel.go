package main

import (
	"bufio"
	"fmt"
	"os"
)

const filename = "/proc/rsc_lkm"

func getDriverInfo() {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Print(line)
				break
			}
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(line)
	}
}

func banner() {
	banner := `
 ____   __    ___  ___  __  __
 || \\ (( \  //   // \\ ||\ ||
 ||_//  \\  ((    ||=|| ||\\||
 || \\ \_))  \\__ || || || \||
    `
	fmt.Println(banner)
}
