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

type JhoveOutput struct {
	Root *JhoveRoot `json:"jhove"`
}

type JhoveRoot struct {
	ToolVersion string         `json:"release"`
	RepInfo     []JhoveRepInfo `json:"repInfo"`
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

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("", getDefaultResponse)
	router.GET("validate/:module", validateFile)
	router.Run()
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
		errorMessage := fmt.Sprintf("error processing file: %s", fileStorePath)
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
		"json",
		fileStorePath,
	)
	jhoveOutput, err := cmd.CombinedOutput()
	if err != nil {
		errorMessage := fmt.Sprintf("error executing JHOVE command: %s", string(jhoveOutput))
		log.Println(errorMessage)
		log.Println(err)
		response := ToolResponse{
			Error: &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	jhoveOutputString := string(jhoveOutput)
	processJhoveOutput(context, jhoveOutputString, module)
}

func processJhoveOutput(context *gin.Context, output string, module string) {
	var parsedJhoveOutput JhoveOutput
	err := json.NewDecoder(strings.NewReader(output)).Decode(&parsedJhoveOutput)
	if err != nil {
		errorMessage := "unable parse JHOVE output"
		log.Println(errorMessage)
		log.Println(err)
		response := ToolResponse{
			ToolVersion:  parsedJhoveOutput.Root.ToolVersion,
			ToolOutput:   output,
			OutputFormat: "text",
			Error:        &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	extractedFeatures := make(map[string]interface{})
	response := ToolResponse{
		ToolVersion:  parsedJhoveOutput.Root.ToolVersion,
		ToolOutput:   output,
		OutputFormat: "json",
		Features:     extractedFeatures,
	}
	if parsedJhoveOutput.Root != nil && len(parsedJhoveOutput.Root.RepInfo) > 0 {
		repInfo := parsedJhoveOutput.Root.RepInfo[0]
		if repInfo.FormatName != nil {
			extractedFeatures["formatName"] = *repInfo.FormatName
		}
		if repInfo.FormatVersion != nil {
			if module == "html" {
				// The html module of JHOVE adds the prefix HTML to format version.
				// This behavior makes it harder to compare to other tool results.
				r := regexp.MustCompile(`HTML ([0-9]+\.[0-9]+)`)
				matches := r.FindStringSubmatch(*repInfo.FormatVersion)
				if len(matches) == 2 {
					extractedFeatures["formatVersion"] = matches[1]
				}
			} else {
				extractedFeatures["formatVersion"] = *repInfo.FormatVersion
			}
		}
		if repInfo.Validation != nil {
			extractedFeatures["wellFormed"] = wellFormedRegEx.MatchString(*repInfo.Validation)
			extractedFeatures["valid"] = validRegEx.MatchString(*repInfo.Validation)
		}
		switch module {
		case "pdf":
			extractedFeatures["mimeType"] = "application/pdf"
			version, ok := extractedFeatures["formatVersion"]
			if ok {
				versionString, ok := version.(string)
				if ok {
					switch versionString {
					case "1.0":
						extractedFeatures["puid"] = "fmt/14"
					case "1.1":
						extractedFeatures["puid"] = "fmt/15"
					case "1.2":
						extractedFeatures["puid"] = "fmt/16"
					case "1.3":
						extractedFeatures["puid"] = "fmt/17"
					case "1.4":
						extractedFeatures["puid"] = "fmt/18"
					case "1.5":
						extractedFeatures["puid"] = "fmt/19"
					case "1.6":
						extractedFeatures["puid"] = "fmt/20"
					case "1.7":
						extractedFeatures["puid"] = "fmt/276"
					}
				}
			}
		case "html":
			extractedFeatures["mimeType"] = "text/html"
			version, ok := extractedFeatures["formatVersion"]
			if ok {
				versionString, ok := version.(string)
				if ok {
					switch versionString {
					case "3.2":
						extractedFeatures["puid"] = "fmt/98"
					case "4.0":
						extractedFeatures["puid"] = "fmt/99"
					case "4.01":
						extractedFeatures["puid"] = "fmt/100"
					}
				}
			}
		case "tiff":
			extractedFeatures["mimeType"] = "image/tiff"
		case "jpeg":
			extractedFeatures["mimeType"] = "image/jpeg"
			version, ok := extractedFeatures["formatVersion"]
			if ok {
				versionString, ok := version.(string)
				if ok {
					switch versionString {
					case "1.00":
						extractedFeatures["puid"] = "fmt/42"
					case "1.01":
						extractedFeatures["puid"] = "fmt/43"
					}
				}
			}
		case "jpeg2000":
			extractedFeatures["mimeType"] = "image/jp2"
			extractedFeatures["puid"] = "x-fmt/392"
		}
	}
	context.JSON(http.StatusOK, response)
}
