package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogging_SuccessfulRequest(t *testing.T) {
	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestLogging_ErrorRequest(t *testing.T) {
	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal error"))
	}))

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/fail", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}

	if body := rec.Body.String(); body != "internal error" {
		t.Errorf("expected body %q, got %q", "internal error", body)
	}
}

func TestLogging_DefaultStatusCode(t *testing.T) {
	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("no explicit status"))
	}))

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/default", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected default status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestWrappedResponseWriter_CapturesBodyOnErrors(t *testing.T) {
	w := &wrappedResponseWriter{
		ResponseWriter: httptest.NewRecorder(),
		statusCode:     http.StatusInternalServerError,
	}
	_, _ = w.Write([]byte("error details"))

	if w.body.String() != "error details" {
		t.Errorf("expected captured body %q, got %q", "error details", w.body.String())
	}
}

func TestWrappedResponseWriter_SkipsBodyOnSuccess(t *testing.T) {
	w := &wrappedResponseWriter{
		ResponseWriter: httptest.NewRecorder(),
		statusCode:     http.StatusOK,
	}
	_, _ = w.Write([]byte("success body"))

	if w.body.Len() != 0 {
		t.Errorf("expected no body captured for success, got %q", w.body.String())
	}
}

func TestLogging_PassesThroughResponse(t *testing.T) {
	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Custom", "test")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("created"))
	}))

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/create", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	if rec.Body.String() != "created" {
		t.Errorf("expected body %q, got %q", "created", rec.Body.String())
	}

	if rec.Header().Get("X-Custom") != "test" {
		t.Errorf("expected X-Custom header %q, got %q", "test", rec.Header().Get("X-Custom"))
	}
}
