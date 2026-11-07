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

type accessLogExtractor func([]any, *response.RequestWriter, *http.Request, time.Duration) []any

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
	extractFields, trackResponse := buildAccessLogExtractor(config.fields)
	attributes := config.attributes

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			started := time.Now()
			writer := w
			var tracked *response.RequestWriter
			if trackResponse {
				if existing, ok := w.(*response.RequestWriter); ok {
					tracked = existing
				} else {
					tracked = response.NewRequestWriter(w, r)
					writer = tracked
				}
			}
			next.ServeHTTP(writer, r)
			args := []any{"event", "request_complete"}
			args = extractFields(args, tracked, r, time.Since(started))
			if attributes != nil {
				args = append(args, attributes(writer, r)...)
			}
			logger.Info("request", args...)
		})
	}
}

// buildAccessLogExtractor compiles the configured field selection once when the
// middleware is constructed. Its returned flag indicates whether a response
// writer must be tracked to populate the selected fields.
func buildAccessLogExtractor(fields []AccessLogField) (accessLogExtractor, bool) {
	extractors := make([]accessLogExtractor, 0, len(fields))
	trackResponse := false

	for _, field := range fields {
		switch field {
		case LogMethod:
			extractors = append(extractors, func(args []any, _ *response.RequestWriter, r *http.Request, _ time.Duration) []any {
				return append(args, "method", r.Method)
			})
		case LogPath:
			extractors = append(extractors, func(args []any, _ *response.RequestWriter, r *http.Request, _ time.Duration) []any {
				return append(args, "path", r.URL.Path)
			})
		case LogStatus:
			trackResponse = true
			extractors = append(extractors, func(args []any, tracked *response.RequestWriter, _ *http.Request, _ time.Duration) []any {
				status := http.StatusOK
				if tracked != nil && tracked.Status != 0 {
					status = tracked.Status
				}
				return append(args, "status", status)
			})
		case LogDuration:
			extractors = append(extractors, func(args []any, _ *response.RequestWriter, _ *http.Request, duration time.Duration) []any {
				return append(args, "duration", duration)
			})
		case LogErrorReference:
			trackResponse = true
			extractors = append(extractors, func(args []any, tracked *response.RequestWriter, _ *http.Request, _ time.Duration) []any {
				if tracked != nil && tracked.ErrorReference != "" {
					return append(args, "error_reference", tracked.ErrorReference)
				}
				return args
			})
		}
	}

	return func(args []any, tracked *response.RequestWriter, r *http.Request, duration time.Duration) []any {
		for _, extractor := range extractors {
			args = extractor(args, tracked, r, duration)
		}
		return args
	}, trackResponse
}

// AccessLogWithAttributes records request completion with application-specific attributes.
func AccessLogWithAttributes(logger *slog.Logger, attributes func(http.ResponseWriter, *http.Request) []any) Middleware {
	return AccessLog(logger, func(config *accessLogConfig) { config.attributes = attributes })
}

// AccessLogWithErrorReference records an error correlation reference when present.
func AccessLogWithErrorReference(logger *slog.Logger) Middleware {
	return AccessLog(logger, WithAccessLogFields(LogMethod, LogPath, LogDuration, LogErrorReference))
}
