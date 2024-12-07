package memorysense

import (
	"runtime"
	"testing"
	"time"
)

func TestGetMemoryStats(t *testing.T) {
	stats, err := GetMemoryStats()
	if err != nil {
		t.Fatalf("Failed to get memory stats: %v", err)
	}

	// Test that Total memory is greater than zero
	if stats.Total == 0 {
		t.Error("Total memory should not be zero")
	}

	// Test that Available memory is not greater than Total
	if stats.Available > stats.Total {
		t.Error("Available memory should not exceed total memory")
	}

	// Test that Used memory is not greater than Total
	if stats.Used > stats.Total {
		t.Error("Used memory should not exceed total memory")
	}

	// Test that Usage is between 0 and 100
	if stats.Usage < 0 || stats.Usage > 100 {
		t.Errorf("Memory usage should be between 0 and 100, got %f", stats.Usage)
	}
}

func TestUnsupportedPlatform(t *testing.T) {
	// Skip test if running on Windows or Linux
	if runtime.GOOS == "windows" || runtime.GOOS == "linux" {
		t.Skip("Skipping test on supported platform")
	}

	_, err := GetMemoryStats()
	if err == nil {
		t.Error("Expected error for unsupported platform")
	}
}

func TestPrintMemoryStats(t *testing.T) {
	err := PrintMemoryStats()
	if err != nil {
		t.Errorf("PrintMemoryStats failed: %v", err)
	}
}

func TestMonitorMemory(t *testing.T) {
	// Create a channel to stop the monitor after a short time
	done := make(chan bool)

	go func() {
		// Stop the monitor after 100ms
		time.Sleep(100 * time.Millisecond)
		done <- true
	}()

	go func() {
		MonitorMemory(50 * time.Millisecond)
	}()

	<-done
}
