package cachesense

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getCacheInfo() (CacheInfo, error) {
	var cacheInfo CacheInfo

	//Read l1d Cache size
	l1d, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cache/index0/size")
	if err != nil {
		return cacheInfo, fmt.Errorf("failed to read L1d cache: %v", err)
	}
	cacheInfo.L1d = parseCacheSize(string(l1d))

	//Read l1i Cache size
	l1i, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cache/index1/size")
	if err != nil {
		return cacheInfo, fmt.Errorf("failed to read L1i cache: %v", err)
	}
	cacheInfo.L1i = parseCacheSize(string(l1i))

	//Read l2 Cache size
	l2, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cache/index2/size")
	if err != nil {
		return cacheInfo, fmt.Errorf("failed to read L2 cache: %v", err)
	}
	cacheInfo.L2 = parseCacheSize(string(l2))

	//Read l3 Cache size
	l3, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cache/index3/size")
	if err != nil {
		return cacheInfo, fmt.Errorf("failed to read L3 cache: %v", err)
	}
	cacheInfo.L3 = parseCacheSize(string(l3))

	return cacheInfo, nil
}

// parseCacheSize converts cache size string (e.g., "32K") to bytes
func parseCacheSize(size string) uint64 {
	size = strings.TrimSpace(size)
	multiplier := uint64(1)

	if strings.HasSuffix(size, "K") {
		multiplier = 1024
		size = strings.TrimSuffix(size, "K")
	} else if strings.HasSuffix(size, "M") {
		multiplier = 1024 * 1024
		size = strings.TrimSuffix(size, "M")
	}

	value, _ := strconv.ParseUint(size, 10, 64)
	return value * multiplier
}
