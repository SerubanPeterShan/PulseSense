package cachesense

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type CacheInfo struct {
	L1d int
	L1i int
	L2  int
	L3  int
}

func GetCacheInfo() (CacheInfo, error) {
	if runtime.GOOS == "windows" {
		return getCacheInfoWindows()
	}
	return getCacheInfoLinux()
}

func getCacheInfoLinux() (CacheInfo, error) {
	var cacheInfo CacheInfo

	//Read l1d Cache size
	l1d, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cache/index0/size")
	if err != nil {
		return cacheInfo, err
	}
	cacheInfo.L1d = parseCacheSize(string(l1d))

	//Read l1i Cache size
	l1i, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cache/index1/size")
	if err != nil {
		return cacheInfo, err
	}
	cacheInfo.L1i = parseCacheSize(string(l1i))

	//Read l2 Cache size
	l2, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cache/index2/size")
	if err != nil {
		return cacheInfo, err
	}
	cacheInfo.L2 = parseCacheSize(string(l2))

	//Read l3 Cache size
	l3, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cache/index3/size")
	if err != nil {
		return cacheInfo, err
	}
	cacheInfo.L3 = parseCacheSize(string(l3))

	return cacheInfo, nil
}

// getCacheInfoWindows returns the sizes of L1, L2, and L3 caches on Windows
func getCacheInfoWindows() (CacheInfo, error) {
	var cacheInfo CacheInfo

	cmd := exec.Command("cmd", "/C", "wmic cpu get L2CacheSize, L3CacheSize /format:csv")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return CacheInfo{}, err
	}

	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "L2CacheSize") || strings.Contains(line, "L3CacheSize") {
			continue
		}
		fields := strings.Split(line, ",")
		if len(fields) >= 2 {
			cacheInfo.L2 = parseCacheSize(strings.TrimSpace(fields[1]) + "K")
			cacheInfo.L3 = parseCacheSize(strings.TrimSpace(fields[2]) + "K")
		}
	}

	if err := scanner.Err(); err != nil {
		return CacheInfo{}, err
	}

	// L1 cache sizes are not directly available via wmic, so we set them to 0
	cacheInfo.L1d = 0
	cacheInfo.L1i = 0

	return cacheInfo, nil
}

// parseCacheSize parses the cache size string and converts it to bytes
func parseCacheSize(sizeStr string) int {
	sizeStr = strings.TrimSpace(sizeStr)
	if strings.HasSuffix(sizeStr, "K") {
		sizeStr = strings.TrimSuffix(sizeStr, "K")
		size, _ := strconv.Atoi(sizeStr)
		return size * 1024 // Convert KB to bytes
	}
	size, _ := strconv.Atoi(sizeStr)
	return size
}
