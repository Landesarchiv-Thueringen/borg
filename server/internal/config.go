package internal

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Tools             []ToolConfig       `yaml:"tools"`
	FileIdentityRules []FileIdentityRule `yaml:"fileIdentity"`
}

type FileIdentityRule struct {
	Conditions []FeatureCondition `yaml:"conditions"`
}

type LocalizationResource struct {
	Endpoint string `yaml:"endpoint"`
}

type ToolConfig struct {
	Id         string           `yaml:"id"`
	Enabled    bool             `yaml:"enabled"`
	Title      string           `yaml:"title"`
	Endpoint   string           `yaml:"endpoint"`
	Triggers   []Trigger        `yaml:"triggers"`
	FeatureSet FeatureSetConfig `yaml:"featureSet"`
}

func (t *ToolConfig) IsTriggered(toolResults map[string]ToolResult) (bool, map[string]ToolFeatureValue) {
	for _, t := range t.Triggers {
		isTriggered, matches := t.IsTriggered(toolResults)
		if isTriggered {
			return true, matches
		}
	}
	return false, nil
}

type Trigger struct {
	Conditions []FeatureCondition `yaml:"conditions"`
}

func (t *Trigger) IsTriggered(toolResults map[string]ToolResult) (bool, map[string]ToolFeatureValue) {
	matches := make(map[string]ToolFeatureValue)
	for _, condition := range t.Conditions {
		isFulFilled := false
		for _, toolResult := range toolResults {
			v, ok := toolResult.Features[condition.Feature]
			if !ok {
				continue
			}
			if condition.IsFulfilled(v.Value) {
				matches[condition.Feature] = v
				isFulFilled = true
				break
			}
		}
		if !isFulFilled {
			return false, matches
		}
	}
	return true, matches
}

type FeatureSetConfig struct {
	Features []FeatureConfig `yaml:"features"`
	Weight   Weight          `yaml:"weight"`
}

func (c *FeatureSetConfig) AreMergeable(
	fs1 map[string]MergeFeatureValue,
	fs2 map[string]ToolFeatureValue,
) (isFulfilled bool, mergeModifier float64) {
	// The merge is always possible if the origin set is empty.
	if len(fs1) == 0 {
		isFulfilled = true
		return
	}
	for _, feature := range c.Features {
		if feature.MergeCondition != nil {
			ok, strongLink := feature.MergeCondition.IsFulfilled(feature.Key, fs1, fs2)
			if !ok {
				isFulfilled = false
				return
			}
			if strongLink {
				mergeModifier += 0.25
			}
		}
	}
	// The merge is possible, if at least on condition is truly fulfilled.
	// That means that both feature sets have values for the condition feature.
	if mergeModifier > 0.0 {
		isFulfilled = true
	}
	return
}

func (c *FeatureSetConfig) GetFeatureConfig(key string) (FeatureConfig, bool) {
	for _, featureConfig := range c.Features {
		if featureConfig.Key == key {
			return featureConfig, true
		}
	}
	return FeatureConfig{}, false
}

type FeatureConfig struct {
	Key               string          `yaml:"key"`
	MergeOrder        uint            `yaml:"mergeOrder"`
	ProvidedByTrigger bool            `yaml:"providedByTrigger"`
	MergeCondition    *MergeCondition `yaml:"mergeCondition"`
}

type MergeCondition struct {
	ExactMatch bool    `yaml:"exactMatch"`
	ValueRegEx *string `yaml:"valueRegEx"`
}

type Weight struct {
	Default            float64             `yaml:"default"`
	ConditionalWeights []ConditionalWeight `yaml:"conditional"`
	ProvidedByTool     bool                `yaml:"providedByTool"`
}

// priority of different providers:
//   - 1. conditional weight
//   - 2. tool provided weight
//   - 3. default weight
func (w *Weight) GetWeight(tr ToolResult) float64 {
	for _, cw := range w.ConditionalWeights {
		if cw.IsFulfilled(tr) {
			return cw.Value
		}
	}
	if w.ProvidedByTool {
		if tr.Score == nil {
			errorMessage := "configuration error: a tool provided weight was set "
			errorMessage += "for a tool that does not support this option"
			log.Fatal(errorMessage)
		}
		return *tr.Score
	}
	return w.Default
}

type ConditionalWeight struct {
	Value      float64            `yaml:"value"`
	Conditions []FeatureCondition `yaml:"conditions"`
}

func (w *ConditionalWeight) IsFulfilled(tr ToolResult) bool {
	for _, c := range w.Conditions {
		v, ok := tr.Features[c.Feature]
		if !ok {
			return false
		}
		if !c.IsFulfilled(v.Value) {
			return false
		}
	}
	return true
}

type FeatureCondition struct {
	Feature string      `yaml:"feature"`
	RegEx   *string     `yaml:"regEx"`
	Value   interface{} `yaml:"value"`
}

func (c *FeatureCondition) IsFulfilled(value interface{}) bool {
	if c.RegEx != nil {
		regEx := regexp.MustCompile(*c.RegEx)
		v, ok := value.(string)
		if !ok {
			log.Fatal("regular expression condition requires a string-typed feature value")
		}
		return regEx.MatchString(v)
	} else if c.Value != nil {
		return value == c.Value
	}
	return false
}

func (c *MergeCondition) IsFulfilled(featureKey string, fs1 map[string]MergeFeatureValue, fs2 map[string]ToolFeatureValue) (isFulfilled bool, strongLink bool) {
	// if the second feature sets doesn't contain any values
	// the first feature set can be empty if merging against an empty set
	if len(fs2) == 0 {
		// merge is not possible because it doesn't add any features but improves the score
		return
	}
	fv1, ok1 := fs1[featureKey]
	fv2, ok2 := fs2[featureKey]
	// if not both feature sets include the feature of the merge condition
	if !ok1 || !ok2 {
		// merge is possible but not a strong link
		isFulfilled = true
		return
	}
	// if value extraction regular expression is configured
	if c.ValueRegEx != nil {
		regEx := regexp.MustCompile(*c.ValueRegEx)
		// extract comparable values from feature strings
		s1, ok1 := fv1.Value.(string)
		s2, ok2 := fv2.Value.(string)
		if !ok1 || !ok2 {
			errorMessage := fmt.Sprintf(
				"configuration error: "+
					"used value extraction string on non string value for key {%s}",
				featureKey)
			log.Fatal(errorMessage)
		}
		m1 := regEx.FindStringSubmatch(s1)
		m2 := regEx.FindStringSubmatch(s2)
		if len(m1) != 2 || len(m2) != 2 {
			isFulfilled = false
			return
		}
		// merge is possible if extracted values are equal
		isFulfilled = m1[1] == m2[1]
		if isFulfilled {
			strongLink = true
		}
		return
	}
	// merge is possible if features are equal
	isFulfilled = fv1.Value == fv2.Value
	if isFulfilled {
		strongLink = true
	}
	return
}

var serverConfig ServerConfig

func ParseConfig() {
	bytes, err := os.ReadFile("config/server_config.yml")
	if err != nil {
		log.Fatal("server config not readable\n" + err.Error())
	}
	var config ServerConfig
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatal("server config couldn't be parsed\n" + err.Error())
	}
	serverConfig = config
}
