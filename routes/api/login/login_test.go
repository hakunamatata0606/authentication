package login

import (
	"authentication/config"
	db "authentication/db/sqlc"
	"authentication/models"
	"authentication/service/auth"
	"authentication/service/password"
	"authentication/service/token"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	r := gin.Default()
	c := config.GetConfig()
	tokenManager := token.NewJwtTokenManager("secret")
	conn, err := sql.Open(c.Db.Driver, c.Db.Addr)
	repo := db.NewRepository(conn)
	require.Nil(t, err)
	pwdManager := password.NewSha256Hash("")
	verifier := auth.NewMysqlAuth(repo, pwdManager)
	r.POST("/login", LoginApi(config.GetConfig(), verifier, tokenManager))

	u := models.UserDetail{
		Username: "foo",
		Password: "bar",
	}

	ujs, err := json.Marshal(u)
	require.Nil(t, err)

	request, err := http.NewRequest("POST", "/login", strings.NewReader(string(ujs)))
	require.Nil(t, err)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	require.Equal(t, http.StatusOK, w.Code)

	jsonMap := map[string]interface{}{}
	err = json.Unmarshal(w.Body.Bytes(), &jsonMap)
	require.Nil(t, err)

	tokenStr, ok := jsonMap["token"].(string)
	require.True(t, ok)

	claims, err := tokenManager.ParseToken(tokenStr)
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

	//failure login

	u = models.UserDetail{
		Username: "foo",
		Password: "dummy",
	}

	ujs, err = json.Marshal(u)
	require.Nil(t, err)

	request, err = http.NewRequest("POST", "/login", strings.NewReader(string(ujs)))
	require.Nil(t, err)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, request)
	require.Equal(t, http.StatusUnauthorized, w.Code)

}
