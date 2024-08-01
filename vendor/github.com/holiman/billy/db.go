// bagdb: Simple datastorage
// Copyright 2021 billy authors
// SPDX-License-Identifier: BSD-3-Clause

package billy

import (
	"fmt"
	"io"
	"sort"
)

// Database represents a `billy` storage.
type Database interface {
	io.Closer

	// Put stores the data to the underlying database, and returns the key needed
	// for later accessing the data.
	// The data is copied by the database, and is safe to modify after the method returns
	Put(data []byte) (uint64, error)

	// Get retrieves the data stored at the given key.
	Get(key uint64) ([]byte, error)

	// Delete marks the data for deletion, which means it will (eventually) be
	// overwritten by other data. After calling Delete with a given key, the results
	// from doing Get(key) is undefined -- it may return the same data, or some other
	// data, or fail with an error.
	Delete(key uint64) error

	// Size returns the storage size of the value belonging to the given key.
	Size(key uint64) uint32

	// Limits returns the smallest and largest slot size.
	Limits() (uint32, uint32)

	// Infos retrieves various internal statistics about the database.
	Infos() *Infos

	// Iterate iterates through all the data in the database, and invokes the
	// given onData method for every element
	Iterate(onData OnDataFn) error
}

// OnDataFn is used to iterate the entire dataset in the database.
// After the method returns, the content of 'data' will be modified by
// the iterator, so it needs to be copied if it is to be used later.
type OnDataFn func(key uint64, size uint32, data []byte)

// SlotSizeFn is a method that acts as a "generator": a closure which, at each
// invocation, should spit out the next slot-size. In order to create a database with three
// shelves invocation of the method should return e.g.
//
//	10, false
//	20, false
//	30, true
//
// OBS! The slot size must take item header size (4 bytes) into account. So if you
// plan to store 120 bytes, then the slot needs to be at least 124 bytes large.
type SlotSizeFn func() (size uint32, done bool)

// SlotSizePowerOfTwo is a SlotSizeFn which arranges the slots in shelves which
// double in size for each level.
func SlotSizePowerOfTwo(min, max uint32) SlotSizeFn {
	v := min
	return func() (uint32, bool) {
		ret := v
		v += v
		return ret, ret >= max
	}
}

// SlotSizeLinear is a SlotSizeFn which arranges the slots in shelves which
// increase linearly.
func SlotSizeLinear(size, count int) SlotSizeFn {
	i := 0
	return func() (uint32, bool) {
		i++
		ret := size * i
		return uint32(ret), i >= count
	}
}

type database struct {
	shelves []*shelf
}

type Options struct {
	Path     string
	Readonly bool
	Snappy   bool // unused for now
}

// Open opens a (new or existing) database, with configurable limits. The given
// slotSizeFn will be used to determine both the shelf sizes and the number of
// shelves. The function must yield values in increasing order.
//
// If shelf already exists, they are opened and read, in order to populate the
// internal gap-list. While doing so, it's a good opportunity for the caller to
// read the data out, (which is probably desirable), which can be done using the
// optional onData callback.
func Open(opts Options, slotSizeFn SlotSizeFn, onData OnDataFn) (Database, error) {
	var (
		db           = &database{}
		prevSlotSize uint32
		prevId       int
		slotSize     uint32
		done         bool
	)
	for !done {
		slotSize, done = slotSizeFn()
		if slotSize <= prevSlotSize {
			return nil, fmt.Errorf("slot sizes must be in increasing order")
		}
		prevSlotSize = slotSize
		shelf, err := openShelf(opts.Path, slotSize, wrapShelfDataFn(len(db.shelves), slotSize, onData), opts.Readonly)
		if err != nil {
			db.Close() // Close shelves
			return nil, err
		}
		db.shelves = append(db.shelves, shelf)

		if id := len(db.shelves) & 0xfff; id < prevId {
			return nil, fmt.Errorf("too many shelves (%d)", len(db.shelves))
		} else {
			prevId = id
		}
	}
	return db, nil
}

// Put stores the data to the underlying database, and returns the key needed
// for later accessing the data.
// The data is copied by the database, and is safe to modify after the method returns
func (db *database) Put(data []byte) (uint64, error) {
	// Search uses binary search to find and return the smallest index i
	// in [0, n) at which f(i) is true,
	index := sort.Search(len(db.shelves), func(i int) bool {
		return len(data)+itemHeaderSize <= int(db.shelves[i].slotSize)
	})
	if index == len(db.shelves) {
		return 0, fmt.Errorf("no shelf found for size %d", len(data))
	}
	if slot, err := db.shelves[index].Put(data); err != nil {
		return 0, err
	} else {
		return slot | uint64(index)<<28, nil
	}
}

// Get retrieves the data stored at the given key.
//
// The key is assumed to be one returned by Put or Iterate (potentially on Open).
// Attempting to access a different key is undefined behavior and may panic.
func (db *database) Get(key uint64) ([]byte, error) {
	id := int(key>>28) & 0xfff
	return db.shelves[id].Get(key & 0x0FFFFFFF)
}

// Delete marks the data for deletion, which means it will (eventually) be
// overwritten by other data. After calling Delete with a given key, the results
// from doing Get(key) is undefined -- it may return the same data, or some other
// data, or fail with an error.
//
// The key is assumed to be one returned by Put or Iterate (potentially on Open).
// Attempting to access a different key is undefined behavior and may panic.
func (db *database) Delete(key uint64) error {
	id := int(key>>28) & 0xfff
	return db.shelves[id].Delete(key & 0x0FFFFFFF)
}

// Size returns the storage size (padding included) of a database entry belonging
// to a key.
//
// The key is assumed to be one returned by Put or Iterate (potentially on Open).
// Attempting to access a different key is undefined behavior and may panic.
func (db *database) Size(key uint64) uint32 {
	id := int(key>>28) & 0xfff
	return db.shelves[id].slotSize
}

func wrapShelfDataFn(shelfId int, shelfSlotSize uint32, onData OnDataFn) onShelfDataFn {
	if onData == nil {
		return nil
	}
	return func(slot uint64, data []byte) {
		key := slot | uint64(shelfId)<<28
		onData(key, shelfSlotSize, data)
	}
}

// Iterate iterates through all the data in the database, and invokes the
// given onData method for every element
func (db *database) Iterate(onData OnDataFn) error {
	var err error
	for i, shelf := range db.shelves {
		if e := shelf.Iterate(wrapShelfDataFn(i, shelf.slotSize, onData)); e != nil {
			err = fmt.Errorf("shelf %d: %w", i, e)
		}
	}
	return err
}

func (db *database) Limits() (uint32, uint32) {
	smallest := db.shelves[0].slotSize
	largest := db.shelves[len(db.shelves)-1].slotSize
	return smallest, largest
}

// Close implements io.Closer
func (db *database) Close() error {
	var err error
	for _, shelf := range db.shelves {
		if e := shelf.Close(); e != nil {
			err = e
		}
	}
	return err
}
