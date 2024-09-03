package logout

import (
	"authentication/config"
	"authentication/middlewares/authorization"
	"authentication/service/blacklist"
	"authentication/service/token"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func LogoutApi(config *config.Config, tokenManager token.TokenManager, bl blacklist.BlackList) gin.HandlerFunc {
	return authorization.HandleWithClaims(func(ctx *gin.Context, cm *token.ClaimMap) {
		tokenExp, ok := (*cm)["__exp"].(float64)
		if !ok {
			ctx.Status(http.StatusInternalServerError)
			return
		}
		token, ok := (*cm)["token"].(string)
		if !ok {
			ctx.Status(http.StatusInternalServerError)
			return
		}
		timeout := (config.Token.RefreshTokenTimeout - config.Token.TokenTimeout) - (int64(tokenExp) - time.Now().Unix())
		bl.Add(ctx, token, "", time.Duration(timeout*int64(time.Second)))
		ctx.Status(http.StatusOK)
	})
}
