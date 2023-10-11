package timeutil

import (
	"database/sql/driver"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// Timestamp is a wrapper for time.Time.
type Timestamp struct {
	t time.Time
}

// Now returns the current Timestamp.
func Now() Timestamp {
	return timeToTimestamp(time.Now())
}

// TimeToTimestamp return the Timestamp that wraps the given time.Time.
func TimeToTimestamp(t time.Time) Timestamp {
	return timeToTimestamp(t)
}

// Time returns the wrapped time.Time.
func (ts Timestamp) Time() time.Time {
	return ts.t
}

// String implements the fmt.Stringer interface.
func (ts Timestamp) String() string {
	return ts.string()
}

// Value implements the driver.Valuer interface.
func (ts Timestamp) Value() (driver.Value, error) {
	return ts.unix(), nil
}

// Scan implements the sql.Scanner interface.
func (ts *Timestamp) Scan(src any) error {
	switch v := src.(type) {
	case int64:
		ts.setUnix(v)
	case []byte:
		ts.setString(string(v))
	default:
		return errors.Errorf("unexpected src type: %T", src)
	}

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (ts Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(ts.string()), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (ts *Timestamp) UnmarshalJSON(b []byte) error {
	return ts.setString(string(b))
}

func timeToTimestamp(t time.Time) Timestamp {
	ts := Timestamp{}
	ts.setTime(t)

	return ts
}

func (ts Timestamp) unix() int64 {
	return ts.t.Unix()
}

func (ts Timestamp) string() string {
	return strconv.FormatInt(ts.unix(), 10)
}

func (ts *Timestamp) setTime(t time.Time) {
	ts.t = t.In(time.UTC)
}

func (ts *Timestamp) setUnix(i int64) {
	ts.setTime(time.Unix(i, 0))
}

func (ts *Timestamp) setString(s string) error {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return errors.Wrapf(err, "failed to convert %s into type int64", s)
	}

	ts.setUnix(i)

	return nil
}
