package test

import (
	"net/http"
	"net/http/httptest"
)

type MonitoringTripper struct {
	TripCount int
}

func NewMonitoringTripper() *MonitoringTripper {
	return &MonitoringTripper{
		TripCount: 0,
	}
}

func (t *MonitoringTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	w.WriteHeader(200)
	w.Write([]byte("OK"))
	t.TripCount++

	return w.Result(), nil
}

func (t *MonitoringTripper) Reset() {
	t.TripCount = 0
}
