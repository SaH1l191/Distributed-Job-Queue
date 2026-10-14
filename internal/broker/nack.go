package broker

import "jobqueue/internal/models"

func (b *Broker) Nac(id string, epoch uint64, reason string, retryable bool) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	job, ok := b.jobs[id]
	if !ok {
		return ErrNotFound
	}
	if job.State != models.StateProcessing {
		return ErrNotProcessing
	}
	if job.LeaseEpoch != epoch {
		return ErrStale
	}
	b.inflight.Delete(job.ID)
	b.failJobLocked(job, reason, retryable)
	return nil
}
