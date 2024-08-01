// bagdb: Simple datastorage
// Copyright 2021 billy authors
// SPDX-License-Identifier: BSD-3-Clause

package billy

import (
	"io"
	"os"
	"time"
)

// store is an interface mimicking an os.File to allow using different storage
// backends for different database instances.
type store interface {
	io.ReaderAt
	io.WriterAt
	io.Closer

	Stat() (os.FileInfo, error)
	Truncate(size int64) error
	Sync() error
}

// fileinfoMock is a mock implementation for returning non-single-file store
// sizes.
type fileinfoMock struct {
	size int64
}

func (s *fileinfoMock) Size() int64 {
	return s.size
}

func (s *fileinfoMock) Name() string       { panic("not implemented") }
func (s *fileinfoMock) Mode() os.FileMode  { panic("not implemented") }
func (s *fileinfoMock) ModTime() time.Time { panic("not implemented") }
func (s *fileinfoMock) IsDir() bool        { panic("not implemented") }
func (s *fileinfoMock) Sys() any           { panic("not implemented") }
