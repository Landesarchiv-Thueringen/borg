package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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

type VeraPDFOutput struct {
	Report Report `json:"report"`
}

type Report struct {
	Jobs []Job `json:"jobs"`
}

type Job struct {
	ValidationResult ValidationResult `json:"validationResult"`
}

type ValidationResult struct {
	ProfileName string `json:"profileName"`
	Compliant   bool   `json:"compliant"`
}

var defaultResponse = "veraPDF API is running"
var workDir = "/borg/tools/verapdf"
var storeDir = "/borg/file-store"
var outputFormat = "json"

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("", getDefaultResponse)
	router.GET("/validate/:profile", validateFile)
	router.Run("0.0.0.0:80")
}

func getDefaultResponse(context *gin.Context) {
	context.String(http.StatusOK, defaultResponse)
}

func validateFile(context *gin.Context) {
	profile := context.Param("profile")
	if profile == "" {
		errorMessage := "no veraPDF profile declared"
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
		log.Println(err.Error())
		response := ToolResponse{
			Error: &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	cmd := exec.Command(
		"/bin/ash",
		filepath.Join(workDir, "bin/verapdf"),
		"-f",
		profile,
		"--format",
		outputFormat,
		"-v",
		fileStorePath,
	)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	veraPDFOutput, err := cmd.Output()
	// exit status 1: file for profile invalid but validation job successful
	if err != nil && err.Error() != "exit status 1" {
		log.Println(err.Error())
		log.Println(stderr.String())
		errorMessage := "error executing verPDF\n" + stderr.String()
		response := ToolResponse{
			Error: &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	veraPDFOutputString := string(veraPDFOutput)
	processVeraPDFOutput(context, veraPDFOutputString)
}

func processVeraPDFOutput(context *gin.Context, output string) {
	var veraPDFOutput VeraPDFOutput
	err := json.NewDecoder(strings.NewReader(output)).Decode(&veraPDFOutput)
	if err != nil {
		errorMessage := "unable parse veraPDF output"
		log.Println(errorMessage)
		log.Println(err.Error())
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
	if veraPDFOutput.Report.Jobs != nil && len(veraPDFOutput.Report.Jobs) > 0 {
		extractedFeatures["valid"] = strconv.FormatBool(
			veraPDFOutput.Report.Jobs[0].ValidationResult.Compliant,
		)
	}
	context.JSON(http.StatusOK, response)
}
