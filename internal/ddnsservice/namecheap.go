package ddnsservice

import (
	"context"
	"fmt"
	"net"

	"github.com/twinklehawk/ddns-updater-go/internal/namecheap"
)

type namecheapDdnsService struct {
	client *namecheap.NamecheapDdnsClient
}

// NewNamecheapDdnsService creates a new [DdnsService] instance for Namecheap.
func NewNamecheapDdnsService(client *namecheap.NamecheapDdnsClient) DdnsService {
	return &namecheapDdnsService{client: client}
}

// See [DdnsService.GetHostIpv4].
//
// Namecheap does not support retrieving the currently configured IP address, so this function
// looks up the IP address associated to the host using a normal DNS lookup.
func (service *namecheapDdnsService) GetHostIpv4(
	ctx context.Context,
	subdomain string,
	domain string,
) (string, error) {
	var host string
	if subdomain != "" {
		host = subdomain + "." + domain
	} else {
		host = domain
	}
	addr, err := net.DefaultResolver.LookupIP(ctx, "ip4", host)
	if err != nil {
		return "", fmt.Errorf("unable to look up IP for host: %w", err)
	}
	if len(addr) == 0 {
		return "", fmt.Errorf("no IP found for host")
	}
	return addr[0].String(), nil
}

// See [DdnsService.UpdateHostIpv4].
func (service *namecheapDdnsService) UpdateHostIpv4(
	ctx context.Context,
	subdomain string,
	domain string,
	ip string,
) error {
	// TODO needs a unit test
	return service.client.UpdateHostIpv4(ctx, subdomain, domain, ip)
}
