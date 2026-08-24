package probe

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RedHatInsights/composer-api/internal/response"
)

func TestPing(t *testing.T) {
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()

	h := New()
	h.Ping(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type %q, got %q", "application/json", ct)
	}

	var resp response.Response
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	body, ok := resp.Body.(map[string]any)
	if !ok {
		t.Fatalf("expected body to be a map, got %T", resp.Body)
	}

	if ping := body["ping"]; ping != "pong" {
		t.Errorf("expected ping %q, got %q", "pong", ping)
	}
}
