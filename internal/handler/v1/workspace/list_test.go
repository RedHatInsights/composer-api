package workspace

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RedHatInsights/composer-api/internal/response"
)

func TestList(t *testing.T) {
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/workspaces", nil)
	rec := httptest.NewRecorder()

	h := New()
	h.List(rec, req)

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

	workspaces, ok := body["workspaces"].([]any)
	if !ok {
		t.Fatalf("expected workspaces to be a list, got %T", body["workspaces"])
	}

	if len(workspaces) != 0 {
		t.Errorf("expected empty workspaces list, got %d items", len(workspaces))
	}
}
