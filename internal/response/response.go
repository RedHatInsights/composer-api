package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Body any `json:"body"`
}

type ErrorDetails struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorDetails `json:"error"`
}

func JSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: ErrorDetails{Code: status, Message: msg}})
}

func JSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{Body: body})
}
