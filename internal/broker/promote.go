package broker

func (b *Broker) PromoteDue() int{
	b.mu.Lock()
	defer b.mu.Unlock()

	now := b.clock.Now()
	promoted := 0
	for {
		id, at, ok := b.delayed.Peek()
		if !ok {
			break
		}
		if at.After(now) {
			break // min-heap: nothing else is due
		}
		b.delayed.PopMin()

		if job, exists := b.jobs[id]; exists {
			b.makeReadyLocked(job)
			promoted++
		}
	}
	return promoted
}
