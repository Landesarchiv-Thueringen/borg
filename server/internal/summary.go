package internal

import "log"

const UNCERTAIN_THRESHOLD = 0.75

// Summary accumulates validation results on the highest level.
//
// All values are calculated with simple rules from the extracted and scored
// feature sets. The aim is to put different values in perspective to allow
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

func GetSummary(sets []FeatureSet, toolResults []ToolResult) Summary {
	var summary Summary
	if len(sets) == 0 {
		summary.FormatUncertain = true
		return summary
	}
	if sets[0].Score < UNCERTAIN_THRESHOLD {
		summary.FormatUncertain = true
	}
	validFeature, ok := sets[0].Features["format:valid"]
	if ok {
		valid, ok := validFeature.Value.(bool)
		if !ok {
			log.Fatal("valid feature has non boolean value")
		} else if valid {
			summary.Valid = true
		} else {
			summary.Invalid = true
		}
	}
	puidFeature, ok := sets[0].Features["format:puid"]
	if ok {
		puid, ok := puidFeature.Value.(string)
		if !ok {
			log.Fatal("PUID feature has non string value")
		} else {
			summary.PUID = puid
		}
	}
	mimeTypeFeature, ok := sets[0].Features["format:mimeType"]
	if ok {
		mimeType, ok := mimeTypeFeature.Value.(string)
		if !ok {
			log.Fatal("MIME type feature has non string value")
		} else {
			summary.MimeType = mimeType
		}
	}
	formatVersionFeature, ok := sets[0].Features["format:version"]
	if ok {
		formatVersion, ok := formatVersionFeature.Value.(string)
		if !ok {
			log.Fatal("format version feature has non string value")
		} else {
			summary.FormatVersion = formatVersion
		}
	}
	for _, result := range toolResults {
		if result.Error != nil {
			summary.Error = true
			break
		}
	}
	return summary
}
