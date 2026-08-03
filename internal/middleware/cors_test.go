package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newCORSHandler(origins []string, maxAge int) http.Handler {
	return CORS(origins, maxAge)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
}

func TestCORS_AllowedOrigin(t *testing.T) {
	handler := newCORSHandler([]string{"https://example.com"}, 3600)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://example.com" {
		t.Errorf("expected Allow-Origin %q, got %q", "https://example.com", got)
	}
	if got := rec.Header().Get("Vary"); got != "Origin" {
		t.Errorf("expected Vary %q, got %q", "Origin", got)
	}
}

func TestCORS_DisallowedOrigin(t *testing.T) {
	handler := newCORSHandler([]string{"https://example.com"}, 3600)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "https://evil.com")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("expected no Allow-Origin header, got %q", got)
	}
}

func TestCORS_PreflightAllowed(t *testing.T) {
	handler := newCORSHandler([]string{"https://example.com"}, 3600)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodOptions, "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Methods"); got != corsAllowedMethods {
		t.Errorf("expected Allow-Methods %q, got %q", corsAllowedMethods, got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Headers"); got != corsAllowedHeaders {
		t.Errorf("expected Allow-Headers %q, got %q", corsAllowedHeaders, got)
	}
	if got := rec.Header().Get("Access-Control-Max-Age"); got != "3600" {
		t.Errorf("expected Max-Age %q, got %q", "3600", got)
	}
}

func TestCORS_PreflightDisallowed(t *testing.T) {
	handler := newCORSHandler([]string{"https://example.com"}, 3600)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodOptions, "/test", nil)
	req.Header.Set("Origin", "https://evil.com")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Methods"); got != "" {
		t.Errorf("expected no Allow-Methods header, got %q", got)
	}
}

func TestCORS_NonCORSOptions(t *testing.T) {
	called := false
	handler := CORS([]string{"https://example.com"}, 3600)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequestWithContext(context.Background(), http.MethodOptions, "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !called {
		t.Error("expected handler to be called for non-CORS OPTIONS request")
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("expected no Allow-Origin header, got %q", got)
	}
}

func TestCORS_NoOrigin(t *testing.T) {
	handler := newCORSHandler([]string{"https://example.com"}, 3600)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("expected no Allow-Origin header, got %q", got)
	}
}

func TestCORS_WildcardOrigin(t *testing.T) {
	handler := newCORSHandler([]string{"*"}, 3600)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "https://anything.com")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://anything.com" {
		t.Errorf("expected Allow-Origin %q, got %q", "https://anything.com", got)
	}
}
