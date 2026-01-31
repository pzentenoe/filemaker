package filemaker

import (
	"errors"
	"net/http"
	"time"
)

const (
	DefaultVersion = "vLatest"
)

type ClientOptions func(*Client) error

func SetURL(url string) ClientOptions {
	return func(c *Client) error {
		if url == "" {
			return errors.New("empty url")
		}
		c.url = url
		return nil
	}
}

func SetVersion(version string) ClientOptions {
	return func(c *Client) error {
		c.version = version
		if c.version == "" {
			c.version = DefaultVersion
		}
		return nil
	}
}

// SetBasicAuth sets the default authentication to Basic Auth with the provided credentials.
func SetBasicAuth(username, password string) ClientOptions {
	return func(c *Client) error {
		if username == "" || password == "" {
			return errors.New("username and password are required")
		}
		c.username = username
		c.password = password
		c.authProvider = WithBasicAuth(username, password)
		return nil
	}
}

// SetAuthProvider sets a custom authentication provider.
func SetAuthProvider(provider AuthProvider) ClientOptions {
	return func(c *Client) error {
		if provider == nil {
			return errors.New("auth provider cannot be nil")
		}
		c.authProvider = provider
		return nil
	}
}

func SetHttpClient(httpClient *http.Client) ClientOptions {
	return func(c *Client) error {
		if httpClient != nil {
			c.httpClient = httpClient
		} else {
			c.httpClient = http.DefaultClient
		}
		return nil
	}
}

// SetRetryConfig sets a custom retry configuration for the client.
func SetRetryConfig(config *RetryConfig) ClientOptions {
	return func(c *Client) error {
		if config == nil {
			c.retryConfig = DefaultRetryConfig()
		} else {
			c.retryConfig = config
		}
		return nil
	}
}

// SetMaxRetries sets the maximum number of retry attempts.
func SetMaxRetries(max int) ClientOptions {
	return func(c *Client) error {
		if c.retryConfig == nil {
			c.retryConfig = DefaultRetryConfig()
		}
		c.retryConfig.MaxRetries = max
		return nil
	}
}

// SetRetryWaitTime sets the minimum and maximum wait time between retries.
func SetRetryWaitTime(min, max time.Duration) ClientOptions {
	return func(c *Client) error {
		if c.retryConfig == nil {
			c.retryConfig = DefaultRetryConfig()
		}
		c.retryConfig.MinWaitTime = min
		c.retryConfig.MaxWaitTime = max
		return nil
	}
}

// DisableRetry disables retry logic by setting MaxRetries to 0.
func DisableRetry() ClientOptions {
	return func(c *Client) error {
		if c.retryConfig == nil {
			c.retryConfig = DefaultRetryConfig()
		}
		c.retryConfig.MaxRetries = 0
		return nil
	}
}

// SetTimeout sets the overall request timeout.
func SetTimeout(timeout time.Duration) ClientOptions {
	return func(c *Client) error {
		if c.httpClient == nil {
			config := DefaultHTTPClientConfig()
			config.Timeout = timeout
			c.httpClient = config.BuildHTTPClient()
		} else {
			c.httpClient.Timeout = timeout
		}
		return nil
	}
}

// SetHTTPClientConfig sets a custom HTTP client configuration.
// This builds a new http.Client with the specified configuration.
func SetHTTPClientConfig(config *HTTPClientConfig) ClientOptions {
	return func(c *Client) error {
		if config == nil {
			config = DefaultHTTPClientConfig()
		}
		c.httpClient = config.BuildHTTPClient()
		return nil
	}
}

// SetDialTimeout sets the TCP connection dial timeout.
// Note: This creates a new HTTP client, discarding any previous client configuration.
// Use SetHTTPClientConfig for more control over multiple transport settings.
func SetDialTimeout(timeout time.Duration) ClientOptions {
	return func(c *Client) error {
		config := DefaultHTTPClientConfig()
		if c.httpClient != nil && c.httpClient.Timeout > 0 {
			config.Timeout = c.httpClient.Timeout
		}
		config.DialTimeout = timeout
		c.httpClient = config.BuildHTTPClient()
		return nil
	}
}

// SetTLSHandshakeTimeout sets the TLS handshake timeout.
// Note: This creates a new HTTP client, discarding any previous client configuration.
// Use SetHTTPClientConfig for more control over multiple transport settings.
func SetTLSHandshakeTimeout(timeout time.Duration) ClientOptions {
	return func(c *Client) error {
		config := DefaultHTTPClientConfig()
		if c.httpClient != nil && c.httpClient.Timeout > 0 {
			config.Timeout = c.httpClient.Timeout
		}
		config.TLSHandshakeTimeout = timeout
		c.httpClient = config.BuildHTTPClient()
		return nil
	}
}

// SetLogger sets a custom logger for the client.
func SetLogger(logger Logger) ClientOptions {
	return func(c *Client) error {
		c.logger = logger
		return nil
	}
}

// EnableLogging enables logging with the specified level.
// Uses the default slog logger with text output.
func EnableLogging(level LogLevel) ClientOptions {
	return func(c *Client) error {
		c.logger = NewDefaultLogger(level)
		return nil
	}
}

// SetMetrics sets a custom Metrics instance for the client.
func SetMetrics(metrics *Metrics) ClientOptions {
	return func(c *Client) error {
		c.metrics = metrics
		return nil
	}
}

// EnableMetrics enables metrics collection.
func EnableMetrics() ClientOptions {
	return func(c *Client) error {
		c.metrics = NewMetrics()
		return nil
	}
}
