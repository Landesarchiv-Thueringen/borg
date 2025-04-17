package internal

func MergeFeatureSets(toolResults map[string]ToolResult) {
	for toolId, tr1 := range toolResults {
		for _, tc2 := range serverConfig.Tools {
			// don't merge feature set with itself
			if toolId == tc2.Id {
				continue
			}
			// check if a result exists for tool configuration
			tr2, ok := toolResults[tc2.Id]
			if !ok {
				continue
			}
			tc2.FeatureSet.AreMergeable(tr1, tr2)
		}
	}
}
