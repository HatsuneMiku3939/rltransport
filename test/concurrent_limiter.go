package test

import (
	"fmt"
	"net/http"
	"sync"
)

type concurrentLimiter struct {
	mu    sync.Mutex
	count int
	limit int
}

func newConcurrentLimiter(limit int) *concurrentLimiter {
	return &concurrentLimiter{
		limit: limit,
	}
}

func (l *concurrentLimiter) Wait(req *http.Request) error {
	if err := req.Context().Err(); err != nil {
		return err
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.count >= l.limit {
		return fmt.Errorf("limit exceeded")
	}

	l.count++
	return nil
}

func (l *concurrentLimiter) Count() int {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.count
}
