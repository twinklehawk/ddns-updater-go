package ipify

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	DefaultBaseUrl = "https://api.ipify.org"
)

type IpifyClient struct {
	httpClient *http.Client
	baseUrl    string
}

func DefaultClient() *IpifyClient {
	return NewClient(DefaultBaseUrl)
}

func NewClient(baseUrl string) *IpifyClient {
	return &IpifyClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseUrl:    baseUrl,
	}
}

type IpifyError struct {
	StatusCode int
	Message    string
}

func (e *IpifyError) Error() string {
	var message string
	if e.Message != "" {
		message = e.Message
	} else {
		message = "Unexpected error"
	}
	return fmt.Sprintf("Error calling IfConfig: %d - %s", e.StatusCode, message)
}

func (client *IpifyClient) GetCurrentIp(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, client.baseUrl, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "text/plain")
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}

	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", &IpifyError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("status %d: %s", resp.StatusCode, string(respData)),
		}
	}

	return string(respData), nil
}
