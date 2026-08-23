package ddnsservice

import (
	"context"
	"testing"
)

func TestGetHostIpv4(t *testing.T) {
	client := NewNamecheapDdnsService(nil)

	resp, err := client.GetHostIpv4(
		context.Background(),
		"",
		"localhost",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp != "127.0.0.1" {
		t.Errorf("expected response '127.0.0.1', got %q", resp)
	}
}
