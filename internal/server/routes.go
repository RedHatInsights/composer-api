package server

import (
	"net/http"

	"github.com/RedHatInsights/composer-api/internal/config"
	"github.com/RedHatInsights/composer-api/internal/handler/probe"
	"github.com/RedHatInsights/composer-api/internal/handler/v1/workspace"
	"github.com/RedHatInsights/composer-api/internal/middleware"
	"github.com/RedHatInsights/composer-api/internal/response"
)

// registerRoutes registers probe routes for liveness/readiness checks
// and delegates versioned route groups to their own registration functions.
func registerRoutes(mux *http.ServeMux, cfg config.Config) {
	commonMiddleware := middleware.Chain(
		middleware.Recover,
		middleware.CORS(cfg.Server.AllowedOrigins, cfg.Server.CORSMaxAge),
		middleware.RequestID,
		middleware.Logging,
		middleware.BodySizeLimit(cfg.Server.MaxBodyBytes),
	)

	// Probe routes — wrapped with Recover so panics in dependency
	// checks (future health expansions) don't crash the process.
	probeHandler := probe.New()
	probeMux := http.NewServeMux()
	probeMux.HandleFunc("GET /ping", probeHandler.Ping)
	probeMux.HandleFunc("GET /health", probeHandler.Health)
	probeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response.WriteError(w, response.NotFound().WithReasonStr("route not found"))
	})
	mux.Handle("/", middleware.Recover(probeMux))

	// v1 Routes
	registerV1Routes(mux, commonMiddleware)
}

// registerV1Routes creates a sub-mux per resource, wraps each with its
// middleware stack, and mounts them under /v1/ on the top-level mux.
// Each resource gets its own sub-mux and middleware stack.
//
// To add a new resource:
//
//	fHandler := v1feature.New()
//	fMux := http.NewServeMux()
//	fMux.HandleFunc("GET /", fHandler.List)
//	fMiddleware := commonMiddleware(fMux)
//	v1Mux.Handle("/features/", http.StripPrefix("/features", fMiddleware))
//
// To add resource-specific middleware, replace commonMiddleware with
// a custom chain, e.g.:
//
//	wsMiddleware := middleware.Chain(commonMiddleware, middleware.Auth, middleware.RBAC)(wsMux)
func registerV1Routes(mux *http.ServeMux, commonMiddleware func(http.Handler) http.Handler) {

	wsHandler := workspace.New()

	wsMux := http.NewServeMux()
	wsMux.HandleFunc("GET /", wsHandler.List)
	wsMiddleware := commonMiddleware(wsMux)

	v1Mux := http.NewServeMux()
	v1Mux.Handle("/workspaces/", http.StripPrefix("/workspaces", wsMiddleware))

	mux.Handle("/v1/", http.StripPrefix("/v1", v1Mux))
}
