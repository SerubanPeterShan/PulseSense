//go:build linux

package processsense

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type linuxMonitor struct{}

func NewProcessMonitor() ProcessMonitor {
	return &linuxMonitor{}
}

func (l *linuxMonitor) GetProcessList() ([]ProcessInfo, error) {
	files, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	var processes []ProcessInfo
	for _, f := range files {
		if pid, err := strconv.ParseInt(f.Name(), 10, 32); err == nil {
			cmdline, err := os.ReadFile(filepath.Join("/proc", f.Name(), "comm"))
			if err != nil {
				continue
			}

			processes = append(processes, ProcessInfo{
				PID:  int32(pid),
				Name: strings.TrimSpace(string(cmdline)),
			})
		}
	}
	return processes, nil
}

func (l *linuxMonitor) GetProcessStats(pid int32) (*ProcessStats, error) {
	statm, err := os.ReadFile(fmt.Sprintf("/proc/%d/statm", pid))
	if err != nil {
		return nil, err
	}

	fields := strings.Fields(string(statm))
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid statm format")
	}

	rss, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return nil, err
	}

	// Convert pages to bytes (multiply by page size, typically 4KB)
	memoryUsage := rss * 4096

	return &ProcessStats{
		MemoryUsage: memoryUsage,
		CPUUsage:    0, // Requires sampling over time
	}, nil
}
