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

var (
	FORMAT_VERSION_LABEL = "Formatversion"
	MIME_TYPE_LABEL      = "Mime-Type"
	PUID_LABEL           = "PUID"
	VALID_LABEL          = "valide"
)

const (
	DEFAULT_RESPONSE = "veraPDF API is running"
	WORK_DIR         = "/borg/tools/verapdf"
	STORE_DIR        = "/borg/file-store"
	TIMEOUT          = 60 * time.Second
)

var toolVersion string

func main() {
	toolVersion = getToolVersion()
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("", getDefaultResponse)
	router.GET("/validate/:profile", validateFile)
	router.Run()
}

func getDefaultResponse(context *gin.Context) {
	context.String(http.StatusOK, DEFAULT_RESPONSE)
}

func getToolVersion() string {
	cmd := exec.Command(
		"/bin/ash",
		filepath.Join(WORK_DIR, "third_party/verapdf"),
		"--version",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	r := regexp.MustCompile(`veraPDF ([0-9]+\.[0-9]+\.[0-9]+)`)
	matches := r.FindStringSubmatch(string(output))
	if len(matches) != 2 {
		log.Fatal("couldn't extract veraPDF version from tool output")
	}
	return matches[1]
}

func validateFile(ginContext *gin.Context) {
	profile := ginContext.Param("profile")
	if profile == "" {
		errorMessage := "no veraPDF profile declared"
		log.Println(errorMessage)
		response := ToolResponse{
			ToolVersion: toolVersion,
			Error:       &errorMessage,
		}
		ginContext.JSON(http.StatusOK, response)
		return
	}
	fileStorePath := filepath.Join(STORE_DIR, ginContext.Query("path"))
	_, err := os.Stat(fileStorePath)
	if err != nil {
		errorMessage := fmt.Sprintf("error processing file: %s", fileStorePath)
		log.Println(errorMessage)
		log.Println(err.Error())
		response := ToolResponse{
			ToolVersion: toolVersion,
			Error:       &errorMessage,
		}
		ginContext.JSON(http.StatusOK, response)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	cmd := exec.CommandContext(
		ctx,
		"/bin/ash",
		filepath.Join(WORK_DIR, "third_party/verapdf"),
		"-f", profile,
		"--format", "json",
		"-v", fileStorePath,
	)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	veraPDFOutput, err := cmd.Output()
	if ctx.Err() == context.DeadlineExceeded {
		errorMessage := fmt.Sprintf("Timeout exceeded after %s.", TIMEOUT)
		log.Println(errorMessage)
		response := ToolResponse{
			ToolVersion: toolVersion,
			Error:       &errorMessage,
		}
		ginContext.JSON(http.StatusOK, response)
		return
	}
	// exit status 1: file for profile invalid but validation job successful
	if err != nil && err.Error() != "exit status 1" {
		log.Println(err.Error())
		log.Println(stderr.String())
		errorMessage := fmt.Sprintf("error executing verPDF: %s", stderr.String())
		response := ToolResponse{
			ToolVersion: toolVersion,
			Error:       &errorMessage,
		}
		ginContext.JSON(http.StatusOK, response)
		return
	}
	veraPDFOutputString := string(veraPDFOutput)
	processVeraPDFOutput(ginContext, veraPDFOutputString, profile)
}

func processVeraPDFOutput(context *gin.Context, output string, profile string) {
	var veraPDFOutput VeraPDFOutput
	err := json.NewDecoder(strings.NewReader(output)).Decode(&veraPDFOutput)
	if err != nil {
		errorMessage := "unable parse veraPDF output"
		log.Println(errorMessage)
		log.Println(err.Error())
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
	if len(veraPDFOutput.Report.Jobs) > 0 {
		extractedFeatures["format:valid"] = ToolFeatureValue{
			Value: veraPDFOutput.Report.Jobs[0].ValidationResult.Compliant,
			Label: &VALID_LABEL,
		}
		switch profile {
		case "1a":
			extractedFeatures["format:puid"] = getPuidFeature("fmt/95")
			extractedFeatures["format:mimeType"] = getMimeTypeFeature("application/pdf")
			extractedFeatures["format:version"] = getVersionFeature("PDF/A-1a")
		case "1b":
			extractedFeatures["format:puid"] = getPuidFeature("fmt/354")
			extractedFeatures["format:mimeType"] = getMimeTypeFeature("application/pdf")
			extractedFeatures["format:version"] = getVersionFeature("PDF/A-1b")
		case "2a":
			extractedFeatures["format:puid"] = getPuidFeature("fmt/476")
			extractedFeatures["format:mimeType"] = getMimeTypeFeature("application/pdf")
			extractedFeatures["format:version"] = getVersionFeature("PDF/A-2a")
		case "2b":
			extractedFeatures["format:puid"] = getPuidFeature("fmt/477")
			extractedFeatures["format:mimeType"] = getMimeTypeFeature("application/pdf")
			extractedFeatures["format:version"] = getVersionFeature("PDF/A-2b")
		case "2u":
			extractedFeatures["format:puid"] = getPuidFeature("fmt/478")
			extractedFeatures["format:mimeType"] = getMimeTypeFeature("application/pdf")
			extractedFeatures["format:version"] = getVersionFeature("PDF/A-2u")
		case "ua1":
			extractedFeatures["format:mimeType"] = getMimeTypeFeature("application/pdf")
			extractedFeatures["format:version"] = getVersionFeature("PDF/UA")
		}
	}
	context.JSON(http.StatusOK, response)
}

func getPuidFeature(value string) ToolFeatureValue {
	return ToolFeatureValue{
		Value: value,
		Label: &PUID_LABEL,
	}
}

func getMimeTypeFeature(value string) ToolFeatureValue {
	return ToolFeatureValue{
		Value: value,
		Label: &MIME_TYPE_LABEL,
	}
}

func getVersionFeature(value string) ToolFeatureValue {
	return ToolFeatureValue{
		Value: value,
		Label: &FORMAT_VERSION_LABEL,
	}
}
