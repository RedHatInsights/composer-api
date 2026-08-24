package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/RedHatInsights/composer-api/internal/response"
)

// Recover catches panics during request handling, logs the error with
// a stack trace, and returns a 500 response without exposing internals.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err == nil {
				return
			}

			stack := string(debug.Stack())

			slog.ErrorContext(r.Context(), "panic during request",
				"error", err,
				"stack", stack,
			)

			response.WriteError(r.Context(), w, response.InternalServerError())
		}()

		next.ServeHTTP(w, r)
	})
}
