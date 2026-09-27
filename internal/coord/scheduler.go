package coord

import (
	"jobqueue/internal/broker"
	"jobqueue/internal/clock"
)

type Scheduler struct {
	clock clock.Clock
	b *broker.Broker
}

func NewScheduler(b *broker.Broker, clock clock.Clock) *Scheduler {
	return &Scheduler{b: b, clock: clock}
}










