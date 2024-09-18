package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Loader defines the interface for loading configuration
type Loader interface {
	LoadConfig(configFile string) (*TerraformHybridConfig, error)
}

// TerraformConfigLoader is the concrete implementation of ConfigLoader
type TerraformConfigLoader struct{}

func NewConfigLoader() Loader {
	return &TerraformConfigLoader{}
}

func (tcl *TerraformConfigLoader) LoadConfig(configFile string) (*TerraformHybridConfig, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("error reading config file %s: %v", configFile, err)
	}

	var config TerraformHybridConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config file %s: %v", configFile, err)
	}

	return &config, nil
}
