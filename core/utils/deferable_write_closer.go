package utils

import (
	"io"
	"sync"
)

// DeferableWriteCloser is to be used in leiu of defer'ing
// Close on an [io.WriteCloser] (// For more background see https://www.joeshaw.org/dont-defer-close-on-writable-files/)
// Callers should *both*
// explicitly call Close and check for errors when done with the underlying writerclose
// *and* defer the Close() to handle returns before the explicit close
//
// For example rather than
//
//	import "os"
//	f, err := os.Create("./foo")
//	if err != nil { return err}
//	defer f.Close()
//	return f.Write([]bytes("hi"))
//
// do
//
//	import "os"
//	f, err := os.Create("./foo")
//	if err != nil {return nil}
//	wc := NewDeferableWriteCloser(f)
//	defer wc.Close()
//	err = wc.Write([]bytes("hi"))
//	if err != nil {return err}
//	return wc.Close()
type DeferableWriteCloser struct {
	mu       sync.Mutex
	closed   bool
	closeErr error
	io.WriteCloser
}

// NewDeferableWriteCloser creates a deferable writercloser. Callers
// should explicit call and defer Close. See DeferabelWriterCloser for details.
func NewDeferableWriteCloser(wc io.WriteCloser) *DeferableWriteCloser {
	return &DeferableWriteCloser{
		WriteCloser: wc,
	}
}

// Close closes the WriterCloser. The underlying Closer
// is Closed exactly once and resulting error is cached.
// Should be called explicitly AND defered
// Thread safe
func (wc *DeferableWriteCloser) Close() error {

	wc.mu.Lock()
	defer wc.mu.Unlock()
	if !wc.closed {
		wc.closeErr = wc.WriteCloser.Close()
		wc.closed = true
	}
	return wc.closeErr

}
