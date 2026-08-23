// Package ddnsservice defines the interface for reading and updating DDNS records
// and contains implementations for supported DDNS providers.
package ddnsservice

import "context"

// A DdnsService reads and writes DDNS records for a host.
type DdnsService interface {
	// GetHostIpv4 retrieves the currently configured IP address for the host (subdomain + domain).
	GetHostIpv4(
		ctx context.Context,
		subdomain string,
		domain string,
	) (string, error)

	// UpdateHostIpv4 updates the IP address for a host (subdomain + domain) to the specified IP address.
	UpdateHostIpv4(
		ctx context.Context,
		subdomain string,
		domain string,
		ip string,
	) error
}
