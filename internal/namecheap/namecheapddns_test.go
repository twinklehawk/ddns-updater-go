package namecheap

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateHostIpv4(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if path := r.URL.Path; path != "/update" {
			t.Errorf("expected path '/update', got %s", path)
		}
		if param := r.URL.Query().Get("host"); param != "test" {
			t.Errorf("expected query param host to be 'test', got %s", param)
		}
		if param := r.URL.Query().Get("domain"); param != "domain.com" {
			t.Errorf("expected query param domain to be 'domain.com', got %s", param)
		}
		if param := r.URL.Query().Get("ip"); param != "192.168.1.1" {
			t.Errorf("expected query param ip to be '192.168.1.1', got %s", param)
		}
		if param := r.URL.Query().Get("password"); param != "test-pass" {
			t.Errorf("expected query param password to be 'test-pass', got %s", param)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-pass")

	err := client.UpdateHostIpv4(
		context.Background(),
		"test",
		"domain.com",
		"192.168.1.1",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateHostIpv4Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Some error"))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-pass")

	err := client.UpdateHostIpv4(
		context.Background(),
		"test",
		"domain.com",
		"192.168.1.1",
	)
	if err == nil {
		t.Fatalf("missing expected error")
	}

	e, ok := err.(*NamecheapDdnsError)
	if !ok {
		t.Fatalf("expected *NamecheapDdnsError type, got %T", e)
	}
	if e.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 response code, got %d", e.StatusCode)
	}
}
