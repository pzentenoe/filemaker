package filemaker

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewGlobalFieldsService(t *testing.T) {
	client, _ := NewClient()
	service := NewGlobalFieldsService(client)

	if service == nil {
		t.Fatal("service is nil")
	}
}

func TestGlobalFieldsService_SetGlobalFields(t *testing.T) {
	tests := []struct {
		name         string
		database     string
		globalFields map[string]any
		token        string
		wantErr      bool
		errType      string
		checkRequest func(*testing.T, *http.Request)
	}{
		{
			name:     "successful set global fields",
			database: "TestDB",
			globalFields: map[string]any{
				"Global::UserName": "john@example.com",
				"Global::Theme":    "dark",
			},
			token:   "test-token",
			wantErr: false,
			checkRequest: func(t *testing.T, r *http.Request) {
				if r.Method != http.MethodPatch {
					t.Errorf("expected PATCH, got %s", r.Method)
				}
				authHeader := r.Header.Get("Authorization")
				if authHeader != "Bearer test-token" {
					t.Errorf("expected Bearer test-token, got %s", authHeader)
				}
				var body map[string]any
				_ = json.NewDecoder(r.Body).Decode(&body)
				if _, ok := body["globalFields"]; !ok {
					t.Error("expected globalFields in body")
				}
			},
		},
		{
			name:     "empty database",
			database: "",
			globalFields: map[string]any{
				"Global::Field": "value",
			},
			token:   "test-token",
			wantErr: true,
			errType: "ValidationError",
		},
		{
			name:         "nil global fields",
			database:     "TestDB",
			globalFields: nil,
			token:        "test-token",
			wantErr:      true,
			errType:      "ValidationError",
		},
		{
			name:         "empty global fields",
			database:     "TestDB",
			globalFields: map[string]any{},
			token:        "test-token",
			wantErr:      true,
			errType:      "ValidationError",
		},
		{
			name:     "empty token",
			database: "TestDB",
			globalFields: map[string]any{
				"Global::Field": "value",
			},
			token:   "",
			wantErr: true,
			errType: "ValidationError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.checkRequest != nil {
					tt.checkRequest(t, r)
				}
				resp := ResponseData{
					Response: Response{RecordID: "1"},
					Messages: []Message{{Code: "0", Message: "OK"}},
				}
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(resp)
			}))
			defer server.Close()

			client, _ := NewClient(SetURL(server.URL))
			service := NewGlobalFieldsService(client)

			resp, err := service.SetGlobalFields(context.Background(), tt.database, tt.globalFields, tt.token)

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

func TestGlobalFieldsService_SetGlobalFields_WithContext(t *testing.T) {
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
	service := NewGlobalFieldsService(client)

	globalFields := map[string]any{
		"Global::Field": "value",
	}

	t.Run("with background context", func(t *testing.T) {
		ctx := context.Background()
		_, err := service.SetGlobalFields(ctx, "TestDB", globalFields, "test-token")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("with nil context", func(t *testing.T) {
		_, err := service.SetGlobalFields(context.TODO(), "TestDB", globalFields, "test-token")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestNewGlobalField(t *testing.T) {
	field := NewGlobalField("TestField", "TestValue")

	if field == nil {
		t.Fatal("field is nil")
	}
	if field.Name != "TestField" {
		t.Errorf("expected name TestField, got %s", field.Name)
	}
	if field.Value != "TestValue" {
		t.Errorf("expected value TestValue, got %v", field.Value)
	}
}

func TestNewGlobalFieldsBuilder(t *testing.T) {
	builder := NewGlobalFieldsBuilder()

	if builder == nil {
		t.Fatal("builder is nil")
	}
	if builder.fields == nil {
		t.Fatal("fields map is nil")
	}
}

func TestGlobalFieldsBuilder_Add(t *testing.T) {
	builder := NewGlobalFieldsBuilder()
	builder.Add("Field1", "Value1")
	builder.Add("Field2", 42)
	builder.Add("Field3", true)

	if len(builder.fields) != 3 {
		t.Errorf("expected 3 fields, got %d", len(builder.fields))
	}
	if builder.fields["Field1"] != "Value1" {
		t.Errorf("expected Field1 to be Value1, got %v", builder.fields["Field1"])
	}
	if builder.fields["Field2"] != 42 {
		t.Errorf("expected Field2 to be 42, got %v", builder.fields["Field2"])
	}
	if builder.fields["Field3"] != true {
		t.Errorf("expected Field3 to be true, got %v", builder.fields["Field3"])
	}
}

func TestGlobalFieldsBuilder_Add_Chaining(t *testing.T) {
	builder := NewGlobalFieldsBuilder().
		Add("Field1", "Value1").
		Add("Field2", "Value2")

	if len(builder.fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(builder.fields))
	}
}

func TestGlobalFieldsBuilder_AddField(t *testing.T) {
	builder := NewGlobalFieldsBuilder()
	field := NewGlobalField("TestField", "TestValue")
	builder.AddField(field)

	if len(builder.fields) != 1 {
		t.Errorf("expected 1 field, got %d", len(builder.fields))
	}
	if builder.fields["TestField"] != "TestValue" {
		t.Errorf("expected TestField to be TestValue, got %v", builder.fields["TestField"])
	}
}

func TestGlobalFieldsBuilder_AddFields(t *testing.T) {
	builder := NewGlobalFieldsBuilder()
	fields := []*GlobalField{
		NewGlobalField("Field1", "Value1"),
		NewGlobalField("Field2", "Value2"),
		NewGlobalField("Field3", "Value3"),
	}
	builder.AddFields(fields...)

	if len(builder.fields) != 3 {
		t.Errorf("expected 3 fields, got %d", len(builder.fields))
	}
	if builder.fields["Field1"] != "Value1" {
		t.Errorf("expected Field1 to be Value1, got %v", builder.fields["Field1"])
	}
	if builder.fields["Field2"] != "Value2" {
		t.Errorf("expected Field2 to be Value2, got %v", builder.fields["Field2"])
	}
	if builder.fields["Field3"] != "Value3" {
		t.Errorf("expected Field3 to be Value3, got %v", builder.fields["Field3"])
	}
}

func TestGlobalFieldsBuilder_Build(t *testing.T) {
	builder := NewGlobalFieldsBuilder()
	builder.Add("Field1", "Value1")
	builder.Add("Field2", "Value2")

	result := builder.Build()

	if len(result) != 2 {
		t.Errorf("expected 2 fields, got %d", len(result))
	}
	if result["Field1"] != "Value1" {
		t.Errorf("expected Field1 to be Value1, got %v", result["Field1"])
	}
	if result["Field2"] != "Value2" {
		t.Errorf("expected Field2 to be Value2, got %v", result["Field2"])
	}
}

func TestGlobalFieldsBuilder_Clear(t *testing.T) {
	builder := NewGlobalFieldsBuilder()
	builder.Add("Field1", "Value1")
	builder.Add("Field2", "Value2")

	if len(builder.fields) != 2 {
		t.Errorf("expected 2 fields before clear, got %d", len(builder.fields))
	}

	builder.Clear()

	if len(builder.fields) != 0 {
		t.Errorf("expected 0 fields after clear, got %d", len(builder.fields))
	}
}

func TestGlobalFieldsBuilder_Clear_Chaining(t *testing.T) {
	builder := NewGlobalFieldsBuilder().
		Add("Field1", "Value1").
		Clear().
		Add("Field2", "Value2")

	if len(builder.fields) != 1 {
		t.Errorf("expected 1 field, got %d", len(builder.fields))
	}
	if _, ok := builder.fields["Field1"]; ok {
		t.Error("expected Field1 to be cleared")
	}
	if builder.fields["Field2"] != "Value2" {
		t.Errorf("expected Field2 to be Value2, got %v", builder.fields["Field2"])
	}
}

func TestGlobalFieldsBuilder_Count(t *testing.T) {
	builder := NewGlobalFieldsBuilder()

	if builder.Count() != 0 {
		t.Errorf("expected count 0, got %d", builder.Count())
	}

	builder.Add("Field1", "Value1")
	if builder.Count() != 1 {
		t.Errorf("expected count 1, got %d", builder.Count())
	}

	builder.Add("Field2", "Value2")
	if builder.Count() != 2 {
		t.Errorf("expected count 2, got %d", builder.Count())
	}

	builder.Clear()
	if builder.Count() != 0 {
		t.Errorf("expected count 0 after clear, got %d", builder.Count())
	}
}

func TestGlobalFieldsBuilder_CompleteWorkflow(t *testing.T) {
	builder := NewGlobalFieldsBuilder().
		Add("Global::UserName", "john@example.com").
		Add("Global::Theme", "dark").
		AddField(NewGlobalField("Global::Language", "en")).
		AddFields(
			NewGlobalField("Global::TimeZone", "UTC"),
			NewGlobalField("Global::Notifications", true),
		)

	if builder.Count() != 5 {
		t.Errorf("expected count 5, got %d", builder.Count())
	}

	fields := builder.Build()

	if len(fields) != 5 {
		t.Errorf("expected 5 fields, got %d", len(fields))
	}

	expected := map[string]any{
		"Global::UserName":      "john@example.com",
		"Global::Theme":         "dark",
		"Global::Language":      "en",
		"Global::TimeZone":      "UTC",
		"Global::Notifications": true,
	}

	for key, expectedValue := range expected {
		if fields[key] != expectedValue {
			t.Errorf("expected %s to be %v, got %v", key, expectedValue, fields[key])
		}
	}
}
