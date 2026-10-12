package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChain(t *testing.T) {
	t.Parallel()
	var calls []string
	record := func(name string) Middleware {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls = append(calls, name+":before")
				next.ServeHTTP(w, r)
				calls = append(calls, name+":after")
			})
		}
	}
	handler := Chain(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { calls = append(calls, "handler") }), record("first"), record("second"))
	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	assert.Equal(t, []string{"first:before", "second:before", "handler", "second:after", "first:after"}, calls)
}

func TestRecoverPanics(t *testing.T) {
	t.Parallel()
	response := httptest.NewRecorder()
	handler := RecoverPanics(func(w http.ResponseWriter, _ *http.Request, recovered any) {
		http.Error(w, "recovered", http.StatusInternalServerError)
	}, nil)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic("example") }))
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))
	assert.Equal(t, http.StatusInternalServerError, response.Code)
}

func TestSecurityHeaders(t *testing.T) {
	t.Parallel()
	response := httptest.NewRecorder()
	handler := SecurityHeaders(http.Header{"X-Content-Type-Options": {"nosniff"}})(http.NotFoundHandler())
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))
	assert.Equal(t, "nosniff", response.Header().Get("X-Content-Type-Options"))
}

func TestRejectCrossOriginWrites(t *testing.T) {
	t.Parallel()
	response := httptest.NewRecorder()
	handler := RejectCrossOriginWrites(slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil)), func(w http.ResponseWriter, _ *http.Request) { http.Error(w, "forbidden", http.StatusForbidden) })(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) }))
	request := httptest.NewRequest(http.MethodPost, "https://wiki.example.test/pages/example", nil)
	request.Header.Set("Sec-Fetch-Site", "cross-site")
	handler.ServeHTTP(response, request)
	assert.Equal(t, http.StatusForbidden, response.Code)
}

func TestIsCrossOriginWrite(t *testing.T) {
	t.Parallel()
	request := httptest.NewRequest(http.MethodPost, "https://wiki.example.test/pages/example", nil)
	request.Header.Set("Origin", "https://attacker.example")
	assert.True(t, IsCrossOriginWrite(request))
}
