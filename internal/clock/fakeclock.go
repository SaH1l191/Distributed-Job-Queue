package clock

import "time"

type FakeClock struct {
	now time.Time
}


func NewFakeClock() *FakeClock {
	return &FakeClock{now: time.Now()}
}

func (f *FakeClock) Now() time.Time {
	return f.now
}

func (f *FakeClock) Advance(d time.Duration) {
	f.now = f.now.Add(d)
}

func (f *FakeClock) NewTimer(d time.Duration) Timer {
	return &fakeTimer{clock: f, duration: d}
}

type fakeTimer struct {
	clock    *FakeClock
	duration time.Duration
}

func (f *fakeTimer) C() <-chan time.Time {
	return time.After(f.duration)
}

func (f *fakeTimer) Stop() bool {
	return true
}

func (f *fakeTimer) Reset(d time.Duration) bool {
	f.duration = d
	return true
}
