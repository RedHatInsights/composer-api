package response

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

type Error struct {
	StatusCode int    `json:"-"`
	Status     string `json:"status"`
	Reason     string `json:"reason"`
}

func (e *Error) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("%d %s: %s", e.StatusCode, e.Status, e.Reason)
	}
	return fmt.Sprintf("%d %s", e.StatusCode, e.Status)
}

func (e *Error) WithReasonStr(reason string) *Error {
	e.Reason = reason
	return e
}

func (e *Error) WithReasonErr(err error) *Error {
	e.Reason = err.Error()
	return e
}

func NewError(statusCode int) *Error {
	return &Error{
		StatusCode: statusCode,
		Status:     http.StatusText(statusCode),
	}
}

func BadRequest() *Error         { return NewError(http.StatusBadRequest) }
func Unauthorized() *Error       { return NewError(http.StatusUnauthorized) }
func Forbidden() *Error          { return NewError(http.StatusForbidden) }
func NotFound() *Error           { return NewError(http.StatusNotFound) }
func MethodNotAllowed() *Error   { return NewError(http.StatusMethodNotAllowed) }
func Conflict() *Error           { return NewError(http.StatusConflict) }
func UnprocessableEntity() *Error { return NewError(http.StatusUnprocessableEntity) }
func TooManyRequests() *Error    { return NewError(http.StatusTooManyRequests) }
func InternalServerError() *Error { return NewError(http.StatusInternalServerError) }
func ServiceUnavailable() *Error { return NewError(http.StatusServiceUnavailable) }

func ToError(v any) *Error {
	switch val := v.(type) {
	case *Error:
		return val
	case error:
		var httpErr *Error
		if errors.As(val, &httpErr) {
			return httpErr
		}
		slog.Error("internal error", "error", val)
		return InternalServerError()
	case string:
		slog.Error("internal error", "error", val)
		return InternalServerError()
	default:
		slog.Error("internal error", "error", fmt.Sprintf("%v", val))
		return InternalServerError()
	}
}

func WriteError(w http.ResponseWriter, err *Error) {
	writeJSON(w, err.StatusCode, err)
}
