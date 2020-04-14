package models

import (
	"bytes"
	"encoding/json"

	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/fxamacker/cbor/v2"
)

// ParseCBOR attempts to coerce the input byte array into valid CBOR
// and then coerces it into a JSON object.
func ParseCBOR(b []byte) (JSON, error) {
	if len(b) == 0 {
		return JSON{}, nil
	}

	var m map[interface{}]interface{}

	if err := cbor.Unmarshal(autoAddMapDelimiters(b), &m); err != nil {
		return JSON{}, err
	}

	coerced, err := utils.CoerceInterfaceMapToStringMap(m)
	if err != nil {
		return JSON{}, err
	}

	jsb, err := json.Marshal(coerced)
	if err != nil {
		return JSON{}, err
	}

	var js JSON
	return js, json.Unmarshal(jsb, &js)
}

// Automatically add missing start map and end map to a CBOR encoded buffer
func autoAddMapDelimiters(b []byte) []byte {
	if len(b) < 2 {
		return b
	}

	if (b[0] >> 5) != 5 {
		var buffer bytes.Buffer
		buffer.Write([]byte{0xbf})
		buffer.Write(b)
		buffer.Write([]byte{0xff})
		return buffer.Bytes()
	}

	return b
}
