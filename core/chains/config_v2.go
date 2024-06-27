package chains

import "github.com/smartcontractkit/chainlink-relay/pkg/types"

type configsV2AsV1[I ID, N Node] struct {
	*configChains
	*configNodes[I, N]
}

type ConfigsV2[I ID, N Node] interface {
	chainConfigsV2
	nodeConfigsV2[I, N]
}

// NewConfigs returns a [Configs] backed by [ConfigsV2].
func NewConfigs[I ID, N Node](cfgs ConfigsV2[I, N]) Configs[I, N] {
	return configsV2AsV1[I, N]{
		newConfigChains[I](cfgs),
		newConfigNodes[I, N](cfgs),
	}
}

// configChains is a generic, immutable Configs for chains.
type configChains struct {
	v2 chainConfigsV2
}

type chainConfigsV2 interface {
	Chains(ids ...string) ([]types.ChainStatus, error)
}

// newConfigChains returns a chains backed by chains.
func newConfigChains[I ID](d chainConfigsV2) *configChains {
	return &configChains{v2: d}
}

func (o *configChains) Chains(offset, limit int, ids ...string) (chains []types.ChainStatus, count int, err error) {
	chains, err = o.v2.Chains(ids...)
	if err != nil {
		return
	}
	count = len(chains)
	if offset < len(chains) {
		chains = chains[offset:]
	} else {
		chains = nil
	}
	if limit > 0 && len(chains) > limit {
		chains = chains[:limit]
	}
	return
}

type nodeConfigsV2[I ID, N Node] interface {
	Node(name string) (N, error)
	Nodes(chainID I) ([]N, error)

	NodeStatus(name string) (types.NodeStatus, error)
	NodeStatuses(chainIDs ...string) (nodes []types.NodeStatus, err error)
}

// configNodes is a generic Configs for nodes.
type configNodes[I ID, N Node] struct {
	nodeConfigsV2[I, N]
}

func newConfigNodes[I ID, N Node](d nodeConfigsV2[I, N]) *configNodes[I, N] {
	return &configNodes[I, N]{d}
}

func (o *configNodes[I, N]) NodeStatusesPaged(offset, limit int, chainIDs ...string) (nodes []types.NodeStatus, count int, err error) {
	nodes, err = o.nodeConfigsV2.NodeStatuses(chainIDs...)
	if err != nil {
		return
	}
	count = len(nodes)
	if offset < len(nodes) {
		nodes = nodes[offset:]
	} else {
		nodes = nil
	}
	if limit > 0 && len(nodes) > limit {
		nodes = nodes[:limit]
	}
	return
}
