// Package config provides YAML configuration loading for the dwm CLI.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the top-level configuration structure parsed from a YAML file.
type Config struct {
	Project      string       `yaml:"project"`
	Environments []string     `yaml:"environments"`
	Connections  []ConnConfig `yaml:"connections"`
	Engines      []string     `yaml:"engines"`
}

// ConnConfig describes a single database connection.
type ConnConfig struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	Folder   string `yaml:"folder"`
}

// Load reads the YAML file at path and returns a parsed Config.
// It returns an error if the file cannot be read or if the YAML is malformed.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config.Load: read %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config.Load: parse %s: %w", path, err)
	}

	return &cfg, nil
}
