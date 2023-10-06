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

// Now returns the current timestamp.
func Now() Timestamp {
	return Timestamp{
		t: time.Now(),
	}
}

// Time returns the underlying time.Time.
func (t Timestamp) Time() time.Time {
	return t.t
}

// String returns the string representation of the timestamp.
func (t Timestamp) String() string {
	return strconv.FormatInt(t.t.Unix(), 10)
}

// Value implements the driver.Valuer interface.
func (t Timestamp) Value() (driver.Value, error) {
	return t.t.Unix(), nil
}

// Scan implements the sql.Scanner interface.
func (t *Timestamp) Scan(src any) error {
	switch v := src.(type) {
	case int64:
		t.t = time.Unix(v, 0).In(time.UTC)

	case []byte:
		i, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return errors.Wrapf(err, "failed to convert %s into type int64", string(v))
		}

		t.t = time.Unix(i, 0).In(time.UTC)

	default:
		return errors.Errorf("unexpected src type: %T", src)
	}

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(t.t.Unix(), 10)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	i, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return errors.Wrapf(err, "failed to convert %s into type int64", string(b))
	}

	t.t = time.Unix(i, 0).In(time.UTC)

	return nil
}
