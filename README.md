[![CodeQL](https://github.com/HatsuneMiku3939/rltransport/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/HatsuneMiku3939/rltransport/actions/workflows/codeql-analysis.yml)
[![Unit Test](https://github.com/HatsuneMiku3939/rltransport/actions/workflows/test.yaml/badge.svg)](https://github.com/HatsuneMiku3939/rltransport/actions/workflows/test.yaml)

# rltransport
The RoundTripper which limits the number of concurrent requests.

## examples

```golang
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/HatsuneMiku3939/rltransport"

	"golang.org/x/time/rate"
)

const (
	// TestBurstSize is the default value for the rate limiter's burst size.
	TestBurstSize = 10
	// TestRefillRate is the default value for the rate limiter's refill rate.
	TestRefillRate = 1.0
	// TestURL is the URL to use for testing.
	TestHost = "http://localhost:8080/"
)

func main() {
	// Create a "tocket bucket" limiter with a burst size of 10 and a refill rate of 1.0/sec.
	limiter := rate.NewLimiter(TestRefillRate, TestBurstSize)

	// Create a new http.Client with the limiter.
	client := &http.Client{
		Transport: &rltransport.RoundTripper{
			Limiter: limiter,
		},
	}

	// Make a request to the server.
	// First 10 requests will be sented immadiately, after that it will be sented by 1.0 req/sec.
	for i := 0; i < 20; i++ {
		res, _ := client.Get(TestHost)
		fmt.Printf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), res.Status)
	}
}
```
