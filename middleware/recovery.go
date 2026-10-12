package middleware

import (
	"net/http"
)

// PanicResponder handles a recovered panic before the response has committed.
type PanicResponder func(http.ResponseWriter, *http.Request, any)

// RecoverPanics calls responder for handler panics and preserves ErrAbortHandler.
func RecoverPanics(responder PanicResponder, committed func(http.ResponseWriter) bool) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					if recovered == http.ErrAbortHandler {
						panic(recovered)
					}
					wasCommitted := committed != nil && committed(w)
					if !wasCommitted && responder != nil {
						responder(w, r, recovered)
					}
					if wasCommitted {
						panic(http.ErrAbortHandler)
					}
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
