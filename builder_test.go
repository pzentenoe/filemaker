package filemaker

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRecordBuilder(t *testing.T) {
	client, _ := NewClient()
	builder := NewRecordBuilder(client, "TestDB", "TestLayout")

	if builder == nil {
		t.Fatal("builder is nil")
	}
	if builder.database != "TestDB" {
		t.Errorf("expected database TestDB, got %s", builder.database)
	}
	if builder.layout != "TestLayout" {
		t.Errorf("expected layout TestLayout, got %s", builder.layout)
	}
	if builder.fieldData == nil {
		t.Fatal("fieldData is nil")
	}
	if builder.portalData == nil {
		t.Fatal("portalData is nil")
	}
}

func TestRecordBuilder_SetField(t *testing.T) {
	client, _ := NewClient()
	builder := NewRecordBuilder(client, "TestDB", "TestLayout")

	builder.SetField("FirstName", "John")
	builder.SetField("Age", 30)

	if builder.fieldData["FirstName"] != "John" {
		t.Errorf("expected FirstName John, got %v", builder.fieldData["FirstName"])
	}
	if builder.fieldData["Age"] != 30 {
		t.Errorf("expected Age 30, got %v", builder.fieldData["Age"])
	}
}

func TestRecordBuilder_SetField_Chaining(t *testing.T) {
	client, _ := NewClient()
	builder := NewRecordBuilder(client, "TestDB", "TestLayout").
		SetField("FirstName", "John").
		SetField("LastName", "Doe")

	if len(builder.fieldData) != 2 {
		t.Errorf("expected 2 fields, got %d", len(builder.fieldData))
	}
}

func TestRecordBuilder_SetFields(t *testing.T) {
	client, _ := NewClient()
	builder := NewRecordBuilder(client, "TestDB", "TestLayout")

	fields := map[string]any{
		"FirstName": "Jane",
		"LastName":  "Smith",
		"Email":     "jane@example.com",
	}
	builder.SetFields(fields)

	if len(builder.fieldData) != 3 {
		t.Errorf("expected 3 fields, got %d", len(builder.fieldData))
	}
	if builder.fieldData["Email"] != "jane@example.com" {
		t.Errorf("expected Email jane@example.com, got %v", builder.fieldData["Email"])
	}
}

func TestRecordBuilder_AddPortalRecord(t *testing.T) {
	client, _ := NewClient()
	builder := NewRecordBuilder(client, "TestDB", "TestLayout")

	portal1 := map[string]any{"ItemName": "Product A"}
	portal2 := map[string]any{"ItemName": "Product B"}

	builder.AddPortalRecord("LineItems", portal1)
	builder.AddPortalRecord("LineItems", portal2)

	if len(builder.portalData["LineItems"]) != 2 {
		t.Errorf("expected 2 portal records, got %d", len(builder.portalData["LineItems"]))
	}
}

func TestRecordBuilder_SetPortalRecords(t *testing.T) {
	client, _ := NewClient()
	builder := NewRecordBuilder(client, "TestDB", "TestLayout")

	records := []map[string]any{
		{"ItemName": "Product A"},
		{"ItemName": "Product B"},
		{"ItemName": "Product C"},
	}
	builder.SetPortalRecords("LineItems", records)

	if len(builder.portalData["LineItems"]) != 3 {
		t.Errorf("expected 3 portal records, got %d", len(builder.portalData["LineItems"]))
	}
}

func TestRecordBuilder_WithScripts(t *testing.T) {
	client, _ := NewClient()
	builder := NewRecordBuilder(client, "TestDB", "TestLayout")

	scripts := NewScriptContext().WithAfter("AfterScript", "param")
	builder.WithScripts(scripts)

	if builder.scripts == nil {
		t.Fatal("scripts is nil")
	}
	if builder.scripts.After.Script != "AfterScript" {
		t.Errorf("expected AfterScript, got %s", builder.scripts.After.Script)
	}
}

func TestRecordBuilder_WithPreRequestScript(t *testing.T) {
	client, _ := NewClient()
	builder := NewRecordBuilder(client, "TestDB", "TestLayout")

	builder.WithPreRequestScript("PreScript", "pre-param")

	if builder.scripts == nil {
		t.Fatal("scripts is nil")
	}
	if builder.scripts.PreRequest.Script != "PreScript" {
		t.Errorf("expected PreScript, got %s", builder.scripts.PreRequest.Script)
	}
}

func TestRecordBuilder_WithAfterScript(t *testing.T) {
	client, _ := NewClient()
	builder := NewRecordBuilder(client, "TestDB", "TestLayout")

	builder.WithAfterScript("AfterScript", "after-param")

	if builder.scripts == nil {
		t.Fatal("scripts is nil")
	}
	if builder.scripts.After.Script != "AfterScript" {
		t.Errorf("expected AfterScript, got %s", builder.scripts.After.Script)
	}
}

func TestRecordBuilder_ForRecord(t *testing.T) {
	client, _ := NewClient()
	builder := NewRecordBuilder(client, "TestDB", "TestLayout")

	builder.ForRecord("123")

	if builder.recordID != "123" {
		t.Errorf("expected recordID 123, got %s", builder.recordID)
	}
}

func TestRecordBuilder_WithModID(t *testing.T) {
	client, _ := NewClient()
	builder := NewRecordBuilder(client, "TestDB", "TestLayout")

	builder.WithModID("5")

	if builder.modID != "5" {
		t.Errorf("expected modID 5, got %s", builder.modID)
	}
}

func TestRecordBuilder_WithDeleteRelated(t *testing.T) {
	client, _ := NewClient()
	builder := NewRecordBuilder(client, "TestDB", "TestLayout")

	builder.WithDeleteRelated("MyPortal")

	if builder.deleteRelated != "MyPortal" {
		t.Error("expected deleteRelated to be MyPortal")
	}
}

func TestRecordBuilder_Create(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			_ = json.NewEncoder(w).Encode(ResponseData{
				Response: Response{Token: "test-token"},
				Messages: []Message{{Code: "0", Message: "OK"}},
			})
			return
		}
		if callCount == 2 {
			_ = json.NewEncoder(w).Encode(ResponseData{
				Response: Response{RecordID: "123"},
				Messages: []Message{{Code: "0", Message: "OK"}},
			})
			return
		}
		_ = json.NewEncoder(w).Encode(ResponseData{
			Messages: []Message{{Code: "0", Message: "OK"}},
		})
	}))
	defer server.Close()

	client, _ := NewClient(
		SetURL(server.URL),
		SetBasicAuth("user", "pass"),
	)
	builder := NewRecordBuilder(client, "TestDB", "TestLayout")
	builder.SetField("Name", "Test")

	resp, err := builder.Create(context.Background())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("response is nil")
	}
}

func TestRecordBuilder_Update(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			_ = json.NewEncoder(w).Encode(ResponseData{
				Response: Response{Token: "test-token"},
				Messages: []Message{{Code: "0", Message: "OK"}},
			})
			return
		}
		if callCount == 2 {
			_ = json.NewEncoder(w).Encode(ResponseData{
				Response: Response{ModID: "2"},
				Messages: []Message{{Code: "0", Message: "OK"}},
			})
			return
		}
		_ = json.NewEncoder(w).Encode(ResponseData{
			Messages: []Message{{Code: "0", Message: "OK"}},
		})
	}))
	defer server.Close()

	client, _ := NewClient(
		SetURL(server.URL),
		SetBasicAuth("user", "pass"),
	)
	builder := NewRecordBuilder(client, "TestDB", "TestLayout")
	builder.ForRecord("123").SetField("Name", "Updated")

	resp, err := builder.Update(context.Background())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("response is nil")
	}
}

func TestRecordBuilder_Update_NoRecordID(t *testing.T) {
	client, _ := NewClient()
	builder := NewRecordBuilder(client, "TestDB", "TestLayout")

	_, err := builder.Update(context.Background())
	if err == nil {
		t.Error("expected error for missing recordID")
	}
	if _, ok := err.(*ValidationError); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestRecordBuilder_Delete(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			_ = json.NewEncoder(w).Encode(ResponseData{
				Response: Response{Token: "test-token"},
				Messages: []Message{{Code: "0", Message: "OK"}},
			})
			return
		}
		_ = json.NewEncoder(w).Encode(ResponseData{
			Messages: []Message{{Code: "0", Message: "OK"}},
		})
	}))
	defer server.Close()

	client, _ := NewClient(
		SetURL(server.URL),
		SetBasicAuth("user", "pass"),
	)
	builder := NewRecordBuilder(client, "TestDB", "TestLayout")
	builder.ForRecord("123")
	resp, err := builder.Delete(context.Background())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("response is nil")
	}
}

func TestRecordBuilder_Update_WithModID(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			_ = json.NewEncoder(w).Encode(ResponseData{
				Response: Response{Token: "test-token"},
				Messages: []Message{{Code: "0", Message: "OK"}},
			})
			return
		}
		if callCount == 2 {
			var payload Payload
			err := json.NewDecoder(r.Body).Decode(&payload)
			if err != nil {
				t.Fatalf("failed to decode payload: %v", err)
			}
			if payload.ModId != "5" {
				t.Errorf("expected modId 5, got %s", payload.ModId)
			}

			_ = json.NewEncoder(w).Encode(ResponseData{
				Response: Response{ModID: "6"},
				Messages: []Message{{Code: "0", Message: "OK"}},
			})
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, _ := NewClient(SetURL(server.URL), SetBasicAuth("u", "p"))
	builder := NewRecordBuilder(client, "DB", "Layout")
	builder.ForRecord("123").WithModID("5").SetField("f", "v")

	_, err := builder.Update(context.Background())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRecordBuilder_Delete_WithDeleteRelated(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
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
		}
		_ = json.NewEncoder(w).Encode(ResponseData{
			Messages: []Message{{Code: "0", Message: "OK"}},
		})
	}))
	defer server.Close()

	client, _ := NewClient(SetURL(server.URL), SetBasicAuth("u", "p"))
	builder := NewRecordBuilder(client, "DB", "Layout")
	builder.ForRecord("123").WithDeleteRelated("MyPortal")

	_, err := builder.Delete(context.Background())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRecordBuilder_Delete_NoRecordID(t *testing.T) {
	client, _ := NewClient()
	builder := NewRecordBuilder(client, "TestDB", "TestLayout")

	_, err := builder.Delete(context.Background())
	if err == nil {
		t.Error("expected error for missing recordID")
	}
}

func TestRecordBuilder_Get(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			_ = json.NewEncoder(w).Encode(ResponseData{
				Response: Response{Token: "test-token"},
				Messages: []Message{{Code: "0", Message: "OK"}},
			})
			return
		}
		_ = json.NewEncoder(w).Encode(ResponseData{
			Response: Response{
				RecordID: "123",
				Data:     []Datum{{RecordID: "123"}},
			},
			Messages: []Message{{Code: "0", Message: "OK"}},
		})
	}))
	defer server.Close()

	client, _ := NewClient(
		SetURL(server.URL),
		SetBasicAuth("user", "pass"),
	)
	builder := NewRecordBuilder(client, "TestDB", "TestLayout")
	builder.ForRecord("123")

	resp, err := builder.Get(context.Background())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("response is nil")
	}
}

func TestRecordBuilder_Get_NoRecordID(t *testing.T) {
	client, _ := NewClient()
	builder := NewRecordBuilder(client, "TestDB", "TestLayout")

	_, err := builder.Get(context.Background())
	if err == nil {
		t.Error("expected error for missing recordID")
	}
}

func TestRecordBuilder_Duplicate(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			_ = json.NewEncoder(w).Encode(ResponseData{
				Response: Response{Token: "test-token"},
				Messages: []Message{{Code: "0", Message: "OK"}},
			})
			return
		}
		_ = json.NewEncoder(w).Encode(ResponseData{
			Response: Response{RecordID: "456"},
			Messages: []Message{{Code: "0", Message: "OK"}},
		})
	}))
	defer server.Close()

	client, _ := NewClient(
		SetURL(server.URL),
		SetBasicAuth("user", "pass"),
	)
	builder := NewRecordBuilder(client, "TestDB", "TestLayout")
	builder.ForRecord("123")

	resp, err := builder.Duplicate(context.Background())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("response is nil")
	}
}

func TestRecordBuilder_Duplicate_NoRecordID(t *testing.T) {
	client, _ := NewClient()
	builder := NewRecordBuilder(client, "TestDB", "TestLayout")

	_, err := builder.Duplicate(context.Background())
	if err == nil {
		t.Error("expected error for missing recordID")
	}
}

func TestNewFindBuilder(t *testing.T) {
	client, _ := NewClient()
	builder := NewFindBuilder(client, "TestDB", "TestLayout")

	if builder == nil {
		t.Fatal("builder is nil")
	}
	if builder.database != "TestDB" {
		t.Errorf("expected database TestDB, got %s", builder.database)
	}
	if builder.layout != "TestLayout" {
		t.Errorf("expected layout TestLayout, got %s", builder.layout)
	}
}

func TestFindBuilder_Where(t *testing.T) {
	client, _ := NewClient()
	builder := NewFindBuilder(client, "TestDB", "TestLayout")

	builder.Where("FirstName", Equal, "John")
	builder.Where("LastName", Equal, "Doe")

	if len(builder.queries) != 1 {
		t.Errorf("expected 1 query group, got %d", len(builder.queries))
	}
	if len(builder.queries[0].queries) != 2 {
		t.Errorf("expected 2 query fields, got %d", len(builder.queries[0].queries))
	}
}

func TestFindBuilder_OrWhere(t *testing.T) {
	client, _ := NewClient()
	builder := NewFindBuilder(client, "TestDB", "TestLayout")

	builder.Where("Status", Equal, "Active")
	builder.OrWhere("Status", Equal, "Pending")

	if len(builder.queries) != 2 {
		t.Errorf("expected 2 query groups, got %d", len(builder.queries))
	}
}

func TestFindBuilder_OrderBy(t *testing.T) {
	client, _ := NewClient()
	builder := NewFindBuilder(client, "TestDB", "TestLayout")

	builder.OrderBy("LastName", Ascending)
	builder.OrderBy("FirstName", Ascending)

	if len(builder.sorters) != 2 {
		t.Errorf("expected 2 sorters, got %d", len(builder.sorters))
	}
}

func TestFindBuilder_Offset(t *testing.T) {
	client, _ := NewClient()
	builder := NewFindBuilder(client, "TestDB", "TestLayout")

	builder.Offset(10)

	if builder.offset != "10" {
		t.Errorf("expected offset 10, got %s", builder.offset)
	}
}

func TestFindBuilder_Limit(t *testing.T) {
	client, _ := NewClient()
	builder := NewFindBuilder(client, "TestDB", "TestLayout")

	builder.Limit(50)

	if builder.limit != "50" {
		t.Errorf("expected limit 50, got %s", builder.limit)
	}
}

func TestFindBuilder_Execute(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			_ = json.NewEncoder(w).Encode(ResponseData{
				Response: Response{Token: "test-token"},
				Messages: []Message{{Code: "0", Message: "OK"}},
			})
			return
		}
		if callCount == 2 {
			_ = json.NewEncoder(w).Encode(ResponseData{
				Response: Response{
					Data: []Datum{{RecordID: "1"}, {RecordID: "2"}},
				},
				Messages: []Message{{Code: "0", Message: "OK"}},
			})
			return
		}
		_ = json.NewEncoder(w).Encode(ResponseData{
			Messages: []Message{{Code: "0", Message: "OK"}},
		})
	}))
	defer server.Close()

	client, _ := NewClient(
		SetURL(server.URL),
		SetBasicAuth("user", "pass"),
	)
	builder := NewFindBuilder(client, "TestDB", "TestLayout")
	builder.Where("Status", Equal, "Active").Limit(10)

	resp, err := builder.Execute(context.Background())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("response is nil")
	}
}

func TestNewSessionBuilder(t *testing.T) {
	client, _ := NewClient()
	builder := NewSessionBuilder(client, "TestDB")

	if builder == nil {
		t.Fatal("builder is nil")
	}
	if builder.database != "TestDB" {
		t.Errorf("expected database TestDB, got %s", builder.database)
	}
}

func TestSessionBuilder_Connect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(ResponseData{
			Response: Response{Token: "new-token"},
			Messages: []Message{{Code: "0", Message: "OK"}},
		})
	}))
	defer server.Close()

	client, _ := NewClient(
		SetURL(server.URL),
		SetBasicAuth("user", "pass"),
	)
	builder := NewSessionBuilder(client, "TestDB")

	resp, err := builder.Connect(context.Background())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("response is nil")
	}
	if builder.token != "new-token" {
		t.Errorf("expected token new-token, got %s", builder.token)
	}
}

func TestSessionBuilder_ConnectWithDatasource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(ResponseData{
			Response: Response{Token: "datasource-token"},
			Messages: []Message{{Code: "0", Message: "OK"}},
		})
	}))
	defer server.Close()

	client, _ := NewClient(
		SetURL(server.URL),
		SetBasicAuth("user", "pass"),
	)
	builder := NewSessionBuilder(client, "TestDB")

	resp, err := builder.ConnectWithDatasource(context.Background())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("response is nil")
	}
	if builder.token != "datasource-token" {
		t.Errorf("expected token datasource-token, got %s", builder.token)
	}
}

func TestSessionBuilder_WithToken(t *testing.T) {
	client, _ := NewClient()
	builder := NewSessionBuilder(client, "TestDB")

	builder.WithToken("custom-token")

	if builder.token != "custom-token" {
		t.Errorf("expected token custom-token, got %s", builder.token)
	}
}

func TestSessionBuilder_Validate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(ResponseData{
			Messages: []Message{{Code: "0", Message: "OK"}},
		})
	}))
	defer server.Close()

	client, _ := NewClient(SetURL(server.URL))
	builder := NewSessionBuilder(client, "TestDB")
	builder.WithToken("valid-token")

	resp, err := builder.Validate(context.Background())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("response is nil")
	}
}

func TestSessionBuilder_Validate_NoToken(t *testing.T) {
	client, _ := NewClient()
	builder := NewSessionBuilder(client, "TestDB")

	_, err := builder.Validate(context.Background())
	if err == nil {
		t.Error("expected error for missing token")
	}
	if _, ok := err.(*ValidationError); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestSessionBuilder_Disconnect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(ResponseData{
			Messages: []Message{{Code: "0", Message: "OK"}},
		})
	}))
	defer server.Close()

	client, _ := NewClient(SetURL(server.URL))
	builder := NewSessionBuilder(client, "TestDB")
	builder.WithToken("active-token")

	resp, err := builder.Disconnect(context.Background())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("response is nil")
	}
}

func TestSessionBuilder_Disconnect_NoToken(t *testing.T) {
	client, _ := NewClient()
	builder := NewSessionBuilder(client, "TestDB")

	_, err := builder.Disconnect(context.Background())
	if err == nil {
		t.Error("expected error for missing token")
	}
}

func TestSessionBuilder_Token(t *testing.T) {
	client, _ := NewClient()
	builder := NewSessionBuilder(client, "TestDB")
	builder.token = "test-token"

	token := builder.Token()
	if token != "test-token" {
		t.Errorf("expected token test-token, got %s", token)
	}
}
