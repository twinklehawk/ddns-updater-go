package namecheap

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

const (
	DefaultBaseUrl = "https://dynamicdns.park-your-domain.com"
)

type NamecheapDdnsClient interface {
	GetHostIpv4(ctx context.Context,
		subdomain string,
		domain string,
	) (string, error)

	UpdateHostIpv4(ctx context.Context,
		subdomain string,
		domain string,
		ip string) error
}

func NewClient(baseUrl string, password string) NamecheapDdnsClient {
	return &internalNamecheapDdnsClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseUrl:    baseUrl,
		password:   password,
	}
}

type internalNamecheapDdnsClient struct {
	httpClient *http.Client
	baseUrl    string
	password   string
}

type NamecheapDdnsError struct {
	StatusCode int
	Message    string
}

func (e *NamecheapDdnsError) Error() string {
	var message string
	if e.Message != "" {
		message = e.Message
	} else {
		message = "Unexpected error"
	}
	return fmt.Sprintf("Error calling Namecheap DDNS: %d - %s", e.StatusCode, message)
}

func (client *internalNamecheapDdnsClient) GetHostIpv4(
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
		return "", &NamecheapDdnsError{StatusCode: 500, Message: "No IP found for host"}
	}
	return addr[0].String(), nil
}

func (client *internalNamecheapDdnsClient) UpdateHostIpv4(
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
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &NamecheapDdnsError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("status %d", resp.StatusCode),
		}
	}

	return nil
}
