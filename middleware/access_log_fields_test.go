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
	if index := strings.LastIndex(logged, "msg=request "); index >= 0 {
		logged = logged[index:]
	}
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

func TestAccessLogResponseTrackingIsSelectedAtConstruction(t *testing.T) {
	for _, test := range []struct {
		name          string
		fields        []AccessLogField
		trackResponse bool
	}{
		{name: "default fields", fields: []AccessLogField{LogMethod, LogPath, LogDuration}},
		{name: "status", fields: []AccessLogField{LogStatus}, trackResponse: true},
		{name: "error reference", fields: []AccessLogField{LogErrorReference}, trackResponse: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, trackResponse := buildAccessLogExtractor(test.fields)
			if trackResponse != test.trackResponse {
				t.Errorf("trackResponse = %t, want %t", trackResponse, test.trackResponse)
			}
		})
	}
}

func TestAccessLogAttributesReceiveOriginalWriterWithoutResponseFields(t *testing.T) {
	var received http.ResponseWriter
	handler := AccessLog(slog.New(slog.NewTextHandler(&strings.Builder{}, nil)), func(config *accessLogConfig) {
		config.attributes = func(w http.ResponseWriter, _ *http.Request) []any {
			received = w
			return nil
		}
	})(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))

	writer := httptest.NewRecorder()
	handler.ServeHTTP(writer, httptest.NewRequest(http.MethodGet, "/", nil))
	if received != writer {
		t.Fatalf("attributes received %T, want original writer", received)
	}
}
