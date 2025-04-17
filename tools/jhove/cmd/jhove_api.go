package main

import (
	"encoding/json"
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
	ToolOutput   string                 `json:"toolOutput"`
	OutputFormat string                 `json:"outputFormat"`
	Features     map[string]interface{} `json:"features"`
	Error        string                 `json:"error"`
}

type JhoveOutput struct {
	Root *JhoveRoot `json:"jhove"`
}

type JhoveRoot struct {
	RepInfo []JhoveRepInfo `json:"repInfo"`
}

type JhoveRepInfo struct {
	FormatName    *string `json:"format"`
	FormatVersion *string `json:"version"`
	Validation    *string `json:"status"`
}

var defaultResponse = "JHOVE API is running"
var storeDir = "/borg/file-store"
var wellFormedRegEx = regexp.MustCompile("Well-Formed")
var validRegEx = regexp.MustCompile("Well-Formed and valid")

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("", getDefaultResponse)
	router.GET("validate/:module", validateFile)
	router.Run()
}

func getDefaultResponse(context *gin.Context) {
	context.String(http.StatusOK, defaultResponse)
}

func validateFile(context *gin.Context) {
	module := context.Param("module")
	if module == "" {
		errorMessage := "no JHOVE module declared"
		log.Println(errorMessage)
		response := ToolResponse{
			Error: errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	fileStorePath := filepath.Join(storeDir, context.Query("path"))
	_, err := os.Stat(fileStorePath)
	if err != nil {
		errorMessage := "error processing file: " + fileStorePath
		log.Println(errorMessage)
		log.Println(err)
		response := ToolResponse{
			Error: errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	cmd := exec.Command(
		"./jhove/jhove",
		"-m",
		module+"-hul",
		"-h",
		"json",
		fileStorePath,
	)
	jhoveOutput, err := cmd.Output()
	if err != nil {
		errorMessage := "error executing JHOVE command"
		log.Println(errorMessage)
		log.Println(err)
		response := ToolResponse{
			Error: errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	jhoveOutputString := string(jhoveOutput)
	processJhoveOutput(context, jhoveOutputString, module)
}

func processJhoveOutput(context *gin.Context, output string, module string) {
	var parsedJhoveOutput JhoveOutput
	err := json.NewDecoder(strings.NewReader(output)).Decode(&parsedJhoveOutput)
	if err != nil {
		errorMessage := "unable parse JHOVE output"
		log.Println(errorMessage)
		log.Println(err)
		response := ToolResponse{
			ToolOutput:   output,
			OutputFormat: "text",
			Error:        errorMessage,
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
	if parsedJhoveOutput.Root != nil && len(parsedJhoveOutput.Root.RepInfo) > 0 {
		repInfo := parsedJhoveOutput.Root.RepInfo[0]
		if repInfo.FormatName != nil {
			extractedFeatures["formatName"] = *repInfo.FormatName
		}
		if repInfo.FormatVersion != nil {
			extractedFeatures["formatVersion"] = *repInfo.FormatVersion
		}
		if repInfo.Validation != nil {
			extractedFeatures["wellFormed"] = wellFormedRegEx.MatchString(*repInfo.Validation)
			extractedFeatures["valid"] = validRegEx.MatchString(*repInfo.Validation)
		}
		switch module {
		case "pdf":
			extractedFeatures["mimeType"] = "application/pdf"
		case "html":
			extractedFeatures["mimeType"] = "text/html"
		case "tiff":
			extractedFeatures["mimeType"] = "image/tiff"
		case "jpeg":
			extractedFeatures["mimeType"] = "image/jpeg"
		case "jpeg2000":
			extractedFeatures["mimeType"] = "image/jp2"
		}
	}
	context.JSON(http.StatusOK, response)
}
