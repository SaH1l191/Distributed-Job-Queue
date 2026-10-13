package coord

import (
	"jobqueue/internal/broker"
	"jobqueue/internal/clock"
)

type Scheduler struct {
	clock clock.Clock
	b     *broker.Broker
}

func NewScheduler(b *broker.Broker, clock clock.Clock) *Scheduler {
	return &Scheduler{b: b, clock: clock}
}

// DI pattern: Scheduler depends on Broker,
//  but Broker does not depend on Scheduler.

// Scheduler is responsible
// for managing the timing and scheduling of jobs,
//  Broker is responsible for managing the
//  state and lifecycle of jobs.

// The Scheduler needs
// to interact with the Broker to move jobs from the
// delayed queue to the ready queue, but the Broker
// does not need to know about the Scheduler's existence
//  or functionality. This separation of concerns allows
//  for better modularity and testability of the code.

func (s *Scheduler) Run() {
	s.loop()
}

func (s *Scheduler) loop() {
	for {
		_, scheduledAt, ok := s.b.DelayedPeek()
		if !ok {
			s.waitForWake()
			continue
		}

		now := s.clock.Now()
		d := scheduledAt.Sub(now)

		if d <= 0 {
			s.promoteDelayed()
			continue
		}

		timer := s.clock.NewTimer(d)
		select {
		case <-timer.C():
			s.promoteDelayed()
		case <-s.b.WakeCh():
			timer.Stop()
		}

	}
}

func (s *Scheduler) waitForWake() {
	<-s.b.WakeCh()
}

func (s *Scheduler) promoteDelayed() {

	for {
		id, scheduledAt, ok := s.b.DelayedPeek()
		if !ok {
			return
		}

		now := s.clock.Now()
		if scheduledAt.After(now) {
			return
		}
		s.b.DelayedPop()

		job := s.b.GetJob(id)
		if job != nil {
			s.b.MakeReadyLocked(job)
		}
	}
}
