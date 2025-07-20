package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the agent-hooks configuration
type Config struct {
	Disable bool `yaml:"disable"`
}

// LoadConfig loads the .agenthooks config file from the current directory or any parent directory
func LoadConfig() (*Config, error) {
	config := &Config{}

	configPath, err := findConfigFile()
	if err != nil {
		return config, err
	}

	if configPath == "" {
		// No config file found, return default config
		return config, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
	}

	return config, nil
}

// findConfigFile searches for .agenthooks config file in current directory and parent directories
func findConfigFile() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	for {
		configPath := filepath.Join(dir, ".agenthooks")
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root directory
			break
		}
		dir = parent
	}

	return "", nil
}
