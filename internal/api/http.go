package api

import (
	"context"
	"fmt"
	broker "jobqueue/internal/broker"
	"log/slog"
	"net/http"
)

type App struct {
	serv    *http.Server
	context context.Context
	logger  *slog.Logger
	broker  *broker.Broker
}

func New(ctx context.Context, logger *slog.Logger, broker *broker.Broker) (*App, error) {
	return &App{
		context: ctx,
		logger:  logger,
		broker:  broker,
	}, nil
}

func (a *App) Start(ctx context.Context) error {
	a.logger.Info("Starting Application")

	mux := http.NewServeMux()
	a.serv = &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	a.logger.Info("Application server initialized")

	go func() {
		a.logger.Info("Starting HTTP server", "addr", ":8080")
		if err := a.serv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error("HTTP server failed", "error", err)
		}
	}()
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	a.logger.Info("Stopping Application")

	if err := a.serv.Shutdown(ctx); err != nil {
		return fmt.Errorf("HTTP server shutdown failed: %w", err)
	}

	a.logger.Info("Application stopped cleanly")
	return nil
}
