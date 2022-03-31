package cbor

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/fxamacker/cbor/v2"
	"github.com/pkg/errors"
)

// ParseDietCBOR attempts to coerce the input byte array into valid CBOR.
// Assumes the input is "diet" CBOR which is like CBOR, except:
// 1. It is guaranteed to always be a map
// 2. It may or may not include the opening and closing markers "{}"
func ParseDietCBOR(b []byte) (map[string]interface{}, error) {
	b = autoAddMapDelimiters(b)

	var m map[interface{}]interface{}
	if err := cbor.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	coerced, err := CoerceInterfaceMapToStringMap(m)
	if err != nil {
		return nil, err
	}

	output, ok := coerced.(map[string]interface{})
	if !ok {
		return nil, errors.New("cbor data cannot be coerced to map")
	}

	return output, nil
}

// ParseStandardCBOR parses CBOR in "standards compliant" mode.
// Literal values are passed through "as-is".
// The input is not assumed to be a map.
// Empty inputs will return nil.
func ParseStandardCBOR(b []byte) (a interface{}, err error) {
	if len(b) == 0 {
		return nil, nil
	}
	if err = cbor.Unmarshal(b, &a); err != nil {
		return nil, err
	}
	return
}

// Automatically add missing start map and end map to a CBOR encoded buffer
func autoAddMapDelimiters(b []byte) []byte {
	if len(b) == 0 || (len(b) > 1 && (b[0]>>5) != 5) {
		var buffer bytes.Buffer
		buffer.Write([]byte{0xbf})
		buffer.Write(b)
		buffer.Write([]byte{0xff})
		return buffer.Bytes()
	}

	return b
}

// CoerceInterfaceMapToStringMap converts map[interface{}]interface{} (interface maps) to
// map[string]interface{} (string maps) and []interface{} with interface maps to string maps.
// Relevant when serializing between CBOR and JSON.
//
// It also handles the CBOR 'bignum' type as documented here: https://tools.ietf.org/html/rfc7049#section-2.4.2
func CoerceInterfaceMapToStringMap(in interface{}) (interface{}, error) {
	switch typed := in.(type) {
	case map[string]interface{}:
		for k, v := range typed {
			coerced, err := CoerceInterfaceMapToStringMap(v)
			if err != nil {
				return nil, err
			}
			typed[k] = coerced
		}
		return typed, nil
	case map[interface{}]interface{}:
		m := map[string]interface{}{}
		for k, v := range typed {
			coercedKey, ok := k.(string)
			if !ok {
				return nil, fmt.Errorf("unable to coerce key %T %v to a string", k, k)
			}
			coerced, err := CoerceInterfaceMapToStringMap(v)
			if err != nil {
				return nil, err
			}
			m[coercedKey] = coerced
		}
		return m, nil
	case []interface{}:
		r := make([]interface{}, len(typed))
		for i, v := range typed {
			coerced, err := CoerceInterfaceMapToStringMap(v)
			if err != nil {
				return nil, err
			}
			r[i] = coerced
		}
		return r, nil
	case big.Int:
		value, _ := (in).(big.Int)
		return &value, nil
	default:
		return in, nil
	}
}
