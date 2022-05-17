package test

import (
	"net/http"
	"testing"

	"github.com/HatsuneMiku3939/rltransport"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// RoundTripperTestSuite is a testify suite for RoundTripper
type RounterTripperTestSuite struct {
	suite.Suite
}

func (s *RounterTripperTestSuite) TestUnlimited() {
	// dummy RoundTripper for test
	monitor := NewMoniteringTripper()

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
		assert.Equalf(s.T(), c.want, monitor.TripCount, "sendCount: %d want %d", c.sendCount, c.want)
	}
}

// TestSimpleLimiter is a test for simple limiter
func (s *RounterTripperTestSuite) TestSimpleLimiter() {
	// dummy RoundTripper for test
	monitor := NewMoniteringTripper()

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
		assert.Equalf(s.T(), c.want, monitor.TripCount, "sendCount: %d want %d", c.sendCount, c.want)
	}
}

func TestRounterTripper(t *testing.T) {
	suite.Run(t, new(RounterTripperTestSuite))
}
