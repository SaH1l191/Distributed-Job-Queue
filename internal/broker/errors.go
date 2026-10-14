package broker

import "errors"

var (
	ErrIDConflict    = errors.New("job id already exists with different payload")
	ErrQueueFull     = errors.New("queue is full")
	ErrNotFound      = errors.New("job not found")
	ErrStale         = errors.New("stale lease epoch")
	ErrNotProcessing = errors.New("job is not processing")
	ErrLeaseExpired  = errors.New("lease already expired")
	ErrTerminal      = errors.New("job already in terminal state")
)
