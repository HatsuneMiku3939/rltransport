package test

import (
	"net/http"
	"net/http/httptest"
)

type MoniteringTripper struct {
	TripCount int
}

func NewMoniteringTripper() *MoniteringTripper {
	return &MoniteringTripper{
		TripCount: 0,
	}
}

func (t *MoniteringTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	w.WriteHeader(200)
	w.Write([]byte("OK"))
	t.TripCount++

	return w.Result(), nil
}

func (t *MoniteringTripper) Reset() {
	t.TripCount = 0
}
