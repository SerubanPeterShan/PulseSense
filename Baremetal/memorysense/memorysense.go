package memorysense

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

// MemoryStats holds memory information
type MemoryStats struct {
	Total     uint64
	Available uint64
	Used      uint64
	Usage     float64
}

// GetMemoryStats returns current memory statistics
func GetMemoryStats() (MemoryStats, error) {

	switch runtime.GOOS {
	case "windows":
		memStats, err := getMemoryStatsWindows()
		if err != nil {
			return MemoryStats{}, err
		}
		return memStats, nil
	case "linux":
		memStats, err := getMemoryStatsLinux()
		if err != nil {
			return MemoryStats{}, err
		}
		return memStats, nil
	default:
		return MemoryStats{}, fmt.Errorf("unsupported platform")
	}
}

func getMemoryStatsWindows() (MemoryStats, error) {
	var memStats MemoryStats
	var memStatusEx struct {
		Length               uint32
		MemoryLoad           uint32
		TotalPhys            uint64
		AvailPhys            uint64
		TotalPageFile        uint64
		AvailPageFile        uint64
		TotalVirtual         uint64
		AvailVirtual         uint64
		AvailExtendedVirtual uint64
	}

	memStatusEx.Length = uint32(unsafe.Sizeof(memStatusEx))
	modkernel32 := syscall.NewLazyDLL("kernel32.dll")
	procGlobalMemoryStatusEx := modkernel32.NewProc("GlobalMemoryStatusEx")
	ret, _, err := procGlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memStatusEx)))
	if ret == 0 {
		return memStats, err
	}

	memStats.Total = memStatusEx.TotalPhys
	memStats.Available = memStatusEx.AvailPhys
	memStats.Used = memStats.Total - memStats.Available
	memStats.Usage = float64(memStats.Used) / float64(memStats.Total) * 100.0

	return memStats, nil
}

func getMemoryStatsLinux() (MemoryStats, error) {
	var memStats MemoryStats

	content, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return memStats, err
	}

	lines := strings.Split(string(content), "\n")
	memInfo := make(map[string]uint64)

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		key := strings.TrimSuffix(fields[0], ":")
		value := fields[1]

		// Convert value to bytes (values in /proc/meminfo are in KB)
		if val, err := strconv.ParseUint(value, 10, 64); err == nil {
			memInfo[key] = val * 1024 // Convert KB to bytes
		}
	}

	memStats.Total = memInfo["MemTotal"]
	memStats.Available = memInfo["MemAvailable"]
	if memStats.Available == 0 {
		// If MemAvailable is not present, calculate from free and buffers/cache
		memStats.Available = memInfo["MemFree"] + memInfo["Buffers"] + memInfo["Cached"]
	}
	memStats.Used = memStats.Total - memStats.Available
	memStats.Usage = float64(memStats.Used) / float64(memStats.Total) * 100.0

	return memStats, nil
}

// PrintMemoryStats prints current memory statistics
func PrintMemoryStats() error {
	stats, err := GetMemoryStats()
	if err != nil {
		return err
	}

	fmt.Printf("Total Memory: %d MB\n", stats.Total/1024/1024)
	fmt.Printf("Available Memory: %d MB\n", stats.Available/1024/1024)
	fmt.Printf("Used Memory: %d MB\n", stats.Used/1024/1024)
	fmt.Printf("Memory Usage: %.2f%%\n", stats.Usage)

	return nil
}

// MonitorMemory continuously prints memory statistics
func MonitorMemory(interval time.Duration) {
	for {
		err := PrintMemoryStats()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		time.Sleep(interval)
	}
}
