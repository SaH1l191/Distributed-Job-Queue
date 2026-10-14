package api

import (
	"encoding/json"
	"time"
)

type SubmitRequest struct {
	ID               string          `json:"id"`
	Payload          json.RawMessage `json:"payload"`
	Priority         int             `json:"priority"`
	Queue            string          `json:"queue"`
	ScheduledAt      time.Time       `json:"scheduled_at"`
	MaxAttempts      int             `json:"max_attempts"`
	ExecutionTimeout time.Duration   `json:"execution_timeout"`
	LeaseDuration    time.Duration   `json:"lease_duration"`
}

type SubmitResponse struct {
	ID        string    `json:"id"`
	State     string    `json:"state"`
	Created   bool      `json:"created"`
	Fingerprint string  `json:"fingerprint"`
	CreatedAt time.Time `json:"created_at"`
}



