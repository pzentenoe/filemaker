package filemaker

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		options []ClientOptions
		check   func(*testing.T, *Client)
	}{
		{
			name:    "client with default values",
			options: []ClientOptions{},
			check: func(t *testing.T, c *Client) {
				if c.version != DefaultVersion {
					t.Errorf("expected version %s, got %s", DefaultVersion, c.version)
				}
				if c.httpClient == nil {
					t.Error("expected non-nil httpClient")
				}
				if c.retryConfig == nil {
					t.Error("expected non-nil retryConfig")
				}
			},
		},
		{
			name: "client with custom options",
			options: []ClientOptions{
				SetURL("http://localhost"),
				SetUsername("user"),
				SetPassword("pass"),
				SetVersion("v2"),
			},
			check: func(t *testing.T, c *Client) {
				if c.url != "http://localhost" {
					t.Errorf("expected url http://localhost, got %s", c.url)
				}
				if c.username != "user" {
					t.Errorf("expected username user, got %s", c.username)
				}
			},
		},
		{
			name: "client with metrics enabled",
			options: []ClientOptions{
				EnableMetrics(),
			},
			check: func(t *testing.T, c *Client) {
				if c.metrics == nil {
					t.Error("expected non-nil metrics")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.options...)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if client == nil {
				t.Fatal("client is nil")
			}
			if tt.check != nil {
				tt.check(t, client)
			}
		})
	}
}

func TestClient_getVersion(t *testing.T) {
	client, _ := NewClient(SetVersion("v2"))
	version := client.getVersion()
	if version != "v2" {
		t.Errorf("expected version v2, got %s", version)
	}
}

func TestClient_getCredentials(t *testing.T) {
	client, _ := NewClient(
		SetVersion("v2"),
		SetUsername("testuser"),
		SetPassword("testpass"),
	)

	version, username, password := client.getCredentials()
	if version != "v2" {
		t.Errorf("expected version v2, got %s", version)
	}
	if username != "testuser" {
		t.Errorf("expected username testuser, got %s", username)
	}
	if password != "testpass" {
		t.Errorf("expected password testpass, got %s", password)
	}
}

func TestClient_getURL(t *testing.T) {
	client, _ := NewClient(SetURL("http://localhost:8080"))
	url := client.getURL()
	if url != "http://localhost:8080" {
		t.Errorf("expected url http://localhost:8080, got %s", url)
	}
}

func TestClient_performRequest(t *testing.T) {
	tests := []struct {
		name      string
		setupOpts *performRequestOptions
		setupFunc func() *httptest.Server
		wantErr   bool
	}{
		{
			name: "successful GET request",
			setupOpts: &performRequestOptions{
				Method: http.MethodGet,
				Path:   "test/path",
			},
			setupFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
				}))
			},
			wantErr: false,
		},
		{
			name: "POST request with body",
			setupOpts: &performRequestOptions{
				Method:      http.MethodPost,
				Path:        "test/path",
				Body:        map[string]string{"key": "value"},
				ContentType: "application/json",
			},
			setupFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
			},
			wantErr: false,
		},
		{
			name: "request with basic auth",
			setupOpts: &performRequestOptions{
				Method:    http.MethodGet,
				Path:      "test/path",
				basicAuth: true,
			},
			setupFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					username, password, ok := r.BasicAuth()
					if !ok || username != "testuser" || password != "testpass" {
						t.Error("basic auth failed")
					}
					w.WriteHeader(http.StatusOK)
				}))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.setupFunc()
			defer server.Close()

			client, _ := NewClient(
				SetURL(server.URL),
				SetUsername("testuser"),
				SetPassword("testpass"),
			)

			ctx := context.Background()
			resp, err := client.performRequest(ctx, tt.setupOpts)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
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

func TestClient_performRequest_Errors(t *testing.T) {
	client, _ := NewClient()

	t.Run("missing url", func(t *testing.T) {
		opts := &performRequestOptions{
			Method: http.MethodGet,
			Path:   "test",
		}
		_, err := client.performRequest(context.Background(), opts)
		if err == nil {
			t.Error("expected error for missing url")
		}
		if _, ok := err.(*ValidationError); !ok {
			t.Errorf("expected ValidationError, got %T", err)
		}
	})

	t.Run("missing method", func(t *testing.T) {
		client, _ := NewClient(SetURL("http://localhost"))
		opts := &performRequestOptions{
			Method: "",
			Path:   "test",
		}
		_, err := client.performRequest(context.Background(), opts)
		if err == nil {
			t.Error("expected error for missing method")
		}
	})
}

func TestClient_NewRequest(t *testing.T) {
	client, _ := NewClient(SetVersion("v2"))
	req, err := client.NewRequest(http.MethodGet, "http://example.com/test")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if req == nil {
		t.Fatal("request is nil")
	}

	httpReq := (*http.Request)(req)
	if httpReq.Method != http.MethodGet {
		t.Errorf("expected method GET, got %s", httpReq.Method)
	}

	userAgent := httpReq.Header.Get("User-Agent")
	if userAgent == "" {
		t.Error("expected User-Agent header")
	}

	accept := httpReq.Header.Get("Accept")
	if accept != "application/json" {
		t.Errorf("expected Accept application/json, got %s", accept)
	}
}

func TestClient_executeQuery(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() *httptest.Server
		options   *performRequestOptions
		wantErr   bool
	}{
		{
			name: "successful query execution",
			setupFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					resp := ResponseData{
						Response: Response{
							RecordID: "123",
							ModID:    "1",
						},
						Messages: []Message{{Code: "0", Message: "OK"}},
					}
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(resp)
				}))
			},
			options: &performRequestOptions{
				Method: http.MethodGet,
				Path:   "test",
			},
			wantErr: false,
		},
		{
			name: "query with FileMaker error",
			setupFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					resp := ResponseData{
						Messages: []Message{{Code: "401", Message: "Record not found"}},
					}
					w.WriteHeader(http.StatusNotFound)
					_ = json.NewEncoder(w).Encode(resp)
				}))
			},
			options: &performRequestOptions{
				Method: http.MethodGet,
				Path:   "test",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.setupFunc()
			defer server.Close()

			client, _ := NewClient(SetURL(server.URL))
			ctx := context.Background()
			_, err := client.executeQuery(ctx, tt.options)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestRequest_setBody(t *testing.T) {
	tests := []struct {
		name         string
		body         any
		gzipCompress bool
		wantErr      bool
	}{
		{
			name:         "set string body",
			body:         "test string",
			gzipCompress: false,
			wantErr:      false,
		},
		{
			name:         "set JSON body",
			body:         map[string]string{"key": "value"},
			gzipCompress: false,
			wantErr:      false,
		},
		{
			name:         "set body with gzip",
			body:         map[string]string{"key": "value"},
			gzipCompress: true,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{}
			req.Header = http.Header{}

			err := req.setBody(tt.body, tt.gzipCompress)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if req.Body == nil {
				t.Error("expected non-nil body")
			}
		})
	}
}

func TestRequest_setBasicAuth(t *testing.T) {
	req := &Request{}
	httpReq := (*http.Request)(req)
	httpReq.Header = http.Header{}

	req.setBasicAuth("testuser", "testpass")

	auth := httpReq.Header.Get("Authorization")
	if auth == "" {
		t.Error("expected Authorization header")
	}
}
