package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ToolResponse struct {
	ToolOutput        string
	OutputFormat      string
	ExtractedFeatures map[string]string
}

type JhoveOutput struct {
	Root *JhoveRoot `json:"jhove"`
}

type JhoveRoot struct {
	RepInfo []JhoveRepInfo `json:"repInfo"`
}

type JhoveRepInfo struct {
	FormatName    *string `json:"format"`
	FormatVersion *string `json:"version"`
	Validation    *string `json:"status"`
}

var defaultResponse = "JHOVE API is running"
var storeDir = "/borg/filestore"
var wellFormedRegEx = regexp.MustCompile("Well-Formed")
var validRegEx = regexp.MustCompile("Well-Formed and valid")

func main() {
	router := gin.Default()
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"*"})
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type"}
	corsConfig.AllowMethods = []string{"GET", "POST"}
	// It's important that the cors configuration is used before declaring the routes.
	router.Use(cors.New(corsConfig))
	router.GET("", getDefaultResponse)
	router.GET("validate/:module", validateFile)
	addr := "0.0.0.0:" + os.Getenv("JHOVE_API_CONTAINER_PORT")
	router.Run(addr)
}

func getDefaultResponse(context *gin.Context) {
	context.String(http.StatusOK, defaultResponse)
}

func validateFile(context *gin.Context) {
	module := context.Param("module")
	if module == "" {
		errorMessage := "no JHOVE module declared"
		log.Println(errorMessage)
		context.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": errorMessage,
		})
		return
	}
	fileStorePath := filepath.Join(storeDir, context.Query("path"))
	_, err := os.Stat(fileStorePath)
	if err != nil {
		log.Println(err)
		context.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": "error processing file: " + fileStorePath,
		})
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
	jhoveOutput, err := cmd.Output()
	if err != nil {
		log.Println(err)
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "error executing JHOVE command",
		})
		return
	}
	jhoveOutputString := string(jhoveOutput)
	processJhoveOutput(context, jhoveOutputString)
}

func processJhoveOutput(context *gin.Context, output string) {
	var parsedJhoveOutput JhoveOutput
	err := json.NewDecoder(strings.NewReader(output)).Decode(&parsedJhoveOutput)
	if err != nil {
		log.Println(err)
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "unable parse JHOVE output",
		})
		return
	}
	extractedFeatures := make(map[string]string)
	response := ToolResponse{
		ToolOutput:        output,
		OutputFormat:      "json",
		ExtractedFeatures: extractedFeatures,
	}
	if parsedJhoveOutput.Root != nil && len(parsedJhoveOutput.Root.RepInfo) > 0 {
		repInfo := parsedJhoveOutput.Root.RepInfo[0]
		if repInfo.FormatVersion != nil {
			// if format name was extracted and version doesn't contain it already
			if repInfo.FormatName != nil &&
				!strings.Contains(*repInfo.FormatVersion, *repInfo.FormatName) {
				extractedFeatures["formatVersion"] = *repInfo.FormatName + " " + *repInfo.FormatVersion
			} else {
				extractedFeatures["formatVersion"] = *repInfo.FormatVersion
			}
		}
		if repInfo.Validation != nil {
			extractedFeatures["wellFormed"] = strconv.FormatBool(
				wellFormedRegEx.MatchString(*repInfo.Validation),
			)
			extractedFeatures["valid"] = strconv.FormatBool(
				validRegEx.MatchString(*repInfo.Validation),
			)
		}
	}
	context.JSON(http.StatusOK, response)
}
