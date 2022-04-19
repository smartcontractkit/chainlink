package evm

import (
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type orm struct {
	*chains.ChainsORM[utils.Big, types.ChainCfg, types.Chain]
	*chains.NodesORM[utils.Big, types.NewNode, types.Node]
}

var _ types.ORM = (*orm)(nil)

// NewORM returns a new EVM ORM
func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.LogConfig) types.ORM {
	q := pg.NewQ(db, lggr.Named("EVMORM"), cfg)
	return &orm{
		chains.NewChainsORM[utils.Big, types.ChainCfg, types.Chain](q, "evm"),
		chains.NewNodesORM[utils.Big, types.NewNode, types.Node](q, "evm", "ws_url", "http_url", "send_only"),
	}
}

func (o *orm) EnabledChainsWithNodes() ([]types.Chain, map[string][]types.Node, error) {
	chains, err := o.EnabledChains()
	if err != nil {
		return nil, nil, err
	}
	nodes, _, err := o.Nodes(0, -1)
	if err != nil {
		return nil, nil, err
	}
	nodemap := make(map[string][]types.Node)
	for _, n := range nodes {
		id := n.EVMChainID.String()
		nodemap[id] = append(nodemap[id], n)
	}
	return chains, nodemap, nil
}
