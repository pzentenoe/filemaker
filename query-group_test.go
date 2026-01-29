package filemaker

import (
	"testing"
)

func TestNewGroupQuery(t *testing.T) {
	tests := []struct {
		name    string
		queries []*queryFieldOperator
		want    int
	}{
		{
			name:    "empty group",
			queries: nil,
			want:    0,
		},
		{
			name: "single query",
			queries: []*queryFieldOperator{
				NewQueryFieldOperator("Name", "John", Equal),
			},
			want: 1,
		},
		{
			name: "multiple queries",
			queries: []*queryFieldOperator{
				NewQueryFieldOperator("Name", "John", Equal),
				NewQueryFieldOperator("Age", "30", GreaterThan),
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := NewGroupQuery(tt.queries...)

			if group == nil {
				t.Fatal("NewGroupQuery() returned nil")
			}

			if len(group.queries) != tt.want {
				t.Errorf("queries length = %v, want %v", len(group.queries), tt.want)
			}
		})
	}
}

func TestGroupQuery_AddQuery(t *testing.T) {
	group := NewGroupQuery()

	if len(group.queries) != 0 {
		t.Errorf("initial queries length = %v, want 0", len(group.queries))
	}

	result := group.AddQuery("Name", Equal, "John")

	if result != group {
		t.Error("AddQuery() should return same instance for chaining")
	}

	if len(group.queries) != 1 {
		t.Errorf("queries length after add = %v, want 1", len(group.queries))
	}

	if group.queries[0].Name != "Name" {
		t.Errorf("query Name = %v, want Name", group.queries[0].Name)
	}

	if group.queries[0].Value != "John" {
		t.Errorf("query Value = %v, want John", group.queries[0].Value)
	}

	if group.queries[0].Operator != Equal {
		t.Errorf("query Operator = %v, want %v", group.queries[0].Operator, Equal)
	}
}

func TestGroupQuery_AddQuery_Chaining(t *testing.T) {
	group := NewGroupQuery().
		AddQuery("Name", Equal, "John").
		AddQuery("Age", GreaterThan, "30").
		AddQuery("Status", Contains, "Active")

	if len(group.queries) != 3 {
		t.Errorf("queries length = %v, want 3", len(group.queries))
	}

	expectedQueries := []struct {
		name     string
		value    string
		operator FieldOperator
	}{
		{"Name", "John", Equal},
		{"Age", "30", GreaterThan},
		{"Status", "Active", Contains},
	}

	for i, expected := range expectedQueries {
		if group.queries[i].Name != expected.name {
			t.Errorf("query[%d] Name = %v, want %v", i, group.queries[i].Name, expected.name)
		}
		if group.queries[i].Value != expected.value {
			t.Errorf("query[%d] Value = %v, want %v", i, group.queries[i].Value, expected.value)
		}
		if group.queries[i].Operator != expected.operator {
			t.Errorf("query[%d] Operator = %v, want %v", i, group.queries[i].Operator, expected.operator)
		}
	}
}

func TestGroupQuery_AddQuery_DifferentOperators(t *testing.T) {
	operators := []FieldOperator{
		Equal, Contains, BeginsWith, EndsWith,
		GreaterThan, GreaterThanEqual, LessThan, LessThanEqual,
	}

	for _, op := range operators {
		t.Run(string(op), func(t *testing.T) {
			group := NewGroupQuery()
			group.AddQuery("Field", op, "value")

			if len(group.queries) != 1 {
				t.Fatalf("queries length = %v, want 1", len(group.queries))
			}

			if group.queries[0].Operator != op {
				t.Errorf("Operator = %v, want %v", group.queries[0].Operator, op)
			}
		})
	}
}
