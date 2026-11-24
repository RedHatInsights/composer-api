package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// wrappedResponseWriter intercepts WriteHeader and Write so we can
// capture the status code and, for error responses, the response body.
type wrappedResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (w *wrappedResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// Write captures the response body only for error status codes (>= 400)
// to avoid buffering successful response payloads in memory.
func (w *wrappedResponseWriter) Write(b []byte) (int, error) {
	if w.statusCode >= 400 {
		w.body.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

func (w *wrappedResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

// Logging is an HTTP middleware that logs every incoming request and its outcome.
// Successful responses (< 400) are logged at INFO level.
// Error responses (>= 400) are logged at ERROR level and include the response body.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := r.Context()

		// Use RequestURI for logging because StripPrefix modifies r.URL.Path
		// before this middleware runs, losing the original client path.
		originalPath := r.RequestURI
		if i := strings.IndexByte(originalPath, '?'); i >= 0 {
			originalPath = originalPath[:i]
		}

		slog.InfoContext(ctx, "request received",
			"method", r.Method,
			"path", originalPath,
			"query", r.URL.RawQuery,
		)

		// Default to 200; overwritten if the handler calls WriteHeader explicitly.
		wrapped := &wrappedResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)

		attrs := []any{
			"method", r.Method,
			"path", originalPath,
			"status", wrapped.statusCode,
			"duration_ms", time.Since(start).Milliseconds(),
		}

		if wrapped.statusCode >= 400 {
			attrs = append(attrs, "response_body", wrapped.body.String())
			slog.ErrorContext(ctx, "request failed", attrs...)
		} else {
			slog.InfoContext(ctx, "request completed", attrs...)
		}
	})
}
