package assist

import (
	"time"

	"bitbucket.org/lewington/autoroller/lg"
)

type Timer struct {
	logEvery        int
	totalRequests   int
	totalTime       time.Duration
	startOfRequests time.Time
}

func (t *Timer) logIfNeeded() {
	if t.totalRequests%t.logEvery == 0 {
		lg.L.Debug("Average time per request: %v", t.timePerRequest())
		t.totalRequests = 0
		t.totalTime = time.Millisecond * 0
		t.logTimeTaken()
	}
}

func (t *Timer) AddRequestTime(elapsedTime time.Duration) {
	mutex.Lock()
	t.totalRequests++
	t.totalTime += elapsedTime
	t.logIfNeeded()
	mutex.Unlock()
}

func (t *Timer) AddInstance() {
	mutex.Lock()
	t.totalRequests++
	if t.totalRequests%t.logEvery == 0 {
		t.logTimeTaken()
	}
	mutex.Unlock()
}

func (t *Timer) logTimeTaken() {
	currentTime := Timestamp()
	lg.L.Debug("Time taken for %v requests: %v", t.logEvery, currentTime.Sub(t.startOfRequests))
	t.startOfRequests = currentTime
}

func (t *Timer) timePerRequest() time.Duration {
	return t.totalTime / time.Duration(t.totalRequests+1)
}

func NewTimer(logEvery int) *Timer {
	return &Timer{logEvery, 0, time.Millisecond * 0, Timestamp()}
}
