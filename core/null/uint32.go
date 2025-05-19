package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// Uint32 encapsulates the value and validity (not null) of a uint32 value,
// to differentiate nil from 0 in json and sql.
type Uint32 struct {
	Uint32 uint32
	Valid  bool
}

// NewUint32 returns an instance of Uint32 with the passed parameters.
func NewUint32(i uint32, valid bool) Uint32 {
	return Uint32{
		Uint32: i,
		Valid:  valid,
	}
}

// Uint32From creates a new Uint32 that will always be valid.
func Uint32From(i uint32) Uint32 {
	return NewUint32(i, true)
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Int.
func (i *Uint32) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		// Unmarshal again, directly to value, to avoid intermediate float64
		err = json.Unmarshal(data, &i.Uint32)
	case string:
		str := x
		if len(str) == 0 {
			i.Valid = false
			return nil
		}
		i.Uint32, err = parse(str)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Uint32", reflect.TypeOf(v).Name())
	}
	i.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Uint32 if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Uint32) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	i.Uint32, err = parse(string(text))
	i.Valid = err == nil
	return err
}

func parse(str string) (uint32, error) {
	v, err := strconv.ParseUint(str, 10, 32)
	return uint32(v), err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Uint32 is null.
func (i Uint32) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatUint(uint64(i.Uint32), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Uint32 is null.
func (i Uint32) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatUint(uint64(i.Uint32), 10)), nil
}

// SetValid changes this Uint32's value and also sets it to be non-null.
func (i *Uint32) SetValid(n uint32) {
	i.Uint32 = n
	i.Valid = true
}

// Value returns this instance serialized for database storage.
func (i Uint32) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}

	// golang's sql driver types as determined by IsValue only supports:
	// []byte, bool, float64, int64, string, time.Time
	// https://golang.org/src/database/sql/driver/types.go
	return int64(i.Uint32), nil
}

// Scan reads the database value and returns an instance.
func (i *Uint32) Scan(value interface{}) error {
	if value == nil {
		*i = Uint32{}
		return nil
	}

	switch typed := value.(type) {
	case int:
		safe := uint32(typed)
		if int(safe) != typed {
			return fmt.Errorf("unable to convert %v of %T to Uint32; overflow", value, value)
		}
		*i = Uint32From(safe)
	case int64:
		safe := uint32(typed)
		if int64(safe) != typed {
			return fmt.Errorf("unable to convert %v of %T to Uint32; overflow", value, value)
		}
		*i = Uint32From(safe)
	case uint:
		safe := uint32(typed)
		if uint(safe) != typed {
			return fmt.Errorf("unable to convert %v of %T to Uint32; overflow", value, value)
		}
		*i = Uint32From(safe)
	case uint32:
		safe := typed
		*i = Uint32From(safe)
	default:
		return fmt.Errorf("unable to convert %v of %T to Uint32", value, value)
	}
	return nil
}
