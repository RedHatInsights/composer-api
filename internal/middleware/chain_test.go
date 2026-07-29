package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChain_AppliesInOrder(t *testing.T) {
	var order []string

	mw := func(name string) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, name)
				next.ServeHTTP(w, r)
			})
		}
	}

	handler := Chain(mw("first"), mw("second"), mw("third"))(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "handler")
		}),
	)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	handler.ServeHTTP(httptest.NewRecorder(), req)

	expected := []string{"first", "second", "third", "handler"}
	if len(order) != len(expected) {
		t.Fatalf("expected %d calls, got %d", len(expected), len(order))
	}
	for i, v := range expected {
		if order[i] != v {
			t.Errorf("position %d: expected %q, got %q", i, v, order[i])
		}
	}
}

func TestChain_Empty(t *testing.T) {
	called := false
	handler := Chain()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	handler.ServeHTTP(httptest.NewRecorder(), req)

	if !called {
		t.Error("expected handler to be called with empty chain")
	}
}
