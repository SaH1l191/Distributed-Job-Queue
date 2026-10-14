package broker

import "time"

func (b *Broker) Backoff(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	exp := b.cfg.BackOffBase
	for i := 1; i < attempt; i++ {
		exp *= 2
		if exp >= b.cfg.BakoffCap {
			exp = b.cfg.BakoffCap
			break
		}
	}
	if exp <= 0 {
		return 0
	}
	n := b.cfg.Rand(int64(exp))

	if n < 0 {
		n = 0
	}
	if n >= int64(exp) {
		n = int64(exp) - 1
	}
	return time.Duration(n)
}
