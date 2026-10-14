package response

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestWriterDiagnostics(t *testing.T) {
	underlying := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/test?token=private", nil)
	writer := NewRequestWriter(underlying, request)
	writer.WriteHeader(http.StatusEarlyHints)
	writer.WriteHeader(http.StatusCreated)
	writer.WriteHeader(http.StatusInternalServerError)
	if writer.Status != http.StatusCreated || writer.Path != "/test" {
		t.Fatalf("unexpected request diagnostics: %#v", writer)
	}
	reference := LogFailure(nil, writer, errors.New("private diagnostic"))
	if reference == "" || writer.ErrorReference != reference {
		t.Fatal("missing error reference")
	}
}

func TestRequestWriterImplicitStatus(t *testing.T) {
	writer := NewRequestWriter(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	if _, err := writer.Write([]byte("ok")); err != nil || writer.Status != http.StatusOK {
		t.Fatalf("implicit write: %v", err)
	}
}
