package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aman-void/go-http-server/internal/config"
)

func main() {

	// Load application configuration
	cfg := config.MustLoad()

	// setup database connections

	// Register HTTP routes.
	router := http.NewServeMux()
	router.HandleFunc(
		"GET /",
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("Hello from the GO http server"))
		},
	)

	// create the HTTP server
	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}

	slog.Info(
		"starting HTTP server",
		slog.String("addr", cfg.HTTPServer.Addr),
	)

	// Listen for a termination signals from the OS
	shutdownSignal := make(chan os.Signal, 1)

	signal.Notify(
		shutdownSignal,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	// start the server in a separate goroutine so the main
	// goroutine can wait for shutdown signals.
	go func() {
		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			slog.Error(
				"server closed unexpectedly",
				slog.Any("error", err),
			)

			os.Exit(1)
		}
	}()

	// Block until shutdown Signal is received
	receivedSignal := <-shutdownSignal

	slog.Info(
		"shutdown signal received",
		slog.Any("signal", receivedSignal.String()),
	)

	// Dive active requests time to complete before forcefully terminating the process.
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error(
			"gracefully shutdown failed",
			slog.Any("error", err),
		)

		os.Exit(1)
	}

	slog.Info("server stopped gracefully")
}
