package multistream

import (
	"io"
	"sync"
)

// lazyServerConn is an io.ReadWriteCloser adapter used for negotiating inbound
// streams (see NegotiateLazy).
//
// This is "lazy" because it doesn't wait for the write half to succeed before
// allowing us to read from the stream.
type lazyServerConn struct {
	waitForHandshake sync.Once
	werr             error

	con io.ReadWriteCloser
}

func (l *lazyServerConn) Write(b []byte) (int, error) {
	l.waitForHandshake.Do(func() { panic("didn't initiate handshake") })
	if l.werr != nil {
		return 0, l.werr
	}
	return l.con.Write(b)
}

func (l *lazyServerConn) Read(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}
	return l.con.Read(b)
}

func (l *lazyServerConn) Close() error {
	// As the server, we MUST flush the handshake on close. Otherwise, if
	// the other side is actually waiting for our close (i.e., reading until
	// EOF), they may get an error even though we received the request.
	//
	// However, we MUST NOT return any errors from Flush. The initiator may
	// have already closed their side for reading. Basically, _we_ don't
	// care about the outcome of this flush, only the other side does.
	_ = l.Flush()
	return l.con.Close()
}

// Flush sends the handshake.
func (l *lazyServerConn) Flush() error {
	l.waitForHandshake.Do(func() { panic("didn't initiate handshake") })
	return l.werr
}
