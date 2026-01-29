package filemaker

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync/atomic"
	"time"
)

// LogLevel represents the severity level of a log message.
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// String returns the string representation of the log level.
func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// toSlogLevel converts LogLevel to slog.Level
func (l LogLevel) toSlogLevel() slog.Level {
	switch l {
	case LogLevelDebug:
		return slog.LevelDebug
	case LogLevelInfo:
		return slog.LevelInfo
	case LogLevelWarn:
		return slog.LevelWarn
	case LogLevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// Logger defines the interface for logging within the FileMaker client.
// Implement this interface to use a custom logger.
type Logger interface {
	// Debug logs a debug message with structured fields.
	Debug(msg string, fields map[string]any)
	// Info logs an informational message with structured fields.
	Info(msg string, fields map[string]any)
	// Warn logs a warning message with structured fields.
	Warn(msg string, fields map[string]any)
	// Error logs an error message with structured fields.
	Error(msg string, fields map[string]any)
	// With returns a new logger with the given fields added to all log entries.
	With(fields map[string]any) Logger
}

// NoOpLogger is a logger that discards all log messages.
type NoOpLogger struct{}

func (n *NoOpLogger) Debug(msg string, fields map[string]any) {}
func (n *NoOpLogger) Info(msg string, fields map[string]any)  {}
func (n *NoOpLogger) Warn(msg string, fields map[string]any)  {}
func (n *NoOpLogger) Error(msg string, fields map[string]any) {}
func (n *NoOpLogger) With(fields map[string]any) Logger       { return n }

// SlogLogger is a logger that uses Go's structured logging package (slog).
// This is the recommended logger for production use.
type SlogLogger struct {
	logger *slog.Logger
	level  LogLevel
}

// NewSlogLogger creates a new SlogLogger with the specified log level and handler.
// If handler is nil, a default TextHandler writing to stdout will be used.
//
// Example with JSON handler:
//
//	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
//	    Level: slog.LevelInfo,
//	    AddSource: true,
//	})
//	logger := filemaker.NewSlogLogger(filemaker.LogLevelInfo, handler)
func NewSlogLogger(level LogLevel, handler slog.Handler) *SlogLogger {
	if handler == nil {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level.toSlogLevel(),
			AddSource: false,
		})
	}

	return &SlogLogger{
		logger: slog.New(handler),
		level:  level,
	}
}

// NewDefaultLogger creates a new SlogLogger with default text output to stdout.
func NewDefaultLogger(level LogLevel) *SlogLogger {
	return NewSlogLogger(level, nil)
}

// NewJSONLogger creates a new SlogLogger with JSON output to stdout.
// This is useful for production environments where logs are aggregated.
func NewJSONLogger(level LogLevel) *SlogLogger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level.toSlogLevel(),
		AddSource: false,
	})
	return NewSlogLogger(level, handler)
}

// fieldsToAttrs converts a map of fields to slog.Attr slice.
func fieldsToAttrs(fields map[string]any) []any {
	if len(fields) == 0 {
		return nil
	}

	attrs := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		attrs = append(attrs, k, v)
	}
	return attrs
}

func (s *SlogLogger) Debug(msg string, fields map[string]any) {
	if s.level <= LogLevelDebug {
		s.logger.Debug(msg, fieldsToAttrs(fields)...)
	}
}

func (s *SlogLogger) Info(msg string, fields map[string]any) {
	if s.level <= LogLevelInfo {
		s.logger.Info(msg, fieldsToAttrs(fields)...)
	}
}

func (s *SlogLogger) Warn(msg string, fields map[string]any) {
	if s.level <= LogLevelWarn {
		s.logger.Warn(msg, fieldsToAttrs(fields)...)
	}
}

func (s *SlogLogger) Error(msg string, fields map[string]any) {
	if s.level <= LogLevelError {
		s.logger.Error(msg, fieldsToAttrs(fields)...)
	}
}

// With returns a new logger with the given fields added to all log entries.
// This is useful for adding context that applies to multiple log statements.
func (s *SlogLogger) With(fields map[string]any) Logger {
	if len(fields) == 0 {
		return s
	}

	return &SlogLogger{
		logger: s.logger.With(fieldsToAttrs(fields)...),
		level:  s.level,
	}
}

// Deprecated: Use NewSlogLogger or NewDefaultLogger instead.
// StandardLogger is kept for backward compatibility.
type StandardLogger struct {
	*SlogLogger
}

// Deprecated: Use NewDefaultLogger instead.
// NewStandardLogger creates a new logger using slog with text output.
// This function is kept for backward compatibility.
func NewStandardLogger(level LogLevel) *StandardLogger {
	return &StandardLogger{
		SlogLogger: NewDefaultLogger(level),
	}
}

// Metrics tracks operational metrics for the FileMaker client.
type Metrics struct {
	RequestsTotal     atomic.Uint64
	RequestsSucceeded atomic.Uint64
	RequestsFailed    atomic.Uint64
	RetriesTotal      atomic.Uint64
	SessionsCreated   atomic.Uint64
	SessionsClosed    atomic.Uint64
}

// NewMetrics creates a new Metrics instance.
func NewMetrics() *Metrics {
	return &Metrics{}
}

// IncrementRequests increments the total request counter.
func (m *Metrics) IncrementRequests() {
	m.RequestsTotal.Add(1)
}

// IncrementRequestsSucceeded increments the successful request counter.
func (m *Metrics) IncrementRequestsSucceeded() {
	m.RequestsSucceeded.Add(1)
}

// IncrementRequestsFailed increments the failed request counter.
func (m *Metrics) IncrementRequestsFailed() {
	m.RequestsFailed.Add(1)
}

// IncrementRetries increments the retry counter.
func (m *Metrics) IncrementRetries() {
	m.RetriesTotal.Add(1)
}

// IncrementSessionsCreated increments the sessions created counter.
func (m *Metrics) IncrementSessionsCreated() {
	m.SessionsCreated.Add(1)
}

// IncrementSessionsClosed increments the sessions closed counter.
func (m *Metrics) IncrementSessionsClosed() {
	m.SessionsClosed.Add(1)
}

// GetRequestsTotal returns the total number of requests.
func (m *Metrics) GetRequestsTotal() uint64 {
	return m.RequestsTotal.Load()
}

// GetRequestsSucceeded returns the number of successful requests.
func (m *Metrics) GetRequestsSucceeded() uint64 {
	return m.RequestsSucceeded.Load()
}

// GetRequestsFailed returns the number of failed requests.
func (m *Metrics) GetRequestsFailed() uint64 {
	return m.RequestsFailed.Load()
}

// GetRetriesTotal returns the total number of retries.
func (m *Metrics) GetRetriesTotal() uint64 {
	return m.RetriesTotal.Load()
}

// GetSessionsCreated returns the number of sessions created.
func (m *Metrics) GetSessionsCreated() uint64 {
	return m.SessionsCreated.Load()
}

// GetSessionsClosed returns the number of sessions closed.
func (m *Metrics) GetSessionsClosed() uint64 {
	return m.SessionsClosed.Load()
}

// Reset resets all metrics to zero.
func (m *Metrics) Reset() {
	m.RequestsTotal.Store(0)
	m.RequestsSucceeded.Store(0)
	m.RequestsFailed.Store(0)
	m.RetriesTotal.Store(0)
	m.SessionsCreated.Store(0)
	m.SessionsClosed.Store(0)
}

// Snapshot returns a snapshot of all metrics.
func (m *Metrics) Snapshot() map[string]uint64 {
	return map[string]uint64{
		"requests_total":     m.GetRequestsTotal(),
		"requests_succeeded": m.GetRequestsSucceeded(),
		"requests_failed":    m.GetRequestsFailed(),
		"retries_total":      m.GetRetriesTotal(),
		"sessions_created":   m.GetSessionsCreated(),
		"sessions_closed":    m.GetSessionsClosed(),
		"active_sessions":    m.GetSessionsCreated() - m.GetSessionsClosed(),
	}
}

// RequestContext adds context-specific information for request tracing.
type RequestContext struct {
	RequestID string
	Operation string
	Database  string
	Layout    string
	StartTime time.Time
}

// NewRequestContext creates a new RequestContext with a generated ID.
func NewRequestContext(operation, database, layout string) *RequestContext {
	return &RequestContext{
		RequestID: generateRequestID(),
		Operation: operation,
		Database:  database,
		Layout:    layout,
		StartTime: time.Now(),
	}
}

// Duration returns the elapsed time since the request started.
func (r *RequestContext) Duration() time.Duration {
	return time.Since(r.StartTime)
}

// Fields returns a map of fields for logging.
func (r *RequestContext) Fields() map[string]any {
	return map[string]any{
		"request_id":  r.RequestID,
		"operation":   r.Operation,
		"database":    r.Database,
		"layout":      r.Layout,
		"duration_ms": r.Duration().Milliseconds(),
	}
}

// SlogAttrs returns slog attributes for structured logging.
func (r *RequestContext) SlogAttrs() []slog.Attr {
	return []slog.Attr{
		slog.String("request_id", r.RequestID),
		slog.String("operation", r.Operation),
		slog.String("database", r.Database),
		slog.String("layout", r.Layout),
		slog.Int64("duration_ms", r.Duration().Milliseconds()),
	}
}

// generateRequestID generates a simple request ID for tracing.
// In production, you might want to use a UUID library or distributed tracing ID.
func generateRequestID() string {
	return fmt.Sprintf("req_%d", time.Now().UnixNano())
}

// contextKey is a private type for context keys to avoid collisions.
type contextKey string

const (
	requestContextKey contextKey = "filemaker_request_context"
	loggerContextKey  contextKey = "filemaker_logger"
)

// WithRequestContext adds RequestContext to the context.
func WithRequestContext(ctx context.Context, reqCtx *RequestContext) context.Context {
	return context.WithValue(ctx, requestContextKey, reqCtx)
}

// GetRequestContext retrieves RequestContext from the context.
func GetRequestContext(ctx context.Context) (*RequestContext, bool) {
	reqCtx, ok := ctx.Value(requestContextKey).(*RequestContext)
	return reqCtx, ok
}

// WithLogger adds a Logger to the context.
// This allows passing logger through context chain.
func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

// GetLogger retrieves Logger from the context.
// If no logger is found, returns a NoOpLogger.
func GetLogger(ctx context.Context) Logger {
	if logger, ok := ctx.Value(loggerContextKey).(Logger); ok {
		return logger
	}
	return &NoOpLogger{}
}

// LogFromContext logs a message using the logger from context.
// If no logger is in context, the message is discarded.
func LogFromContext(ctx context.Context, level LogLevel, msg string, fields map[string]any) {
	logger := GetLogger(ctx)
	switch level {
	case LogLevelDebug:
		logger.Debug(msg, fields)
	case LogLevelInfo:
		logger.Info(msg, fields)
	case LogLevelWarn:
		logger.Warn(msg, fields)
	case LogLevelError:
		logger.Error(msg, fields)
	}
}
