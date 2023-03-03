package multiplex

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"

	pool "github.com/libp2p/go-buffer-pool"
	"go.uber.org/multierr"
)

var (
	ErrStreamReset  = errors.New("stream reset")
	ErrStreamClosed = errors.New("closed stream")
)

// streamID is a convenience type for operating on stream IDs
type streamID struct {
	id        uint64
	initiator bool
}

// header computes the header for the given tag
func (id *streamID) header(tag uint64) uint64 {
	header := id.id<<3 | tag
	if !id.initiator {
		header--
	}
	return header
}

type Stream struct {
	id     streamID
	name   string
	dataIn chan []byte
	mp     *Multiplex

	extra []byte

	// exbuf is for holding the reference to the beginning of the extra slice
	// for later memory pool freeing
	exbuf []byte

	rDeadline, wDeadline pipeDeadline

	clLock                        sync.Mutex
	writeCancelErr, readCancelErr error
	writeCancel, readCancel       chan struct{}
}

func (s *Stream) Name() string {
	return s.name
}

// tries to preload pending data
func (s *Stream) preloadData() {
	select {
	case read, ok := <-s.dataIn:
		if !ok {
			return
		}
		s.extra = read
		s.exbuf = read
	default:
	}
}

func (s *Stream) waitForData() error {
	select {
	case read, ok := <-s.dataIn:
		if !ok {
			return io.EOF
		}
		s.extra = read
		s.exbuf = read
		return nil
	case <-s.readCancel:
		// This is the only place where it's safe to return these.
		s.returnBuffers()
		return s.readCancelErr
	case <-s.rDeadline.wait():
		return errTimeout
	}
}

func (s *Stream) returnBuffers() {
	if s.exbuf != nil {
		pool.Put(s.exbuf)
		s.exbuf = nil
		s.extra = nil
	}
	for {
		select {
		case read, ok := <-s.dataIn:
			if !ok {
				return
			}
			if read == nil {
				continue
			}
			pool.Put(read)
		default:
			return
		}
	}
}

func (s *Stream) Read(b []byte) (int, error) {
	select {
	case <-s.readCancel:
		return 0, s.readCancelErr
	default:
	}

	if s.extra == nil {
		err := s.waitForData()
		if err != nil {
			return 0, err
		}
	}
	n := 0
	for s.extra != nil && n < len(b) {
		read := copy(b[n:], s.extra)
		n += read
		if read < len(s.extra) {
			s.extra = s.extra[read:]
		} else {
			if s.exbuf != nil {
				pool.Put(s.exbuf)
			}
			s.extra = nil
			s.exbuf = nil
			s.preloadData()
		}
	}
	return n, nil
}

func (s *Stream) Write(b []byte) (int, error) {
	var written int
	for written < len(b) {
		wl := len(b) - written
		if wl > MaxMessageSize {
			wl = MaxMessageSize
		}

		n, err := s.write(b[written : written+wl])
		if err != nil {
			return written, err
		}

		written += n
	}

	return written, nil
}

func (s *Stream) write(b []byte) (int, error) {
	select {
	case <-s.writeCancel:
		return 0, s.writeCancelErr
	default:
	}

	err := s.mp.sendMsg(s.wDeadline.wait(), s.writeCancel, s.id.header(messageTag), b)
	if err != nil {
		return 0, err
	}

	return len(b), nil
}

func (s *Stream) cancelWrite(err error) bool {
	s.wDeadline.close()

	s.clLock.Lock()
	defer s.clLock.Unlock()
	select {
	case <-s.writeCancel:
		return false
	default:
		s.writeCancelErr = err
		close(s.writeCancel)
		return true
	}
}

func (s *Stream) cancelRead(err error) bool {
	// Always unregister for reading first, even if we're already closed (or
	// already closing). When handleIncoming calls this, it expects the
	// stream to be unregistered by the time it returns.
	s.mp.chLock.Lock()
	delete(s.mp.channels, s.id)
	s.mp.chLock.Unlock()

	s.rDeadline.close()

	s.clLock.Lock()
	defer s.clLock.Unlock()
	select {
	case <-s.readCancel:
		return false
	default:
		s.readCancelErr = err
		close(s.readCancel)
		return true
	}
}

func (s *Stream) CloseWrite() error {
	if !s.cancelWrite(ErrStreamClosed) {
		// Check if we closed the stream _nicely_. If so, we don't need
		// to report an error to the user.
		if s.writeCancelErr == ErrStreamClosed {
			return nil
		}
		// Closed for some other reason. Report it.
		return s.writeCancelErr
	}

	ctx, cancel := context.WithTimeout(context.Background(), ResetStreamTimeout)
	defer cancel()

	err := s.mp.sendMsg(ctx.Done(), nil, s.id.header(closeTag), nil)
	// We failed to close the stream after 2 minutes, something is probably wrong.
	if err != nil && !s.mp.isShutdown() {
		log.Warnf("Error closing stream: %s; killing connection", err.Error())
		s.mp.Close()
	}
	return err
}

func (s *Stream) CloseRead() error {
	s.cancelRead(ErrStreamClosed)
	return nil
}

func (s *Stream) Close() error {
	return multierr.Combine(s.CloseRead(), s.CloseWrite())
}

func (s *Stream) Reset() error {
	s.cancelRead(ErrStreamReset)

	if s.cancelWrite(ErrStreamReset) {
		// Send a reset in the background.
		go s.mp.sendResetMsg(s.id.header(resetTag), true)
	}

	return nil
}

func (s *Stream) SetDeadline(t time.Time) error {
	s.rDeadline.set(t)
	s.wDeadline.set(t)
	return nil
}

func (s *Stream) SetReadDeadline(t time.Time) error {
	s.rDeadline.set(t)
	return nil
}

func (s *Stream) SetWriteDeadline(t time.Time) error {
	s.wDeadline.set(t)
	return nil
}
