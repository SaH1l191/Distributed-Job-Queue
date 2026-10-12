package clock

import "time"


type Timer interface {
    C() <-chan time.Time
    Stop() bool
    Reset(d time.Duration) bool
}

type Clock interface{ 
	Now() time.Time
	NewTimer(d time.Duration) Timer
}

//eg : 
// func waitForTimeout(c Clock, d time.Duration) {
//     timer := c.NewTimer(d)
//     defer timer.Stop()

//     <-timer.C()
// }