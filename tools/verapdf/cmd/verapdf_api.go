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

const (
	DEFAULT_RESPONSE = "veraPDF API is running"
	WORK_DIR         = "/borg/tools/verapdf"
	STORE_DIR        = "/borg/file-store"
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

func validateFile(context *gin.Context) {
	profile := context.Param("profile")
	if profile == "" {
		errorMessage := "no veraPDF profile declared"
		log.Println(errorMessage)
		response := ToolResponse{
			ToolVersion: toolVersion,
			Error:       &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	fileStorePath := filepath.Join(STORE_DIR, context.Query("path"))
	_, err := os.Stat(fileStorePath)
	if err != nil {
		errorMessage := fmt.Sprintf("error processing file: %s", fileStorePath)
		log.Println(errorMessage)
		log.Println(err.Error())
		response := ToolResponse{
			ToolVersion: toolVersion,
			Error:       &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	cmd := exec.Command(
		"/bin/ash",
		filepath.Join(WORK_DIR, "third_party/verapdf"),
		"-f", profile,
		"--format", "json",
		"-v", fileStorePath,
	)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	veraPDFOutput, err := cmd.Output()
	// exit status 1: file for profile invalid but validation job successful
	if err != nil && err.Error() != "exit status 1" {
		log.Println(err.Error())
		log.Println(stderr.String())
		errorMessage := fmt.Sprintf("error executing verPDF: %s", stderr.String())
		response := ToolResponse{
			ToolVersion: toolVersion,
			Error:       &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	veraPDFOutputString := string(veraPDFOutput)
	processVeraPDFOutput(context, veraPDFOutputString, profile)
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
	extractedFeatures := make(map[string]interface{})
	response := ToolResponse{
		ToolVersion:  toolVersion,
		ToolOutput:   output,
		OutputFormat: "json",
		Features:     extractedFeatures,
	}
	if len(veraPDFOutput.Report.Jobs) > 0 {
		extractedFeatures["valid"] =
			veraPDFOutput.Report.Jobs[0].ValidationResult.Compliant
		switch profile {
		case "1a":
			extractedFeatures["puid"] = "fmt/95"
			extractedFeatures["mimeType"] = "application/pdf"
			extractedFeatures["formatVersion"] = "PDF/A-1a"
		case "1b":
			extractedFeatures["puid"] = "fmt/354"
			extractedFeatures["mimeType"] = "application/pdf"
			extractedFeatures["formatVersion"] = "PDF/A-1b"
		case "2a":
			extractedFeatures["puid"] = "fmt/476"
			extractedFeatures["mimeType"] = "application/pdf"
			extractedFeatures["formatVersion"] = "PDF/A-2a"
		case "2b":
			extractedFeatures["puid"] = "fmt/477"
			extractedFeatures["mimeType"] = "application/pdf"
			extractedFeatures["formatVersion"] = "PDF/A-2b"
		case "2u":
			extractedFeatures["puid"] = "fmt/478"
			extractedFeatures["mimeType"] = "application/pdf"
			extractedFeatures["formatVersion"] = "PDF/A-2u"
		case "ua1":
			extractedFeatures["mimeType"] = "application/pdf"
			extractedFeatures["formatVersion"] = "PDF/UA"
		}
	}
	context.JSON(http.StatusOK, response)
}
