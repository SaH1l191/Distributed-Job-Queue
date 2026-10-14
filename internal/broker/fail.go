package broker

import "jobqueue/internal/models"

func (b *Broker) failJobLocked(job *models.Job, reason string, retryable bool) {
	job.LastError = reason
	job.WorkerID = ""

	if retryable && job.Attempt < job.MaxAttempts {
		job.State = models.StateScheduled
		backoff := b.Backoff(job.Attempt)
		b.putDelayedLocked(job, b.ClockNow().Add(backoff))
	}
	job.State = models.StateDead
	b.dlq = append(b.dlq, job.ID)

}
