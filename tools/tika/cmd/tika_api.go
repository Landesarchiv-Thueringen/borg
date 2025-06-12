package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ToolResponse struct {
	ToolVersion  string                      `json:"toolVersion"`
	ToolOutput   string                      `json:"toolOutput"`
	OutputFormat string                      `json:"outputFormat"`
	Features     map[string]ToolFeatureValue `json:"features"`
	Error        *string                     `json:"error"`
}

type ToolFeatureValue struct {
	Value interface{} `json:"value"`
	Label *string     `json:"label"`
}

type TikaOutput struct {
	MimeType    *string `json:"Content-Type"`
	Encoding    *string `json:"Content-Encoding"`
	PDFVersion  *string `json:"pdf:PDFVersion"`
	PDFAVersion *string `json:"pdfa:PDFVersion"`
}

var (
	FORMAT_VERSION_LABEL = "Formatversion"
	MIME_TYPE_LABEL      = "Mime-Type"
	TEXT_ENCODING_LABEL  = "Zeichenkodierung"
)

const (
	DEFAULT_RESPONSE = "Tika API is running"
	WORK_DIR         = "/borg/tools/tika"
	STORE_DIR        = "/borg/file-store"
	TIME_OUT         = 30 * time.Second
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

func extractMetadata(ginContext *gin.Context) {
	fileStorePath := filepath.Join(STORE_DIR, ginContext.Query("path"))
	_, err := os.Stat(fileStorePath)
	if err != nil {
		errorMessage := fmt.Sprintf("error processing file: %s", fileStorePath)
		log.Println(errorMessage)
		log.Println(err)
		response := ToolResponse{
			Error: &errorMessage,
		}
		ginContext.JSON(http.StatusOK, response)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT)
	defer cancel()
	cmd := exec.CommandContext(
		ctx,
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
	if ctx.Err() == context.DeadlineExceeded {
		errorMessage := fmt.Sprintf("Timeout exceeded after %s.", TIME_OUT)
		response := ToolResponse{
			ToolVersion: toolVersion,
			Error:       &errorMessage,
		}
		ginContext.JSON(http.StatusOK, response)
		return
	}
	if err != nil {
		errorMessage := fmt.Sprintf("error executing Tika command: %s", stderr.String())
		log.Println(errorMessage)
		log.Println(err)
		response := ToolResponse{
			ToolVersion: toolVersion,
			Error:       &errorMessage,
		}
		ginContext.JSON(http.StatusOK, response)
		return
	}
	tikaOutputString := string(tikaOutput)
	processTikaOutput(ginContext, tikaOutputString)
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
	extractedFeatures := make(map[string]ToolFeatureValue)
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
		extractedFeatures["format:mimeType"] = ToolFeatureValue{
			Value: mimeType,
			Label: &MIME_TYPE_LABEL,
		}
	}
	if parsedTikaOutput.Encoding != nil {
		extractedFeatures["text:encoding"] = ToolFeatureValue{
			Value: *parsedTikaOutput.Encoding,
			Label: &TEXT_ENCODING_LABEL,
		}
	}
	// use PDF/A version if existing
	if parsedTikaOutput.PDFAVersion != nil {
		extractedFeatures["format:version"] = ToolFeatureValue{
			Value: "PDF/" + *parsedTikaOutput.PDFAVersion,
			Label: &FORMAT_VERSION_LABEL,
		}
	} else if parsedTikaOutput.PDFVersion != nil {
		// no PDF/A version --> use normal version info
		extractedFeatures["format:version"] = ToolFeatureValue{
			Value: *parsedTikaOutput.PDFVersion,
			Label: &FORMAT_VERSION_LABEL,
		}
	}
	context.JSON(http.StatusOK, response)
}
