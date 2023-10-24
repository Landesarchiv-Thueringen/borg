package main

import (
	"encoding/json"
	"lath/borg/internal/config"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
	context.JSON(http.StatusOK, identificationResults)
}

func runFileIdentificationTools(fileName string) []ToolResponse {
	var results []ToolResponse
	// for every identification tool
	for _, tool := range serverConfig.FormatIdentificationTools {
		toolResponse := getToolResponse(tool.ToolName, tool.ToolVersion, tool.Endpoint, fileName)
		results = append(results, toolResponse)
	}
	return results
}

func getToolResponse(
	toolName string,
	toolVersion string,
	endpoint string,
	fileName string,
) ToolResponse {
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
		return toolResponse
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
		return toolResponse
	}
	// process request response
	processToolResponse(response, &toolResponse)
	return toolResponse
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
