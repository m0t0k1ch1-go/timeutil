package timeutil

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

var (
	_ fmt.Stringer             = Timestamp{}
	_ driver.Valuer            = Timestamp{}
	_ sql.Scanner              = &Timestamp{}
	_ encoding.TextMarshaler   = Timestamp{}
	_ json.MarshalerTo         = Timestamp{}
	_ json.Marshaler           = Timestamp{}
	_ graphql.Marshaler        = Timestamp{}
	_ encoding.TextUnmarshaler = &Timestamp{}
	_ json.UnmarshalerFrom     = &Timestamp{}
	_ json.Unmarshaler         = &Timestamp{}
	_ graphql.Unmarshaler      = &Timestamp{}
)

// Timestamp represents a point in time in UTC.
type Timestamp struct {
	t time.Time
}

// NewTimestamp returns a new [Timestamp] from a [time.Time].
func NewTimestamp(t time.Time) Timestamp {
	var ts Timestamp
	ts.setTime(t)

	return ts
}

func (ts *Timestamp) setTime(t time.Time) {
	ts.t = t.In(time.UTC)
}

// NewTimestampFromUnix returns a new [Timestamp] from an int64 representing a Unix timestamp in seconds.
func NewTimestampFromUnix(sec int64) Timestamp {
	return NewTimestamp(time.Unix(sec, 0))
}

func (ts *Timestamp) setString(s string) error {
	if len(s) == 0 {
		return errors.New("invalid decimal string: empty")
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid decimal string: %w", err)
	}

	ts.setTime(time.Unix(i, 0))

	return nil
}

// Time returns the underlying [time.Time].
func (ts Timestamp) Time() time.Time {
	return ts.t
}

// Unix returns an int64 representing the Unix timestamp in seconds.
func (ts Timestamp) Unix() int64 {
	return ts.t.Unix()
}

// String implements [fmt.Stringer].
// It encodes ts as a decimal string representing the Unix timestamp in seconds.
func (ts Timestamp) String() string {
	return strconv.FormatInt(ts.Unix(), 10)
}

// Value implements [driver.Valuer].
// It encodes ts as an int64 representing the Unix timestamp in seconds.
func (ts Timestamp) Value() (driver.Value, error) {
	return ts.Unix(), nil
}

// Scan implements [sql.Scanner].
// It decodes one of the following values representing a Unix timestamp in seconds into ts:
//   - int64
//   - uint64 (<= [math.MaxInt64])
//   - decimal string bytes
func (ts *Timestamp) Scan(src any) error {
	if src == nil {
		return errors.New("unsupported source: nil")
	}

	switch src := src.(type) {
	case int64:
		ts.setTime(time.Unix(src, 0))

		return nil

	case uint64:
		if src > math.MaxInt64 {
			return errors.New("invalid uint64 source: exceeds int64 range")
		}

		ts.setTime(time.Unix(int64(src), 0))

		return nil

	case []byte:
		return ts.setString(string(src))

	default:
		return fmt.Errorf("unsupported source type: %T", src)
	}
}

// MarshalText implements [encoding.TextMarshaler].
// It encodes ts as a decimal string representing the Unix timestamp in seconds.
func (ts Timestamp) MarshalText() ([]byte, error) {
	return []byte(ts.String()), nil
}

// MarshalJSONTo implements [json.MarshalerTo].
// It encodes ts as an unquoted decimal string representing the Unix timestamp in seconds and writes it to enc.
func (ts Timestamp) MarshalJSONTo(enc *jsontext.Encoder) error {
	return json.MarshalEncode(enc, ts.Unix())
}

// MarshalJSON implements [json.Marshaler].
// It is like [Timestamp.MarshalJSONTo] but returns the encoded bytes instead of writing them to a [jsontext.Encoder].
func (ts Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(ts)
}

// MarshalGQL implements [graphql.Marshaler].
// It encodes ts as a quoted decimal string representing the Unix timestamp in seconds and writes it to w.
func (ts Timestamp) MarshalGQL(w io.Writer) {
	_, _ = io.WriteString(w, strconv.Quote(ts.String()))
}

// UnmarshalText implements [encoding.TextUnmarshaler].
// It decodes a decimal string representing a Unix timestamp in seconds into ts.
func (ts *Timestamp) UnmarshalText(text []byte) error {
	return ts.setString(string(text))
}

// UnmarshalJSONFrom implements [json.UnmarshalerFrom].
// It decodes an unquoted decimal string representing a Unix timestamp in seconds from dec into ts.
func (ts *Timestamp) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	switch k := dec.PeekKind(); k {
	case jsontext.KindNumber:
		v, err := dec.ReadValue()
		if err != nil {
			return fmt.Errorf("failed to read value: %w", err)
		}

		return ts.UnmarshalText(v)

	default:
		return fmt.Errorf("unsupported json token kind: %v", k)
	}
}

// UnmarshalJSON implements [json.Unmarshaler].
// It is like [Timestamp.UnmarshalJSONFrom] but decodes b instead of reading from a [jsontext.Decoder].
func (ts *Timestamp) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, ts)
}

// UnmarshalGQL implements [graphql.Unmarshaler].
// It decodes a decimal string representing a Unix timestamp in seconds into ts.
func (ts *Timestamp) UnmarshalGQL(v any) error {
	if v == nil {
		return errors.New("unsupported value: nil")
	}

	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("unsupported value type: %T", v)
	}

	return ts.UnmarshalText([]byte(s))
}
