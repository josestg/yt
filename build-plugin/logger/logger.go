package logger

import (
	"context"
	"io"
	"log/slog"
)

type logContextKey struct{}

func FromContext(ctx context.Context) *slog.Logger {
	log, ok := ctx.Value(logContextKey{}).(*slog.Logger)
	if !ok {
		panic("missing logger from context")
	}
	return log
}

func WithContext(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, logContextKey{}, log)
}

func New(w io.Writer) *slog.Logger {
	opts := &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				nano := a.Value.Time().UnixMicro()
				return slog.Int64(a.Key, nano)
			}
			return a
		},
	}
	return slog.New(slog.NewTextHandler(w, opts))
}
