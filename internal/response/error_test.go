package response

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
)

func TestFactoryFunctions(t *testing.T) {
	tests := []struct {
		name       string
		factory    func() *Error
		wantCode   int
		wantStatus string
	}{
		{"BadRequest", BadRequest, 400, "Bad Request"},
		{"Unauthorized", Unauthorized, 401, "Unauthorized"},
		{"Forbidden", Forbidden, 403, "Forbidden"},
		{"NotFound", NotFound, 404, "Not Found"},
		{"MethodNotAllowed", MethodNotAllowed, 405, "Method Not Allowed"},
		{"Conflict", Conflict, 409, "Conflict"},
		{"UnprocessableEntity", UnprocessableEntity, 422, "Unprocessable Entity"},
		{"TooManyRequests", TooManyRequests, 429, "Too Many Requests"},
		{"InternalServerError", InternalServerError, 500, "Internal Server Error"},
		{"ServiceUnavailable", ServiceUnavailable, 503, "Service Unavailable"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.factory()
			if err.StatusCode != tt.wantCode {
				t.Errorf("StatusCode = %d, want %d", err.StatusCode, tt.wantCode)
			}
			if err.Status != tt.wantStatus {
				t.Errorf("Status = %q, want %q", err.Status, tt.wantStatus)
			}
			if err.Reason != "" {
				t.Errorf("Reason = %q, want empty", err.Reason)
			}
		})
	}
}

func TestWithReasonStr(t *testing.T) {
	err := NotFound().WithReasonStr("user not found")
	if err.Reason != "user not found" {
		t.Errorf("Reason = %q, want %q", err.Reason, "user not found")
	}
	if err.StatusCode != 404 {
		t.Errorf("StatusCode = %d, want 404", err.StatusCode)
	}
}

func TestWithReasonErr(t *testing.T) {
	cause := errors.New("connection refused")
	err := ServiceUnavailable().WithReasonErr(cause)
	if err.Reason != "connection refused" {
		t.Errorf("Reason = %q, want %q", err.Reason, "connection refused")
	}
}

func TestErrorInterface(t *testing.T) {
	var _ error = NotFound()

	err := NotFound().WithReasonStr("page missing")
	want := "404 Not Found: page missing"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}

	errNoReason := BadRequest()
	wantNoReason := "400 Bad Request"
	if got := errNoReason.Error(); got != wantNoReason {
		t.Errorf("Error() = %q, want %q", got, wantNoReason)
	}
}

func TestToError(t *testing.T) {
	t.Run("from *Error", func(t *testing.T) {
		original := NotFound().WithReasonStr("gone")
		result := ToError(context.Background(), original)
		if result != original {
			t.Error("expected same pointer")
		}
	})

	t.Run("from wrapped *Error", func(t *testing.T) {
		original := Forbidden().WithReasonStr("no access")
		wrapped := fmt.Errorf("wrapped: %w", original)
		result := ToError(context.Background(), wrapped)
		if result.StatusCode != 403 {
			t.Errorf("StatusCode = %d, want 403", result.StatusCode)
		}
	})

	t.Run("from plain error", func(t *testing.T) {
		result := ToError(context.Background(), errors.New("something broke"))
		if result.StatusCode != 500 {
			t.Errorf("StatusCode = %d, want 500", result.StatusCode)
		}
		if result.Reason != "" {
			t.Errorf("Reason = %q, want empty (internal details not leaked)", result.Reason)
		}
	})

	t.Run("from string", func(t *testing.T) {
		result := ToError(context.Background(), "bad thing")
		if result.StatusCode != 500 {
			t.Errorf("StatusCode = %d, want 500", result.StatusCode)
		}
		if result.Reason != "" {
			t.Errorf("Reason = %q, want empty (internal details not leaked)", result.Reason)
		}
	})

	t.Run("from other type", func(t *testing.T) {
		result := ToError(context.Background(), 42)
		if result.StatusCode != 500 {
			t.Errorf("StatusCode = %d, want 500", result.StatusCode)
		}
		if result.Reason != "" {
			t.Errorf("Reason = %q, want empty (internal details not leaked)", result.Reason)
		}
	})
}

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()
	WriteError(context.Background(), w, NotFound().WithReasonStr("workspace not found"))

	if w.Code != 404 {
		t.Errorf("status code = %d, want 404", w.Code)
	}

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/json")
	}

	var body struct {
		Status string `json:"status"`
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body.Status != "Not Found" {
		t.Errorf("status = %q, want %q", body.Status, "Not Found")
	}
	if body.Reason != "workspace not found" {
		t.Errorf("reason = %q, want %q", body.Reason, "workspace not found")
	}
}
