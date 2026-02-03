package filemaker

import (
	"testing"
)

func TestValidateRequired(t *testing.T) {
	tests := []struct {
		name    string
		field   string
		value   string
		message string
		wantErr bool
	}{
		{
			name:    "valid value",
			field:   "testField",
			value:   "test",
			message: "test field is required",
			wantErr: false,
		},
		{
			name:    "empty value",
			field:   "testField",
			value:   "",
			message: "test field is required",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRequired(tt.field, tt.value, tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateRequired() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				valErr, ok := err.(*ValidationError)
				if !ok {
					t.Errorf("validateRequired() error type = %T, want *ValidationError", err)
				}
				if valErr.Field != tt.field {
					t.Errorf("ValidationError.Field = %v, want %v", valErr.Field, tt.field)
				}
			}
		})
	}
}

func TestValidateDatabase(t *testing.T) {
	tests := []struct {
		name     string
		database string
		wantErr  bool
	}{
		{
			name:     "valid database",
			database: "MyDatabase",
			wantErr:  false,
		},
		{
			name:     "empty database",
			database: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDatabase(tt.database)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateDatabase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateLayout(t *testing.T) {
	tests := []struct {
		name    string
		layout  string
		wantErr bool
	}{
		{
			name:    "valid layout",
			layout:  "MyLayout",
			wantErr: false,
		},
		{
			name:    "empty layout",
			layout:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLayout(tt.layout)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateLayout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   "abc123",
			wantErr: false,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateRepetition(t *testing.T) {
	tests := []struct {
		name       string
		repetition int
		wantErr    bool
	}{
		{
			name:       "valid repetition",
			repetition: 1,
			wantErr:    false,
		},
		{
			name:       "valid repetition greater than 1",
			repetition: 5,
			wantErr:    false,
		},
		{
			name:       "invalid repetition zero",
			repetition: 0,
			wantErr:    true,
		},
		{
			name:       "invalid repetition negative",
			repetition: -1,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRepetition(tt.repetition)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateRepetition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGlobalFields(t *testing.T) {
	tests := []struct {
		name    string
		fields  map[string]any
		wantErr bool
	}{
		{
			name:    "valid fields",
			fields:  map[string]any{"field1": "value1"},
			wantErr: false,
		},
		{
			name:    "empty fields",
			fields:  map[string]any{},
			wantErr: true,
		},
		{
			name:    "nil fields",
			fields:  nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateGlobalFields(tt.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateGlobalFields() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateFileData(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "valid data",
			data:    []byte("test data"),
			wantErr: false,
		},
		{
			name:    "empty data",
			data:    []byte{},
			wantErr: true,
		},
		{
			name:    "nil data",
			data:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFileData(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFileData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateAction(t *testing.T) {
	tests := []struct {
		name    string
		action  string
		wantErr bool
	}{
		{
			name:    "valid action create",
			action:  "create",
			wantErr: false,
		},
		{
			name:    "valid action edit",
			action:  "edit",
			wantErr: false,
		},
		{
			name:    "valid action delete",
			action:  "delete",
			wantErr: false,
		},
		{
			name:    "invalid action",
			action:  "invalid",
			wantErr: true,
		},
		{
			name:    "empty action",
			action:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAction(tt.action)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
