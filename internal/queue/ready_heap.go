package queue

// ReadyKey is the sort key for the ready queue: higher priority first,
// then FIFO by sequence number within the same priority.
type ReadyKey struct {
	Priority int
	Seq      uint64
}

// NewReadyHeap creates a ready queue indexed heap.
// Higher priority jobs come first; within the same priority, lower seq (older) comes first.
func NewReadyHeap() *IndexedHeap[ReadyKey] {
	return NewIndexedHeap(func(a ReadyKey, b ReadyKey) bool {
		if a.Priority == b.Priority {
			return a.Seq < b.Seq // FIFO within priority
		}
		return a.Priority > b.Priority // higher priority first
	})
}
