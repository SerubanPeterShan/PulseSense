package memorysense

// MemoryStats holds memory information
type MemoryStats struct {
	Total     uint64  // Total physical memory in bytes
	Available uint64  // Available memory in bytes
	Used      uint64  // Used memory in bytes
	Usage     float64 // Memory usage percentage
}

// GetMemoryStats returns current memory statistics
func GetMemoryStats() (MemoryStats, error) {
	return getMemoryStats()
}
