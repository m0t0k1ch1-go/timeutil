package timeutil_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/m0t0k1ch1-go/timeutil/v4"
	"github.com/m0t0k1ch1-go/timeutil/v4/internal/testutil"
)

func TestClock(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		clk := timeutil.NewClock()

		ts := clk.Now()

		time.Sleep(1 * time.Second)

		testutil.Equal(t, true, clk.Now().Time().After(ts.Time()))
	})
}

func TestMockClock(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ts := timeutil.Now()

		clk := timeutil.NewMockClock(ts)

		testutil.Equal(t, ts, clk.Now(), cmp.AllowUnexported(timeutil.Timestamp{}))

		time.Sleep(1 * time.Second)

		testutil.Equal(t, ts, clk.Now(), cmp.AllowUnexported(timeutil.Timestamp{}))
	})
}
