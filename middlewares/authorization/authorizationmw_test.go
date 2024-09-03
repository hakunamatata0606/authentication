package authorization

import (
	"authentication/config"
	db "authentication/db/sqlc"
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

func TestAuthorization(t *testing.T) {
	r := gin.Default()

	tokenManager := token.NewJwtTokenManager("secret")
	c := config.GetConfig()
	conn, err := sql.Open(c.Db.Driver, c.Db.Addr)
	repo := db.NewRepository(conn)
	require.Nil(t, err)
	pwdManager := password.NewSha256Hash("")
	rbl := blacklist.NewRedisBlackList(&redis.Options{
		Addr: c.Redis.Addr,
	})
	verifier := auth.NewMysqlAuth(repo, pwdManager)
	r.POST("/login", login.LoginApi(config.GetConfig(), verifier, tokenManager))
	protectedGroup := r.Group("/")
	protectedGroup.Use(AuthorizationMW(tokenManager, rbl))
	protectedGroup.GET("/protected", HandleWithClaims(func(ctx *gin.Context, cm *token.ClaimMap) {
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
	require.Equal(t, "user", role)

	bearer := "Bearer " + tokenStr

	request, err = http.NewRequest("GET", "/protected", nil)
	require.Nil(t, err)

	request.Header.Add("Authorization", bearer)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, request)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "hello foo", w.Body.String())

	//failure redirect

	bearer = "Bearer " + "lala"

	request, err = http.NewRequest("GET", "/protected", nil)
	require.Nil(t, err)

	request.Header.Add("Authorization", bearer)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, request)
	require.Equal(t, http.StatusTemporaryRedirect, w.Code)
	require.Equal(t, "/refresh_token", w.Result().Header.Get("Location"))

	//failure
	request, err = http.NewRequest("GET", "/protected", nil)
	require.Nil(t, err)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, request)
	require.Equal(t, http.StatusUnauthorized, w.Code)
}
