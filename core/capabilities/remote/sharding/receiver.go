package sharding

import (
	"slices"

	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func isPeerInDON(peer p2ptypes.PeerID, members []p2ptypes.PeerID) bool {
	return slices.Contains(members, peer)
}
