package response

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteJSON(context.Background(), rec, http.StatusOK, map[string]string{"key": "value"})

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type %q, got %q", "application/json", ct)
	}

	var resp Response
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	body, ok := resp.Body.(map[string]any)
	if !ok {
		t.Fatalf("expected body to be a map, got %T", resp.Body)
	}

	if body["key"] != "value" {
		t.Errorf("expected key %q, got %q", "value", body["key"])
	}
}

func TestWriteJSON_CustomStatus(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteJSON(context.Background(), rec, http.StatusNotFound, map[string]string{"error": "not found"})

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestWriteJSON_NilBody(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteJSON(context.Background(), rec, http.StatusOK, nil)

	var resp Response
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Body != nil {
		t.Errorf("expected nil body, got %v", resp.Body)
	}
}
