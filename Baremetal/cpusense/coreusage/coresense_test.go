package coresense

import (
	"testing"
)

// TestGetNumCPU tests the GetNumCPU function
func TestGetNumCPU(t *testing.T) {
	numCPU := GetNumCPU()
	if numCPU <= 0 {
		t.Errorf("Expected number of CPUs to be greater than 0, got %d", numCPU)
	}
}

// TestGetCPUUsageByOtherServices tests the GetCPUUsageByOtherServices function
func TestGetCPUUsageByOtherServices(t *testing.T) {
	usage, err := GetCPUUsageByOtherServices()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if usage < 0 {
		t.Errorf("Expected CPU usage to be non-negative, got %f", usage)
	}
}
