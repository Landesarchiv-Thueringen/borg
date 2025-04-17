package internal

import (
	"log"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Tools []ToolConfig `yaml:"tools"`
}

type ToolConfig struct {
	Id         string           `yaml:"id"`
	Title      string           `yaml:"title"`
	Endpoint   string           `yaml:"endpoint"`
	Triggers   []Trigger        `yaml:"triggers"`
	FeatureSet FeatureSetConfig `yaml:"featureSet"`
}

func (t *ToolConfig) IsTriggered(toolResults map[string]ToolResult) bool {
	for _, t := range t.Triggers {
		if t.IsTriggered(toolResults) {
			return true
		}
	}
	return false
}

type Trigger struct {
	Conditions []FeatureCondition `yaml:"conditions"`
}

func (t *Trigger) IsTriggered(toolResults map[string]ToolResult) bool {
	for _, condition := range t.Conditions {
		isFulFilled := false
		for _, toolResult := range toolResults {
			f, ok := toolResult.Features[condition.Feature]
			if !ok {
				continue
			}
			if condition.IsFulfilled(f) {
				isFulFilled = true
				break
			}
		}
		if !isFulFilled {
			return false
		}
	}
	return true
}

type FeatureSetConfig struct {
	Features        []Feature        `yaml:"features"`
	Weight          Weight           `yaml:"weight"`
	MergeConditions []MergeCondition `yaml:"mergeConditions"`
}

func (s *FeatureSetConfig) AreMergeable(
	fs1 map[string]interface{},
	fs2 map[string]interface{},
) bool {
	for _, condition := range s.MergeConditions {
		if !condition.IsFulfilled(fs1, fs2) {
			return false
		}
	}
	return true
}

type Feature struct {
	Key string `yaml:"key"`
}

type Weight struct {
	Default            float64             `yaml:"default"`
	ConditionalWeights []ConditionalWeight `yaml:"conditional"`
	ProvidedByTool     bool                `yaml:"providedByTool"`
}

type ConditionalWeight struct {
	Value      float64            `yaml:"value"`
	Conditions []FeatureCondition `yaml:"conditions"`
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

func (c *MergeCondition) IsFulfilled(fs1 map[string]interface{}, fs2 map[string]interface{}) bool {
	fv1, ok1 := fs1[c.Feature]
	fv2, ok2 := fs2[c.Feature]
	// if not both feature sets include the feature of the merge condition
	if !ok1 || !ok2 {
		// merge is possible
		return true
	}
	// if value extraction regular expression is configured
	if c.ValueRegEx != nil {
		regEx := regexp.MustCompile(*c.ValueRegEx)
		// extract comparable values from feature strings
		s1, ok1 := fv1.(string)
		s2, ok2 := fv2.(string)
		if !ok1 || !ok2 {
			log.Fatal("configuration faulty: used value extraction string on non string value")
		}
		// merge is possible if extracted values are equal
		return regEx.FindString(s1) == regEx.FindString(s2)
	}
	// merge is possible if features are equal
	return fv1 == fv2
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
