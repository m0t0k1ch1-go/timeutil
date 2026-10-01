package timeutil_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/timeutil/v5"
)

func TestClock(t *testing.T) {
	clk := timeutil.NewClock()

	before := timeutil.NewTimestamp(time.Now())
	time.Sleep(time.Second)
	got := clk.Now()
	time.Sleep(time.Second)
	after := timeutil.NewTimestamp(time.Now())

	require.True(t, before.Before(got))
	require.True(t, got.Before(after))
	require.Equal(t, 0, got.Time().Nanosecond())
	require.Equal(t, time.UTC, got.Time().Location())
}

func TestMockClock(t *testing.T) {
	now := time.Unix(0, 0)
	clk := timeutil.NewMockClock(timeutil.NewTimestamp(now))

	got := clk.Now()
	require.True(t, got.Time().Equal(now))
	require.Equal(t, time.UTC, got.Time().Location())

	now = time.Unix(1231006505, 0)
	clk.Set(timeutil.NewTimestamp(now))

	got = clk.Now()
	require.True(t, got.Time().Equal(now))
	require.Equal(t, time.UTC, got.Time().Location())
}
