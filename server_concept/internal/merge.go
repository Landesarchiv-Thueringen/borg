package internal

import (
	"fmt"
	"slices"
	"sort"
)

type FeatureSet struct {
	Features        map[string]interface{} `json:"features"`
	SupportingTools []string               `json:"supportingTools"`
	Score           float64                `json:"score"`
}

type ByScore []FeatureSet

func (a ByScore) Len() int           { return len(a) }
func (a ByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByScore) Less(i, j int) bool { return a[i].Score < a[j].Score }

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

func normalizeSetScore(sets []FeatureSet) []FeatureSet {
	var normalizedSets []FeatureSet
	totalScore := 0.0
	for _, s := range sets {
		totalScore += s.Score
	}
	if totalScore == 0.0 {
		return normalizedSets
	}
	for _, s := range sets {
		normalizedScore := s.Score / totalScore
		if normalizedScore < s.Score {
			s.Score = normalizedScore
		}
		normalizedSets = append(normalizedSets, s)
	}
	return normalizedSets
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
	featureValues := make(map[string][]FeatureValue)
	for _, tr := range m.toolResults {
		tc := getToolConfig(tr.Id)
		for k, v := range tr.Features {
			featureValue := FeatureValue{
				Value: v,
			}
			featureConfig, ok := tc.FeatureSet.GetFeatureConfig(k)
			if ok {
				featureValue.MergeOrder = featureConfig.MergeOrder
			}
			featureValues[k] = append(
				featureValues[k],
				featureValue,
			)
		}
	}
	for key, featureValues := range featureValues {
		sort.Sort(sort.Reverse(ByOrder(featureValues)))
		features[key] = featureValues[0].Value
	}
	supportingTools := make([]string, 0)
	score := 0.0
	for index, tc := range m.toolConfigs {
		supportingTools = append(supportingTools, tc.Id)
		score += tc.FeatureSet.Weight.GetWeight(m.toolResults[index])
	}
	return FeatureSet{
		Features:        features,
		SupportingTools: supportingTools,
		Score:           score,
	}
}

func MergeFeatureSets(toolResults map[string]ToolResult) []FeatureSet {
	var mergedSets []FeatureSet
	for toolId, tr1 := range toolResults {
		// don't merge tool results without any extracted features
		if len(tr1.Features) == 0 {
			continue
		}
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
	revisedSets := normalizeSetScore(filterDuplicateSets(mergedSets))
	sort.Sort(sort.Reverse(ByScore(revisedSets)))
	return revisedSets
}

func getToolConfig(id string) ToolConfig {
	for _, config := range serverConfig.Tools {
		if config.Id == id {
			return config
		}
	}
	panic(fmt.Sprintf("faulty configuration: no tool configuration for id: %s", id))
}
