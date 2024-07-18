

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ToolResponse struct {
	ToolOutput        *string
	OutputFormat      *string
	ExtractedFeatures *map[string]string
	Error             *string
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
var outputFormat = "json"

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("", getDefaultResponse)
	router.GET("validate/:module", validateFile)
	router.Run("0.0.0.0:80")
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
			Error: &errorMessage,
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
			Error: &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	cmd := exec.Command(
		"./jhove/jhove",
		"-m",
		module+"-hul",
		"-h",
		outputFormat,
		fileStorePath,
	)
	jhoveOutput, err := cmd.Output()
	if err != nil {
		errorMessage := "error executing JHOVE command"
		log.Println(errorMessage)
		log.Println(err)
		response := ToolResponse{
			Error: &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	jhoveOutputString := string(jhoveOutput)
	processJhoveOutput(context, jhoveOutputString)
}

func processJhoveOutput(context *gin.Context, output string) {
	var parsedJhoveOutput JhoveOutput
	err := json.NewDecoder(strings.NewReader(output)).Decode(&parsedJhoveOutput)
	if err != nil {
		errorMessage := "unable parse JHOVE output"
		log.Println(errorMessage)
		log.Println(err)
		response := ToolResponse{
			ToolOutput:   &output,
			OutputFormat: &outputFormat,
			Error:        &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	extractedFeatures := make(map[string]string)
	response := ToolResponse{
		ToolOutput:        &output,
		OutputFormat:      &outputFormat,
		ExtractedFeatures: &extractedFeatures,
	}
	if parsedJhoveOutput.Root != nil && len(parsedJhoveOutput.Root.RepInfo) > 0 {
		repInfo := parsedJhoveOutput.Root.RepInfo[0]
		if repInfo.FormatVersion != nil {
			// if format name was extracted and version doesn't contain it already
			if repInfo.FormatName != nil &&
				!strings.Contains(*repInfo.FormatVersion, *repInfo.FormatName) {
				extractedFeatures["formatVersion"] = *repInfo.FormatName + " " + *repInfo.FormatVersion
			} else {
				extractedFeatures["formatVersion"] = *repInfo.FormatVersion
			}
		}
		if repInfo.Validation != nil {
			extractedFeatures["wellFormed"] = strconv.FormatBool(
				wellFormedRegEx.MatchString(*repInfo.Validation),
			)
			extractedFeatures["valid"] = strconv.FormatBool(
				validRegEx.MatchString(*repInfo.Validation),
			)
		}
	}
	context.JSON(http.StatusOK, response)
}
