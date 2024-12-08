package disksense

import (
	"fmt"
	"math"
	"runtime"
	"strings"
	"testing"
)

func TestGetDiskInfo(t *testing.T) {
	fmt.Printf("\n=== Disk Information Test (%s) ===\n", runtime.GOOS)

	disks, err := GetDiskInfo()
	if err != nil {
		t.Fatalf("GetDiskInfo failed: %v", err)
	}

	fmt.Printf("Found %d disk(s)\n", len(disks))

	if len(disks) == 0 {
		t.Fatal("Expected at least one disk, got none")
	}

	for i, disk := range disks {
		fmt.Printf("\nDisk %d:\n", i+1)
		fmt.Printf("Path: %s\n", disk.Path)
		fmt.Printf("Mount Point: %s\n", disk.MountPoint)
		fmt.Printf("Filesystem: %s\n", disk.FileSystem)
		fmt.Printf("Total: %s\n", disk.TotalReadable())
		fmt.Printf("Used: %s (%.1f%%)\n", disk.UsedReadable(), disk.Usage)
		fmt.Printf("Free: %s\n", disk.FreeReadable())

		t.Run(fmt.Sprintf("Disk_%d", i), func(t *testing.T) {
			// Basic validation
			if disk.Path == "" {
				t.Error("Disk path should not be empty")
			}
			if disk.MountPoint == "" {
				t.Error("Mount point should not be empty")
			}
			if disk.FileSystem == "" {
				t.Error("Filesystem type should not be empty")
			}

			// Space validation
			if disk.Total == 0 {
				t.Error("Total disk space should not be zero")
			}
			if disk.Usage < 0 || disk.Usage > 100 {
				t.Errorf("Disk usage should be between 0 and 100, got %.2f", disk.Usage)
			}
			if disk.Used > disk.Total {
				t.Error("Used space should not exceed total space")
			}
			if disk.Free > disk.Total {
				t.Error("Free space should not exceed total space")
			}

			// Check that Used + Free equals Total (within 1% margin for rounding)
			totalFromParts := disk.Used + disk.Free
			diffPercentage := math.Abs(float64(totalFromParts-disk.Total)) / float64(disk.Total) * 100
			if diffPercentage > 1.0 {
				t.Errorf("Space calculation mismatch: Total=%d, Used+Free=%d (%.2f%% difference)",
					disk.Total, totalFromParts, diffPercentage)
			}
		})
	}
}

func TestPlatformSpecific(t *testing.T) {
	fmt.Printf("\n=== Platform Specific Tests (%s) ===\n", runtime.GOOS)

	disks, err := GetDiskInfo()
	if err != nil {
		t.Fatalf("GetDiskInfo failed: %v", err)
	}

	switch runtime.GOOS {
	case "windows":
		testWindowsDisks(t, disks)
	case "linux":
		testLinuxDisks(t, disks)
	}
}

func testWindowsDisks(t *testing.T, disks []DiskInfo) {
	for _, disk := range disks {
		// Test Windows path format
		if len(disk.Path) < 2 || disk.Path[1] != ':' {
			t.Errorf("Invalid Windows path format: %s", disk.Path)
		}

		// Test filesystem type
		validFS := map[string]bool{"NTFS": true, "FAT32": true, "FAT": true, "exFAT": true}
		if !validFS[disk.FileSystem] {
			fmt.Printf("Note: Unusual Windows filesystem: %s\n", disk.FileSystem)
		}
	}
}

func testLinuxDisks(t *testing.T, disks []DiskInfo) {
	for _, disk := range disks {
		// Test Linux path format
		if !strings.HasPrefix(disk.Path, "/dev/") && !strings.HasPrefix(disk.Path, "UUID=") {
			t.Errorf("Invalid Linux device path: %s", disk.Path)
		}

		// Test mount point format
		if !strings.HasPrefix(disk.MountPoint, "/") {
			t.Errorf("Invalid Linux mount point: %s", disk.MountPoint)
		}

		// Test filesystem type
		validFS := map[string]bool{
			"ext4": true, "ext3": true, "ext2": true,
			"xfs": true, "btrfs": true, "zfs": true,
		}
		if !validFS[disk.FileSystem] {
			fmt.Printf("Note: Unusual Linux filesystem: %s\n", disk.FileSystem)
		}
	}
}

func TestLowSpaceDetection(t *testing.T) {
	fmt.Printf("\n=== Low Space Detection Test ===\n")

	disks, err := GetDiskInfo()
	if err != nil {
		t.Fatalf("GetDiskInfo failed: %v", err)
	}

	threshold := 90.0 // 90% usage threshold
	for _, disk := range disks {
		isLow := disk.IsLowSpace(threshold)
		fmt.Printf("Disk %s: %.1f%% used - Low Space: %v\n",
			disk.Path, disk.Usage, isLow)
	}
}
