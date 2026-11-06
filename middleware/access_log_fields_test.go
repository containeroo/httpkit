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

func TestAccessLogSelectedFields(t *testing.T) {
	var output strings.Builder
	logger := slog.New(slog.NewTextHandler(&output, nil))
	handler := AccessLog(logger, WithAccessLogFields(LogMethod, LogStatus, LogDuration, LogErrorReference))(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		response.LogFailure(logger, w, errors.New("failed"))
		w.WriteHeader(http.StatusInternalServerError)
	}))
	request := httptest.NewRequest(http.MethodPost, "/private?secret=hidden", nil)
	handler.ServeHTTP(httptest.NewRecorder(), request)
	logged := output.String()
	for _, expected := range []string{"method=POST", "status=500", "duration=", "error_reference="} {
		if !strings.Contains(logged, expected) {
			t.Errorf("missing %s in %s", expected, logged)
		}
	}
	if strings.Contains(logged, "path=") || strings.Contains(logged, "secret=hidden") {
		t.Errorf("unexpected path in %s", logged)
	}
}

func TestAccessLogDefaultFields(t *testing.T) {
	var output strings.Builder
	logger := slog.New(slog.NewTextHandler(&output, nil))
	handler := AccessLog(logger)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) }))
	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/health?secret=hidden", nil))
	logged := output.String()
	for _, expected := range []string{"method=GET", "path=/health", "duration="} {
		if !strings.Contains(logged, expected) {
			t.Errorf("missing %s in %s", expected, logged)
		}
	}
	if strings.Contains(logged, "status=") || strings.Contains(logged, "secret=hidden") {
		t.Errorf("unexpected field in %s", logged)
	}
}

func TestAccessLogImplicitStatus(t *testing.T) {
	var output strings.Builder
	logger := slog.New(slog.NewTextHandler(&output, nil))
	handler := AccessLog(logger, WithAccessLogFields(LogStatus))(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok")) }))
	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	if !strings.Contains(output.String(), "status=200") {
		t.Fatal(output.String())
	}
}
