package ipprovider

import "context"

// A CurrentIpProvider determines the current IP address.
type CurrentIpProvider interface {
	// GetCurrentIp returns the current WAN IP address.
	GetCurrentIp(ctx context.Context) (string, error)
}
