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
	"io"
)

// headerMagic is used to disambiguate between this package and the original
// steakknife implementation.
// Since the key hashing algorithm has changed, the format is no longer
// binary compatible
var version = []byte("v02\n")
var headerMagic = append([]byte{0, 0, 0, 0, 0, 0, 0, 0}, version...)

// counter is a utility to count bytes written
type counter struct {
	bytes int
}

func (c *counter) Write(p []byte) (n int, err error) {
	count := len(p)
	c.bytes += count
	return count, nil
}

// conforms to encoding.BinaryMarshaler

// MarshallToWriter marshalls the filter into the given io.Writer
// Binary layout (Little Endian):
//
//	 k	1 uint64
//	 n	1 uint64
//	 m	1 uint64
//	 keys	[k]uint64
//	 bits	[(m+63)/64]uint64
//	 hash	sha384 (384 bits == 48 bytes)
//
//	 size = (3 + k + (m+63)/64) * 8 bytes
//
func (f *Filter) MarshallToWriter(out io.Writer) (int, [sha512.Size384]byte, error) {
	var (
		c      = &counter{0}
		hasher = sha512.New384()
		mw     = io.MultiWriter(out, hasher, c)
		hash   [sha512.Size384]byte
	)
	f.lock.RLock()
	defer f.lock.RUnlock()

	if _, err := mw.Write(headerMagic); err != nil {
		return c.bytes, hash, err
	}
	if err := binary.Write(mw, binary.LittleEndian, []uint64{f.K(), f.n, f.m}); err != nil {
		return c.bytes, hash, err
	}
	if err := binary.Write(mw, binary.LittleEndian, f.keys); err != nil {
		return c.bytes, hash, err
	}
	// Write it in chunks of 5% (but at least 4K). Otherwise, the binary.Write will allocate a
	// same-size slice of bytes, doubling the memory usage
	var chunkSize = len(f.bits) / 20
	if chunkSize < 512 {
		chunkSize = 512 // Min 4K bytes (512 uint64s)
	}
	buf := make([]byte, chunkSize*8)
	for start := 0; start < len(f.bits); {
		end := start + chunkSize
		if end > len(f.bits) {
			end = len(f.bits)
		}
		for i, x := range f.bits[start:end] {
			binary.LittleEndian.PutUint64(buf[8*i:], x)
		}
		if _, err := mw.Write(buf[0 : (end-start)*8]); err != nil {
			return c.bytes, hash, err
		}
		start = end
	}
	// Now we stop using the multiwriter, pick out the hash of what we've
	// written so far, and then write the hash to the output
	hashbytes := hasher.Sum(nil)
	copy(hash[:], hashbytes[:sha512.Size384])
	err := binary.Write(out, binary.LittleEndian, hashbytes)
	return c.bytes + len(hashbytes), hash, err
}

// MarshalBinary converts a Filter into []bytes
func (f *Filter) MarshalBinary() (data []byte, err error) {
	buf := new(bytes.Buffer)
	_, _, err = f.MarshallToWriter(buf)
	if err != nil {
		return nil, err
	}
	data = buf.Bytes()
	return data, nil
}
