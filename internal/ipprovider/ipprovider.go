package ipprovider

import "context"

//go:generate go tool mockgen -destination=../mocks/mock_ipprovider.go -package=mocks . CurrentIpProvider

// A CurrentIpProvider determines the current IP address.
type CurrentIpProvider interface {
	// GetCurrentIp returns the current WAN IP address.
	GetCurrentIp(ctx context.Context) (string, error)
}
