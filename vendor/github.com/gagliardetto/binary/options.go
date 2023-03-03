// Copyright 2021 github.com/gagliardetto
// This file has been modified by github.com/gagliardetto
//
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

package bin

import "encoding/binary"

type option struct {
	OptionalField bool
	SizeOfSlice   *int
	Order         binary.ByteOrder
}

var LE binary.ByteOrder = binary.LittleEndian
var BE binary.ByteOrder = binary.BigEndian

var defaultByteOrder = binary.LittleEndian

func newDefaultOption() *option {
	return &option{
		OptionalField: false,
		Order:         defaultByteOrder,
	}
}

func (o *option) clone() *option {
	out := &option{
		OptionalField: o.OptionalField,
		SizeOfSlice:   o.SizeOfSlice,
		Order:         o.Order,
	}
	return out
}

func (o *option) isOptional() bool {
	return o.OptionalField
}

func (o *option) hasSizeOfSlice() bool {
	return o.SizeOfSlice != nil
}

func (o *option) getSizeOfSlice() int {
	return *o.SizeOfSlice
}

func (o *option) setSizeOfSlice(size int) *option {
	o.SizeOfSlice = &size
	return o
}
func (o *option) setIsOptional(isOptional bool) *option {
	o.OptionalField = isOptional
	return o
}

type Encoding int

const (
	EncodingBin Encoding = iota
	EncodingCompactU16
	EncodingBorsh
)

func (enc Encoding) String() string {
	switch enc {
	case EncodingBin:
		return "Bin"
	case EncodingCompactU16:
		return "CompactU16"
	case EncodingBorsh:
		return "Borsh"
	default:
		return ""
	}
}

func (en Encoding) IsBorsh() bool {
	return en == EncodingBorsh
}

func (en Encoding) IsBin() bool {
	return en == EncodingBin
}

func (en Encoding) IsCompactU16() bool {
	return en == EncodingCompactU16
}

func isValidEncoding(enc Encoding) bool {
	switch enc {
	case EncodingBin, EncodingCompactU16, EncodingBorsh:
		return true
	default:
		return false
	}
}
