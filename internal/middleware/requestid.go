package middleware

import (
	"context"
	"net/http"

	"github.com/RedHatInsights/composer-api/internal/logger"
	"github.com/google/uuid"
)

type contextKey string

const requestIDKey contextKey = "request_id"

// RequestID generates a unique ID for each request and adds it to the
// request context for use in logging and tracing. The ID is also added
// to the logger context so it appears automatically in all downstream
// slog calls that use slog.InfoContext/ErrorContext.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		ctx = logger.AddContextValue(ctx, "request_id", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestID retrieves the request ID from the context.
func GetRequestID(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey).(string)
	return id
}
