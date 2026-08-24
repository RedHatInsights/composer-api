package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestID_AddsToContext(t *testing.T) {
	var capturedID string

	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID = GetRequestID(r.Context())
	}))

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if capturedID == "" {
		t.Error("expected request ID in context, got empty string")
	}
}

func TestRequestID_UniquePerRequest(t *testing.T) {
	var ids []string

	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ids = append(ids, GetRequestID(r.Context()))
	}))

	for range 3 {
		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}

	if ids[0] == ids[1] || ids[1] == ids[2] {
		t.Errorf("expected unique IDs, got %v", ids)
	}
}

func TestRequestID_NoResponseHeader(t *testing.T) {
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if h := rec.Header().Get("X-Request-ID"); h != "" {
		t.Errorf("expected no X-Request-ID header, got %q", h)
	}
}

func TestGetRequestID_EmptyContext(t *testing.T) {
	id := GetRequestID(context.Background())
	if id != "" {
		t.Errorf("expected empty string, got %q", id)
	}
}
