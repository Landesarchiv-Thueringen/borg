package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type ToolResponse struct {
	ToolVersion  string                 `json:"toolVersion"`
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

const (
	DEFAULT_RESPONSE = "Tika API is running"
	WORK_DIR         = "/borg/tools/tika"
	STORE_DIR        = "/borg/file-store"
)

var toolVersion string

func main() {
	toolVersion = getToolVersion()
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("", getDefaultResponse)
	router.GET("/extract-metadata", extractMetadata)
	router.Run()
}

func getDefaultResponse(context *gin.Context) {
	context.String(http.StatusOK, DEFAULT_RESPONSE)
}

func getToolVersion() string {
	cmd := exec.Command(
		"java",
		"-jar",
		filepath.Join(WORK_DIR, "third_party/tika-app-2.9.2.jar"),
		"--version",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	r := regexp.MustCompile(`Apache Tika ([0-9]+\.[0-9]+\.[0-9]+)`)
	matches := r.FindStringSubmatch(string(output))
	if len(matches) != 2 {
		log.Fatal("couldn't extract tika version from tool output")
	}
	return matches[1]
}

func extractMetadata(context *gin.Context) {
	fileStorePath := filepath.Join(STORE_DIR, context.Query("path"))
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
		filepath.Join(WORK_DIR, "third_party/tika-app-2.9.2.jar"),
		"--metadata",
		"--json",
		fileStorePath,
	)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	tikaOutput, err := cmd.Output()
	if err != nil {
		errorMessage := fmt.Sprintf("error executing Tika command: %s", stderr.String())
		log.Println(errorMessage)
		log.Println(err)
		response := ToolResponse{
			ToolVersion: toolVersion,
			Error:       &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	tikaOutputString := string(tikaOutput)
	log.Println(tikaOutputString)
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
			ToolVersion:  toolVersion,
			ToolOutput:   output,
			OutputFormat: "text",
			Error:        &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	extractedFeatures := make(map[string]interface{})
	response := ToolResponse{
		ToolVersion:  toolVersion,
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
		extractedFeatures["format:mimeType"] = mimeType
	}
	if parsedTikaOutput.Encoding != nil {
		extractedFeatures["text:encoding"] = *parsedTikaOutput.Encoding
	}
	// use PDF/A version if existing
	if parsedTikaOutput.PDFAVersion != nil {
		extractedFeatures["format:version"] = "PDF/" + *parsedTikaOutput.PDFAVersion
	} else if parsedTikaOutput.PDFVersion != nil {
		// no PDF/A version --> use normal version info
		extractedFeatures["format:version"] = *parsedTikaOutput.PDFVersion
	}
	context.JSON(http.StatusOK, response)
}
