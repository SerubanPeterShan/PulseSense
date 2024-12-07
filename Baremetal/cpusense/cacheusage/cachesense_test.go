package cachesense

import (
	"os"
	"os/exec"
	"testing"
)

// Mock the os.ReadFile function for testing
var readFile = os.ReadFile

// Mock the exec.Command function for testing
var execCommand = exec.Command

func TestGetCacheInfoLinux(t *testing.T) {
	// Mock the os.ReadFile function
	readFile = func(filename string) ([]byte, error) {
		switch filename {
		case "/sys/devices/system/cpu/cpu0/cache/index0/size":
			return []byte("32K"), nil
		case "/sys/devices/system/cpu/cpu0/cache/index1/size":
			return []byte("32K"), nil
		case "/sys/devices/system/cpu/cpu0/cache/index2/size":
			return []byte("256K"), nil
		case "/sys/devices/system/cpu/cpu0/cache/index3/size":
			return []byte("8192K"), nil
		default:
			return nil, os.ErrNotExist
		}
	}

	cacheInfo, err := getCacheInfoLinux()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cacheInfo.L1d != 32*1024 {
		t.Errorf("Expected L1d cache size to be 32K, got %d", cacheInfo.L1d)
	}
	if cacheInfo.L1i != 32*1024 {
		t.Errorf("Expected L1i cache size to be 32K, got %d", cacheInfo.L1i)
	}
	if cacheInfo.L2 != 256*1024 {
		t.Errorf("Expected L2 cache size to be 256K, got %d", cacheInfo.L2)
	}
	if cacheInfo.L3 != 8192*1024 {
		t.Errorf("Expected L3 cache size to be 8192K, got %d", cacheInfo.L3)
	}
}

func TestGetCacheInfoWindows(t *testing.T) {
	// Mock the exec.Command function
	execCommand = func(_ string, _ ...string) *exec.Cmd {
		cmd := exec.Command("echo", "Node,L2CacheSize,L3CacheSize\n0,256,8192")
		return cmd
	}

	cacheInfo, err := getCacheInfoWindows()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cacheInfo.L1d != 0 {
		t.Errorf("Expected L1d cache size to be 0, got %d", cacheInfo.L1d)
	}
	if cacheInfo.L1i != 0 {
		t.Errorf("Expected L1i cache size to be 0, got %d", cacheInfo.L1i)
	}
	if cacheInfo.L2 != 256*1024 {
		t.Errorf("Expected L2 cache size to be 256K, got %d", cacheInfo.L2)
	}
	if cacheInfo.L3 != 8192*1024 {
		t.Errorf("Expected L3 cache size to be 8192K, got %d", cacheInfo.L3)
	}
}
