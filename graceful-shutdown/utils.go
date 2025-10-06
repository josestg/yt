package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
)

func replyJSON(w http.ResponseWriter, r *http.Request, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log := LogFromContext(r.Context())
		log.Error("cannot write response", "error", err.Error())
	}
}

func readJSONBody(w http.ResponseWriter, r *http.Request, dest any) bool {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		log := LogFromContext(r.Context())
		log.Warn("cannot read request body", "error", err.Error())
		replyJSON(w, r, http.StatusBadRequest, map[string]any{"message": "Request Invalid"})
		return false
	}
	return true
}

type logContextKey struct{}

func LogFromContext(ctx context.Context) *slog.Logger {
	log, ok := ctx.Value(logContextKey{}).(*slog.Logger)
	if !ok {
		panic("missing logger from context")
	}
	return log
}

func newLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				nano := a.Value.Time().UnixMicro()
				return slog.Int64(a.Key, nano)
			}
			return a
		},
	}

	return slog.New(slog.NewTextHandler(os.Stdout, opts))
}
