package timeutil_test

import (
	"database/sql/driver"
	"encoding/json"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/timeutil/v4"
)

func TestTimestampValue(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   timeutil.Timestamp
			out  driver.Value
		}{
			{
				"positive",
				timeutil.NewTimestamp(time.Unix(1231006505, 0)),
				int64(1231006505),
			},
			{
				"negative",
				timeutil.NewTimestamp(time.Unix(-1231006505, 0)),
				int64(-1231006505),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				v, err := tc.in.Value()
				require.NoError(t, err)

				require.Equal(t, tc.out, v)
			})
		}
	})
}

func TestTimestampScan(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
		}{
			{
				"nil",
				nil,
			},
			{
				"string",
				"1231006505",
			},
			{
				"too large uint64",
				uint64(math.MaxInt64) + 1,
			},
			{
				"invalid []byte",
				[]byte("invalid"),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var ts timeutil.Timestamp
				{
					err := ts.Scan(tc.in)
					require.Error(t, err)
				}
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			out  timeutil.Timestamp
		}{
			{
				"positive int64",
				int64(1231006505),
				timeutil.NewTimestamp(time.Unix(1231006505, 0)),
			},
			{
				"negative int64",
				int64(-1231006505),
				timeutil.NewTimestamp(time.Unix(-1231006505, 0)),
			},
			{
				"uint64",
				uint64(1231006505),
				timeutil.NewTimestamp(time.Unix(1231006505, 0)),
			},
			{
				"[]byte",
				[]byte("1231006505"),
				timeutil.NewTimestamp(time.Unix(1231006505, 0)),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var ts timeutil.Timestamp
				{
					err := ts.Scan(tc.in)
					require.NoError(t, err)
				}

				require.Equal(t, tc.out, ts)
			})
		}
	})
}

func TestTimestampMarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   timeutil.Timestamp
			out  []byte
		}{
			{
				"positive",
				timeutil.NewTimestamp(time.Unix(1231006505, 0)),
				[]byte("1231006505"),
			},
			{
				"negative",
				timeutil.NewTimestamp(time.Unix(-1231006505, 0)),
				[]byte("-1231006505"),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				b, err := json.Marshal(tc.in)
				require.NoError(t, err)

				require.Equal(t, tc.out, b)
			})
		}
	})
}

func TestTimestampUnmarshalJSON(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
		}{
			{
				"invalid",
				[]byte("invalid"),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var ts timeutil.Timestamp
				{
					err := json.Unmarshal(tc.in, &ts)
					require.Error(t, err)
				}
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			out  timeutil.Timestamp
		}{
			{
				"positive",
				[]byte("1231006505"),
				timeutil.NewTimestamp(time.Unix(1231006505, 0)),
			},
			{
				"negative",
				[]byte("-1231006505"),
				timeutil.NewTimestamp(time.Unix(-1231006505, 0)),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var ts timeutil.Timestamp
				{
					err := json.Unmarshal(tc.in, &ts)
					require.NoError(t, err)
				}

				require.Equal(t, tc.out, ts)
			})
		}
	})
}
