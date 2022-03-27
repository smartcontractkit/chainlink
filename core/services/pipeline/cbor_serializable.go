package pipeline

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/cbor"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// CBORSerializable implements type-safe serialization.
type CBORSerializable struct {
	Val   interface{}
	Valid bool
}

// Scan is used by sql driver to read "bytea" value deserialize into CBORSerializable.
func (s *CBORSerializable) Scan(value interface{}) error {
	if value == nil {
		*s = CBORSerializable{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.Errorf("CBORSerializable#Scan received a value of type %T", value)
	}
	if s == nil {
		*s = CBORSerializable{}
	}
	parsed, err := cbor.ParseStandardCBOR(bytes)
	if err == nil {
		coerced, err2 := cbor.CoerceInterfaceMapToStringMap(parsed)
		if err2 != nil {
			err = errors.Wrap(err2, "error in cbor.CoerceInterfaceMapToStringMap")
		} else {
			s.Val = coerced
		}
	}
	s.Valid = err == nil
	return err
}

// Value is used by sql driver to serialize the CBORSerializable into "bytea" value.
func (s CBORSerializable) Value() (driver.Value, error) {
	if !s.Valid {
		return nil, nil
	}
	return cbor.EncodeCBOR(s.Val)
}

// Empty returns true if this CBORSerializable is either nil or invalid.
func (s *CBORSerializable) Empty() bool {
	return s == nil || !s.Valid
}

// UnmarshalJSON implements custom unmarshaling logic, used by web presenters.
// TODO: change web presenters to use cbor?
func (s *CBORSerializable) UnmarshalJSON(bs []byte) error {
	if s == nil {
		*s = CBORSerializable{}
	}
	str := string(bs)
	if str == "" || str == "null" {
		s.Valid = false
		return nil
	}

	err := json.Unmarshal(bs, &s.Val)
	s.Valid = err == nil
	return err
}

// MarshalJSON implements custom marshaling logic, used by web presenters.
// TODO: change web presenters to use cbor?
func (s CBORSerializable) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	switch x := s.Val.(type) {
	case []byte:
		// Don't need to HEX encode if it is a valid JSON string
		if json.Valid(x) {
			return json.Marshal(string(x))
		}

		// Don't need to HEX encode if it is already HEX encoded value
		if utils.IsHexBytes(x) {
			return json.Marshal(string(x))
		}

		return json.Marshal(hex.EncodeToString(x))
	default:
		return json.Marshal(s.Val)
	}
}
