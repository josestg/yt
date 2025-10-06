package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var addr = "localhost:8080"
var delay = 30 * time.Second
var gracePeriod = 15 * time.Second

func main() {
	flag.StringVar(&addr, "addr", addr, "http server address")
	flag.DurationVar(&delay, "delay", delay, "delay simulation")
	flag.DurationVar(&gracePeriod, "grace-period", delay, "wait timeout before executing force shutdown")
	flag.Parse()
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "exit: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	log := newLogger()
	log.Info("app started", "pid", os.Getpid())
	defer log.Info("app stopped")

	mux := http.NewServeMux()
	mux.Handle("POST /api/v1/orders", orderHandler())

	srv := http.Server{
		Addr:    addr,
		Handler: handlerWithLogContext(mux, log),
	}

	var wg sync.WaitGroup
	wg.Go(func() {
		log.Info("server is listening", "addr", srv.Addr)
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Warn("unexpected error from listen and serve", "error", err.Error())
		}
	})

	shutdownRequest := make(chan os.Signal, 1)
	signal.Notify(shutdownRequest, syscall.SIGTERM)
	wg.Go(func() {
		sig := <-shutdownRequest
		log.Info("shutdown request received", "signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), gracePeriod)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			log.Warn("shutdown failed, continue with force shutdown", "error", err.Error())
			err = srv.Close()
			if err != nil {
				log.Error("force shutdown also failed", "error", err.Error())
			}
			return
		}
		log.Info("server shutdown gracefully")
	})

	wg.Wait()
	return nil
}

func orderHandler() http.HandlerFunc {
	type Req struct {
		ProductID int64  `json:"product_id"`
		Price     string `json:"price"`
		Count     int64  `json:"count"`
	}

	type Res struct {
		OrderID   int64  `json:"order_id"`
		Status    string `json:"status"`
		Timestamp int64  `json:"timestamp"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		var req Req
		if ok := readJSONBody(w, r, &req); !ok {
			return
		}

		ctx := r.Context()
		log := LogFromContext(ctx)

		log.Info("order request received", "request_body", req)
		select {
		case <-ctx.Done():
			log.Warn("context done", "error", ctx.Err())
			return
		case <-time.After(delay): // simulate latency
			// do something important....
			// fake result.
			res := Res{
				OrderID:   rand.Int64N(100_000_000),
				Status:    "completed",
				Timestamp: time.Now().UnixMilli(),
			}
			replyJSON(w, r, http.StatusOK, res)
			log.Info("order request completed", "response_body", res, "elapsed", time.Since(started).String())
		}
	}
}
