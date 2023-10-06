package timeutil_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/m0t0k1ch1-go/timeutil"
	"github.com/m0t0k1ch1-go/timeutil/internal/testutil"
)

func TestScan(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			Name   string
			Input  any
			Output int64
		}{
			{
				Name:   "int64(1231006505)",
				Input:  int64(1231006505),
				Output: 1231006505,
			},
			{
				Name:   "[]byte(1231006505)",
				Input:  []byte("1231006505"),
				Output: 1231006505,
			},
		}

		for _, tc := range tcs {
			t.Run(tc.Name, func(t *testing.T) {
				ts := new(timeutil.Timestamp)
				if err := ts.Scan(tc.Input); err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.Output, ts.Time().Unix())
			})
		}
	})
}

func TestMarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			Name   string
			Input  timeutil.Timestamp
			Output string
		}{
			{
				Name:   "1231006505",
				Input:  timeutil.Time(time.Unix(1231006505, 0)),
				Output: "1231006505",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.Name, func(t *testing.T) {
				b, err := json.Marshal(tc.Input)
				if err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.Output, string(b))
			})
		}
	})
}

func TestUnmarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			Name   string
			Input  []byte
			Output int64
		}{
			{
				Name:   "1231006505",
				Input:  []byte("1231006505"),
				Output: 1231006505,
			},
		}

		for _, tc := range tcs {
			t.Run(tc.Name, func(t *testing.T) {
				var ts timeutil.Timestamp
				if err := json.Unmarshal(tc.Input, &ts); err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.Output, ts.Time().Unix())
			})
		}
	})
}
