package middleware

import (
	"net/http"
	"slices"
	"strings"
)

// CORSOptions selects explicitly trusted origins for cross-origin browser requests.
type CORSOptions struct {
	// AllowedOrigins contains complete origins, including scheme and port.
	AllowedOrigins []string
	// AllowedMethods contains permitted preflight HTTP methods.
	AllowedMethods []string
	// AllowedHeaders contains permitted preflight request headers.
	AllowedHeaders []string
	// AllowCredentials permits cookies or HTTP credentials from allowed origins.
	AllowCredentials bool
}

// CORS validates an origin against an explicit allowlist before adding CORS headers.
// An empty allowlist disables CORS, without blocking ordinary same-origin requests.
func CORS(options CORSOptions) Middleware {
	origins := slices.Clone(options.AllowedOrigins)
	methods := strings.Join(options.AllowedMethods, ", ")
	headers := strings.Join(options.AllowedHeaders, ", ")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "" || !slices.Contains(origins, origin) {
				// Denied preflight requests must not reach the application handler.
				if origin != "" && r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
					w.WriteHeader(http.StatusForbidden)
					return
				}
				next.ServeHTTP(w, r)
				return
			}
			w.Header().Add("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if options.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
				method := r.Header.Get("Access-Control-Request-Method")
				if !slices.Contains(options.AllowedMethods, method) {
					w.WriteHeader(http.StatusForbidden)
					return
				}
				if methods != "" {
					w.Header().Set("Access-Control-Allow-Methods", methods)
				}
				if headers != "" {
					w.Header().Set("Access-Control-Allow-Headers", headers)
				}
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
