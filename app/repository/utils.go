package repository

import (
	"context"
	"time"
)

func InitContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return ctx, cancel
}

func intPtr(i int64) *int64 {
	return &i
}
