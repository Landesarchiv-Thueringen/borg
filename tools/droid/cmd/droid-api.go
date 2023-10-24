package main

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ToolResponse struct {
	ToolOutput        string
	ExtractedFeatures map[string]string
}

var defaultResponse = "droid API is running"
var workDir = "/borg/tools/droid"
var storeDir = "/borg/filestore"
var signatureFilePath = filepath.Join(workDir, "bin/DROID_SignatureFile_V114.xml")
var containerSignatureFilePath = filepath.Join(workDir, "bin/container-signature-20230822.xml")

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
	fileStorePath := filepath.Join(storeDir, context.Query("path"))
	_, err := os.Stat(fileStorePath)
	if err != nil {
		log.Println(err)
		context.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": "error processing file: " + fileStorePath,
		})
		return
	}
	cmd := exec.Command(
		"/bin/ash",
		"/borg/tools/droid/bin/droid-binary-6.7.0-bin/droid.sh",
		"-Ns",
		signatureFilePath,
		"-Nc",
		containerSignatureFilePath,
		"-Nr",
		fileStorePath,
	)
	droidOutput, err := cmd.Output()
	if err != nil {
		log.Println(err)
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "error executing DROID command",
		})
		return
	}
	droidOutputString := string(droidOutput)
	csvReader := csv.NewReader(strings.NewReader(droidOutputString))
	formats, err := csvReader.ReadAll()
	if err != nil {
		log.Println(err)
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "unable to DROID csv output",
		})
		return
	}
	extractedFeatures := make(map[string]string)
	response := ToolResponse{
		ToolOutput: droidOutputString,
	}
	// TODO: implement returning multiple results
	if len(formats) > 1 {
		extractedFeatures["puid"] = formats[1][1]
	}
	response.ExtractedFeatures = extractedFeatures
	context.JSON(http.StatusOK, response)
}
