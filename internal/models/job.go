package models

import (
	"encoding/json"
	"time"
)

type JobState uint8

const (
	StateScheduled JobState = iota + 1
	StateReady
	StateProcessing
	StateCompleted
	StateDead
	StateCancelled
)

type Job struct {
	ID      string          `json:"id"` //for idempotent producer
	Payload json.RawMessage `json:"payload"`

	Priority int      `json:"priority"` // 0 > 1 > 2
	State    JobState `json:"state"`
	Queue    string   `json:"queue"`

	Fingerprint string `json:"fingerprint"`

	Attempt     int `json:"attempt"`
	MaxAttempts int `json:"max_attempts"`

	CreatedAt        time.Time     `json:"created_at"`
	ScheduledAt      time.Time     `json:"scheduled_at"`
	ExecutionTimeout time.Duration `json:"execution_timeout"`
	LeaseDuration    time.Duration `json:"lease_duration"`

	//ownership of worker over job iff there
	WorkerID       string    `json:"worker_id"`
	LeaseEpoch     uint64    `json:"lease_epoch"`
	LeaseExpiresAt time.Time `json:"lease_expires_at"`
	AttemptStarted time.Time `json:"attempt_started"`

	LastError string `json:"last_error"`

	CompletedAt *time.Time `json:"completed_at"`
}
