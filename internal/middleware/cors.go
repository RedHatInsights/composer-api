package middleware

import (
	"net/http"
	"slices"
	"strconv"
)

const (
	corsAllowedMethods = "GET, POST, PUT, PATCH, DELETE, OPTIONS"
	corsAllowedHeaders = "Accept, Authorization, Content-Type, X-Correlation-ID"
	corsExposedHeaders = "X-Correlation-ID"
)

// CORS returns a middleware that applies a strict, browser-correct CORS
// policy. It adds Access-Control headers for allowed origins only,
// short-circuits preflight requests, and leaves non-browser clients
// unaffected.
func CORS(origins []string, maxAgeSec int) func(http.Handler) http.Handler {
	allowAll := slices.Contains(origins, "*")
	allowed := make(map[string]struct{}, len(origins))
	for _, o := range origins {
		allowed[o] = struct{}{}
	}
	maxAge := strconv.Itoa(maxAgeSec)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if origin == "" {
				next.ServeHTTP(w, r)
				return
			}

			_, ok := allowed[origin]
			if !allowAll && !ok {
				next.ServeHTTP(w, r)
				return
			}

			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Vary", "Origin")
			w.Header().Set("Access-Control-Expose-Headers", corsExposedHeaders)

			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", corsAllowedMethods)
				w.Header().Set("Access-Control-Allow-Headers", corsAllowedHeaders)
				w.Header().Set("Access-Control-Max-Age", maxAge)
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
