package main

import (
	"encoding/json"
	"lath/borg/internal/config"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FileAnalysis struct {
	Summary                   map[string]Feature `json:"summary"`
	FileIdentificationResults []ToolResponse     `json:"fileIdentificationResults"`
	FileValidationResults     []ToolResponse     `json:"fileValidationResults"`
}

type ToolResponse struct {
	ToolName          string                 `json:"toolName"`
	ToolVersion       string                 `json:"toolVersion"`
	FeatureConfig     []config.FeatureConfig `json:"-"`
	ToolOutput        *string                `json:"toolOutput"`
	OutputFormat      *string                `json:"outputFormat"`
	ExtractedFeatures *map[string]string     `json:"extractedFeatures"`
	Error             *string                `json:"error"`
}

type ToolConfidence struct {
	ToolName      string                 `json:"toolName"`
	Confidence    float64                `json:"confidence"`
	FeatureConfig []config.FeatureConfig `json:"-"`
}

type FeatureValue struct {
	Value string           `json:"value"`
	Score float64          `json:"score"`
	Tools []ToolConfidence `json:"tools"`
}

type Feature struct {
	Key    string         `json:"key"`
	Values []FeatureValue `json:"values"`
}

// implement sorting interface for feature values
type ByScore []FeatureValue

func (featureValues ByScore) Len() int {
	return len(featureValues)
}

func (featureValues ByScore) Less(i, j int) bool {
	// sort reversed, biggest score first
	return featureValues[i].Score > featureValues[j].Score
}

func (featureValues ByScore) Swap(i, j int) {
	featureValues[i], featureValues[j] = featureValues[j], featureValues[i]
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
	summary := summarizeToolResults(identificationResults, validationResults)
	fileAnalysis := FileAnalysis{
		Summary:                   summary,
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
		go getToolResponse(
			tool.ToolName,
			tool.ToolVersion,
			tool.Endpoint,
			tool.Features,
			fileName,
			rc,
		)
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
				go getToolResponse(
					tool.ToolName,
					tool.ToolVersion,
					tool.Endpoint,
					tool.Features,
					fileName,
					rc,
				)
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
	toolFeatureConfig []config.FeatureConfig,
	fileName string,
	rc chan ToolResponse,
) {
	toolResponse := ToolResponse{
		ToolName:      toolName,
		ToolVersion:   toolVersion,
		FeatureConfig: toolFeatureConfig,
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
	// tool request was successful
	if response.StatusCode == http.StatusOK {
		var parsedResponse ToolResponse
		err := json.NewDecoder(response.Body).Decode(&parsedResponse)
		if err != nil {
			errorMessage := "error parsing tool response"
			log.Println(errorMessage)
			log.Println(err)
			toolResponse.Error = &errorMessage
		} else {
			toolResponse.ToolOutput = parsedResponse.ToolOutput
			toolResponse.OutputFormat = parsedResponse.OutputFormat
			toolResponse.ExtractedFeatures = parsedResponse.ExtractedFeatures
			toolResponse.Error = parsedResponse.Error
		}
	} else {
		// tool request was not successful
		errorMessage := "tool request error"
		toolResponse.Error = &errorMessage
		bytes, err := httputil.DumpResponse(response, true)
		if err == nil {
			responseString := string(bytes)
			toolResponse.ToolOutput = &responseString
		}
	}
}

func summarizeToolResults(
	identificationResults []ToolResponse,
	validationResults []ToolResponse,
) map[string]Feature {
	tools := append(identificationResults, validationResults...)
	summary := make(map[string]Feature)
	// for every tool response
	for _, tool := range tools {
		// if no extracted features exist for tool
		if tool.ExtractedFeatures == nil {
			// continue with next tool
			continue
		}
		// for every extracted feature
		for featureKey, featureValue := range *tool.ExtractedFeatures {
			f, ok := summary[featureKey]
			// feature exists in summary
			if ok {
				// add current tool to feature
				valueExists := false
				for i, v := range f.Values {
					// extracted value exists already, that means another tool extracted the same value
					if v.Value == featureValue {
						// add tool to tools that extracted current value for feature
						v.Tools = append(v.Tools, getToolConfidence(tool, featureKey))
						// overwrite original tool list
						summary[featureKey].Values[i].Tools = v.Tools
						valueExists = true
						break
					}
				}
				// value for key doesn't exist already --> add it to value list
				if !valueExists {
					f.Values = append(f.Values, getFeatureValue(featureKey, featureValue, tool))
					// overwrite feature in summary map
					summary[featureKey] = f
				}
			} else {
				// feature doesn't exist already in summary --> add it to summary
				summary[featureKey] = getFeature(featureKey, featureValue, tool)
			}
		}
	}
	calculateFeatureValueScore(&summary)
	sortFeatureValues(&summary)
	// only after first score calculation, can global feature conditions be applied
	correctToolConfidence(&summary)
	calculateFeatureValueScore(&summary)
	sortFeatureValues(&summary)
	return summary
}

func getToolConfidence(tool ToolResponse, featureKey string) ToolConfidence {
	confidence := 1.0
	for _, featureConfig := range tool.FeatureConfig {
		if featureConfig.Key == featureKey {
			confidence = featureConfig.Confidence.DefaultValue
			break
		}
	}
	toolConfidence := ToolConfidence{
		ToolName:      tool.ToolName,
		Confidence:    confidence,
		FeatureConfig: tool.FeatureConfig,
	}
	return toolConfidence
}

func getCorrectedToolConfidence(
	featureKey string,
	toolConfidence ToolConfidence,
	scoredFeatures *map[string]Feature,
) ToolConfidence {
	for _, featureConfig := range toolConfidence.FeatureConfig {
		// if feature configuration doesn't belong to currentlu corrected feature
		if featureKey != featureConfig.Key {
			continue
		}
		if len(featureConfig.Confidence.Conditions) > 0 {
			for _, condition := range featureConfig.Confidence.Conditions {
				scoredFeature, ok := (*scoredFeatures)[condition.GlobalFeature]
				if ok {
					regex := regexp.MustCompile(condition.RegEx)
					// the first value has the highest score --> voted truth
					if regex.MatchString(scoredFeature.Values[0].Value) {
						toolConfidence.Confidence = condition.Value
						return toolConfidence
					}
				}
			}
		}
	}
	return toolConfidence
}

func getFeatureValue(featureKey string, featureValue string, tool ToolResponse) FeatureValue {
	tools := []ToolConfidence{
		getToolConfidence(tool, featureKey),
	}
	value := FeatureValue{
		Value: featureValue,
		Score: 0.0,
		Tools: tools,
	}
	return value
}

func getFeature(featureKey string, featureValue string, tool ToolResponse) Feature {
	tools := []ToolConfidence{
		getToolConfidence(tool, featureKey),
	}
	values := []FeatureValue{
		{Value: featureValue, Tools: tools},
	}
	feature := Feature{
		Key:    featureKey,
		Values: values,
	}
	return feature
}

func calculateFeatureValueScore(features *map[string]Feature) {
	for featureKey, feauture := range *features {
		totalFeatureConfidence := 0.0
		totalValueConfidence := make(map[string]float64)
		for _, featureValue := range feauture.Values {
			totalValueConfidence[featureValue.Value] = 0.0
			for _, tool := range featureValue.Tools {
				totalFeatureConfidence += tool.Confidence
				totalValueConfidence[featureValue.Value] += tool.Confidence
			}
		}
		for valueIndex, featureValue := range feauture.Values {
			// if only one tool has extracted the feature
			if totalValueConfidence[featureValue.Value] == totalFeatureConfidence {
				// total confidence is equal to tool confidence
				(*features)[featureKey].Values[valueIndex].Score = totalValueConfidence[featureValue.Value]
			} else {
				// if multiple tools have extracted the feature, calculate the ratio
				(*features)[featureKey].Values[valueIndex].Score =
					totalValueConfidence[featureValue.Value] / totalFeatureConfidence
			}
		}
	}
}

func sortFeatureValues(features *map[string]Feature) {
	for featureKey := range *features {
		sort.Sort(ByScore((*features)[featureKey].Values))
	}
}

func correctToolConfidence(scoredFeatures *map[string]Feature) {
	for featureKey, feature := range *scoredFeatures {
		for featureValueIndex, featureValue := range feature.Values {
			for toolIndex, toolConfidence := range featureValue.Tools {
				log.Println(toolConfidence)
				(*scoredFeatures)[featureKey].Values[featureValueIndex].Tools[toolIndex] =
					getCorrectedToolConfidence(featureKey, toolConfidence, scoredFeatures)
			}
		}
	}
}
