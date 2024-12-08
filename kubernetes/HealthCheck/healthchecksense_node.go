package healthchecksense

import "net/http"

type NodeHealth struct {
	HealthCheckSense
	nodeName     string
	nodeIP       string
	nodeCapacity map[string]string
}

func NewNodeHealth(nodeName string, livenessProbe, readynessProbe *ProbeConfig) *NodeHealth {
	nh := &NodeHealth{
		nodeName:     nodeName,
		nodeCapacity: make(map[string]string),
	}
	nh.isReady.Store(true)
	nh.isHealthy.Store(true)

	if livenessProbe == nil {
		nh.livenessProbe = DefaultProbeConfig("/node/healthz")
	} else {
		nh.livenessProbe = *livenessProbe
	}

	if readynessProbe == nil {
		nh.readinessProbe = DefaultProbeConfig("/node/ready")
	} else {
		nh.readinessProbe = *readynessProbe
	}

	return nh
}

func (nh *NodeHealth) Start(port string) error {
	http.HandleFunc(nh.livenessProbe.Path, nh.healthz)
	http.HandleFunc(nh.readinessProbe.Path, nh.ready)
	return http.ListenAndServe(":"+port, nil)
}

func (ph *NodeHealth) healthz(w http.ResponseWriter, r *http.Request) {
	if ph.isHealthy.Load() {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("node healthy"))
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("node unhealthy"))
	}
}

func (ph *NodeHealth) ready(w http.ResponseWriter, r *http.Request) {
	if ph.isReady.Load() {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("node ready"))
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("node not ready"))
	}
}
