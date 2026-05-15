package asana

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   "test-token-12345",
			wantErr: false,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: false, // Error happens during API call, not initialization
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.token)
			if client == nil {
				t.Fatal("NewClient returned nil")
			}
			if client.baseURL != "https://app.asana.com/api/1.0" {
				t.Errorf("unexpected baseURL: %s", client.baseURL)
			}
		})
	}
}

// Skip API call tests that require network
func TestClientInitialization(t *testing.T) {
	client := NewClient("test-token")
	
	if client.apiToken != "test-token" {
		t.Errorf("apiToken not set correctly: %s", client.apiToken)
	}
	
	if client.baseURL != "https://app.asana.com/api/1.0" {
		t.Errorf("baseURL not set correctly: %s", client.baseURL)
	}
	
	if client.http == nil {
		t.Error("http client not initialized")
	}
}