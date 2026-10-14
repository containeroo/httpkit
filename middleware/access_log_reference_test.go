package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/containeroo/httpkit/response"
)

func TestAccessLogWithErrorReference(t *testing.T) {
	var output strings.Builder
	logger := slog.New(slog.NewTextHandler(&output, nil))
	handler := AccessLogWithErrorReference(logger)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		response.LogFailure(logger, w, errors.New("internal failure"))
		w.WriteHeader(http.StatusInternalServerError)
	}))
	request := httptest.NewRequest(http.MethodGet, "/private?token=hidden", nil)
	handler.ServeHTTP(response.NewRequestWriter(httptest.NewRecorder(), request), request)
	if !strings.Contains(output.String(), "error_reference=") || strings.Contains(output.String(), "token=hidden") {
		t.Fatal(output.String())
	}
}
