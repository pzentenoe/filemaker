package filemaker

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNewScriptService(t *testing.T) {
	client, _ := NewClient()
	service := NewScriptService(client)

	if service == nil {
		t.Fatal("service is nil")
	}
}

func TestScriptService_Execute(t *testing.T) {
	tests := []struct {
		name         string
		database     string
		layout       string
		script       string
		scriptParam  string
		token        string
		wantErr      bool
		errType      string
		checkRequest func(*testing.T, *http.Request)
	}{
		{
			name:        "successful script execution",
			database:    "TestDB",
			layout:      "TestLayout",
			script:      "TestScript",
			scriptParam: "test-param",
			token:       "test-token",
			wantErr:     false,
			checkRequest: func(t *testing.T, r *http.Request) {
				authHeader := r.Header.Get("Authorization")
				if authHeader != "Bearer test-token" {
					t.Errorf("expected Bearer test-token, got %s", authHeader)
				}
				if r.Method != http.MethodGet {
					t.Errorf("expected GET, got %s", r.Method)
				}
				scriptParam := r.URL.Query().Get("script.param")
				if scriptParam != "test-param" {
					t.Errorf("expected script.param test-param, got %s", scriptParam)
				}
			},
		},
		{
			name:     "script execution without parameters",
			database: "TestDB",
			layout:   "TestLayout",
			script:   "NoParamScript",
			token:    "test-token",
			wantErr:  false,
			checkRequest: func(t *testing.T, r *http.Request) {
				scriptParam := r.URL.Query().Get("script.param")
				if scriptParam != "" {
					t.Error("expected no script.param")
				}
			},
		},
		{
			name:     "script with special characters",
			database: "TestDB",
			layout:   "TestLayout",
			script:   "Script With Spaces",
			token:    "test-token",
			wantErr:  false,
		},
		{
			name:     "empty database",
			database: "",
			layout:   "TestLayout",
			script:   "TestScript",
			token:    "test-token",
			wantErr:  true,
			errType:  "ValidationError",
		},
		{
			name:     "empty layout",
			database: "TestDB",
			layout:   "",
			script:   "TestScript",
			token:    "test-token",
			wantErr:  true,
			errType:  "ValidationError",
		},
		{
			name:     "empty script",
			database: "TestDB",
			layout:   "TestLayout",
			script:   "",
			token:    "test-token",
			wantErr:  true,
			errType:  "ValidationError",
		},
		{
			name:     "empty token",
			database: "TestDB",
			layout:   "TestLayout",
			script:   "TestScript",
			token:    "",
			wantErr:  true,
			errType:  "ValidationError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.checkRequest != nil {
					tt.checkRequest(t, r)
				}
				resp := ResponseData{
					Response: Response{
						RecordID: "1",
					},
					Messages: []Message{{Code: "0", Message: "OK"}},
				}
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(resp)
			}))
			defer server.Close()

			client, _ := NewClient(SetURL(server.URL))
			service := NewScriptService(client)

			resp, err := service.Execute(context.Background(), tt.database, tt.layout, tt.script, tt.scriptParam, tt.token)

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

func TestScriptService_Execute_WithContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := ResponseData{
			Response: Response{RecordID: "1"},
			Messages: []Message{{Code: "0", Message: "OK"}},
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, _ := NewClient(SetURL(server.URL))
	service := NewScriptService(client)

	t.Run("with background context", func(t *testing.T) {
		ctx := context.Background()
		resp, err := service.Execute(ctx, "TestDB", "TestLayout", "TestScript", "", "test-token")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("response is nil")
		}
	})

	t.Run("with nil context", func(t *testing.T) {
		resp, err := service.Execute(context.TODO(), "TestDB", "TestLayout", "TestScript", "", "test-token")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("response is nil")
		}
	})
}

func TestScriptService_ExecuteAfterAction(t *testing.T) {
	tests := []struct {
		name         string
		database     string
		layout       string
		recordID     string
		script       string
		scriptParam  string
		token        string
		action       string
		wantErr      bool
		errType      string
		checkRequest func(*testing.T, *http.Request)
	}{
		{
			name:        "execute after create",
			database:    "TestDB",
			layout:      "TestLayout",
			recordID:    "1",
			script:      "AfterCreateScript",
			scriptParam: "param",
			token:       "test-token",
			action:      "create",
			wantErr:     false,
			checkRequest: func(t *testing.T, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST, got %s", r.Method)
				}
				script := r.URL.Query().Get("script")
				if script != "AfterCreateScript" {
					t.Errorf("expected script AfterCreateScript, got %s", script)
				}
			},
		},
		{
			name:     "execute after edit",
			database: "TestDB",
			layout:   "TestLayout",
			recordID: "123",
			script:   "AfterEditScript",
			token:    "test-token",
			action:   "edit",
			wantErr:  false,
			checkRequest: func(t *testing.T, r *http.Request) {
				if r.Method != http.MethodPatch {
					t.Errorf("expected PATCH, got %s", r.Method)
				}
			},
		},
		{
			name:     "execute after delete",
			database: "TestDB",
			layout:   "TestLayout",
			recordID: "123",
			script:   "AfterDeleteScript",
			token:    "test-token",
			action:   "delete",
			wantErr:  false,
			checkRequest: func(t *testing.T, r *http.Request) {
				if r.Method != http.MethodDelete {
					t.Errorf("expected DELETE, got %s", r.Method)
				}
			},
		},
		{
			name:     "invalid action",
			database: "TestDB",
			layout:   "TestLayout",
			recordID: "123",
			script:   "TestScript",
			token:    "test-token",
			action:   "invalid",
			wantErr:  true,
			errType:  "ValidationError",
		},
		{
			name:     "empty database",
			database: "",
			layout:   "TestLayout",
			recordID: "123",
			script:   "TestScript",
			token:    "test-token",
			action:   "edit",
			wantErr:  true,
			errType:  "ValidationError",
		},
		{
			name:     "empty layout",
			database: "TestDB",
			layout:   "",
			recordID: "123",
			script:   "TestScript",
			token:    "test-token",
			action:   "edit",
			wantErr:  true,
			errType:  "ValidationError",
		},
		{
			name:     "empty recordID",
			database: "TestDB",
			layout:   "TestLayout",
			recordID: "",
			script:   "TestScript",
			token:    "test-token",
			action:   "edit",
			wantErr:  true,
			errType:  "ValidationError",
		},
		{
			name:     "empty script",
			database: "TestDB",
			layout:   "TestLayout",
			recordID: "123",
			script:   "",
			token:    "test-token",
			action:   "edit",
			wantErr:  true,
			errType:  "ValidationError",
		},
		{
			name:     "empty token",
			database: "TestDB",
			layout:   "TestLayout",
			recordID: "123",
			script:   "TestScript",
			token:    "",
			action:   "edit",
			wantErr:  true,
			errType:  "ValidationError",
		},
		{
			name:     "empty action",
			database: "TestDB",
			layout:   "TestLayout",
			recordID: "123",
			script:   "TestScript",
			token:    "test-token",
			action:   "",
			wantErr:  true,
			errType:  "ValidationError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.checkRequest != nil {
					tt.checkRequest(t, r)
				}
				resp := ResponseData{
					Response: Response{RecordID: "123"},
					Messages: []Message{{Code: "0", Message: "OK"}},
				}
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(resp)
			}))
			defer server.Close()

			client, _ := NewClient(SetURL(server.URL))
			service := NewScriptService(client)

			resp, err := service.ExecuteAfterAction(context.Background(), tt.database, tt.layout, tt.recordID, tt.script, tt.scriptParam, tt.token, tt.action)

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

func TestScriptParameter_ToQueryParams(t *testing.T) {
	tests := []struct {
		name     string
		param    *ScriptParameter
		prefix   string
		expected map[string]string
	}{
		{
			name: "with script and param",
			param: &ScriptParameter{
				Script:      "TestScript",
				ScriptParam: "test-param",
			},
			prefix: "script",
			expected: map[string]string{
				"script":       "TestScript",
				"script.param": "test-param",
			},
		},
		{
			name: "with script only",
			param: &ScriptParameter{
				Script: "TestScript",
			},
			prefix: "script",
			expected: map[string]string{
				"script": "TestScript",
			},
		},
		{
			name:     "empty script",
			param:    &ScriptParameter{},
			prefix:   "script",
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := tt.param.ToQueryParams(tt.prefix)
			for key, expectedValue := range tt.expected {
				if params.Get(key) != expectedValue {
					t.Errorf("expected %s=%s, got %s", key, expectedValue, params.Get(key))
				}
			}
			for key := range params {
				if _, ok := tt.expected[key]; !ok {
					t.Errorf("unexpected parameter: %s", key)
				}
			}
		})
	}
}

func TestScriptParameter_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		param    *ScriptParameter
		expected string
	}{
		{
			name: "with script and param",
			param: &ScriptParameter{
				Script:      "TestScript",
				ScriptParam: "test-param",
			},
			expected: `{"script":"TestScript","script.param":"test-param"}`,
		},
		{
			name: "with script only",
			param: &ScriptParameter{
				Script: "TestScript",
			},
			expected: `{"script":"TestScript"}`,
		},
		{
			name:     "empty script",
			param:    &ScriptParameter{},
			expected: `{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.param)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(data))
			}
		})
	}
}

func TestNewScriptContext(t *testing.T) {
	sc := NewScriptContext()
	if sc == nil {
		t.Fatal("script context is nil")
	}
}

func TestScriptContext_WithPreRequest(t *testing.T) {
	sc := NewScriptContext()
	sc.WithPreRequest("PreRequestScript", "param1")

	if sc.PreRequest == nil {
		t.Fatal("PreRequest is nil")
	}
	if sc.PreRequest.Script != "PreRequestScript" {
		t.Errorf("expected PreRequestScript, got %s", sc.PreRequest.Script)
	}
	if sc.PreRequest.ScriptParam != "param1" {
		t.Errorf("expected param1, got %s", sc.PreRequest.ScriptParam)
	}
}

func TestScriptContext_WithPreSort(t *testing.T) {
	sc := NewScriptContext()
	sc.WithPreSort("PreSortScript", "param2")

	if sc.PreSort == nil {
		t.Fatal("PreSort is nil")
	}
	if sc.PreSort.Script != "PreSortScript" {
		t.Errorf("expected PreSortScript, got %s", sc.PreSort.Script)
	}
	if sc.PreSort.ScriptParam != "param2" {
		t.Errorf("expected param2, got %s", sc.PreSort.ScriptParam)
	}
}

func TestScriptContext_WithAfter(t *testing.T) {
	sc := NewScriptContext()
	sc.WithAfter("AfterScript", "param3")

	if sc.After == nil {
		t.Fatal("After is nil")
	}
	if sc.After.Script != "AfterScript" {
		t.Errorf("expected AfterScript, got %s", sc.After.Script)
	}
	if sc.After.ScriptParam != "param3" {
		t.Errorf("expected param3, got %s", sc.After.ScriptParam)
	}
}

func TestScriptContext_ToQueryParams(t *testing.T) {
	tests := []struct {
		name     string
		context  *ScriptContext
		expected url.Values
	}{
		{
			name: "all scripts with params",
			context: &ScriptContext{
				PreRequest: &ScriptParameter{
					Script:      "PreRequestScript",
					ScriptParam: "prereq-param",
				},
				PreSort: &ScriptParameter{
					Script:      "PreSortScript",
					ScriptParam: "presort-param",
				},
				After: &ScriptParameter{
					Script:      "AfterScript",
					ScriptParam: "after-param",
				},
			},
			expected: url.Values{
				"script.prerequest":       []string{"PreRequestScript"},
				"script.prerequest.param": []string{"prereq-param"},
				"script.presort":          []string{"PreSortScript"},
				"script.presort.param":    []string{"presort-param"},
				"script":                  []string{"AfterScript"},
				"script.param":            []string{"after-param"},
			},
		},
		{
			name: "only after script",
			context: &ScriptContext{
				After: &ScriptParameter{
					Script: "AfterScript",
				},
			},
			expected: url.Values{
				"script": []string{"AfterScript"},
			},
		},
		{
			name:     "empty context",
			context:  &ScriptContext{},
			expected: url.Values{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := tt.context.ToQueryParams()
			for key, expectedValues := range tt.expected {
				actualValues := params[key]
				if len(actualValues) != len(expectedValues) {
					t.Errorf("key %s: expected %v, got %v", key, expectedValues, actualValues)
					continue
				}
				for i, expectedValue := range expectedValues {
					if actualValues[i] != expectedValue {
						t.Errorf("key %s: expected %v, got %v", key, expectedValues, actualValues)
					}
				}
			}
			for key := range params {
				if _, ok := tt.expected[key]; !ok {
					t.Errorf("unexpected parameter: %s", key)
				}
			}
		})
	}
}

func TestScriptContext_ChainedMethods(t *testing.T) {
	sc := NewScriptContext().
		WithPreRequest("PreScript", "pre-param").
		WithPreSort("SortScript", "sort-param").
		WithAfter("AfterScript", "after-param")

	params := sc.ToQueryParams()

	expected := map[string]string{
		"script.prerequest":       "PreScript",
		"script.prerequest.param": "pre-param",
		"script.presort":          "SortScript",
		"script.presort.param":    "sort-param",
		"script":                  "AfterScript",
		"script.param":            "after-param",
	}

	for key, expectedValue := range expected {
		if params.Get(key) != expectedValue {
			t.Errorf("expected %s=%s, got %s", key, expectedValue, params.Get(key))
		}
	}
}
