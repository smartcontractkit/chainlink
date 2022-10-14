package chains

import (
	"github.com/pkg/errors"

	v2 "github.com/smartcontractkit/chainlink/core/config/v2"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

type ormImmut[I ID, C Config, N Node] struct {
	*chainsORMImmut[I, C]
	*nodesORMImmut[I, N]
}

type ChainConfig[I ID, C Config, N Node] interface {
	chainData[I, C]
	nodeData[I, N]
}

// NewORMImmut returns an ORM backed by q, for the tables <prefix>_chains and <prefix>_nodes with column <prefix>_chain_id.
// Additional Node fields should be included in nodeCols.
func NewORMImmut[I ID, C Config, N Node](chainConfigs ChainConfig[I, C, N]) ORM[I, C, N] {
	return ormImmut[I, C, N]{
		newChainsORMImmut[I, C](chainConfigs),
		newNodesORMImmut[I, N](chainConfigs),
	}
}

func (o ormImmut[I, C, N]) SetupNodes(_ []N, _ []I) error {
	return v2.ErrUnsupported
}

func (o ormImmut[I, C, N]) EnsureChains(_ []I, _ ...pg.QOpt) error {
	return v2.ErrUnsupported
}

// chainsORMImmut is a generic, immutable ORM for chains.
type chainsORMImmut[I ID, C Config] struct {
	data chainData[I, C]
}

type chainData[I ID, C Config] interface {
	// Chains returns a slice of DBChain for ids, or all if none are provided.
	Chains(ids ...I) []DBChain[I, C]
}

// newChainsORMImmut returns an chainsORM backed by q, for the table <prefix>_chains.
func newChainsORMImmut[I ID, C Config](d chainData[I, C]) *chainsORMImmut[I, C] {
	return &chainsORMImmut[I, C]{data: d}
}

func (o *chainsORMImmut[I, C]) Chain(id I, _ ...pg.QOpt) (dbchain DBChain[I, C], err error) {
	chains := o.data.Chains(id)
	if len(chains) == 0 {
		err = errors.Errorf("chain not found: %v", id)
		return
	} else if len(chains) > 1 {
		err = errors.Errorf("more than one chain found: %v", id)
		return
	}
	dbchain = chains[0]
	return
}

func (o *chainsORMImmut[I, C]) GetChainsByIDs(ids []I) (chains []DBChain[I, C], err error) {
	return o.data.Chains(ids...), nil
}

func (o *chainsORMImmut[I, C]) CreateChain(id I, config C, _ ...pg.QOpt) (chain DBChain[I, C], err error) {
	err = v2.ErrUnsupported
	return
}

func (o *chainsORMImmut[I, C]) UpdateChain(id I, enabled bool, config C, _ ...pg.QOpt) (chain DBChain[I, C], err error) {
	err = v2.ErrUnsupported
	return
}

// StoreString saves a string value into the config for the given chain and key
func (o *chainsORMImmut[I, C]) StoreString(chainID I, name, val string) error {
	return v2.ErrUnsupported
}

// Clear deletes a config value for the given chain and key
func (o *chainsORMImmut[I, C]) Clear(chainID I, name string) error {
	return v2.ErrUnsupported
}

func (o *chainsORMImmut[I, C]) DeleteChain(id I, _ ...pg.QOpt) error {
	return v2.ErrUnsupported
}

func (o *chainsORMImmut[I, C]) Chains(offset, limit int, _ ...pg.QOpt) (chains []DBChain[I, C], count int, err error) {
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

func (o *chainsORMImmut[I, C]) EnabledChains(_ ...pg.QOpt) (chains []DBChain[I, C], err error) {
	err = v2.ErrUnsupported
	return
}

// nodesORMImmut is a generic ORM for nodes.
type nodesORMImmut[I ID, N Node] struct {
	data nodeData[I, N]
}

type nodeData[I ID, N Node] interface {
	Node(name string) (N, error)
	Nodes() []N
	NodesByID(...I) []N
}

func newNodesORMImmut[I ID, N Node](d nodeData[I, N]) *nodesORMImmut[I, N] {
	return &nodesORMImmut[I, N]{data: d}
}

func (o *nodesORMImmut[I, N]) CreateNode(data N, _ ...pg.QOpt) (node N, err error) {
	err = v2.ErrUnsupported
	return
}

func (o *nodesORMImmut[I, N]) DeleteNode(id int32, _ ...pg.QOpt) error {
	return v2.ErrUnsupported
}

func (o *nodesORMImmut[I, N]) NodeNamed(name string, _ ...pg.QOpt) (node N, err error) {
	return o.data.Node(name)
}

func (o *nodesORMImmut[I, N]) Nodes(offset, limit int, _ ...pg.QOpt) (nodes []N, count int, err error) {
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

func (o *nodesORMImmut[I, N]) NodesForChain(chainID I, offset, limit int, _ ...pg.QOpt) (nodes []N, count int, err error) {
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

func (o *nodesORMImmut[I, N]) GetNodesByChainIDs(chainIDs []I, _ ...pg.QOpt) (nodes []N, err error) {
	nodes = o.data.NodesByID(chainIDs...)
	return
}
