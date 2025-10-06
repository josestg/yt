package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/josestg/yt/build-plugin/cache"
	"github.com/josestg/yt/build-plugin/fib"
	"github.com/josestg/yt/build-plugin/logger"
)

var addr = "localhost:8080"
var gracePeriod = 15 * time.Second

func main() {
	flag.StringVar(&addr, "addr", addr, "http server address")
	flag.DurationVar(&gracePeriod, "grace-period", gracePeriod, "wait timeout before executing force shutdown")
	flag.Parse()
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "exit: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	log := logger.New(os.Stdout)
	log.Info("app started", "pid", os.Getpid())
	defer log.Info("app stopped")

	mux := http.NewServeMux()
	mux.Handle("GET /api/v1/fib/{n}", fib.NewHandler(cache.NopCache, 0))

	srv := http.Server{
		Addr:    addr,
		Handler: mux,
		BaseContext: func(listener net.Listener) context.Context {
			return logger.WithContext(context.Background(), log)
		},
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
