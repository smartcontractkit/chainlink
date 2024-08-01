// bagdb: Simple datastorage
// Copyright 2021 billy authors
// SPDX-License-Identifier: BSD-3-Clause

package billy

import (
	"errors"
	"io"
	"os"
	"sync"
)

// memoryStore implements store for in-memory ephemeral data persistence.
type memoryStore struct {
	buffer []byte
	lock   sync.Mutex
}

// ReadAt implements io.ReaderAt of the store interface.
func (ms *memoryStore) ReadAt(p []byte, off int64) (int, error) {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	return ms.readAt(p, off)
}

// readAt is the actual implementation of Read/ReadAt.
func (ms *memoryStore) readAt(p []byte, off int64) (int, error) {
	if off < 0 {
		return 0, errors.New("memoryStore.ReadAt: negative offset")
	}
	if off >= int64(len(ms.buffer)) {
		return 0, io.EOF
	}
	copy(p, ms.buffer[off:])
	read := len(p)
	fail := error(nil)

	if off += int64(read); off > int64(len(ms.buffer)) {
		read -= int(off - int64(len(ms.buffer)))
		fail = io.EOF
	}
	return read, fail
}

// WriteAt implements io.WriterAt of the store interface.
func (ms *memoryStore) WriteAt(p []byte, off int64) (int, error) {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	return ms.writeAt(p, off)
}

// writeAt is the actual implementation of Write/WriteAt.
func (ms *memoryStore) writeAt(p []byte, off int64) (int, error) {
	if off < 0 {
		return 0, errors.New("memoryStore.WriteAt: negative offset")
	}
	if off+int64(len(p)) > int64(len(ms.buffer)) {
		ms.buffer = append(ms.buffer, make([]byte, off+int64(len(p))-int64(len(ms.buffer)))...)
	}
	copy(ms.buffer[off:], p)
	return len(p), nil
}

// Close implements io.Closer of the store interface.
func (ms *memoryStore) Close() error {
	return nil
}

// Stat implements the store interface.
func (ms *memoryStore) Stat() (os.FileInfo, error) {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	return &fileinfoMock{size: int64(len(ms.buffer))}, nil
}

// Truncate implements the store interface.
func (ms *memoryStore) Truncate(size int64) error {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	if int64(len(ms.buffer)) >= size {
		ms.buffer = ms.buffer[:size]
	} else {
		ms.buffer = append(ms.buffer, make([]byte, size-int64(len(ms.buffer)))...)
	}
	return nil
}

// Sync implements the store interface.
func (ms *memoryStore) Sync() error {
	return nil
}
