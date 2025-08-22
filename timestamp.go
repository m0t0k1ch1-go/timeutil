package timeutil

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"time"
)

// Timestamp represents a point in time in UTC.
type Timestamp struct {
	t time.Time
}

// NewTimestamp returns a new Timestamp.
func NewTimestamp(t time.Time) Timestamp {
	var ts Timestamp
	ts.setTime(t)

	return ts
}

func (ts *Timestamp) setTime(t time.Time) {
	ts.t = t.In(time.UTC)
}

// NewTimestampFromUnix returns a new Timestamp from an int64 Unix timestamp in seconds.
func NewTimestampFromUnix(sec int64) Timestamp {
	return NewTimestamp(time.Unix(sec, 0))
}

// Time returns the underlying time.Time.
func (ts Timestamp) Time() time.Time {
	return ts.t
}

// Unix returns the Unix timestamp in seconds as an int64.
func (ts Timestamp) Unix() int64 {
	return ts.t.Unix()
}

// String implements fmt.Stringer.
// It returns the Unix timestamp in seconds as a decimal string.
func (ts Timestamp) String() string {
	return strconv.FormatInt(ts.Unix(), 10)
}

// Value implements driver.Valuer.
// It returns the Unix timestamp in seconds as an int64.
func (ts Timestamp) Value() (driver.Value, error) {
	return ts.Unix(), nil
}

// Scan implements sql.Scanner.
// It accepts a Unix timestamp in seconds as one of:
//   - int64
//   - uint64 (<= math.MaxInt64)
//   - []byte (decimal string)
func (ts *Timestamp) Scan(src any) error {
	if src == nil {
		return errors.New("invalid source: nil")
	}

	switch v := src.(type) {

	case int64:
		ts.setTime(time.Unix(v, 0))

		return nil

	case uint64:
		if v > math.MaxInt64 {
			return errors.New("invalid source: exceeds int64 range")
		}

		ts.setTime(time.Unix(int64(v), 0))

		return nil

	case []byte:
		if len(v) == 0 {
			return errors.New("invalid source: empty []byte")
		}

		i, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return fmt.Errorf("invalid source: %w", err)
		}

		ts.setTime(time.Unix(i, 0))

		return nil

	default:
		return fmt.Errorf("unsupported source type: %T", src)
	}
}

// MarshalJSON implements json.Marshaler.
// It returns the Unix timestamp in seconds as a JSON number.
func (ts Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(ts.Unix())
}

// MarshalGQL implements graphql.Marshaler.
// It writes the Unix timestamp in seconds as a quoted decimal string.
func (ts Timestamp) MarshalGQL(w io.Writer) {
	_, _ = io.WriteString(w, strconv.Quote(strconv.FormatInt(ts.Unix(), 10)))
}

// UnmarshalJSON implements json.Unmarshaler.
// It accepts a JSON number (integer Unix timestamp in seconds).
func (ts *Timestamp) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return errors.New("invalid json value: empty")
	}
	if string(b) == "null" {
		return errors.New("invalid json value: null")
	}

	var i int64
	if err := json.Unmarshal(b, &i); err != nil {
		return fmt.Errorf("invalid json number: %w", err)
	}

	ts.setTime(time.Unix(i, 0))

	return nil
}

// UnmarshalGQL implements graphql.Unmarshaler.
// It accepts a GraphQL String (decimal Unix timestamp in seconds).
func (ts *Timestamp) UnmarshalGQL(v any) error {
	if v == nil {
		return errors.New("invalid graphql value: nil")
	}

	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("unsupported graphql value type: %T", v)
	}
	if len(s) == 0 {
		return errors.New("invalid graphql string: empty")
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid graphql string: %w", err)
	}

	ts.setTime(time.Unix(i, 0))

	return nil
}
