package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RedHatInsights/composer-api/internal/config"
	"github.com/RedHatInsights/composer-api/internal/logger"
	"github.com/RedHatInsights/composer-api/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	logger.Init(os.Stdout, cfg.Log.Level, cfg.Log.Pretty)

	// Root context is canceled on SIGINT/SIGTERM or if the server
	// goroutine exits, triggering the graceful shutdown sequence.
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	srv := server.New(ctx, cfg)

	var startErr error
	go func() {
		defer cancel()
		slog.Info("starting server", "port", cfg.Server.Port, "log_level", cfg.Log.Level)
		if err := srv.Start(); err != nil {
			startErr = err
			slog.Error("server error", "error", err)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down server")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown error", "error", err)
		os.Exit(1)
	}

	if err := srv.Close(shutdownCtx); err != nil {
		slog.Error("failed to close dependencies", "error", err)
		os.Exit(1)
	}

	if startErr != nil {
		os.Exit(1)
	}

	slog.Info("server stopped")
}
