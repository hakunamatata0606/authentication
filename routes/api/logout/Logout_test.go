package logout

import (
	"authentication/config"
	db "authentication/db/sqlc"
	"authentication/middlewares/authorization"
	"authentication/models"
	"authentication/routes/api/login"
	"authentication/service/auth"
	"authentication/service/blacklist"
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
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestLogout(t *testing.T) {
	r := gin.Default()

	tokenManager := token.NewJwtTokenManager("secret")
	c := config.GetConfig()
	conn, err := sql.Open(c.Db.Driver, c.Db.Addr)
	repo := db.NewRepository(conn)
	require.Nil(t, err)
	pwdManager := password.NewSha256Hash("")
	verifier := auth.NewMysqlAuth(repo, pwdManager)
	rbl := blacklist.NewRedisBlackList(&redis.Options{
		Addr: "localhost:6379",
	})
	r.POST("/login", login.LoginApi(config.GetConfig(), verifier, tokenManager))
	protectedGroup := r.Group("/")
	protectedGroup.Use(authorization.AuthorizationMW(tokenManager, rbl))

	protectedGroup.POST("/logout", LogoutApi(config.GetConfig(), tokenManager, rbl))

	protectedGroup.GET("/protected", authorization.HandleWithClaims(func(ctx *gin.Context, cm *token.ClaimMap) {
		u, ok := (*cm)["user"].(string)
		if !ok {
			ctx.Status(http.StatusInternalServerError)
			return
		}
		ctx.String(http.StatusOK, "hello "+u)
	}))

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
	require.Equal(t, []interface{}([]interface{}{"user"}), role)

	tokenStr, ok = jsonMap["refresh_token"].(string)
	require.True(t, ok)

	claims, err = tokenManager.ParseToken(tokenStr)
	require.Nil(t, err)
	username, ok = claims["user"].(string)
	require.True(t, ok)
	require.Equal(t, "foo", username)
	role, ok = claims["role"]
	require.True(t, ok)
	require.Equal(t, []interface{}([]interface{}{"user"}), role)

	request, err = http.NewRequest("GET", "/protected", nil)
	require.Nil(t, err)
	bearer := "Bearer " + tokenStr
	request.Header.Add("Authorization", bearer)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, request)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "hello foo", w.Body.String())

	request, err = http.NewRequest("POST", "/logout", nil)
	request.Header.Add("Authorization", bearer)
	require.Nil(t, err)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, request)
	require.Equal(t, http.StatusOK, w.Code)

	request, err = http.NewRequest("GET", "/protected", nil)
	require.Nil(t, err)
	request.Header.Add("Authorization", bearer)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, request)
	require.Equal(t, http.StatusUnauthorized, w.Code)

}
