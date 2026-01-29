package filemaker

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestLogLevel_String(t *testing.T) {
	tests := []struct {
		level LogLevel
		want  string
	}{
		{LogLevelDebug, "DEBUG"},
		{LogLevelInfo, "INFO"},
		{LogLevelWarn, "WARN"},
		{LogLevelError, "ERROR"},
		{LogLevel(99), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.level.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogLevel_toSlogLevel(t *testing.T) {
	tests := []struct {
		level LogLevel
		want  slog.Level
	}{
		{LogLevelDebug, slog.LevelDebug},
		{LogLevelInfo, slog.LevelInfo},
		{LogLevelWarn, slog.LevelWarn},
		{LogLevelError, slog.LevelError},
		{LogLevel(99), slog.LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			if got := tt.level.toSlogLevel(); got != tt.want {
				t.Errorf("toSlogLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoOpLogger(t *testing.T) {
	logger := &NoOpLogger{}
	fields := map[string]any{"key": "value"}

	logger.Debug("test", fields)
	logger.Info("test", fields)
	logger.Warn("test", fields)
	logger.Error("test", fields)

	newLogger := logger.With(fields)
	if newLogger != logger {
		t.Error("With() should return same instance")
	}
}

func TestNewSlogLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewTextHandler(buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := NewSlogLogger(LogLevelInfo, handler)
	if logger == nil {
		t.Fatal("NewSlogLogger() returned nil")
	}

	if logger.level != LogLevelInfo {
		t.Errorf("level = %v, want %v", logger.level, LogLevelInfo)
	}

	logger.Info("test message", map[string]any{"key": "value"})

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("output should contain message, got: %s", output)
	}
}

func TestNewSlogLogger_NilHandler(t *testing.T) {
	logger := NewSlogLogger(LogLevelInfo, nil)
	if logger == nil {
		t.Fatal("NewSlogLogger() with nil handler returned nil")
	}

	logger.Info("test", nil)
}

func TestNewDefaultLogger(t *testing.T) {
	logger := NewDefaultLogger(LogLevelDebug)
	if logger == nil {
		t.Fatal("NewDefaultLogger() returned nil")
	}

	if logger.level != LogLevelDebug {
		t.Errorf("level = %v, want %v", logger.level, LogLevelDebug)
	}
}

func TestNewJSONLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewJSONHandler(buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := NewSlogLogger(LogLevelInfo, handler)
	if logger == nil {
		t.Fatal("NewJSONLogger() returned nil")
	}

	logger.Info("test message", map[string]any{"key": "value"})

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("output should contain message, got: %s", output)
	}
	if !strings.Contains(output, `"key":"value"`) {
		t.Errorf("output should be JSON with key-value, got: %s", output)
	}
}

func TestSlogLogger_Debug(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewTextHandler(buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := NewSlogLogger(LogLevelDebug, handler)
	logger.Debug("debug message", map[string]any{"field": "value"})

	output := buf.String()
	if !strings.Contains(output, "debug message") {
		t.Errorf("output should contain debug message, got: %s", output)
	}
}

func TestSlogLogger_Debug_BelowLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewTextHandler(buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := NewSlogLogger(LogLevelInfo, handler)
	logger.Debug("debug message", map[string]any{"field": "value"})

	output := buf.String()
	if strings.Contains(output, "debug message") {
		t.Errorf("debug message should not be logged at info level")
	}
}

func TestSlogLogger_Info(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewTextHandler(buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := NewSlogLogger(LogLevelInfo, handler)
	logger.Info("info message", map[string]any{"field": "value"})

	output := buf.String()
	if !strings.Contains(output, "info message") {
		t.Errorf("output should contain info message, got: %s", output)
	}
}

func TestSlogLogger_Warn(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewTextHandler(buf, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})

	logger := NewSlogLogger(LogLevelWarn, handler)
	logger.Warn("warn message", map[string]any{"field": "value"})

	output := buf.String()
	if !strings.Contains(output, "warn message") {
		t.Errorf("output should contain warn message, got: %s", output)
	}
}

func TestSlogLogger_Error(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewTextHandler(buf, &slog.HandlerOptions{
		Level: slog.LevelError,
	})

	logger := NewSlogLogger(LogLevelError, handler)
	logger.Error("error message", map[string]any{"field": "value"})

	output := buf.String()
	if !strings.Contains(output, "error message") {
		t.Errorf("output should contain error message, got: %s", output)
	}
}

func TestSlogLogger_With(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewTextHandler(buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := NewSlogLogger(LogLevelInfo, handler)
	contextLogger := logger.With(map[string]any{"request_id": "123"})

	contextLogger.Info("test message", nil)

	output := buf.String()
	if !strings.Contains(output, "request_id") {
		t.Errorf("output should contain contextual field, got: %s", output)
	}
}

func TestSlogLogger_With_EmptyFields(t *testing.T) {
	logger := NewDefaultLogger(LogLevelInfo)
	same := logger.With(map[string]any{})

	if same != logger {
		t.Error("With() with empty fields should return same logger")
	}
}

func TestFieldsToAttrs(t *testing.T) {
	tests := []struct {
		name   string
		fields map[string]any
		want   int
	}{
		{
			name:   "nil fields",
			fields: nil,
			want:   0,
		},
		{
			name:   "empty fields",
			fields: map[string]any{},
			want:   0,
		},
		{
			name:   "single field",
			fields: map[string]any{"key": "value"},
			want:   2,
		},
		{
			name:   "multiple fields",
			fields: map[string]any{"key1": "value1", "key2": "value2"},
			want:   4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs := fieldsToAttrs(tt.fields)
			if len(attrs) != tt.want {
				t.Errorf("fieldsToAttrs() length = %v, want %v", len(attrs), tt.want)
			}
		})
	}
}

func TestNewStandardLogger(t *testing.T) {
	logger := NewStandardLogger(LogLevelInfo)
	if logger == nil {
		t.Fatal("NewStandardLogger() returned nil")
	}

	if logger.SlogLogger == nil {
		t.Error("StandardLogger.SlogLogger is nil")
	}
}

func TestMetrics_NewMetrics(t *testing.T) {
	m := NewMetrics()
	if m == nil {
		t.Fatal("NewMetrics() returned nil")
	}

	if m.GetRequestsTotal() != 0 {
		t.Error("initial requests total should be 0")
	}
}

func TestMetrics_IncrementRequests(t *testing.T) {
	m := NewMetrics()
	m.IncrementRequests()
	m.IncrementRequests()

	if got := m.GetRequestsTotal(); got != 2 {
		t.Errorf("GetRequestsTotal() = %v, want 2", got)
	}
}

func TestMetrics_IncrementRequestsSucceeded(t *testing.T) {
	m := NewMetrics()
	m.IncrementRequestsSucceeded()

	if got := m.GetRequestsSucceeded(); got != 1 {
		t.Errorf("GetRequestsSucceeded() = %v, want 1", got)
	}
}

func TestMetrics_IncrementRequestsFailed(t *testing.T) {
	m := NewMetrics()
	m.IncrementRequestsFailed()

	if got := m.GetRequestsFailed(); got != 1 {
		t.Errorf("GetRequestsFailed() = %v, want 1", got)
	}
}

func TestMetrics_IncrementRetries(t *testing.T) {
	m := NewMetrics()
	m.IncrementRetries()

	if got := m.GetRetriesTotal(); got != 1 {
		t.Errorf("GetRetriesTotal() = %v, want 1", got)
	}
}

func TestMetrics_IncrementSessionsCreated(t *testing.T) {
	m := NewMetrics()
	m.IncrementSessionsCreated()

	if got := m.GetSessionsCreated(); got != 1 {
		t.Errorf("GetSessionsCreated() = %v, want 1", got)
	}
}

func TestMetrics_IncrementSessionsClosed(t *testing.T) {
	m := NewMetrics()
	m.IncrementSessionsClosed()

	if got := m.GetSessionsClosed(); got != 1 {
		t.Errorf("GetSessionsClosed() = %v, want 1", got)
	}
}

func TestMetrics_Reset(t *testing.T) {
	m := NewMetrics()
	m.IncrementRequests()
	m.IncrementRequestsSucceeded()
	m.IncrementRequestsFailed()
	m.IncrementRetries()
	m.IncrementSessionsCreated()
	m.IncrementSessionsClosed()

	m.Reset()

	if m.GetRequestsTotal() != 0 ||
		m.GetRequestsSucceeded() != 0 ||
		m.GetRequestsFailed() != 0 ||
		m.GetRetriesTotal() != 0 ||
		m.GetSessionsCreated() != 0 ||
		m.GetSessionsClosed() != 0 {
		t.Error("Reset() should reset all metrics to 0")
	}
}

func TestMetrics_Snapshot(t *testing.T) {
	m := NewMetrics()
	m.IncrementRequests()
	m.IncrementRequestsSucceeded()
	m.IncrementSessionsCreated()
	m.IncrementSessionsCreated()

	snapshot := m.Snapshot()

	if snapshot["requests_total"] != 1 {
		t.Errorf("snapshot requests_total = %v, want 1", snapshot["requests_total"])
	}

	if snapshot["requests_succeeded"] != 1 {
		t.Errorf("snapshot requests_succeeded = %v, want 1", snapshot["requests_succeeded"])
	}

	if snapshot["active_sessions"] != 2 {
		t.Errorf("snapshot active_sessions = %v, want 2", snapshot["active_sessions"])
	}
}

func TestNewRequestContext(t *testing.T) {
	rc := NewRequestContext("Create", "TestDB", "TestLayout")

	if rc.Operation != "Create" {
		t.Errorf("Operation = %v, want Create", rc.Operation)
	}

	if rc.Database != "TestDB" {
		t.Errorf("Database = %v, want TestDB", rc.Database)
	}

	if rc.Layout != "TestLayout" {
		t.Errorf("Layout = %v, want TestLayout", rc.Layout)
	}

	if rc.RequestID == "" {
		t.Error("RequestID should not be empty")
	}

	if rc.StartTime.IsZero() {
		t.Error("StartTime should not be zero")
	}
}

func TestRequestContext_Duration(t *testing.T) {
	rc := NewRequestContext("Test", "DB", "Layout")
	time.Sleep(10 * time.Millisecond)

	duration := rc.Duration()
	if duration < 10*time.Millisecond {
		t.Errorf("Duration() = %v, want >= 10ms", duration)
	}
}

func TestRequestContext_Fields(t *testing.T) {
	rc := NewRequestContext("Create", "TestDB", "TestLayout")

	fields := rc.Fields()

	if fields["operation"] != "Create" {
		t.Errorf("fields[operation] = %v, want Create", fields["operation"])
	}

	if fields["database"] != "TestDB" {
		t.Errorf("fields[database] = %v, want TestDB", fields["database"])
	}

	if fields["layout"] != "TestLayout" {
		t.Errorf("fields[layout] = %v, want TestLayout", fields["layout"])
	}

	if fields["request_id"] == "" {
		t.Error("fields[request_id] should not be empty")
	}
}

func TestRequestContext_SlogAttrs(t *testing.T) {
	rc := NewRequestContext("Create", "TestDB", "TestLayout")

	attrs := rc.SlogAttrs()

	if len(attrs) != 5 {
		t.Errorf("SlogAttrs() length = %v, want 5", len(attrs))
	}
}

func TestWithRequestContext(t *testing.T) {
	ctx := context.Background()
	rc := NewRequestContext("Test", "DB", "Layout")

	newCtx := WithRequestContext(ctx, rc)

	retrieved, ok := GetRequestContext(newCtx)
	if !ok {
		t.Fatal("GetRequestContext() failed to retrieve context")
	}

	if retrieved != rc {
		t.Error("Retrieved context does not match original")
	}
}

func TestGetRequestContext_NotFound(t *testing.T) {
	ctx := context.Background()

	_, ok := GetRequestContext(ctx)
	if ok {
		t.Error("GetRequestContext() should return false for context without RequestContext")
	}
}

func TestWithLogger(t *testing.T) {
	ctx := context.Background()
	logger := NewDefaultLogger(LogLevelInfo)

	newCtx := WithLogger(ctx, logger)

	retrieved := GetLogger(newCtx)
	if retrieved != logger {
		t.Error("Retrieved logger does not match original")
	}
}

func TestGetLogger_NotFound(t *testing.T) {
	ctx := context.Background()

	logger := GetLogger(ctx)

	if _, ok := logger.(*NoOpLogger); !ok {
		t.Error("GetLogger() should return NoOpLogger when not found")
	}
}

func TestLogFromContext(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewTextHandler(buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := NewSlogLogger(LogLevelDebug, handler)
	ctx := WithLogger(context.Background(), logger)

	tests := []struct {
		level LogLevel
		msg   string
	}{
		{LogLevelDebug, "debug message"},
		{LogLevelInfo, "info message"},
		{LogLevelWarn, "warn message"},
		{LogLevelError, "error message"},
	}

	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			buf.Reset()
			LogFromContext(ctx, tt.level, tt.msg, map[string]any{"key": "value"})

			output := buf.String()
			if !strings.Contains(output, tt.msg) {
				t.Errorf("output should contain %s, got: %s", tt.msg, output)
			}
		})
	}
}

func TestLogFromContext_NoLogger(t *testing.T) {
	ctx := context.Background()
	LogFromContext(ctx, LogLevelInfo, "test", nil)
}
