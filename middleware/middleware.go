// Package middleware provides composable HTTP handler wrappers.
package middleware

import (
	"net/http"
	"slices"
)

// Middleware wraps an HTTP handler with cross-cutting behavior.
type Middleware func(http.Handler) http.Handler

// Chain wraps handler with middleware in declaration order.
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range slices.Backward(middlewares) {
		handler = middleware(handler)
	}
	return handler
}
