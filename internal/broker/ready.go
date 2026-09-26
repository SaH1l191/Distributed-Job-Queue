package broker

import (
	"jobqueue/internal/models"
	"jobqueue/internal/queue"
)

func (b *Broker) notifyReadyLocked() {
	close(b.readyCh) //closing a channel wakes all receivers.
	b.readyCh = make(chan struct{})
	//the next notification has a new channel to close
}//readyCh = "something changed" signal

func (b *Broker) makeReadyLocked(queueName string, job *models.Job) {
	h := b.ready[queueName]

	if h == nil {
		h = queue.NewIndexedHeap(func(a ReadyKey, b ReadyKey) bool {
			return BetterReady(a, b)
		})
		b.ready[queueName] = h
	}
	b.readySeq++
	key := ReadyKey{
		Priority: job.Priority,
		Seq:      b.readySeq,
	}
	job.State = models.StateReady
	h.Put(job.ID, key)

	b.notifyReadyLocked()
}
