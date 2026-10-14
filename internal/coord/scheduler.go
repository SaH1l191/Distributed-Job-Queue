package coord

import (
	"context"
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

func (s *Scheduler) Run(ctx context.Context) {
	s.loop(ctx)
}

func (s *Scheduler) loop(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			return
		}

		at, ok := s.b.NextDelayedAt()
		if !ok {
			if !s.waitForWake(ctx) { //queues empty , instead of waiting , return and let the broker handle the wakeup
				return
			}
			continue
		}

		d := at.Sub(s.clock.Now())
		if d <= 0 {
			s.b.PromoteDue()
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

func (s *Scheduler) waitForWake(ctx context.Context) bool {
	select {
	case <-s.b.WakeCh():
		return true
	case <-ctx.Done():
		return false
	}
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

		job, _ := s.b.GetJob(id)
		if job != nil {
			s.b.MakeReadyLocked(job)
		}
	}
}
