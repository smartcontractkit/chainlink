package chains

import (
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
)

type ChainConfigs interface {
	Chains(offset, limit int, ids ...string) ([]types.ChainStatus, int, error)
}

type NodeConfigs[I ID, N Node] interface {
	Node(name string) (N, error)
	Nodes(chainID I) (nodes []N, err error)

	NodeStatus(name string) (types.NodeStatus, error)
	NodeStatusesPaged(offset, limit int, chainIDs ...string) (nodes []types.NodeStatus, count int, err error)
}

// Configs holds chain and node configurations.
type Configs[I ID, N Node] interface {
	ChainConfigs
	NodeConfigs[I, N]
}
