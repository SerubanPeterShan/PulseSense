package processsense

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type windowsMonitor struct{}

func NewProcessMonitor() ProcessMonitor {
	return &windowsMonitor{}
}

func (w *windowsMonitor) GetProcessList() ([]ProcessInfo, error) {
	// Use tasklist command to get processes
	cmd := exec.Command("tasklist", "/FO", "CSV", "/NH")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to execute tasklist: %v", err)
	}

	var processes []ProcessInfo
	for _, line := range strings.Split(out.String(), "\n") {
		// Skip empty lines
		if len(line) == 0 {
			continue
		}

		// Parse CSV format
		fields := strings.Split(strings.Trim(line, "\"\r\n"), "\",\"")
		if len(fields) < 2 {
			continue
		}

		// Get process name and PID
		name := strings.Trim(fields[0], "\"")
		pid, err := strconv.ParseInt(strings.Trim(fields[1], "\""), 10, 32)
		if err != nil {
			continue
		}

		processes = append(processes, ProcessInfo{
			PID:  int32(pid),
			Name: name,
		})
	}

	return processes, nil
}

func (w *windowsMonitor) GetProcessStats(pid int32) (*ProcessStats, error) {
	// Use wmic to get process memory info
	cmd := exec.Command("wmic", "process", "where",
		fmt.Sprintf("ProcessId=%d", pid),
		"get", "WorkingSetSize", "/format:value")

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to get process stats: %v", err)
	}

	// Parse the output
	output := strings.TrimSpace(out.String())
	memoryStr := strings.TrimPrefix(output, "WorkingSetSize=")
	memory, err := strconv.ParseUint(strings.TrimSpace(memoryStr), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse memory value: %v", err)
	}

	return &ProcessStats{
		MemoryUsage: memory,
		CPUUsage:    0, // Would need multiple samples over time to calculate CPU usage
	}, nil
}
