package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	c := GetConfig()
	require.Equal(t, int64(120), c.Token.TokenTimeout)
	require.Equal(t, int64(360), c.Token.RefreshTokenTimeout)
	require.Equal(t, "mysql", c.Db.Driver)
	require.Equal(t, "bao:123@tcp(172.17.0.2:3306)/test", c.Db.Addr)
	require.Equal(t, "localhost:6379", c.Redis.Addr)
}
