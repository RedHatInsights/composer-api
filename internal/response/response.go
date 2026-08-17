package response

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type Response struct {
	Body any `json:"body"`
}

func WriteJSON(ctx context.Context, w http.ResponseWriter, status int, body any) {
	encodeJSON(ctx, w, status, Response{Body: body})
}

func encodeJSON(ctx context.Context, w http.ResponseWriter, status int, v any) {
	data, err := json.Marshal(v)
	if err != nil {
		slog.ErrorContext(ctx, "failed to marshal response", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(`{"error":{"code":500,"message":"internal server error"}}`)); err != nil {
			slog.ErrorContext(ctx, "failed to write error response", "error", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.ErrorContext(ctx, "failed to write response", "error", err)
	}
}
