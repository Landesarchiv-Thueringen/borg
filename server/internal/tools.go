package internal

import (
	"encoding/json"
	"fmt"
	"lath/borg/internal/config"
	"log"
	"net/http"
	"net/http/httputil"
	"regexp"
)

type ToolResult struct {
	ToolName string `json:"toolName"`
	// ToolType is the kind of tool. Possible values are "identification" and
	// "validation".
	ToolType string `json:"toolType"`
	// ToolVersion is the version number of the utilized tool as by the tool's
	// own versioning scheme.
	ToolVersion string `json:"toolVersion"`
	// ToolOutput is the tool's raw output string.
	ToolOutput string `json:"toolOutput"`
	// OutputFormat is the format expected for ToolOutput. Possible values are
	// "text", "json", and "csv".
	OutputFormat string `json:"outputFormat"`
	// Features is a list of features as extracted from the tool's output.
	Features map[string]interface{} `json:"features"`
	// Score is the from the tool supplied confidence of the result.
	Score *float64 `json:"score"`
	// Error is an error emitted from the tool in case of failure.
	Error string `json:"error"`
}

type toolResponse struct {
	ToolOutput   string                 `json:"toolOutput"`
	OutputFormat string                 `json:"outputFormat"`
	Features     map[string]interface{} `json:"features"`
	Score        *float64               `json:"score"`
	Error        string                 `json:"error"`
}

func RunIdentificationTools(filename string) []ToolResult {
	var responseChannels []chan ToolResult
	// for every identification tool
	for _, tool := range serverConfig.FormatIdentificationTools {
		rc := make(chan ToolResult)
		responseChannels = append(responseChannels, rc)
		// request tool results concurrent
		go func() {
			response := getToolResponse(tool.Endpoint, filename)
			rc <- ToolResult{
				ToolName:     tool.ToolName,
				ToolVersion:  tool.ToolVersion,
				ToolType:     "identification",
				ToolOutput:   response.ToolOutput,
				OutputFormat: response.OutputFormat,
				Features:     response.Features,
				Score:        response.Score,
				Error:        response.Error,
			}
		}()
	}
	// gather all tool responses
	var results []ToolResult
	for _, rc := range responseChannels {
		toolResponse := <-rc
		results = append(results, toolResponse)
	}
	return results
}

func RunValidationTools(filename string, identificationResults []ToolResult) []ToolResult {
	var responseChannels []chan ToolResult
	// for every validation tool
	for _, tool := range serverConfig.FormatValidationTools {
		// for every possible trigger of current validation tool
		for _, trigger := range tool.ToolTrigger {
			if checkToolTrigger(trigger, identificationResults) {
				rc := make(chan ToolResult)
				responseChannels = append(responseChannels, rc)
				// request tool results concurrent
				go func() {
					response := getToolResponse(tool.Endpoint, filename)
					rc <- ToolResult{
						ToolName:     tool.ToolName,
						ToolVersion:  tool.ToolVersion,
						ToolType:     "validation",
						ToolOutput:   response.ToolOutput,
						OutputFormat: response.OutputFormat,
						Features:     response.Features,
						Score:        response.Score,
						Error:        response.Error,
					}
				}()
				// don't check other triggers, tool response already requested
				break
			}
		}
	}
	// gather all tool responses
	results := identificationResults
	for _, rc := range responseChannels {
		toolResponse := <-rc
		results = append(results, toolResponse)
	}
	return results
}

// returns true if the trigger fires
func checkToolTrigger(trigger config.ToolTrigger, identificationResults []ToolResult) bool {
	// j, _ = json.MarshalIndent(identificationResults, "", "\t")
	// fmt.Printf("identificationResults%s \n", j)
	regex := regexp.MustCompile(trigger.RegEx)
	for _, toolResponse := range identificationResults {
		if toolResponse.Features != nil {
			features := toolResponse.Features
			featureValue, ok := features[trigger.Feature]
			if ok {
				v, ok := featureValue.(string)
				if ok && regex.MatchString(v) {
					return true
				}
			}
		}
	}
	return false
}

func getToolResponse(
	endpoint string,
	filename string,
) toolResponse {
	// create http get request
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Println(err)
		return toolResponse{Error: "error creating request: " + endpoint}
	}
	// add file path URL parameter
	query := req.URL.Query()
	query.Add("path", filename)
	req.URL.RawQuery = query.Encode()
	// send get request
	response, err := http.Get(req.URL.String())
	if err != nil {
		log.Println(err)
		return toolResponse{Error: "error requesting: " + req.URL.String()}
	}
	// process request response
	return processToolResponse(response)
}

func processToolResponse(response *http.Response) toolResponse {
	if response.StatusCode != http.StatusOK {
		toolResponse := toolResponse{
			Error: fmt.Sprintf("tool request error: %d", response.StatusCode),
		}
		bytes, err := httputil.DumpResponse(response, true)
		if err == nil {
			responseString := string(bytes)
			toolResponse.ToolOutput = responseString
			toolResponse.OutputFormat = "text"
		}
		return toolResponse
	}
	var parsedResponse toolResponse
	err := json.NewDecoder(response.Body).Decode(&parsedResponse)
	if err != nil {
		errorMessage := "error parsing tool response"
		log.Println(errorMessage)
		log.Println(err)
		return toolResponse{Error: errorMessage}
	}
	return parsedResponse
}
