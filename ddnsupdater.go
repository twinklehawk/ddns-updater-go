package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/twinklehawk/ddns-updater-go/internal/config"
	"github.com/twinklehawk/ddns-updater-go/internal/ddnsservice"
	"github.com/twinklehawk/ddns-updater-go/internal/ipprovider"
)

// A ddnsUpdater updates DNS records for hosts to the current IP address.
// The [ddnsUpdater.processDdnsEntry] function is the entry point for updating DDNS data.
type ddnsUpdater struct {
	ipProviders  []ipprovider.CurrentIpProvider
	ddnsServices map[string]ddnsservice.DdnsService
}

func (updater ddnsUpdater) processDdnsEntry(ddnsEntry config.DdnsConfig) error {
	service := updater.ddnsServices[ddnsEntry.Provider]
	if service == nil {
		return fmt.Errorf("ddns provider not supported: %s", ddnsEntry.Provider)
	}

	currentIp, err := updater.getCurrentIpAddress()
	if err != nil {
		return err
	}
	slog.Info("retrieved current IP address: " + currentIp)

	failedCount := 0
	for _, subdomain := range ddnsEntry.Subdomains {
		err := processDdnsHost(subdomain, ddnsEntry.Domain, currentIp, service)
		if err != nil {
			slog.Error("failed to get IP for subdomain", slog.String("subdomain", subdomain), slog.Any("error", err))
			failedCount++
			continue
		}
	}
	if failedCount != 0 {
		return fmt.Errorf("failed to update %d DDNS hosts", failedCount)
	}

	return nil
}

func (updater ddnsUpdater) getCurrentIpAddress() (string, error) {
	var currentIp string
	for _, provider := range updater.ipProviders {
		ip, err := provider.GetCurrentIp(context.Background())
		if err != nil {
			slog.Warn("failed to get IP from provider", slog.Any("error", err))
		} else {
			currentIp = ip
			break
		}
	}
	if currentIp == "" {
		return "", fmt.Errorf("unable to get current ip")
	}
	return currentIp, nil
}

func processDdnsHost(
	subdomain string,
	domain string,
	ip string,
	service ddnsservice.DdnsService,
) error {
	configuredIp, err := service.GetHostIpv4(
		context.Background(),
		subdomain,
		domain,
	)
	if err != nil {
		return fmt.Errorf("failed to get current IP for subdomain %s: %w", subdomain, err)
	}
	if configuredIp == ip {
		slog.Info("skipping unchanged domain " + subdomain)
		return nil
	}
	slog.Info("updating IP for changed domain " + subdomain)
	err = service.UpdateHostIpv4(context.Background(), subdomain, domain, ip)
	if err != nil {
		return fmt.Errorf("failed to update IP for subdomain %s: %w", subdomain, err)
	}
	return nil
}
