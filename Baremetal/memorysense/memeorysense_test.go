package memorysense

import (
	"fmt"
	"runtime"
	"testing"
)

func TestMemoryStats(t *testing.T) {
	fmt.Printf("\n=== Memory Statistics Test ===\n")

	stats, err := GetMemoryStats()
	if err != nil {
		t.Fatalf("Failed to get memory stats: %v", err)
	}

	// Print current stats
	fmt.Printf("Total Memory: %.2f GB\n", float64(stats.Total)/(1024*1024*1024))
	fmt.Printf("Available Memory: %.2f GB\n", float64(stats.Available)/(1024*1024*1024))
	fmt.Printf("Used Memory: %.2f GB\n", float64(stats.Used)/(1024*1024*1024))
	fmt.Printf("Memory Usage: %.1f%%\n", stats.Usage)

	// Validation tests
	testCases := []struct {
		name    string
		test    func() bool
		message string
	}{
		{
			name:    "Total memory check",
			test:    func() bool { return stats.Total > 0 },
			message: "Total memory should be greater than zero",
		},
		{
			name:    "Available memory check",
			test:    func() bool { return stats.Available <= stats.Total },
			message: "Available memory should not exceed total memory",
		},
		{
			name:    "Used memory check",
			test:    func() bool { return stats.Used <= stats.Total },
			message: "Used memory should not exceed total memory",
		},
		{
			name:    "Usage percentage check",
			test:    func() bool { return stats.Usage >= 0 && stats.Usage <= 100 },
			message: "Memory usage should be between 0 and 100",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.test() {
				t.Error(tc.message)
			}
		})
	}
}

func TestPlatformSpecific(t *testing.T) {
	fmt.Printf("\n=== Platform Specific Test ===\n")
	fmt.Printf("Testing on: %s\n", runtime.GOOS)

	stats, err := GetMemoryStats()
	if err != nil {
		t.Fatalf("Failed to get memory stats: %v", err)
	}

	switch runtime.GOOS {
	case "windows":
		// Windows specific tests
		if stats.Total < 1024*1024*1024 { // Less than 1GB
			t.Error("Windows system should have at least 1GB memory")
		}
	case "linux":
		// Linux specific tests
		if stats.Total < 1024*1024 { // Less than 1MB
			t.Error("Linux system should have at least 1MB memory")
		}
	}
}

func TestUnsupportedPlatform(t *testing.T) {
	if runtime.GOOS == "windows" || runtime.GOOS == "linux" {
		t.Skip("Skipping test on supported platform")
	}

	_, err := GetMemoryStats()
	if err == nil {
		t.Error("Expected error for unsupported platform")
	}
}

func BenchmarkGetMemoryStats(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GetMemoryStats()
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}

// Helper function to print memory stats
func TestMemoryStatsOutput(t *testing.T) {
	fmt.Printf("\n=== Memory Stats Output Test ===\n")

	stats, err := GetMemoryStats()
	if err != nil {
		t.Fatalf("Failed to get memory stats: %v", err)
	}

	// Print with units
	units := []struct {
		name string
		div  float64
	}{
		{"GB", 1024 * 1024 * 1024},
		{"MB", 1024 * 1024},
		{"KB", 1024},
	}

	for _, unit := range units {
		fmt.Printf("\nMemory in %s:\n", unit.name)
		fmt.Printf("Total: %.2f %s\n", float64(stats.Total)/unit.div, unit.name)
		fmt.Printf("Available: %.2f %s\n", float64(stats.Available)/unit.div, unit.name)
		fmt.Printf("Used: %.2f %s\n", float64(stats.Used)/unit.div, unit.name)
	}
}
