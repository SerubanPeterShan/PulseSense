package processsense

type ProcessInfo struct {
	PID  int32
	Name string
}

type ProcessStats struct {
	MemoryUsage uint64  // in bytes
	CPUUsage    float64 // percentage
}

// Interface declarations
type ProcessMonitor interface {
	GetProcessList() ([]ProcessInfo, error)
	GetProcessStats(pid int32) (*ProcessStats, error)
}
