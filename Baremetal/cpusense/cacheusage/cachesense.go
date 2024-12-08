package cachesense

// CacheInfo holds information about CPU cache sizes
type CacheInfo struct {
	L1d uint64 // L1 Data Cache size in bytes
	L1i uint64 // L1 Instruction Cache size in bytes
	L2  uint64 // L2 Cache size in bytes
	L3  uint64 // L3 Cache size in bytes
}

// GetCacheInfo returns cache sizes for the current platform
func GetCacheInfo() (CacheInfo, error) {
	return getCacheInfo()
}
