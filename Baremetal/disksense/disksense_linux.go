package disksense

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
)

func getDiskInfo() ([]DiskInfo, error) {
	var disks []DiskInfo
	mountCount := 0

	// Read /proc/mounts for filesystem info
	file, err := os.Open("/proc/mounts")
	if err != nil {
		return nil, fmt.Errorf("failed to read mount points: %v", err)
	}
	defer file.Close()

	// Skip special filesystems
	skipFS := map[string]bool{
		"proc":       true,
		"sysfs":      true,
		"devpts":     true,
		"devtmpfs":   true,
		"tmpfs":      true,
		"securityfs": true,
		"pstore":     true,
		"debugfs":    true,
		"cgroup":     true,
		"cgroup2":    true,
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 6 {
			continue
		}

		device := fields[0]
		mountPoint := fields[1]
		fsType := fields[2]

		// Skip special filesystems
		if skipFS[fsType] {
			continue
		}

		// Skip bind mounts and virtual filesystems
		if strings.HasPrefix(device, "/dev/") || strings.HasPrefix(device, "UUID=") {
			var stat syscall.Statfs_t
			err := syscall.Statfs(mountPoint, &stat)
			if err != nil {
				if os.IsPermission(err) {
					return nil, fmt.Errorf("access denied to mount point %s: %v", mountPoint, err)
				}
				continue
			}

			// Skip if total size is 0
			total := uint64(stat.Blocks) * uint64(stat.Bsize)
			if total == 0 {
				continue
			}

			free := uint64(stat.Bfree) * uint64(stat.Bsize)
			available := uint64(stat.Bavail) * uint64(stat.Bsize)
			used := total - free
			usage := float64(used) / float64(total) * 100

			disk := DiskInfo{
				Path:       device,
				Total:      total,
				Free:       available, // Use available instead of free for better accuracy
				Used:       used,
				Usage:      usage,
				MountPoint: mountPoint,
				FileSystem: fsType,
			}
			disks = append(disks, disk)
			mountCount++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading mount points: %v", err)
	}

	if mountCount == 0 {
		return nil, fmt.Errorf("no valid mount points found")
	}

	return disks, nil
}

// Helper function to check if filesystem is remote
func isRemoteFS(fsType string) bool {
	remoteFS := map[string]bool{
		"nfs":   true,
		"cifs":  true,
		"smbfs": true,
		"ncpfs": true,
		"sshfs": true,
	}
	return remoteFS[fsType]
}
