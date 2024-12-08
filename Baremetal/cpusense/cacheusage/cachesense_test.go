package cachesense

import (
	"fmt"
	"runtime"
	"testing"
)

func TestGetCacheInfo(t *testing.T) {
	fmt.Printf("\n=== Cache Information Test (%s) ===\n", runtime.GOOS)

	cacheInfo, err := GetCacheInfo()
	if err != nil {
		t.Fatalf("Failed to get cache info: %v", err)
	}

	// Print cache information
	fmt.Printf("L1 Data Cache: %s\n", formatSize(cacheInfo.L1d))
	fmt.Printf("L1 Instruction Cache: %s\n", formatSize(cacheInfo.L1i))
	fmt.Printf("L2 Cache: %s\n", formatSize(cacheInfo.L2))
	fmt.Printf("L3 Cache: %s\n", formatSize(cacheInfo.L3))

	// Validation tests
	testCases := []struct {
		name  string
		size  uint64
		level string
	}{
		{"L1 Data Cache", cacheInfo.L1d, "L1d"},
		{"L1 Instruction Cache", cacheInfo.L1i, "L1i"},
		{"L2 Cache", cacheInfo.L2, "L2"},
		{"L3 Cache", cacheInfo.L3, "L3"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validateCacheSize(t, tc.size, tc.level)
		})
	}
}

func TestPlatformSpecific(t *testing.T) {
	fmt.Printf("\n=== Platform Specific Cache Tests (%s) ===\n", runtime.GOOS)

	cacheInfo, err := GetCacheInfo()
	if err != nil {
		t.Fatalf("Failed to get cache info: %v", err)
	}

	switch runtime.GOOS {
	case "windows":
		testWindowsCache(t, cacheInfo)
	case "linux":
		testLinuxCache(t, cacheInfo)
	}
}

func testWindowsCache(t *testing.T, info CacheInfo) {
	// Windows specific validations
	if info.L2 == 0 {
		t.Error("L2 cache size should not be zero on Windows")
	}
	fmt.Printf("Windows CPU Cache Configuration:\n")
	fmt.Printf("L2 Cache: %s\n", formatSize(info.L2))
	fmt.Printf("L3 Cache: %s\n", formatSize(info.L3))
}

func testLinuxCache(t *testing.T, info CacheInfo) {
	// Linux specific validations
	if info.L1d == 0 {
		t.Error("L1 data cache size should not be zero on Linux")
	}
	fmt.Printf("Linux CPU Cache Configuration:\n")
	fmt.Printf("L1 Data: %s\n", formatSize(info.L1d))
	fmt.Printf("L1 Instruction: %s\n", formatSize(info.L1i))
	fmt.Printf("L2: %s\n", formatSize(info.L2))
	fmt.Printf("L3: %s\n", formatSize(info.L3))
}

func validateCacheSize(t *testing.T, size uint64, level string) {
	// Common cache size ranges (in bytes)
	ranges := map[string]struct{ min, max uint64 }{
		"L1d": {8 * 1024, 64 * 1024},          // 8KB - 64KB
		"L1i": {8 * 1024, 64 * 1024},          // 8KB - 64KB
		"L2":  {128 * 1024, 24 * 1024 * 1024}, // 128KB - 24MB
		"L3":  {512 * 1024, 64 * 1024 * 1024}, // 512KB - 64MB
	}

	if r, ok := ranges[level]; ok && size > 0 {
		if size < r.min || size > r.max {
			t.Errorf("%s cache size %d is outside expected range (%d-%d)",
				level, size, r.min, r.max)
		}
	}
}

// Helper function to format cache sizes
func formatSize(bytes uint64) string {
	if bytes == 0 {
		return "N/A"
	}
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
