package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ProcessInfo stores process information
type ProcessInfo struct {
	PID       int
	PPID      int
	Name      string
	UID       int
	State     string
	StartTime int64  // Process start time (Unix timestamp)
	CmdLine   string // Full command line
}

// ProcessMap is used to store a mapping of process information
type ProcessMap struct {
	sync.RWMutex
	processes map[int]*ProcessInfo
}

func NewProcessMap() *ProcessMap {
	return &ProcessMap{
		processes: make(map[int]*ProcessInfo),
	}
}

func (pm *ProcessMap) Add(pid int, info *ProcessInfo) {
	pm.Lock()
	defer pm.Unlock()
	pm.processes[pid] = info
}

func (pm *ProcessMap) Get(pid int) (*ProcessInfo, bool) {
	pm.RLock()
	defer pm.RUnlock()
	info, exists := pm.processes[pid]
	return info, exists
}

func (pm *ProcessMap) Len() int {
	pm.RLock()
	defer pm.RUnlock()
	return len(pm.processes)
}

func (pm *ProcessMap) GetPIDs() []int {
	pm.RLock()
	defer pm.RUnlock()
	pids := make([]int, 0, len(pm.processes))
	for pid := range pm.processes {
		pids = append(pids, pid)
	}
	return pids
}

// getBootTime retrieves system boot time
func getBootTime() (int64, error) {
	content, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return 0, err
	}

	for _, line := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(line, "btime") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				return strconv.ParseInt(fields[1], 10, 64)
			}
		}
	}
	return 0, fmt.Errorf("boot time not found in /proc/stat")
}

// getProcessStartTime retrieves process start time
func getProcessStartTime(pid int) (int64, error) {
	statFile := fmt.Sprintf("/proc/%d/stat", pid)
	content, err := ioutil.ReadFile(statFile)
	if err != nil {
		return 0, err
	}

	fields := strings.Fields(string(content))
	if len(fields) < 22 {
		return 0, fmt.Errorf("invalid stat file format")
	}

	startTime, err := strconv.ParseInt(fields[21], 10, 64)
	if err != nil {
		return 0, err
	}

	bootTime, err := getBootTime()
	if err != nil {
		return 0, err
	}

	// Convert start time to Unix timestamp
	clockTicks := float64(1)
	if hz, err := getClockTicks(); err == nil {
		clockTicks = float64(hz)
	}
	return bootTime + int64(float64(startTime)/clockTicks), nil
}

// getClockTicks retrieves system clock frequency
func getClockTicks() (int64, error) {
	cmd := exec.Command("getconf", "CLK_TCK")
	output, err := cmd.Output()
	if err != nil {
		return 100, nil // Default value
	}
	return strconv.ParseInt(strings.TrimSpace(string(output)), 10, 64)
}

// getCmdLine retrieves the full command line of a process
func getCmdLine(pid int) string {
	cmdline, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return ""
	}
	// Replace \0 with space
	return strings.ReplaceAll(string(cmdline), "\x00", " ")
}

// getProcProcesses retrieves process information from the /proc filesystem
func getProcProcesses(minAge time.Duration) (*ProcessMap, error) {
	procMap := NewProcessMap()
	files, err := ioutil.ReadDir("/proc")
	if err != nil {
		return nil, fmt.Errorf("error reading /proc: %v", err)
	}

	now := time.Now().Unix()

	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		pid, err := strconv.Atoi(f.Name())
		if err != nil {
			continue
		}

		startTime, err := getProcessStartTime(pid)
		if err != nil {
			continue
		}

		// Skip processes running for less than minAge
		if now-startTime < int64(minAge.Seconds()) {
			continue
		}

		statusFile := filepath.Join("/proc", f.Name(), "status")
		content, err := ioutil.ReadFile(statusFile)
		if err != nil {
			continue
		}

		info := &ProcessInfo{
			PID:       pid,
			StartTime: startTime,
			CmdLine:   getCmdLine(pid),
		}

		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			switch {
			case strings.HasPrefix(line, "Name:"):
				info.Name = strings.TrimSpace(strings.TrimPrefix(line, "Name:"))
			case strings.HasPrefix(line, "PPid:"):
				info.PPID, _ = strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "PPid:")))
			case strings.HasPrefix(line, "Uid:"):
				uidFields := strings.Fields(strings.TrimPrefix(line, "Uid:"))
				if len(uidFields) > 0 {
					info.UID, _ = strconv.Atoi(uidFields[0])
				}
			case strings.HasPrefix(line, "State:"):
				stateParts := strings.Fields(strings.TrimPrefix(line, "State:"))
				if len(stateParts) > 0 {
					info.State = string(stateParts[0][0])
				}
			}
		}

		procMap.Add(pid, info)
	}

	return procMap, nil
}

// getPsProcesses retrieves process information from the ps command
func getPsProcesses(minAge time.Duration) (*ProcessMap, error) {
	procMap := NewProcessMap()

	// Use ps command to get more detailed information, including start time
	cmd := exec.Command("ps", "ax", "-o", "pid,ppid,uid,stat,lstart,comm")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error executing ps command: %v", err)
	}

	now := time.Now()
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	scanner.Scan() // Skip header line

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 6 {
			continue
		}

		pid, err := strconv.Atoi(fields[0])
		if err != nil {
			continue
		}

		// Parse start time
		timeStr := strings.Join(fields[4:9], " ")
		startTime, err := time.Parse("Mon Jan 2 15:04:05 2006", timeStr)
		if err != nil {
			continue
		}

		// Skip processes running for less than minAge
		if now.Sub(startTime) < minAge {
			continue
		}

		ppid, _ := strconv.Atoi(fields[1])
		uid, _ := strconv.Atoi(fields[2])

		info := &ProcessInfo{
			PID:       pid,
			PPID:      ppid,
			UID:       uid,
			State:     string(fields[3][0]),
			StartTime: startTime.Unix(),
			Name:      fields[len(fields)-1],
			CmdLine:   getCmdLine(pid),
		}

		procMap.Add(pid, info)
	}

	return procMap, nil
}

// compareProcesses compares process information obtained by two methods
func compareProcesses(procProcesses, psProcesses *ProcessMap) {

	procPIDs := procProcesses.GetPIDs()
	_ = psProcesses.GetPIDs()

	// Check for hidden processes
	hiddenProcesses := make([]int, 0)
	for _, pid := range procPIDs {
		if _, exists := psProcesses.Get(pid); !exists {
			hiddenProcesses = append(hiddenProcesses, pid)
		}
	}

	if len(hiddenProcesses) > 0 {
		fmt.Printf("[Warning] Found %d hidden processes:\n", len(hiddenProcesses))
		for _, pid := range hiddenProcesses {
			if procInfo, exists := procProcesses.Get(pid); exists {
				fmt.Printf("PID: %d\nName: %s\nUID: %d\nStart Time: %s\nCommand: %s\n\n",
					pid,
					procInfo.Name,
					procInfo.UID,
					time.Unix(procInfo.StartTime, 0).Format("2006-01-02 15:04:05"),
					procInfo.CmdLine)
			}
		}
		fmt.Println("[Warning] The ps command may have been replaced. ")
	}

}

func psCheck() {

	fmt.Print("\n" + "===  Ps integrity Analysis Results === " + "\n\n")
	// Only check processes running for more than 5 minutes
	minAge := 1 * time.Minute
	procProcesses, err := getProcProcesses(minAge)
	if err != nil {
		log.Fatalf("Error getting processes from /proc: %v", err)
	}
	psProcesses, err := getPsProcesses(minAge)
	if err != nil {
		log.Fatalf("Error getting processes from ps command: %v", err)
	}
	compareProcesses(procProcesses, psProcesses)
}
