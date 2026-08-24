package config

import (
	"os"
	"testing"
)

func TestLoadConfigFromFile(t *testing.T) {
	configData := `
ddns:
  domain: example.com
  provider: namecheap
  subdomains:
    - test
    - test2
namecheap:
  password: test-pass
`
	file, err := createConfigFile(configData)
	if err != nil {
		t.Fatalf("failed to create temp config file: %v", err)
	}
	defer os.Remove(file)

	config, err := LoadConfig(file, nil)
	if err != nil {
		t.Fatalf("unexpected error when reading config file: %v", err)
	}
	if config == nil {
		t.Fatalf("expected non-nil config, received nil")
	}

	if config.Ddns.Domain != "example.com" {
		t.Errorf("expected 'example.com' domain, found %s", config.Ddns.Domain)
	}
	if config.Ddns.Provider != "namecheap" {
		t.Errorf("expected 'namecheap' provider, found %s", config.Ddns.Provider)
	}
	if len := len(config.Ddns.Subdomains); len != 2 {
		t.Errorf("expected 2 subdomains, found %d", len)
	}
	if config.Ddns.Subdomains[0] != "test" {
		t.Errorf("expected 'test' domain, found %s", config.Ddns.Subdomains[0])
	}
	if config.Ddns.Subdomains[1] != "test2" {
		t.Errorf("expected 'test2' subdomain, found %s", config.Ddns.Subdomains[1])
	}
	if config.Namecheap.Password != "test-pass" {
		t.Errorf("expected 'test-pass' namecheap password, found %s", config.Namecheap.Password)
	}
}

func TestLoadConfigEnvOverridesFileValues(t *testing.T) {
	configData := `
ddns:
  domain: example.com
  provider: namecheap
  subdomains:
    - test
    - test2
namecheap:
  password: test-pass
`
	file, err := createConfigFile(configData)
	if err != nil {
		t.Fatalf("failed to create temp config file: %v", err)
	}
	defer os.Remove(file)

	envkey := "TEST_ENV_KEY"
	err = os.Setenv(envkey, "test-pass-2")
	if err != nil {
		t.Fatalf("failed to set envvar: %v", err)
	}
	defer os.Unsetenv(envkey)

	envEntries := []EnvConfigEntry{
		{Env: envkey, Handler: func(value string, config *Config) { config.Namecheap.Password = value }},
	}
	config, err := LoadConfig(file, envEntries)
	if err != nil {
		t.Fatalf("unexpected error when reading config file: %v", err)
	}
	if config == nil {
		t.Fatalf("expected non-nil config, received nil")
	}

	if config.Namecheap.Password != "test-pass-2" {
		t.Errorf("expected 'test-pass-2' namecheap password, found %s", config.Namecheap.Password)
	}
}

func createConfigFile(data string) (string, error) {
	file, err := os.CreateTemp("", "loader-test")
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}
