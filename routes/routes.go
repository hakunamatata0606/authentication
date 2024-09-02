package routes

import (
	"authentication/middlewares/logging"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	router.Use(logging.LogMW())

	//dummy api to check if be is ok
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	return router
}
