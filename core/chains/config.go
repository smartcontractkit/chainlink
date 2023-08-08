package chains

import (
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
)

type ChainConfig interface {
	ChainStatus() (types.ChainStatus, error)
}

type NodeConfigs[N Node] interface {
	//Node(name string) (N, error)
	Nodes(names ...string) (nodes []N, err error)

	NodeStatus(name string) (types.NodeStatus, error)
	NodeStatusesPaged(offset, limit int) (nodes []types.NodeStatus, count int, err error)
}

// Statuser holds chain and node configurations.
type Statuser[N Node] interface {
	ChainConfig
	NodeConfigs[N]
}
