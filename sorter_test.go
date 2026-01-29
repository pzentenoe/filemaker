package filemaker

import (
	"encoding/json"
	"testing"
)

func TestNewSorter(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		sortOrder SortOrder
	}{
		{
			name:      "ascending sorter",
			fieldName: "Name",
			sortOrder: Ascending,
		},
		{
			name:      "descending sorter",
			fieldName: "Date",
			sortOrder: Descending,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sorter := NewSorter(tt.fieldName, tt.sortOrder)

			if sorter == nil {
				t.Fatal("NewSorter() returned nil")
			}

			if sorter.FieldName != tt.fieldName {
				t.Errorf("FieldName = %v, want %v", sorter.FieldName, tt.fieldName)
			}

			if sorter.SortOrder != tt.sortOrder {
				t.Errorf("SortOrder = %v, want %v", sorter.SortOrder, tt.sortOrder)
			}
		})
	}
}

func TestSortOrder_String(t *testing.T) {
	tests := []struct {
		name      string
		sortOrder SortOrder
		want      string
	}{
		{
			name:      "ascending",
			sortOrder: Ascending,
			want:      "ascend",
		},
		{
			name:      "descending",
			sortOrder: Descending,
			want:      "descend",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sortOrder.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSorter_JSON(t *testing.T) {
	sorter := NewSorter("TestField", Ascending)

	jsonData, err := json.Marshal(sorter)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	expected := `{"fieldName":"TestField","sortOrder":"ascend"}`
	if string(jsonData) != expected {
		t.Errorf("JSON = %v, want %v", string(jsonData), expected)
	}

	var decoded Sorter
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if decoded.FieldName != "TestField" {
		t.Errorf("Decoded FieldName = %v, want TestField", decoded.FieldName)
	}

	if decoded.SortOrder != Ascending {
		t.Errorf("Decoded SortOrder = %v, want %v", decoded.SortOrder, Ascending)
	}
}
