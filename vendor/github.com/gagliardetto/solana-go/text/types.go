// Copyright 2020 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package text

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
)

type SafeString string

func (ss SafeString) TextEncode(encoder *Encoder, option *Option) error {
	return encoder.ToWriter(string(ss), option.indent, option.fgColor)
}

type Bool bool

func (b Bool) TextEncode(encoder *Encoder, option *Option) error {
	return encoder.ToWriter(fmt.Sprintf("%t", bool(b)), option.indent, option.fgColor)
}

type HexBytes []byte

func (o HexBytes) TextEncode(encoder *Encoder, option *Option) error {
	return encoder.ToWriter(hex.EncodeToString(o), option.indent, option.fgColor)
}

type Varint16 int16

func (o Varint16) TextEncode(encoder *Encoder, option *Option) error {
	return encoder.ToWriter(fmt.Sprintf("%d", int(o)), option.indent, option.fgColor)
}

type Varuint16 uint16

func (o Varuint16) TextEncode(encoder *Encoder, option *Option) error {
	return encoder.ToWriter(fmt.Sprintf("%d", int(o)), option.indent, option.fgColor)

}

type Varuint32 uint32

func (o Varuint32) TextEncode(encoder *Encoder, option *Option) error {
	return encoder.ToWriter(fmt.Sprintf("%d", int(o)), option.indent, option.fgColor)
}

type Varint32 int32

func (o Varint32) TextEncode(encoder *Encoder, option *Option) error {
	return encoder.ToWriter(fmt.Sprintf("%d", int(o)), option.indent, option.fgColor)

}

type JSONFloat64 float64

func (f JSONFloat64) TextEncode(encoder *Encoder, option *Option) error {
	return encoder.ToWriter(fmt.Sprintf("%f", float64(f)), option.indent, option.fgColor)
}

type Int64 int64

func (i Int64) TextEncode(encoder *Encoder, option *Option) error {
	return encoder.ToWriter(fmt.Sprintf("%d", int64(i)), option.indent, option.fgColor)
}

type Uint64 uint64

func (i Uint64) TextEncode(encoder *Encoder, option *Option) error {
	return encoder.ToWriter(fmt.Sprintf("%d", uint64(i)), option.indent, option.fgColor)
}

// uint128
type Uint128 struct {
	Lo uint64
	Hi uint64
}

func (i Uint128) BigInt() *big.Int {
	buf := make([]byte, 16)
	binary.BigEndian.PutUint64(buf[:], i.Hi)
	binary.BigEndian.PutUint64(buf[8:], i.Lo)
	value := (&big.Int{}).SetBytes(buf)
	return value
}

func (i Uint128) DecimalString() string {
	return i.BigInt().String()
}

func (i Uint128) TextEncode(encoder *Encoder, option *Option) error {
	return encoder.ToWriter(i.BigInt().String(), option.indent, option.fgColor)
}

// Int128
type Int128 Uint128

func (i Int128) BigInt() *big.Int {
	comp := byte(0x80)
	buf := make([]byte, 16)
	binary.BigEndian.PutUint64(buf[:], i.Hi)
	binary.BigEndian.PutUint64(buf[8:], i.Lo)

	var value *big.Int
	if (buf[0] & comp) == comp {
		buf = twosComplement(buf)
		value = (&big.Int{}).SetBytes(buf)
		value = value.Neg(value)
	} else {
		value = (&big.Int{}).SetBytes(buf)
	}
	return value
}

func (i Int128) DecimalString() string {
	return i.BigInt().String()
}

func (i Int128) TextEncode(encoder *Encoder, option *Option) error {
	return encoder.ToWriter(i.BigInt().String(), option.indent, option.fgColor)
}

type Float128 Uint128

func (f Float128) TextEncode(encoder *Encoder, option *Option) error {
	return encoder.ToWriter(Uint128(f).DecimalString(), option.indent, option.fgColor)
}

// Blob

// Blob is base64 encoded data
// https://github.com/EOSIO/fc/blob/0e74738e938c2fe0f36c5238dbc549665ddaef82/include/fc/variant.hpp#L47
type Blob string

// Data returns decoded base64 data
func (b Blob) Data() ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(b))
}

// String returns the blob as a string
func (b Blob) String() string {
	return string(b)
}

func twosComplement(v []byte) []byte {
	buf := make([]byte, len(v))
	for i, b := range v {
		buf[i] = b ^ byte(0xff)
	}
	one := big.NewInt(1)
	value := (&big.Int{}).SetBytes(buf)
	return value.Add(value, one).Bytes()
}
