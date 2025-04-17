package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

type ToolResult struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	// ToolVersion is the version number of the utilized tool as by the tool's
	// own versioning scheme.
	// ToolVersion string `json:"toolVersion"`
	// ToolOutput is the tool's raw output string.
	ToolOutput string `json:"toolOutput"`
	// OutputFormat is the format expected for ToolOutput. Possible values are
	// "text", "json", and "csv".
	OutputFormat string `json:"outputFormat"`
	// Features is a list of features as extracted from the tool's output.
	Features map[string]interface{} `json:"features"`
	// Score is the from the tool supplied confidence of the result.
	// Score *float64 `json:"score"`
	// Error is an error emitted from the tool in case of failure.
	Error string `json:"error"`
}

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
			rc <- ToolResult{
				Id:    tool.Id,
				Title: tool.Title,
				// ToolVersion:  tool.ToolVersion,
				ToolOutput:   response.ToolOutput,
				OutputFormat: response.OutputFormat,
				Features:     response.Features,
				// Score:        response.Score,
				Error: response.Error,
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
	for _, tool := range serverConfig.Tools {
		if len(tool.Triggers) == 0 || !tool.IsTriggered(identificationResults) {
			continue
		}
		rc := make(chan ToolResult)
		responseChannels = append(responseChannels, rc)
		// request tool results concurrent
		go func() {
			response := getToolResult(tool.Endpoint, filename)
			rc <- ToolResult{
				Id:    tool.Id,
				Title: tool.Title,
				// ToolVersion:  tool.ToolVersion,
				ToolOutput:   response.ToolOutput,
				OutputFormat: response.OutputFormat,
				Features:     response.Features,
				// Score:        response.Score,
				Error: response.Error,
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
		return ToolResult{Error: fmt.Sprintf("error creating request: %s", endpoint)}
	}
	// add file path URL parameter
	query := req.URL.Query()
	query.Add("path", filename)
	req.URL.RawQuery = query.Encode()
	// send get request
	response, err := http.Get(req.URL.String())
	if err != nil {
		log.Println(err)
		return ToolResult{Error: fmt.Sprintf("error requesting: %s", req.URL.String())}
	}
	// process request response
	return processToolResponse(response)
}

func processToolResponse(response *http.Response) ToolResult {
	if response.StatusCode != http.StatusOK {
		toolResponse := ToolResult{
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
	var result ToolResult
	err := json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		errorMessage := "error parsing tool response"
		log.Println(errorMessage)
		log.Println(err)
		return ToolResult{Error: errorMessage}
	}
	return result
}
