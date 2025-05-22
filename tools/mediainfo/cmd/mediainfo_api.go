package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type ToolResponse struct {
	ToolVersion  string                 `json:"toolVersion"`
	ToolOutput   string                 `json:"toolOutput"`
	OutputFormat string                 `json:"outputFormat"`
	Features     map[string]interface{} `json:"features"`
	Error        *string                `json:"error"`
}

const (
	defaultResponse = "MediaInfo API is running"
	workDir         = "/borg/tools/magika"
	storeDir        = "/borg/file-store"
)

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("", getDefaultResponse)
	router.GET("/extract-metadata", extractMetadata)
	router.Run()
}

func getDefaultResponse(context *gin.Context) {
	context.String(http.StatusOK, defaultResponse)
}

func extractMetadata(context *gin.Context) {
	fileStorePath := filepath.Join(storeDir, context.Query("path"))
	_, err := os.Stat(fileStorePath)
	if err != nil {
		log.Println(err)
		errorMessage := fmt.Sprintf("error processing file: %s", fileStorePath)
		response := ToolResponse{
			ToolVersion: "1.223",
			Error:       &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	cmd := exec.Command(
		"mediainfo",
		"--Output=XML",
		fileStorePath,
	)
	output, err := cmd.CombinedOutput()
	outputString := string(output)
	extractedFeatures := make(map[string]interface{})
	response := ToolResponse{
		ToolVersion:  "1.223",
		ToolOutput:   outputString,
		OutputFormat: "xml",
		Features:     extractedFeatures,
	}
	context.JSON(http.StatusOK, response)
}
