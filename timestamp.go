package timeutil

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
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

// Time returns the underlying time.Time.
func (ts Timestamp) Time() time.Time {
	return ts.t
}

// String implements fmt.Stringer.
// It returns the Unix timestamp in seconds as a decimal string.
func (ts Timestamp) String() string {
	return strconv.FormatInt(ts.t.Unix(), 10)
}

// Value implements driver.Valuer.
// It returns the Unix timestamp in seconds as an int64.
func (ts Timestamp) Value() (driver.Value, error) {
	return ts.t.Unix(), nil
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
	return json.Marshal(ts.t.Unix())
}

// UnmarshalJSON implements json.Unmarshaler.
// It accepts a JSON number representing a Unix timestamp in seconds.
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
