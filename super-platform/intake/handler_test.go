package intake

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerResponses(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		contentType   string
		body          string
		wantStatus    int
		wantBody      string
		wantAllowPost bool
	}{
		{
			name:          "method not allowed",
			method:        http.MethodGet,
			contentType:   "application/json",
			body:          "{}",
			wantStatus:    http.StatusMethodNotAllowed,
			wantBody:      `{"error":"method not allowed"}`,
			wantAllowPost: true,
		},
		{
			name:        "unsupported media type",
			method:      http.MethodPost,
			contentType: "text/plain",
			body:        "{}",
			wantStatus:  http.StatusUnsupportedMediaType,
			wantBody:    `{"error":"content type must be application/json"}`,
		},
		{
			name:        "empty body",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        "",
			wantStatus:  http.StatusBadRequest,
			wantBody:    `{"error":"request body is required"}`,
		},
		{
			name:        "invalid json",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        `{"timestamp":`,
			wantStatus:  http.StatusBadRequest,
			wantBody:    `{"error":"invalid JSON"}`,
		},
		{
			name:        "multiple json objects",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        `{"timestamp":"2026-03-05T10:30:45Z"}{"another":1}`,
			wantStatus:  http.StatusBadRequest,
			wantBody:    `{"error":"request body must contain only one JSON object"}`,
		},
		{
			name:        "timestamp missing",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        `{"event":"scan_complete"}`,
			wantStatus:  http.StatusBadRequest,
			wantBody:    `{"error":"timestamp is required"}`,
		},
		{
			name:        "timestamp not string",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        `{"timestamp":123}`,
			wantStatus:  http.StatusBadRequest,
			wantBody:    `{"error":"timestamp must be a string"}`,
		},
		{
			name:        "timestamp not rfc3339",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        `{"timestamp":"2026/03/05 10:30:45"}`,
			wantStatus:  http.StatusBadRequest,
			wantBody:    `{"error":"timestamp must be RFC3339"}`,
		},
		{
			name:        "accepted",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        `{"timestamp":"2026-03-05T10:30:45Z","event":"scan_complete"}`,
			wantStatus:  http.StatusAccepted,
			wantBody:    `{"ok":true}`,
		},
		{
			name:        "content type with charset",
			method:      http.MethodPost,
			contentType: "application/json; charset=utf-8",
			body:        `{"timestamp":"2026-03-05T10:30:45Z"}`,
			wantStatus:  http.StatusAccepted,
			wantBody:    `{"ok":true}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := NewHandler()

			req := httptest.NewRequest(tc.method, EventsPath, strings.NewReader(tc.body))
			if tc.contentType != "" {
				req.Header.Set("Content-Type", tc.contentType)
			}
			rr := httptest.NewRecorder()

			h.ServeHTTP(rr, req)

			if rr.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d", rr.Code, tc.wantStatus)
			}
			if strings.TrimSpace(rr.Body.String()) != tc.wantBody {
				t.Fatalf("body = %q, want %q", strings.TrimSpace(rr.Body.String()), tc.wantBody)
			}
			if tc.wantAllowPost && rr.Header().Get("Allow") != http.MethodPost {
				t.Fatalf("Allow header = %q, want %q", rr.Header().Get("Allow"), http.MethodPost)
			}
		})
	}
}
