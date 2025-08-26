package test

import "github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"

// P2PIDs is a slice of p2pkey.PeerID with convenient transform methods.
type P2PIDs []p2pkey.PeerID

// Strings returns the string representation of the p2p IDs.
func (ps P2PIDs) Strings() []string {
	out := make([]string, len(ps))
	for i, p := range ps {
		out[i] = p.String()
	}
	return out
}

// Bytes32 returns the byte representation of the p2p IDs.
func (ps P2PIDs) Bytes32() [][32]byte {
	out := make([][32]byte, len(ps))
	for i, p := range ps {
		out[i] = p
	}
	return out
}

// Unique returns a new slice with duplicate p2p IDs removed.
func (ps P2PIDs) Unique() P2PIDs {
	dedup := make(map[p2pkey.PeerID]struct{})
	var out []p2pkey.PeerID
	for _, p := range ps {
		if _, exists := dedup[p]; !exists {
			out = append(out, p)

			dedup[p] = struct{}{}
		}
	}
	return out
}
