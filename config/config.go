package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type WebHook struct {
	URL    string `yaml:"url"`
	Secret string `yaml:"secret"`
}

type Config struct {
	WebHook *WebHook `yaml:"webhook"`
}

// NewConfig loads config from config_env.yml file
func NewConfig(cf string) (*Config, error) {

	c := &Config{}

	confFile, err := ioutil.ReadFile(cf)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(confFile, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
