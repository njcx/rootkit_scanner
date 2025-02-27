package main

import (
	"fmt"
	"github.com/mitchellh/go-ps"
)

func main() {
	// 获取所有进程
	processes, err := ps.Processes()
	if err != nil {
		fmt.Printf("Error getting processes: %v\n", err)
		return
	}

	for _, process := range processes {
		fmt.Printf("PID: %d\n", process.Pid())
		fmt.Printf("Parent PID: %d\n", process.PPid())
		fmt.Printf("Executable: %s\n", process.Executable())
		fmt.Println("------------------------")
	}
}
