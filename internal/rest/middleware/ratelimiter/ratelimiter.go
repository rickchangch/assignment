package ratelimiter

import (
	"context"
)

// nolint: lll
type Ratelimiter interface {
	AllowByTokenBucket(ctx context.Context, userID string, cost int) (bool, int, error)
}
