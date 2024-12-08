package memorysense

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func getMemoryStats() (MemoryStats, error) {
	cmd := exec.Command("wmic", "OS", "get", "TotalVisibleMemorySize,FreePhysicalMemory", "/Value")
	output, err := cmd.Output()
	if err != nil {
		return MemoryStats{}, fmt.Errorf("failed to execute wmic: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	var total, free uint64

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "TotalVisibleMemorySize=") {
			val := strings.TrimPrefix(line, "TotalVisibleMemorySize=")
			total, _ = strconv.ParseUint(strings.TrimSpace(val), 10, 64)
		} else if strings.HasPrefix(line, "FreePhysicalMemory=") {
			val := strings.TrimPrefix(line, "FreePhysicalMemory=")
			free, _ = strconv.ParseUint(strings.TrimSpace(val), 10, 64)
		}
	}

	// Convert KB to bytes
	total *= 1024
	free *= 1024
	used := total - free

	return MemoryStats{
		Total:     total,
		Available: free,
		Used:      used,
		Usage:     float64(used) / float64(total) * 100,
	}, nil
}
