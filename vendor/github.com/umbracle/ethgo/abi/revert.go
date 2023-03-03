package abi

import (
	"bytes"
	"fmt"
)

var revertId = []byte{0x8, 0xC3, 0x79, 0xA0}

func UnpackRevertError(b []byte) (string, error) {
	if !bytes.HasPrefix(b, revertId) {
		return "", fmt.Errorf("revert error prefix not found")
	}

	b = b[4:]
	tt := MustNewType("tuple(string)")
	vals, err := tt.Decode(b)
	if err != nil {
		return "", err
	}
	revVal := vals.(map[string]interface{})["0"].(string)
	return revVal, nil
}
