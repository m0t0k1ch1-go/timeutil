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
	return Timestamp{
		t: time.Now(),
	}
}

// Time retuns the Timestamp corresponding to the given time.Time.
func Time(t time.Time) Timestamp {
	return Timestamp{
		t: t,
	}
}

// Time returns the underlying time.Time.
func (ts Timestamp) Time() time.Time {
	return ts.t
}

// String returns the string representation of the Timestamp.
func (ts Timestamp) String() string {
	return strconv.FormatInt(ts.t.Unix(), 10)
}

// Value implements the driver.Valuer interface.
func (ts Timestamp) Value() (driver.Value, error) {
	return ts.t.Unix(), nil
}

// Scan implements the sql.Scanner interface.
func (ts *Timestamp) Scan(src any) error {
	switch v := src.(type) {
	case int64:
		ts.t = time.Unix(v, 0).In(time.UTC)

	case []byte:
		i, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return errors.Wrapf(err, "failed to convert %s into type int64", string(v))
		}

		ts.t = time.Unix(i, 0).In(time.UTC)

	default:
		return errors.Errorf("unexpected src type: %T", src)
	}

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (ts Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(ts.t.Unix(), 10)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (ts *Timestamp) UnmarshalJSON(b []byte) error {
	i, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return errors.Wrapf(err, "failed to convert %s into type int64", string(b))
	}

	ts.t = time.Unix(i, 0).In(time.UTC)

	return nil
}
