package internal

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Tools []Tool `yaml:"tools"`
}

type Tool struct {
	Id         string        `yaml:"id"`
	Title      string        `yaml:"title"`
	Endpoint   string        `yaml:"endpoint"`
	Trigger    []TriggerItem `yaml:"trigger"`
	FeatureSet FeatureSet    `yaml:"featureSet"`
}

type TriggerItem struct {
	Conditions []Condition `yaml:"conditions"`
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
