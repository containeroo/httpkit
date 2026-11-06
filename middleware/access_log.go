package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/containeroo/httpkit/response"
)

// AccessLogField identifies an HTTP request completion attribute.
type AccessLogField string

const (
	// LogMethod records the HTTP method.
	LogMethod AccessLogField = "method"
	// LogPath records the URL path without its query string.
	LogPath AccessLogField = "path"
	// LogStatus records the final HTTP response status.
	LogStatus AccessLogField = "status"
	// LogDuration records elapsed request processing time.
	LogDuration AccessLogField = "duration"
	// LogErrorReference records a correlation reference when present.
	LogErrorReference AccessLogField = "error_reference"
)

// AccessLogOption configures request completion logging.
type AccessLogOption func(*accessLogConfig)

type accessLogConfig struct {
	fields     []AccessLogField
	attributes func(http.ResponseWriter, *http.Request) []any
}

// WithAccessLogFields selects the fields emitted by AccessLog, replacing the defaults.
func WithAccessLogFields(fields ...AccessLogField) AccessLogOption {
	selected := append([]AccessLogField(nil), fields...)
	return func(config *accessLogConfig) { config.fields = selected }
}

// AccessLog records request completion with optional field selection.
func AccessLog(logger *slog.Logger, options ...AccessLogOption) Middleware {
	config := accessLogConfig{fields: []AccessLogField{LogMethod, LogPath, LogDuration}}
	for _, option := range options {
		if option != nil {
			option(&config)
		}
	}
	if logger == nil {
		logger = slog.Default()
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			started := time.Now()
			writer := w
			var tracked *response.RequestWriter
			if containsAccessLogField(config.fields, LogStatus) || containsAccessLogField(config.fields, LogErrorReference) {
				if existing, ok := w.(*response.RequestWriter); ok {
					tracked = existing
				} else {
					tracked = response.NewRequestWriter(w, r)
					writer = tracked
				}
			}
			next.ServeHTTP(writer, r)
			args := []any{"event", "request_complete"}
			for _, field := range config.fields {
				switch field {
				case LogMethod:
					args = append(args, "method", r.Method)
				case LogPath:
					args = append(args, "path", r.URL.Path)
				case LogDuration:
					args = append(args, "duration", time.Since(started))
				case LogStatus:
					status := http.StatusOK
					if tracked != nil && tracked.Status != 0 {
						status = tracked.Status
					}
					args = append(args, "status", status)
				case LogErrorReference:
					if tracked != nil && tracked.ErrorReference != "" {
						args = append(args, "error_reference", tracked.ErrorReference)
					}
				}
			}
			if config.attributes != nil {
				args = append(args, config.attributes(writer, r)...)
			}
			logger.Info("request", args...)
		})
	}
}

func containsAccessLogField(fields []AccessLogField, field AccessLogField) bool {
	for _, candidate := range fields {
		if candidate == field {
			return true
		}
	}
	return false
}

// AccessLogWithAttributes records request completion with application-specific attributes.
func AccessLogWithAttributes(logger *slog.Logger, attributes func(http.ResponseWriter, *http.Request) []any) Middleware {
	return AccessLog(logger, func(config *accessLogConfig) { config.attributes = attributes })
}

// AccessLogWithErrorReference records an error correlation reference when present.
func AccessLogWithErrorReference(logger *slog.Logger) Middleware {
	return AccessLog(logger, WithAccessLogFields(LogMethod, LogPath, LogDuration, LogErrorReference))
}
