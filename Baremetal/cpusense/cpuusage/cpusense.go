package cpuusagesense

import (
	"fmt"
	"runtime"
	"time"
)

func GetCPUUsage() float64 {
	var cpuUsage float64

	cpuUsage = float64(runtime.NumGoroutine()) / float64(runtime.NumCPU()) * 100
	return cpuUsage
}

func MonitorCPU() {
	for {
		cpuUsage := GetCPUUsage()
		fmt.Printf("CPU Usage: %f\n", cpuUsage)
		time.Sleep(1 * time.Second)
	}
}
