package websocket

import (
	"net"
	"time"
)

// GracefulCloseTimeout is the time to wait trying to gracefully close a
// connection before simply cutting it.
var GracefulCloseTimeout = 100 * time.Millisecond

var _ net.Conn = (*Conn)(nil)
