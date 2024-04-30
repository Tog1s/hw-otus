package config

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Logger Logger
	// Storage Storage
	// DB DB
	// Server Server
}

type Logger struct {
	Level  string
	Output string
}

func New(path string) (*Config, error) {
	configData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error while reading config file %s: %w", path, err)
	}
	var config *Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return nil, fmt.Errorf("error while parse yaml file %s: %w", path, err)
	}
	return config, nil
}
