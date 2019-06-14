package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type USmall struct {
	Uint32 uint32
	Valid  bool
}

func NewUSmall(i uint32, valid bool) USmall {
	return USmall{
		Uint32: i,
		Valid:  valid,
	}
}

// USmallFrom creates a new USmall that will always be valid.
func USmallFrom(i uint32) USmall {
	return NewUSmall(i, true)
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Int.
// It also supports unmarshalling a sql.NullInt64.
func (i *USmall) UnmarshalJSON(data []byte) error {
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
		str := string(x)
		if len(str) == 0 {
			i.Valid = false
			return nil
		}
		i.Uint32, err = parse(str)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.USmall", reflect.TypeOf(v).Name())
	}
	i.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null USmall if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *USmall) UnmarshalText(text []byte) error {
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
// It will encode null if this USmall is null.
func (i USmall) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatUint(uint64(i.Uint32), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this USmall is null.
func (i USmall) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatUint(uint64(i.Uint32), 10)), nil
}

// SetValid changes this USmall's value and also sets it to be non-null.
func (i *USmall) SetValid(n uint32) {
	i.Uint32 = n
	i.Valid = true
}

// Value returns this instance serialized for database storage.
func (i USmall) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}
	return i.Uint32, nil
}

// Scan reads the database value and returns an instance.
func (i *USmall) Scan(value interface{}) error {
	if value == nil {
		*i = USmall{}
		return nil
	}

	switch typed := value.(type) {
	case int:
		safe := uint32(typed)
		if int(safe) != typed {
			return fmt.Errorf("Unable to convert %v of %T to USmall; overflow", value, value)
		}
		*i = USmallFrom(safe)
	case int64:
		safe := uint32(typed)
		if int64(safe) != typed {
			return fmt.Errorf("Unable to convert %v of %T to USmall; overflow", value, value)
		}
		*i = USmallFrom(safe)
	default:
		return fmt.Errorf("Unable to convert %v of %T to USmall", value, value)
	}
	return nil
}
