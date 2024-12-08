package healthchecksense

import (
	"sync/atomic"
)

type ProbeConfig struct {
	Path                string
	InitialDelaySeconds int
	PeriodSeconds       int
	TimeoutSeconds      int
	FailureThreshold    int
	SuccessThreshold    int
}

type HealthCheckSense struct {
	isReady        atomic.Bool
	isHealthy      atomic.Bool
	livenessProbe  ProbeConfig
	readinessProbe ProbeConfig
}

// DefaultProbeConfig returns default probe configuration
func DefaultProbeConfig(path string) ProbeConfig {
	return ProbeConfig{
		Path:                path,
		InitialDelaySeconds: 3,
		PeriodSeconds:       10,
		TimeoutSeconds:      1,
		FailureThreshold:    3,
		SuccessThreshold:    1,
	}
}
