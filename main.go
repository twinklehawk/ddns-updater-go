// Ddns-updater-go updates the IP address for hosts.
// The command assumes it is running on a computer using the desired IP address for the DDNS record.
// For each configured host, ddns-updater-go:
//   - determines the current WAN IP address
//   - compares the configured IP address for the host matches to the current IP address
//   - updates the configured IP address for the host to the current IP address if different
//
// Ddns-updater-go configuration is loaded from config.yaml in the current directory.
// See [config] for details on configuration.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

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

	slog.Info("processing ddns entries")

	config, err := config.ReadConfig("")
	if err != nil {
		slog.Error("unable to read config", slog.Any("error", err))
		return
	}

	ipProviders := buildIpProviders()
	ddnsClients := buildDdnsClients(config)

	currentIp, err := getCurrentIpAddress(ipProviders)
	if err != nil {
		slog.Error("unable to get current IP address", slog.Any("error", err))
		return
	}
	slog.Info("retrieved current IP address: " + currentIp)

	slog.Info("processing ddns for domain: " + config.Ddns.Domain)
	processDdnsEntry(config.Ddns, currentIp, ddnsClients)

	slog.Info("finished processing ddns entries")
}

func buildIpProviders() []ipprovider.CurrentIpProvider {
	return []ipprovider.CurrentIpProvider{
		ipify.NewClient(""),
		ifconfig.NewClient(""),
	}
}

func buildDdnsClients(config *config.Config) map[string]ddnsservice.DdnsService {
	ddnsClients := make(map[string]ddnsservice.DdnsService)
	ddnsClients["namecheap"] = ddnsservice.NewNamecheapDdnsService(namecheap.NewClient("", config.Namecheap.Password))
	return ddnsClients
}

func getCurrentIpAddress(providers []ipprovider.CurrentIpProvider) (string, error) {
	var currentIp string
	for _, provider := range providers {
		ip, err := provider.GetCurrentIp(context.Background())
		if err != nil {
			slog.Warn("failed to get IP from provider", slog.Any("error", err))
		} else {
			currentIp = ip
		}
	}
	if currentIp == "" {
		return "", fmt.Errorf("unable to get current ip")
	}
	return currentIp, nil
}

func processDdnsEntry(
	ddnsEntry config.DdnsConfig,
	currentIp string,
	clients map[string]ddnsservice.DdnsService,
) error {
	client := clients[ddnsEntry.Provider]
	if client == nil {
		return fmt.Errorf("ddns provider not supported: %s", ddnsEntry.Provider)
	}
	for _, subdomain := range ddnsEntry.Subdomains {
		// get currently configured IP
		configuredIp, err := client.GetHostIpv4(
			context.Background(),
			subdomain,
			ddnsEntry.Domain,
		)
		if err != nil {
			slog.Error("failed to get current IP for subdomain", slog.String("subdomain", subdomain), slog.Any("error", err))
			continue
		}
		if configuredIp == currentIp {
			slog.Info("skipping unchanged domain", slog.String("subdomain", subdomain))
			continue
		}
		slog.Info("updating IP for changed domain", slog.String("subdomain", subdomain))
		err = client.UpdateHostIpv4(context.Background(), subdomain, ddnsEntry.Domain, currentIp)
		if err != nil {
			slog.Error("failed to update current IP for subdomain", slog.String("subdomain", subdomain), slog.Any("error", err))
		}
	}
	return nil
}
