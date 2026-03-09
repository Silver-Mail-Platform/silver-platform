package intake

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewMuxRoutesEventsPath(t *testing.T) {
	mux := NewMux()

	req := httptest.NewRequest(http.MethodPost, EventsPath, strings.NewReader(`{"timestamp":"2026-03-05T10:30:45Z"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusAccepted)
	}
	if strings.TrimSpace(rr.Body.String()) != `{"ok":true}` {
		t.Fatalf("body = %q, want %q", strings.TrimSpace(rr.Body.String()), `{"ok":true}`)
	}
}

func TestNewMuxUnknownPath(t *testing.T) {
	mux := NewMux()

	req := httptest.NewRequest(http.MethodPost, "/unknown", strings.NewReader(`{"timestamp":"2026-03-05T10:30:45Z"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
}
