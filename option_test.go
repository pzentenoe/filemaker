package filemaker

import (
	"net/http"
	"testing"
	"time"
)

func TestSetURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid url",
			url:     "https://filemaker.example.com",
			wantErr: false,
		},
		{
			name:    "empty url",
			url:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{}
			err := SetURL(tt.url)(client)

			if (err != nil) != tt.wantErr {
				t.Errorf("SetURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && client.url != tt.url {
				t.Errorf("client.url = %v, want %v", client.url, tt.url)
			}
		})
	}
}

func TestSetBasicAuth(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{
			name:     "valid credentials",
			username: "admin",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "empty username",
			username: "",
			password: "password123",
			wantErr:  true,
		},
		{
			name:     "empty password",
			username: "admin",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{}
			err := SetBasicAuth(tt.username, tt.password)(client)

			if (err != nil) != tt.wantErr {
				t.Errorf("SetBasicAuth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && client.authProvider == nil {
				t.Error("client.authProvider should not be nil")
			}
		})
	}
}

func TestSetAuthProvider(t *testing.T) {
	tests := []struct {
		name     string
		provider AuthProvider
		wantErr  bool
	}{
		{
			name:     "valid provider",
			provider: WithBasicAuth("user", "pass"),
			wantErr:  false,
		},
		{
			name:     "nil provider",
			provider: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{}
			err := SetAuthProvider(tt.provider)(client)

			if (err != nil) != tt.wantErr {
				t.Errorf("SetAuthProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && client.authProvider == nil {
				t.Error("client.authProvider should not be nil")
			}
		})
	}
}

func TestSetVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{
			name:    "custom version",
			version: "v1",
			want:    "v1",
		},
		{
			name:    "empty version uses default",
			version: "",
			want:    DefaultVersion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{}
			err := SetVersion(tt.version)(client)

			if err != nil {
				t.Errorf("SetVersion() error = %v", err)
				return
			}

			if client.version != tt.want {
				t.Errorf("client.version = %v, want %v", client.version, tt.want)
			}
		})
	}
}

func TestSetHttpClient(t *testing.T) {
	customClient := &http.Client{Timeout: 5 * time.Second}

	tests := []struct {
		name   string
		client *http.Client
		want   *http.Client
	}{
		{
			name:   "custom http client",
			client: customClient,
			want:   customClient,
		},
		{
			name:   "nil client uses default",
			client: nil,
			want:   http.DefaultClient,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{}
			err := SetHttpClient(tt.client)(client)

			if err != nil {
				t.Errorf("SetHttpClient() error = %v", err)
				return
			}

			if client.httpClient != tt.want {
				t.Errorf("client.httpClient = %v, want %v", client.httpClient, tt.want)
			}
		})
	}
}

func TestSetRetryConfig(t *testing.T) {
	customConfig := &RetryConfig{MaxRetries: 5}

	tests := []struct {
		name   string
		config *RetryConfig
	}{
		{
			name:   "custom retry config",
			config: customConfig,
		},
		{
			name:   "nil config uses default",
			config: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{}
			err := SetRetryConfig(tt.config)(client)

			if err != nil {
				t.Errorf("SetRetryConfig() error = %v", err)
				return
			}

			if tt.config != nil && client.retryConfig != tt.config {
				t.Error("client.retryConfig should match custom config")
			}

			if tt.config == nil && client.retryConfig == nil {
				t.Error("client.retryConfig should not be nil when using default")
			}
		})
	}
}

func TestSetMaxRetries(t *testing.T) {
	client := &Client{}
	err := SetMaxRetries(5)(client)

	if err != nil {
		t.Fatalf("SetMaxRetries() error = %v", err)
	}

	if client.retryConfig == nil {
		t.Fatal("retryConfig should not be nil")
	}

	if client.retryConfig.MaxRetries != 5 {
		t.Errorf("MaxRetries = %v, want 5", client.retryConfig.MaxRetries)
	}
}

func TestSetRetryWaitTime(t *testing.T) {
	client := &Client{}
	min := 2 * time.Second
	max := 20 * time.Second

	err := SetRetryWaitTime(min, max)(client)

	if err != nil {
		t.Fatalf("SetRetryWaitTime() error = %v", err)
	}

	if client.retryConfig == nil {
		t.Fatal("retryConfig should not be nil")
	}

	if client.retryConfig.MinWaitTime != min {
		t.Errorf("MinWaitTime = %v, want %v", client.retryConfig.MinWaitTime, min)
	}

	if client.retryConfig.MaxWaitTime != max {
		t.Errorf("MaxWaitTime = %v, want %v", client.retryConfig.MaxWaitTime, max)
	}
}

func TestDisableRetry(t *testing.T) {
	client := &Client{}
	err := DisableRetry()(client)

	if err != nil {
		t.Fatalf("DisableRetry() error = %v", err)
	}

	if client.retryConfig == nil {
		t.Fatal("retryConfig should not be nil")
	}

	if client.retryConfig.MaxRetries != 0 {
		t.Errorf("MaxRetries = %v, want 0", client.retryConfig.MaxRetries)
	}
}

func TestSetTimeout(t *testing.T) {
	timeout := 15 * time.Second

	tests := []struct {
		name       string
		httpClient *http.Client
	}{
		{
			name:       "with existing client",
			httpClient: &http.Client{},
		},
		{
			name:       "without existing client",
			httpClient: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{httpClient: tt.httpClient}
			err := SetTimeout(timeout)(client)

			if err != nil {
				t.Fatalf("SetTimeout() error = %v", err)
			}

			if client.httpClient == nil {
				t.Fatal("httpClient should not be nil")
			}

			if client.httpClient.Timeout != timeout {
				t.Errorf("Timeout = %v, want %v", client.httpClient.Timeout, timeout)
			}
		})
	}
}

func TestSetHTTPClientConfig(t *testing.T) {
	config := &HTTPClientConfig{
		Timeout:     15 * time.Second,
		DialTimeout: 5 * time.Second,
	}

	tests := []struct {
		name   string
		config *HTTPClientConfig
	}{
		{
			name:   "custom config",
			config: config,
		},
		{
			name:   "nil config uses default",
			config: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{}
			err := SetHTTPClientConfig(tt.config)(client)

			if err != nil {
				t.Fatalf("SetHTTPClientConfig() error = %v", err)
			}

			if client.httpClient == nil {
				t.Fatal("httpClient should not be nil")
			}

			if tt.config != nil && client.httpClient.Timeout != tt.config.Timeout {
				t.Errorf("Timeout = %v, want %v", client.httpClient.Timeout, tt.config.Timeout)
			}
		})
	}
}

func TestSetDialTimeout(t *testing.T) {
	client := &Client{}
	timeout := 5 * time.Second

	err := SetDialTimeout(timeout)(client)

	if err != nil {
		t.Fatalf("SetDialTimeout() error = %v", err)
	}

	if client.httpClient == nil {
		t.Fatal("httpClient should not be nil")
	}
}

func TestSetTLSHandshakeTimeout(t *testing.T) {
	client := &Client{}
	timeout := 5 * time.Second

	err := SetTLSHandshakeTimeout(timeout)(client)

	if err != nil {
		t.Fatalf("SetTLSHandshakeTimeout() error = %v", err)
	}

	if client.httpClient == nil {
		t.Fatal("httpClient should not be nil")
	}
}

func TestSetLogger(t *testing.T) {
	logger := NewDefaultLogger(LogLevelInfo)
	client := &Client{}

	err := SetLogger(logger)(client)

	if err != nil {
		t.Fatalf("SetLogger() error = %v", err)
	}

	if client.logger != logger {
		t.Error("client.logger should match provided logger")
	}
}

func TestEnableLogging(t *testing.T) {
	client := &Client{}

	err := EnableLogging(LogLevelDebug)(client)

	if err != nil {
		t.Fatalf("EnableLogging() error = %v", err)
	}

	if client.logger == nil {
		t.Fatal("logger should not be nil")
	}
}

func TestSetMetrics(t *testing.T) {
	metrics := NewMetrics()
	client := &Client{}

	err := SetMetrics(metrics)(client)

	if err != nil {
		t.Fatalf("SetMetrics() error = %v", err)
	}

	if client.metrics != metrics {
		t.Error("client.metrics should match provided metrics")
	}
}

func TestEnableMetrics(t *testing.T) {
	client := &Client{}

	err := EnableMetrics()(client)

	if err != nil {
		t.Fatalf("EnableMetrics() error = %v", err)
	}

	if client.metrics == nil {
		t.Fatal("metrics should not be nil")
	}
}
