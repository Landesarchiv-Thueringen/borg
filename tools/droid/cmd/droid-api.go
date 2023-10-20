package main

import (
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var defaultResponse = "droid API is running"

func main() {
	router := gin.Default()
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"*"})
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type"}
	corsConfig.AllowMethods = []string{"GET", "POST"}
	// It's important that the cors configuration is used before declaring the routes.
	router.Use(cors.New(corsConfig))
	router.GET("", getDefaultResponse)
	addr := "0.0.0.0:" + os.Getenv("DROID_API_CONTAINER_PORT")
	router.Run(addr)
}

func getDefaultResponse(context *gin.Context) {
	context.String(http.StatusOK, defaultResponse)
}
