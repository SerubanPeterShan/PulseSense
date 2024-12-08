// healthchecksense_test.go
package healthchecksense

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPodHealthCheck(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		isHealthy      bool
		isReady        bool
		expectedStatus int
		expectedBody   string
	}{
		{"Pod Healthy", "/pod/healthz", true, true, http.StatusOK, "pod healthy"},
		{"Pod Unhealthy", "/pod/healthz", false, true, http.StatusServiceUnavailable, "pod unhealthy"},
		{"Pod Ready", "/pod/ready", true, true, http.StatusOK, "pod ready"},
		{"Pod Not Ready", "/pod/ready", true, false, http.StatusServiceUnavailable, "pod not ready"},
	}

	pod := NewPodHealth("test-pod", "default", nil, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pod.isHealthy.Store(tt.isHealthy)
			pod.isReady.Store(tt.isReady)

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/pod/healthz":
					pod.healthz(w, r)
				case "/pod/ready":
					pod.ready(w, r)
				}
			}))
			defer server.Close()

			resp, err := http.Get(server.URL + tt.path)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			body, _ := io.ReadAll(resp.Body)
			if string(body) != tt.expectedBody {
				t.Errorf("Expected body %q, got %q", tt.expectedBody, string(body))
			}
		})
	}
}

func TestNodeHealthCheck(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		isHealthy      bool
		isReady        bool
		expectedStatus int
		expectedBody   string
	}{
		{"Node Healthy", "/node/healthz", true, true, http.StatusOK, "node healthy"},
		{"Node Unhealthy", "/node/healthz", false, true, http.StatusServiceUnavailable, "node unhealthy"},
		{"Node Ready", "/node/ready", true, true, http.StatusOK, "node ready"},
		{"Node Not Ready", "/node/ready", true, false, http.StatusServiceUnavailable, "node not ready"},
	}

	node := NewNodeHealth("test-node", nil, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node.isHealthy.Store(tt.isHealthy)
			node.isReady.Store(tt.isReady)

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/node/healthz":
					node.healthz(w, r)
				case "/node/ready":
					node.ready(w, r)
				}
			}))
			defer server.Close()

			resp, err := http.Get(server.URL + tt.path)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			body, _ := io.ReadAll(resp.Body)
			if string(body) != tt.expectedBody {
				t.Errorf("Expected body %q, got %q", tt.expectedBody, string(body))
			}
		})
	}
}
