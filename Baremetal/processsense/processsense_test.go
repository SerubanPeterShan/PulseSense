package processsense

import (
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestRealProcessList(t *testing.T) {
	fmt.Printf("\n=== Real Process List Test (%s) ===\n", runtime.GOOS)

	monitor := NewProcessMonitor()
	processes, err := monitor.GetProcessList()
	if err != nil {
		t.Fatalf("Failed to get process list: %v", err)
	}

	fmt.Printf("Total processes found: %d\n", len(processes))

	// Print first 10 processes
	fmt.Println("\nFirst 10 processes:")
	count := 0
	for _, p := range processes {
		if count >= 10 {
			break
		}
		fmt.Printf("PID: %d, Name: %s\n", p.PID, p.Name)
		count++
	}
}

func TestCurrentProcessStats(t *testing.T) {
	fmt.Printf("\n=== Current Process Stats Test (%s) ===\n", runtime.GOOS)

	monitor := NewProcessMonitor()
	currentPID := int32(os.Getpid())

	// Get initial stats
	stats1, err := monitor.GetProcessStats(currentPID)
	if err != nil {
		t.Fatalf("Failed to get process stats: %v", err)
	}

	fmt.Printf("\nCurrent Process (PID: %d)\n", currentPID)
	fmt.Printf("Memory Usage: %.2f MB\n", float64(stats1.MemoryUsage)/(1024*1024))
	fmt.Printf("CPU Usage: %.1f%%\n", stats1.CPUUsage)

	// Validate stats
	if stats1.MemoryUsage == 0 {
		t.Error("Memory usage should not be zero")
	}
}

func TestProcessStatsOverTime(t *testing.T) {
	fmt.Printf("\n=== Process Stats Over Time Test (%s) ===\n", runtime.GOOS)

	monitor := NewProcessMonitor()
	currentPID := int32(os.Getpid())

	// Get stats multiple times
	for i := 0; i < 3; i++ {
		stats, err := monitor.GetProcessStats(currentPID)
		if err != nil {
			t.Fatalf("Failed to get process stats: %v", err)
		}

		fmt.Printf("\nSample %d:\n", i+1)
		fmt.Printf("Memory Usage: %.2f MB\n", float64(stats.MemoryUsage)/(1024*1024))
		fmt.Printf("CPU Usage: %.1f%%\n", stats.CPUUsage)

		time.Sleep(time.Second)
	}
}

func TestInvalidProcess(t *testing.T) {
	fmt.Printf("\n=== Invalid Process Test (%s) ===\n", runtime.GOOS)

	monitor := NewProcessMonitor()

	// Try to get stats for an invalid PID
	invalidPID := int32(999999)
	_, err := monitor.GetProcessStats(invalidPID)

	if err == nil {
		t.Error("Expected error for invalid PID, got none")
	} else {
		fmt.Printf("Expected error received: %v\n", err)
	}
}

func TestSystemProcesses(t *testing.T) {
	fmt.Printf("\n=== System Processes Test (%s) ===\n", runtime.GOOS)

	monitor := NewProcessMonitor()
	processes, err := monitor.GetProcessList()
	if err != nil {
		t.Fatalf("Failed to get process list: %v", err)
	}

	// Track some statistics
	var systemProcesses int
	var userProcesses int
	var largeProcessCount int

	for _, p := range processes {
		stats, err := monitor.GetProcessStats(p.PID)
		if err != nil {
			continue
		}

		if stats.MemoryUsage > 100*1024*1024 { // More than 100MB
			fmt.Printf("Large process: PID=%d, Name=%s, Memory=%.2f MB\n",
				p.PID, p.Name, float64(stats.MemoryUsage)/(1024*1024))
			largeProcessCount++
			if largeProcessCount >= 4 {
				break
			}
		}

		if p.PID < 1000 {
			systemProcesses++
		} else {
			userProcesses++
		}
	}

	fmt.Printf("\nProcess Statistics:\n")
	fmt.Printf("System Processes: %d\n", systemProcesses)
	fmt.Printf("User Processes: %d\n", userProcesses)
}
