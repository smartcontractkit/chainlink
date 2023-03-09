package types

import (
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/db"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

// ORM manages cosmos chains and nodes.
type ORM interface {
	Chain(string, ...pg.QOpt) (DBChain, error)
	Chains(offset, limit int, qopts ...pg.QOpt) ([]DBChain, int, error)
	GetChainsByIDs(ids []string) (chains []DBChain, err error)
	EnabledChains(...pg.QOpt) ([]DBChain, error)

	GetNodesByChainIDs(chainIDs []string, qopts ...pg.QOpt) (nodes []db.Node, err error)
	NodeNamed(string, ...pg.QOpt) (db.Node, error)
	Nodes(offset, limit int, qopts ...pg.QOpt) (nodes []db.Node, count int, err error)
	NodesForChain(chainID string, offset, limit int, qopts ...pg.QOpt) (nodes []db.Node, count int, err error)

	EnsureChains([]string, ...pg.QOpt) error
}

type DBChain = chains.DBChain[string, *db.ChainCfg]

// NewNode defines a new node to create.
type NewNode struct {
	Name          string `json:"name"`
	CosmosChainID string `json:"cosmosChainId"`
	TendermintURL string `json:"tendermintURL" db:"tendermint_url"`
}
