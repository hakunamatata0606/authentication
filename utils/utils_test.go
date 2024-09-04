package utils

import (
	"authentication/service/token"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUtils(t *testing.T) {
	expected := []interface{}{"user", "admin"}
	claims := token.ClaimMap{
		"role": expected,
	}
	roles := GetRolesFromClaims(&claims)
	require.Equal(t, []string{"user", "admin"}, roles)
}
