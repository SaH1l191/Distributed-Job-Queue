package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"jobqueue/internal/app"
	"jobqueue/pkg/logger"
)

func main() {
	l := logger.InitLogger()
	ctx := context.Background()

	a, err := app.New(ctx, l)
	if err != nil {
		l.Error("failed to create application", "error", err)
		os.Exit(1)
	}

	if err := a.Start(ctx); err != nil {
		l.Error("failed to start application", "error", err)
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := a.Stop(ctx); err != nil {
		l.Error("failed to stop application", "error", err)
	}
}
