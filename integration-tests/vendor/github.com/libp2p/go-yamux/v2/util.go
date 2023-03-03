package yamux

import (
	"io"
	"sync"

	pool "github.com/libp2p/go-buffer-pool"
)

// asyncSendErr is used to try an async send of an error
func asyncSendErr(ch chan error, err error) {
	if ch == nil {
		return
	}
	select {
	case ch <- err:
	default:
	}
}

// asyncNotify is used to signal a waiting goroutine
func asyncNotify(ch chan struct{}) {
	select {
	case ch <- struct{}{}:
	default:
	}
}

// min computes the minimum of a set of values
func min(values ...uint32) uint32 {
	m := values[0]
	for _, v := range values[1:] {
		if v < m {
			m = v
		}
	}
	return m
}

// The segmented buffer looks like:
//
//      |     data      | empty space       |
//       < window (10)                     >
//       < len (5)     > < cap (5)         >
//                       < pending (4)    >
//
// As data is read, the buffer gets updated like so:
//
//         |     data   | empty space       |
//          < window (8)                   >
//          < len (3)  > < cap (5)         >
//                       < pending (4)    >
//
// It can then grow as follows (given a "max" of 10):
//
//
//         |     data   | empty space          |
//          < window (10)                     >
//          < len (3)  > < cap (7)            >
//                       < pending (4)    >
//
// Data can then be written into the pending space, expanding len, and shrinking
// cap and pending:
//
//         |     data       | empty space      |
//          < window (10)                     >
//          < len (5)      > < cap (5)        >
//                           < pending (2)>
//
type segmentedBuffer struct {
	cap     uint32
	pending uint32
	len     uint32
	bm      sync.Mutex
	b       [][]byte
}

// NewSegmentedBuffer allocates a ring buffer.
func newSegmentedBuffer(initialCapacity uint32) segmentedBuffer {
	return segmentedBuffer{cap: initialCapacity, b: make([][]byte, 0)}
}

// Len is the amount of data in the receive buffer.
func (s *segmentedBuffer) Len() int {
	s.bm.Lock()
	len := s.len
	s.bm.Unlock()
	return int(len)
}

// Cap is the remaining capacity in the receive buffer.
//
// Note: this is _not_ the same as go's 'cap' function. The total size of the
// buffer is len+cap.
func (s *segmentedBuffer) Cap() uint32 {
	s.bm.Lock()
	cap := s.cap
	s.bm.Unlock()
	return cap
}

// If the space to write into + current buffer size has grown to half of the window size,
// grow up to that max size, and indicate how much additional space was reserved.
func (s *segmentedBuffer) GrowTo(max uint32, force bool) (bool, uint32) {
	s.bm.Lock()
	defer s.bm.Unlock()

	currentWindow := s.cap + s.len
	if currentWindow >= max {
		return force, 0
	}
	delta := max - currentWindow

	if delta < (max/2) && !force {
		return false, 0
	}

	s.cap += delta
	return true, delta
}

func (s *segmentedBuffer) TryReserve(space uint32) bool {
	s.bm.Lock()
	defer s.bm.Unlock()
	if s.cap < s.pending+space {
		return false
	}
	s.pending += space
	return true
}

func (s *segmentedBuffer) Read(b []byte) (int, error) {
	s.bm.Lock()
	defer s.bm.Unlock()
	if len(s.b) == 0 {
		return 0, io.EOF
	}
	n := copy(b, s.b[0])
	if n == len(s.b[0]) {
		pool.Put(s.b[0])
		s.b[0] = nil
		s.b = s.b[1:]
	} else {
		s.b[0] = s.b[0][n:]
	}
	if n > 0 {
		s.len -= uint32(n)
	}
	return n, nil
}

func (s *segmentedBuffer) Append(input io.Reader, length int) error {
	dst := pool.Get(length)
	n := 0
	read := 0
	var err error
	for n < length && err == nil {
		read, err = input.Read(dst[n:])
		n += read
	}
	if err == io.EOF {
		if length == n {
			err = nil
		} else {
			err = io.ErrUnexpectedEOF
		}
	}

	s.bm.Lock()
	defer s.bm.Unlock()
	if n > 0 {
		s.len += uint32(n)
		s.cap -= uint32(n)
		s.pending = s.pending - uint32(length)
		s.b = append(s.b, dst[0:n])
	}
	return err
}
