package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ToolResponse struct {
	ToolOutput   string            `json:"toolOutput"`
	OutputFormat string            `json:"outputFormat"`
	Features     map[string]string `json:"features"`
	Error        string            `json:"error"`
}

const defaultResponse = "DROID API is running"
const workDir = "/borg/tools/droid"
const storeDir = "/borg/file-store"

var signatureFilePath = filepath.Join(workDir, "third_party/DROID_SignatureFile_V114.xml")
var containerSignatureFilePath = filepath.Join(workDir, "third_party/container-signature-20230822.xml")

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

// identifyFileFormat executes DROID and parses the output of the command.
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
		"/bin/ash",
		"/borg/tools/droid/third_party/droid.sh",
		"-Ns",
		signatureFilePath,
		"-Nc",
		containerSignatureFilePath,
		fileStorePath,
	)
	droidOutput, err := cmd.Output()
	if err != nil {
		log.Println(err)
		response := ToolResponse{
			Error: "error executing DROID command",
		}
		context.JSON(http.StatusOK, response)
		return
	}
	droidOutputString := string(droidOutput)
	csvReader := csv.NewReader(strings.NewReader(droidOutputString))
	// Deactivate field number check per row.
	// DROID can return invalid csv.
	csvReader.FieldsPerRecord = -1
	formatTable, err := csvReader.ReadAll()
	if err != nil || len(formatTable) < 2 {
		if err != nil {
			log.Println(err.Error())
		}
		response := ToolResponse{
			ToolOutput:   droidOutputString,
			OutputFormat: "csv",
			Error:        "unable to parse DROID csv output",
		}
		context.JSON(http.StatusOK, response)
		return
	}
	features, err := extractFeatures(formatTable)
	if err != nil {
		log.Println(err.Error())
		response := ToolResponse{
			ToolOutput:   droidOutputString,
			OutputFormat: "csv",
			Error:        "unable to parse DROID csv output",
		}
		context.JSON(http.StatusOK, response)
		return
	}
	response := ToolResponse{
		ToolOutput:   droidOutputString,
		OutputFormat: "csv",
		Features:     features,
	}
	context.JSON(http.StatusOK, response)
}

// extractFeatures extracts all relevant information from parsed DROID output.
// Extracts only features from the first detected format.
func extractFeatures(formatTable [][]string) (map[string]string, error) {
	features := make(map[string]string)
	keyMap := getKeyMap(formatTable[0])
	formatNumberAsString, err := extractFeature("FORMAT_COUNT", formatTable[1], keyMap)
	// key and value errors prevent further processing
	if err != nil {
		return features, fmt.Errorf("extractFeatures: unexpected csv layout: %w", err)
	}
	formatNumber, err := strconv.Atoi(formatNumberAsString)
	if err != nil {
		return features, fmt.Errorf("extractFeatures: failed to extract format number: %w", err)
	}
	// if no formats were identified
	if formatNumber == 0 {
		return features, nil
	}
	// extract the relevant features
	// only key errors prevent further processing
	// value errors are expected, not all features exist for all files
	var keyError *KeyNotFoundError
	// PUID
	puid, err := extractFeature("PUID", formatTable[1], keyMap)
	if err == nil {
		if puid != "" {
			features["puid"] = puid
		}
	} else if errors.As(err, &keyError) {
		return features, fmt.Errorf("extractFeatures: unexpected csv layout: %w", keyError)
	}
	// MIME type
	mimeType, err := extractFeature("MIME_TYPE", formatTable[1], keyMap)
	if err == nil {
		if mimeType != "" {
			features["mimeType"] = mimeType
		}
	} else if errors.As(err, &keyError) {
		return features, fmt.Errorf("extractFeatures: unexpected csv layout: %w", keyError)
	}
	// format name
	formatName, err := extractFeature("FORMAT_NAME", formatTable[1], keyMap)
	if err == nil {
		if formatName != "" {
			features["formatName"] = formatName
		}
	} else if errors.As(err, &keyError) {
		return features, fmt.Errorf("extractFeatures: unexpected csv layout: %w", keyError)
	}
	// format version
	formatVersion, err := extractFeature("FORMAT_VERSION", formatTable[1], keyMap)
	if err == nil {
		if formatVersion != "" {
			features["formatVersion"] = formatVersion
			// add prefix to format version if format name contains PDF/A
			if strings.Contains(formatName, "PDF/A") {
				features["formatVersion"] = "PDF/A-" + features["formatVersion"]
			}
		}

	} else if errors.As(err, &keyError) {
		return features, fmt.Errorf("extractFeatures: unexpected csv layout: %w", keyError)
	}
	return features, nil
}

// getKeyMap maps the column keys to their index.
func getKeyMap(header []string) map[string]int {
	m := make(map[string]int)
	for index, columnHeader := range header {
		m[columnHeader] = index
	}
	return m
}

// KeyNotFoundError represents the absence of an expected key.
type KeyNotFoundError struct {
	key string
}

func (e *KeyNotFoundError) Error() string {
	return fmt.Sprintf("key [%q] does not exist", e.key)
}

// KeyNotFoundError represents the absence of an value foe an existing key.
type ValueNotFoundError struct {
	key string
}

func (e *ValueNotFoundError) Error() string {
	return fmt.Sprintf("value for key [%q] does not exist", e.key)
}

// extractFeature tries to extract feature with give key.
func extractFeature(key string, formatRow []string, keyMap map[string]int) (string, error) {
	valueIndex, ok := keyMap[key]
	if !ok {
		return "", &KeyNotFoundError{key: key}
	}
	if valueIndex >= len(formatRow) {
		return "", &ValueNotFoundError{key: key}
	}
	return formatRow[valueIndex], nil
}
