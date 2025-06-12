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

var (
	VALID_LABEL = "valide"
)

const (
	TOOL_VERSION     = "2.1.5"
	STORE_DIR        = "/borg/file-store"
	DEFAULT_RESPONSE = "OOXML-Validator API is running"
	TIME_OUT         = 30 * time.Second
)

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("", getDefaultResponse)
	router.GET("validate", validate)
	router.Run()
}

// getDefaultResponse is the test endpoint for checking if the service is running.
func getDefaultResponse(context *gin.Context) {
	context.String(http.StatusOK, DEFAULT_RESPONSE)
}

// validate is the API endpoint for validating a file with OOXML-Validator.
func validate(context *gin.Context) {
	path := filepath.Join(STORE_DIR, context.Query("path"))
	valid, output, err := validateFile(path)
	if err != nil {
		errorMessage := err.Error()
		response := ToolResponse{
			ToolVersion: TOOL_VERSION,
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
	response := ToolResponse{
		ToolVersion:  TOOL_VERSION,
		ToolOutput:   output,
		OutputFormat: "text",
		Features:     extractedFeatures,
	}
	context.JSON(http.StatusOK, response)
}

// validateFile uses OOXML-Validator to determine whether a given file is a valid Open Office XML document.
//
// It returns
// - a boolean indicating whether the file is valid OOXML
// - the command's combined stdout and stderr output
// - an error if validation failed for unforeseen reasons.
func validateFile(path string) (valid bool, output string, err error) {
	_, err = os.Stat(path)
	if err != nil {
		err = fmt.Errorf("error processing file %s: %w", path, err)
		log.Println(err)
		return false, "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT)
	defer cancel()
	cmd := exec.CommandContext(ctx, "third_party/OOXMLValidatorCLI", path)
	outputBytes, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		errorMessage := fmt.Sprintf("Timeout exceeded after %s.", TIME_OUT)
		log.Println(string(output))
		log.Println(errorMessage)
		return false, string(output), errors.New(errorMessage)
	}
	output = string(outputBytes)
	if err != nil {
		err = fmt.Errorf("error executing OOXML-Validator command: %w", err)
		log.Println(err)
		return false, output, err
	}
	valid = output == "[]"
	return valid, output, nil
}
