package timeutil_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/m0t0k1ch1-go/timeutil/v2"
	"github.com/m0t0k1ch1-go/timeutil/v2/internal/testutil"
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
				timeutil.TimeToTimestamp(time.Unix(1231006505, 0)),
				int64(1231006505),
			},
			{
				"negative",
				timeutil.TimeToTimestamp(time.Unix(-1231006505, 0)),
				int64(-1231006505),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				v, err := tc.in.Value()
				if err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.out, v)
			})
		}
	})
}

func TestTimestampScan(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			out  timeutil.Timestamp
		}{
			{
				name: "positive int64",
				in:   int64(1231006505),
				out:  timeutil.TimeToTimestamp(time.Unix(1231006505, 0)),
			},
			{
				name: "negative int64",
				in:   int64(-1231006505),
				out:  timeutil.TimeToTimestamp(time.Unix(-1231006505, 0)),
			},
			{
				name: "[]byte",
				in:   []byte("1231006505"),
				out:  timeutil.TimeToTimestamp(time.Unix(1231006505, 0)),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				ts := timeutil.Timestamp{}
				if err := (&ts).Scan(tc.in); err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.out, ts, cmp.AllowUnexported(timeutil.Timestamp{}))
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
				name: "positive",
				in:   timeutil.TimeToTimestamp(time.Unix(1231006505, 0)),
				out:  []byte("1231006505"),
			},
			{
				name: "negative",
				in:   timeutil.TimeToTimestamp(time.Unix(-1231006505, 0)),
				out:  []byte("-1231006505"),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				b, err := json.Marshal(tc.in)
				if err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.out, b)
			})
		}
	})
}

func TestTimestampUnmarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			out  timeutil.Timestamp
		}{
			{
				name: "positive",
				in:   []byte("1231006505"),
				out:  timeutil.TimeToTimestamp(time.Unix(1231006505, 0)),
			},
			{
				name: "negative",
				in:   []byte("-1231006505"),
				out:  timeutil.TimeToTimestamp(time.Unix(-1231006505, 0)),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var ts timeutil.Timestamp
				if err := json.Unmarshal(tc.in, &ts); err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.out, ts, cmp.AllowUnexported(timeutil.Timestamp{}))
			})
		}
	})
}
