package rltransport

import (
	"context"
	"net/http"
)

// Limiter is an interface that rate limiter must implement.
type Limiter interface {
	// Wait blocks until the request can be sent.
	Wait(req *http.Request) error
}

// LimiterFunc adapts a function into a Limiter.
type LimiterFunc func(req *http.Request) error

// Wait blocks until the request can be sent.
func (f LimiterFunc) Wait(req *http.Request) error {
	return f(req)
}

// ContextLimiter is an interface for limiters that only need request context.
type ContextLimiter interface {
	// Wait blocks until the request can be sent.
	Wait(ctx context.Context) error
}

// NewContextLimiter adapts a context-only limiter to the request-aware Limiter
// interface. It treats nil as an unlimited limiter.
func NewContextLimiter(limiter ContextLimiter) Limiter {
	if limiter == nil {
		return &unlimitedLimiter{}
	}

	return &contextLimiter{
		limiter: limiter,
	}
}

type contextLimiter struct {
	limiter ContextLimiter
}

func (l *contextLimiter) Wait(req *http.Request) error {
	return l.limiter.Wait(req.Context())
}

// unlimitedLimiter is a limiter that always allows requests to be sent.
type unlimitedLimiter struct{}

// Wait always succeeds.
func (l *unlimitedLimiter) Wait(_ *http.Request) error {
	return nil
}
