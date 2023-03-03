package peerstore

import core "github.com/libp2p/go-libp2p-core/peerstore"

// Deprecated: use github.com/libp2p/go-libp2p-core/peerstore.ErrNotFound instead.
var ErrNotFound = core.ErrNotFound

var (
	// Deprecated: use github.com/libp2p/go-libp2p-core/peerstore.AddressTTL instead.
	AddressTTL = core.AddressTTL

	// Deprecated: use github.com/libp2p/go-libp2p-core/peerstore.TempAddrTTL instead.
	TempAddrTTL = core.TempAddrTTL

	// Deprecated: use github.com/libp2p/go-libp2p-core/peerstore.ProviderAddrTTL instead.
	ProviderAddrTTL = core.ProviderAddrTTL

	// Deprecated: use github.com/libp2p/go-libp2p-core/peerstore.RecentlyConnectedAddrTTL instead.
	RecentlyConnectedAddrTTL = core.RecentlyConnectedAddrTTL

	// Deprecated: use github.com/libp2p/go-libp2p-core/peerstore.OwnObservedAddrTTL instead.
	OwnObservedAddrTTL = core.OwnObservedAddrTTL
)

const (
	// Deprecated: use github.com/libp2p/go-libp2p-core/peerstore.PermanentAddrTTL instead.
	PermanentAddrTTL = core.PermanentAddrTTL

	// Deprecated: use github.com/libp2p/go-libp2p-core/peerstore.ConnectedAddrTTL instead.
	ConnectedAddrTTL = core.ConnectedAddrTTL
)

// Deprecated: use github.com/libp2p/go-libp2p-core/peerstore.Peerstore instead.
type Peerstore = core.Peerstore

// Deprecated: use github.com/libp2p/go-libp2p-core/peerstore.PeerMetadata instead.
type PeerMetadata = core.PeerMetadata

// Deprecated: use github.com/libp2p/go-libp2p-core/peerstore.AddrBook instead.
type AddrBook = core.AddrBook

// Deprecated: use github.com/libp2p/go-libp2p-core/peerstore.KeyBook instead.
type KeyBook = core.KeyBook

// Deprecated: use github.com/libp2p/go-libp2p-core/peerstore.ProtoBook instead.
type ProtoBook = core.ProtoBook
