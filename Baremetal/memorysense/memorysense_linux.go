package memorysense

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getMemoryStats() (MemoryStats, error) {
	content, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return MemoryStats{}, fmt.Errorf("failed to read /proc/meminfo: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	memInfo := make(map[string]uint64)

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		key := strings.TrimSuffix(fields[0], ":")
		value, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}
		memInfo[key] = value * 1024
	}

	total := memInfo["MemTotal"]
	available := memInfo["MemAvailable"]
	if available == 0 {
		available = memInfo["MemFree"] + memInfo["Buffers"] + memInfo["Cached"]
	}
	used := total - available

	return MemoryStats{
		Total:     total,
		Available: available,
		Used:      used,
		Usage:     float64(used) / float64(total) * 100,
	}, nil
}
