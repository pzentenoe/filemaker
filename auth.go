package filemaker

import (
	"context"
	"fmt"
	"net/http"
)

const (
	sessionAuthPath string = "fmi/data/%s/databases/%s/sessions"
)

type ConnectionDatasource struct {
	FmDataSource []FmDatasource `json:"fmDataSource"`
}
type FmDatasource struct {
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Connect establishes a session with the FileMaker database using basic authentication.
// It uses context.Background() internally. For custom context, use ConnectWithContext.
func (c *Client) Connect(database string) (*ResponseData, error) {
	return c.ConnectWithContext(context.Background(), database)
}

// ConnectWithContext establishes a session with the FileMaker database using basic authentication
// with a custom context for cancellation and timeout support.
func (c *Client) ConnectWithContext(ctx context.Context, database string) (*ResponseData, error) {
	version := c.getVersion()
	path := fmt.Sprintf(sessionAuthPath, version, database)

	options := &performRequestOptions{
		Method:      http.MethodPost,
		Path:        path,
		Body:        "{}",
		ContentType: "application/json",
		basicAuth:   true,
	}

	return c.executeQuery(ctx, options)
}

// ConnectWithDatasource establishes a session with external data source authentication.
// It uses context.Background() internally. For custom context, use ConnectWithDatasourceContext.
func (c *Client) ConnectWithDatasource(database string) (*ResponseData, error) {
	return c.ConnectWithDatasourceContext(context.Background(), database)
}

// ConnectWithDatasourceContext establishes a session with external data source authentication
// with a custom context for cancellation and timeout support.
func (c *Client) ConnectWithDatasourceContext(ctx context.Context, database string) (*ResponseData, error) {
	version, username, password := c.getCredentials()
	path := fmt.Sprintf(sessionAuthPath, version, database)
	datasource := FmDatasource{
		Database: database,
		Username: username,
		Password: password,
	}

	fileMakerConnection := ConnectionDatasource{FmDataSource: []FmDatasource{
		datasource,
	}}

	options := &performRequestOptions{
		Method:    http.MethodPost,
		Path:      path,
		Body:      fileMakerConnection,
		basicAuth: true,
	}

	return c.executeQuery(ctx, options)
}

// Disconnect logs out of a FileMaker session and invalidates the token.
// It uses context.Background() internally. For custom context, use DisconnectWithContext.
func (c *Client) Disconnect(database, token string) (*ResponseData, error) {
	return c.DisconnectWithContext(context.Background(), database, token)
}

// DisconnectWithContext logs out of a FileMaker session with a custom context
// for cancellation and timeout support.
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
// It uses context.Background() internally. For custom context, use ValidateSessionWithContext.
// Returns nil error if the session is valid, otherwise returns an error indicating the session is invalid.
func (c *Client) ValidateSession(database, token string) (*ResponseData, error) {
	return c.ValidateSessionWithContext(context.Background(), database, token)
}

// ValidateSessionWithContext checks if a session token is still valid with a custom context
// for cancellation and timeout support.
// This endpoint is useful for checking session validity before performing operations,
// or for implementing session keep-alive mechanisms.
// Returns nil error if the session is valid, otherwise returns an error.
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
