package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"jobqueue/internal/api"
	"jobqueue/internal/broker"
	"jobqueue/internal/clock"
	"jobqueue/internal/coord"
	"jobqueue/pkg/logger"
)

func main() {
	l := logger.InitLogger()
	ctx := context.Background()

	cfg := broker.Config{
		QueueCapacity:           1000,
		DefaultMaxAttempts:      5,
		DefaultLeaseDuration:    30 * time.Second,
		DefaultExecutionTimeout: 5 * time.Minute,
	}

	clk := clock.NewRealClock()
	b := broker.NewBroker(clk, cfg)
	l.Info("Broker initialized")

	scheduler := coord.NewScheduler(b, clk)
	go scheduler.Run()
	l.Info("Scheduler started")

	a, err := api.New(ctx, l ,b)
	if err != nil {
		l.Error("failed to create application", "error", err)
		os.Exit(1)
	}

	if err := a.Start(ctx); err != nil {
		l.Error("failed to start application", "error", err)
		os.Exit(1)
	}

	//buffered chan blocking op, receiving ctrl+c,stopping
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit //receive something into chan and til then wait here

	if err := a.Stop(ctx); err != nil {
		l.Error("failed to stop application", "error", err)
	}
}
