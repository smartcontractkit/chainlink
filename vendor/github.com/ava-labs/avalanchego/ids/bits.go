// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package ids

import (
	"bytes"
	"math/bits"
)

// NumBits is the number of bits this patricia tree manages
const NumBits = 256

// BitsPerByte is the number of bits per byte
const BitsPerByte = 8

// EqualSubset takes in two indices and two ids and returns if the ids are
// equal from bit start to bit end (non-inclusive). Bit indices are defined as:
// [7 6 5 4 3 2 1 0] [15 14 13 12 11 10 9 8] ... [255 254 253 252 251 250 249 248]
// Where index 7 is the MSB of byte 0.
func EqualSubset(start, stop int, id1, id2 ID) bool {
	stop--
	if start > stop || stop < 0 {
		return true
	}
	if stop >= NumBits {
		return false
	}

	startIndex := start / BitsPerByte
	stopIndex := stop / BitsPerByte

	// If there is a series of bytes between the first byte and the last byte, they must be equal
	if startIndex+1 < stopIndex && !bytes.Equal(id1[startIndex+1:stopIndex], id2[startIndex+1:stopIndex]) {
		return false
	}

	startBit := uint(start % BitsPerByte) // Index in the byte that the first bit is at
	stopBit := uint(stop % BitsPerByte)   // Index in the byte that the last bit is at

	startMask := -1 << startBit          // 111...0... The number of 0s is equal to startBit
	stopMask := (1 << (stopBit + 1)) - 1 // 000...1... The number of 1s is equal to stopBit+1

	if startIndex == stopIndex {
		// If we are looking at the same byte, both masks need to be applied
		mask := startMask & stopMask

		// The index here could be startIndex or stopIndex, as they are equal
		b1 := mask & int(id1[startIndex])
		b2 := mask & int(id2[startIndex])

		return b1 == b2
	}

	start1 := startMask & int(id1[startIndex])
	start2 := startMask & int(id2[startIndex])

	stop1 := stopMask & int(id1[stopIndex])
	stop2 := stopMask & int(id2[stopIndex])

	return start1 == start2 && stop1 == stop2
}

// FirstDifferenceSubset takes in two indices and two ids and returns the index
// of the first difference between the ids inside bit start to bit end
// (non-inclusive). Bit indices are defined above
func FirstDifferenceSubset(start, stop int, id1, id2 ID) (int, bool) {
	stop--
	if start > stop || stop < 0 || stop >= NumBits {
		return 0, false
	}

	startIndex := start / BitsPerByte
	stopIndex := stop / BitsPerByte

	startBit := uint(start % BitsPerByte) // Index in the byte that the first bit is at
	stopBit := uint(stop % BitsPerByte)   // Index in the byte that the last bit is at

	startMask := -1 << startBit          // 111...0... The number of 0s is equal to startBit
	stopMask := (1 << (stopBit + 1)) - 1 // 000...1... The number of 1s is equal to stopBit+1

	if startIndex == stopIndex {
		// If we are looking at the same byte, both masks need to be applied
		mask := startMask & stopMask

		// The index here could be startIndex or stopIndex, as they are equal
		b1 := mask & int(id1[startIndex])
		b2 := mask & int(id2[startIndex])

		if b1 == b2 {
			return 0, false
		}

		bd := b1 ^ b2
		return bits.TrailingZeros8(uint8(bd)) + startIndex*BitsPerByte, true
	}

	// Check the first byte, may have some bits masked
	start1 := startMask & int(id1[startIndex])
	start2 := startMask & int(id2[startIndex])

	if start1 != start2 {
		bd := start1 ^ start2
		return bits.TrailingZeros8(uint8(bd)) + startIndex*BitsPerByte, true
	}

	// Check all the interior bits
	for i := startIndex + 1; i < stopIndex; i++ {
		b1 := int(id1[i])
		b2 := int(id2[i])
		if b1 != b2 {
			bd := b1 ^ b2
			return bits.TrailingZeros8(uint8(bd)) + i*BitsPerByte, true
		}
	}

	// Check the last byte, may have some bits masked
	stop1 := stopMask & int(id1[stopIndex])
	stop2 := stopMask & int(id2[stopIndex])

	if stop1 != stop2 {
		bd := stop1 ^ stop2
		return bits.TrailingZeros8(uint8(bd)) + stopIndex*BitsPerByte, true
	}

	// No difference was found
	return 0, false
}
