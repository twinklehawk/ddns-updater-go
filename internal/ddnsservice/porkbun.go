package ddnsservice

import (
	"context"
	"fmt"

	"github.com/greatliontech/porkbun-go"
)

//go:generate go tool mockgen -destination=../mocks/mock_porkbun.go -package=mocks . PorkbunClient

type PorkbunClient interface {
	// RetrieveDNSRecordsByNameType retrieves DNS records by domain, subdomain and type.
	RetrieveDNSRecordsByNameType(
		ctx context.Context,
		domain string,
		recordType porkbun.RecordType,
		subdomain string,
	) (*porkbun.RetrieveDNSRecordsResponse, error)

	// EditDNSRecordByNameType edits all DNS records matching a subdomain and type.
	EditDNSRecordByNameType(
		ctx context.Context,
		domain string,
		recordType porkbun.RecordType,
		subdomain string,
		record *porkbun.EditDNSRecordByNameTypeRequest,
	) (*porkbun.EditDNSRecordByNameTypeResponse, error)
}

type porkbunDdnsService struct {
	client PorkbunClient
}

// NewPorkbunDdnsService creates a new [DdnsService] for Porkbun.
func NewPorkbunDdnsService(client *porkbun.Client) DdnsService {
	return &porkbunDdnsService{client: client}
}

// See [DdnsService.GetHostIpv4].
func (service *porkbunDdnsService) GetHostIpv4(
	ctx context.Context,
	subdomain string,
	domain string,
) (string, error) {
	resp, err := service.client.RetrieveDNSRecordsByNameType(
		ctx,
		domain,
		porkbun.RecordTypeA,
		subdomain,
	)
	if err != nil {
		return "", fmt.Errorf("unable to fetch dns records for subdomain %s: %w", subdomain, err)
	}
	switch len(resp.Records) {
	case 0:
		return "", nil
	case 1:
		return resp.Records[0].Content, nil
	default:
		return "", fmt.Errorf("multiple A records found for subdomain %s", subdomain)
	}
}

// See [DdnsService.UpdateHostIpv4].
func (service *porkbunDdnsService) UpdateHostIpv4(
	ctx context.Context,
	subdomain string,
	domain string,
	ip string,
) error {
	req := porkbun.EditDNSRecordByNameTypeRequest{
		Content: ip,
		TTL:     600,
		Prio:    0,
		Notes:   nil,
	}
	// TODO need to create record if it doesn't already exist?
	_, err := service.client.EditDNSRecordByNameType(ctx, domain, porkbun.RecordTypeA, subdomain, &req)
	if err != nil {
		return fmt.Errorf("failed to update A record for subdomain %s: %w", subdomain, err)
	}
	return nil
}
