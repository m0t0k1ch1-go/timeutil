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

func TestNewTimestamp(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   time.Time
			want int64
		}{
			{
				"Unix epoch in JST",
				time.Date(1970, 1, 1, 9, 0, 0, 0, time.FixedZone("JST", 9*60*60)),
				0,
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				ts := timeutil.NewTimestamp(tc.in)
				require.Equal(t, tc.want, ts.Unix())
				require.Equal(t, time.UTC, ts.Time().Location())
			})
		}
	})
}

func TestNewTimestampFromUnix(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   int64
			want int64
		}{
			{
				"zero",
				0,
				0,
			},
			{
				"positive",
				1231006505,
				1231006505,
			},
			{
				"negative",
				-1231006505,
				-1231006505,
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				ts := timeutil.NewTimestampFromUnix(tc.in)
				require.Equal(t, tc.want, ts.Unix())
				require.Equal(t, time.UTC, ts.Time().Location())
			})
		}
	})
}

func TestTimestamp_String(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   timeutil.Timestamp
			want string
		}{
			{
				"zero",
				timeutil.NewTimestampFromUnix(0),
				"0",
			},
			{
				"positive",
				timeutil.NewTimestampFromUnix(1231006505),
				"1231006505",
			},
			{
				"negative",
				timeutil.NewTimestampFromUnix(-1231006505),
				"-1231006505",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				s := tc.in.String()
				require.Equal(t, tc.want, s)
			})
		}
	})
}

func TestTimestamp_Value(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   timeutil.Timestamp
			want driver.Value
		}{
			{
				"zero",
				timeutil.NewTimestampFromUnix(0),
				int64(0),
			},
			{
				"positive",
				timeutil.NewTimestampFromUnix(1231006505),
				int64(1231006505),
			},
			{
				"negative",
				timeutil.NewTimestampFromUnix(-1231006505),
				int64(-1231006505),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				v, err := tc.in.Value()
				require.NoError(t, err)

				require.Equal(t, tc.want, v)
			})
		}
	})
}

func TestTimestamp_Scan(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			want string
		}{
			{
				"nil",
				nil,
				"invalid source: nil",
			},
			{
				"time.Time",
				time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
				"unsupported source type: time.Time",
			},
			{
				"uint64: exceeds int64 range",
				uint64(math.MaxInt64) + 1,
				"invalid source: exceeds int64 range",
			},
			{
				"[]byte: empty",
				[]byte{},
				"invalid source: empty []byte",
			},
			{
				"[]byte: invalid",
				[]byte("invalid"),
				"invalid source",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var ts timeutil.Timestamp
				{
					err := ts.Scan(tc.in)
					require.ErrorContains(t, err, tc.want)
				}
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			want int64
		}{
			{
				"int64: positive",
				int64(1231006505),
				1231006505,
			},
			{
				"int64: negative",
				int64(-1231006505),
				-1231006505,
			},
			{
				"uint64",
				uint64(1231006505),
				1231006505,
			},
			{
				"[]byte",
				[]byte("1231006505"),
				1231006505,
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var ts timeutil.Timestamp
				err := ts.Scan(tc.in)
				require.NoError(t, err)

				require.Equal(t, tc.want, ts.Unix())
				require.Equal(t, time.UTC, ts.Time().Location())
			})
		}
	})
}

func TestTimestamp_JSONMarshal(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   timeutil.Timestamp
			want []byte
		}{
			{
				"zero",
				timeutil.NewTimestampFromUnix(0),
				[]byte(`0`),
			},
			{
				"positive",
				timeutil.NewTimestampFromUnix(1231006505),
				[]byte(`1231006505`),
			},
			{
				"negative",
				timeutil.NewTimestampFromUnix(-1231006505),
				[]byte(`-1231006505`),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				b, err := json.Marshal(tc.in)
				require.NoError(t, err)

				require.Equal(t, tc.want, b)
			})
		}
	})
}

func TestTimestamp_JSONUnmarshal(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want string
		}{
			{
				"null",
				[]byte(`null`),
				"invalid json value: null",
			},
			{
				"number: exceeds int64 range",
				[]byte(`9223372036854775808`),
				"invalid json number",
			},
			{
				"number: fractional",
				[]byte(`1231006505.0`),
				"invalid json number",
			},
			{
				"number: exponential",
				[]byte(`1231006505e0`),
				"invalid json number",
			},
			{
				"string: empty",
				[]byte(`""`),
				"invalid json number",
			},
			{
				"string: positive decimal",
				[]byte(`"1231006505"`),
				"invalid json number",
			},
			{
				"string: negative decimal",
				[]byte(`"-1231006505"`),
				"invalid json number",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var ts timeutil.Timestamp
				err := json.Unmarshal(tc.in, &ts)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want int64
		}{
			{
				"number: zero",
				[]byte(`0`),
				0,
			},
			{
				"number: positive",
				[]byte(`1231006505`),
				1231006505,
			},
			{
				"number: negative",
				[]byte(`-1231006505`),
				-1231006505,
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var ts timeutil.Timestamp
				err := json.Unmarshal(tc.in, &ts)
				require.NoError(t, err)

				require.Equal(t, tc.want, ts.Unix())
				require.Equal(t, time.UTC, ts.Time().Location())
			})
		}
	})
}
