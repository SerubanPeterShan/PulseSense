package disksense

import (
	"os"
	"runtime"
	"testing"
	"time"
)

func TestGetDiskInfo(t *testing.T) {
	disks, err := GetDiskInfo()
	if err != nil {
		t.Fatalf("GetDiskInfo failed: %v", err)
	}

	if len(disks) == 0 {
		t.Error("Expected at least one disk, got none")
	}

	for _, disk := range disks {
		// Test disk path
		if disk.Path == "" {
			t.Error("Disk path should not be empty")
		}

		// Test mount point
		if disk.MountPoint == "" {
			t.Error("Mount point should not be empty")
		}

		// Test total space
		if disk.Total == 0 {
			t.Error("Total disk space should not be zero")
		}

		// Test usage calculations
		if disk.Usage < 0 || disk.Usage > 100 {
			t.Errorf("Disk usage should be between 0 and 100, got %f", disk.Usage)
		}

		// Test used space
		if disk.Used > disk.Total {
			t.Error("Used space should not exceed total space")
		}

		// Test free space
		if disk.Free > disk.Total {
			t.Error("Free space should not exceed total space")
		}
	}
}

func TestPlatformSpecific(t *testing.T) {
	if runtime.GOOS == "windows" {
		testWindowsDiskInfo(t)
	} else {
		testLinuxDiskInfo(t)
	}
}

func testWindowsDiskInfo(t *testing.T) {
	disks, err := getWindowsDiskInfo()
	if err != nil {
		t.Fatalf("getWindowsDiskInfo failed: %v", err)
	}

	for _, disk := range disks {
		// Test Windows-specific path format
		if len(disk.Path) < 3 || disk.Path[1:3] != ":\\" {
			t.Errorf("Invalid Windows path format: %s", disk.Path)
		}
	}
}

func testLinuxDiskInfo(t *testing.T) {
	disks, err := getLinuxDiskInfo()
	if err != nil {
		t.Fatalf("getLinuxDiskInfo failed: %v", err)
	}

	for _, disk := range disks {
		// Test Linux-specific path format
		if disk.Path[0] != '/' {
			t.Errorf("Invalid Linux path format: %s", disk.Path)
		}
	}
}

func TestGetDiskInfoSingle(t *testing.T) {
	// Test with current directory
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	disk, err := getDiskInfo(currentDir)
	if err != nil {
		t.Fatalf("getDiskInfo failed: %v", err)
	}

	// Test disk info
	if disk.Path == "" {
		t.Error("Disk path should not be empty")
	}
	if disk.Total == 0 {
		t.Error("Total disk space should not be zero")
	}
	if disk.Usage < 0 || disk.Usage > 100 {
		t.Errorf("Disk usage should be between 0 and 100, got %f", disk.Usage)
	}
}

func TestMonitorDiskInfo(t *testing.T) {
	// Create a channel to stop the monitor after a short time
	done := make(chan bool)

	go func() {
		// Stop the monitor after 100ms
		time.Sleep(100 * time.Millisecond)
		done <- true
	}()

	go func() {
		MonitorDiskInfo(50 * time.Millisecond)
	}()

	<-done
}
