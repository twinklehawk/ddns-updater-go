// Ddns-updater-go updates the IP address for hosts.
// The command assumes it is running on a computer using the desired IP address for the DDNS record.
// For each configured host, ddns-updater-go:
//   - determines the current WAN IP address
//   - compares the configured IP address for the host matches to the current IP address
//   - updates the configured IP address for the host to the current IP address if different
//
// Ddns-updater-go configuration is loaded from config.yaml in the current directory and
// environment variables.
// See the config package for details on configuration.
package main

import (
	"log/slog"
	"os"

	"github.com/greatliontech/porkbun-go"
	"github.com/twinklehawk/ddns-updater-go/internal/config"
	"github.com/twinklehawk/ddns-updater-go/internal/ddnsservice"
	"github.com/twinklehawk/ddns-updater-go/internal/ifconfig"
	"github.com/twinklehawk/ddns-updater-go/internal/ipify"
	"github.com/twinklehawk/ddns-updater-go/internal/ipprovider"
	"github.com/twinklehawk/ddns-updater-go/internal/namecheap"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := run(); err != nil {
		slog.Error("run failed", slog.Any("error", err))
		os.Exit(1)
	}
}

func run() error {
	slog.Info("loading config")
	config, err := config.LoadConfig("", getEnvConfigMap())
	if err != nil {
		return err
	}

	updater := ddnsUpdater{
		ipProviders:  buildIpProviders(),
		ddnsServices: buildDdnsServices(config),
	}

	slog.Info("processing ddns for domain: " + config.Ddns.Domain)
	err = updater.processDdnsEntry(config.Ddns)
	if err != nil {
		return err
	}
	slog.Info("finished processing ddns entries")

	return nil
}

func buildIpProviders() []ipprovider.CurrentIpProvider {
	return []ipprovider.CurrentIpProvider{
		ipify.NewClient(""),
		ifconfig.NewClient(""),
	}
}

func buildDdnsServices(config *config.Config) map[string]ddnsservice.DdnsService {
	ddnsServices := make(map[string]ddnsservice.DdnsService)
	namecheapClient := namecheap.NewClient("", config.Namecheap.Password)
	ddnsServices["namecheap"] = ddnsservice.NewNamecheapDdnsService(namecheapClient)
	porkbunClient := porkbun.NewClient(config.Porkbun.ApiKey, config.Porkbun.SecretApiKey)
	ddnsServices["porkbun"] = ddnsservice.NewPorkbunDdnsService(porkbunClient)
	return ddnsServices
}

func getEnvConfigMap() []config.EnvConfigEntry {
	return []config.EnvConfigEntry{
		{
			Env:     "NAMECHEAP_PASSWORD",
			Handler: func(s string, c *config.Config) { c.Namecheap.Password = s },
		},
		{
			Env:     "PORKBUN_API_KEY",
			Handler: func(s string, c *config.Config) { c.Porkbun.ApiKey = s },
		},
		{
			Env:     "PORKBUN_SECRET_API_KEY",
			Handler: func(s string, c *config.Config) { c.Porkbun.SecretApiKey = s },
		},
	}
}
