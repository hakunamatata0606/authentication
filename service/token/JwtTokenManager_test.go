package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestJwt(t *testing.T) {
	jwtToken := NewJwtTokenManager("secret")
	claims := ClaimMap{
		"foo": "bar",
		"abc": "xyz",
	}

	tokenStr, err := jwtToken.CreateToken(claims, 2)
	require.Nil(t, err)
	claims, err = jwtToken.ParseToken(tokenStr)
	require.Nil(t, err)
	require.Equal(t, "bar", claims["foo"])
	require.Equal(t, "xyz", claims["abc"])

	_, ok := claims["dummy"]
	require.Equal(t, false, ok)

	time.Sleep(time.Duration(3 * time.Second))
	_, err = jwtToken.ParseToken(tokenStr)
	require.NotNil(t, err)
	require.Equal(t, "token timeout", err.Error())

}
