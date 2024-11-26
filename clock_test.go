package timeutil_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/timeutil/v4"
)

func TestClock(t *testing.T) {
	clk := timeutil.NewClock()

	ts := clk.Now()

	time.Sleep(1 * time.Second)

	require.True(t, clk.Now().Time().After(ts.Time()))
}

func TestMockClock(t *testing.T) {
	ts := timeutil.Now()

	clk := timeutil.NewMockClock(ts)

	require.Equal(t, ts, clk.Now())

	time.Sleep(1 * time.Second)

	require.Equal(t, ts, clk.Now())
}
