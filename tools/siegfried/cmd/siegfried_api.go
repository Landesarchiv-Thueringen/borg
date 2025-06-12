package main

import (
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

type SiegfriedResult struct {
	Version     string       `json:"siegfried"`
	FileResults []FileResult `json:"files"`
}

type FileResult struct {
	IdentMatches []IdentMatch `json:"matches"`
}

type IdentMatch struct {
	NameSpace     string `json:"ns"`
	Id            string `json:"id"`
	FormatName    string `json:"format"`
	FormatVersion string `json:"version"`
	MimeType      string `json:"mime"`
}

var (
	FORMAT_NAME_LABEL    = "Formatname"
	FORMAT_VERSION_LABEL = "Formatversion"
	MIME_TYPE_LABEL      = "Mime-Type"
	PUID_LABEL           = "PUID"
)

const (
	DEFAULT_RESPONSE = "Siegfried API is running"
	WORK_DIR         = "/borg/tools/siegfried"
	STORE_DIR        = "/borg/file-store"
	TIME_OUT         = 30 * time.Second
)

var toolVersion string

func main() {
	toolVersion = getToolVersion()
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("", getDefaultResponse)
	router.GET("/identify", identifyFileFormat)
	router.Run()
}

func getToolVersion() string {
	cmd := exec.Command(
		"./third_party/sf",
		"-v",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	outputString := string(output)
	regEx := regexp.MustCompile(`siegfried\s*(.+)`)
	matches := regEx.FindStringSubmatch(outputString)
	if len(matches) != 2 {
		log.Fatal("couldn't extract Siegfried version from tool output")
	}
	return matches[1]
}

func getDefaultResponse(context *gin.Context) {
	context.String(http.StatusOK, DEFAULT_RESPONSE)
}

func identifyFileFormat(ginContext *gin.Context) {
	fileStorePath := filepath.Join(STORE_DIR, ginContext.Query("path"))
	_, err := os.Stat(fileStorePath)
	if err != nil {
		log.Println(err)
		errorMessage := fmt.Sprintf("error processing file: %s", fileStorePath)
		response := ToolResponse{
			ToolVersion: toolVersion,
			Error:       &errorMessage,
		}
		ginContext.JSON(http.StatusOK, response)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT)
	defer cancel()
	cmd := exec.CommandContext(
		ctx,
		"./third_party/sf",
		"-json",
		fileStorePath,
	)
	output, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		errorMessage := fmt.Sprintf("Timeout exceeded after %s.", TIME_OUT)
		response := ToolResponse{
			ToolVersion: toolVersion,
			Error:       &errorMessage,
		}
		ginContext.JSON(http.StatusOK, response)
		return
	}
	outputString := string(output)
	if err != nil {
		log.Println(err)
		errorMessage := fmt.Sprintf(
			"error executing Siegfried command: %s", outputString)
		response := ToolResponse{
			ToolVersion: toolVersion,
			Error:       &errorMessage,
		}
		ginContext.JSON(http.StatusOK, response)
		return
	}
	var result SiegfriedResult
	err = json.Unmarshal(output, &result)
	if err != nil {
		errorMessage := err.Error()
		response := ToolResponse{
			ToolVersion: toolVersion,
			Error:       &errorMessage,
		}
		ginContext.JSON(http.StatusOK, response)
		return
	}
	features := make(map[string]ToolFeatureValue)
	if len(result.FileResults) > 0 && len(result.FileResults[0].IdentMatches) > 0 {
		match := result.FileResults[0].IdentMatches[0]
		if match.NameSpace == "pronom" {
			if match.Id != "UNKNOWN" {
				features["format:puid"] = ToolFeatureValue{
					Value: match.Id,
					Label: &PUID_LABEL,
				}
			}
			if len(match.MimeType) > 0 {
				features["format:mimeType"] = ToolFeatureValue{
					Value: match.MimeType,
					Label: &MIME_TYPE_LABEL,
				}
			}
			if len(match.FormatVersion) > 0 {
				version := match.FormatVersion
				// add prefix to format version if format name contains PDF/A
				if strings.Contains(match.FormatName, "PDF/A") {
					version = "PDF/A-" + version
				}
				features["format:version"] = ToolFeatureValue{
					Value: version,
					Label: &FORMAT_VERSION_LABEL,
				}
			}
			if len(match.FormatName) > 0 {
				features["format:name"] = ToolFeatureValue{
					Value: match.FormatName,
					Label: &FORMAT_NAME_LABEL,
				}
			}
		}
	}
	response := ToolResponse{
		ToolVersion:  result.Version,
		ToolOutput:   outputString,
		OutputFormat: "json",
		Features:     features,
	}
	ginContext.JSON(http.StatusOK, response)
}
