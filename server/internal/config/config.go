

package config

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
	DefaultValue float64               `yaml:"defaultValue"`
	Conditions   []ConfidenceCondition `yaml:"conditions"`
}

type FeatureConfig struct {
	Key        string           `yaml:"key"`
	Confidence ConfidenceConfig `yaml:"confidence"`
}

type FormatIdentificationTool struct {
	ToolName    string          `yaml:"toolName"`
	ToolVersion string          `yaml:"toolVersion"`
	Endpoint    string          `yaml:"endpoint"`
	Features    []FeatureConfig `yaml:"features"`
}

type ToolTrigger struct {
	Feature string `yaml:"feature"`
	RegEx   string `yaml:"regEx"`
}

type FormatValidationTool struct {
	ToolName    string          `yaml:"toolName"`
	ToolVersion string          `yaml:"toolVersion"`
	Endpoint    string          `yaml:"endpoint"`
	ToolTrigger []ToolTrigger   `yaml:"trigger"`
	Features    []FeatureConfig `yaml:"features"`
}

type ServerConfig struct {
	FormatIdentificationTools []FormatIdentificationTool `yaml:"formatIdentificationTools"`
	FormatValidationTools     []FormatValidationTool     `yaml:"formatValidationTools"`
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
