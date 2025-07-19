package logger

import "encoding/json"

// PrettyLogger is a wrapper around an existing Logger implementation
// that adds pretty-printing capabilities for structured log data.
type PrettyLogger struct {
	Logger
}

// DebugPretty formats the log message and the associated data as pretty-printed JSON.
func (p PrettyLogger) DebugPretty(msg string, data any) {
	p.pretty(DebugLevel, msg, data)
}

// InfoPretty formats the log message and the associated data as pretty-printed JSON.
func (p PrettyLogger) InfoPretty(msg string, data any) {
	p.pretty(InfoLevel, msg, data)
}

// Marshals the provided data into a pretty-printed JSON string and logs it at the specified level.
func (p PrettyLogger) pretty(level Level, msg string, data any) {
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		p.Error("failed to marshal data", "error", err)
		return
	}

	switch level {
	case DebugLevel:
		p.Debug(msg + "\n" + string(prettyJSON))
	case InfoLevel:
		p.Info(msg + "\n" + string(prettyJSON))
	}
}
