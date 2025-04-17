package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type ToolResponse struct {
	ToolOutput   string                 `json:"toolOutput"`
	OutputFormat string                 `json:"outputFormat"`
	Features     map[string]interface{} `json:"features"`
	Error        *string                `json:"error"`
}

type TikaOutput struct {
	MimeType    *string `json:"Content-Type"`
	Encoding    *string `json:"Content-Encoding"`
	PDFVersion  *string `json:"pdf:PDFVersion"`
	PDFAVersion *string `json:"pdfa:PDFVersion"`
}

var defaultResponse = "Tika API is running"
var workDir = "/borg/tools/tika"
var storeDir = "/borg/file-store"

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
		errorMessage := fmt.Sprintf("error processing file: %s", fileStorePath)
		log.Println(errorMessage)
		log.Println(err)
		response := ToolResponse{
			Error: &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	cmd := exec.Command(
		"java",
		"-jar",
		filepath.Join(workDir, "third_party/tika-app-2.9.2.jar"),
		"--metadata",
		"--json",
		fileStorePath,
	)
	tikaOutput, err := cmd.CombinedOutput()
	if err != nil {
		errorMessage := fmt.Sprintf("error executing Tika command: %s", string(tikaOutput))
		log.Println(errorMessage)
		log.Println(err)
		response := ToolResponse{
			Error: &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	tikaOutputString := string(tikaOutput)
	processTikaOutput(context, tikaOutputString)
}

func processTikaOutput(context *gin.Context, output string) {
	var parsedTikaOutput TikaOutput
	err := json.NewDecoder(strings.NewReader(output)).Decode(&parsedTikaOutput)
	if err != nil {
		errorMessage := "unable parse Tika output"
		log.Println(errorMessage)
		log.Println(err)
		response := ToolResponse{
			ToolOutput:   output,
			OutputFormat: "text",
			Error:        &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	extractedFeatures := make(map[string]interface{})
	response := ToolResponse{
		ToolOutput:   output,
		OutputFormat: "json",
		Features:     extractedFeatures,
	}
	if parsedTikaOutput.MimeType != nil {
		// removes charset from MIME-Type if existing, example: text/x-yaml; charset=ISO-8859-1
		mimeType := strings.Split(*parsedTikaOutput.MimeType, ";")[0]
		// text/x-web-markdown is not the official Mime type
		// https://www.iana.org/assignments/media-types/media-types.xhtml
		if mimeType == "text/x-web-markdown" {
			mimeType = "text/markdown"
		}
		extractedFeatures["mimeType"] = mimeType
	}
	if parsedTikaOutput.Encoding != nil {
		extractedFeatures["encoding"] = *parsedTikaOutput.Encoding
	}
	// use PDF/A version if existing
	if parsedTikaOutput.PDFAVersion != nil {
		extractedFeatures["formatVersion"] = "PDF/" + *parsedTikaOutput.PDFAVersion
	} else if parsedTikaOutput.PDFVersion != nil {
		// no PDF/A version --> use normal version info
		extractedFeatures["formatVersion"] = *parsedTikaOutput.PDFVersion
	}
	context.JSON(http.StatusOK, response)
}
