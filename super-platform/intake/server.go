package intake

import (
	"context"
	"log/slog"
	"net/http"
)

type LoggingIngestor struct {
	logger *slog.Logger
}

func NewLoggingIngestor(logger *slog.Logger) *LoggingIngestor {
	if logger == nil {
		logger = slog.Default()
	}
	return &LoggingIngestor{logger: logger}
}

func (i *LoggingIngestor) Ingest(_ context.Context, body map[string]any) error {
	i.logger.Info("ingested event", "body", body)
	return nil
}

func NewMux(ingestor Ingestor, logger *slog.Logger) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/v1/silver/events", NewHandler(ingestor, logger))
	return mux
}

