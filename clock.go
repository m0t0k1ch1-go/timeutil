package timeutil

// Clock is an interface for getting the current Timestamp.
type Clock interface {
	Now() Timestamp
}

type clock struct{}

// NewClock returns a new Clock.
func NewClock() Clock {
	return &clock{}
}

// Now returns the current Timestamp.
func (clk *clock) Now() Timestamp {
	return Now()
}

// MockClock is a mock implementation of Clock.
type MockClock struct {
	ts Timestamp
}

// NewMockClock returns a new MockClock.
func NewMockClock(ts Timestamp) *MockClock {
	return &MockClock{
		ts: ts,
	}
}

// Now returns the mock Timestamp.
func (clk *MockClock) Now() Timestamp {
	return clk.ts
}

// Set sets the mock Timestamp.
func (clk *MockClock) Set(ts Timestamp) {
	clk.ts = ts
}
