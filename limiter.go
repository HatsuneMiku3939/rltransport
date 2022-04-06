package rltransport

import (
	"context"
)

// Limiter is an interface that rate limiter must implement.
type Limiter interface {
	// Wait blocks until the request can be sent.
	Wait(ctx context.Context) error
}

// unlimitedLimiter is a limiter that always allows requests to be sent.
type unlimitedLimiter struct{}

// Wait always success.
func (l *unlimitedLimiter) Wait(ctx context.Context) error {
	return nil
}
