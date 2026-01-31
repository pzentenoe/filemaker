package filemaker

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Internal structs for testing to match the unexported ones in auth.go
type testConnectionDatasource struct {
	FmDataSource []testFmDatasource `json:"fmDataSource"`
}
type testFmDatasource struct {
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestConnect(t *testing.T) {
	tests := []struct {
		name           string
		database       string
		serverResponse ResponseData
		statusCode     int
		wantErr        bool
		errType        string
	}{
		{
			name:     "successful connection",
			database: "TestDB",
			serverResponse: ResponseData{
				Response: Response{Token: "test-token-123"},
				Messages: []Message{{Code: "0", Message: "OK"}},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:     "invalid credentials",
			database: "TestDB",
			serverResponse: ResponseData{
				Messages: []Message{{Code: "212", Message: "Invalid user account or password"}},
			},
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
			errType:    "FileMakerError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				username, password, ok := r.BasicAuth()
				if !ok || username == "" || password == "" {
					w.WriteHeader(http.StatusUnauthorized)
					_ = json.NewEncoder(w).Encode(ResponseData{
						Messages: []Message{{Code: "212", Message: "Invalid credentials"}},
					})
					return
				}

				w.WriteHeader(tt.statusCode)
				_ = json.NewEncoder(w).Encode(tt.serverResponse)
			}))
			defer server.Close()

			client, _ := NewClient(
				SetURL(server.URL),
				SetBasicAuth("testuser", "testpass"),
			)

			resp, err := client.CreateSession(context.Background(), tt.database)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if resp.Response.Token == "" {
				t.Error("expected non-empty token")
			}
		})
	}
}

func TestConnectWithContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := ResponseData{
			Response: Response{Token: "test-token"},
			Messages: []Message{{Code: "0", Message: "OK"}},
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, _ := NewClient(
		SetURL(server.URL),
		SetBasicAuth("testuser", "testpass"),
		DisableRetry(),
	)

	t.Run("with background context", func(t *testing.T) {
		ctx := context.Background()
		_, err := client.CreateSession(ctx, "TestDB")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("with nil context", func(t *testing.T) {
		_, err := client.CreateSession(context.TODO(), "TestDB")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("with canceled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := client.CreateSession(ctx, "TestDB")
		if err == nil {
			t.Error("expected error with canceled context")
		}
	})
}

func TestConnectWithDatasource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check Basic Auth (Master credentials)
		username, password, ok := r.BasicAuth()
		if !ok || username != "masterUser" || password != "masterPass" {
			t.Errorf("expected master credentials in Basic Auth, got %s:%s", username, password)
		}

		// Check JSON Body (External credentials)
		var payload testConnectionDatasource
		_ = json.NewDecoder(r.Body).Decode(&payload)

		if len(payload.FmDataSource) != 1 {
			t.Errorf("expected 1 datasource, got %d", len(payload.FmDataSource))
		}
		ds := payload.FmDataSource[0]
		if ds.Username != "externalUser" || ds.Password != "externalPass" {
			t.Errorf("expected external credentials in body, got %s:%s", ds.Username, ds.Password)
		}

		resp := ResponseData{
			Response: Response{Token: "datasource-token"},
			Messages: []Message{{Code: "0", Message: "OK"}},
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, _ := NewClient(
		SetURL(server.URL),
		SetBasicAuth("masterUser", "masterPass"),
	)

	resp, err := client.CreateSession(context.Background(), "TestDB", WithCustomDatasource("externalUser", "externalPass"))

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if resp.Response.Token != "datasource-token" {
		t.Errorf("expected token datasource-token, got %s", resp.Response.Token)
	}
}

func TestDisconnect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		resp := ResponseData{
			Messages: []Message{{Code: "0", Message: "OK"}},
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, _ := NewClient(SetURL(server.URL))
	_, err := client.Disconnect("TestDB", "test-token")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDisconnectWithContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := ResponseData{
			Messages: []Message{{Code: "0", Message: "OK"}},
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, _ := NewClient(SetURL(server.URL))
	ctx := context.Background()
	_, err := client.DisconnectWithContext(ctx, "TestDB", "test-token")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateSession(t *testing.T) {
	tests := []struct {
		name       string
		database   string
		token      string
		statusCode int
		wantErr    bool
		errType    string
	}{
		{
			name:       "valid session",
			database:   "TestDB",
			token:      "valid-token",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "empty database",
			database:   "",
			token:      "valid-token",
			statusCode: http.StatusOK,
			wantErr:    true,
			errType:    "ValidationError",
		},
		{
			name:       "empty token",
			database:   "TestDB",
			token:      "",
			statusCode: http.StatusOK,
			wantErr:    true,
			errType:    "ValidationError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				authHeader := r.Header.Get("Authorization")
				expectedAuth := "Bearer " + tt.token
				if tt.token != "" && authHeader != expectedAuth {
					t.Errorf("expected Authorization %s, got %s", expectedAuth, authHeader)
				}

				resp := ResponseData{
					Messages: []Message{{Code: "0", Message: "OK"}},
				}
				w.WriteHeader(tt.statusCode)
				_ = json.NewEncoder(w).Encode(resp)
			}))
			defer server.Close()

			client, _ := NewClient(SetURL(server.URL))
			_, err := client.ValidateSession(tt.database, tt.token)

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
		})
	}
}

func TestValidateSessionWithContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := ResponseData{
			Messages: []Message{{Code: "0", Message: "OK"}},
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, _ := NewClient(SetURL(server.URL), DisableRetry())

	t.Run("with context", func(t *testing.T) {
		ctx := context.Background()
		_, err := client.ValidateSessionWithContext(ctx, "TestDB", "valid-token")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("with timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := client.ValidateSessionWithContext(ctx, "TestDB", "valid-token")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestConnectWithOAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := r.Header.Get("X-FM-Data-OAuth-Request-Id")
		identifier := r.Header.Get("X-FM-Data-OAuth-Identifier")

		if requestId != "req-id-123" {
			t.Errorf("expected request id req-id-123, got %s", requestId)
		}
		if identifier != "ident-456" {
			t.Errorf("expected identifier ident-456, got %s", identifier)
		}

		if r.Header.Get("Authorization") != "" {
			t.Error("expected no Authorization header")
		}

		resp := ResponseData{
			Response: Response{Token: "oauth-token"},
			Messages: []Message{{Code: "0", Message: "OK"}},
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, _ := NewClient(SetURL(server.URL))
	resp, err := client.CreateSession(context.Background(), "TestDB", WithOAuth("req-id-123", "ident-456"))

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp.Response.Token != "oauth-token" {
		t.Errorf("expected token oauth-token, got %s", resp.Response.Token)
	}
}
