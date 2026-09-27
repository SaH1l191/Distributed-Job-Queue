package broker

import (
	"jobqueue/internal/models"
	"time"
)


//JOB => READY → PROCESSING
//inflightHeap uses a min-heap on LeaseExpiresAt.



func minTime(a time.Time, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func (b *Broker) tryClaim(workerID string) (j models.Job, ok bool) {
	now := b.clock.Now()
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.ready.Len() == 0 {
		return models.Job{}, false
	}

	id, _, _ := b.ready.PopMin()
	job := b.jobs[id]

	b.epochSeq++ //global fencing token wrt the job and the worker processing that job
	job.State = models.StateProcessing
	job.Attempt++
	job.WorkerID = workerID
	job.LeaseEpoch = b.epochSeq
	job.AttemptStarted = now
	job.LeaseExpiresAt = minTime(now.Add(job.LeaseDuration), now.Add(job.ExecutionTimeout))
	b.inflight.Put(job.ID, job.LeaseExpiresAt)

	return *job, true
}
