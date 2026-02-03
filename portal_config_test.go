package filemaker

import "testing"

func TestNewPortalConfig(t *testing.T) {
	config := NewPortalConfig("TestPortal")

	if config.Name != "TestPortal" {
		t.Errorf("expected Name to be 'TestPortal', got '%s'", config.Name)
	}

	if config.Offset != 1 {
		t.Errorf("expected default Offset to be 1, got %d", config.Offset)
	}

	if config.Limit != 50 {
		t.Errorf("expected default Limit to be 50, got %d", config.Limit)
	}
}

func TestPortalConfig_WithOffset(t *testing.T) {
	config := NewPortalConfig("TestPortal").WithOffset(10)

	if config.Offset != 10 {
		t.Errorf("expected Offset to be 10, got %d", config.Offset)
	}
}

func TestPortalConfig_WithLimit(t *testing.T) {
	config := NewPortalConfig("TestPortal").WithLimit(25)

	if config.Limit != 25 {
		t.Errorf("expected Limit to be 25, got %d", config.Limit)
	}
}

func TestPortalConfig_Chaining(t *testing.T) {
	config := NewPortalConfig("Orders").
		WithOffset(5).
		WithLimit(20)

	if config.Name != "Orders" {
		t.Errorf("expected Name to be 'Orders', got '%s'", config.Name)
	}

	if config.Offset != 5 {
		t.Errorf("expected Offset to be 5, got %d", config.Offset)
	}

	if config.Limit != 20 {
		t.Errorf("expected Limit to be 20, got %d", config.Limit)
	}
}

func TestSearchService_SetPortalConfigs(t *testing.T) {
	client, _ := NewClient(SetURL("https://test.com"))
	service := NewSearchService("TestDB", "TestLayout", client)

	config1 := NewPortalConfig("Portal1").WithOffset(1).WithLimit(10)
	config2 := NewPortalConfig("Portal2").WithOffset(5).WithLimit(20)

	service.SetPortalConfigs(config1, config2)

	if len(service.portalConfigs) != 2 {
		t.Errorf("expected 2 portal configs, got %d", len(service.portalConfigs))
	}

	if len(service.searchData.Portal) != 2 {
		t.Errorf("expected 2 portal names in searchData, got %d", len(service.searchData.Portal))
	}

	if service.searchData.Portal[0] != "Portal1" {
		t.Errorf("expected first portal to be 'Portal1', got '%s'", service.searchData.Portal[0])
	}

	if service.searchData.Portal[1] != "Portal2" {
		t.Errorf("expected second portal to be 'Portal2', got '%s'", service.searchData.Portal[1])
	}
}

func TestSearchService_SetPortalConfigs_ExtractsNames(t *testing.T) {
	client, _ := NewClient(SetURL("https://test.com"))
	service := NewSearchService("TestDB", "TestLayout", client)

	// Configure portals with pagination
	service.SetPortalConfigs(
		NewPortalConfig("RelatedOrders"),
		NewPortalConfig("RelatedPayments"),
	)

	// Verify portal names are extracted correctly
	expectedNames := []string{"RelatedOrders", "RelatedPayments"}
	if len(service.searchData.Portal) != len(expectedNames) {
		t.Fatalf("expected %d portal names, got %d", len(expectedNames), len(service.searchData.Portal))
	}

	for i, name := range expectedNames {
		if service.searchData.Portal[i] != name {
			t.Errorf("expected portal[%d] to be '%s', got '%s'", i, name, service.searchData.Portal[i])
		}
	}
}
