package login

import (
	"authentication/config"
	"authentication/models"
	"authentication/service/auth"
	"authentication/service/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginApi(config *config.Config, verifier auth.VerifyAuth, tokenManager token.TokenManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.UserDetail
		if err := ctx.BindJSON(&user); err != nil {
			ctx.Status(http.StatusBadRequest)
			ctx.Abort()
			return
		}
		claims, err := verifier.Verify(ctx, &user)
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

		refreshToken, err := tokenManager.CreateToken(claims, config.Token.RefreshTokenTimeout)
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
