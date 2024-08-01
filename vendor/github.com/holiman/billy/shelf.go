// bagdb: Simple datastorage
// Copyright 2021 billy authors
// SPDX-License-Identifier: BSD-3-Clause

package billy

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

const (
	curVersion     = uint16(0)
	itemHeaderSize = 4 // size of the per-item header
	maxSlotSize    = uint64(0xffffffff)
	// minSlotSize is the minimum size of a slot. It needs to fit the header,
	// and then some actual data too.
	minSlotSize = itemHeaderSize * 2
)

var (
	ErrClosed      = errors.New("shelf closed")
	ErrOversized   = errors.New("data too large for shelf")
	ErrBadIndex    = errors.New("bad index")
	ErrEmptyData   = errors.New("empty data")
	ErrReadonly    = errors.New("read-only mode")
	ErrCorruptData = errors.New("corrupt data")
)

// shelf represents a collection of similarly-sized items. The shelf uses
// a number of slots, where each slot is of the exact same size.
type shelf struct {
	slotSize uint32 // Size of the slots, up to 4GB

	// gaps is a slice of indices to slots that are free to use. The
	// gaps are always sorted lowest numbers first.
	gaps   sortedUniqueInts
	gapsMu sync.Mutex // Mutex for operating on 'gaps' and 'count'.
	count  uint64     // count holds the number of items on the shelf.

	f      store        // f is the file where data is persisted.
	fileMu sync.RWMutex // Mutex for file operations on 'f' (rw versus Close) and closed.

	closed   bool
	readonly bool
}

var (
	Magic           = [5]byte{'b', 'i', 'l', 'l', 'y'}
	ShelfHeaderSize = binary.Size(shelfHeader{})
)

// shelfHeader is the file-header for a shelf file. It has a 'magic' "billy" prefix,
// followed by version and slotsize.
type shelfHeader struct {
	Magic    [5]byte // "billy'
	Version  uint16
	Slotsize uint32
}

// openShelf opens a (new or existing) shelf with the given slot size.
// If the shelf already exists, it's opened and read, which populates the
// internal gap-list.
// The onData callback is optional, and can be nil.
func openShelf(path string, slotSize uint32, onData onShelfDataFn, readonly bool) (*shelf, error) {
	if slotSize < minSlotSize {
		return nil, fmt.Errorf("slot size %d smaller than minimum (%d)", slotSize, minSlotSize)
	}
	if path != "" { // empty path == in-memory database
		if finfo, err := os.Stat(path); err != nil {
			return nil, err
		} else if !finfo.IsDir() {
			return nil, fmt.Errorf("not a directory: '%v'", path)
		}
	}
	var (
		fileSize int
		h        = shelfHeader{Magic, curVersion, slotSize}
		fname    = fmt.Sprintf("bkt_%08d.bag", slotSize)
		flags    = os.O_RDWR | os.O_CREATE
	)
	if readonly {
		flags = os.O_RDONLY
	}
	var (
		f   store
		err error
	)
	if path != "" {
		f, err = os.OpenFile(filepath.Join(path, fname), flags, 0666)
		if err != nil {
			return nil, err
		}
	} else {
		f = new(memoryStore)
	}
	if stat, err := f.Stat(); err != nil {
		_ = f.Close()
		return nil, err
	} else {
		fileSize = int(stat.Size())
	}
	if fileSize == 0 {
		a := new(bytes.Buffer)
		if err = binary.Write(a, binary.BigEndian, &h); err == nil {
			_, err = f.WriteAt(a.Bytes(), 0)
		}
	} else {
		b := make([]byte, binary.Size(h))
		if _, err = f.ReadAt(b, 0); err == nil {
			err = binary.Read(bytes.NewReader(b), binary.BigEndian, &h)
		}
	}
	if err != nil {
		_ = f.Close()
		return nil, err
	}
	switch {
	case h.Magic != Magic:
		err = errors.New("missing magic")
	case h.Version != curVersion:
		err = fmt.Errorf("wrong version: %d", h.Version)
	case h.Slotsize != slotSize:
		err = fmt.Errorf("wrong slotsize, file:%d, need:%d", h.Slotsize, slotSize)
	}
	if err != nil {
		_ = f.Close()
		return nil, err
	}
	dataSize := fileSize
	if fileSize >= int(ShelfHeaderSize) {
		dataSize = fileSize - int(ShelfHeaderSize)
	}
	sh := &shelf{
		slotSize: slotSize,
		count:    uint64((dataSize + int(slotSize) - 1) / int(slotSize)),
		f:        f,
		readonly: readonly,
	}
	// Compact + iterate
	if err := sh.compact(onData); err != nil {
		_ = f.Close()
		return nil, err
	}
	return sh, nil
}

func (s *shelf) Close() error {
	// We don't need the gapsMu until later, but order matters: all places
	// which require both mutexes first obtain gapsMu, and _then_ fileMu.
	// If one place uses a different order, then a deadlock is possible
	s.gapsMu.Lock()
	defer s.gapsMu.Unlock()
	s.fileMu.Lock()
	defer s.fileMu.Unlock()
	if s.closed {
		return nil
	}
	s.closed = true
	if s.readonly {
		return nil
	}
	var err error
	setErr := func(e error) {
		if err == nil && e != nil {
			err = e
		}
	}
	// Before closing the file, we overwrite all gaps with
	// blank space in the headers. Later on, when opening, we can reconstruct the
	// gaps by skimming through the slots and checking the headers.
	hdr := make([]byte, 4)
	for _, gap := range s.gaps {
		setErr(s.writeSlot(hdr, gap))
	}
	s.gaps = s.gaps[:0]
	setErr(s.f.Sync())
	setErr(s.f.Close())
	return err
}

// Update overwrites the existing data at the given slot. This operation is more
// efficient than Delete + Put, since it does not require managing slot availability
// but instead just overwrites in-place.
func (s *shelf) Update(data []byte, slot uint64) error {
	if s.readonly {
		return ErrReadonly
	}
	if len(data) == 0 {
		return ErrEmptyData
	}
	if have, max := uint32(len(data)+itemHeaderSize), s.slotSize; have > max {
		return ErrOversized
	}
	return s.update(data, slot)
}

// Put writes the given data and returns a slot identifier. The caller may
// modify the data after this method returns.
func (s *shelf) Put(data []byte) (uint64, error) {
	if s.readonly {
		return 0, ErrReadonly
	}
	if len(data) == 0 {
		return 0, ErrEmptyData
	}
	if have, max := uint32(len(data)+itemHeaderSize), s.slotSize; have > max {
		return 0, ErrOversized
	}
	slot := s.getSlot()
	return slot, s.update(data, slot)
}

// update writes the data to the given slot.
func (s *shelf) update(data []byte, slot uint64) error {
	// Read-lock to prevent file from being closed while writing to it
	s.fileMu.RLock()
	defer s.fileMu.RUnlock()
	if s.closed {
		return ErrClosed
	}
	buf := make([]byte, s.slotSize)
	binary.BigEndian.PutUint32(buf, uint32(len(data))) // Write header
	copy(buf[itemHeaderSize:], data)                   // Write data
	return s.writeSlot(buf, slot)
}

// Delete marks the data at the given slot of deletion.
// Delete does not touch the disk. When the shelf is Close():d, any remaining
// gaps will be marked as such in the backing file.
// NOTE: If a Get-operation is performed _after_ Delete, then the results
// are undefined. It may return the original value or a new value, if a new
// value has been written into the slot.
// It will _not_ return any kind of "MissingItem" error in this scenario.
func (s *shelf) Delete(slot uint64) error {
	if s.readonly {
		return ErrReadonly
	}
	// Mark gap
	s.gapsMu.Lock()
	defer s.gapsMu.Unlock()
	// Can't delete outside of the file
	if slot >= s.count {
		return fmt.Errorf("%w: shelf %d, slot %d, tail %d", ErrBadIndex, s.slotSize, slot, s.count)
	}
	// We try to keep writes going to the early parts of the file, to have the
	// possibility of trimming the file when/if the tail becomes unused.
	s.gaps.Append(slot)

	// s.count is the first empty location. If the gaps has reached to one below
	// the tail, then we can start truncating
	if lastGap := s.gaps[len(s.gaps)-1]; lastGap+1 == s.count {
		// we can delete a portion of the file
		s.fileMu.Lock()
		defer s.fileMu.Unlock()
		if s.closed { // Undo (not really important, but correct) and back out again
			s.gaps = s.gaps[:0]
			return ErrClosed
		}
		for len(s.gaps) > 0 && s.gaps[len(s.gaps)-1]+1 == s.count {
			s.gaps = s.gaps[:len(s.gaps)-1]
			s.count--
		}
		if err := s.f.Truncate(int64(ShelfHeaderSize) + int64(s.count*uint64(s.slotSize))); err != nil {
			return err
		}
	}
	return nil
}

// Get returns the data at the given slot. If the slot has been deleted, the returndata
// this method is undefined: it may return the original data, or some newer data
// which has been written into the slot after Delete was called.
func (s *shelf) Get(slot uint64) ([]byte, error) {
	// Read-lock to prevent file from being closed while reading from it
	s.fileMu.RLock()
	defer s.fileMu.RUnlock()
	if s.closed {
		return nil, ErrClosed
	}
	data, err := s.readSlot(make([]byte, s.slotSize), slot)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBadIndex, err)
	}
	return data, nil
}

// readSlot is a convenience function to the data from a slot.
// It
//   - expects the given 'buf' to be correctly sized (len = s.slotSize),
//   - expects the fileMu to be R-locked, to prevent file from being closed during/before reading.
//   - returns a subslice of buf containing the live data.
func (s *shelf) readSlot(buf []byte, slot uint64) ([]byte, error) {
	// Read the entire slot at once -- this might mean we read a bit more
	// than strictly necessary, but it saves us one syscall.
	if _, err := s.f.ReadAt(buf, int64(ShelfHeaderSize)+int64(slot)*int64(s.slotSize)); err != nil {
		return nil, err
	}
	size := binary.BigEndian.Uint32(buf) + itemHeaderSize
	if size > uint32(s.slotSize) {
		return nil, fmt.Errorf("%w: item size %d, slot size %d", ErrCorruptData, size, s.slotSize)
	}
	return buf[itemHeaderSize:size], nil
}

// writeSlot writes the given data to the slot. This method assumes that the
// fileMu is read-locked.
func (s *shelf) writeSlot(data []byte, slot uint64) error {
	_, err := s.f.WriteAt(data, int64(ShelfHeaderSize)+int64(slot)*int64(s.slotSize))
	return err
}

func (s *shelf) getSlot() uint64 {
	var slot uint64
	// Locate the first free slot
	s.gapsMu.Lock()
	defer s.gapsMu.Unlock()
	if nGaps := len(s.gaps); nGaps > 0 {
		slot = s.gaps[0]
		s.gaps = s.gaps[1:]
		return slot
	}
	// No gaps available: Expand the tail
	slot = s.count
	s.count++
	return slot
}

// onShelfDataFn is used to iterate the entire dataset in the shelf.
// After the method returns, the content of 'data' will be modified by
// the iterator, so it needs to be copied if it is to be used later.
type onShelfDataFn func(slot uint64, data []byte)

// Iterate iterates through the elements on the shelf, and invokes the onData
// callback for each item.
func (s *shelf) Iterate(onData onShelfDataFn) error {
	s.gapsMu.Lock()
	defer s.gapsMu.Unlock()

	s.fileMu.RLock()
	defer s.fileMu.RUnlock()
	if s.closed {
		return ErrClosed
	}

	var (
		buf     = make([]byte, s.slotSize)
		nextGap = uint64(0xffffffffffffffff)
		gapIdx  = 0
	)
	if gapIdx < len(s.gaps) {
		nextGap = s.gaps[gapIdx]
	}
	for slot := uint64(0); slot < s.count; slot++ {
		if slot == nextGap {
			// We've reached a gap. Skip it
			gapIdx++
			if gapIdx < len(s.gaps) {
				nextGap = s.gaps[gapIdx]
			}
			// implicit else: leave 'nextGap' as is, we're already past it now
			// and won't hit this clause again
			continue
		}
		data, err := s.readSlot(buf, slot)
		if err != nil {
			return err
		}
		onData(slot, data)
	}
	return nil
}

// compact moves data 'up' to fill gaps, and truncates the file afterwards.
// This operation must only be performed during the opening of the shelf.
func (s *shelf) compact(onData onShelfDataFn) error {
	s.gapsMu.Lock()
	defer s.gapsMu.Unlock()
	s.fileMu.RLock()
	defer s.fileMu.RUnlock()

	buf := make([]byte, s.slotSize)
	// nextGap searches upwards from the given slot (inclusive),
	// to find the first gap.
	nextGap := func(slot uint64) (uint64, error) {
		for ; slot < s.count; slot++ {
			data, err := s.readSlot(buf, slot)
			if err != nil {
				return 0, err
			}
			if len(data) == 0 { // We've found a gap
				break
			}
			if onData != nil {
				onData(slot, data)
			}
		}
		return slot, nil
	}
	// prevData searches downwards from the given slot (inclusive), to find
	// the next data-filled slot.
	prevData := func(slot, gap uint64) (uint64, error) {
		for ; slot > gap && slot > 0; slot-- {
			data, err := s.readSlot(buf, slot)
			if err != nil {
				return 0, err
			}
			if len(data) != 0 {
				// We've found a slot of data. Copy it to the gap
				if err := s.writeSlot(buf, gap); err != nil {
					return 0, err
				}
				if onData != nil {
					onData(gap, data)
				}
				break
			}
		}
		return slot, nil
	}
	var (
		gapped = uint64(0)
		filled = s.count
		empty  = s.count == 0
		err    error
	)
	// The compaction / iteration goes through the file two directions:
	// - forwards: search for gaps,
	// - backwards: searh for data to move into the gaps
	// The two searches happen in turns, and if both find a match, the
	// data is moved from the slot to the gap. Once the two searches cross eachother,
	// the algorithm is finished.
	// This algorithm reads minimal number of items and performs minimal
	// number of writes.
	s.gaps = make([]uint64, 0)
	if empty {
		return nil
	}
	if s.readonly {
		// Don't (try to) mutate the file in readonly mode, but still
		// iterate for the ondata callbacks.
		for gapped <= s.count {
			gapped, err = nextGap(gapped)
			if err != nil {
				return err
			}
			gapped++
		}
		return nil
	}
	filled--
	firstTail := s.count
	for gapped <= filled {
		// Find next gap. If we've reached the tail, we're done here.
		if gapped, err = nextGap(gapped); err != nil {
			return err
		}
		if gapped >= s.count {
			break
		}
		// We have a gap. Now, find the last piece of data to move to that gap
		if filled, err = prevData(filled, gapped); err != nil {
			return err
		}
		// dataSlot is now the empty area
		s.count = filled
		gapped++
		filled--
	}
	if firstTail != s.count {
		// Some gc was performed. gapSlot is the first empty slot now
		if err := s.f.Truncate(int64(ShelfHeaderSize) + int64(s.count*uint64(s.slotSize))); err != nil {
			return fmt.Errorf("truncation failed: %v", err)
		}
	}
	return nil
}

// stats returns the total number of slots in the shelf and the gaps within.
func (s *shelf) stats() (uint64, uint64) {
	s.gapsMu.Lock()
	defer s.gapsMu.Unlock()

	return s.count, uint64(len(s.gaps))
}

// sortedUniqueInts is a helper structure to maintain an ordered slice
// of gaps. We keep them ordered to make writes prefer early slots, to increase
// the chance of trimming the end of files upon deletion.
type sortedUniqueInts []uint64

func (u *sortedUniqueInts) Append(elem uint64) {
	s := *u
	size := len(s)
	idx := sort.Search(size, func(i int) bool {
		return elem <= s[i]
	})
	if idx < size && s[idx] == elem {
		return // Elem already there
	}
	*u = append(s[:idx], append([]uint64{elem}, s[idx:]...)...)
}
