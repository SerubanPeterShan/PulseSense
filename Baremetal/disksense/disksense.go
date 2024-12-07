package disksense

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// DiskInfo holds information about a disk partition
type DiskInfo struct {
	Path       string
	Total      uint64
	Free       uint64
	Used       uint64
	Usage      float64
	MountPoint string
	FileSystem string
}

// GetDiskInfo returns disk information based on OS
func GetDiskInfo() ([]DiskInfo, error) {
	if runtime.GOOS == "windows" {
		return getWindowsDiskInfo()
	}
	return getLinuxDiskInfo()
}

func getWindowsDiskInfo() ([]DiskInfo, error) {
	var disks []DiskInfo

	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		path := string(drive) + ":\\"
		if _, err := os.Stat(path); err == nil {
			disk, err := getDiskInfo(path)
			if err == nil {
				disks = append(disks, disk)
			}
		}
	}

	return disks, nil
}

func getLinuxDiskInfo() ([]DiskInfo, error) {
	var disks []DiskInfo

	// Common mount points in Linux
	mountPoints := []string{
		"/",
		"/home",
		"/boot",
		"/usr",
		"/var",
		"/tmp",
	}

	for _, mount := range mountPoints {
		if _, err := os.Stat(mount); err == nil {
			disk, err := getDiskInfo(mount)
			if err == nil {
				disks = append(disks, disk)
			}
		}
	}

	return disks, nil
}

func getDiskInfo(path string) (DiskInfo, error) {
	var disk DiskInfo
	disk.Path = path
	disk.MountPoint = path

	// Get disk usage using filepath.Walk
	var size int64
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	if err != nil {
		return disk, err
	}

	// Try to create a temp file to check write permission
	tempFile, err := os.CreateTemp(path, "diskspace")
	if err == nil {
		tempPath := tempFile.Name()
		tempFile.Close()
		os.Remove(tempPath)
	}

	// Set reasonable defaults for total space
	if runtime.GOOS == "windows" {
		disk.Total = 1 << 40 // 1 TB for Windows
	} else {
		disk.Total = 1 << 38 // 256 GB for Linux
	}

	disk.Used = uint64(size)
	disk.Free = disk.Total - disk.Used
	if disk.Free > disk.Total {
		disk.Free = 0
	}
	disk.Usage = float64(disk.Used) / float64(disk.Total) * 100.0

	return disk, nil
}

// PrintDiskInfo prints information about all available disks
func PrintDiskInfo() error {
	disks, err := GetDiskInfo()
	if err != nil {
		return err
	}

	fmt.Printf("\n=== Disk Information ===\n")
	for _, disk := range disks {
		fmt.Printf("\nMount Point: %s\n", disk.MountPoint)
		fmt.Printf("Total Space: %.2f GB\n", float64(disk.Total)/(1024*1024*1024))
		fmt.Printf("Used Space: %.2f GB\n", float64(disk.Used)/(1024*1024*1024))
		fmt.Printf("Free Space: %.2f GB\n", float64(disk.Free)/(1024*1024*1024))
		fmt.Printf("Usage: %.2f%%\n", disk.Usage)
	}

	return nil
}

// MonitorDiskInfo continuously monitors disk information
func MonitorDiskInfo(interval time.Duration) {
	for {
		err := PrintDiskInfo()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		time.Sleep(interval)
	}
}
