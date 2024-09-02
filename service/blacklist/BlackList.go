package blacklist

import (
	"context"
	"time"
)

type BlackList interface {
	Add(ctx context.Context, key string, value string, timeout time.Duration) error
	IsExist(ctx context.Context, key string) error
}
