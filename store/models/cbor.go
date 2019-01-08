package models

import (
	"bytes"
	"encoding/json"

	"github.com/smartcontractkit/chainlink/utils"
	"github.com/ugorji/go/codec"
)

// ParseCBOR attempts to coerce the input byte array into valid CBOR
// and then coerces it into a JSON object.
func ParseCBOR(b []byte) (JSON, error) {
	var m map[interface{}]interface{}

	cbor := codec.NewDecoderBytes(autoAddMapDelimiters(b), new(codec.CborHandle))
	if err := cbor.Decode(&m); err != nil {
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

	var buffer bytes.Buffer
	if (b[0] >> 5) != 5 {
		buffer.Write([]byte{0xbf})
	}
	buffer.Write(b)

	if b[len(b)-1] != 0xff {
		buffer.Write([]byte{0xff})
	}
	return buffer.Bytes()
}
