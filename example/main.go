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

	// Will be printed:
	// [2022-04-06 20:11:09] 200 OK
	// [2022-04-06 20:11:09] 200 OK
	// [2022-04-06 20:11:09] 200 OK
	// [2022-04-06 20:11:09] 200 OK
	// [2022-04-06 20:11:09] 200 OK
	// [2022-04-06 20:11:09] 200 OK
	// [2022-04-06 20:11:09] 200 OK
	// [2022-04-06 20:11:09] 200 OK
	// [2022-04-06 20:11:09] 200 OK
	// [2022-04-06 20:11:09] 200 OK
	// [2022-04-06 20:11:10] 200 OK  ## <-- First 10 requests will be sented immadiately.
	// [2022-04-06 20:11:11] 200 OK
	// [2022-04-06 20:11:12] 200 OK
	// [2022-04-06 20:11:13] 200 OK
	// [2022-04-06 20:11:14] 200 OK
	// [2022-04-06 20:11:15] 200 OK
	// [2022-04-06 20:11:16] 200 OK
	// [2022-04-06 20:11:17] 200 OK
	// [2022-04-06 20:11:18] 200 OK
	// [2022-04-06 20:11:19] 200 OK
}
