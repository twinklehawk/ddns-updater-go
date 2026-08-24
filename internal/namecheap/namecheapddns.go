// Package namecheap provides a client to update Namecheap DDNS IP addresses.
package namecheap

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	defaultBaseUrl = "https://dynamicdns.park-your-domain.com"
)

// NamecheapDdnsClient implements calls to update Namecheap DDNS records.
type NamecheapDdnsClient struct {
	httpClient *http.Client
	baseUrl    string
	password   string
}

// NewClient builds a new [NamecheapDdnsClient] instance.
//
// If baseUrl is empty, the default ipify URL is used.
func NewClient(baseUrl string, password string) *NamecheapDdnsClient {
	if baseUrl == "" {
		baseUrl = defaultBaseUrl
	}
	return &NamecheapDdnsClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseUrl:    baseUrl,
		password:   password,
	}
}

// NamecheapDdnsError indicates an error calling the Namecheap DDNS API.
type NamecheapDdnsError struct {
	StatusCode int
	Message    string
}

// Error returns the error string for a [NamecheapDdnsError] instance.
func (e *NamecheapDdnsError) Error() string {
	var message string
	if e.Message != "" {
		message = e.Message
	} else {
		message = "Unexpected error"
	}
	return fmt.Sprintf("Error calling Namecheap DDNS: %d - %s", e.StatusCode, message)
}

// UpdateHostIpv4 updates the IP address for a host (subdomain + domain) to the specified IP address.
// Namecheap only supports DDNS with IPv4 addresses.
func (client *NamecheapDdnsClient) UpdateHostIpv4(
	ctx context.Context,
	subdomain string,
	domain string,
	ip string,
) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, client.baseUrl+"/update", http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	q.Add("host", subdomain)
	q.Add("domain", domain)
	q.Add("password", client.password)
	q.Add("ip", ip)
	req.URL.RawQuery = q.Encode()

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		err = errors.Join(err, resp.Body.Close())
	}()

	if resp.StatusCode != http.StatusOK {
		return &NamecheapDdnsError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("status %d", resp.StatusCode),
		}
	}

	return err
}
