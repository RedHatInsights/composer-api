package middleware

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBodySizeLimit_WithinLimit(t *testing.T) {
	handler := BodySizeLimit(1024)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("unexpected error reading body: %v", err)
		}
		if len(body) != 10 {
			t.Errorf("expected body length 10, got %d", len(body))
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/test", strings.NewReader("0123456789"))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestBodySizeLimit_ExceedsLimit(t *testing.T) {
	handler := BodySizeLimit(5)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.ReadAll(r.Body)
		if err == nil {
			t.Fatal("expected error reading oversized body, got nil")
		}
		w.WriteHeader(http.StatusRequestEntityTooLarge)
	}))

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/test", strings.NewReader("0123456789"))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("expected status %d, got %d", http.StatusRequestEntityTooLarge, rec.Code)
	}
}
