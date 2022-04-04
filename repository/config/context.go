package config

import (
	"context"
	"time"
)

func TimeoutContext(seconds time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), seconds*time.Second)
}
