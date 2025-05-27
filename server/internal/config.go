package internal

import (
	"log"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Tools             []ToolConfig       `yaml:"tools"`
	FileIdentityRules []FileIdentityRule `yaml:"fileIdentity"`
	Localization      LocalizationConfig `yaml:"localization"`
}

type FileIdentityRule struct {
	Conditions []FeatureCondition `yaml:"conditions"`
}

type LocalizationConfig struct {
	Resources []LocalizationResource `yaml:"resources"`
	Dict      map[string]string      `yaml:"dict"`
}

type LocalizationResource struct {
	Endpoint string `yaml:"endpoint"`
}

type ToolConfig struct {
	Id         string           `yaml:"id"`
	Title      string           `yaml:"title"`
	Endpoint   string           `yaml:"endpoint"`
	Triggers   []Trigger        `yaml:"triggers"`
	FeatureSet FeatureSetConfig `yaml:"featureSet"`
}

func (t *ToolConfig) IsTriggered(toolResults map[string]ToolResult) (bool, map[string]interface{}) {
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

func (t *Trigger) IsTriggered(toolResults map[string]ToolResult) (bool, map[string]interface{}) {
	matches := make(map[string]interface{})
	for _, condition := range t.Conditions {
		isFulFilled := false
		for _, toolResult := range toolResults {
			v, ok := toolResult.Features[condition.Feature]
			if !ok {
				continue
			}
			if condition.IsFulfilled(v) {
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
	Features        []Feature        `yaml:"features"`
	Weight          Weight           `yaml:"weight"`
	MergeConditions []MergeCondition `yaml:"mergeConditions"`
}

func (c *FeatureSetConfig) AreMergeable(
	fs1 map[string]FeatureValue,
	fs2 map[string]interface{},
) (isFulfilled bool, mergeModifier float64) {
	// The merge is always possible if the origin set is empty.
	if len(fs1) == 0 {
		isFulfilled = true
		return
	}
	for _, condition := range c.MergeConditions {
		ok, strongLink := condition.IsFulfilled(fs1, fs2)
		if !ok {
			isFulfilled = false
			return
		}
		if strongLink {
			mergeModifier += 0.25
		}
	}
	// The merge is possible, if at least on condition is truly fulfilled.
	// That means that both feature sets have values for the condition feature.
	if mergeModifier > 0.0 {
		isFulfilled = true
	}
	return
}

func (c *FeatureSetConfig) GetFeatureConfig(key string) (Feature, bool) {
	for _, featureConfig := range c.Features {
		if featureConfig.Key == key {
			return featureConfig, true
		}
	}
	return Feature{}, false
}

type Feature struct {
	Key               string `yaml:"key"`
	MergeOrder        uint   `yaml:"mergeOrder"`
	ProvidedByTrigger bool   `yaml:"providedByTrigger"`
}

type FeatureValue struct {
	Value           interface{} `json:"value"`
	MergeOrder      uint        `json:"-"`
	SupportingTools []string    `json:"supportingTools"`
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
		if !c.IsFulfilled(v) {
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

type MergeCondition struct {
	Feature    string  `yaml:"feature"`
	ValueRegEx *string `yaml:"valueRegEx"`
}

func (c *MergeCondition) IsFulfilled(fs1 map[string]FeatureValue, fs2 map[string]interface{}) (isFulfilled bool, strongLink bool) {
	// if the second feature sets doesn't contain any values
	// the first feature set can be empty if merging against an empty set
	if len(fs2) == 0 {
		// merge is not possible because it doesn't add any features but improves the score
		return
	}
	fv1, ok1 := fs1[c.Feature]
	fv2, ok2 := fs2[c.Feature]
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
		s2, ok2 := fv2.(string)
		if !ok1 || !ok2 {
			log.Fatal("configuration faulty: used value extraction string on non string value")
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
	isFulfilled = fv1.Value == fv2
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
