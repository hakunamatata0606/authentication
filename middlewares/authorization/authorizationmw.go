package authorization

import (
	"authentication/service/blacklist"
	"authentication/service/token"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthorizationMW(tokenManager token.TokenManager, bl blacklist.BlackList) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auths, ok := ctx.Request.Header["Authorization"]
		if !ok || len(auths) == 0 {
			ctx.Status(http.StatusUnauthorized)
			ctx.Abort()
			return
		}
		auth := strings.Split(auths[0], " ")
		if len(auth) != 2 || auth[0] != "Bearer" {
			ctx.Status(http.StatusUnauthorized)
			ctx.Abort()
			return
		}
		token := auth[1]
		if bl.IsExist(ctx, token) == nil {
			ctx.Status(http.StatusUnauthorized)
			ctx.Abort()
			return
		}
		claims, err := tokenManager.ParseToken(token)
		if err != nil {
			ctx.Redirect(http.StatusTemporaryRedirect, "/refresh_token")
			ctx.Abort()
			return
		}
		ctx.Set("tokenClaims", claims)
		ctx.Next()
	}
}

func HandleWithClaims(handle func(*gin.Context, *token.ClaimMap)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c, ok := ctx.Get("tokenClaims")
		if !ok {
			ctx.Status(http.StatusInternalServerError)
			ctx.Abort()
			return
		}
		claims, ok := c.(token.ClaimMap)
		if !ok {
			ctx.Status(http.StatusInternalServerError)
			ctx.Abort()
			return
		}
		handle(ctx, &claims)
	}
}
