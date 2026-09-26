package queue

import "time"

// NewInflightHeap creates a min-heap ordered by lease expiry time.
// The job whose lease expires soonest is always at the root.
// Used by the reaper to efficiently find which job's lease expires next.
func NewInflightHeap() *IndexedHeap[time.Time] {
	return NewIndexedHeap(func(a time.Time, b time.Time) bool {
		return a.Before(b) // earliest expiry first (min-heap)
	})
}
