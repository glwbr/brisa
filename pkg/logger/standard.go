package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"maps"
	"os"
)

// StandardLogger is a Logger implementation using Go's standard log package.
// Supports method chaining for contextual logging:
//
//	NewStandardLogger(os.Stdout, InfoLevel).
//	    WithField("service", "payments").
//	    WithFields(map[string]any{"version": "1.2.3"}).
//	    Info("transaction processed", "amount", 19.99)
//	// Output: [INFO] transaction processed service=payments version=1.2.3 amount=19.99
type StandardLogger struct {
	logger *log.Logger
	level  Level
	fields map[string]any
}

// NewStandardLogger creates a new StandardLogger writing to the specified output.
// If out is nil, os.Stderr is used. The level sets the minimum logging threshold.
func NewStandardLogger(out io.Writer, level Level) *StandardLogger {
	if out == nil {
		out = os.Stderr
	}
	return &StandardLogger{
		logger: log.New(out, "", log.LstdFlags),
		level:  level,
		fields: make(map[string]any),
	}
}

// Debug logs a debug message with optional key-value pairs.
// Example: logger.Debug("started", "count", 42)
func (l *StandardLogger) Debug(msg string, args ...any) { l.logAt(DebugLevel, msg, args...) }

// Info logs an informational message with optional key-value pairs.
func (l *StandardLogger) Info(msg string, args ...any) { l.logAt(InfoLevel, msg, args...) }

// Warn logs a warning message with optional key-value pairs.
func (l *StandardLogger) Warn(msg string, args ...any) { l.logAt(WarnLevel, msg, args...) }

// Error logs an error message with optional key-value pairs.
func (l *StandardLogger) Error(msg string, args ...any) { l.logAt(ErrorLevel, msg, args...) }

// WithContext returns a new logger with context support (currently unimplemented).
func (l *StandardLogger) WithContext(ctx context.Context) Logger { return l.clone() }

// WithField returns a new logger with an additional field. The field will be included
// in all subsequent log entries made with the returned logger.
//
// Example:
//
//	logger := NewStandardLogger(os.Stdout, InfoLevel)
//	userLogger := logger.WithField("user_id", 42)
//	userLogger.Info("profile updated")
//	// Output: [INFO] profile updated user_id=42
func (l *StandardLogger) WithField(key string, value any) Logger {
	newLogger := l.clone()
	newLogger.fields[key] = value
	return newLogger
}

// WithFields returns a new logger with additional fields. These fields will be included
// in all subsequent log entries made with the returned logger.
//
// Example:
//
//	logger := NewStandardLogger(os.Stdout, InfoLevel)
//	contextLogger := logger.WithFields(map[string]any{
//	    "service":  "auth",
//	    "requestID": "12345",
//	})
//	contextLogger.Info("user logged in", "user_id", 42)
//	// Output: [INFO] user logged in service=auth requestID=12345 user_id=42
func (l *StandardLogger) WithFields(fields map[string]any) Logger {
	newLogger := l.clone()
	maps.Copy(newLogger.fields, fields)
	return newLogger
}

// clone creates a deep copy of the logger with the same configuration and fields.
func (l *StandardLogger) clone() *StandardLogger {
	newFields := make(map[string]any, len(l.fields))
	maps.Copy(newFields, l.fields)
	return &StandardLogger{
		logger: l.logger,
		level:  l.level,
		fields: newFields,
	}
}

// formatFields combines base fields with additional arguments into a formatted string.
// If args has an odd length, a placeholder value is added for the last key.
func formatFields(base map[string]any, args []any) string {
	if len(args)%2 != 0 {
		args = append(args, "<MISSING_VALUE>")
	}

	fields := make(map[string]any, len(base)+len(args)/2)
	maps.Copy(fields, base)

	for i := 0; i < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			key = fmt.Sprintf("invalid_key_%d", i)
		}
		fields[key] = args[i+1]
	}

	result := ""
	for k, v := range fields {
		result += " " + k + "=" + stringify(v)
	}
	return result
}

// logAt logs a message at the specified level if it meets the logger's level threshold.
func (l *StandardLogger) logAt(target Level, msg string, args ...any) {
	if l.level <= target {
		l.log(target.String(), msg, args...)
	}
}

// log performs the actual logging operation with formatted fields.
func (l *StandardLogger) log(level, msg string, args ...any) {
	fieldString := formatFields(l.fields, args)
	l.logger.Printf("[%s] %s%s", level, msg, fieldString)
}

// stringify converts any value to a loggable string representation.
// Nil values become "<nil>", all others use default formatting.
func stringify(v any) string {
	if v == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", v)
}
