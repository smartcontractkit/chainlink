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
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math"
)

const (
	MMin        = 2 // MMin is the minimum Bloom filter bits count
	KMin        = 1 // KMin is the minimum number of keys
	Uint64Bytes = 8 // Uint64Bytes is the number of bytes in type uint64
)

// OptimalK calculates the optimal k value for creating a new Bloom filter
// maxn is the maximum anticipated number of elements
func OptimalK(m, maxN uint64) uint64 {
	return uint64(math.Ceil(float64(m) * math.Ln2 / float64(maxN)))
}

// OptimalM calculates the optimal m value for creating a new Bloom filter
// p is the desired false positive probability
// optimal m = ceiling( - n * ln(p) / ln(2)**2 )
func OptimalM(maxN uint64, p float64) uint64 {
	return uint64(math.Ceil(-float64(maxN) * math.Log(p) / (math.Ln2 * math.Ln2)))
}

// New Filter with CSPRNG keys
//
// m is the size of the Bloom filter, in bits, >= 2
//
// k is the number of random keys, >= 1
func New(m, k uint64) (*Filter, error) {
	return NewWithKeys(m, newRandKeys(m, k))
}

func newRandKeys(m uint64, k uint64) []uint64 {
	keys := make([]uint64, k)
	if err := binary.Read(crand.Reader, binary.LittleEndian, keys); err != nil {
		panic(fmt.Sprintf("Cannot read %d bytes from CSRPNG crypto/rand.Read (err=%v)",
			Uint64Bytes, err))
	}
	return keys
}

// NewCompatible Filter compatible with f
func (f *Filter) NewCompatible() (*Filter, error) {
	return NewWithKeys(f.m, f.keys)
}

// NewOptimal Bloom filter with random CSPRNG keys
func NewOptimal(maxN uint64, p float64) (*Filter, error) {
	m := OptimalM(maxN, p)
	k := OptimalK(m, maxN)
	return New(m, k)
}

// uniqueKeys is true if all keys are unique
func uniqueKeys(keys []uint64) bool {
	for j := 0; j < len(keys)-1; j++ {
		for i := j + 1; i < len(keys); i++ {
			if keys[i] == keys[j] {
				return false
			}
		}
	}
	return true
}

// NewWithKeys creates a new Filter from user-supplied origKeys
func NewWithKeys(m uint64, origKeys []uint64) (f *Filter, err error) {
	var (
		bits []uint64
		keys []uint64
	)
	if bits, err = newBits(m); err != nil {
		return nil, err
	}
	if keys, err = newKeysCopy(origKeys); err != nil {
		return nil, err
	}
	return &Filter{
		m:    m,
		n:    0,
		bits: bits,
		keys: keys,
	}, nil
}

func newBits(m uint64) ([]uint64, error) {
	if m < MMin {
		return nil, fmt.Errorf("number of bits in the filter must be >= %d (was %d)", MMin, m)
	}
	return make([]uint64, (m+63)/64), nil
}

func newKeysCopy(origKeys []uint64) (keys []uint64, err error) {
	if len(origKeys) < KMin {
		return nil, fmt.Errorf("keys must have length %d or greater (was %d)", KMin, len(origKeys))
	}
	if !uniqueKeys(origKeys) {
		return nil, fmt.Errorf("Bloom filter keys must be unique")
	}
	keys = append(keys, origKeys...)
	return keys, err
}
