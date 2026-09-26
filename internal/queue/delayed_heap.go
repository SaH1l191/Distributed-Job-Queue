package queue

import "time"

// NewDelayedHeap creates a min-heap ordered by ScheduledAt.
// The earliest scheduled job is always at the root.
// Used for the delayed queue: jobs wait here until their ScheduledAt is reached.
func NewDelayedHeap() *IndexedHeap[time.Time] {
	return NewIndexedHeap(func(a time.Time, b time.Time) bool {
		return a.Before(b) // earliest time first (min-heap)
	})
}
