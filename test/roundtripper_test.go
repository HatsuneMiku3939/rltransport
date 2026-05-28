package test

import (
	"net/http"
	"sync"
	"testing"

	"github.com/HatsuneMiku3939/rltransport"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// RoundTripperTestSuite is a testify suite for RoundTripper
type RoundTripperTestSuite struct {
	suite.Suite
}

func (s *RoundTripperTestSuite) TestUnlimited() {
	// dummy RoundTripper for test
	monitor := NewMonitoringTripper()

	// Create a new http client with unlimited limit
	client := &http.Client{
		Transport: &rltransport.RoundTripper{
			Transport: monitor,
		},
	}

	// test cases
	cases := []struct {
		sendCount int
		want      int
	}{
		{
			sendCount: 0,
			want:      0,
		},
		{
			sendCount: 1,
			want:      1,
		},
		{
			sendCount: 100000,
			want:      100000,
		},
	}

	for _, c := range cases {
		monitor.Reset()

		// Send requests
		for i := 0; i < c.sendCount; i++ {
			client.Get("http://example.com")
		}

		// Check the result
		assert.Equalf(s.T(), c.want, monitor.Count(), "sendCount: %d want %d", c.sendCount, c.want)
	}
}

// TestSimpleLimiter is a test for simple limiter
func (s *RoundTripperTestSuite) TestSimpleLimiter() {
	// dummy RoundTripper for test
	monitor := NewMonitoringTripper()

	// simple limiter for test
	limiter := &simpleLimiter{
		Count: 0,
		Limit: 100,
	}

	// Create a new http client with limiter
	client := &http.Client{
		Transport: &rltransport.RoundTripper{
			Transport: monitor,
			Limiter:   limiter,
		},
	}

	// test cases
	cases := []struct {
		sendCount int
		want      int
	}{
		{
			sendCount: 0,
			want:      0,
		},
		{
			sendCount: 1,
			want:      1,
		},
		// send 100000 requests, limited to 100 by simple limiter
		{
			sendCount: 100000,
			want:      limiter.Limit,
		},
	}

	for _, c := range cases {
		failCount := 0
		monitor.Reset()
		limiter.Reset()

		// Send requests
		for i := 0; i < c.sendCount; i++ {
			_, err := client.Get("http://example.com")
			if err != nil {
				failCount++
			}
		}

		// Check the result
		assert.Equalf(s.T(), c.want, monitor.Count(), "sendCount: %d want %d", c.sendCount, c.want)
	}
}

func (s *RoundTripperTestSuite) TestNew() {
	limiter := &simpleLimiter{
		Count: 0,
		Limit: 1,
	}

	rt := rltransport.New(limiter)
	assert.Same(s.T(), limiter, rt.Limiter)
	assert.NotNil(s.T(), rt.Transport)

	monitor := NewMonitoringTripper()
	rt.Transport = monitor

	client := &http.Client{
		Transport: rt,
	}

	_, err := client.Get("http://example.com")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, monitor.Count())
}

func (s *RoundTripperTestSuite) TestNewWithNilLimiter() {
	rt := rltransport.New(nil)
	assert.NotNil(s.T(), rt.Limiter)
	assert.NotNil(s.T(), rt.Transport)

	monitor := NewMonitoringTripper()
	rt.Transport = monitor

	client := &http.Client{
		Transport: rt,
	}

	_, err := client.Get("http://example.com")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, monitor.Count())
}

func (s *RoundTripperTestSuite) TestNewWithTransport() {
	limiter := &simpleLimiter{
		Count: 0,
		Limit: 1,
	}
	monitor := NewMonitoringTripper()

	cases := []struct {
		name                 string
		limiter              rltransport.Limiter
		transport            http.RoundTripper
		wantLimiter          rltransport.Limiter
		wantTransport        http.RoundTripper
		wantDefaultLimiter   bool
		wantDefaultTransport bool
		wantRequest          bool
	}{
		{
			name:          "custom limiter and transport",
			limiter:       limiter,
			transport:     monitor,
			wantLimiter:   limiter,
			wantTransport: monitor,
			wantRequest:   true,
		},
		{
			name:               "nil limiter and custom transport",
			transport:          monitor,
			wantTransport:      monitor,
			wantDefaultLimiter: true,
			wantRequest:        true,
		},
		{
			name:                 "custom limiter and nil transport",
			limiter:              limiter,
			wantLimiter:          limiter,
			wantDefaultTransport: true,
		},
	}

	for _, c := range cases {
		s.Run(c.name, func() {
			monitor.Reset()
			limiter.Reset()

			rt := rltransport.NewWithTransport(c.limiter, c.transport)

			if c.wantDefaultLimiter {
				assert.NotNil(s.T(), rt.Limiter)
			} else {
				assert.Same(s.T(), c.wantLimiter, rt.Limiter)
			}

			if c.wantDefaultTransport {
				assert.Same(s.T(), http.DefaultTransport, rt.Transport)
			} else {
				assert.Same(s.T(), c.wantTransport, rt.Transport)
			}

			if c.wantRequest {
				client := &http.Client{
					Transport: rt,
				}

				_, err := client.Get("http://example.com")
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), 1, monitor.Count())
			}
		})
	}
}

func (s *RoundTripperTestSuite) TestConcurrentRoundTrip() {
	const (
		workerCount       = 16
		requestsPerWorker = 64
	)
	requestCount := workerCount * requestsPerWorker

	monitor := NewMonitoringTripper()
	limiter := newConcurrentLimiter(requestCount)
	client := &http.Client{
		Transport: &rltransport.RoundTripper{
			Transport: monitor,
			Limiter:   limiter,
		},
	}

	var wg sync.WaitGroup
	start := make(chan struct{})
	errs := make(chan error, requestCount)

	for worker := 0; worker < workerCount; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start

			for request := 0; request < requestsPerWorker; request++ {
				resp, err := client.Get("http://example.com")
				if err != nil {
					errs <- err
					continue
				}
				if resp != nil {
					_ = resp.Body.Close()
				}
			}
		}()
	}

	close(start)
	wg.Wait()
	close(errs)

	for err := range errs {
		assert.NoError(s.T(), err)
	}
	assert.Equal(s.T(), requestCount, monitor.Count())
	assert.Equal(s.T(), requestCount, limiter.Count())
}

func TestRoundTripper(t *testing.T) {
	suite.Run(t, new(RoundTripperTestSuite))
}
