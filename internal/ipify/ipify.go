// Package ipify provides a client for calling Ipify APIs.
// A client instance can be created by calling [NewClient].
package ipify

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultBaseUrl = "https://api.ipify.org"
)

// IpifyClient implements calls to Ipify APIs.
type IpifyClient struct {
	httpClient *http.Client
	baseUrl    string
}

// NewClient builds a new [IpifyClient] instance.
//
// If baseUrl is empty, the default ipify URL is used.
func NewClient(baseUrl string) *IpifyClient {
	if baseUrl == "" {
		baseUrl = defaultBaseUrl
	}
	return &IpifyClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseUrl:    baseUrl,
	}
}

// IpifyError indicates an error calling an Ipify API.
type IpifyError struct {
	// StatusCode is the status code returned from the API or 500 if not known.
	StatusCode int
	// Message is the response body returned from the API if one was provided.
	Message string
}

// Error returns the error string for an [IpifyError] instance.
func (e *IpifyError) Error() string {
	var message string
	if e.Message != "" {
		message = e.Message
	} else {
		message = "unexpected error"
	}
	return fmt.Sprintf("error calling Ipify: %d - %s", e.StatusCode, message)
}

// GetCurrentIp returns the current WAN IP address
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

	defer func() {
		err = errors.Join(err, resp.Body.Close())
	}()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", &IpifyError{
			StatusCode: resp.StatusCode,
			Message:    string(respData),
		}
	}

	return string(respData), nil
}
