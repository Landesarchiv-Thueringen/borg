package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	STORE_DIR        = "/borg/file-store"
	DEFAULT_RESPONSE = "ODF Validator API is running"
	TIME_OUT         = 30 * time.Second
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

var (
	MIME_TYPE_LABEL = "Mime-Type"
	VALID_LABEL     = "valide"
)

var toolVersion string

func main() {
	toolVersion = getToolVersion()
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("", getDefaultResponse)
	router.GET("validate", validate)
	router.Run()
}

func getToolVersion() string {
	cmd := exec.Command(
		"java",
		"-jar",
		"third_party/odfvalidator-0.12.0-jar-with-dependencies.jar",
		"-V",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	r := regexp.MustCompile(`odfvalidator v([0-9]+\.[0-9]+\.[0-9]+)`)
	matches := r.FindStringSubmatch(string(output))
	if len(matches) != 2 {
		log.Fatal("couldn't extract ODF Validator version from tool output")
	}
	return matches[1]
}

// getDefaultResponse is the test endpoint for checking if the service is running.
func getDefaultResponse(context *gin.Context) {
	context.String(http.StatusOK, DEFAULT_RESPONSE)
}

// validate is the API endpoint for validating a file with ODF Validator.
func validate(context *gin.Context) {
	path := filepath.Join(STORE_DIR, context.Query("path"))
	valid, output, err := validateFile(path)
	if err != nil {
		errorMessage := err.Error()
		response := ToolResponse{
			ToolVersion: toolVersion,
			Error:       &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	extractedFeatures := make(map[string]ToolFeatureValue)
	extractedFeatures["format:valid"] = ToolFeatureValue{
		Value: valid,
		Label: &VALID_LABEL,
	}
	r := regexp.MustCompile(`Media Type:\s*([a-zA-Z0-9.+/-]+)`)
	matches := r.FindStringSubmatch(output)
	if len(matches) == 2 {
		extractedFeatures["format:mimeType"] = ToolFeatureValue{
			Value: matches[1],
			Label: &MIME_TYPE_LABEL,
		}
	}
	response := ToolResponse{
		ToolVersion:  toolVersion,
		ToolOutput:   output,
		OutputFormat: "text",
		Features:     extractedFeatures,
	}
	context.JSON(http.StatusOK, response)
}

// validateFile uses ODF Validator to determine whether a given file is a valid ODF document.
//
// It returns
// - a boolean indicating whether the file is valid ODF
// - the command's combined stdout and stderr output
// - an error if validation failed for unforeseen reasons.
func validateFile(path string) (bool, string, error) {
	_, err := os.Stat(path)
	if err != nil {
		errorMessage := "error processing file: " + path
		log.Println(errorMessage)
		log.Println(err)
		return false, "", errors.New(errorMessage)
	}
	// -v for verbose output to extract the MIME type
	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT)
	defer cancel()
	cmd := exec.CommandContext(
		ctx,
		"java",
		"-jar",
		"third_party/odfvalidator-0.12.0-jar-with-dependencies.jar",
		"-v",
		"-c",
		"-e",
		path,
	)
	output, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		errorMessage := fmt.Sprintf("Timeout exceeded after %s.", TIME_OUT)
		log.Println(errorMessage)
		return false, "", errors.New(errorMessage)
	}
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// Determined the given file to be invalid.
			if exitError.ExitCode() == 2 {
				return false, string(output), nil
			}
		}
		errorMessage := fmt.Sprintf("error executing ODF-Validator command: %v", err)
		log.Println(errorMessage)
		return false, "", errors.New(errorMessage)
	}
	return true, string(output), nil
}
