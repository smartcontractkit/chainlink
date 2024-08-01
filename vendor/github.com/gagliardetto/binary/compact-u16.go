// Copyright 2021 github.com/gagliardetto
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

import (
	"io"
)

// EncodeCompactU16Length encodes a "Compact-u16" length into the provided slice pointer.
// See https://docs.solana.com/developing/programming-model/transactions#compact-u16-format
// See https://github.com/solana-labs/solana/blob/2ef2b6daa05a7cff057e9d3ef95134cee3e4045d/web3.js/src/util/shortvec-encoding.ts
func EncodeCompactU16Length(bytes *[]byte, ln int) {
	rem_len := ln
	for {
		elem := rem_len & 0x7f
		rem_len >>= 7
		if rem_len == 0 {
			*bytes = append(*bytes, byte(elem))
			break
		} else {
			elem |= 0x80
			*bytes = append(*bytes, byte(elem))
		}
	}
}

// DecodeCompactU16Length decodes a "Compact-u16" length from the provided byte slice.
func DecodeCompactU16Length(bytes []byte) int {
	v, _, _ := DecodeCompactU16(bytes)
	return v
}

func DecodeCompactU16(bytes []byte) (int, int, error) {
	ln := 0
	size := 0
	for {
		if len(bytes) == 0 {
			return 0, 0, io.ErrUnexpectedEOF
		}
		elem := int(bytes[0])
		bytes = bytes[1:]
		ln |= (elem & 0x7f) << (size * 7)
		size += 1
		if (elem & 0x80) == 0 {
			break
		}
	}
	return ln, size, nil
}

// DecodeCompactU16LengthFromByteReader decodes a "Compact-u16" length from the provided io.ByteReader.
func DecodeCompactU16LengthFromByteReader(reader io.ByteReader) (int, error) {
	ln := 0
	size := 0
	for {
		elemByte, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}
		elem := int(elemByte)
		ln |= (elem & 0x7f) << (size * 7)
		size += 1
		if (elem & 0x80) == 0 {
			break
		}
	}
	return ln, nil
}
