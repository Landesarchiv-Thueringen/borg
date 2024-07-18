

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

const defaultResponse = "DROID API is running"
const workDir = "/borg/tools/droid"
const storeDir = "/borg/file-store"

var signatureFilePath = filepath.Join(workDir, "third_party/DROID_SignatureFile_V114.xml")
var containerSignatureFilePath = filepath.Join(workDir, "third_party/container-signature-20230822.xml")
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
		"/borg/tools/droid/third_party/droid.sh",
		"-Ns",
		signatureFilePath,
		"-Nc",
		containerSignatureFilePath,
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
	log.Println(droidOutputString)
	csvReader := csv.NewReader(strings.NewReader(droidOutputString))
	formats, err := csvReader.ReadAll()
	if err != nil {
		log.Println(err.Error())
		errorMessage := "unable to parse DROID csv output"
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
	if len(formats) == 0 || len(formats[1]) < 18 {
		errorMessage := "unable to parse DROID csv output"
		response := ToolResponse{
			ToolOutput:   &droidOutputString,
			OutputFormat: &outputFormat,
			Error:        &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	if formats[1][14] != "" {
		extractedFeatures["puid"] = formats[1][14]
	}
	if formats[1][15] != "" {
		extractedFeatures["mimeType"] = formats[1][15]
	}
	if formats[1][16] != "" {
		extractedFeatures["formatName"] = formats[1][16]
	}
	if formats[1][17] != "" {
		extractedFeatures["formatVersion"] = formats[1][17]
	}
	context.JSON(http.StatusOK, response)
}
