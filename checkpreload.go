package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func checkLDPreloadEnv() (bool, string) {
	ldPreload := os.Getenv("LD_PRELOAD")
	if ldPreload != "" {
		return true, ldPreload
	}
	return false, ""
}

func checkLDPreloadConfig() (bool, []string) {
	suspiciousFiles := []string{}

	// Check /etc/ld.so.preload
	if _, err := os.Stat("/etc/ld.so.preload"); err == nil {
		content, err := os.ReadFile("/etc/ld.so.preload")
		if err == nil && len(content) > 0 {
			suspiciousFiles = append(suspiciousFiles, "/etc/ld.so.preload")
		}
	}

	// Check /etc/ld.so.conf.d/ directory
	files, err := filepath.Glob("/etc/ld.so.conf.d/*.conf")
	if err == nil {
		for _, file := range files {
			content, err := os.ReadFile(file)
			if err == nil && strings.Contains(string(content), "LD_PRELOAD") {
				suspiciousFiles = append(suspiciousFiles, file)
			}
		}
	}

	return len(suspiciousFiles) > 0, suspiciousFiles
}

func ldPreloadCheck() {
	fmt.Print("===  LD_PRELOAD hook Analysis Results === " + "\n\n")
	// Check environment variable
	if hasEnvHook, envPath := checkLDPreloadEnv(); hasEnvHook {
		fmt.Printf("[Warning] LD_PRELOAD environment variable found: %s\n", envPath)
	}
	// Check configuration files
	if hasConfigHook, configFiles := checkLDPreloadConfig(); hasConfigHook {
		fmt.Println("[Warning] Suspicious preload configuration files found:")
		for _, file := range configFiles {
			fmt.Printf("  - %s\n", file)
		}
	}

}
