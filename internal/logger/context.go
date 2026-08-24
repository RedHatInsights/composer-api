package logger

import (
	"context"
	"log/slog"
)

type ctxKey struct{}

// ContextHandler is a slog.Handler that automatically includes
// key-value pairs stored in the context via AddContextValue.
type ContextHandler struct {
	slog.Handler
}

func (c ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(ctxKey{}).([]slog.Attr); ok {
		r.AddAttrs(attrs...)
	}
	return c.Handler.Handle(ctx, r)
}

func (c ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return ContextHandler{Handler: c.Handler.WithAttrs(attrs)}
}

func (c ContextHandler) WithGroup(name string) slog.Handler {
	return ContextHandler{Handler: c.Handler.WithGroup(name)}
}

// AddContextValue returns a new context that includes the given
// key-value pair. Any slog call using the returned context will
// automatically log this value.
func AddContextValue(parent context.Context, key string, value any) context.Context {
	attr := slog.Any(key, value)

	if existing, ok := parent.Value(ctxKey{}).([]slog.Attr); ok {
		attrs := make([]slog.Attr, 0, len(existing)+1)
		attrs = append(attrs, existing...)
		attrs = append(attrs, attr)
		return context.WithValue(parent, ctxKey{}, attrs)
	}

	return context.WithValue(parent, ctxKey{}, []slog.Attr{attr})
}
