package timeutil

import (
	"time"
)

var (
	_ Clock = &MockClock{}
)

// Clock represents a time source.
type Clock interface {
	// Now returns a [Timestamp] representing the current time.
	Now() Timestamp
}

type clock struct{}

// NewClock returns a new [Clock].
func NewClock() Clock {
	return &clock{}
}

func (clk clock) Now() Timestamp {
	return NewTimestamp(time.Now())
}

// MockClock represents a controllable, fixed-time source for testing.
type MockClock struct {
	ts Timestamp
}

// NewMockClock returns a new [MockClock] initialized with the given [Timestamp].
func NewMockClock(ts Timestamp) *MockClock {
	return &MockClock{
		ts: ts,
	}
}

// Now implements [Clock.Now].
// It returns the stored [Timestamp].
func (clk *MockClock) Now() Timestamp {
	return clk.ts
}

// Set updates the stored [Timestamp] to ts.
func (clk *MockClock) Set(ts Timestamp) {
	clk.ts = ts
}
