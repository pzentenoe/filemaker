package filemaker

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

// HTTPClientConfig defines configuration for the HTTP client with comprehensive timeout settings.
type HTTPClientConfig struct {
	// Timeout is the overall request timeout (includes connection, TLS, headers, and body).
	// Default: 30 seconds
	Timeout time.Duration

	// DialTimeout is the maximum time to wait for a TCP connection to be established.
	// Default: 10 seconds
	DialTimeout time.Duration

	// TLSHandshakeTimeout is the maximum time to wait for a TLS handshake.
	// Default: 10 seconds
	TLSHandshakeTimeout time.Duration

	// ResponseHeaderTimeout is the maximum time to wait for response headers.
	// Default: 10 seconds
	ResponseHeaderTimeout time.Duration

	// IdleConnTimeout is the maximum time an idle connection remains in the pool.
	// Default: 90 seconds
	IdleConnTimeout time.Duration

	// MaxIdleConns is the maximum number of idle connections across all hosts.
	// Default: 100
	MaxIdleConns int

	// MaxIdleConnsPerHost is the maximum number of idle connections per host.
	// Default: 10
	MaxIdleConnsPerHost int

	// MaxConnsPerHost is the maximum number of connections per host (0 = unlimited).
	// Default: 0 (unlimited)
	MaxConnsPerHost int

	// DisableKeepAlives disables HTTP keep-alives.
	// Default: false
	DisableKeepAlives bool

	// DisableCompression disables automatic gzip compression.
	// Default: false
	DisableCompression bool

	// TLSConfig provides custom TLS configuration.
	// Default: nil (uses default TLS config)
	TLSConfig *tls.Config
}

// DefaultHTTPClientConfig returns an HTTPClientConfig with sensible defaults.
func DefaultHTTPClientConfig() *HTTPClientConfig {
	return &HTTPClientConfig{
		Timeout:               30 * time.Second,
		DialTimeout:           10 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		IdleConnTimeout:       90 * time.Second,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		MaxConnsPerHost:       0,
		DisableKeepAlives:     false,
		DisableCompression:    false,
		TLSConfig:             nil,
	}
}

// BuildHTTPClient creates an http.Client with the specified configuration.
func (c *HTTPClientConfig) BuildHTTPClient() *http.Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   c.DialTimeout,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   c.TLSHandshakeTimeout,
		ResponseHeaderTimeout: c.ResponseHeaderTimeout,
		IdleConnTimeout:       c.IdleConnTimeout,
		MaxIdleConns:          c.MaxIdleConns,
		MaxIdleConnsPerHost:   c.MaxIdleConnsPerHost,
		MaxConnsPerHost:       c.MaxConnsPerHost,
		DisableKeepAlives:     c.DisableKeepAlives,
		DisableCompression:    c.DisableCompression,
		TLSClientConfig:       c.TLSConfig,
		ForceAttemptHTTP2:     true,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   c.Timeout,
	}
}
