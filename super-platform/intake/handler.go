package intake

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"time"
)

type Ingestor interface {
	Ingest(ctx context.Context, body map[string]any) error
}

type Handler struct {
	ingestor Ingestor
	logger   *slog.Logger
}

func NewHandler(ingestor Ingestor, logger *slog.Logger) *Handler {
	if ingestor == nil {
		panic("ingestor cannot be nil")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{
		ingestor: ingestor,
		logger:   logger,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	if !isJSONContentType(r.Header.Get("Content-Type")) {
		writeJSON(w, http.StatusUnsupportedMediaType, map[string]string{"error": "content type must be application/json"})
		return
	}

	var body map[string]any
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, io.EOF) {
			writeJSON(w, status, map[string]string{"error": "request body is required"})
			return
		}
		writeJSON(w, status, map[string]string{"error": "invalid JSON"})
		return
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "request body must contain only one JSON object"})
		return
	}

	tsRaw, ok := body["timestamp"]
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "timestamp is required"})
		return
	}
	ts, ok := tsRaw.(string)
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "timestamp must be a string"})
		return
	}
	if _, err := time.Parse(time.RFC3339, ts); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "timestamp must be RFC3339"})
		return
	}

	if err := h.ingestor.Ingest(r.Context(), body); err != nil {
		h.logger.Error("ingest failed", "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	h.logger.Info("event accepted", "timestamp", ts)
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "accepted"})
}

func isJSONContentType(contentType string) bool {
	if contentType == "" {
		return false
	}
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}
	return mediaType == "application/json"
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

