package timeutil

import (
	"time"
)

// Clock represents a time source that returns Timestamps.
type Clock interface {
	Now() Timestamp
}

type clock struct{}

// NewClock returns a new Clock.
func NewClock() Clock {
	return &clock{}
}

// Now returns the current time as a Timestamp.
func (clk clock) Now() Timestamp {
	return NewTimestamp(time.Now())
}

// MockClock implements Clock.
// It represents a controllable, fixed-time clock for testing.
type MockClock struct {
	ts Timestamp
}

// NewMockClock returns a new MockClock initialized with the given Timestamp.
func NewMockClock(ts Timestamp) *MockClock {
	return &MockClock{
		ts: ts,
	}
}

// Now returns the stored Timestamp.
func (clk *MockClock) Now() Timestamp {
	return clk.ts
}

// Set updates the stored Timestamp.
func (clk *MockClock) Set(ts Timestamp) {
	clk.ts = ts
}
