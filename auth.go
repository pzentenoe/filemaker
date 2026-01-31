package filemaker

import (
	"context"
	"fmt"
	"net/http"
)

const (
	sessionAuthPath string = "fmi/data/%s/databases/%s/sessions"
	contentTypeJSON        = "application/json"
)

// AuthProvider defines the strategy for authentication.
// It configures the request options (headers, body, auth type) for the session creation.
type AuthProvider func(client *Client, database string) (*performRequestOptions, error)

// WithBasicAuth uses the provided username and password for authentication.
func WithBasicAuth(username, password string) AuthProvider {

	return func(c *Client, database string) (*performRequestOptions, error) {
		version := c.getVersion()
		path := fmt.Sprintf(sessionAuthPath, version, database)

		return &performRequestOptions{
				Method:      http.MethodPost,
				Path:        path,
				Body:        "{}",
				ContentType: "application/json",
				Username:    username,
				Password:    password,
				BasicAuth:   true,
			},
			nil
	}

}

// WithCustomDatasource uses provided credentials to authenticate against an external data source.
// The provided username and password are used for the external data source (JSON body).
// The client's configured credentials (SetBasicAuth) are used for the FileMaker file authentication (HTTP Header).
func WithCustomDatasource(externalUsername, externalPassword string) AuthProvider {
	return func(c *Client, database string) (*performRequestOptions, error) {
		version := c.getVersion()
		path := fmt.Sprintf(sessionAuthPath, version, database)

		// Use passed credentials for the external data source payload
		datasource := fmDatasource{
			Database: database,
			Username: externalUsername,
			Password: externalPassword,
		}

		payload := connectionDatasource{FmDataSource: []fmDatasource{datasource}}

		// Use client's stored "master" credentials for the HTTP Basic Auth header
		c.mu.RLock()
		defer c.mu.RUnlock()
		masterUser := c.username
		masterPass := c.password

		return &performRequestOptions{
				Method:    http.MethodPost,
				Path:      path,
				Body:      payload,
				Username:  masterUser, // HTTP Basic Auth User (Master File)
				Password:  masterPass, // HTTP Basic Auth Pass (Master File)
				BasicAuth: true,
			},
			nil
	}
}

// WithOAuth authenticates using OAuth headers.
func WithOAuth(requestId, identifier string) AuthProvider {
	return func(c *Client, database string) (*performRequestOptions, error) {
		version := c.getVersion()
		path := fmt.Sprintf(sessionAuthPath, version, database)

		headers := http.Header{}
		headers.Set("X-FM-Data-OAuth-Request-Id", requestId)
		headers.Set("X-FM-Data-OAuth-Identifier", identifier)

		return &performRequestOptions{
			Method:      http.MethodPost,
			Path:        path,
			Body:        "{}",
			ContentType: contentTypeJSON,
			Headers:     headers,
			BasicAuth:   false,
		}, nil
	}
}

// WithFMID authenticates using a Claris ID token.
func WithFMID(fmidToken string) AuthProvider {
	return func(c *Client, database string) (*performRequestOptions, error) {
		version := c.getVersion()
		path := fmt.Sprintf(sessionAuthPath, version, database)

		headers := http.Header{}
		headers.Set("Authorization", "FMID "+fmidToken)

		return &performRequestOptions{
			Method:      http.MethodPost,
			Path:        path,
			Body:        "{}",
			ContentType: contentTypeJSON,
			Headers:     headers,
			BasicAuth:   false,
		}, nil
	}
}

// CreateSession establishes a new session using the provided authentication provider.
// If no auth provider is passed, it uses the client's default authentication provider (configured via SetBasicAuth, etc.).
func (c *Client) CreateSession(ctx context.Context, database string, auth ...AuthProvider) (*ResponseData, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if len(auth) > 1 {
		return nil, &ValidationError{
			Field:   "auth",
			Message: "only one authentication provider allowed",
		}
	}

	var provider AuthProvider
	if len(auth) > 0 {
		provider = auth[0]
	}

	// Use default if no specific provider is passed
	if provider == nil {
		provider = c.authProvider
	}

	if provider == nil {
		return nil, &ValidationError{
			Field:   "auth",
			Message: "authentication provider is required (configure client with SetBasicAuth or similar)",
		}
	}

	options, err := provider(c, database)
	if err != nil {
		return nil, err
	}

	return c.executeQuery(ctx, options)
}

type connectionDatasource struct {
	FmDataSource []fmDatasource `json:"fmDataSource"`
}

type fmDatasource struct {
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Disconnect logs out of a FileMaker session.
func (c *Client) Disconnect(database, token string) (*ResponseData, error) {
	return c.DisconnectWithContext(context.Background(), database, token)
}

// DisconnectWithContext logs out of a FileMaker session with context.
func (c *Client) DisconnectWithContext(ctx context.Context, database, token string) (*ResponseData, error) {
	version := c.getVersion()
	path := fmt.Sprintf(sessionAuthPath+"/%s", version, database, token)

	options := &performRequestOptions{
		Method: http.MethodDelete,
		Path:   path,
	}

	return c.executeQuery(ctx, options)
}

// ValidateSession checks if a session token is still valid.
func (c *Client) ValidateSession(database, token string) (*ResponseData, error) {
	return c.ValidateSessionWithContext(context.Background(), database, token)
}

// ValidateSessionWithContext checks if a session token is still valid with context.
func (c *Client) ValidateSessionWithContext(ctx context.Context, database, token string) (*ResponseData, error) {
	if database == "" {
		return nil, &ValidationError{
			Field:   "database",
			Message: "database name is required",
		}
	}

	if token == "" {
		return nil, &ValidationError{
			Field:   "token",
			Message: "session token is required",
		}
	}

	version := c.getVersion()
	path := fmt.Sprintf(sessionAuthPath+"/%s", version, database, token)

	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+token)

	options := &performRequestOptions{
		Method:  http.MethodGet,
		Path:    path,
		Headers: headers,
	}

	return c.executeQuery(ctx, options)
}
