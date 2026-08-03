package middleware

import "net/http"

// BodySizeLimit returns a middleware that restricts the maximum number
// of bytes that can be read from a request body. Requests exceeding the
// limit receive an error when the handler reads past maxBytes.
func BodySizeLimit(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next.ServeHTTP(w, r)
		})
	}
}
