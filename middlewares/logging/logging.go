package logging

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

func LogMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		hostname := ctx.ClientIP()

		uri := ctx.Request.RequestURI

		log.Printf("%s %s -->", hostname, uri)
		ctx.Next()

		took := time.Since(start)
		log.Printf("%s %s <-- %d [%s]", hostname, uri, ctx.Writer.Status(), took)
	}
}
