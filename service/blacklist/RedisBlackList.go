package blacklist

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisBlackList struct {
	client *redis.Client
}

func NewRedisBlackList(opt *redis.Options) BlackList {
	rdb := redis.NewClient(opt)
	return &RedisBlackList{
		client: rdb,
	}
}

func (rbl *RedisBlackList) Add(ctx context.Context, key string, value string, timeout time.Duration) error {
	err := rbl.client.Set(ctx, key, value, timeout).Err()
	return err
}

func (rbl *RedisBlackList) IsExist(ctx context.Context, key string) error {
	_, err := rbl.client.Get(ctx, key).Result()
	return err
}
