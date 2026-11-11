package middleware

import (
	"net/http"
	"strings"
)

// TrailingSlash redirects requests with a trailing slash to the same
// path without it (e.g. /ping/ → /ping). The root path "/" is left as-is.
func TrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) > 1 && strings.HasSuffix(r.URL.Path, "/") {
			target := strings.TrimRight(r.URL.Path, "/")
			if r.URL.RawQuery != "" {
				target += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, target, http.StatusMovedPermanently)
			return
		}
		next.ServeHTTP(w, r)
	})
}
