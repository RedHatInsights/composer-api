package server

import (
	"net/http"

	"github.com/RedHatInsights/composer-api/internal/config"
	"github.com/RedHatInsights/composer-api/internal/handler/probe"
	"github.com/RedHatInsights/composer-api/internal/handler/v1/workspace"
	"github.com/RedHatInsights/composer-api/internal/middleware"
	"github.com/RedHatInsights/composer-api/internal/response"
)

// mountResource registers resourceHandler under prefix on mux using
// http.StripPrefix. It also registers the bare prefix (without trailing
// slash) to prevent Go's ServeMux from issuing a broken redirect after
// StripPrefix.
func mountResource(mux *http.ServeMux, prefix string, resourceHandler http.Handler) {
	mux.Handle(prefix+"/", http.StripPrefix(prefix, resourceHandler))
	mux.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
		r2 := r.Clone(r.Context())
		r2.URL.Path = "/"
		r2.URL.RawPath = ""
		resourceHandler.ServeHTTP(w, r2)
	})
}

// registerRoutes registers probe routes for liveness/readiness checks
// and delegates versioned route groups to their own registration functions.
// It returns the top-level handler wrapped with TrailingSlash so that
// paths with and without a trailing slash are handled consistently.
func registerRoutes(mux *http.ServeMux, cfg config.Config) http.Handler {
	commonMiddleware := middleware.Chain(
		middleware.Recover,
		middleware.CORS(cfg.Server.AllowedOrigins, cfg.Server.CORSMaxAge),
		middleware.RequestID,
		middleware.Logging,
		middleware.BodySizeLimit(cfg.Server.MaxBodyBytes),
	)

	probeHandler := probe.New()
	probeMux := http.NewServeMux()
	probeMux.HandleFunc("GET /ping", probeHandler.Ping)
	probeMux.HandleFunc("GET /health", probeHandler.Health)
	probeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response.WriteError(r.Context(), w, response.NotFound().WithReasonStr("route not found"))
	})
	mux.Handle("/", commonMiddleware(probeMux))

	// v1 Routes
	registerV1Routes(mux, commonMiddleware)

	return middleware.TrailingSlash(mux)
}

// registerV1Routes creates a sub-mux per resource, wraps each with its
// middleware stack, and mounts them under /v1/ on the top-level mux.
// Each resource gets its own sub-mux and middleware stack.
//
// To add a new resource:
//
//	f := v1feature.New()
//	fMux := http.NewServeMux()
//	fMux.HandleFunc("GET /{$}", f.List)
//	fHandler := commonMiddleware(fMux)
//	mountResource(v1Mux, "/features", fHandler)
//
// To add resource-specific middleware, build a custom chain:
//
//	authMiddleware := middleware.Chain(commonMiddleware, middleware.Auth, middleware.RBAC)
//	fHandler := authMiddleware(fMux)
//	mountResource(v1Mux, "/features", fHandler)
func registerV1Routes(mux *http.ServeMux, commonMiddleware func(http.Handler) http.Handler) {
	v1Mux := http.NewServeMux()
	v1Mux.Handle("/", commonMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.WriteError(r.Context(), w, response.NotFound().WithReasonStr("route not found"))
	})))

	ws := workspace.New()
	wsMux := http.NewServeMux()
	wsMux.HandleFunc("GET /{$}", ws.List)
	wsMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response.WriteError(r.Context(), w, response.NotFound().WithReasonStr("route not found"))
	})

	wsHandler := commonMiddleware(wsMux)
	mountResource(v1Mux, "/workspaces", wsHandler)
	mux.Handle("/v1/", http.StripPrefix("/v1", v1Mux))

}
