// Package config handles reading configuration from a yaml file and env vars.
package config

// Config is the configuration for the ddns-updater-go command.
type Config struct {
	Ddns      DdnsConfig
	Namecheap NamecheapConfig
	Porkbun   PorkbunConfig
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

// PorkbunConfig is the configuration for calling porkbun.
type PorkbunConfig struct {
	ApiKey       string
	SecretApiKey string
}
