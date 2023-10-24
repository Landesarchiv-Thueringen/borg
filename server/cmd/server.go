package main

import (
	"encoding/json"
	"lath/borg/internal/config"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FileAnalysis struct {
	FileIdentificationResults []ToolResponse `json:"fileIdentificationResults"`
	FileValidationResults     []ToolResponse `json:"fileValidationResults"`
}

type ToolResponse struct {
	ToolName          string             `json:"toolName"`
	ToolVersion       string             `json:"toolVersion"`
	ToolOutput        *string            `json:"toolOutput"`
	ExtractedFeatures *map[string]string `json:"extractedFeatures"`
	Error             *string            `json:"error"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type Feature struct {
	Key   string  `json:"key"`
	Value string  `json:"value"`
	Score float64 `json:"score"`
}

var defaultResponse = "borg server is running"
var storePath = "/borg/filestore"
var serverConfig config.ServerConfig

func main() {
	serverConfig = config.ParseConfig()
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
	router.POST("analyse-file", analyseFile)
	addr := "0.0.0.0:" + os.Getenv("BORG_SERVER_CONTAINER_PORT")
	router.Run(addr)
}

func getDefaultResponse(context *gin.Context) {
	context.String(http.StatusOK, defaultResponse)
}

func analyseFile(context *gin.Context) {
	file, err := context.FormFile("file")
	// no file received
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "no file received",
		})
		return
	}
	// generate unique file name for storing
	fileName := uuid.New().String() + "_" + file.Filename
	fileStorePath := filepath.Join(storePath, fileName)
	err = context.SaveUploadedFile(file, fileStorePath)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "unable to save file",
		})
		return
	}
	defer os.Remove(fileStorePath)
	identificationResults := runFileIdentificationTools(fileName)
	validationResults := runFileValidationTools(fileName, identificationResults)
	fileAnalysis := FileAnalysis{
		FileIdentificationResults: identificationResults,
		FileValidationResults:     validationResults,
	}
	context.JSON(http.StatusOK, fileAnalysis)
}

func runFileIdentificationTools(fileName string) []ToolResponse {
	var responseChannels []chan ToolResponse
	// for every identification tool
	for _, tool := range serverConfig.FormatIdentificationTools {
		rc := make(chan ToolResponse)
		responseChannels = append(responseChannels, rc)
		// request tool results concurrent
		go getToolResponse(tool.ToolName, tool.ToolVersion, tool.Endpoint, fileName, rc)
	}
	// gather all tool responses
	var results []ToolResponse
	for _, rc := range responseChannels {
		toolResponse := <-rc
		results = append(results, toolResponse)
	}
	return results
}

func runFileValidationTools(fileName string, identificationResults []ToolResponse) []ToolResponse {
	var responseChannels []chan ToolResponse
	// for every validation tool
	for _, tool := range serverConfig.FormatValidationTools {
		// for every possible trigger of current validation tool
		for _, trigger := range tool.ToolTrigger {
			if checkToolTrigger(trigger, identificationResults) {
				rc := make(chan ToolResponse)
				responseChannels = append(responseChannels, rc)
				// request tool results concurrent
				go getToolResponse(tool.ToolName, tool.ToolVersion, tool.Endpoint, fileName, rc)
				// don't check other triggers, tool response already requested
				break
			}
		}
	}
	// gather all tool responses
	var results []ToolResponse
	for _, rc := range responseChannels {
		toolResponse := <-rc
		results = append(results, toolResponse)
	}
	return results
}

// returns true if the trigger fires
func checkToolTrigger(trigger config.ToolTrigger, identificationResults []ToolResponse) bool {
	regex := regexp.MustCompile(trigger.RegEx)
	for _, toolResponse := range identificationResults {
		if toolResponse.ExtractedFeatures != nil {
			features := *toolResponse.ExtractedFeatures
			featureValue, ok := features[trigger.Feature]
			if ok && regex.MatchString(featureValue) {
				return true
			}
		}
	}
	return false
}

func getToolResponse(
	toolName string,
	toolVersion string,
	endpoint string,
	fileName string,
	rc chan ToolResponse,
) {
	toolResponse := ToolResponse{
		ToolName:    toolName,
		ToolVersion: toolVersion,
	}
	// create http get request
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Println(err)
		errorMessage := "error creating request: " + endpoint
		toolResponse.Error = &errorMessage
		rc <- toolResponse
		return
	}
	// add file path URL parameter
	query := req.URL.Query()
	query.Add("path", fileName)
	req.URL.RawQuery = query.Encode()
	// send get request
	response, err := http.Get(req.URL.String())
	if err != nil {
		log.Println(err)
		errorMessage := "error requesting: " + req.URL.String()
		toolResponse.Error = &errorMessage
		rc <- toolResponse
		return
	}
	// process request response
	processToolResponse(response, &toolResponse)
	rc <- toolResponse
}

func processToolResponse(response *http.Response, toolResponse *ToolResponse) {
	// identification tool request was successful
	if response.StatusCode == http.StatusOK {
		var parsedResponse ToolResponse
		err := json.NewDecoder(response.Body).Decode(&parsedResponse)
		if err != nil {
			log.Println(err)
			errorMessage := "error parsing tool response"
			toolResponse.Error = &errorMessage
		} else {
			toolResponse.ToolOutput = parsedResponse.ToolOutput
			toolResponse.ExtractedFeatures = parsedResponse.ExtractedFeatures
		}
	} else {
		// error occurred during file identification
		var errorResponse ErrorResponse
		err := json.NewDecoder(response.Body).Decode(&errorResponse)
		if err != nil {
			log.Println(err)
			errorMessage := "error parsing tool error response"
			toolResponse.Error = &errorMessage
		} else {
			toolResponse.Error = &errorResponse.Message
		}
	}
}
