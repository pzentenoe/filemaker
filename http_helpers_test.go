package filemaker

import (
	"context"
	"net/http"
	"testing"
)

func TestBearerAuthHeader(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  string
	}{
		{
			name:  "valid token",
			token: "test-token-123",
			want:  "Bearer test-token-123",
		},
		{
			name:  "empty token",
			token: "",
			want:  "Bearer ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := bearerAuthHeader(tt.token)
			got := headers.Get("Authorization")
			if got != tt.want {
				t.Errorf("bearerAuthHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithContentType(t *testing.T) {
	tests := []struct {
		name        string
		headers     http.Header
		contentType string
		want        string
	}{
		{
			name:        "add to existing headers",
			headers:     http.Header{"X-Custom": []string{"value"}},
			contentType: "application/json",
			want:        "application/json",
		},
		{
			name:        "create new headers",
			headers:     nil,
			contentType: "text/plain",
			want:        "text/plain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := withContentType(tt.headers, tt.contentType)
			got := headers.Get("Content-Type")
			if got != tt.want {
				t.Errorf("withContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnsureContext(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		wantNil bool
	}{
		{
			name:    "nil context returns background",
			ctx:     nil,
			wantNil: false,
		},
		{
			name:    "non-nil context is preserved",
			ctx:     context.Background(),
			wantNil: false,
		},
		{
			name:    "todo context is preserved",
			ctx:     context.TODO(),
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ensureContext(tt.ctx)
			if (got == nil) != tt.wantNil {
				t.Errorf("ensureContext() returned nil = %v, want nil = %v", got == nil, tt.wantNil)
			}
		})
	}
}
