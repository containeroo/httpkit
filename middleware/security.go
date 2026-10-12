package middleware

import (
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

// SecurityHeaders adds the supplied response headers to every response.
func SecurityHeaders(headers http.Header) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for name, values := range headers {
				w.Header()[name] = append([]string(nil), values...)
			}
			next.ServeHTTP(w, r)
		})
	}
}

// IsSafeMethod reports whether method is safe for cross-site requests.
func IsSafeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	default:
		return false
	}
}

// IsCrossOriginWrite reports whether browser metadata identifies an unsafe request as external.
func IsCrossOriginWrite(r *http.Request) bool {
	switch strings.ToLower(strings.TrimSpace(r.Header.Get("Sec-Fetch-Site"))) {
	case "cross-site", "same-site":
		return true
	case "same-origin":
		return false
	}
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if origin == "" {
		return false
	}
	parsed, err := url.Parse(origin)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return true
	}
	return !strings.EqualFold(parsed.Host, r.Host)
}

// RejectCrossOriginWrites blocks browser-originated unsafe requests.
func RejectCrossOriginWrites(logger *slog.Logger, reject func(http.ResponseWriter, *http.Request)) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if IsSafeMethod(r.Method) || !IsCrossOriginWrite(r) {
				next.ServeHTTP(w, r)
				return
			}
			logger.Warn("cross-site write rejected", "event", "cross_site_write_rejected", "method", r.Method, "path", r.URL.Path, "origin", r.Header.Get("Origin"), "sec_fetch_site", r.Header.Get("Sec-Fetch-Site"), "remote_addr", r.RemoteAddr)
			reject(w, r)
		})
	}
}
