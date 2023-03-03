// Package sec provides secure connection and transport interfaces for libp2p.
package sec

import (
	"context"
	"net"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

// SecureConn is an authenticated, encrypted connection.
type SecureConn interface {
	net.Conn
	network.ConnSecurity
}

// A SecureTransport turns inbound and outbound unauthenticated,
// plain-text, native connections into authenticated, encrypted connections.
type SecureTransport interface {
	// SecureInbound secures an inbound connection.
	SecureInbound(ctx context.Context, insecure net.Conn) (SecureConn, error)

	// SecureOutbound secures an outbound connection.
	SecureOutbound(ctx context.Context, insecure net.Conn, p peer.ID) (SecureConn, error)
}

// A SecureMuxer is a wrapper around SecureTransport which can select security protocols
// and open outbound connections with simultaneous open.
type SecureMuxer interface {
	// SecureInbound secures an inbound connection.
	// The returned boolean indicates whether the connection should be trated as a server
	// connection; in the case of SecureInbound it should always be true.
	SecureInbound(ctx context.Context, insecure net.Conn) (SecureConn, bool, error)

	// SecureOutbound secures an outbound connection.
	// The returned boolean indicates whether the connection should be treated as a server
	// connection due to simultaneous open.
	SecureOutbound(ctx context.Context, insecure net.Conn, p peer.ID) (SecureConn, bool, error)
}
