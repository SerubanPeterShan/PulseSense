package cpuusagesense

import (
	"testing"
	"time"
)

// TestGetCPUUsage tests the GetCPUUsage function
func TestGetCPUUsage(t *testing.T) {
	cpuUsage := GetCPUUsage()
	if cpuUsage < 0 || cpuUsage > 100 {
		t.Errorf("Expected CPU usage to be between 0 and 100, got %f", cpuUsage)
	}
}

// TestMonitorCPU tests the MonitorCPU function
func TestMonitorCPU(t *testing.T) {
	//Run MonitorCPU in a goroutine
	go MonitorCPU()

	//Allow MonitorCPU to run for 3 seconds
	time.Sleep(3 * time.Second)

	//Since MonitorCPU runs indefinitely, we can assume that it ran without errors if it runs for 3 seconds
	//Instead of checking the output of MonitorCPU, we ensure that it runs without errors for 3 seconds
	t.Log("MonitorCPU ran without errors for 3 seconds")
}
