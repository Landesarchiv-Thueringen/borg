package internal

const formatUncertainThreshold = 0.75
const validThreshold = 0.75

var requiredFormatFeatures = []string{"mimeType", "puid"}

// Summary accumulates validation results on the highest level.
//
// All values are calculated with simple rules from the extracted and scored
// feature values. The aim is to put different values in perspective to allow
// easy reasoning on the results. I.e., flags concerning validity are only set
// if we are reasonably sure about the determined file format.
type Summary struct {
	// Valid means the file could be identified as valid by one or more suitable
	// validators.
	Valid bool `json:"valid"`
	// Invalid means the file could be identified as invalid by one or more
	// suitable validators.
	Invalid bool `json:"invalid"`
	// FormatUncertain means the file format could not be identified with
	// sufficient confidence.
	FormatUncertain bool `json:"formatUncertain"`
	// ValidityConflict means there have been conflicting validation results
	// from tools with sufficient confidence.
	ValidityConflict bool `json:"validityConflict"`
	// Error means that one or more tools aborted with an error.
	Error bool `json:"error"`
	// PUID is the extracted PUID with the highest score.
	PUID string `json:"puid"`
	// MimeType is the extracted mime type with the highest score.
	MimeType string `json:"mimeType"`
	// FormatVersion is the extracted format version with the highest score.
	FormatVersion string `json:"formatVersion"`
}

func GetSummary(features map[string][]FeatureValue, toolResults []ToolResult) Summary {
	return Summary{
		Valid:            isValid(features),
		Invalid:          isInvalid(features),
		FormatUncertain:  isFormatUncertain(features),
		ValidityConflict: hasValidityConflict(features),
		Error:            hasError(toolResults),
		PUID:             getStringValue(features, "puid"),
		MimeType:         getStringValue(features, "mimeType"),
		FormatVersion:    getStringValue(features, "formatVersion"),
	}
}

func isValid(features map[string][]FeatureValue) bool {
	values, ok := features["valid"]
	if ok {
		b, ok := values[0].Value.(bool)
		return ok &&
			b &&
			values[0].Score >= validThreshold &&
			!isFormatUncertain(features)
	}
	return false
}

func isInvalid(features map[string][]FeatureValue) bool {
	values, ok := features["valid"]
	if ok {
		b, ok := values[0].Value.(bool)
		return ok &&
			!b &&
			values[0].Score >= validThreshold &&
			!isFormatUncertain(features)
	}
	return false
}

func isFormatUncertain(features map[string][]FeatureValue) bool {
	for _, key := range requiredFormatFeatures {
		values, ok := features[key]
		if !ok || values[0].Score < formatUncertainThreshold {
			return true
		}
	}
	return false
}

func hasValidityConflict(features map[string][]FeatureValue) bool {
	values, ok := features["valid"]
	if !ok ||
		values[0].Score >= validThreshold ||
		isFormatUncertain(features) {
		return false
	}
	// We have some results for validity, but not with a sufficiently high
	// score...
	for _, confidence := range values[0].SupportingTools {
		if confidence >= validThreshold {
			// ...however, there is at least one tool, that would have produced
			// a sufficiently high score, that was challenged by another tool.
			return true
		}
	}
	return false
}

func hasError(toolResults []ToolResult) bool {
	for _, r := range toolResults {
		if r.Error != "" {
			return true
		}
	}
	return false
}

func getStringValue(features map[string][]FeatureValue, key string) string {
	values, ok := features[key]
	if !ok {
		return ""
	}
	stringValue, ok := values[0].Value.(string)
	if !ok {
		return ""
	}
	return stringValue
}
