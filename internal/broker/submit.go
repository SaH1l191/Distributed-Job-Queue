package broker

import "jobqueue/internal/models"

func (b *Broker) Submit(in SubmitInput) (models.Job, bool, error) {
	now := b.clock.Now()
	b.mu.Lock()
	defer b.mu.Unlock()

	if in.MaxAttempts <= 0 {
		in.MaxAttempts = b.cfg.DefaultMaxAttempts
	}
	if in.LeaseDuration <= 0 {
		in.LeaseDuration = b.cfg.DefaultLeaseDuration
	}
	if in.ExecutionTimeout <= 0 {
		in.ExecutionTimeout = b.cfg.DefaultExecutionTimeout
	}

	if b.cfg.QueueCapacity > 0 && len(b.jobs) >= b.cfg.QueueCapacity {
		return models.Job{}, false, ErrQueueFull
	}

	fp := Fingerprint(in.Queue, in.Priority, in.Payload)

	if ex, ok := b.jobs[in.ID]; ok {
		if ex.Fingerprint != fp {
			return models.Job{}, false, ErrIDConflict
		}
		return *ex, false, nil
	}

	job := &models.Job{
		ID:               in.ID,
		Payload:          in.Payload,
		Priority:         in.Priority,
		State:            models.StateScheduled,
		Fingerprint:      fp,
		Attempt:          0,
		MaxAttempts:      in.MaxAttempts,
		CreatedAt:        now,
		ScheduledAt:      in.ScheduledAt,
		ExecutionTimeout: in.ExecutionTimeout,
		LeaseDuration:    in.LeaseDuration,
	}
	b.jobs[job.ID] = job
	//either scheduled or ready state
	if !job.ScheduledAt.After(now) {
		//ready here
		b.makeReadyLocked(job)
	} else {
		//dealyed schedule at
		b.putDelayedLocked(job,job.ScheduledAt)
	}
	//emit to observers

	return *job, true, nil
}
