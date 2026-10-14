package broker

import "jobqueue/internal/models"

func (b *Broker) Ack(id string, epoch uint64, result []byte) error {
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
	now := b.clock.Now()
	job.State = models.StateCompleted
	job.CompletedAt = &now
	job.Result = result
	job.LastError = ""
	job.WorkerID = ""
	b.inflight.Delete(job.ID)
	return nil
}