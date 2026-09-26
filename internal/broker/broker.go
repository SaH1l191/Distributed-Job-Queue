package broker

import (
	"jobqueue/internal/clock"
	"jobqueue/internal/models"
	"jobqueue/internal/queue"
	"sync"
	"time"
)

type ReadyKey struct {
	Priority int
	Seq      uint64
}

func BetterReady(a ReadyKey, b ReadyKey) bool {
	if a.Priority == b.Priority {
		return a.Seq < b.Seq
	}
	return a.Priority > b.Priority
}

type Config struct {
	QueueCapacity           int
	DefaultMaxAttempts      int
	DefaultLeaseDuration    time.Duration
	DefaultExecutionTimeout time.Duration
}

type SubmitInput struct {
	ID               string
	Payload          []byte
	Priority         int
	Queue            string
	Fingerprint      string
	ScheduledAt      time.Time
	MaxAttempts      int
	ExecutionTimeout time.Duration
	LeaseDuration    time.Duration
}

type Broker struct {
	mu       sync.Mutex
	clock    clock.Clock
	cfg      Config
	jobs     map[string]*models.Job
	delayed  *queue.IndexedHeap[time.Time]
	ready    map[string]*queue.IndexedHeap[ReadyKey]
	inflight *queue.IndexedHeap[time.Time]

	dlq []string
	// completed *completedSet

	epochSeq uint64 //monotonically increasing fencing toekn

	//fifo seq for ready jobs
	readySeq uint64

	// workers map[string]*Worker

	//notificaitons
	readyCh chan struct{}
	wakeCh  chan struct{}
}

func NewBroker(clock clock.Clock, cfg Config) *Broker {

	return &Broker{
		clock: clock,
		cfg:   cfg,
		jobs:  make(map[string]*models.Job),
		delayed: queue.NewIndexedHeap(
			func(a time.Time, b time.Time) bool {
				return a.Before(b)
			},
		),
		ready: make(map[string]*queue.IndexedHeap[ReadyKey]),
		inflight: queue.NewIndexedHeap(
			func(a time.Time, b time.Time) bool {
				return a.Before(b)
			},
		),
		dlq:     make([]string, 0),
		readyCh: make(chan struct{}),
		wakeCh:  make(chan struct{}),
	}

}
