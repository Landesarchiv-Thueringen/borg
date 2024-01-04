package main

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type ToolResponse struct {
	ToolOutput        *string
	OutputFormat      *string
	ExtractedFeatures *map[string]string
	Error             *string
}

var defaultResponse = "DROID API is running"
var workDir = "/borg/tools/droid"
var storeDir = "/borg/filestore"
var signatureFilePath = filepath.Join(workDir, "bin/DROID_SignatureFile_V114.xml")
var containerSignatureFilePath = filepath.Join(workDir, "bin/container-signature-20230822.xml")
var outputFormat = "csv"

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("", getDefaultResponse)
	router.GET("/identify-file-format", identifyFileFormat)
	router.Run("0.0.0.0:80")
}

func getDefaultResponse(context *gin.Context) {
	context.String(http.StatusOK, defaultResponse)
}

func identifyFileFormat(context *gin.Context) {
	fileStorePath := filepath.Join(storeDir, context.Query("path"))
	_, err := os.Stat(fileStorePath)
	if err != nil {
		log.Println(err)
		errorMessage := "error processing file: " + fileStorePath
		response := ToolResponse{
			Error: &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	cmd := exec.Command(
		"/bin/ash",
		"/borg/tools/droid/bin/droid.sh",
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
		errorMessage := "error executing DROID command"
		response := ToolResponse{
			Error: &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	droidOutputString := string(droidOutput)
	csvReader := csv.NewReader(strings.NewReader(droidOutputString))
	formats, err := csvReader.ReadAll()
	if err != nil {
		log.Println(err.Error())
		errorMessage := "unable to DROID csv output"
		response := ToolResponse{
			ToolOutput:   &droidOutputString,
			OutputFormat: &outputFormat,
			Error:        &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	extractedFeatures := make(map[string]string)
	response := ToolResponse{
		ToolOutput:        &droidOutputString,
		OutputFormat:      &outputFormat,
		ExtractedFeatures: &extractedFeatures,
	}
	// TODO: discuss returning multiple results
	if len(formats) > 1 && formats[1][1] != "" {
		extractedFeatures["puid"] = formats[1][1]
	}
	context.JSON(http.StatusOK, response)
}
