// Package config handles reading configuration from a yaml file.
package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

const (
	defaultConfigFile = "./config.yaml"
)

// Config is the configuration for the ddns-updater-go command.
type Config struct {
	Ddns      DdnsConfig
	Namecheap NamecheapConfig
}

// DdnsConfig is the configuration for the domains that the command should update.
type DdnsConfig struct {
	// Domain is the top level domain, for example "example.com".
	Domain string
	// Provider is the DDNS provider, or what service provides the nameservers the domain uses.
	// This value should be set to a supported DDNS provider and is used to determine how to update the domain.
	Provider string
	// Subdomains is the list of subdomains that should have DDNS records updated with the current IP address.
	Subdomains []string
}

// NamecheapConfig is the configuration for calling namecheap.
type NamecheapConfig struct {
	// Password is the namecheap password to use when updating the DDNS record for a domain.
	Password string
}

// ReadConfig reads configuration from a YAML file and parses the contents into a Config instance.
//
// configFile is the path to the YAML file to load configuration from.
// If configFile is empty, the default of "./config.yaml" will be used.
func ReadConfig(configFile string) (*Config, error) {
	if configFile == "" {
		configFile = defaultConfigFile
	}
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configFile, err)
	}
	c := Config{}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configFile, err)
	}
	return &c, nil
}
