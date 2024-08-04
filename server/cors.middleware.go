package server

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		access_control_origin, ok := os.LookupEnv("ACCESS_CONTROL_ORIGIN")
		if !ok {
			log.Fatal("Missing ACCESS_CONTROL_ORIGIN env")
		}

		ctx.Writer.Header().Set("Access-Control-Allow-Origin", access_control_origin)
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	}
}
