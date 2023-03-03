// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package cb58

import (
	"bytes"
	"errors"
	"math"

	"github.com/mr-tron/base58/base58"

	"github.com/ava-labs/avalanchego/utils/hashing"
)

const (
	checksumLen = 4
)

var (
	errEncodingOverFlow = errors.New("encoding overflow")
	errMissingChecksum  = errors.New("input string is smaller than the checksum size")
	errBadChecksum      = errors.New("invalid input checksum")
)

// Encode [bytes] to a string using cb58 format.
// [bytes] may be nil, in which case it will be treated the same as an empty
// slice.
func Encode(bytes []byte) (string, error) {
	bytesLen := len(bytes)
	if bytesLen > math.MaxInt32-checksumLen {
		return "", errEncodingOverFlow
	}
	checked := make([]byte, bytesLen+checksumLen)
	copy(checked, bytes)
	copy(checked[len(bytes):], hashing.Checksum(bytes, checksumLen))
	return base58.Encode(checked), nil
}

// Decode [str] to bytes from cb58.
func Decode(str string) ([]byte, error) {
	decodedBytes, err := base58.Decode(str)
	if err != nil {
		return nil, err
	}
	if len(decodedBytes) < checksumLen {
		return nil, errMissingChecksum
	}
	// Verify the checksum
	rawBytes := decodedBytes[:len(decodedBytes)-checksumLen]
	checksum := decodedBytes[len(decodedBytes)-checksumLen:]
	if !bytes.Equal(checksum, hashing.Checksum(rawBytes, checksumLen)) {
		return nil, errBadChecksum
	}
	return rawBytes, nil
}
