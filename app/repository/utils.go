package repository

import (
	"context"
	"time"
)

func InitContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	return ctx, cancel
}
