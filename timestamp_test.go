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
			want int64
		}{
			{
				"int64: zero",
				int64(0),
				0,
			},
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
				"decimal string bytes: unsigned",
				[]byte("1231006505"),
				1231006505,
			},
			{
				"decimal string bytes: signed positive",
				[]byte("+1231006505"),
				1231006505,
			},
			{
				"decimal string bytes: signed negative",
				[]byte("-1231006505"),
				-1231006505,
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
			want int64
		}{
			{
				"decimal string: zero",
				[]byte("0"),
				0,
			},
			{
				"decimal string: unsigned",
				[]byte("1231006505"),
				1231006505,
			},
			{
				"decimal string: signed positive",
				[]byte("+1231006505"),
				1231006505,
			},
			{
				"decimal string: signed negative",
				[]byte("-1231006505"),
				-1231006505,
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var ts timeutil.Timestamp
				err := ts.UnmarshalText(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want, ts.Unix())
				require.Equal(t, time.UTC, ts.Time().Location())
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
				"empty",
				[]byte{},
				"",
			},
			{
				"null",
				[]byte(`null`),
				"unsupported json token kind: null",
			},
			{
				"quoted decimal string",
				[]byte(`"0"`),
				"unsupported json token kind: string",
			},
			{
				"unquoted decimal string: signed positive",
				[]byte(`+1231006505`),
				"unsupported json token kind: invalid",
			},
			{
				"unquoted decimal string: truncated",
				[]byte(`0.`),
				"failed to read value",
			},
			{
				"unquoted decimal string: fractional",
				[]byte(`1231006505.0`),
				"invalid decimal string",
			},
			{
				"unquoted decimal string: exponential",
				[]byte(`1231006505e0`),
				"invalid decimal string",
			},
			{
				"unquoted decimal string: exceeds int64 range",
				[]byte(`9223372036854775808`),
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
			want int64
		}{
			{
				"unquoted decimal string: zero",
				[]byte(`0`),
				0,
			},
			{
				"unquoted decimal string: unsigned",
				[]byte(`1231006505`),
				1231006505,
			},
			{
				"unquoted decimal string: signed negative",
				[]byte(`-1231006505`),
				-1231006505,
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				for _, dec := range decs {
					t.Run(dec.name, func(t *testing.T) {
						var ts timeutil.Timestamp
						err := dec.unmarshal(tc.in, &ts)
						require.NoError(t, err)
						require.Equal(t, tc.want, ts.Unix())
						require.Equal(t, time.UTC, ts.Time().Location())
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
			want int64
		}{
			{
				"decimal string: zero",
				"0",
				0,
			},
			{
				"decimal string: unsigned",
				"1231006505",
				1231006505,
			},
			{
				"decimal string: signed positive",
				"+1231006505",
				1231006505,
			},
			{
				"decimal string: signed negative",
				"-1231006505",
				-1231006505,
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var ts timeutil.Timestamp
				err := ts.UnmarshalGQL(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want, ts.Unix())
				require.Equal(t, time.UTC, ts.Time().Location())
			})
		}
	})
}
