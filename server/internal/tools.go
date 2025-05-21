package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"net/http"
	"net/http/httputil"
	"slices"
	"sort"
)

type ToolResult struct {
	Id    string `json:"id"`
	Title string `json:"title"`
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
	Error *string `json:"error"`
}

type ByTitle []ToolResult

func (a ByTitle) Len() int           { return len(a) }
func (a ByTitle) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTitle) Less(i, j int) bool { return a[i].Title < a[j].Title }

func RunIdentificationTools(filename string) map[string]ToolResult {
	var responseChannels []chan ToolResult
	// for every identification tool
	for _, tool := range serverConfig.Tools {
		if len(tool.Triggers) > 0 {
			continue
		}
		rc := make(chan ToolResult)
		responseChannels = append(responseChannels, rc)
		// request tool results concurrent
		go func() {
			response := getToolResult(tool.Endpoint, filename)
			features := make(map[string]interface{})
			if len(response.Features) > 0 {
				features = response.Features
			}
			rc <- ToolResult{
				Id:           tool.Id,
				Title:        tool.Title,
				ToolVersion:  response.ToolVersion,
				ToolOutput:   response.ToolOutput,
				OutputFormat: response.OutputFormat,
				Features:     features,
				Score:        response.Score,
				Error:        response.Error,
			}
		}()
	}
	// gather all tool responses
	results := make(map[string]ToolResult)
	for _, rc := range responseChannels {
		toolResponse := <-rc
		results[toolResponse.Id] = toolResponse
	}
	return results
}

func RunTriggeredTools(
	filename string,
	identificationResults map[string]ToolResult,
) map[string]ToolResult {
	var responseChannels []chan ToolResult
	// for every identification tool
	for _, toolConfig := range serverConfig.Tools {
		isTriggered, matches := toolConfig.IsTriggered(identificationResults)
		if len(toolConfig.Triggers) == 0 || !isTriggered {
			continue
		}
		rc := make(chan ToolResult)
		responseChannels = append(responseChannels, rc)
		// request tool results concurrent
		go func() {
			response := getToolResult(toolConfig.Endpoint, filename)
			features := make(map[string]interface{})
			if len(response.Features) > 0 {
				features = response.Features
			}
			// get feature values from tool trigger
			for _, featureConfig := range toolConfig.FeatureSet.Features {
				if featureConfig.ProvidedByTrigger {
					v, ok := matches[featureConfig.Key]
					if ok {
						features[featureConfig.Key] = v
					}
				}
			}
			rc <- ToolResult{
				Id:           toolConfig.Id,
				Title:        toolConfig.Title,
				ToolVersion:  response.ToolVersion,
				ToolOutput:   response.ToolOutput,
				OutputFormat: response.OutputFormat,
				Features:     features,
				Score:        response.Score,
				Error:        response.Error,
			}
		}()
	}
	// gather all tool responses
	results := make(map[string]ToolResult)
	for _, rc := range responseChannels {
		toolResponse := <-rc
		results[toolResponse.Id] = toolResponse
	}
	return results
}

func getToolResult(
	endpoint string,
	filename string,
) ToolResult {
	// create http get request
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Println(err)
		errorMessage := fmt.Sprintf("error creating request: %s", endpoint)
		return ToolResult{Error: &errorMessage}
	}
	// add file path URL parameter
	query := req.URL.Query()
	query.Add("path", filename)
	req.URL.RawQuery = query.Encode()
	// send get request
	response, err := http.Get(req.URL.String())
	if err != nil {
		log.Println(err)
		errorMessage := fmt.Sprintf("error requesting: %s", req.URL.String())
		return ToolResult{Error: &errorMessage}
	}
	// process request response
	return processToolResponse(response)
}

func processToolResponse(response *http.Response) ToolResult {
	if response.StatusCode != http.StatusOK {
		errorMessage := fmt.Sprintf("tool request error: %d", response.StatusCode)
		toolResponse := ToolResult{
			Error: &errorMessage,
		}
		bytes, err := httputil.DumpResponse(response, true)
		if err == nil {
			responseString := string(bytes)
			toolResponse.ToolOutput = responseString
			toolResponse.OutputFormat = "text"
		}
		return toolResponse
	}
	var result ToolResult
	err := json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		errorMessage := "error parsing tool response"
		log.Println(errorMessage)
		log.Println(err)
		return ToolResult{Error: &errorMessage}
	}
	return result
}

func CombineToolResults(
	identResults map[string]ToolResult,
	triggeredResults map[string]ToolResult,
) map[string]ToolResult {
	toolResults := make(map[string]ToolResult)
	for k, v := range identResults {
		toolResults[k] = v
	}
	for k, v := range triggeredResults {
		toolResults[k] = v
	}
	return toolResults
}

func GetSortedToolResults(
	identResults map[string]ToolResult,
	triggeredResults map[string]ToolResult,
) []ToolResult {
	r1 := slices.Collect(maps.Values(identResults))
	r2 := slices.Collect(maps.Values(triggeredResults))
	sort.Sort(ByTitle(r1))
	sort.Sort(ByTitle(r2))
	return append(r1, r2...)
}
