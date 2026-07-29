package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

// New creates an HTTP server. The provided context is used as the base
// context for all incoming requests, so canceling it signals in-flight
// handlers to stop work during shutdown.
func New(ctx context.Context, port string) *Server {
	mux := http.NewServeMux()
	registerRoutes(mux)

	return &Server{
		httpServer: &http.Server{
			BaseContext:       func(_ net.Listener) context.Context { return ctx },
			Addr:              net.JoinHostPort("", port),
			Handler:           mux,
			ReadHeaderTimeout: 10 * time.Second,
			ReadTimeout:       30 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       120 * time.Second,
			MaxHeaderBytes:    64 * 1024,
		},
	}
}

// Start begins listening for HTTP requests. It blocks until the server
// is shut down and returns nil on graceful close.
func (s *Server) Start() error {
	err := s.httpServer.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

// Shutdown stops the server from accepting new connections and waits
// for in-flight requests to complete, subject to the context deadline.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// Close releases resources held by stateful dependencies (database
// pools, cache clients, message queues, etc.).
// TODO: close dependencies here as they are added.
func (s *Server) Close(ctx context.Context) error {
	return nil
}
