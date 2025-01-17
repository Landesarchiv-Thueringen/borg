package internal

var serverConfig ServerConfig
var toolConfidences map[string]map[string]float64

func init() {
	serverConfig = ParseConfig()
	populateToolConfidences()
}

func populateToolConfidences() {
	toolConfidences = make(map[string]map[string]float64)
	for _, tool := range serverConfig.Tools {
		toolConfidences[tool.ToolName] = make(map[string]float64)
		for _, feature := range tool.Features {
			toolConfidences[tool.ToolName][feature.Key] = feature.Confidence.DefaultValue
		}
	}
}
