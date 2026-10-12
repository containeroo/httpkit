package response

import "net/http"

// StatusWriter records the final HTTP status while preserving response-controller support.
type StatusWriter struct {
	http.ResponseWriter
	Status int
}

// NewStatusWriter wraps w with final-status diagnostics.
func NewStatusWriter(w http.ResponseWriter) *StatusWriter {
	return &StatusWriter{ResponseWriter: w}
}

// Unwrap returns the underlying response writer.
func (w *StatusWriter) Unwrap() http.ResponseWriter { return w.ResponseWriter }

// WriteHeader records the first final status before passing it through.
func (w *StatusWriter) WriteHeader(status int) {
	if w.Status != 0 {
		return
	}
	if status >= http.StatusOK {
		w.Status = status
	}
	w.ResponseWriter.WriteHeader(status)
}

// Write records an implicit successful status before writing data.
func (w *StatusWriter) Write(data []byte) (int, error) {
	if w.Status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(data)
}

// FlushError flushes the underlying response after recording an implicit status.
func (w *StatusWriter) FlushError() error {
	if w.Status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return http.NewResponseController(w.ResponseWriter).Flush()
}
