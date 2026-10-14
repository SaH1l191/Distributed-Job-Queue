package broker

import (
	"jobqueue/internal/models"
	"time"
)

func (b *Broker) notifyReadyLocked() {
	close(b.readyCh) //closing a channel wakes all receivers.
	b.readyCh = make(chan struct{})
	//the next notification has a new channel to close
}//readyCh = "something changed" signal
// When a worker calls claim() and the ready queue is empty, it needs to wait until a job arrives. Instead of busy-spinning (polling every few milliseconds), it blocks on this channel:
//close and replace => broadcast, losers jsut find empty heap and wait again ,
//no polling -> no wasted CPU

func (b *Broker) signalWakeLocked() {
	close(b.wakeCh)
	b.wakeCh = make(chan struct{})
}

func (b *Broker) putDelayedLocked(job *models.Job, at time.Time) {
	_, headBefore, hadHead := b.delayed.Peek()
	b.delayed.Put(job.ID, at)
	job.ScheduledAt = at

	_, headAfter, _ := b.delayed.Peek()
	if !hadHead || headAfter.Before(headBefore) {
		b.signalWakeLocked()
	}
}

func (b *Broker) makeReadyLocked(job *models.Job) {
	b.readySeq++
	key := ReadyKey{
		Priority: job.Priority,
		Seq:      b.readySeq,
	}
	job.State = models.StateReady
	b.ready.Put(job.ID, key)

	b.notifyReadyLocked()
}
