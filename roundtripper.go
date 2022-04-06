package rltransport

import (
	"net/http"
	"sync"
)

// RoundTripper implements the http.RoundTripper interface.
type RoundTripper struct {
	// once ensures that the logic to initialize the default client runs at
	// most once, in a single thread.
	once sync.Once

	// Limiter is used to rate limit the number of requests that can be made
	// to the underlying client.
	Limiter Limiter

	// Transport is the underlying RoundTripper that will be used to make
	// the actual HTTP requests.
	Transport http.RoundTripper
}

// RoundTrip satisfies the http.RoundTripper interface.
func (rt *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Ensure that the Transport is initialized.
	rt.once.Do(rt.init)

	// Wait for the rate limiter  within context of the request (if limiter is given).
	if err := rt.Limiter.Wait(req.Context()); err != nil {
		return nil, err
	}

	// Execute the request.
	resp, err := rt.Transport.RoundTrip(req)

	return resp, err
}

// init initializes the underlying transport.
func (rt *RoundTripper) init() {
	if rt.Transport == nil {
		rt.Transport = http.DefaultTransport
	}

	if rt.Limiter == nil {
		rt.Limiter = &unlimitedLimiter{}
	}
}
