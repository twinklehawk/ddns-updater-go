package config

import (
	"fmt"
	"log/slog"
	"os"

	"go.yaml.in/yaml/v4"
)

const (
	defaultConfigFile = "./config.yaml"
)

// EnvValueHandler updates a value in the passed in [Config] instance using the
// passed in environment variable value.
type EnvValueHandler func(value string, config *Config)

// An EnvConfigEntry contains an environment variable key and a handler to
// update the config with the env value.
type EnvConfigEntry struct {
	// Env is the name of the environment variable.
	Env string
	// Handler is a function to update a [Config] instance using the
	// environment variable value referenced by Env if the environment variable is set.
	Handler EnvValueHandler
}

// LoadConfig loads configuration into a [Config] instance.
// The configuration is first loaded from a YAML file, then values are overridden from environment variables.
//
// If configFile is empty, the default config file './config.yaml' is used.
func LoadConfig(configFile string, envEntries []EnvConfigEntry) (*Config, error) {
	if configFile == "" {
		configFile = defaultConfigFile
	}
	config, err := loadConfigFromFile(configFile)
	if err != nil {
		return nil, err
	}
	overrideConfigFromEnv(config, envEntries)
	return config, nil
}

func loadConfigFromFile(configFile string) (*Config, error) {
	slog.Debug("reading config file" + configFile)
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

func overrideConfigFromEnv(config *Config, envEntries []EnvConfigEntry) {
	for _, entry := range envEntries {
		val, isSet := os.LookupEnv(entry.Env)
		if isSet {
			slog.Debug("overriding value for " + entry.Env + " from env")
			entry.Handler(val, config)
		}
	}
}
