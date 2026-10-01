package timeutil_test

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json/v2"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/timeutil/v5"
)

func TestTimestamp(t *testing.T) {
	var ts timeutil.Timestamp
	require.Implements(t, (*fmt.Stringer)(nil), &ts)
	require.Implements(t, (*driver.Valuer)(nil), &ts)
	require.Implements(t, (*sql.Scanner)(nil), &ts)
	require.Implements(t, (*encoding.TextMarshaler)(nil), &ts)
	require.Implements(t, (*json.MarshalerTo)(nil), &ts)
	require.Implements(t, (*json.Marshaler)(nil), &ts)
	require.Implements(t, (*graphql.Marshaler)(nil), &ts)
	require.Implements(t, (*encoding.TextUnmarshaler)(nil), &ts)
	require.Implements(t, (*json.UnmarshalerFrom)(nil), &ts)
	require.Implements(t, (*json.Unmarshaler)(nil), &ts)
	require.Implements(t, (*graphql.Unmarshaler)(nil), &ts)
}

func TestNewTimestamp(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   time.Time
			want time.Time
		}{
			{
				"Unix epoch in JST",
				time.Date(1970, 1, 1, 9, 0, 0, 0, time.FixedZone("JST", 9*60*60)),
				time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				"positive with sub-second part",
				time.Date(2009, 1, 3, 18, 15, 5, 999_999_999, time.UTC),
				time.Date(2009, 1, 3, 18, 15, 5, 0, time.UTC),
			},
			{
				"negative with sub-second part",
				time.Unix(-1231006505, 999_999_999),
				time.Date(1930, 12, 29, 5, 44, 55, 0, time.UTC),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				ts := timeutil.NewTimestamp(tc.in)
				require.Equal(t, tc.want, ts.Time())
			})
		}
	})
}

func TestNewTimestampFromUnix(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   int64
			want time.Time
		}{
			{
				"zero",
				0,
				time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				"positive",
				1231006505,
				time.Date(2009, 1, 3, 18, 15, 5, 0, time.UTC),
			},
			{
				"negative",
				-1231006505,
				time.Date(1930, 12, 29, 5, 44, 55, 0, time.UTC),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				ts := timeutil.NewTimestampFromUnix(tc.in)
				require.Equal(t, tc.want, ts.Time())
			})
		}
	})
}

func TestTimestamp_Before(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		type input struct {
			ts  timeutil.Timestamp
			ts2 timeutil.Timestamp
		}

		tcs := []struct {
			name string
			in   input
			want bool
		}{
			{
				"ts == ts2",
				input{
					timeutil.NewTimestampFromUnix(1231006505),
					timeutil.NewTimestampFromUnix(1231006505),
				},
				false,
			},
			{
				"ts and ts2 within the same second",
				input{
					timeutil.NewTimestamp(time.Date(2009, 1, 3, 18, 15, 5, 100_000_000, time.UTC)),
					timeutil.NewTimestamp(time.Date(2009, 1, 3, 18, 15, 5, 900_000_000, time.UTC)),
				},
				false,
			},
			{
				"ts < ts2",
				input{
					timeutil.NewTimestampFromUnix(1231006505),
					timeutil.NewTimestampFromUnix(1231006506),
				},
				true,
			},
			{
				"ts > ts2",
				input{
					timeutil.NewTimestampFromUnix(1231006505),
					timeutil.NewTimestampFromUnix(1231006504),
				},
				false,
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				ok := tc.in.ts.Before(tc.in.ts2)
				require.Equal(t, tc.want, ok)
			})
		}
	})
}

func TestTimestamp_After(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		type input struct {
			ts  timeutil.Timestamp
			ts2 timeutil.Timestamp
		}

		tcs := []struct {
			name string
			in   input
			want bool
		}{
			{
				"ts == ts2",
				input{
					timeutil.NewTimestampFromUnix(1231006505),
					timeutil.NewTimestampFromUnix(1231006505),
				},
				false,
			},
			{
				"ts and ts2 within the same second",
				input{
					timeutil.NewTimestamp(time.Date(2009, 1, 3, 18, 15, 5, 100_000_000, time.UTC)),
					timeutil.NewTimestamp(time.Date(2009, 1, 3, 18, 15, 5, 900_000_000, time.UTC)),
				},
				false,
			},
			{
				"ts > ts2",
				input{
					timeutil.NewTimestampFromUnix(1231006505),
					timeutil.NewTimestampFromUnix(1231006504),
				},
				true,
			},
			{
				"ts < ts2",
				input{
					timeutil.NewTimestampFromUnix(1231006505),
					timeutil.NewTimestampFromUnix(1231006506),
				},
				false,
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				ok := tc.in.ts.After(tc.in.ts2)
				require.Equal(t, tc.want, ok)
			})
		}
	})
}

func TestTimestamp_Add(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		type input struct {
			ts timeutil.Timestamp
			d  time.Duration
		}

		tcs := []struct {
			name string
			in   input
			want time.Time
		}{
			{
				"zero",
				input{
					timeutil.NewTimestampFromUnix(1231006505),
					0,
				},
				time.Date(2009, 1, 3, 18, 15, 5, 0, time.UTC),
			},
			{
				"positive",
				input{
					timeutil.NewTimestampFromUnix(1231006505),
					time.Second,
				},
				time.Date(2009, 1, 3, 18, 15, 6, 0, time.UTC),
			},
			{
				"negative",
				input{
					timeutil.NewTimestampFromUnix(1231006505),
					-time.Second,
				},
				time.Date(2009, 1, 3, 18, 15, 4, 0, time.UTC),
			},
			{
				"positive with sub-second part",
				input{
					timeutil.NewTimestampFromUnix(1231006505),
					1500 * time.Millisecond,
				},
				time.Date(2009, 1, 3, 18, 15, 6, 0, time.UTC),
			},
			{
				"negative with sub-second part",
				input{
					timeutil.NewTimestampFromUnix(1231006505),
					-1500 * time.Millisecond,
				},
				time.Date(2009, 1, 3, 18, 15, 3, 0, time.UTC),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				tsBefore := tc.in.ts
				tsAfter := tsBefore.Add(tc.in.d)
				require.Equal(t, tc.in.ts.Time(), tsBefore.Time())
				require.Equal(t, tc.want, tsAfter.Time())
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
				"unsupported source: nil",
			},
			{
				"time",
				time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
				"unsupported source type: time.Time",
			},
			{
				"uint64: exceeds int64 range",
				uint64(math.MaxInt64) + 1,
				"invalid uint64 source: exceeds int64 range",
			},
			{
				"bytes: empty",
				[]byte{},
				"invalid decimal string: empty",
			},
			{
				"bytes: invalid",
				[]byte("invalid"),
				"invalid decimal string",
			},
			{
				"decimal string bytes: fractional",
				[]byte("1231006505.0"),
				"invalid decimal string",
			},
			{
				"decimal string bytes: exponential",
				[]byte("1231006505e0"),
				"invalid decimal string",
			},
			{
				"decimal string bytes: contains underscores",
				[]byte("1_231_006_505"),
				"invalid decimal string",
			},
			{
				"decimal string bytes: exceeds int64 range",
				[]byte("9223372036854775808"),
				"invalid decimal string",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var ts timeutil.Timestamp
				err := ts.Scan(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			want time.Time
		}{
			{
				"int64: zero",
				int64(0),
				time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				"int64: positive",
				int64(1231006505),
				time.Date(2009, 1, 3, 18, 15, 5, 0, time.UTC),
			},
			{
				"int64: negative",
				int64(-1231006505),
				time.Date(1930, 12, 29, 5, 44, 55, 0, time.UTC),
			},
			{
				"uint64",
				uint64(1231006505),
				time.Date(2009, 1, 3, 18, 15, 5, 0, time.UTC),
			},
			{
				"decimal string bytes: unsigned",
				[]byte("1231006505"),
				time.Date(2009, 1, 3, 18, 15, 5, 0, time.UTC),
			},
			{
				"decimal string bytes: signed positive",
				[]byte("+1231006505"),
				time.Date(2009, 1, 3, 18, 15, 5, 0, time.UTC),
			},
			{
				"decimal string bytes: signed negative",
				[]byte("-1231006505"),
				time.Date(1930, 12, 29, 5, 44, 55, 0, time.UTC),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var ts timeutil.Timestamp
				err := ts.Scan(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want, ts.Time())
			})
		}
	})
}

func TestTimestamp_ValueScanRoundTrip(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   timeutil.Timestamp
		}{
			{
				"zero",
				timeutil.NewTimestampFromUnix(0),
			},
			{
				"positive",
				timeutil.NewTimestampFromUnix(1231006505),
			},
			{
				"negative",
				timeutil.NewTimestampFromUnix(-1231006505),
			},
			{
				"positive with sub-second part",
				timeutil.NewTimestamp(time.Date(2009, 1, 3, 18, 15, 5, 999_999_999, time.UTC)),
			},
			{
				"negative with sub-second part",
				timeutil.NewTimestamp(time.Unix(-1231006505, 999_999_999)),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				v, err := tc.in.Value()
				require.NoError(t, err)

				var ts timeutil.Timestamp
				err = ts.Scan(v)
				require.NoError(t, err)
				require.Equal(t, tc.in.Time(), ts.Time())
			})
		}
	})
}

func TestTimestamp_MarshalText(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   timeutil.Timestamp
			want []byte
		}{
			{
				"zero",
				timeutil.NewTimestampFromUnix(0),
				[]byte("0"),
			},
			{
				"positive",
				timeutil.NewTimestampFromUnix(1231006505),
				[]byte("1231006505"),
			},
			{
				"negative",
				timeutil.NewTimestampFromUnix(-1231006505),
				[]byte("-1231006505"),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				b, err := tc.in.MarshalText()
				require.NoError(t, err)
				require.Equal(t, tc.want, b)
			})
		}
	})
}

func TestTimestamp_JSONMarshaling(t *testing.T) {
	encs := []struct {
		name    string
		marshal func(timeutil.Timestamp) ([]byte, error)
	}{
		{
			"json.Marshal",
			func(ts timeutil.Timestamp) ([]byte, error) {
				return json.Marshal(ts)
			},
		},
		{
			"MarshalJSON",
			func(ts timeutil.Timestamp) ([]byte, error) {
				return ts.MarshalJSON()
			},
		},
	}

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
				for _, enc := range encs {
					t.Run(enc.name, func(t *testing.T) {
						b, err := enc.marshal(tc.in)
						require.NoError(t, err)
						require.Equal(t, tc.want, b)
					})
				}
			})
		}
	})
}

func TestTimestamp_MarshalGQL(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   timeutil.Timestamp
			want string
		}{
			{
				"zero",
				timeutil.NewTimestampFromUnix(0),
				`"0"`,
			},
			{
				"positive",
				timeutil.NewTimestampFromUnix(1231006505),
				`"1231006505"`,
			},
			{
				"negative",
				timeutil.NewTimestampFromUnix(-1231006505),
				`"-1231006505"`,
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var buf bytes.Buffer
				tc.in.MarshalGQL(&buf)
				require.Equal(t, tc.want, buf.String())
			})
		}
	})
}

func TestTimestamp_UnmarshalText(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want string
		}{
			{
				"nil",
				nil,
				"invalid decimal string: empty",
			},
			{
				"bytes: empty",
				[]byte{},
				"invalid decimal string: empty",
			},
			{
				"string bytes: invalid",
				[]byte("invalid"),
				"invalid decimal string",
			},
			{
				"decimal string bytes: fractional",
				[]byte("1231006505.0"),
				"invalid decimal string",
			},
			{
				"decimal string bytes: exponential",
				[]byte("1231006505e0"),
				"invalid decimal string",
			},
			{
				"decimal string bytes: contains underscores",
				[]byte("1_231_006_505"),
				"invalid decimal string",
			},
			{
				"decimal string bytes: exceeds int64 range",
				[]byte("9223372036854775808"),
				"invalid decimal string",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var ts timeutil.Timestamp
				err := ts.UnmarshalText(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want time.Time
		}{
			{
				"decimal string bytes: zero",
				[]byte("0"),
				time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				"decimal string bytes: unsigned",
				[]byte("1231006505"),
				time.Date(2009, 1, 3, 18, 15, 5, 0, time.UTC),
			},
			{
				"decimal string bytes: signed positive",
				[]byte("+1231006505"),
				time.Date(2009, 1, 3, 18, 15, 5, 0, time.UTC),
			},
			{
				"decimal string bytes: signed negative",
				[]byte("-1231006505"),
				time.Date(1930, 12, 29, 5, 44, 55, 0, time.UTC),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var ts timeutil.Timestamp
				err := ts.UnmarshalText(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want, ts.Time())
			})
		}
	})
}

func TestTimestamp_JSONUnmarshaling(t *testing.T) {
	decs := []struct {
		name      string
		unmarshal func([]byte, *timeutil.Timestamp) error
	}{
		{
			"json.Unmarshal",
			func(b []byte, ts *timeutil.Timestamp) error {
				return json.Unmarshal(b, ts)
			},
		},
		{
			"UnmarshalJSON",
			func(b []byte, ts *timeutil.Timestamp) error {
				return ts.UnmarshalJSON(b)
			},
		},
	}

	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want string
		}{
			{
				"nil",
				nil,
				"",
			},
			{
				"bytes: empty",
				[]byte{},
				"",
			},
			{
				"null",
				[]byte(`null`),
				"unsupported json token kind: null",
			},
			{
				"unquoted decimal string bytes: signed positive",
				[]byte(`+1231006505`),
				"unsupported json token kind: invalid",
			},
			{
				"unquoted decimal string bytes: truncated",
				[]byte(`0.`),
				"failed to read token",
			},
			{
				"unquoted decimal string bytes: fractional",
				[]byte(`1231006505.0`),
				"invalid decimal string",
			},
			{
				"unquoted decimal string bytes: exponential",
				[]byte(`1231006505e0`),
				"invalid decimal string",
			},
			{
				"unquoted decimal string bytes: exceeds int64 range",
				[]byte(`9223372036854775808`),
				"invalid decimal string",
			},
			{
				"quoted decimal string bytes: truncated",
				[]byte(`"0`),
				"failed to read token",
			},
			{
				"quoted string bytes: empty",
				[]byte(`""`),
				"invalid decimal string: empty",
			},
			{
				"quoted decimal string bytes: fractional",
				[]byte(`"1231006505.0"`),
				"invalid decimal string",
			},
			{
				"quoted decimal string bytes: exponential",
				[]byte(`"1231006505e0"`),
				"invalid decimal string",
			},
			{
				"quoted string bytes: exceeds int64 range",
				[]byte(`"9223372036854775808"`),
				"invalid decimal string",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				for _, dec := range decs {
					t.Run(dec.name, func(t *testing.T) {
						var ts timeutil.Timestamp
						err := dec.unmarshal(tc.in, &ts)
						require.ErrorContains(t, err, tc.want)
					})
				}
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want time.Time
		}{
			{
				"unquoted decimal string bytes: zero",
				[]byte(`0`),
				time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				"unquoted decimal string bytes: unsigned",
				[]byte(`1231006505`),
				time.Date(2009, 1, 3, 18, 15, 5, 0, time.UTC),
			},
			{
				"unquoted decimal string bytes: signed negative",
				[]byte(`-1231006505`),
				time.Date(1930, 12, 29, 5, 44, 55, 0, time.UTC),
			},
			{
				"quoted decimal string bytes: zero",
				[]byte(`"0"`),
				time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				"quoted decimal string bytes: unsigned",
				[]byte(`"1231006505"`),
				time.Date(2009, 1, 3, 18, 15, 5, 0, time.UTC),
			},
			{
				"quoted decimal string bytes: signed positive",
				[]byte(`"+1231006505"`),
				time.Date(2009, 1, 3, 18, 15, 5, 0, time.UTC),
			},
			{
				"quoted decimal string bytes: signed negative",
				[]byte(`"-1231006505"`),
				time.Date(1930, 12, 29, 5, 44, 55, 0, time.UTC),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				for _, dec := range decs {
					t.Run(dec.name, func(t *testing.T) {
						var ts timeutil.Timestamp
						err := dec.unmarshal(tc.in, &ts)
						require.NoError(t, err)
						require.Equal(t, tc.want, ts.Time())
					})
				}
			})
		}
	})
}

func TestTimestamp_UnmarshalGQL(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			want string
		}{
			{
				"nil",
				nil,
				"unsupported value: nil",
			},
			{
				"int",
				int(0),
				"unsupported value type: int",
			},
			{
				"string: empty",
				"",
				"invalid decimal string: empty",
			},
			{
				"string: invalid",
				"invalid",
				"invalid decimal string",
			},
			{
				"decimal string: fractional",
				"1231006505.0",
				"invalid decimal string",
			},
			{
				"decimal string: exponential",
				"1231006505e0",
				"invalid decimal string",
			},
			{
				"decimal string: contains underscores",
				"1_231_006_505",
				"invalid decimal string",
			},
			{
				"decimal string: exceeds int64 range",
				"9223372036854775808",
				"invalid decimal string",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var ts timeutil.Timestamp
				err := ts.UnmarshalGQL(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			want time.Time
		}{
			{
				"decimal string: zero",
				"0",
				time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				"decimal string: unsigned",
				"1231006505",
				time.Date(2009, 1, 3, 18, 15, 5, 0, time.UTC),
			},
			{
				"decimal string: signed positive",
				"+1231006505",
				time.Date(2009, 1, 3, 18, 15, 5, 0, time.UTC),
			},
			{
				"decimal string: signed negative",
				"-1231006505",
				time.Date(1930, 12, 29, 5, 44, 55, 0, time.UTC),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var ts timeutil.Timestamp
				err := ts.UnmarshalGQL(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want, ts.Time())
			})
		}
	})
}
