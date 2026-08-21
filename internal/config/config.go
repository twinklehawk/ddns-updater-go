package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

const (
	yamlFile = "./config.yaml"
)

type Config struct {
	Namecheap NamecheapConfig
	Ddns      DdnsConfig
}

type NamecheapConfig struct {
	Password string
}

type DdnsConfig struct {
	Domain     string
	Provider   string
	Subdomains []string
}

func ReadConfig() (*Config, error) {
	data, err := os.ReadFile(yamlFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", yamlFile, err)
	}
	c := Config{}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", yamlFile, err)
	}
	return &c, nil
}
