package types

import (
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

type ORM interface {
	Chain(string, ...pg.QOpt) (DBChain, error)
	Chains(offset, limit int, qopts ...pg.QOpt) ([]DBChain, int, error)
	GetChainsByIDs(ids []string) (chains []DBChain, err error)

	GetNodesByChainIDs(chainIDs []string, qopts ...pg.QOpt) (nodes []db.Node, err error)
	NodeNamed(string, ...pg.QOpt) (db.Node, error)
	Nodes(offset, limit int, qopts ...pg.QOpt) (nodes []db.Node, count int, err error)
	NodesForChain(chainID string, offset, limit int, qopts ...pg.QOpt) (nodes []db.Node, count int, err error)

	EnsureChains([]string, ...pg.QOpt) error
}

type DBChain = chains.DBChain[string, *db.ChainCfg]
