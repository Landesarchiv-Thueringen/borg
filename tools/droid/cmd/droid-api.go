package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"

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
	router.GET("/identify-file-format", identifyFileFormat)
	addr := "0.0.0.0:" + os.Getenv("DROID_API_CONTAINER_PORT")
	router.Run(addr)
}

func getDefaultResponse(context *gin.Context) {
	context.String(http.StatusOK, defaultResponse)
}

func identifyFileFormat(context *gin.Context) {
	fileStorePath := context.Query("path")
	_, err := os.Stat(fileStorePath)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, err)
	}
	cmd := exec.Command("./bin/droid-binary-6.7.0-bin/droid.sh")
	out, err := cmd.Output()
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, err)
	}
	log.Println(out)
}
