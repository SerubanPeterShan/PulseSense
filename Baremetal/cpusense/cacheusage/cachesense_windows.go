package cachesense

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func getCacheInfo() (CacheInfo, error) {
	var cacheInfo CacheInfo

	// Get cache info using wmic
	cmd := exec.Command("wmic", "cpu", "get", "L2CacheSize,L3CacheSize", "/format:list")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return CacheInfo{}, fmt.Errorf("failed to execute wmic command: %v", err)
	}

	// Parse the output which will be in format:
	// L2CacheSize=256
	// L3CacheSize=6144
	output := strings.TrimSpace(out.String())
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "L2CacheSize":
			if size, err := strconv.ParseUint(value, 10, 64); err == nil {
				cacheInfo.L2 = size * 1024 // Convert KB to bytes
			}
		case "L3CacheSize":
			if size, err := strconv.ParseUint(value, 10, 64); err == nil {
				cacheInfo.L3 = size * 1024 // Convert KB to bytes
			}
		}
	}

	// Validate that we got at least some cache information
	if cacheInfo.L2 == 0 && cacheInfo.L3 == 0 {
		return CacheInfo{}, fmt.Errorf("no cache information found in wmic output")
	}

	// Note: Windows doesn't provide L1 cache info through wmic
	cacheInfo.L1d = 0
	cacheInfo.L1i = 0

	return cacheInfo, nil
}
