package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
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

const (
	defaultResponse = "MediaInfo API is running"
	workDir         = "/borg/tools/magika"
	storeDir        = "/borg/file-store"
)

var (
	toolVersion string
	dict        map[string]string
)

func main() {
	toolVersion = getToolVersion()
	dict = readLocalizationCsv()
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("", getDefaultResponse)
	router.GET("/localization-dict", getLocalizationDict)
	router.GET("/extract-metadata", extractMetadata)
	router.Run()
}

func getToolVersion() string {
	cmd := exec.Command(
		"mediainfo",
		"--Version",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	outputString := string(output)
	regEx := regexp.MustCompile(`MediaInfoLib\s*-\s*(.+)`)
	matches := regEx.FindStringSubmatch(outputString)
	if len(matches) != 2 {
		log.Fatal("couldn't extract MediaInfo version from tool output")
	}
	return matches[1]
}

func getDefaultResponse(context *gin.Context) {
	context.String(http.StatusOK, defaultResponse)
}

func getLocalizationDict(context *gin.Context) {
	context.JSON(http.StatusOK, dict)
}

func extractMetadata(context *gin.Context) {
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
		"mediainfo",
		"--Full",
		"--Output=XML",
		fileStorePath,
	)
	output, err := cmd.CombinedOutput()
	outputString := string(output)
	if err != nil {
		log.Println(err)
		errorMessage := fmt.Sprintf(
			"error executing MediaInfo command: %s", outputString)
		response := ToolResponse{
			ToolVersion: toolVersion,
			Error:       &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	features, err := extractFeatures(outputString)
	if err != nil {
		errorMessage := err.Error()
		response := ToolResponse{
			ToolVersion: toolVersion,
			Error:       &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	mimeType, ok := features["av_container:internetmediatype"]
	if ok {
		features["format:mimeType"] = mimeType
	}
	response := ToolResponse{
		ToolVersion:  toolVersion,
		ToolOutput:   outputString,
		OutputFormat: "xml",
		Features:     features,
	}
	context.JSON(http.StatusOK, response)
}

func extractFeatures(output string) (map[string]ToolFeatureValue, error) {
	features := make(map[string]ToolFeatureValue)
	decoder := xml.NewDecoder(strings.NewReader(output))
CATEGORY_LOOP:
	for {
		token, err := decoder.Token()
		if err != nil {
			if err != io.EOF {
				log.Println(err)
				return features, err
			}
			break
		}
		startElement, ok := token.(xml.StartElement)
		if ok && startElement.Name.Local == "track" {
			var category string
			for _, attr := range startElement.Attr {
				if attr.Name.Local == "type" {
					category = attr.Value
					break
				}
			}
			if category == "General" {
				category = "av_container"
			}
			for {
				innerToken, err := decoder.Token()
				if err != nil {
					break
				}
				switch innerElement := innerToken.(type) {
				case xml.StartElement:
					value, ok := getInnerValue(decoder, innerElement)
					if !ok {
						continue
					}
					featureKey := strings.ToLower(innerElement.Name.Local)
					key := fmt.Sprintf(
						"%s:%s",
						strings.ToLower(category),
						featureKey,
					)
					features[key] = getFeatureValue(featureKey, value)
				case xml.EndElement:
					if innerElement.Name.Local == "track" {
						continue CATEGORY_LOOP
					}
				}
			}
		}
	}
	return features, nil
}

// getInnerValue returns the value of a flat xml element.
// ok is false if the element has child elements.
// Consumes the whole xml element.
func getInnerValue(decoder *xml.Decoder, element xml.StartElement) (value string, ok bool) {
	ok = true
	tagName := element.Name.Local
	for {
		nextToken, err := decoder.Token()
		if err != nil {
			ok = false
			break
		}
		switch nextElement := nextToken.(type) {
		case xml.StartElement:
			ok = false
		case xml.EndElement:
			if nextElement.Name.Local == tagName {
				return
			}
		case xml.CharData:
			if ok {
				value = string(nextElement)
			}
		}
	}
	return
}

func readLocalizationCsv() map[string]string {
	file, err := os.Open("de.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	dict := make(map[string]string)
	csvReader := csv.NewReader(file)
	csvReader.Comma = ';'
	for {
		row, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				return dict
			}
			continue
		}
		if len(row) != 2 {
			continue
		}
		key := strings.ToLower(row[0])
		dict[key] = row[1]
	}
}

var keyRegEx = regexp.MustCompile(`^(.+?)_(string)(\d*)$`)

func getFeatureValue(key string, value string) ToolFeatureValue {
	// strings.LastIndex(key, "_string")
	matches := keyRegEx.FindStringSubmatch(key)
	if len(matches) > 2 {
		subKey := matches[1]
		label, ok := dict[subKey]
		if ok {
			if len(matches) == 4 && len(matches[3]) == 0 {
				if label == "" {
					log.Println(key)
				}
				label += " (string)"
			} else if len(matches) == 4 && len(matches[3]) > 0 {
				label += fmt.Sprintf(" (string %s)", matches[3])
			}
			return ToolFeatureValue{
				Value: value,
				Label: &label,
			}
		}
	}
	label, ok := dict[key]
	if ok && len(label) > 0 {
		return ToolFeatureValue{
			Value: value,
			Label: &label,
		}
	}
	return ToolFeatureValue{
		Value: value,
	}
}
