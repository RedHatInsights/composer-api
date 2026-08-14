package response

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Response struct {
	Body any `json:"body"`
}

func JSON(w http.ResponseWriter, status int, body any) {
	writeJSON(w, status, Response{Body: body})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	data, err := json.Marshal(v)
	if err != nil {
		slog.Error("failed to marshal response", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(`{"error":{"code":500,"message":"internal server error"}}`)); err != nil {
			slog.Error("failed to write error response", "error", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error("failed to write response", "error", err)
	}
}
