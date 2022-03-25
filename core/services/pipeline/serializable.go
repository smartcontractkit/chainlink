package pipeline

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/cbor"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Serializable implements type-safe serialization.
type Serializable struct {
	Val   interface{}
	Valid bool
}

// NewValidSerializable creates new valid Serializable object.
func NewValidSerializable(value interface{}) *Serializable {
	return &Serializable{
		Val:   value,
		Valid: true,
	}
}

// Scan is used by sql driver to read "bytea" value deserialize into Serializable.
func (s *Serializable) Scan(value interface{}) error {
	if value == nil {
		*s = Serializable{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.Errorf("Serializable#Scan received a value of type %T", value)
	}
	if s == nil {
		*s = Serializable{}
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

// Value is used by sql driver to serialize the Serializable into "bytea" value.
func (s Serializable) Value() (driver.Value, error) {
	if !s.Valid {
		return nil, nil
	}
	return cbor.EncodeCBOR(s.Val)
}

// Empty returns true if this Serializable is either nil or invalid.
func (s *Serializable) Empty() bool {
	return s == nil || !s.Valid
}

// String renders the Serializable to human-friendly JSON string for debugging purposes.
// For an invalid Serializable it outputs "invalid".
// Any []byte values will be rendered as hex-encoded strings.
// Any [][]byte values will be rendered as list of hex-encoded strings.
func (s Serializable) String() string {
	if s.Valid {
		processedVal := replaceBytesToHex(s.Val)
		b, err := json.Marshal(processedVal)
		if err == nil {
			return string(b)
		}
	}
	return "invalid"
}

func replaceBytesToHex(val interface{}) interface{} {
	if val == nil {
		return "nil"
	}
	switch value := val.(type) {
	case []byte:
		return utils.StringToHex(string(value))
	case [][]byte:
		var list []string
		for _, bytes := range value {
			list = append(list, utils.StringToHex(string(bytes)))
		}
		return list
	case []interface{}:
		var list []interface{}
		for _, item := range value {
			list = append(list, replaceBytesToHex(item))
		}
		return list
	case map[string]interface{}:
		mapCopy := utils.CopyMap(value)
		for k, v := range mapCopy {
			mapCopy[k] = replaceBytesToHex(v)
		}
		return mapCopy
	default:
		return val
	}
}
