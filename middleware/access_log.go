package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/containeroo/httpkit/response"
)

// AccessLog records standard request completion attributes.
func AccessLog(logger *slog.Logger) Middleware {
	return AccessLogWithAttributes(logger, nil)
}

// AccessLogWithAttributes records request completion with application-specific attributes.
func AccessLogWithAttributes(logger *slog.Logger, attributes func(http.ResponseWriter, *http.Request) []any) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			started := time.Now()
			next.ServeHTTP(w, r)
			args := []any{"event", "request_complete", "method", r.Method, "path", r.URL.Path, "duration", time.Since(started)}
			if attributes != nil {
				args = append(args, attributes(w, r)...)
			}
			logger.Info("request", args...)
		})
	}
}

// AccessLogWithErrorReference logs the reference associated with an HTTP 500 response.
func AccessLogWithErrorReference(logger *slog.Logger) Middleware {
	return AccessLogWithAttributes(logger, func(w http.ResponseWriter, _ *http.Request) []any {
		request, ok := w.(*response.RequestWriter)
		if !ok || request.ErrorReference == "" {
			return nil
		}
		return []any{"error_reference", request.ErrorReference}
	})
}
