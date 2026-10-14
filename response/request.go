package response

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
)

// RequestWriter records response status and failure correlation for request diagnostics.
type RequestWriter struct {
	http.ResponseWriter
	// Status is the first final response status, or zero before commit.
	Status int
	// ErrorReference identifies an internal error logged while handling the request.
	ErrorReference string
	// Method is the incoming request method.
	Method string
	// Path is the URL path without query parameters.
	Path string
}

// NewRequestWriter wraps the HTTP response writer with request diagnostics.
func NewRequestWriter(w http.ResponseWriter, r *http.Request) *RequestWriter {
	return &RequestWriter{ResponseWriter: w, Method: r.Method, Path: r.URL.Path}
}

// Unwrap provides access to the underlying writer for response controllers.
func (w *RequestWriter) Unwrap() http.ResponseWriter { return w.ResponseWriter }

// WriteHeader records the first final HTTP status.
func (w *RequestWriter) WriteHeader(status int) {
	if w.Status != 0 {
		return
	}
	if status >= http.StatusOK {
		w.Status = status
	}
	w.ResponseWriter.WriteHeader(status)
}

// Write records an implicit success before writing data.
func (w *RequestWriter) Write(data []byte) (int, error) {
	if w.Status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(data)
}

// FlushError flushes the underlying writer with correct status tracking.
func (w *RequestWriter) FlushError() error {
	if w.Status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return http.NewResponseController(w.ResponseWriter).Flush()
}

// LogFailure records a private error with a public correlation reference.
// The caller owns the response format and must not expose err to the client.
func LogFailure(logger *slog.Logger, w http.ResponseWriter, err error) string {
	if logger == nil {
		logger = slog.Default()
	}
	if request, ok := w.(*RequestWriter); ok {
		logger = logger.With("method", request.Method, "path", request.Path)
	}
	var random [16]byte
	reference := "unavailable"
	if _, e := rand.Read(random[:]); e == nil {
		reference = hex.EncodeToString(random[:])
	}
	if request, ok := w.(*RequestWriter); ok {
		request.ErrorReference = reference
	}
	logger.Error("request failed", "event", "request_failed", "error_reference", reference, "error", err)
	return reference
}
