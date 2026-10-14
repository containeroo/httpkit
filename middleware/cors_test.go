package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSPreflight(t *testing.T) {
	handler := CORS(CORSOptions{AllowedOrigins: []string{"https://app.example"}, AllowedMethods: []string{"GET", "POST"}, AllowedHeaders: []string{"Content-Type"}, AllowCredentials: true})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(200) }))
	for _, tc := range []struct {
		origin, method string
		want           int
	}{{"https://app.example", "POST", 204}, {"https://other.example", "POST", 403}, {"https://app.example", "DELETE", 403}} {
		req := httptest.NewRequest(http.MethodOptions, "/api", nil)
		req.Header.Set("Origin", tc.origin)
		req.Header.Set("Access-Control-Request-Method", tc.method)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, req)
		if response.Code != tc.want {
			t.Errorf("origin %q method %q got %d", tc.origin, tc.method, response.Code)
		}
		if tc.want == 204 && response.Header().Get("Access-Control-Allow-Origin") != tc.origin {
			t.Fatal("missing CORS header")
		}
	}
}
