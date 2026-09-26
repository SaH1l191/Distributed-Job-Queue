package broker

import "jobqueue/internal/models"

func (b *Broker) Submit(in SubmitInput) (models.Job, bool, error) {
	now := b.clock.Now()
	b.mu.Lock()
	defer b.mu.Unlock()

	if ex, ok := b.jobs[in.ID]; ok {
		if ex.Fingerprint != in.Fingerprint {
			return models.Job{}, false, ErrIDConflict
		}
		return *ex, false, nil
	}

	job := &models.Job{
		ID:               in.ID,
		Payload:          in.Payload,
		Priority:         in.Priority,
		State:            models.StateScheduled,
		Fingerprint:      in.Fingerprint,
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
		b.makeReadyLocked(in.Queue, job)
	} else {
		//dealyed schedule at
		b.delayed.Put(job.ID, job.ScheduledAt)
	}
	//emit to observers

	return *job, true, nil
}
