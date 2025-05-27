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

var (
	FORMAT_NAME_LABEL    = "Formatname"
	FORMAT_VERSION_LABEL = "Formatversion"
	MIME_TYPE_LABEL      = "Mime-Type"
	PUID_LABEL           = "PUID"
	VALID_LABEL          = "valide"
	WELL_FORMED_LABEL    = "wohlgeformt"
)

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
			ToolOutput:   output,
			OutputFormat: "text",
			Error:        &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	extractedFeatures := make(map[string]ToolFeatureValue)
	response := ToolResponse{
		ToolVersion:  parsedJhoveOutput.Root.ToolVersion,
		ToolOutput:   output,
		OutputFormat: "json",
		Features:     extractedFeatures,
	}
	if parsedJhoveOutput.Root != nil && len(parsedJhoveOutput.Root.RepInfo) > 0 {
		repInfo := parsedJhoveOutput.Root.RepInfo[0]
		if repInfo.FormatName != nil {
			extractedFeatures["format:name"] = getFormatNameFeature(*repInfo.FormatName)
		}
		if repInfo.FormatVersion != nil {
			if module == "html" {
				// The html module of JHOVE adds the prefix HTML to format version.
				// This behavior makes it harder to compare to other tool results.
				r := regexp.MustCompile(`HTML ([0-9]+\.[0-9]+)`)
				matches := r.FindStringSubmatch(*repInfo.FormatVersion)
				if len(matches) == 2 {
					extractedFeatures["format:version"] = getVersionFeature(matches[1])
				}
			} else {
				extractedFeatures["format:version"] = getVersionFeature(*repInfo.FormatVersion)
			}
		}
		if repInfo.Validation != nil {
			extractedFeatures["format:wellFormed"] = getWellFormedFeature(wellFormedRegEx.MatchString(*repInfo.Validation))
			extractedFeatures["format:valid"] = getValidFeature(validRegEx.MatchString(*repInfo.Validation))
		}
		switch module {
		case "pdf":
			extractedFeatures["format:mimeType"] = ToolFeatureValue{
				Value: "application/pdf",
				Label: &MIME_TYPE_LABEL,
			}
			version, ok := extractedFeatures["format:version"]
			if ok {
				versionString, ok := version.Value.(string)
				if ok {
					switch versionString {
					case "1.0":
						extractedFeatures["format:puid"] = getPuidFeature("fmt/14")
					case "1.1":
						extractedFeatures["format:puid"] = getPuidFeature("fmt/15")
					case "1.2":
						extractedFeatures["format:puid"] = getPuidFeature("fmt/16")
					case "1.3":
						extractedFeatures["format:puid"] = getPuidFeature("fmt/17")
					case "1.4":
						extractedFeatures["format:puid"] = getPuidFeature("fmt/18")
					case "1.5":
						extractedFeatures["format:puid"] = getPuidFeature("fmt/19")
					case "1.6":
						extractedFeatures["format:puid"] = getPuidFeature("fmt/20")
					case "1.7":
						extractedFeatures["format:puid"] = getPuidFeature("fmt/276")
					}
				}
			}
		case "html":
			extractedFeatures["format:mimeType"] = getMimeTypeFeature("text/html")
			version, ok := extractedFeatures["format:version"]
			if ok {
				versionString, ok := version.Value.(string)
				if ok {
					switch versionString {
					case "3.2":
						extractedFeatures["format:puid"] = getPuidFeature("fmt/98")
					case "4.0":
						extractedFeatures["format:puid"] = getPuidFeature("fmt/99")
					case "4.01":
						extractedFeatures["format:puid"] = getPuidFeature("fmt/100")
					}
				}
			}
		case "tiff":
			extractedFeatures["format:mimeType"] = getMimeTypeFeature("image/tiff")
		case "jpeg":
			extractedFeatures["format:mimeType"] = getMimeTypeFeature("image/jpeg")
			version, ok := extractedFeatures["format:version"]
			if ok {
				versionString, ok := version.Value.(string)
				if ok {
					switch versionString {
					case "1.00":
						extractedFeatures["format:puid"] = getPuidFeature("fmt/42")
					case "1.01":
						extractedFeatures["format:puid"] = getPuidFeature("fmt/43")
					}
				}
			}
		case "jpeg2000":
			extractedFeatures["format:mimeType"] = getMimeTypeFeature("image/jp2")
			extractedFeatures["format:puid"] = getPuidFeature("x-fmt/392")
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

func getFormatNameFeature(value string) ToolFeatureValue {
	return ToolFeatureValue{
		Value: value,
		Label: &FORMAT_NAME_LABEL,
	}
}

func getWellFormedFeature(value bool) ToolFeatureValue {
	return ToolFeatureValue{
		Value: value,
		Label: &WELL_FORMED_LABEL,
	}
}

func getValidFeature(value bool) ToolFeatureValue {
	return ToolFeatureValue{
		Value: value,
		Label: &VALID_LABEL,
	}
}
