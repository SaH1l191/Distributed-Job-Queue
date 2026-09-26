package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"jobqueue/internal/api"
	"jobqueue/pkg/logger"
)

func main() {
	l := logger.InitLogger()
	ctx := context.Background()

	a, err := api.New(ctx, l)
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
