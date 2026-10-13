package clock

import "time"

type Timer interface {
	C() <-chan time.Time
	Stop() bool
	Reset(d time.Duration) bool
}

type Clock interface {
	Now() time.Time
	NewTimer(d time.Duration) Timer
}

//eg :
// func waitForTimeout(c Clock, d time.Duration) {
//     timer := c.NewTimer(d)
//     defer timer.Stop()

//     <-timer.C()
// }

type realTimer struct {
	t *time.Timer
}

func (r *realTimer) C() <-chan time.Time {
	return r.t.C
}
func (r *realTimer) Stop() bool {
	return r.t.Stop()
}
func (r *realTimer) Reset(d time.Duration) bool {
	return r.t.Reset(d)
}

type RealClock struct{}

func NewRealClock() Clock {
	return &RealClock{}
}

func (r *RealClock) Now() time.Time {
	return time.Now()
}

func (r *RealClock) NewTimer(d time.Duration) Timer {
	return &realTimer{t: time.NewTimer(d)}
}
