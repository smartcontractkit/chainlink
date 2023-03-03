package peerstream_multiplex

import (
	"net"

	"github.com/libp2p/go-libp2p-core/mux"

	mp "github.com/libp2p/go-mplex"
)

// DefaultTransport has default settings for Transport
var DefaultTransport = &Transport{}

// Transport implements mux.Multiplexer that constructs
// mplex-backed muxed connections.
type Transport struct{}

func (t *Transport) NewConn(nc net.Conn, isServer bool) (mux.MuxedConn, error) {
	return (*conn)(mp.NewMultiplex(nc, isServer)), nil
}

var _ mux.Multiplexer = &Transport{}
