package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/gin-gonic/gin"
)

type ToolResponse struct {
	ToolVersion  string                 `json:"toolVersion"`
	ToolOutput   string                 `json:"toolOutput"`
	OutputFormat string                 `json:"outputFormat"`
	Features     map[string]interface{} `json:"features"`
	Score        *float64               `json:"score"`
	Error        *string                `json:"error"`
}

type Output struct {
	Description string   `json:"description"`
	Extensions  []string `json:"extensions"`
	Group       string   `json:"group"`
	IsText      bool     `json:"is_text"`
	Label       string   `json:"label"`
	MimeType    string   `json:"mime_type"`
}

type Value struct {
	// deep learning model output
	DL Output `json:"dl"`
	// overall tool output
	Output Output  `json:"output"`
	Score  float64 `json:"score"`
}

type Result struct {
	Status string `json:"status"`
	Value  Value  `json:"value"`
}

// output documentation:
// https://github.com/google/magika/blob/main/docs/magika_output.md
type Data struct {
	Path   string `json:"path"`
	Result Result `json:"result"`
}

const (
	defaultResponse = "Magika API is running"
	workDir         = "/borg/tools/magika"
	storeDir        = "/borg/file-store"
)

var toolVersion string

func main() {
	toolVersion = getToolVersion()
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("", getDefaultResponse)
	router.GET("/identify-file-format", identifyFileFormat)
	router.Run()
}

func getDefaultResponse(context *gin.Context) {
	context.String(http.StatusOK, defaultResponse)
}

func getToolVersion() string {
	cmd := exec.Command(
		"magika",
		"--json",
		"--version",
	)
	magikaOutput, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	r := regexp.MustCompile(`magika ([0-9]+\.[0-9]+\.[0-9]+) ([a-zA-Z0-9_]+)`)
	matches := r.FindStringSubmatch(string(magikaOutput))
	if len(matches) != 3 {
		log.Fatal("couldn't extract magika version from tool output")
	}
	return fmt.Sprintf("cli: %s, model: %s", matches[1], matches[2])
}

// identifyFileFormat executes Magika and parses the output of the command.
func identifyFileFormat(context *gin.Context) {
	fileStorePath := filepath.Join(storeDir, context.Query("path"))
	_, err := os.Stat(fileStorePath)
	if err != nil {
		log.Println(err)
		errorMessage := fmt.Sprintf("error processing file: %s", fileStorePath)
		response := ToolResponse{
			ToolVersion: toolVersion,
			Error:       &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	cmd := exec.Command(
		"magika",
		"--json",
		fileStorePath,
	)
	magikaOutput, err := cmd.CombinedOutput()
	magikaOutputString := string(magikaOutput)
	if err != nil {
		log.Println(magikaOutputString)
		log.Println(err)
		errorMessage := fmt.Sprintf("error executing Magika command: %s", magikaOutputString)
		response := ToolResponse{
			ToolVersion: toolVersion,
			ToolOutput:  magikaOutputString,
			Error:       &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	var data []Data
	err = json.Unmarshal(magikaOutput, &data)
	if err != nil {
		log.Println(err)
		errorMessage := "unable to parse Magika JSON output"
		response := ToolResponse{
			ToolVersion:  toolVersion,
			ToolOutput:   magikaOutputString,
			OutputFormat: "json",
			Error:        &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	// one file is analyzed at a time -> only first result is relevant
	features := extractFeatures(data[0])
	var score *float64
	if data[0].Result.Status == "ok" {
		score = &(data[0].Result.Value.Score)
	}
	response := ToolResponse{
		ToolVersion:  toolVersion,
		ToolOutput:   magikaOutputString,
		OutputFormat: "json",
		Features:     features,
		Score:        score,
	}
	context.JSON(http.StatusOK, response)
}

func extractFeatures(data Data) map[string]interface{} {
	features := make(map[string]interface{})
	if data.Result.Status == "ok" {
		if data.Result.Value.Output.MimeType != "" {
			features["format:mimeType"] = data.Result.Value.Output.MimeType
			// image/jpeg2000 is not the official Mime type
			// https://www.iana.org/assignments/media-types/media-types.xhtml
			if features["format:mimeType"] == "image/jpeg2000" {
				features["format:mimeType"] = "image/jp2"
			}
		}
		features["format:isText"] = data.Result.Value.Output.IsText
	}
	return features
}
