package middleware

import (
	"net/http"
	"strings"
)

// TrailingSlash redirects requests with a trailing slash to the same
// path without it (e.g. /ping/ → /ping). The root path "/" is left as-is.
func TrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		escaped := r.URL.EscapedPath()
		if len(escaped) > 1 && strings.HasSuffix(escaped, "/") {
			target := "/" + strings.Trim(escaped, "/")
			if r.URL.RawQuery != "" {
				target += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, target, http.StatusMovedPermanently)
			return
		}
		next.ServeHTTP(w, r)
	})
}
