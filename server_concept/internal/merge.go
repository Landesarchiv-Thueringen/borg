package internal

import (
	"fmt"
	"slices"
	"sort"
)

type FeatureSet struct {
	Features        map[string]FeatureValue `json:"features"`
	SupportingTools []string                `json:"supportingTools"`
	Score           float64                 `json:"score"`
}

func (s *FeatureSet) FulFillesAny(fileIdentityRules []FileIdentityRule) bool {
	for _, rule := range fileIdentityRules {
		if s.FulFilles(rule) {
			return true
		}
	}
	return false
}

func (s *FeatureSet) FulFilles(fileIdentityRule FileIdentityRule) bool {
	for _, condition := range fileIdentityRule.Conditions {
		v, ok := s.Features[condition.Feature]
		if !ok || !condition.IsFulfilled(v.Value) {
			return false
		}
	}
	return true
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

func applyFileIdentityRules(sets []FeatureSet) []FeatureSet {
	for i, s := range sets {
		if s.FulFillesAny(serverConfig.FileIdentityRules) {
			return setFileIdentity(sets, i)
		}
	}
	return sets
}

func setFileIdentity(sets []FeatureSet, setIndex int) []FeatureSet {
	var adjustedSets []FeatureSet
	for i, s := range sets {
		if i == setIndex {
			s.Score = 1.0
		} else {
			s.Score = 0.0
		}
		adjustedSets = append(adjustedSets, s)
	}
	return adjustedSets
}

type Merge struct {
	toolConfigs      []ToolConfig
	toolResults      []ToolResult
	AccumulatedScore float64
}

func (m *Merge) MergeIfPossible(tc2 ToolConfig, tr2 ToolResult) {
	isMergeable, mergeModifier := m.IsMergeable(tc2, tr2)
	if isMergeable {
		if len(m.toolConfigs) == 0 {
			m.AccumulatedScore = tc2.FeatureSet.Weight.GetWeight(tr2)
		} else {
			m.AccumulatedScore += mergeModifier * tc2.FeatureSet.Weight.GetWeight(tr2)
		}
		m.toolConfigs = append(m.toolConfigs, tc2)
		m.toolResults = append(m.toolResults, tr2)
	}
}

func (m *Merge) IsMergeable(tc2 ToolConfig, tr2 ToolResult) (isMergeable bool, mergeModifier float64) {
	if tr2.Error != nil {
		return
	}
	mergedResults := m.GetMergedToolResults()
	// check if all conditions of the new set are met
	isMergeable, mergeModifier = tc2.FeatureSet.AreMergeable(mergedResults.Features, tr2.Features)
	if !isMergeable {
		return
	}
	// check if all conditions of already merged sets are met
	for _, tc := range m.toolConfigs {
		subsetMergeable, _ := tc.FeatureSet.AreMergeable(mergedResults.Features, tr2.Features)
		if !subsetMergeable {
			return
		}
	}
	return
}

func (m *Merge) GetMergedToolResults() FeatureSet {
	features := make(map[string]FeatureValue)
	featureValues := make(map[string][]FeatureValue)
	for _, tr := range m.toolResults {
		tc := getToolConfig(tr.Id)
		for k, v := range tr.Features {
			featureValue := FeatureValue{
				Value:           v,
				SupportingTools: []string{tc.Id},
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
	for key, values := range featureValues {
		for i, v := range values {
			if i == 0 {
				features[key] = v
			} else {
				if features[key].Value == v.Value {
					tools := append(features[key].SupportingTools, v.SupportingTools...)
					mergeOrder := features[key].MergeOrder
					if v.MergeOrder > mergeOrder {
						mergeOrder = v.MergeOrder
					}
					features[key] = FeatureValue{
						Value:           v.Value,
						SupportingTools: tools,
						MergeOrder:      mergeOrder,
					}
				} else if v.MergeOrder > features[key].MergeOrder {
					features[key] = v
				}
			}
		}
	}
	supportingTools := make([]string, 0)
	for _, tc := range m.toolConfigs {
		supportingTools = append(supportingTools, tc.Id)
	}
	return FeatureSet{
		Features:        features,
		SupportingTools: supportingTools,
		Score:           m.AccumulatedScore,
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
	revisedSets := filterDuplicateSets(mergedSets)
	revisedSets = normalizeSetScore(revisedSets)
	revisedSets = applyFileIdentityRules(revisedSets)
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
