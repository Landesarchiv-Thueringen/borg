package internal

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfidenceCondition struct {
	GlobalFeature string  `yaml:"globalFeature"`
	RegEx         string  `yaml:"regEx"`
	Value         float64 `yaml:"value"`
}

type ConfidenceConfig struct {
	ProvidedByTool bool                  `yaml:"providedByTool"`
	DefaultValue   float64               `yaml:"defaultValue"`
	Conditions     []ConfidenceCondition `yaml:"conditions"`
}

type FeatureConfig struct {
	Key        string           `yaml:"key"`
	Confidence ConfidenceConfig `yaml:"confidence"`
}

type ToolTrigger struct {
	Feature string `yaml:"feature"`
	RegEx   string `yaml:"regEx"`
}

type ToolConfig struct {
	ToolName    string          `yaml:"toolName"`
	ToolVersion string          `yaml:"toolVersion"`
	Endpoint    string          `yaml:"endpoint"`
	ToolTrigger []ToolTrigger   `yaml:"trigger"`
	Features    []FeatureConfig `yaml:"features"`
}

type ServerConfig struct {
	Tools []ToolConfig `yaml:"tools"`
}

func ParseConfig() ServerConfig {
	bytes, err := os.ReadFile("config/server_config.yml")
	if err != nil {
		log.Fatal("server config not readable\n" + err.Error())
	}
	var config ServerConfig
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatal("server config couldn't be parsed\n" + err.Error())
	}
	return config
}
