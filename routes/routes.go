package routes

import (
	"authentication/config"
	db "authentication/db/sqlc"
	"authentication/middlewares/authorization"
	"authentication/middlewares/logging"
	"authentication/routes/api/login"
	"authentication/routes/api/logout"
	refreshtoken "authentication/routes/api/refresh_token"
	"authentication/service/auth"
	"authentication/service/blacklist"
	"authentication/service/password"
	"authentication/service/token"
	"authentication/utils"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

func GetRouter() *gin.Engine {
	cfg := config.GetConfig()

	router := gin.New()
	router.Use(gin.Recovery())

	router.Use(logging.LogMW())

	//dummy api to check if be is ok
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	tokenManager := token.NewJwtTokenManager(cfg.Token.Secret)
	conn, err := sql.Open(cfg.Db.Driver, cfg.Db.Addr)
	if err != nil {
		log.Fatal(err)
	}
	repo := db.NewRepository(conn)
	pwdManager := password.NewSha256Hash("")
	verifier := auth.NewMysqlAuth(repo, pwdManager)
	rbl := blacklist.NewRedisBlackList(&redis.Options{
		Addr: cfg.Redis.Addr,
	})

	router.POST("/login", login.LoginApi(cfg, verifier, tokenManager))

	authenticatedGroup := router.Group("/")
	authenticatedGroup.Use(authorization.AuthorizationMW(tokenManager, rbl))
	{
		authenticatedGroup.POST("/refresh_token", refreshtoken.RefreshTokenApi(cfg, tokenManager))
		authenticatedGroup.POST("/logout", logout.LogoutApi(cfg, tokenManager, rbl))
		authenticatedGroup.GET("/protected/user", authorization.HandleWithClaims(func(ctx *gin.Context, cm *token.ClaimMap) {
			user, ok := (*cm)["user"].(string)
			if !ok {
				ctx.Status(http.StatusInternalServerError)
				return
			}
			roles := utils.GetRolesFromClaims(cm)
			log.Println("aaaa ", roles)
			for _, role := range roles {
				if role == "user" {
					ctx.String(http.StatusOK, "hello "+user)
					return
				}
			}
			ctx.Status(http.StatusForbidden)
		}))

		authenticatedGroup.GET("/protected/admin", authorization.HandleWithClaims(func(ctx *gin.Context, cm *token.ClaimMap) {
			user, ok := (*cm)["user"].(string)
			if !ok {
				ctx.Status(http.StatusInternalServerError)
				return
			}
			roles := utils.GetRolesFromClaims(cm)
			log.Println("aaaa ", cm)
			for _, role := range roles {
				if role == "admin" {
					ctx.String(http.StatusOK, "hello admin "+user)
					return
				}
			}
			ctx.Status(http.StatusForbidden)
		}))
	}

	return router
}
