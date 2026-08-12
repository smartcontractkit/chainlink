// Package types re-exports the DON-to-DON peer types now living in capabilities/libs/x/rage, so
// callers that only need the type/interface definitions (not the peer implementation itself, which
// moved) can keep importing this path unchanged.
package types

import (
	"github.com/smartcontractkit/capabilities/libs/x/rage"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"
)

type (
	PeerID       = ragetypes.PeerID
	Peer         = rage.Peer
	DonPair      = rage.DonPair
	SharedPeer   = rage.SharedPeer
	PeerWrapper  = rage.PeerWrapper
	Signer       = rage.Signer
	Message      = rage.Message
	StreamConfig = rage.StreamConfig
)
