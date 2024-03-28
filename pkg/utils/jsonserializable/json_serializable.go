package jsonserializable

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	pkgerrors "github.com/pkg/errors"
)

type JSONSerializable struct {
	Val   interface{}
	Valid bool
}

func ReinterpretJSONNumbers(val interface{}) (interface{}, error) {
	switch v := val.(type) {
	case json.Number:
		return getJSONNumberValue(v)
	case []interface{}:
		s := make([]interface{}, len(v))
		for i, vv := range v {
			ival, ierr := ReinterpretJSONNumbers(vv)
			if ierr != nil {
				return nil, ierr
			}
			s[i] = ival
		}
		return s, nil
	case map[string]interface{}:
		m := make(map[string]interface{}, len(v))
		for k, vv := range v {
			ival, ierr := ReinterpretJSONNumbers(vv)
			if ierr != nil {
				return nil, ierr
			}
			m[k] = ival
		}
		return m, nil
	}
	return val, nil
}

// UnmarshalJSON implements custom unmarshaling logic
func (js *JSONSerializable) UnmarshalJSON(bs []byte) error {
	if js == nil {
		*js = JSONSerializable{}
	}
	if len(bs) == 0 {
		js.Valid = false
		return nil
	}

	var decoded interface{}
	d := json.NewDecoder(bytes.NewReader(bs))
	d.UseNumber()
	if err := d.Decode(&decoded); err != nil {
		return err
	}

	if decoded != nil {
		reinterpreted, err := ReinterpretJSONNumbers(decoded)
		if err != nil {
			return err
		}

		*js = JSONSerializable{
			Valid: true,
			Val:   reinterpreted,
		}
	}

	return nil
}

// MarshalJSON implements custom marshaling logic
func (js JSONSerializable) MarshalJSON() ([]byte, error) {
	if !js.Valid {
		return json.Marshal(nil)
	}
	jsWithHex := replaceBytesWithHex(js.Val)
	return json.Marshal(jsWithHex)
}

func (js *JSONSerializable) Scan(value interface{}) error {
	if value == nil {
		*js = JSONSerializable{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return pkgerrors.Errorf("JSONSerializable#Scan received a value of type %T", value)
	}
	if js == nil {
		*js = JSONSerializable{}
	}
	return js.UnmarshalJSON(bytes)
}

func (js JSONSerializable) Value() (driver.Value, error) {
	if !js.Valid {
		return nil, nil
	}
	return js.MarshalJSON()
}

func (js *JSONSerializable) Empty() bool {
	return js == nil || !js.Valid
}

// replaceBytesWithHex replaces all []byte with hex-encoded strings
func replaceBytesWithHex(val interface{}) interface{} {
	switch value := val.(type) {
	case nil:
		return value
	case []byte:
		return stringToHex(string(value))
	case common.Address:
		return value.Hex()
	case common.Hash:
		return value.Hex()
	case [][]byte:
		var list []string
		for _, bytes := range value {
			list = append(list, stringToHex(string(bytes)))
		}
		return list
	case []common.Address:
		var list []string
		for _, addr := range value {
			list = append(list, addr.Hex())
		}
		return list
	case []common.Hash:
		var list []string
		for _, hash := range value {
			list = append(list, hash.Hex())
		}
		return list
	case []interface{}:
		if value == nil {
			return value
		}
		var list []interface{}
		for _, item := range value {
			list = append(list, replaceBytesWithHex(item))
		}
		return list
	case map[string]interface{}:
		if value == nil {
			return value
		}
		m := make(map[string]interface{})
		for k, v := range value {
			m[k] = replaceBytesWithHex(v)
		}
		return m
	default:
		// This handles solidity types: bytes1..bytes32,
		// which map to [1]uint8..[32]uint8 when decoded.
		// We persist them as hex strings, and we know ETH ABI encoders
		// can parse hex strings, same as BytesParam does.
		if s := uint8ArrayToSlice(value); s != nil {
			return replaceBytesWithHex(s)
		}
		return value
	}
}

// uint8ArrayToSlice converts [N]uint8 array to slice.
func uint8ArrayToSlice(arr interface{}) interface{} {
	t := reflect.TypeOf(arr)
	if t.Kind() != reflect.Array || t.Elem().Kind() != reflect.Uint8 {
		return nil
	}
	v := reflect.ValueOf(arr)
	s := reflect.MakeSlice(reflect.SliceOf(t.Elem()), v.Len(), v.Len())
	reflect.Copy(s, v)
	return s.Interface()
}

func getJSONNumberValue(value json.Number) (interface{}, error) {
	var result interface{}

	bn, ok := new(big.Int).SetString(value.String(), 10)
	if ok {
		if bn.IsInt64() {
			result = bn.Int64()
		} else if bn.IsUint64() {
			result = bn.Uint64()
		} else {
			result = bn
		}
	} else {
		f, err := value.Float64()
		if err != nil {
			return nil, pkgerrors.Errorf("failed to parse json.Value: %v", err)
		}
		result = f
	}

	return result, nil
}

func stringToHex(in string) string {
	str := hex.EncodeToString([]byte(in))
	if len(str) < 2 || len(str) > 1 && strings.ToLower(str[0:2]) != "0x" {
		str = "0x" + str
	}
	return str
}
