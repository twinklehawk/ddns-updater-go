package ifconfig

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCurrentIp(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if path := r.URL.Path; path != "/ip" {
			t.Errorf("expected path '/ip', got %s", path)
		}
		if accept := r.Header.Get("Accept"); accept != "text/plain" {
			t.Errorf("expected accept header text/plain, got %s", accept)
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("192.168.1.1"))
		if err != nil {
			t.Fatalf("unable to write response body: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL)

	resp, err := client.GetCurrentIp(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp != "192.168.1.1" {
		t.Errorf("expected response '192.168.1.1', got %q", resp)
	}
}

func TestGetCurrentIpRequestFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Some error"))
		if err != nil {
			t.Fatalf("unable to write response body: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL)

	resp, err := client.GetCurrentIp(context.Background())
	if resp != "" {
		t.Errorf("expected empty response, got %s", resp)
	}
	if err == nil {
		t.Fatalf("missing expected error")
	}

	e, ok := err.(*IfConfigError)
	if !ok {
		t.Fatalf("expected *IfConfigError type, got %T", e)
	}
	if e.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 response code, got %d", e.StatusCode)
	}
	if e.Message != "Some error" {
		t.Errorf("unexpected response body, got %s", e.Message)
	}
}
