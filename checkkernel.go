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
 ____    ___     ___   ______ __ __ __ ______  __    ___  ___  __  __ __  __  ____ ____ 
 || \\  // \\   // \\  | || | || // || | || | (( \  //   // \\ ||\ || ||\ || ||    || \\
 ||_// ((   )) ((   ))   ||   ||<<  ||   ||    \\  ((    ||=|| ||\\|| ||\\|| ||==  ||_//
 || \\  \\_//   \\_//    ||   || \\ ||   ||   \_))  \\__ || || || \|| || \|| ||___ || \\
    `
	fmt.Println(banner)
}
