package disksense

import (
	"fmt"
)

// DriveType represents the type of storage device
type DriveType int

const (
	Unknown DriveType = iota
	HDD
	SSD
	Network
	Removable
	OpticalDrive
)

// DiskStatus represents the health status of the disk
type DiskStatus struct {
	IsHealthy bool
	ReadOnly  bool
	Error     string
}

// DiskInfo holds information about a disk partition
type DiskInfo struct {
	Path       string
	Label      string    // Disk label/name
	DriveType  DriveType // Type of drive (SSD/HDD/etc)
	Total      uint64    // Total size in bytes
	Free       uint64    // Free space in bytes
	Used       uint64    // Used space in bytes
	Usage      float64   // Usage percentage
	MountPoint string    // Mount point path
	FileSystem string    // Filesystem type
	Status     DiskStatus
}

// ByteSize formats bytes into human readable string
func (d DiskInfo) ByteSize(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(bytes)/float64(div), "KMGTPE"[exp])
}

// TotalReadable returns human readable total size
func (d DiskInfo) TotalReadable() string {
	return d.ByteSize(d.Total)
}

// FreeReadable returns human readable free space
func (d DiskInfo) FreeReadable() string {
	return d.ByteSize(d.Free)
}

// UsedReadable returns human readable used space
func (d DiskInfo) UsedReadable() string {
	return d.ByteSize(d.Used)
}

// IsLowSpace returns true if available space is less than threshold percent
func (d DiskInfo) IsLowSpace(threshold float64) bool {
	return d.Usage > (100 - threshold)
}

// GetDiskInfo returns disk information based on OS
func GetDiskInfo() ([]DiskInfo, error) {
	return getDiskInfo()
}

// GetDiskByPath returns disk information for a specific path
func GetDiskByPath(path string) (*DiskInfo, error) {
	disks, err := getDiskInfo()
	if err != nil {
		return nil, err
	}

	for _, disk := range disks {
		if disk.Path == path {
			return &disk, nil
		}
	}
	return nil, fmt.Errorf("disk not found: %s", path)
}
