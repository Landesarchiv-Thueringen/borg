package internal

import (
	"fmt"
	"slices"
)

type FeatureSet struct {
	Features        map[string]interface{} `json:"features"`
	SupportingTools []string               `json:"supportingTools"`
}

func (s1 *FeatureSet) IsEqual(s2 FeatureSet) bool {
	if len(s1.SupportingTools) != len(s2.SupportingTools) {
		return false
	}
	t1 := s1.SupportingTools
	t2 := s2.SupportingTools
	slices.Sort(t1)
	slices.Sort(t2)
	for index, toolId := range t1 {
		if toolId != t2[index] {
			return false
		}
	}
	return true
}

func filterDuplicateSets(sets []FeatureSet) []FeatureSet {
	var filteredSets []FeatureSet
	for _, s := range sets {
		setExistsAlready := false
		for _, fs := range filteredSets {
			if s.IsEqual(fs) {
				setExistsAlready = true
				break
			}
		}
		if !setExistsAlready {
			filteredSets = append(filteredSets, s)
		}
	}
	return filteredSets
}

type Merge struct {
	toolConfigs []ToolConfig
	toolResults []ToolResult
}

func (m *Merge) MergeIfPossible(tc2 ToolConfig, tr2 ToolResult) {
	if m.IsMergeable(tc2, tr2) {
		m.toolConfigs = append(m.toolConfigs, tc2)
		m.toolResults = append(m.toolResults, tr2)
	}
}

func (m *Merge) IsMergeable(tc2 ToolConfig, tr2 ToolResult) bool {
	if tr2.Error != nil {
		return false
	}
	mergedResults := m.GetMergedToolResults()
	// check if all conditions of the new set are met
	if !tc2.FeatureSet.AreMergeable(mergedResults.Features, tr2.Features) {
		return false
	}
	// check if all conditions of already merged sets are met
	for _, tc := range m.toolConfigs {
		if !tc.FeatureSet.AreMergeable(mergedResults.Features, tr2.Features) {
			return false
		}
	}
	return true
}

func (m *Merge) GetMergedToolResults() FeatureSet {
	features := make(map[string]interface{})
	for _, tr := range m.toolResults {
		for k, v := range tr.Features {
			features[k] = v
		}
	}
	supportingTools := make([]string, 0)
	for _, tc := range m.toolConfigs {
		supportingTools = append(supportingTools, tc.Id)
	}
	return FeatureSet{
		Features:        features,
		SupportingTools: supportingTools,
	}
}

func MergeFeatureSets(toolResults map[string]ToolResult) []FeatureSet {
	var mergedSets []FeatureSet
	for toolId, tr1 := range toolResults {
		// don't merge results with errors
		if tr1.Error != nil {
			continue
		}
		var m Merge
		tc1 := getToolConfig(toolId)
		m.MergeIfPossible(tc1, tr1)
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
			m.MergeIfPossible(tc2, tr2)
		}
		mergedSets = append(mergedSets, m.GetMergedToolResults())
	}
	return filterDuplicateSets(mergedSets)
}

func getToolConfig(id string) ToolConfig {
	for _, config := range serverConfig.Tools {
		if config.Id == id {
			return config
		}
	}
	panic(fmt.Sprintf("faulty configuration: no tool configuration for id: %s", id))
}
