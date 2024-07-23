package main

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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
	router.Run("0.0.0.0:80")
}

func getDefaultResponse(context *gin.Context) {
	context.String(http.StatusOK, defaultResponse)
}

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
	log.Println(droidOutputString)
	csvReader := csv.NewReader(strings.NewReader(droidOutputString))
	formats, err := csvReader.ReadAll()
	if err != nil {
		log.Println(err.Error())
		response := ToolResponse{
			ToolOutput:   droidOutputString,
			OutputFormat: "text",
			Error:        "unable to parse DROID csv output",
		}
		context.JSON(http.StatusOK, response)
		return
	}
	if len(formats) == 0 || len(formats[1]) < 18 {
		response := ToolResponse{
			ToolOutput:   droidOutputString,
			OutputFormat: "csv",
			Error:        "unable to parse DROID csv output",
		}
		context.JSON(http.StatusOK, response)
		return
	}
	features := make(map[string]string)
	if formats[1][14] != "" {
		features["puid"] = formats[1][14]
	}
	if formats[1][15] != "" {
		features["mimeType"] = formats[1][15]
	}
	if formats[1][16] != "" {
		features["formatName"] = formats[1][16]
	}
	if formats[1][17] != "" {
		features["formatVersion"] = formats[1][17]
	}
	response := ToolResponse{
		ToolOutput:   droidOutputString,
		OutputFormat: "csv",
		Features:     features,
	}
	context.JSON(http.StatusOK, response)
}
