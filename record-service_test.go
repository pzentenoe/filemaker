package filemaker

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRecordService(t *testing.T) {
	client, _ := NewClient()
	service := NewRecordService("TestDB", "TestLayout", client)

	if service == nil {
		t.Fatal("service is nil")
	}
	if service.database != "TestDB" {
		t.Errorf("expected database TestDB, got %s", service.database)
	}
	if service.layout != "TestLayout" {
		t.Errorf("expected layout TestLayout, got %s", service.layout)
	}
}

func TestRecordService_Create(t *testing.T) {
	tests := []struct {
		name       string
		payload    *Payload
		wantErr    bool
		statusCode int
		response   ResponseData
	}{
		{
			name: "successful record creation",
			payload: &Payload{
				FieldData: map[string]any{
					"FirstName": "John",
					"LastName":  "Doe",
				},
			},
			wantErr:    false,
			statusCode: http.StatusOK,
			response: ResponseData{
				Response: Response{
					RecordID: "123",
					ModID:    "1",
				},
				Messages: []Message{{Code: "0", Message: "OK"}},
			},
		},
		{
			name: "creation with portal data",
			payload: &Payload{
				FieldData: map[string]any{
					"Name": "Company",
				},
				PortalData: map[string]any{
					"Employees": []map[string]any{
						{"Name": "Alice"},
					},
				},
			},
			wantErr:    false,
			statusCode: http.StatusOK,
			response: ResponseData{
				Response: Response{RecordID: "456"},
				Messages: []Message{{Code: "0", Message: "OK"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callCount := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				callCount++
				if callCount == 1 {
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(ResponseData{
						Response: Response{Token: "test-token"},
						Messages: []Message{{Code: "0", Message: "OK"}},
					})
					return
				}
				if callCount == 2 {
					w.WriteHeader(tt.statusCode)
					_ = json.NewEncoder(w).Encode(tt.response)
					return
				}
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(ResponseData{
					Messages: []Message{{Code: "0", Message: "OK"}},
				})
			}))
			defer server.Close()

			client, _ := NewClient(
				SetURL(server.URL),
				SetBasicAuth("user", "pass"),
			)

			service := NewRecordService("TestDB", "TestLayout", client)
			resp, err := service.Create(context.Background(), tt.payload)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if resp.Response.RecordID == "" {
				t.Error("expected non-empty recordID")
			}
		})
	}
}

func TestRecordService_Edit(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(ResponseData{
				Response: Response{Token: "test-token"},
				Messages: []Message{{Code: "0", Message: "OK"}},
			})
			return
		}
		if callCount == 2 {
			if r.Method != http.MethodPatch {
				t.Errorf("expected PATCH, got %s", r.Method)
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(ResponseData{
				Response: Response{ModID: "2"},
				Messages: []Message{{Code: "0", Message: "OK"}},
			})
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ResponseData{
			Messages: []Message{{Code: "0", Message: "OK"}},
		})
	}))
	defer server.Close()

	client, _ := NewClient(
		SetURL(server.URL),
		SetBasicAuth("user", "pass"),
	)

	service := NewRecordService("TestDB", "TestLayout", client)
	payload := &Payload{
		FieldData: map[string]any{"FirstName": "Jane"},
	}

	resp, err := service.Edit(context.Background(), "123", payload)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp.Response.ModID != "2" {
		t.Errorf("expected modID 2, got %s", resp.Response.ModID)
	}
}

func TestRecordService_Duplicate(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(ResponseData{
				Response: Response{Token: "test-token"},
				Messages: []Message{{Code: "0", Message: "OK"}},
			})
			return
		}
		if callCount == 2 {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(ResponseData{
				Response: Response{RecordID: "456"},
				Messages: []Message{{Code: "0", Message: "OK"}},
			})
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ResponseData{
			Messages: []Message{{Code: "0", Message: "OK"}},
		})
	}))
	defer server.Close()

	client, _ := NewClient(
		SetURL(server.URL),
		SetBasicAuth("user", "pass"),
	)

	service := NewRecordService("TestDB", "TestLayout", client)
	resp, err := service.Duplicate(context.Background(), "123")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp.Response.RecordID != "456" {
		t.Errorf("expected recordID 456, got %s", resp.Response.RecordID)
	}
}

func TestRecordService_Delete(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(ResponseData{
				Response: Response{Token: "test-token"},
				Messages: []Message{{Code: "0", Message: "OK"}},
			})
			return
		}
		if callCount == 2 {
			if r.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", r.Method)
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(ResponseData{
				Messages: []Message{{Code: "0", Message: "OK"}},
			})
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ResponseData{
			Messages: []Message{{Code: "0", Message: "OK"}},
		})
	}))
	defer server.Close()

	client, _ := NewClient(
		SetURL(server.URL),
		SetBasicAuth("user", "pass"),
	)

	service := NewRecordService("TestDB", "TestLayout", client)
	_, err := service.Delete(context.Background(), "123", "")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRecordService_Delete_WithDeleteRelated(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(ResponseData{
				Response: Response{Token: "test-token"},
				Messages: []Message{{Code: "0", Message: "OK"}},
			})
			return
		}
		if callCount == 2 {
			if r.URL.Query().Get("deleteRelated") != "MyPortal" {
				t.Errorf("expected deleteRelated=MyPortal, got %s", r.URL.Query().Get("deleteRelated"))
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(ResponseData{
				Messages: []Message{{Code: "0", Message: "OK"}},
			})
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ResponseData{
			Messages: []Message{{Code: "0", Message: "OK"}},
		})
	}))
	defer server.Close()

	client, _ := NewClient(SetURL(server.URL), SetBasicAuth("u", "p"))
	service := NewRecordService("DB", "Layout", client)
	_, err := service.Delete(context.Background(), "123", "MyPortal")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRecordService_GetById(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(ResponseData{
				Response: Response{Token: "test-token"},
				Messages: []Message{{Code: "0", Message: "OK"}},
			})
			return
		}
		if callCount == 2 {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(ResponseData{
				Response: Response{
					RecordID: "123",
					ModID:    "1",
					Data: []Datum{
						{
							RecordID: "123",
							ModID:    "1",
							FieldData: map[string]any{
								"FirstName": "John",
							},
						},
					},
				},
				Messages: []Message{{Code: "0", Message: "OK"}},
			})
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ResponseData{
			Messages: []Message{{Code: "0", Message: "OK"}},
		})
	}))
	defer server.Close()

	client, _ := NewClient(
		SetURL(server.URL),
		SetBasicAuth("user", "pass"),
	)

	service := NewRecordService("TestDB", "TestLayout", client)
	resp, err := service.GetById(context.Background(), "123")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp.Response.RecordID != "123" {
		t.Errorf("expected recordID 123, got %s", resp.Response.RecordID)
	}
}

func TestRecordService_List(t *testing.T) {
	tests := []struct {
		name    string
		offset  string
		limit   string
		sorters []*Sorter
		check   func(*testing.T, *http.Request)
	}{
		{
			name:    "list without sorting",
			offset:  "1",
			limit:   "10",
			sorters: nil,
			check: func(t *testing.T, r *http.Request) {
				if r.URL.Query().Get("_offset") != "1" {
					t.Errorf("expected offset 1, got %s", r.URL.Query().Get("_offset"))
				}
				if r.URL.Query().Get("_limit") != "10" {
					t.Errorf("expected limit 10, got %s", r.URL.Query().Get("_limit"))
				}
			},
		},
		{
			name:   "list with sorting",
			offset: "0",
			limit:  "20",
			sorters: []*Sorter{
				NewSorter("FirstName", Ascending),
				NewSorter("LastName", Descending),
			},
			check: func(t *testing.T, r *http.Request) {
				sortParam := r.URL.Query().Get("_sort")
				if sortParam == "" {
					t.Error("expected _sort parameter")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callCount := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				callCount++
				if callCount == 1 {
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(ResponseData{
						Response: Response{Token: "test-token"},
						Messages: []Message{{Code: "0", Message: "OK"}},
					})
					return
				}
				if callCount == 2 {
					if tt.check != nil {
						tt.check(t, r)
					}
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(ResponseData{
						Response: Response{
							Data: []Datum{
								{RecordID: "1"},
								{RecordID: "2"},
							},
						},
						Messages: []Message{{Code: "0", Message: "OK"}},
					})
					return
				}
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(ResponseData{
					Messages: []Message{{Code: "0", Message: "OK"}},
				})
			}))
			defer server.Close()

			client, _ := NewClient(
				SetURL(server.URL),
				SetBasicAuth("user", "pass"),
			)

			service := NewRecordService("TestDB", "TestLayout", client)
			resp, err := service.List(context.Background(), tt.offset, tt.limit, tt.sorters...)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if len(resp.Response.Data) == 0 {
				t.Error("expected data, got empty array")
			}
		})
	}
}

func TestSortersToJson(t *testing.T) {
	tests := []struct {
		name     string
		sorters  []*Sorter
		expected string
	}{
		{
			name:     "empty sorters",
			sorters:  []*Sorter{},
			expected: "",
		},
		{
			name:     "nil sorters",
			sorters:  nil,
			expected: "",
		},
		{
			name: "single sorter",
			sorters: []*Sorter{
				NewSorter("Name", Ascending),
			},
			expected: `[{"fieldName":"Name","sortOrder":"ascend"}]`,
		},
		{
			name: "multiple sorters",
			sorters: []*Sorter{
				NewSorter("FirstName", Ascending),
				NewSorter("LastName", Descending),
			},
			expected: `[{"fieldName":"FirstName","sortOrder":"ascend"},{"fieldName":"LastName","sortOrder":"descend"}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sortersToJson(tt.sorters...)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestRecordService_withAuth_ContextHandling(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(ResponseData{
				Response: Response{Token: "test-token"},
				Messages: []Message{{Code: "0", Message: "OK"}},
			})
			return
		}
		if callCount == 2 {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(ResponseData{
				Response: Response{RecordID: "123"},
				Messages: []Message{{Code: "0", Message: "OK"}},
			})
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ResponseData{
			Messages: []Message{{Code: "0", Message: "OK"}},
		})
	}))
	defer server.Close()

	client, _ := NewClient(
		SetURL(server.URL),
		SetBasicAuth("user", "pass"),
	)

	service := NewRecordService("TestDB", "TestLayout", client)
	payload := &Payload{
		FieldData: map[string]any{"Name": "Test"},
	}

	resp, err := service.Create(context.TODO(), payload)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("response is nil")
	}
}
