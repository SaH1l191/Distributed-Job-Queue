package api

import (
	"encoding/json"
	"jobqueue/internal/broker"
	"net/http"
)

// HandleSubmit is the POST /jobs producer endpoint.
// It lives in package api (same folder as routes.go) so no import cycle is possible.
func HandleSubmit(b *broker.Broker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SubmitRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		job, created, err := b.Submit(
			broker.SubmitInput{
				ID:               req.ID,
				Payload:          req.Payload,
				Priority:         req.Priority,
				ScheduledAt:      req.ScheduledAt,
				MaxAttempts:      req.MaxAttempts,
				ExecutionTimeout: req.ExecutionTimeout,
				LeaseDuration:    req.LeaseDuration,
			},
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if created {
			w.WriteHeader(http.StatusCreated)
		}
		json.NewEncoder(w).Encode(SubmitResponse{
			ID:        job.ID,
			State:     job.State.String(),
			CreatedAt: job.CreatedAt,
		})

	}
}
