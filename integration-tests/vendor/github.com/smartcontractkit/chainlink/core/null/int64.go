package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

// Int64 encapsulates the value and validity (not null) of a int64 value,
// to differentiate nil from 0 in json and sql.
type Int64 struct {
	Int64 int64
	Valid bool
}

// NewInt64 returns an instance of Int64 with the passed parameters.
func NewInt64(i int64, valid bool) Int64 {
	return Int64{
		Int64: i,
		Valid: valid,
	}
}

// Int64From creates a new Int64 that will always be valid.
func Int64From(i int64) Int64 {
	return NewInt64(i, true)
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Int.
func (i *Int64) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		// Unmarshal again, directly to value, to avoid intermediate float64
		err = json.Unmarshal(data, &i.Int64)
	case string:
		str := string(x)
		if len(str) == 0 {
			i.Valid = false
			return nil
		}
		i.Int64, err = parse64(str)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Int64", reflect.TypeOf(v).Name())
	}
	i.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int64 if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Int64) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	i.Int64, err = parse64(string(text))
	i.Valid = err == nil
	return err
}

func parse64(str string) (int64, error) {
	v, err := strconv.ParseInt(str, 10, 64)
	return v, err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Int64 is null.
func (i Int64) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(int64(i.Int64), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Int64 is null.
func (i Int64) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(int64(i.Int64), 10)), nil
}

// SetValid changes this Int64's value and also sets it to be non-null.
func (i *Int64) SetValid(n int64) {
	i.Int64 = n
	i.Valid = true
}

// Value returns this instance serialized for database storage.
func (i Int64) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}

	// golang's sql driver types as determined by IsValue only supports:
	// []byte, bool, float64, int64, string, time.Time
	// https://golang.org/src/database/sql/driver/types.go
	return int64(i.Int64), nil
}

// Scan reads the database value and returns an instance.
func (i *Int64) Scan(value interface{}) error {
	if value == nil {
		*i = Int64{}
		return nil
	}

	switch typed := value.(type) {
	case int:
		safe := int64(typed)
		*i = Int64From(safe)
	case int32:
		safe := int64(typed)
		*i = Int64From(safe)
	case int64:
		safe := int64(typed)
		*i = Int64From(safe)
	case uint:
		if typed > uint(math.MaxInt64) {
			return fmt.Errorf("unable to convert %v of %T to Int64; overflow", value, value)
		}
		safe := int64(typed)
		*i = Int64From(safe)
	case uint64:
		if typed > uint64(math.MaxInt64) {
			return fmt.Errorf("unable to convert %v of %T to Int64; overflow", value, value)
		}
		safe := int64(typed)
		*i = Int64From(safe)
	default:
		return fmt.Errorf("unable to convert %v of %T to Int64", value, value)
	}
	return nil
}
