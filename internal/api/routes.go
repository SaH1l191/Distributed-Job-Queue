package api

import (
	"jobqueue/internal/broker"
	"log/slog"
	"net/http"
	"time"
)


func SetupRoutes(mux *http.ServeMux, b *broker.Broker) {
	// Producer
	mux.HandleFunc("POST /jobs", HandleSubmit(b))
	
	// mux.HandleFunc("GET /jobs/{id}", handlers.HandleGetJob(b))
	// mux.HandleFunc("DELETE /jobs/{id}", handlers.HandleCancel(b))

	// Worker
	// mux.HandleFunc("POST /claim", handlers.HandleClaim(b))
	// mux.HandleFunc("POST /jobs/{id}/renew", handlers.HandleRenew(b))
	// mux.HandleFunc("POST /jobs/{id}/ack", handlers.HandleAck(b))
	// mux.HandleFunc("POST /jobs/{id}/nack", handlers.HandleNack(b))
}


func loggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Info("HTTP request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start),
			"remote", r.RemoteAddr,
		)
	})
}