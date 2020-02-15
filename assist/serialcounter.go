package assist

import "sync"

var mutex = &sync.Mutex{}

type SerialCounter struct {
	count int
}

// Next increments the count by one and returns the count
// behind a monitor.
func (s *SerialCounter) Next() int {
	mutex.Lock()
	s.count++
	c := s.count
	mutex.Unlock()
	return c
}

// NewSerialCounter initializes a SerialCounter with a starting
// count of 0.
func NewSerialCounter() *SerialCounter {
	return &SerialCounter{
		0,
	}
}
