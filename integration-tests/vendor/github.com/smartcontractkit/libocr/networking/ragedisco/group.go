package ragedisco

import (
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"
)

type group struct {
	oracleNodes       []ragetypes.PeerID
	bootstrapperNodes []ragetypes.PeerInfo
}

func (g *group) oracleIDs() []ragetypes.PeerID {
	return g.oracleNodes
}

func (g *group) bootstrapperIDs() (ps []ragetypes.PeerID) {
	for _, inf := range g.bootstrapperNodes {
		ps = append(ps, inf.ID)
	}
	return
}

func (g *group) peerIDs() []ragetypes.PeerID {
	return append(g.oracleIDs(), g.bootstrapperIDs()...)
}

func (g *group) hasOracle(hpid ragetypes.PeerID) bool {
	for _, pid := range g.oracleNodes {
		if pid == hpid {
			return true
		}
	}
	return false
}
