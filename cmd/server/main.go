// Small HTTP microservice: /health for liveness, /process to transform a payload
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"thesisPrototype/internal/config"
	"thesisPrototype/internal/handler"
	"thesisPrototype/internal/logging"
	"thesisPrototype/internal/service"
)

func main() {
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		slog.Default().Error("invalid config", "error", err)
		os.Exit(1)
	}

	logger := logging.NewLogger(logging.LevelFromString(cfg.LogLevel))
	slog.SetDefault(logger)

	proc := service.NewProcessService()
	h := &handler.Handler{
		Processor: proc,
		Logger:    logger,
		Timeout:   cfg.RequestTimeout,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("POST /process", h.Process)
	mux.HandleFunc("GET /process", h.Process)

	srv := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	// wait for SIGTERM/SIGINT then shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
		sig := <-sigCh
		logger.Info("shutdown signal received", "signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("server shutdown error", "error", err)
			os.Exit(1)
		}
		logger.Info("server stopped gracefully")
	}()

	logger.Info("server listening", "addr", cfg.HTTPAddr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}
