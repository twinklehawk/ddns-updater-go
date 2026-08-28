package ddnsservice

import (
	"context"
	"testing"

	"github.com/greatliontech/porkbun-go"
	"github.com/twinklehawk/ddns-updater-go/internal/mocks"
	"go.uber.org/mock/gomock"
)

func TestGetHostIpv4_Success(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockClient := mocks.NewMockPorkbunClient(ctrl)
	recordsResponse := porkbun.RetrieveDNSRecordsResponse{
		Status: "success",
		Records: []porkbun.DNSRecord{{
			ID:      "1",
			Name:    "test.test.com",
			Type:    "A",
			Content: "192.168.1.2",
			TTL:     "600",
		}},
	}
	mockClient.EXPECT().RetrieveDNSRecordsByNameType(
		gomock.Any(),
		"test.com",
		porkbun.RecordTypeA,
		"test",
	).Return(&recordsResponse, nil)
	service := porkbunDdnsService{client: mockClient}

	ip, err := service.GetHostIpv4(context.Background(), "test", "test.com")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ip != "192.168.1.2" {
		t.Errorf("expected ip '192.168.1.2', found %s", ip)
	}
}

func TestGetHostIpv4_ReturnsEmptyWhenNoRecord(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockClient := mocks.NewMockPorkbunClient(ctrl)
	recordsResponse := porkbun.RetrieveDNSRecordsResponse{
		Status:  "success",
		Records: []porkbun.DNSRecord{},
	}
	mockClient.EXPECT().RetrieveDNSRecordsByNameType(
		gomock.Any(),
		"test.com",
		porkbun.RecordTypeA,
		"test",
	).Return(&recordsResponse, nil)
	service := porkbunDdnsService{client: mockClient}

	ip, err := service.GetHostIpv4(context.Background(), "test", "test.com")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ip != "" {
		t.Errorf("expected ip '', found %s", ip)
	}
}

func TestGetHostIpv4_FailsWhenMultipleRecords(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockClient := mocks.NewMockPorkbunClient(ctrl)
	recordsResponse := porkbun.RetrieveDNSRecordsResponse{
		Status: "success",
		Records: []porkbun.DNSRecord{
			{
				ID:      "1",
				Name:    "test.test.com",
				Type:    "A",
				Content: "192.168.1.2",
				TTL:     "600",
			},
			{
				ID:      "2",
				Name:    "test.test.com",
				Type:    "A",
				Content: "192.168.1.3",
				TTL:     "600",
			},
		},
	}
	mockClient.EXPECT().RetrieveDNSRecordsByNameType(
		gomock.Any(),
		"test.com",
		porkbun.RecordTypeA,
		"test",
	).Return(&recordsResponse, nil)
	service := porkbunDdnsService{client: mockClient}

	ip, err := service.GetHostIpv4(context.Background(), "test", "test.com")

	if err == nil {
		t.Errorf("expected error but none found")
	}
	if ip != "" {
		t.Errorf("expected empty ip, found %s", ip)
	}
}

func TestUpdateHostIpv4_Success(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockClient := mocks.NewMockPorkbunClient(ctrl)
	expectedRequest := porkbun.EditDNSRecordByNameTypeRequest{
		Content: "192.168.1.1",
		TTL:     600,
		Prio:    0,
	}
	mockClient.EXPECT().EditDNSRecordByNameType(
		gomock.Any(),
		"test.com",
		porkbun.RecordTypeA,
		"test",
		&expectedRequest,
	).Return(&porkbun.EditDNSRecordByNameTypeResponse{}, nil)
	service := porkbunDdnsService{client: mockClient}

	err := service.UpdateHostIpv4(
		context.Background(),
		"test",
		"test.com",
		"192.168.1.1",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
