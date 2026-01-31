package filemaker

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewContainerService(t *testing.T) {
	client, _ := NewClient()
	service := NewContainerService(client)

	if service == nil {
		t.Fatal("service is nil")
	}
}

func TestContainerService_UploadData(t *testing.T) {
	tests := []struct {
		name      string
		database  string
		layout    string
		recordID  string
		fieldName string
		filename  string
		data      []byte
		token     string
		wantErr   bool
		errType   string
	}{
		{
			name:      "successful upload",
			database:  "TestDB",
			layout:    "TestLayout",
			recordID:  "123",
			fieldName: "Photo",
			filename:  "test.jpg",
			data:      []byte("test data"),
			token:     "test-token",
			wantErr:   false,
		},
		{
			name:      "empty database",
			database:  "",
			layout:    "TestLayout",
			recordID:  "123",
			fieldName: "Photo",
			filename:  "test.jpg",
			data:      []byte("test data"),
			token:     "test-token",
			wantErr:   true,
			errType:   "ValidationError",
		},
		{
			name:      "empty layout",
			database:  "TestDB",
			layout:    "",
			recordID:  "123",
			fieldName: "Photo",
			filename:  "test.jpg",
			data:      []byte("test data"),
			token:     "test-token",
			wantErr:   true,
			errType:   "ValidationError",
		},
		{
			name:      "empty recordID",
			database:  "TestDB",
			layout:    "TestLayout",
			recordID:  "",
			fieldName: "Photo",
			filename:  "test.jpg",
			data:      []byte("test data"),
			token:     "test-token",
			wantErr:   true,
			errType:   "ValidationError",
		},
		{
			name:      "empty fieldName",
			database:  "TestDB",
			layout:    "TestLayout",
			recordID:  "123",
			fieldName: "",
			filename:  "test.jpg",
			data:      []byte("test data"),
			token:     "test-token",
			wantErr:   true,
			errType:   "ValidationError",
		},
		{
			name:      "empty filename",
			database:  "TestDB",
			layout:    "TestLayout",
			recordID:  "123",
			fieldName: "Photo",
			filename:  "",
			data:      []byte("test data"),
			token:     "test-token",
			wantErr:   true,
			errType:   "ValidationError",
		},
		{
			name:      "empty data",
			database:  "TestDB",
			layout:    "TestLayout",
			recordID:  "123",
			fieldName: "Photo",
			filename:  "test.jpg",
			data:      []byte{},
			token:     "test-token",
			wantErr:   true,
			errType:   "ValidationError",
		},
		{
			name:      "empty token",
			database:  "TestDB",
			layout:    "TestLayout",
			recordID:  "123",
			fieldName: "Photo",
			filename:  "test.jpg",
			data:      []byte("test data"),
			token:     "",
			wantErr:   true,
			errType:   "ValidationError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST, got %s", r.Method)
				}

				authHeader := r.Header.Get("Authorization")
				if tt.token != "" && authHeader != "Bearer "+tt.token {
					t.Errorf("expected Bearer %s, got %s", tt.token, authHeader)
				}

				contentType := r.Header.Get("Content-Type")
				if !strings.Contains(contentType, "multipart/form-data") {
					t.Errorf("expected multipart/form-data, got %s", contentType)
				}

				resp := ResponseData{
					Response: Response{ModID: "2"},
					Messages: []Message{{Code: "0", Message: "OK"}},
				}
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(resp)
			}))
			defer server.Close()

			client, _ := NewClient(SetURL(server.URL), DisableRetry())
			service := NewContainerService(client)

			response, err := service.UploadData(context.Background(), tt.database, tt.layout, tt.recordID, tt.fieldName, tt.filename, tt.data, tt.token)

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
			if response == nil {
				t.Fatal("response is nil")
			}
		})
	}
}

func TestContainerService_UploadDataWithRepetition(t *testing.T) {
	tests := []struct {
		name       string
		repetition int
		wantErr    bool
		errType    string
	}{
		{
			name:       "repetition 1",
			repetition: 1,
			wantErr:    false,
		},
		{
			name:       "repetition 3",
			repetition: 3,
			wantErr:    false,
		},
		{
			name:       "repetition 0",
			repetition: 0,
			wantErr:    true,
			errType:    "ValidationError",
		},
		{
			name:       "negative repetition",
			repetition: -1,
			wantErr:    true,
			errType:    "ValidationError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := "/fmi/data/vLatest/databases/TestDB/layouts/TestLayout/records/123/containers/Photo/"
				if !strings.Contains(r.URL.Path, expectedPath) {
					t.Errorf("expected path to contain %s, got %s", expectedPath, r.URL.Path)
				}

				resp := ResponseData{
					Response: Response{ModID: "2"},
					Messages: []Message{{Code: "0", Message: "OK"}},
				}
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(resp)
			}))
			defer server.Close()

			client, _ := NewClient(SetURL(server.URL), DisableRetry())
			service := NewContainerService(client)

			data := []byte("test file content")
			_, err := service.UploadDataWithRepetition(context.Background(), "TestDB", "TestLayout", "123", "Photo", "test.txt", data, "test-token", tt.repetition)

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

func TestContainerService_UploadFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-upload-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	testContent := []byte("test file content for upload")
	if _, err := tmpFile.Write(testContent); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	_ = tmpFile.Close()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		bodyBytes, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(bodyBytes), "test file content for upload") {
			t.Error("expected file content in multipart body")
		}

		resp := ResponseData{
			Response: Response{ModID: "2"},
			Messages: []Message{{Code: "0", Message: "OK"}},
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, _ := NewClient(SetURL(server.URL), DisableRetry())
	service := NewContainerService(client)

	response, err := service.UploadFile(context.Background(), "TestDB", "TestLayout", "123", "Photo", tmpFile.Name(), "test-token")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if response == nil {
		t.Fatal("response is nil")
	}
}

func TestContainerService_UploadFile_InvalidPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not reach server with invalid file path")
	}))
	defer server.Close()

	client, _ := NewClient(SetURL(server.URL), DisableRetry())
	service := NewContainerService(client)

	_, err := service.UploadFile(context.Background(), "TestDB", "TestLayout", "123", "Photo", "/nonexistent/file.txt", "test-token")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
	if _, ok := err.(*ValidationError); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestContainerService_UploadFile_EmptyPath(t *testing.T) {
	client, _ := NewClient(DisableRetry())
	service := NewContainerService(client)

	_, err := service.UploadFile(context.Background(), "TestDB", "TestLayout", "123", "Photo", "", "test-token")
	if err == nil {
		t.Error("expected error for empty file path")
	}
	if _, ok := err.(*ValidationError); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestContainerService_UploadFileWithRepetition(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-repetition-*.jpg")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	testContent := []byte("image data")
	_, _ = tmpFile.Write(testContent)
	_ = tmpFile.Close()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedFilename := filepath.Base(tmpFile.Name())
		bodyBytes, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(bodyBytes), expectedFilename) {
			t.Errorf("expected filename %s in multipart body", expectedFilename)
		}

		resp := ResponseData{
			Response: Response{ModID: "2"},
			Messages: []Message{{Code: "0", Message: "OK"}},
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, _ := NewClient(SetURL(server.URL), DisableRetry())
	service := NewContainerService(client)

	response, err := service.UploadFileWithRepetition(context.Background(), "TestDB", "TestLayout", "123", "Photo", tmpFile.Name(), "test-token", 2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if response == nil {
		t.Fatal("response is nil")
	}
}

func TestContainerService_UploadData_WithContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := ResponseData{
			Response: Response{ModID: "2"},
			Messages: []Message{{Code: "0", Message: "OK"}},
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, _ := NewClient(SetURL(server.URL), DisableRetry())
	service := NewContainerService(client)

	data := []byte("test data")

	t.Run("with background context", func(t *testing.T) {
		ctx := context.Background()
		_, err := service.UploadData(ctx, "TestDB", "TestLayout", "123", "Photo", "test.jpg", data, "test-token")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("with nil context", func(t *testing.T) {
		_, err := service.UploadData(context.TODO(), "TestDB", "TestLayout", "123", "Photo", "test.jpg", data, "test-token")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestNewContainerFileInfo(t *testing.T) {
	info := NewContainerFileInfo("Photo", "/path/to/file.jpg")

	if info == nil {
		t.Fatal("info is nil")
	}
	if info.FieldName != "Photo" {
		t.Errorf("expected FieldName Photo, got %s", info.FieldName)
	}
	if info.FilePath != "/path/to/file.jpg" {
		t.Errorf("expected FilePath /path/to/file.jpg, got %s", info.FilePath)
	}
	if info.Repetition != 1 {
		t.Errorf("expected Repetition 1, got %d", info.Repetition)
	}
}

func TestNewContainerDataInfo(t *testing.T) {
	data := []byte("test data")
	info := NewContainerDataInfo("Photo", "test.jpg", data)

	if info == nil {
		t.Fatal("info is nil")
	}
	if info.FieldName != "Photo" {
		t.Errorf("expected FieldName Photo, got %s", info.FieldName)
	}
	if info.Filename != "test.jpg" {
		t.Errorf("expected Filename test.jpg, got %s", info.Filename)
	}
	if len(info.Data) != len(data) {
		t.Errorf("expected Data length %d, got %d", len(data), len(info.Data))
	}
	if info.Repetition != 1 {
		t.Errorf("expected Repetition 1, got %d", info.Repetition)
	}
}

func TestContainerFileInfo_WithRepetition(t *testing.T) {
	info := NewContainerFileInfo("Photo", "/path/to/file.jpg")
	info.WithRepetition(3)

	if info.Repetition != 3 {
		t.Errorf("expected Repetition 3, got %d", info.Repetition)
	}
}

func TestContainerFileInfo_WithRepetition_Chaining(t *testing.T) {
	info := NewContainerFileInfo("Photo", "/path/to/file.jpg").WithRepetition(5)

	if info.Repetition != 5 {
		t.Errorf("expected Repetition 5, got %d", info.Repetition)
	}
}
