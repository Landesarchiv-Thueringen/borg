package internal

import (
	"lath/borg/internal/config"
	"regexp"
	"sort"
)

type FeatureValue struct {
	// Value is a string representing the feature value, e.g.,
	// "application/pdf".
	Value interface{} `json:"value"`
	// Score is a number between 0 and 1 that ranks the value against other
	// values. The sum of scores of all values is at most one, but can be lower,
	// depending on supporting tools.
	Score float64 `json:"score"`
	// SupportingTools is a list of all tools that support this value and their
	// associated weights / confidences.
	SupportingTools map[string]float64 `json:"supportingTools"`
}

// implement sorting interface for feature values
type byScore []FeatureValue

func (featureValues byScore) Len() int {
	return len(featureValues)
}

func (featureValues byScore) Less(i, j int) bool {
	// sort reversed, biggest score first
	return featureValues[i].Score > featureValues[j].Score
}

func (featureValues byScore) Swap(i, j int) {
	featureValues[i], featureValues[j] = featureValues[j], featureValues[i]
}

func AccumulateFeatures(results []ToolResult) map[string][]FeatureValue {
	features := make(map[string][]FeatureValue)
	// for every tool response
	for _, r := range results {
		// for every extracted feature
		for featureKey, featureValue := range r.Features {
			// add current tool to feature
			index := -1
			for i, v := range features[featureKey] {
				// extracted value exists already, that means another tool extracted the same value
				if v.Value == featureValue {
					index = i
					break
				}
			}
			// value for key doesn't exist already --> add it to value list
			if index < 0 {
				features[featureKey] = append(features[featureKey], FeatureValue{
					Value:           featureValue,
					SupportingTools: make(map[string]float64),
				})
				index = len(features[featureKey]) - 1
			}
			// add tool to tools that extracted current value for feature
			features[featureKey][index].SupportingTools[r.ToolName] = toolConfidences[r.ToolName][featureKey]
		}
	}
	calculateFeatureValueScore(features)
	sortFeatureValues(features)
	// only after first score calculation, can global feature conditions be applied
	correctToolConfidence(features)
	calculateFeatureValueScore(features)
	sortFeatureValues(features)
	return features
}

func getFeatureConfig(toolName string) []config.FeatureConfig {
	for _, tool := range serverConfig.FormatIdentificationTools {
		if tool.ToolName == toolName {
			return tool.Features
		}
	}
	for _, tool := range serverConfig.FormatValidationTools {
		if tool.ToolName == toolName {
			return tool.Features
		}
	}
	panic("no such tool: " + toolName)
}

func getCorrectedToolConfidence(
	toolName string,
	featureKey string,
	toolConfidence float64,
	scoredFeatures map[string][]FeatureValue,
) float64 {
	featureConfig := getFeatureConfig(toolName)
	for _, featureConfig := range featureConfig {
		// if feature configuration doesn't belong to currently corrected feature
		if featureKey != featureConfig.Key {
			continue
		}
		if len(featureConfig.Confidence.Conditions) > 0 {
			for _, condition := range featureConfig.Confidence.Conditions {
				scoredFeature, ok := scoredFeatures[condition.GlobalFeature]
				if ok {
					regex := regexp.MustCompile(condition.RegEx)
					stringValue, ok := scoredFeature[0].Value.(string)
					// the first value has the highest score --> voted truth
					if ok && regex.MatchString(stringValue) {
						return condition.Value
					}
				}
			}
		}
	}
	return toolConfidence
}

func calculateFeatureValueScore(featureValues map[string][]FeatureValue) {
	for key, values := range featureValues {
		totalFeatureConfidence := 0.0
		totalValueConfidence := make(map[interface{}]float64)
		for _, featureValue := range values {
			totalValueConfidence[featureValue.Value] = 0.0
			for _, toolConfidence := range featureValue.SupportingTools {
				totalFeatureConfidence += toolConfidence
				totalValueConfidence[featureValue.Value] += toolConfidence
			}
		}
		for valueIndex, featureValue := range values {
			if totalFeatureConfidence == 0 {
				featureValues[key][valueIndex].Score = 0
			} else {
				valueConfidence := totalValueConfidence[featureValue.Value]
				normalizedValueConfidence := valueConfidence / totalFeatureConfidence
				featureValues[key][valueIndex].Score = min(valueConfidence, normalizedValueConfidence)
			}
		}
	}
}

func sortFeatureValues(featuresValues map[string][]FeatureValue) {
	for key := range featuresValues {
		sort.Sort(byScore(featuresValues[key]))
	}
}

func correctToolConfidence(scoredFeatureValues map[string][]FeatureValue) {
	for key, values := range scoredFeatureValues {
		for featureValueIndex, featureValue := range values {
			for toolName, toolConfidence := range featureValue.SupportingTools {
				scoredFeatureValues[key][featureValueIndex].SupportingTools[toolName] =
					getCorrectedToolConfidence(toolName, key, toolConfidence, scoredFeatureValues)
			}
		}
	}
}
