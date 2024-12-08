package healthchecksense

import "net/http"

type PodHealth struct {
	HealthCheckSense
	podName      string
	namespace    string
	containerIDs []string
}

func NewPodHealth(podName, namespace string, livenessProbe, readynessProbe *ProbeConfig) *PodHealth {
	ph := &PodHealth{
		podName:   podName,
		namespace: namespace,
	}
	ph.isReady.Store(true)
	ph.isHealthy.Store(true)

	if livenessProbe == nil {
		ph.livenessProbe = DefaultProbeConfig("/pod/healthz")
	} else {
		ph.livenessProbe = *livenessProbe
	}

	if readynessProbe == nil {
		ph.readinessProbe = DefaultProbeConfig("/pod/ready")
	} else {
		ph.readinessProbe = *readynessProbe
	}

	return ph
}

func (ph *PodHealth) Start(port string) error {
	http.HandleFunc(ph.livenessProbe.Path, ph.healthz)
	http.HandleFunc(ph.readinessProbe.Path, ph.ready)
	return http.ListenAndServe(":"+port, nil)
}

func (ph *PodHealth) healthz(w http.ResponseWriter, r *http.Request) {
	if ph.isHealthy.Load() {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pod healthy"))
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("pod unhealthy"))
	}
}

func (ph *PodHealth) ready(w http.ResponseWriter, r *http.Request) {
	if ph.isReady.Load() {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pod ready"))
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("pod not ready"))
	}
}
