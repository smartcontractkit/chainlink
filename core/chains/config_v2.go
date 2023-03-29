package chains

type configsV2AsV1[I ID, N Node] struct {
	*configChains[I]
	*configNodes[I, N]
}

type ConfigsV2[I ID, N Node] interface {
	chainConfigsV2[I]
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
type configChains[I ID] struct {
	v2 chainConfigsV2[I]
}

type chainConfigsV2[I ID] interface {
	Chains(ids ...I) ([]ChainConfig, error)
}

// newConfigChains returns a chains backed by chains.
func newConfigChains[I ID](d chainConfigsV2[I]) *configChains[I] {
	return &configChains[I]{v2: d}
}

func (o *configChains[I]) Chains(offset, limit int, ids ...I) (chains []ChainConfig, count int, err error) {
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

// configNodes is a generic Configs for nodes.
type configNodes[I ID, N Node] struct {
	v2 nodeConfigsV2[I, N]
}

type nodeConfigsV2[I ID, N Node] interface {
	Node(name string) (N, error)
	Nodes() []N
	NodesByID(...I) []N
}

func newConfigNodes[I ID, N Node](d nodeConfigsV2[I, N]) *configNodes[I, N] {
	return &configNodes[I, N]{v2: d}
}

func (o *configNodes[I, N]) NodeNamed(name string) (node N, err error) {
	return o.v2.Node(name)
}

func (o *configNodes[I, N]) Nodes(offset, limit int) (nodes []N, count int, err error) {
	nodes = o.v2.Nodes()
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

func (o *configNodes[I, N]) NodesForChain(chainID I, offset, limit int) (nodes []N, count int, err error) {
	nodes = o.v2.NodesByID(chainID)
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

func (o *configNodes[I, N]) GetNodesByChainIDs(chainIDs []I) (nodes []N, err error) {
	nodes = o.v2.NodesByID(chainIDs...)
	return
}
