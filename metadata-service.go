package filemaker

import (
	"context"
	"fmt"
	"net/http"
)

const (
	databasesPath      = "fmi/data/%s/databases"
	layoutsPath        = "fmi/data/%s/databases/%s/layouts"
	layoutMetadataPath = "fmi/data/%s/databases/%s/layouts/%s"
	scriptsPath        = "fmi/data/%s/databases/%s/scripts"
	productInfoPath    = "fmi/data/%s/productInfo"
)

// MetadataService provides methods for discovering and retrieving metadata
// about FileMaker databases, layouts, scripts, and product information.
type MetadataService interface {
	// GetDatabases retrieves a list of all hosted databases accessible by the authenticated user.
	GetDatabases(ctx context.Context) (*ResponseData, error)

	// GetLayouts retrieves a list of all layouts in the specified database.
	GetLayouts(ctx context.Context, database, token string) (*ResponseData, error)

	// GetLayoutMetadata retrieves detailed metadata for a specific layout including
	// field definitions, value lists, and portal information.
	GetLayoutMetadata(ctx context.Context, database, layout, token string) (*ResponseData, error)

	// GetScripts retrieves a list of all scripts in the specified database.
	GetScripts(ctx context.Context, database, token string) (*ResponseData, error)

	// GetProductInfo retrieves FileMaker Server product information including
	// version, name, and date/time formats.
	GetProductInfo(ctx context.Context) (*ResponseData, error)
}

type metadataService struct {
	client *Client
}

// NewMetadataService creates a new MetadataService instance.
func NewMetadataService(client *Client) MetadataService {
	return &metadataService{
		client: client,
	}
}

// GetDatabases retrieves a list of all hosted databases accessible by the authenticated user.
// This endpoint does not require a session token.
func (m *metadataService) GetDatabases(ctx context.Context) (*ResponseData, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	m.client.mu.RLock()
	version := m.client.version
	username := m.client.username
	password := m.client.password
	m.client.mu.RUnlock()

	path := fmt.Sprintf(databasesPath, version)

	options := &performRequestOptions{
		Method:    http.MethodGet,
		Path:      path,
		BasicAuth: true,
		Username:  username,
		Password:  password,
	}

	return m.client.executeQuery(ctx, options)
}

// GetLayouts retrieves a list of all layouts in the specified database.
// Requires a valid session token obtained from Connect or ConnectWithDatasource.
func (m *metadataService) GetLayouts(ctx context.Context, database, token string) (*ResponseData, error) {
	if ctx == nil {
		ctx = context.Background()
	}

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

	m.client.mu.RLock()
	version := m.client.version
	m.client.mu.RUnlock()

	path := fmt.Sprintf(layoutsPath, version, database)

	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+token)

	options := &performRequestOptions{
		Method:  http.MethodGet,
		Path:    path,
		Headers: headers,
	}

	return m.client.executeQuery(ctx, options)
}

// GetLayoutMetadata retrieves detailed metadata for a specific layout.
// The response includes field definitions, value lists, portal information,
// and relationships.
func (m *metadataService) GetLayoutMetadata(ctx context.Context, database, layout, token string) (*ResponseData, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if database == "" {
		return nil, &ValidationError{
			Field:   "database",
			Message: "database name is required",
		}
	}

	if layout == "" {
		return nil, &ValidationError{
			Field:   "layout",
			Message: "layout name is required",
		}
	}

	if token == "" {
		return nil, &ValidationError{
			Field:   "token",
			Message: "session token is required",
		}
	}

	m.client.mu.RLock()
	version := m.client.version
	m.client.mu.RUnlock()

	path := fmt.Sprintf(layoutMetadataPath, version, database, layout)

	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+token)

	options := &performRequestOptions{
		Method:  http.MethodGet,
		Path:    path,
		Headers: headers,
	}

	return m.client.executeQuery(ctx, options)
}

// GetScripts retrieves a list of all scripts in the specified database.
// Only scripts that are configured to be accessible from the Data API are returned.
func (m *metadataService) GetScripts(ctx context.Context, database, token string) (*ResponseData, error) {
	if ctx == nil {
		ctx = context.Background()
	}

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

	m.client.mu.RLock()
	version := m.client.version
	m.client.mu.RUnlock()

	path := fmt.Sprintf(scriptsPath, version, database)

	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+token)

	options := &performRequestOptions{
		Method:  http.MethodGet,
		Path:    path,
		Headers: headers,
	}

	return m.client.executeQuery(ctx, options)
}

// GetProductInfo retrieves FileMaker Server product information.
// This includes version, build number, server name, date/time formats, and locale.
// This endpoint does not require a session token.
func (m *metadataService) GetProductInfo(ctx context.Context) (*ResponseData, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	m.client.mu.RLock()
	version := m.client.version
	username := m.client.username
	password := m.client.password
	m.client.mu.RUnlock()

	path := fmt.Sprintf(productInfoPath, version)

	options := &performRequestOptions{
		Method:    http.MethodGet,
		Path:      path,
		BasicAuth: true,
		Username:  username,
		Password:  password,
	}

	return m.client.executeQuery(ctx, options)
}
