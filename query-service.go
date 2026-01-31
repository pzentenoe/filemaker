package filemaker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// RecordService defines the interface for CRUD operations on FileMaker records.
type RecordService interface {
	Create(ctx context.Context, payload *Payload) (*ResponseData, error)
	Edit(ctx context.Context, recordId string, payload *Payload) (*ResponseData, error)
	Duplicate(ctx context.Context, recordId string) (*ResponseData, error)
	Delete(ctx context.Context, recordId string) (*ResponseData, error)
	GetById(ctx context.Context, recordId string) (*ResponseData, error)
	List(ctx context.Context, offset, limit string, sorters ...*Sorter) (*ResponseData, error)
}

const recordsPath = "fmi/data/%s/databases/%s/layouts/%s/records"

type recordService struct {
	database string
	layout   string
	client   *Client
}

// NewRecordService creates a new instance of RecordService for a specific database and layout.
func NewRecordService(database, layout string, client *Client) *recordService {
	return &recordService{
		database: database,
		layout:   layout,
		client:   client,
	}
}

// Payload represents the data structure for creating or editing records.
type Payload struct {
	FieldData  any `json:"fieldData"`
	PortalData any `json:"portalData,omitempty"`
}

// withAuth wraps a function with authentication session management.
// It handles context initialization, session creation, and cleanup.
func (s *recordService) withAuth(ctx context.Context, fn func(context.Context, string) (*ResponseData, error)) (*ResponseData, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	auth, err := s.client.ConnectWithContext(ctx, s.database)
	if err != nil {
		return nil, err
	}
	defer func(client *Client, ctx context.Context, database, token string) {
		_, _ = client.DisconnectWithContext(ctx, database, token)
	}(s.client, ctx, s.database, auth.Response.Token)

	return fn(ctx, auth.Response.Token)
}

// Create creates a new record in the FileMaker database.
func (s *recordService) Create(ctx context.Context, payload *Payload) (*ResponseData, error) {
	return s.withAuth(ctx, func(ctx context.Context, token string) (*ResponseData, error) {
		path := fmt.Sprintf(recordsPath, s.client.getVersion(), s.database, s.layout)
		options := &performRequestOptions{
			Method:      http.MethodPost,
			Path:        path,
			ContentType: "application/json",
			Body:        payload,
			Headers: http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer %s", token)},
			},
		}
		return s.client.executeQuery(ctx, options)
	})
}

// Edit updates an existing record in the FileMaker database.
func (s *recordService) Edit(ctx context.Context, recordId string, payload *Payload) (*ResponseData, error) {
	return s.withAuth(ctx, func(ctx context.Context, token string) (*ResponseData, error) {
		path := fmt.Sprintf(recordsPath+"/%s", s.client.getVersion(), s.database, s.layout, recordId)
		options := &performRequestOptions{
			Method:      http.MethodPatch,
			Path:        path,
			ContentType: "application/json",
			Body:        payload,
			Headers: http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer %s", token)},
			},
		}
		return s.client.executeQuery(ctx, options)
	})
}

// Duplicate creates a copy of an existing record.
func (s *recordService) Duplicate(ctx context.Context, recordId string) (*ResponseData, error) {
	return s.withAuth(ctx, func(ctx context.Context, token string) (*ResponseData, error) {
		path := fmt.Sprintf(recordsPath+"/%s", s.client.getVersion(), s.database, s.layout, recordId)
		options := &performRequestOptions{
			Method:      http.MethodPost,
			Path:        path,
			ContentType: "application/json",
			Headers: http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer %s", token)},
			},
		}
		return s.client.executeQuery(ctx, options)
	})
}

// Delete removes a record from the FileMaker database.
func (s *recordService) Delete(ctx context.Context, recordId string) (*ResponseData, error) {
	return s.withAuth(ctx, func(ctx context.Context, token string) (*ResponseData, error) {
		path := fmt.Sprintf(recordsPath+"/%s", s.client.getVersion(), s.database, s.layout, recordId)
		options := &performRequestOptions{
			Method: http.MethodDelete,
			Path:   path,
			Headers: http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer %s", token)},
			},
		}
		return s.client.executeQuery(ctx, options)
	})
}

// GetById retrieves a single record by its ID.
func (s *recordService) GetById(ctx context.Context, recordId string) (*ResponseData, error) {
	return s.withAuth(ctx, func(ctx context.Context, token string) (*ResponseData, error) {
		path := fmt.Sprintf(recordsPath+"/%s", s.client.getVersion(), s.database, s.layout, recordId)
		options := &performRequestOptions{
			Method: http.MethodGet,
			Path:   path,
			Headers: http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer %s", token)},
			},
		}
		return s.client.executeQuery(ctx, options)
	})
}

// List retrieves a range of records with optional sorting.
// offset and limit control pagination, sorters define the sort order.
func (s *recordService) List(ctx context.Context, offset, limit string, sorters ...*Sorter) (*ResponseData, error) {
	return s.withAuth(ctx, func(ctx context.Context, token string) (*ResponseData, error) {
		path := fmt.Sprintf(recordsPath, s.client.getVersion(), s.database, s.layout)

		params := url.Values{}
		params.Add("_offset", offset)
		params.Add("_limit", limit)

		sortersStr := sortersToJson(sorters...)
		if sortersStr != "" {
			params.Add("_sort", sortersStr)
		}

		options := &performRequestOptions{
			Method: http.MethodGet,
			Path:   path,
			Params: params,
			Headers: http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer %s", token)},
			},
		}
		return s.client.executeQuery(ctx, options)
	})
}

func sortersToJson(sorters ...*Sorter) string {
	if len(sorters) > 0 {
		data, err := json.Marshal(sorters)
		if err != nil {
			return ""
		}
		return string(data)
	}
	return ""
}
