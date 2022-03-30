package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/warrenhodg/urlmon/healthcheck"
	"github.com/warrenhodg/urlmon/httpserver"
	"github.com/warrenhodg/urlmon/urlchecker"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// waitForTermination waits for a signal
// to terminate the app
func waitForTermination(ctx context.Context, logger *zap.Logger) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		logger.Info("shutdown request due to context done")
		return
	case <-sigs:
		logger.Info("shutdown request due to signal received")
		return
	}
}

// failReadiness fails the readiness check for a specified duration
// then returns to complete other shutdown logic
func failReadiness(ctx context.Context, logger *zap.Logger, duration time.Duration) {
	logger.
		With(zap.Duration("duration", duration)).
		Info("failing readiness as part of normal termination sequence")

	healthcheck.SetReady(false)
	t := time.NewTimer(duration)
	defer t.Stop()

	select {
	case <-ctx.Done():
		return
	case <-t.C:
		return
	}
}

func setupMetrics() {
	http.DefaultServeMux.Handle("/metrics", promhttp.Handler())
}

func run(logger *zap.Logger) error {
	logger.Info("app start")
	defer logger.Info("app graceful stop")

	ctx := context.Background()

	cfg := GetConfig()

	svrOptions := httpserver.
		DefaultOptions().
		ListenAddr(cfg.listenAddr)
	svr := httpserver.New(logger, svrOptions)

	healthcheck.Init(logger)

	setupMetrics()

	ucOptions := urlchecker.
		DefaultOptions().
		Urls(cfg.urls).
		CheckPeriod(cfg.checkPeriod).
		Workers(cfg.workers)
	uc := urlchecker.New(logger, ucOptions)

	err := svr.Run()
	if err != nil {
		return fmt.Errorf("unable to run httpserver: %w", err)
	}
	defer svr.Shutdown(ctx)

	uc.Run(ctx)
	defer uc.Shutdown()

	waitForTermination(ctx, logger)

	failReadiness(ctx, logger, cfg.shutdownReadyFailDuration)

	return nil
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("unable to create logger: %s", err.Error())
		return
	}

	err = run(logger)
	if err != nil {
		logger.With(zap.Error(err)).Error("an error occurred")
	}
}
