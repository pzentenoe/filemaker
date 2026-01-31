package filemaker

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewSearchService(t *testing.T) {
	client, _ := NewClient()
	service := NewSearchService("TestDB", "TestLayout", client)

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

func TestSearchService_GroupQueries(t *testing.T) {
	t.Run("single query group", func(t *testing.T) {
		search := NewSearchService("test", "test_layout", nil)
		search.GroupQueries(
			NewGroupQuery(
				NewQueryFieldOperator("nombre", "pablo", Equal),
				NewQueryFieldOperator("apellido", "zenteno", Equal),
			),
		)
		b, _ := json.Marshal(search.searchData.QueryGroup)
		expected := "[{\"apellido\":\"==zenteno\",\"nombre\":\"==pablo\"}]"
		if string(b) != expected {
			t.Errorf("expected %s, got %s", expected, string(b))
		}
	})

	t.Run("multiple query groups", func(t *testing.T) {
		search := NewSearchService("test", "test_layout", nil)
		search.GroupQueries(
			NewGroupQuery(
				NewQueryFieldOperator("nombre", "pablo", Equal),
			),
			NewGroupQuery(
				NewQueryFieldOperator("apellido", "zenteno", Equal),
			),
		)
		if len(search.searchData.QueryGroup) != 2 {
			t.Errorf("expected 2 query groups, got %d", len(search.searchData.QueryGroup))
		}
	})
}

func TestSearchService_SetOffset(t *testing.T) {
	search := NewSearchService("test", "test_layout", nil)
	search.SetOffset("10")

	if search.searchData.Offset != "10" {
		t.Errorf("expected offset 10, got %s", search.searchData.Offset)
	}
}

func TestSearchService_SetLimit(t *testing.T) {
	search := NewSearchService("test", "test_layout", nil)
	search.SetLimit("50")

	if search.searchData.Limit != "50" {
		t.Errorf("expected limit 50, got %s", search.searchData.Limit)
	}
}

func TestSearchService_SetPortals(t *testing.T) {
	search := NewSearchService("test", "test_layout", nil)
	portals := []string{"Portal1", "Portal2"}
	search.SetPortals(portals)

	if len(search.searchData.Portal) != 2 {
		t.Errorf("expected 2 portals, got %d", len(search.searchData.Portal))
	}
}

func TestSearchService_Sorters(t *testing.T) {
	search := NewSearchService("test", "test_layout", nil)
	sorters := []*Sorter{
		NewSorter("FirstName", Ascending),
		NewSorter("LastName", Descending),
	}
	search.Sorters(sorters...)

	if len(search.searchData.Sort) != 2 {
		t.Errorf("expected 2 sorters, got %d", len(search.searchData.Sort))
	}
}

func TestSearchService_Do(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(*searchService)
		wantErr   bool
		checkResp func(*testing.T, *ResponseData)
	}{
		{
			name: "successful search",
			setupFunc: func(s *searchService) {
				s.GroupQueries(
					NewGroupQuery(
						NewQueryFieldOperator("Name", "John", Equal),
					),
				)
				s.SetLimit("10")
				s.SetOffset("0")
			},
			wantErr: false,
			checkResp: func(t *testing.T, resp *ResponseData) {
				if len(resp.Response.Data) == 0 {
					t.Error("expected data, got empty")
				}
			},
		},
		{
			name: "search with sorting",
			setupFunc: func(s *searchService) {
				s.GroupQueries(
					NewGroupQuery(
						NewQueryFieldOperator("Status", "Active", Equal),
					),
				)
				s.Sorters(NewSorter("Name", Ascending))
			},
			wantErr: false,
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
					if r.Method != http.MethodPost {
						t.Errorf("expected POST, got %s", r.Method)
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
				SetUsername("user"),
				SetPassword("pass"),
			)

			service := NewSearchService("TestDB", "TestLayout", client)
			if tt.setupFunc != nil {
				tt.setupFunc(service)
			}

			resp, err := service.Do(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.checkResp != nil {
				tt.checkResp(t, resp)
			}
		})
	}
}

func TestSearchService_Do_WithContext(t *testing.T) {
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
					Data: []Datum{{RecordID: "1"}},
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
		SetUsername("user"),
		SetPassword("pass"),
	)

	service := NewSearchService("TestDB", "TestLayout", client)
	service.GroupQueries(
		NewGroupQuery(
			NewQueryFieldOperator("ID", "1", Equal),
		),
	)

	resp, err := service.Do(context.TODO())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("response is nil")
	}
}
