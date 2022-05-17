package test

import (
	"context"
	"fmt"
)

type simpleLimiter struct {
	Count int
	Limit int
}

func (l *simpleLimiter) Wait(ctx context.Context) error {
	if l.Count >= l.Limit {
		return fmt.Errorf("limit exceeded")
	}

	l.Count++
	return nil
}

func (l *simpleLimiter) Reset() {
	l.Count = 0
}
