package refreshtoken

import (
	"authentication/config"
	"authentication/models"
	"authentication/service/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RefreshTokenApi(config *config.Config, tokenManager token.TokenManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var data models.RefreshTokenData
		err := ctx.BindJSON(&data)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			ctx.Abort()
			return
		}
		claims, err := tokenManager.ParseToken(data.RefreshToken)
		if err != nil {
			ctx.Status(http.StatusUnauthorized)
			ctx.Abort()
			return
		}
		token, err := tokenManager.CreateToken(claims, config.Token.TokenTimeout)
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			ctx.Abort()
			return
		}
		refreshToken, err := tokenManager.CreateToken(claims, config.Token.TokenTimeout)
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			ctx.Abort()
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"token":         token,
			"refresh_token": refreshToken,
		})
	}
}
