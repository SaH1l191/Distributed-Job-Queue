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
	BackOffBase             time.Duration
	BakoffCap               time.Duration
}

type SubmitInput struct {
	ID               string
	Payload          []byte
	Priority         int
	Queue            string
	ScheduledAt      time.Time
	MaxAttempts      int
	ExecutionTimeout time.Duration
	LeaseDuration    time.Duration
}

type Broker struct {
	mu       sync.Mutex
	clock    clock.Clock //fake injectable clock for testing/prod
	cfg      Config
	jobs     map[string]*models.Job
	delayed  *queue.IndexedHeap[time.Time]
	ready    *queue.IndexedHeap[ReadyKey]
	inflight *queue.IndexedHeap[time.Time]

	dlq []string
	// completed *completedSet

	epochSeq uint64 //monotonically increasing fencing toekn(incremented per claim)

	//fifo seq for ready jobs
	readySeq uint64

	// workers map[string]*Worker

	//notificaitons
	readyCh chan struct{} // Closed when jobs become READY (long-poll)
	wakeCh  chan struct{} // Signaled when earlier delayed job inserted
}

////dedup detction sha256 -> calc fingerprint of queue+priority+payload

func NewBroker(clock clock.Clock, cfg Config) *Broker {

	if cfg.BackOffBase <= 0 {
		cfg.BackOffBase = 500 * time.Millisecond
	}
	if cfg.BakoffCap <= 0 {
		cfg.BakoffCap = 30 * time.Second
	}

	return &Broker{
		clock: clock,
		cfg:   cfg,
		jobs:  make(map[string]*models.Job),
		delayed: queue.NewIndexedHeap(
			func(a time.Time, b time.Time) bool {
				return a.Before(b)
			},
		),
		ready: queue.NewIndexedHeap(func(a ReadyKey, b ReadyKey) bool {
			return BetterReady(a, b)
		}),
		inflight: queue.NewIndexedHeap(
			func(a time.Time, b time.Time) bool {
				return a.Before(b)
			},
		),
		dlq:     make([]string, 0),
		readyCh: make(chan struct{}), //Recalculate when the next delayed job is due.
		wakeCh:  make(chan struct{}), //Recalculate when the next delayed job is due.
	}

}

func (b *Broker) DelayedPeek() (string, time.Time, bool) {
	return b.delayed.Peek()
}

func (b *Broker) NextDelayedAt() (time.Time, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	_, at, ok := b.delayed.Peek()
	return at, ok
}


func (b *Broker) DelayedPop() string {
	id, _, _ := b.delayed.PopMin()
	return id
}

func (b *Broker) ClockNow() time.Time {
	return b.clock.Now()
}

func (b *Broker) GetJob(id string) (*models.Job, bool) {
	j, ok := b.jobs[id]
	if !ok {
		return nil, false
	}
	return j, true
}

func (b *Broker) MakeReadyLocked(job *models.Job) {
	b.makeReadyLocked(job)
}

func (b *Broker) WakeCh() <-chan struct{} {
	return b.wakeCh
}
