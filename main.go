package main

import (
	"fmt"
	"github.com/prometheus/procfs"
)

func main() {
	// 获取procfs实例
	fs, err := procfs.NewFS("/proc")
	if err != nil {
		fmt.Printf("Error creating procfs: %v\n", err)
		return
	}

	// 获取所有进程
	procs, err := fs.AllProcs()
	if err != nil {
		fmt.Printf("Error getting processes: %v\n", err)
		return
	}

	for _, proc := range procs {
		stat, err := proc.Stat()
		if err != nil {
			continue
		}

		cmdline, _ := proc.CmdLine()

		fmt.Printf("PID: %d\n", proc.PID)
		fmt.Printf("Command: %v\n", cmdline)
		fmt.Printf("State: %s\n", stat.State)
		fmt.Printf("Parent PID: %d\n", stat.PPID)
		fmt.Println("------------------------")
	}
}
