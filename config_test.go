package filemaker

import (
	"crypto/tls"
	"net/http"
	"testing"
	"time"
)

func TestDefaultHTTPClientConfig(t *testing.T) {
	config := DefaultHTTPClientConfig()

	if config.Timeout != 30*time.Second {
		t.Errorf("Timeout = %v, want 30s", config.Timeout)
	}

	if config.DialTimeout != 10*time.Second {
		t.Errorf("DialTimeout = %v, want 10s", config.DialTimeout)
	}

	if config.TLSHandshakeTimeout != 10*time.Second {
		t.Errorf("TLSHandshakeTimeout = %v, want 10s", config.TLSHandshakeTimeout)
	}

	if config.ResponseHeaderTimeout != 10*time.Second {
		t.Errorf("ResponseHeaderTimeout = %v, want 10s", config.ResponseHeaderTimeout)
	}

	if config.IdleConnTimeout != 90*time.Second {
		t.Errorf("IdleConnTimeout = %v, want 90s", config.IdleConnTimeout)
	}

	if config.MaxIdleConns != 100 {
		t.Errorf("MaxIdleConns = %v, want 100", config.MaxIdleConns)
	}

	if config.MaxIdleConnsPerHost != 10 {
		t.Errorf("MaxIdleConnsPerHost = %v, want 10", config.MaxIdleConnsPerHost)
	}

	if config.MaxConnsPerHost != 0 {
		t.Errorf("MaxConnsPerHost = %v, want 0", config.MaxConnsPerHost)
	}

	if config.DisableKeepAlives {
		t.Error("DisableKeepAlives should be false")
	}

	if config.DisableCompression {
		t.Error("DisableCompression should be false")
	}

	if config.TLSConfig != nil {
		t.Error("TLSConfig should be nil")
	}
}

func TestHTTPClientConfig_BuildHTTPClient(t *testing.T) {
	config := DefaultHTTPClientConfig()
	client := config.BuildHTTPClient()

	if client == nil {
		t.Fatal("BuildHTTPClient() returned nil")
	}

	if client.Timeout != config.Timeout {
		t.Errorf("client.Timeout = %v, want %v", client.Timeout, config.Timeout)
	}

	if client.Transport == nil {
		t.Fatal("client.Transport is nil")
	}
}

func TestHTTPClientConfig_BuildHTTPClient_CustomTimeout(t *testing.T) {
	config := &HTTPClientConfig{
		Timeout:               5 * time.Second,
		DialTimeout:           2 * time.Second,
		TLSHandshakeTimeout:   2 * time.Second,
		ResponseHeaderTimeout: 2 * time.Second,
		IdleConnTimeout:       30 * time.Second,
		MaxIdleConns:          50,
		MaxIdleConnsPerHost:   5,
		MaxConnsPerHost:       10,
	}

	client := config.BuildHTTPClient()

	if client.Timeout != 5*time.Second {
		t.Errorf("client.Timeout = %v, want 5s", client.Timeout)
	}
}

func TestHTTPClientConfig_BuildHTTPClient_DisableKeepAlives(t *testing.T) {
	config := DefaultHTTPClientConfig()
	config.DisableKeepAlives = true

	client := config.BuildHTTPClient()

	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatal("client.Transport is not *http.Transport")
	}

	if !transport.DisableKeepAlives {
		t.Error("DisableKeepAlives should be true")
	}
}

func TestHTTPClientConfig_BuildHTTPClient_DisableCompression(t *testing.T) {
	config := DefaultHTTPClientConfig()
	config.DisableCompression = true

	client := config.BuildHTTPClient()

	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatal("client.Transport is not *http.Transport")
	}

	if !transport.DisableCompression {
		t.Error("DisableCompression should be true")
	}
}

func TestHTTPClientConfig_BuildHTTPClient_CustomTLSConfig(t *testing.T) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	config := DefaultHTTPClientConfig()
	config.TLSConfig = tlsConfig

	client := config.BuildHTTPClient()

	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatal("client.Transport is not *http.Transport")
	}

	if transport.TLSClientConfig != tlsConfig {
		t.Error("TLSClientConfig does not match")
	}

	if !transport.TLSClientConfig.InsecureSkipVerify {
		t.Error("TLSClientConfig.InsecureSkipVerify should be true")
	}
}

func TestHTTPClientConfig_BuildHTTPClient_MaxConnsPerHost(t *testing.T) {
	config := DefaultHTTPClientConfig()
	config.MaxConnsPerHost = 20

	client := config.BuildHTTPClient()

	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatal("client.Transport is not *http.Transport")
	}

	if transport.MaxConnsPerHost != 20 {
		t.Errorf("MaxConnsPerHost = %v, want 20", transport.MaxConnsPerHost)
	}
}

func TestHTTPClientConfig_BuildHTTPClient_HTTP2(t *testing.T) {
	config := DefaultHTTPClientConfig()
	client := config.BuildHTTPClient()

	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatal("client.Transport is not *http.Transport")
	}

	if !transport.ForceAttemptHTTP2 {
		t.Error("ForceAttemptHTTP2 should be true")
	}
}
