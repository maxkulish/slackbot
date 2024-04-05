// Package config provides functionality to load configuration settings.
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	WebHook string `yaml:"webhook"`
}

// NewConfig loads config from config_env.yml file
func NewConfig(cf string) (*Config, error) {
	c := &Config{}

	confFile, err := os.ReadFile(cf)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(confFile, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
