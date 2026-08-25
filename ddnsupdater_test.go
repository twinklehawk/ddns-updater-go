package main

import (
	"testing"

	"github.com/twinklehawk/ddns-updater-go/internal/config"
	"github.com/twinklehawk/ddns-updater-go/internal/ddnsservice"
	"github.com/twinklehawk/ddns-updater-go/internal/ipprovider"
	"github.com/twinklehawk/ddns-updater-go/internal/mocks"
	"go.uber.org/mock/gomock"
)

func TestProcessDdnsEntry_Success(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockIpProvider := mocks.NewMockCurrentIpProvider(ctrl)
	mockIpProvider.EXPECT().GetCurrentIp(gomock.Any()).Return("192.168.1.2", nil)
	ipProviders := []ipprovider.CurrentIpProvider{mockIpProvider}
	mockDdnsService := mocks.NewMockDdnsService(ctrl)
	mockDdnsService.EXPECT().GetHostIpv4(gomock.Any(), "test", "test.com").Return("192.168.1.1", nil)
	mockDdnsService.EXPECT().UpdateHostIpv4(gomock.Any(), "test", "test.com", "192.168.1.2").Return(nil)
	ddnsServices := make(map[string]ddnsservice.DdnsService)
	ddnsServices["test"] = mockDdnsService
	entry := config.DdnsConfig{Domain: "test.com", Provider: "test", Subdomains: []string{"test"}}
	updater := ddnsUpdater{ipProviders: ipProviders, ddnsServices: ddnsServices}

	err := updater.processDdnsEntry(entry)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProcessDdnsEntry_UnchangedIp(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockIpProvider := mocks.NewMockCurrentIpProvider(ctrl)
	mockIpProvider.EXPECT().GetCurrentIp(gomock.Any()).Return("192.168.1.1", nil)
	ipProviders := []ipprovider.CurrentIpProvider{mockIpProvider}
	mockDdnsService := mocks.NewMockDdnsService(ctrl)
	mockDdnsService.EXPECT().GetHostIpv4(gomock.Any(), "test", "test.com").Return("192.168.1.1", nil)
	// no update call should be made since IPs are equal
	ddnsServices := make(map[string]ddnsservice.DdnsService)
	ddnsServices["test"] = mockDdnsService
	entry := config.DdnsConfig{Domain: "test.com", Provider: "test", Subdomains: []string{"test"}}
	updater := ddnsUpdater{ipProviders: ipProviders, ddnsServices: ddnsServices}

	err := updater.processDdnsEntry(entry)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
