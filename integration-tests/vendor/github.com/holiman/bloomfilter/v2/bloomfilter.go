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
// Copyright © 2020 Martin Holst Swende, continued on the work of Barry Allard

package v2

import (
	"errors"
	"hash"
	"sync"
)

var (
	errHashMismatch = errors.New("hash mismatch, bloom filter corruption or wrong version")
)

// Filter is an opaque Bloom filter type
type Filter struct {
	lock sync.RWMutex
	bits []uint64
	keys []uint64
	m    uint64 // number of bits the "bits" field should recognize
	n    uint64 // number of inserted elements
}

// M is the size of Bloom filter, in bits
func (f *Filter) M() uint64 {
	return f.m
}

// K is the count of keys
func (f *Filter) K() uint64 {
	return uint64(len(f.keys))
}

// Add a hashable item, v, to the filter
func (f *Filter) Add(v hash.Hash64) {
	f.AddHash(v.Sum64())
}

// rotation sets how much to rotate the hash on each filter iteration. This
// is somewhat randomly set to a prime on the lower segment of 64. At 17, the cycle
// does not repeat for quite a while, but even for low number of filters the
// changes are quite rapid
const rotation = 17

// Adds an already hashes item to the filter.
// Identical to Add (but slightly faster)
func (f *Filter) AddHash(hash uint64) {
	f.lock.Lock()
	defer f.lock.Unlock()
	var (
		i uint64
	)
	for n := 0; n < len(f.keys); n++ {
		hash = ((hash << rotation) | (hash >> (64 - rotation))) ^ f.keys[n]
		i = hash % f.m
		f.bits[i>>6] |= 1 << uint(i&0x3f)
	}
	f.n++
}

// ContainsHash tests if f contains the (already hashed) key
// Identical to Contains but slightly faster
func (f *Filter) ContainsHash(hash uint64) bool {
	f.lock.RLock()
	defer f.lock.RUnlock()
	var (
		i uint64
		r = uint64(1)
	)
	for n := 0; n < len(f.keys) && r != 0; n++ {
		hash = ((hash << rotation) | (hash >> (64 - rotation))) ^ f.keys[n]
		i = hash % f.m
		r &= (f.bits[i>>6] >> uint(i&0x3f)) & 1
	}
	return r != 0
}

// Contains tests if f contains v
// false: f definitely does not contain value v
// true:  f maybe contains value v
func (f *Filter) Contains(v hash.Hash64) bool {
	return f.ContainsHash(v.Sum64())
}

// Copy f to a new Bloom filter
func (f *Filter) Copy() (*Filter, error) {
	f.lock.RLock()
	defer f.lock.RUnlock()

	out, err := f.NewCompatible()
	if err != nil {
		return nil, err
	}
	copy(out.bits, f.bits)
	out.n = f.n
	return out, nil
}

// UnionInPlace merges Bloom filter f2 into f
func (f *Filter) UnionInPlace(f2 *Filter) error {
	if !f.IsCompatible(f2) {
		return errors.New("incompatible bloom filters")
	}

	f.lock.Lock()
	defer f.lock.Unlock()

	for i, bitword := range f2.bits {
		f.bits[i] |= bitword
	}
	// Also update the counters
	f.n += f2.n
	return nil
}

// Union merges f2 and f2 into a new Filter out
func (f *Filter) Union(f2 *Filter) (out *Filter, err error) {
	if !f.IsCompatible(f2) {
		return nil, errors.New("incompatible bloom filters")
	}

	f.lock.RLock()
	defer f.lock.RUnlock()

	out, err = f.NewCompatible()
	if err != nil {
		return nil, err
	}
	for i, bitword := range f2.bits {
		out.bits[i] = f.bits[i] | bitword
	}
	// Also update the counters
	out.n = f.n + f2.n
	return out, nil
}
