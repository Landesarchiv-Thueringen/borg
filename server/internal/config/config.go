package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type FormatIdentificationTool struct {
	ToolName    string `yaml:"toolName"`
	ToolVersion string `yaml:"toolVersion"`
	Endpoint    string `yaml:"endpoint"`
}

type ServerConfig struct {
	FormatIdentificationTools []FormatIdentificationTool `yaml:"formatIdentificationTools"`
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
