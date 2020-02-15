package auth

import (
	"time"

	"bitbucket.org/lewington/autoroller/assist"
)

// TimeThrottle allows a certain number of incorrect logins
// and then initiates a timeout during which any attempt
// will be denied.
type TimeThrottle struct {
	throttleAt        int
	failCount         int
	throttleFor       time.Duration
	releaseThrottleAt time.Time
}

// Allow returns true if there is currently no
// timeout in progress.
func (t *TimeThrottle) Allow() bool {
	return assist.Timestamp().After(t.releaseThrottleAt)
}

// Fail records a login failure.
func (t *TimeThrottle) Fail() {
	t.failCount++

	if t.failCount >= t.throttleAt {
		t.startTimeout()
		t.failCount = 0
	}
}

// Succeed records a successful login.
func (t *TimeThrottle) Succeed() {
	t.failCount = 0
}

func (t *TimeThrottle) startTimeout() {
	t.releaseThrottleAt = assist.Timestamp().Add(t.throttleFor)
}

// NewTimeThrottle initializes a TimeThrottle.
func NewTimeThrottle(mistakesAllowed int, timeout time.Duration) Throttle {
	return &TimeThrottle{
		mistakesAllowed,
		0,
		timeout,
		assist.Timestamp().Add(time.Minute * -1),
	}
}
