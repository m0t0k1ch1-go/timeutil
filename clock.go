package timeutil

type Clock interface {
	Now() Timestamp
}

type clock struct{}

func NewClock() Clock {
	return &clock{}
}

func (clk *clock) Now() Timestamp {
	return Now()
}

type MockClock struct {
	ts Timestamp
}

func NewMockClock(ts Timestamp) *MockClock {
	return &MockClock{
		ts: ts,
	}
}

func (clk *MockClock) Now() Timestamp {
	return clk.ts
}

func (clk *MockClock) Set(ts Timestamp) {
	clk.ts = ts
}
