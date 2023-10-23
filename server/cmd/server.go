package main

import (
	"lath/borg/internal/config"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var defaultResponse = "borg server is running"
var serverConfig config.ServerConfig

func main() {
	serverConfig = config.ParseConfig()
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
	router.POST("analyse-file", analyseFile)
	addr := "0.0.0.0:" + os.Getenv("BORG_SERVER_CONTAINER_PORT")
	router.Run(addr)
}

func getDefaultResponse(context *gin.Context) {
	context.String(http.StatusOK, defaultResponse)
}

func analyseFile(context *gin.Context) {
	file, err := context.FormFile("file")
	// no file received
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "no file received",
		})
		return
	}
	// generate unique file name for storing
	fileStorePath := "filestore/" + uuid.New().String() + "_" + file.Filename
	err = context.SaveUploadedFile(file, fileStorePath)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "unable to save file",
		})
		return
	}
	defer os.Remove(fileStorePath)
	runFileIdentificationTools(fileStorePath)
}

func runFileIdentificationTools(fileStorePath string) {
	for _, tool := range serverConfig.FormatIdentificationTools {
		response, err := http.Get(tool.Endpoint)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(response)
	}
}
