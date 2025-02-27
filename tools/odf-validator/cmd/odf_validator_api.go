package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const storeDir = "/borg/file-store"
const defaultResponse = "ODF Validator API is running"

type ToolResponse struct {
	ToolOutput   string                 `json:"toolOutput"`
	OutputFormat string                 `json:"outputFormat"`
	Features     map[string]interface{} `json:"features"`
	Error        string                 `json:"error"`
}

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("", getDefaultResponse)
	router.GET("validate", validate)
	router.Run()
}

// getDefaultResponse is the test endpoint for checking if the service is running.
func getDefaultResponse(context *gin.Context) {
	context.String(http.StatusOK, defaultResponse)
}

// validate is the API endpoint for validating a file with ODF Validator.
func validate(context *gin.Context) {
	path := filepath.Join(storeDir, context.Query("path"))
	valid, output, err := validateFile(path)
	if err != nil {
		response := ToolResponse{
			Error: err.Error(),
		}
		context.JSON(http.StatusOK, response)
		return
	}
	extractedFeatures := make(map[string]interface{})
	extractedFeatures["valid"] = valid
	response := ToolResponse{
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
	cmd := exec.Command(
		"java",
		"-jar",
		"third_party/odfvalidator-0.12.0-jar-with-dependencies.jar",
		"-c",
		"-e",
		path,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// Determined the given file to be invalid.
			if exitError.ExitCode() == 2 {
				return false, string(output), nil
			}
		}
		errorMessage := "error executing ODF-Validator command"
		log.Println(string(output))
		log.Println(errorMessage)
		log.Println(err)
		return false, string(output), errors.New(errorMessage)
	}
	return true, string(output), nil
}
