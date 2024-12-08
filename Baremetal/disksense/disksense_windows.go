package disksense

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

func getDiskInfo() ([]DiskInfo, error) {
	var disks []DiskInfo
	driveCount := 0

	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		path := string(drive) + ":\\"

		// Check if drive exists
		if _, err := os.Stat(path); err != nil {
			if os.IsPermission(err) {
				return nil, fmt.Errorf("access denied to path %s: %v", path, err)
			}
			continue
		}

		var freeBytesAvailable, totalBytes, totalFreeBytes uint64
		kernel32 := syscall.NewLazyDLL("kernel32.dll")
		getDiskFreeSpaceEx := kernel32.NewProc("GetDiskFreeSpaceExW")

		pathPtr, err := syscall.UTF16PtrFromString(path)
		if err != nil {
			return nil, fmt.Errorf("failed to convert path %s: %v", path, err)
		}

		ret, _, err := getDiskFreeSpaceEx.Call(
			uintptr(unsafe.Pointer(pathPtr)),
			uintptr(unsafe.Pointer(&freeBytesAvailable)),
			uintptr(unsafe.Pointer(&totalBytes)),
			uintptr(unsafe.Pointer(&totalFreeBytes)),
		)

		if ret == 0 {
			// Check if it's a permission error
			if err == syscall.ERROR_ACCESS_DENIED {
				return nil, fmt.Errorf("access denied to disk %s", path)
			}
			// Skip unreadable drives but continue
			continue
		}

		// Skip drives with 0 total bytes (usually empty card readers)
		if totalBytes == 0 {
			continue
		}

		used := totalBytes - totalFreeBytes
		usage := float64(used) / float64(totalBytes) * 100

		// Get filesystem type
		fsType := "Unknown"
		volInfo, err := getVolumeInformation(path)
		if err == nil {
			fsType = volInfo
		}

		disk := DiskInfo{
			Path:       path,
			Total:      totalBytes,
			Free:       totalFreeBytes,
			Used:       used,
			Usage:      usage,
			MountPoint: path,
			FileSystem: fsType,
		}
		disks = append(disks, disk)
		driveCount++
	}

	if driveCount == 0 {
		return nil, fmt.Errorf("no accessible drives found")
	}

	return disks, nil
}

func getVolumeInformation(path string) (string, error) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getVolumeInformation := kernel32.NewProc("GetVolumeInformationW")

	var fsNameBuffer [256]uint16
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return "", err
	}

	ret, _, err := getVolumeInformation.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		0,
		0,
		0,
		0,
		0,
		uintptr(unsafe.Pointer(&fsNameBuffer[0])),
		256,
	)

	if ret == 0 {
		return "", fmt.Errorf("failed to get volume information: %v", err)
	}

	return syscall.UTF16ToString(fsNameBuffer[:]), nil
}
