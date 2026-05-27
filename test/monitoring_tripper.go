package test

import (
	"net/http"
	"net/http/httptest"
	"sync"
)

type MonitoringTripper struct {
	mu        sync.Mutex
	tripCount int
}

func NewMonitoringTripper() *MonitoringTripper {
	return &MonitoringTripper{}
}

func (t *MonitoringTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	w.WriteHeader(200)
	w.Write([]byte("OK"))

	t.mu.Lock()
	t.tripCount++
	t.mu.Unlock()

	return w.Result(), nil
}

func (t *MonitoringTripper) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.tripCount = 0
}

func (t *MonitoringTripper) Count() int {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.tripCount
}
