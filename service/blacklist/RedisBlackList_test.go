package blacklist

import (
	"authentication/config"
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestRedisBlackList(t *testing.T) {
	c := config.GetConfig()
	rbl := NewRedisBlackList(&redis.Options{
		Addr: c.Redis.Addr,
	})
	key := "aloha"

	ctx := context.Background()
	err := rbl.Add(ctx, key, "dummy", time.Duration(5*time.Second))
	require.Nil(t, err)

	err = rbl.IsExist(ctx, key)
	require.Nil(t, err)

	time.Sleep(5 * time.Second)
	err = rbl.IsExist(ctx, key)
	require.NotNil(t, err)
}
