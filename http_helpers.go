package filemaker

import (
	"context"
	"net/http"
)

// bearerAuthHeader creates an HTTP header with Bearer token authentication.
func bearerAuthHeader(token string) http.Header {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+token)
	return headers
}

// withContentType adds a Content-Type header to existing headers.
// If headers is nil, a new http.Header is created.
func withContentType(headers http.Header, contentType string) http.Header {
	if headers == nil {
		headers = http.Header{}
	}
	headers.Set("Content-Type", contentType)
	return headers
}

// ensureContext returns the provided context if not nil, otherwise returns context.Background().
func ensureContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

// jsonContentType is a constant for JSON content type.
const jsonContentType = "application/json"
