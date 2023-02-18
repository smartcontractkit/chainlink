package utils

import (
	"io"
	"sync"
)

// DeferableWriterCloser is to be used in leiu of defer'ing
// Close on an io.WriterClose (// For more background see https://www.joeshaw.org/dont-defer-close-on-writable-files/)
// Callers should *both*
// explicitly call Close and check for errors when done with the underlying writerclose
// *and* defer the Close() to handle returns before the explicit close
//
// For example rather than
// f, err := os.Create(...)
// if err != nil {...}
// defer f.Close()
// f.Write(...)
// ...
// return
//
// do
//
// f, err := os.Create(...)
// if err != nil {...}
// wc := NewDeferableWriterCloser(f)
// defer wc.Close()
// wc.Write(...)
// ...
// err = wc.Close()
// if err != nil {...}
// return
type DeferableWriterCloser struct {
	mu sync.Mutex
	io.WriteCloser
}

// NewDeferableWriterCloser creates a deferable writercloser. Callers
// should explicit call and defer Close. See DeferabelWriterCloser for details.
func NewDeferableWriterCloser(wc io.WriteCloser) *DeferableWriterCloser {
	return &DeferableWriterCloser{
		WriteCloser: wc,
	}
}

// Close closes the WriterCloser.
// Should be called explicitly AND defered
// Thread safe
func (wc *DeferableWriterCloser) Close() error {
	var err error
	wc.mu.Lock()
	defer wc.mu.Unlock()
	if wc.WriteCloser != nil {
		err = wc.WriteCloser.Close()
		wc.WriteCloser = nil
	}
	return err
}
