// Package ifconfig provides a client for calling IfConfig APIs.
// A client instance can be created by calling [NewClient].
package ifconfig

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultBaseUrl = "https://ifconfig.me"
)

// IfConfigClient implements calls to IfConfig APIs.
type IfConfigClient struct {
	httpClient *http.Client
	baseUrl    string
}

// NewClient builds a new [IfConfigClient] instance.
//
// If baseUrl is empty, the default ifconfig URL is used.
func NewClient(baseUrl string) *IfConfigClient {
	if baseUrl == "" {
		baseUrl = defaultBaseUrl
	}
	return &IfConfigClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseUrl:    baseUrl,
	}
}

// IfConfigError indicates an error calling an IfConfig API.
type IfConfigError struct {
	// StatusCode is the status code returned from the API or 500 if not known.
	StatusCode int
	// Message is the response body returned from the API if one was provided.
	Message string
}

// Error returns the error string for an [IfConfigError] instance.
func (e *IfConfigError) Error() string {
	var message string
	if e.Message != "" {
		message = e.Message
	} else {
		message = "Unexpected error"
	}
	return fmt.Sprintf("Error calling IfConfig: %d - %s", e.StatusCode, message)
}

// GetCurrentIp returns the current WAN IP address
func (client *IfConfigClient) GetCurrentIp(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, client.baseUrl+"/ip", http.NoBody)
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
		return "", &IfConfigError{
			StatusCode: resp.StatusCode,
			Message:    string(respData),
		}
	}

	return string(respData), nil
}
