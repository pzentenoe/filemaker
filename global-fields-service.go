package filemaker

import (
	"context"
	"fmt"
	"net/http"
)

const (
	globalFieldsPath = "fmi/data/%s/databases/%s/globals"
)

// GlobalFieldsService provides methods for setting global field values in FileMaker.
// Global fields are fields that maintain their value across all records in a database session.
type GlobalFieldsService interface {
	// SetGlobalFields sets one or more global field values for the current session.
	// The globalFields parameter is a map where keys are field names and values are the values to set.
	// Requires a valid session token.
	SetGlobalFields(ctx context.Context, database string, globalFields map[string]any, token string) (*ResponseData, error)
}

type globalFieldsService struct {
	client *Client
}

// NewGlobalFieldsService creates a new GlobalFieldsService instance.
func NewGlobalFieldsService(client *Client) GlobalFieldsService {
	return &globalFieldsService{
		client: client,
	}
}

// SetGlobalFields sets one or more global field values for the current session.
// Global fields are database-wide fields that maintain their value for the duration of the session.
// This is useful for setting preferences, user context, or other session-specific data.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - database: Name of the FileMaker database
//   - globalFields: Map of field names to values (e.g., {"MyGlobalField": "value"})
//   - token: Valid session token from Connect or ConnectWithDatasource
//
// Example:
//
//	globalFields := map[string]any{
//	    "MyGlobal::UserName": "john@example.com",
//	    "MyGlobal::Theme": "dark",
//	    "MyGlobal::Language": "en",
//	}
//	response, err := service.SetGlobalFields(ctx, "MyDatabase", globalFields, token)
//
// Returns the response from the FileMaker Data API.
func (g *globalFieldsService) SetGlobalFields(ctx context.Context, database string, globalFields map[string]any, token string) (*ResponseData, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if database == "" {
		return nil, &ValidationError{
			Field:   "database",
			Message: "database name is required",
		}
	}

	if len(globalFields) == 0 {
		return nil, &ValidationError{
			Field:   "globalFields",
			Message: "at least one global field must be specified",
		}
	}

	if token == "" {
		return nil, &ValidationError{
			Field:   "token",
			Message: "session token is required",
		}
	}

	g.client.mu.RLock()
	version := g.client.version
	g.client.mu.RUnlock()

	path := fmt.Sprintf(globalFieldsPath, version, database)

	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+token)

	// Wrap global fields in the expected request format
	body := map[string]any{
		"globalFields": globalFields,
	}

	options := &performRequestOptions{
		Method:      http.MethodPatch,
		Path:        path,
		Headers:     headers,
		Body:        body,
		ContentType: "application/json",
	}

	return g.client.executeQuery(ctx, options)
}

// GlobalField represents a single global field with its name and value.
// This is a helper type for building global field maps programmatically.
type GlobalField struct {
	Name  string
	Value any
}

// NewGlobalField creates a new GlobalField instance.
func NewGlobalField(name string, value any) *GlobalField {
	return &GlobalField{
		Name:  name,
		Value: value,
	}
}

// GlobalFieldsBuilder provides a fluent interface for building global field maps.
type GlobalFieldsBuilder struct {
	fields map[string]any
}

// NewGlobalFieldsBuilder creates a new GlobalFieldsBuilder instance.
func NewGlobalFieldsBuilder() *GlobalFieldsBuilder {
	return &GlobalFieldsBuilder{
		fields: make(map[string]any),
	}
}

// Add adds a global field to the builder.
func (b *GlobalFieldsBuilder) Add(name string, value any) *GlobalFieldsBuilder {
	b.fields[name] = value
	return b
}

// AddField adds a GlobalField to the builder.
func (b *GlobalFieldsBuilder) AddField(field *GlobalField) *GlobalFieldsBuilder {
	b.fields[field.Name] = field.Value
	return b
}

// AddFields adds multiple GlobalFields to the builder.
func (b *GlobalFieldsBuilder) AddFields(fields ...*GlobalField) *GlobalFieldsBuilder {
	for _, field := range fields {
		b.fields[field.Name] = field.Value
	}
	return b
}

// Build returns the map of global fields.
func (b *GlobalFieldsBuilder) Build() map[string]any {
	return b.fields
}

// Clear removes all fields from the builder.
func (b *GlobalFieldsBuilder) Clear() *GlobalFieldsBuilder {
	b.fields = make(map[string]any)
	return b
}

// Count returns the number of fields in the builder.
func (b *GlobalFieldsBuilder) Count() int {
	return len(b.fields)
}
