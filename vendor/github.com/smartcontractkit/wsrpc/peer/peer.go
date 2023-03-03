package peer

import (
	"context"

	"github.com/smartcontractkit/wsrpc/credentials"
)

// Peer contains the information of the peer for an RPC.
type Peer struct {
	// The public key that the peer connected with
	PublicKey credentials.StaticSizedPublicKey
}

type peerKey struct{}

// NewContext creates a new context with peer information attached.
func NewContext(ctx context.Context, p *Peer) context.Context {
	return context.WithValue(ctx, peerKey{}, p)
}

// FromContext returns the peer information in ctx if it exists.
func FromContext(ctx context.Context) (p *Peer, ok bool) {
	p, ok = ctx.Value(peerKey{}).(*Peer)

	return
}

// NewCallContext creates a new context with the peer information of the RPC
// call.
//
// Used when making server side RPC calls to the client so we can identify the
// correct client to send the RPC call.
func NewCallContext(ctx context.Context, pubkey credentials.StaticSizedPublicKey) context.Context {
	return NewContext(ctx, &Peer{PublicKey: pubkey})
}
