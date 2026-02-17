package pkg

import (
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"
)

func PeerIDsToBytes(p2pIDs []p2pkey.PeerID) [][32]byte {
	out := make([][32]byte, len(p2pIDs))
	for i, p2pID := range p2pIDs {
		out[i] = p2pID
	}
	return out
}
