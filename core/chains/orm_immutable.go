package chains

import (
	"github.com/pkg/errors"

	v2 "github.com/smartcontractkit/chainlink/core/config/v2"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

type configs[I ID, C Config, N Node] struct {
	*configChains[I, C]
	*configNodes[I, N]
}

type Configs[I ID, C Config, N Node] interface {
	chains[I, C]
	nodes[I, N]
}

// NewORMImmut returns an ORM backed by q, for the tables <prefix>_chains and <prefix>_nodes with column <prefix>_chain_id.
// Additional Node fields should be included in nodeCols.
func NewORMImmut[I ID, C Config, N Node](chainConfigs Configs[I, C, N]) ORM[I, C, N] {
	return configs[I, C, N]{
		newConfigChains[I, C](chainConfigs),
		newConfigNodes[I, N](chainConfigs),
	}
}

func (o configs[I, C, N]) EnsureChains(_ []I, _ ...pg.QOpt) error {
	return v2.ErrUnsupported
}

// configChains is a generic, immutable ORM for chains.
type configChains[I ID, C Config] struct {
	data chains[I, C]
}

type chains[I ID, C Config] interface {
	// Chains returns a slice of ChainConfig for ids, or all if none are provided.
	Chains(ids ...I) []ChainConfig[I, C]
}

// newConfigChains returns a chains backed by chains.
func newConfigChains[I ID, C Config](d chains[I, C]) *configChains[I, C] {
	return &configChains[I, C]{data: d}
}

func (o *configChains[I, C]) Chain(id I, _ ...pg.QOpt) (cc ChainConfig[I, C], err error) {
	chains := o.data.Chains(id)
	if len(chains) == 0 {
		err = errors.Errorf("chain not found: %v", id)
		return
	} else if len(chains) > 1 {
		err = errors.Errorf("more than one chain found: %v", id)
		return
	}
	cc = chains[0]
	return
}

func (o *configChains[I, C]) GetChainsByIDs(ids []I) (chains []ChainConfig[I, C], err error) {
	return o.data.Chains(ids...), nil
}

func (o *configChains[I, C]) Chains(offset, limit int, _ ...pg.QOpt) (chains []ChainConfig[I, C], count int, err error) {
	chains = o.data.Chains()
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

// configNodes is a generic ORM for nodes.
type configNodes[I ID, N Node] struct {
	data nodes[I, N]
}

type nodes[I ID, N Node] interface {
	Node(name string) (N, error)
	Nodes() []N
	NodesByID(...I) []N
}

func newConfigNodes[I ID, N Node](d nodes[I, N]) *configNodes[I, N] {
	return &configNodes[I, N]{data: d}
}

func (o *configNodes[I, N]) NodeNamed(name string, _ ...pg.QOpt) (node N, err error) {
	return o.data.Node(name)
}

func (o *configNodes[I, N]) Nodes(offset, limit int, _ ...pg.QOpt) (nodes []N, count int, err error) {
	nodes = o.data.Nodes()
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

func (o *configNodes[I, N]) NodesForChain(chainID I, offset, limit int, _ ...pg.QOpt) (nodes []N, count int, err error) {
	nodes = o.data.NodesByID(chainID)
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

func (o *configNodes[I, N]) GetNodesByChainIDs(chainIDs []I, _ ...pg.QOpt) (nodes []N, err error) {
	nodes = o.data.NodesByID(chainIDs...)
	return
}
