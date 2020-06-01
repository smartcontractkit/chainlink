package models

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignature(t *testing.T) {
	str := "0xb7a987222fc36c4c8ed1b91264867a422769998aadbeeb1c697586a04fa2b616025b5ca936ec5bdb150999e298b6ecf09251d3c4dd1306dedec0692e7037584800"
	signature, err := NewSignature(str)
	assert.NoError(t, err)

	assert.Equal(t, []byte{
		0xb7, 0xa9, 0x87, 0x22, 0x2f, 0xc3, 0x6c, 0x4c, 0x8e, 0xd1, 0xb9, 0x12,
		0x64, 0x86, 0x7a, 0x42, 0x27, 0x69, 0x99, 0x8a, 0xad, 0xbe, 0xeb, 0x1c,
		0x69, 0x75, 0x86, 0xa0, 0x4f, 0xa2, 0xb6, 0x16, 0x02, 0x5b, 0x5c, 0xa9,
		0x36, 0xec, 0x5b, 0xdb, 0x15, 0x09, 0x99, 0xe2, 0x98, 0xb6, 0xec, 0xf0,
		0x92, 0x51, 0xd3, 0xc4, 0xdd, 0x13, 0x06, 0xde, 0xde, 0xc0, 0x69, 0x2e,
		0x70, 0x37, 0x58, 0x48, 0x00,
	}, signature.Bytes())

	bi, _ := (new(big.Int)).SetString("b7a987222fc36c4c8ed1b91264867a422769998aadbeeb1c697586a04fa2b616025b5ca936ec5bdb150999e298b6ecf09251d3c4dd1306dedec0692e7037584800", 16)
	assert.Equal(t, bi, signature.Big())

	assert.Equal(t, str, signature.String())

	assert.Equal(t, str, signature.String())

	zerosignature := Signature{}
	err = json.Unmarshal([]byte(`"0xb7a987222fc36c4c8ed1b91264867a422769998aadbeeb1c697586a04fa2b616025b5ca936ec5bdb150999e298b6ecf09251d3c4dd1306dedec0692e7037584800"`), &zerosignature)
	assert.NoError(t, err)
	assert.Equal(t, str, zerosignature.String())

	zerosignature = Signature{}
	err = zerosignature.UnmarshalText([]byte(str))
	assert.NoError(t, err)
	assert.Equal(t, str, zerosignature.String())
}
