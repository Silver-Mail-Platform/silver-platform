package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Silver-Mail-Platform/super-platform/intake"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ingestor := intake.NewLoggingIngestor(logger)
	mux := intake.NewMux(ingestor, logger)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	logger.Info("starting intake server", "addr", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server exited with error", "error", err)
		os.Exit(1)
	}
}
