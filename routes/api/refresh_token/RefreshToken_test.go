package refreshtoken

import (
	"authentication/config"
	"authentication/models"
	"authentication/service/token"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestRefreshToken(t *testing.T) {
	r := gin.Default()

	tokenManager := token.NewJwtTokenManager("secret")
	r.POST("/refresh_token", RefreshTokenApi(config.GetConfig(), tokenManager))

	claims := token.ClaimMap{
		"user":  "foo",
		"role":  "user",
		"email": "aloha@example.com",
	}

	token, err := tokenManager.CreateToken(claims, 60)
	require.Nil(t, err)

	u := models.RefreshTokenData{
		RefreshToken: token,
	}

	ujs, err := json.Marshal(u)
	require.Nil(t, err)

	request, err := http.NewRequest("POST", "/refresh_token", strings.NewReader(string(ujs)))
	require.Nil(t, err)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	require.Equal(t, http.StatusOK, w.Code)

	jsonMap := map[string]interface{}{}
	err = json.Unmarshal(w.Body.Bytes(), &jsonMap)
	require.Nil(t, err)

	tokenStr, ok := jsonMap["token"].(string)
	require.True(t, ok)

	claims, err = tokenManager.ParseToken(tokenStr)
	require.Nil(t, err)
	username, ok := claims["user"].(string)
	require.True(t, ok)
	require.Equal(t, "foo", username)
	role, ok := claims["role"]
	require.True(t, ok)
	require.Equal(t, "user", role)

	tokenStr, ok = jsonMap["refresh_token"].(string)
	require.True(t, ok)

	claims, err = tokenManager.ParseToken(tokenStr)
	require.Nil(t, err)
	username, ok = claims["user"].(string)
	require.True(t, ok)
	require.Equal(t, "foo", username)
	role, ok = claims["role"]
	require.True(t, ok)
	require.Equal(t, "user", role)
}
