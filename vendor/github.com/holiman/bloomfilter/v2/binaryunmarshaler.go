// Package bloomfilter is face-meltingly fast, thread-safe,
// marshalable, unionable, probability- and
// optimal-size-calculating Bloom filter in go
//
// https://github.com/steakknife/bloomfilter
//
// Copyright © 2014, 2015, 2018 Barry Allard
//
// MIT license
//
package v2

import (
	"bytes"
	"crypto/sha512"
	"encoding/binary"
	"fmt"
	"hash"
	"io"
)

func unmarshalBinaryHeader(r io.Reader) (k, n, m uint64, err error) {
	magic := make([]byte, len(headerMagic))
	if _, err := io.ReadFull(r, magic); err != nil {
		return 0, 0, 0, err
	}
	if !bytes.Equal(magic, headerMagic) {
		return 0, 0, 0, fmt.Errorf("incompatible version (wrong magic), got %x", magic)
	}
	var knm = make([]uint64, 3)
	err = binary.Read(r, binary.LittleEndian, knm)
	if err != nil {
		return 0, 0, 0, err
	}
	k = knm[0]
	n = knm[1]
	m = knm[2]
	if k < KMin {
		return 0, 0, 0, fmt.Errorf("keys must have length %d or greater (was %d)", KMin, k)
	}
	if m < MMin {
		return 0, 0, 0, fmt.Errorf("number of bits in the filter must be >= %d (was %d)", MMin, m)
	}
	return k, n, m, err
}

func unmarshalBinaryBits(r io.Reader, m uint64) (bits []uint64, err error) {
	bits, err = newBits(m)
	if err != nil {
		return bits, err
	}
	bs := make([]byte, 8)
	for i := 0; i < len(bits) && err == nil; i++ {
		_, err = io.ReadFull(r, bs)
		bits[i] = binary.LittleEndian.Uint64(bs)
	}
	if err != nil {
		return nil, err
	}
	return bits, nil
}

func unmarshalBinaryKeys(r io.Reader, k uint64) (keys []uint64, err error) {
	keys = make([]uint64, k)
	err = binary.Read(r, binary.LittleEndian, keys)
	return keys, err
}

// hashingReader can be used to read from a reader, and simultaneously
// do a hash on the bytes that were read
type hashingReader struct {
	reader io.Reader
	hasher hash.Hash
	tot    int64
}

func (h *hashingReader) Read(p []byte) (n int, err error) {
	n, err = h.reader.Read(p)
	h.tot += int64(n)
	if err != nil {
		return n, err
	}
	_, _ = h.hasher.Write(p[:n])
	return n, err
}

// UnmarshalBinary converts []bytes into a Filter
// conforms to encoding.BinaryUnmarshaler
func (f *Filter) UnmarshalBinary(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	_, err = f.UnmarshalFromReader(buf)
	return err
}

func (f *Filter) UnmarshalFromReader(input io.Reader) (n int64, err error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	buf := &hashingReader{
		reader: input,
		hasher: sha512.New384(),
	}
	var k uint64
	k, f.n, f.m, err = unmarshalBinaryHeader(buf)
	if err != nil {
		return buf.tot, err
	}

	f.keys, err = unmarshalBinaryKeys(buf, k)
	if err != nil {
		return buf.tot, err
	}
	f.bits, err = unmarshalBinaryBits(buf, f.m)
	if err != nil {
		return buf.tot, err
	}

	// Only the hash remains to be read now
	// so abort the hasher at this point
	gotHash := buf.hasher.Sum(nil)
	expHash := make([]byte, sha512.Size384)
	err = binary.Read(buf, binary.LittleEndian, expHash)
	if err != nil {
		return buf.tot, err
	}
	if !bytes.Equal(gotHash, expHash) {
		return buf.tot, errHashMismatch
	}
	return buf.tot, nil
}
