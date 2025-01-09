package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ToolResponse struct {
	ToolOutput   string            `json:"toolOutput"`
	OutputFormat string            `json:"outputFormat"`
	Features     map[string]string `json:"features"`
	Error        string            `json:"error"`
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

const defaultResponse = "Magika API is running"
const workDir = "/borg/tools/magika"
const storeDir = "/borg/file-store"

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("", getDefaultResponse)
	router.GET("/identify-file-format", identifyFileFormat)
	router.Run()
}

func getDefaultResponse(context *gin.Context) {
	context.String(http.StatusOK, defaultResponse)
}

// identifyFileFormat executes Magika and parses the output of the command.
func identifyFileFormat(context *gin.Context) {
	fileStorePath := filepath.Join(storeDir, context.Query("path"))
	_, err := os.Stat(fileStorePath)
	if err != nil {
		log.Println(err)
		response := ToolResponse{
			Error: "error processing file: " + fileStorePath,
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
		response := ToolResponse{
			Error: "error executing Magika command",
		}
		context.JSON(http.StatusOK, response)
		return
	}
	var data []Data
	err = json.Unmarshal(magikaOutput, &data)
	if err != nil {
		log.Println(err)
		response := ToolResponse{
			ToolOutput:   magikaOutputString,
			OutputFormat: "json",
			Error:        "unable to parse Magika JSON output",
		}
		context.JSON(http.StatusOK, response)
		return
	}
	// one file is analyzed at a time -> only first result is relevant
	features := extractFeatures(data[0])
	response := ToolResponse{
		ToolOutput:   magikaOutputString,
		OutputFormat: "json",
		Features:     features,
	}
	context.JSON(http.StatusOK, response)
}

func extractFeatures(data Data) map[string]string {
	features := make(map[string]string)
	if data.Result.Status == "ok" {
		if data.Result.Value.Output.MimeType != "" {
			features["mimeType"] = data.Result.Value.Output.MimeType
		}
		features["isText"] = strconv.FormatBool(data.Result.Value.Output.IsText)
	}
	return features
}
