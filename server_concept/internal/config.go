package internal

import (
	"log"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Tools []Tool `yaml:"tools"`
}

type Tool struct {
	Id         string     `yaml:"id"`
	Title      string     `yaml:"title"`
	Endpoint   string     `yaml:"endpoint"`
	Triggers   []Trigger  `yaml:"triggers"`
	FeatureSet FeatureSet `yaml:"featureSet"`
}

func (t *Tool) IsTriggered(toolResults []ToolResult) bool {
	for _, t := range t.Triggers {
		if t.IsTriggered(toolResults) {
			return true
		}
	}
	return false
}

type Trigger struct {
	Conditions []Condition `yaml:"conditions"`
}

func (t *Trigger) IsTriggered(toolResults []ToolResult) bool {
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

type FeatureSet struct {
	Features []Feature `yaml:"features"`
	Weight   Weight    `yaml:"weight"`
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
	Value      float64     `yaml:"value"`
	Conditions []Condition `yaml:"conditions"`
}

type Condition struct {
	Feature string      `yaml:"feature"`
	RegEx   *string     `yaml:"regEx"`
	Value   interface{} `yaml:"value"`
}

func (c *Condition) IsFulfilled(value interface{}) bool {
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
