package coresense

import (
	"bytes"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func GetNumCPU() int {
	return runtime.NumCPU()
}

// GetCPUUsageByOtherServices returns the CPU usage by other services
func GetCPUUsageByOtherServices() (float64, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "wmic cpu get loadpercentage")
	} else {
		cmd = exec.Command("sh", "-c", "ps -A -o %cpu | awk '{s+=$1} END {print s}'")
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 0, err
	}

	usageStr := strings.TrimSpace(out.String())
	if runtime.GOOS == "windows" {
		lines := strings.Split(usageStr, "\n")
		if len(lines) < 2 {
			return 0, nil
		}
		usageStr = strings.TrimSpace(lines[1])
	}

	usage, err := strconv.ParseFloat(usageStr, 64)
	if err != nil {
		return 0, err
	}

	return usage, nil
}
