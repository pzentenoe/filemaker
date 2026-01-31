package filemaker

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewMetadataService(t *testing.T) {
	client, _ := NewClient()
	service := NewMetadataService(client)

	if service == nil {
		t.Fatal("service is nil")
	}
}

func TestMetadataService_GetDatabases(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			t.Error("expected basic auth")
		}
		if username == "" || password == "" {
			t.Error("expected non-empty credentials")
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		resp := ResponseData{
			Response: Response{
				Data: []Datum{
					{RecordID: "1"},
					{RecordID: "2"},
				},
			},
			Messages: []Message{{Code: "0", Message: "OK"}},
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, _ := NewClient(
		SetURL(server.URL),
		SetBasicAuth("user", "pass"),
	)
	service := NewMetadataService(client)

	resp, err := service.GetDatabases(context.Background())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("response is nil")
	}
	if len(resp.Response.Data) != 2 {
		t.Errorf("expected 2 databases, got %d", len(resp.Response.Data))
	}
}

func TestMetadataService_GetDatabases_WithContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := ResponseData{
			Response: Response{Data: []Datum{{RecordID: "1"}}},
			Messages: []Message{{Code: "0", Message: "OK"}},
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, _ := NewClient(
		SetURL(server.URL),
		SetBasicAuth("user", "pass"),
	)
	service := NewMetadataService(client)

	t.Run("with background context", func(t *testing.T) {
		ctx := context.Background()
		_, err := service.GetDatabases(ctx)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("with nil context", func(t *testing.T) {
		_, err := service.GetDatabases(context.TODO())
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestMetadataService_GetLayouts(t *testing.T) {
	tests := []struct {
		name     string
		database string
		token    string
		wantErr  bool
		errType  string
	}{
		{
			name:     "successful get layouts",
			database: "TestDB",
			token:    "test-token",
			wantErr:  false,
		},
		{
			name:     "empty database",
			database: "",
			token:    "test-token",
			wantErr:  true,
			errType:  "ValidationError",
		},
		{
			name:     "empty token",
			database: "TestDB",
			token:    "",
			wantErr:  true,
			errType:  "ValidationError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				authHeader := r.Header.Get("Authorization")
				if authHeader != "Bearer "+tt.token {
					t.Errorf("expected Bearer %s, got %s", tt.token, authHeader)
				}

				resp := ResponseData{
					Response: Response{
						Data: []Datum{
							{RecordID: "Layout1"},
							{RecordID: "Layout2"},
						},
					},
					Messages: []Message{{Code: "0", Message: "OK"}},
				}
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(resp)
			}))
			defer server.Close()

			client, _ := NewClient(SetURL(server.URL))
			service := NewMetadataService(client)

			resp, err := service.GetLayouts(context.Background(), tt.database, tt.token)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if tt.errType == "ValidationError" {
					if _, ok := err.(*ValidationError); !ok {
						t.Errorf("expected ValidationError, got %T", err)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if resp == nil {
				t.Fatal("response is nil")
			}
		})
	}
}

func TestMetadataService_GetLayoutMetadata(t *testing.T) {
	tests := []struct {
		name     string
		database string
		layout   string
		token    string
		wantErr  bool
		errType  string
	}{
		{
			name:     "successful get layout metadata",
			database: "TestDB",
			layout:   "TestLayout",
			token:    "test-token",
			wantErr:  false,
		},
		{
			name:     "empty database",
			database: "",
			layout:   "TestLayout",
			token:    "test-token",
			wantErr:  true,
			errType:  "ValidationError",
		},
		{
			name:     "empty layout",
			database: "TestDB",
			layout:   "",
			token:    "test-token",
			wantErr:  true,
			errType:  "ValidationError",
		},
		{
			name:     "empty token",
			database: "TestDB",
			layout:   "TestLayout",
			token:    "",
			wantErr:  true,
			errType:  "ValidationError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				authHeader := r.Header.Get("Authorization")
				if authHeader != "Bearer "+tt.token {
					t.Errorf("expected Bearer %s, got %s", tt.token, authHeader)
				}

				resp := ResponseData{
					Response: Response{RecordID: "metadata"},
					Messages: []Message{{Code: "0", Message: "OK"}},
				}
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(resp)
			}))
			defer server.Close()

			client, _ := NewClient(SetURL(server.URL))
			service := NewMetadataService(client)

			resp, err := service.GetLayoutMetadata(context.Background(), tt.database, tt.layout, tt.token)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if tt.errType == "ValidationError" {
					if _, ok := err.(*ValidationError); !ok {
						t.Errorf("expected ValidationError, got %T", err)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if resp == nil {
				t.Fatal("response is nil")
			}
		})
	}
}

func TestMetadataService_GetScripts(t *testing.T) {
	tests := []struct {
		name     string
		database string
		token    string
		wantErr  bool
		errType  string
	}{
		{
			name:     "successful get scripts",
			database: "TestDB",
			token:    "test-token",
			wantErr:  false,
		},
		{
			name:     "empty database",
			database: "",
			token:    "test-token",
			wantErr:  true,
			errType:  "ValidationError",
		},
		{
			name:     "empty token",
			database: "TestDB",
			token:    "",
			wantErr:  true,
			errType:  "ValidationError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				authHeader := r.Header.Get("Authorization")
				if authHeader != "Bearer "+tt.token {
					t.Errorf("expected Bearer %s, got %s", tt.token, authHeader)
				}

				resp := ResponseData{
					Response: Response{
						Data: []Datum{
							{RecordID: "Script1"},
							{RecordID: "Script2"},
						},
					},
					Messages: []Message{{Code: "0", Message: "OK"}},
				}
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(resp)
			}))
			defer server.Close()

			client, _ := NewClient(SetURL(server.URL))
			service := NewMetadataService(client)

			resp, err := service.GetScripts(context.Background(), tt.database, tt.token)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if tt.errType == "ValidationError" {
					if _, ok := err.(*ValidationError); !ok {
						t.Errorf("expected ValidationError, got %T", err)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if resp == nil {
				t.Fatal("response is nil")
			}
		})
	}
}

func TestMetadataService_GetProductInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			t.Error("expected basic auth")
		}
		if username == "" || password == "" {
			t.Error("expected non-empty credentials")
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		resp := ResponseData{
			Response: Response{RecordID: "product-info"},
			Messages: []Message{{Code: "0", Message: "OK"}},
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, _ := NewClient(
		SetURL(server.URL),
		SetBasicAuth("user", "pass"),
	)
	service := NewMetadataService(client)

	resp, err := service.GetProductInfo(context.Background())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("response is nil")
	}
}

func TestMetadataService_GetProductInfo_WithContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := ResponseData{
			Response: Response{RecordID: "info"},
			Messages: []Message{{Code: "0", Message: "OK"}},
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, _ := NewClient(
		SetURL(server.URL),
		SetBasicAuth("user", "pass"),
	)
	service := NewMetadataService(client)

	t.Run("with background context", func(t *testing.T) {
		ctx := context.Background()
		_, err := service.GetProductInfo(ctx)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("with nil context", func(t *testing.T) {
		_, err := service.GetProductInfo(context.TODO())
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
