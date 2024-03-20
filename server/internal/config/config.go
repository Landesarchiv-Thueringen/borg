/* BorgFormat - File format identification and validation
 * Copyright (C) 2024 Landesarchiv Thüringen
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

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
