package filemaker

import (
	"testing"
)

func TestNewQueryFieldOperator(t *testing.T) {
	qf := NewQueryFieldOperator("Name", "John", Equal)

	if qf == nil {
		t.Fatal("NewQueryFieldOperator() returned nil")
	}

	if qf.Name != "Name" {
		t.Errorf("Name = %v, want Name", qf.Name)
	}

	if qf.Value != "John" {
		t.Errorf("Value = %v, want John", qf.Value)
	}

	if qf.Operator != Equal {
		t.Errorf("Operator = %v, want %v", qf.Operator, Equal)
	}
}

func TestQueryFieldOperator_valueWithOp(t *testing.T) {
	tests := []struct {
		name     string
		operator FieldOperator
		value    string
		want     string
	}{
		{
			name:     "equal",
			operator: Equal,
			value:    "John",
			want:     "==John",
		},
		{
			name:     "contains",
			operator: Contains,
			value:    "John",
			want:     "==*John*",
		},
		{
			name:     "begins with",
			operator: BeginsWith,
			value:    "John",
			want:     "==John*",
		},
		{
			name:     "ends with",
			operator: EndsWith,
			value:    "John",
			want:     "==*John",
		},
		{
			name:     "greater than",
			operator: GreaterThan,
			value:    "100",
			want:     ">100",
		},
		{
			name:     "greater than equal",
			operator: GreaterThanEqual,
			value:    "100",
			want:     ">=100",
		},
		{
			name:     "less than",
			operator: LessThan,
			value:    "100",
			want:     "<100",
		},
		{
			name:     "less than equal",
			operator: LessThanEqual,
			value:    "100",
			want:     "<=100",
		},
		{
			name:     "unknown operator",
			operator: FieldOperator("unknown"),
			value:    "test",
			want:     "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qf := NewQueryFieldOperator("field", tt.value, tt.operator)
			if got := qf.valueWithOp(); got != tt.want {
				t.Errorf("valueWithOp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldOperator_Constants(t *testing.T) {
	tests := []struct {
		name     string
		operator FieldOperator
		want     string
	}{
		{"equal", Equal, "eq"},
		{"contains", Contains, "cn"},
		{"begins with", BeginsWith, "bw"},
		{"ends with", EndsWith, "ew"},
		{"greater than", GreaterThan, "gt"},
		{"greater than equal", GreaterThanEqual, "gte"},
		{"less than", LessThan, "lt"},
		{"less than equal", LessThanEqual, "lte"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.operator) != tt.want {
				t.Errorf("operator = %v, want %v", tt.operator, tt.want)
			}
		})
	}
}
