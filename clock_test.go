package timeutil_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/timeutil/v4"
)

func TestClock(t *testing.T) {
	clk := timeutil.NewClock()

	before := time.Now()
	got := clk.Now()
	after := time.Now()

	require.False(t, got.Time().Before(before))
	require.False(t, got.Time().After(after))
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
