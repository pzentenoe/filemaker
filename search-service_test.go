package filemaker

import (
	"encoding/json"
	"testing"
)

func Test_searchService_GroupQuery(t *testing.T) {
	t.Run("Test query group", func(t *testing.T) {
		search := NewSearchService("test", "test_layout", nil)
		search.GroupQuery(
			NewQueryGroup(
				NewQueryFieldOperator("nombre", "pablo", Equal),
				NewQueryFieldOperator("apellido", "zenteno", Equal),
			),
		)
		b, _ := json.Marshal(search.seachData.QueryGroup)
		if string(b) != "[{\"apellido\":\"==zenteno\",\"nombre\":\"==pablo\"}]" {
			t.Errorf("JSON was incorrect, got: %s, want: %s", string(b), "[{\"apellido\":\"==zenteno\",\"nombre\":\"==pablo\"}]")
		}
	})

	t.Run("Test query group with 1 query", func(t *testing.T) {
		search := NewSearchService("test", "test_layout", nil)
		search.GroupQuery(
			NewQueryGroup(
				NewQueryFieldOperator("nombre", "pablo", Equal),
			),
		)
		b, _ := json.Marshal(search.seachData.QueryGroup)

		if string(b) != "[{\"nombre\":\"==pablo\"}]" {
			t.Errorf("JSON was incorrect, got: %s, want: %s", string(b), "[{\"nombre\":\"==pablo\"}]")
		}
	})
}
