package timeutil_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/m0t0k1ch1-go/timeutil"
	"github.com/m0t0k1ch1-go/timeutil/internal/testutil"
)

func TestTimestampScan(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			out  int64
		}{
			{
				name: "int64(1231006505)",
				in:   int64(1231006505),
				out:  1231006505,
			},
			{
				name: "[]byte(1231006505)",
				in:   []byte("1231006505"),
				out:  1231006505,
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				ts := new(timeutil.Timestamp)
				if err := ts.Scan(tc.in); err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.out, ts.Time().Unix())
			})
		}
	})
}

func TestTimestampMarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   timeutil.Timestamp
			out  string
		}{
			{
				name: "1231006505",
				in:   timeutil.Time(time.Unix(1231006505, 0)),
				out:  "1231006505",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				b, err := json.Marshal(tc.in)
				if err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.out, string(b))
			})
		}
	})
}

func TestTimestampUnmarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			out  int64
		}{
			{
				name: "1231006505",
				in:   []byte("1231006505"),
				out:  1231006505,
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var ts timeutil.Timestamp
				if err := json.Unmarshal(tc.in, &ts); err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.out, ts.Time().Unix())
			})
		}
	})
}
