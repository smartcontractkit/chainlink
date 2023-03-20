package types

import (
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

type ORM interface {
	chains.ChainConfigs[string, *db.ChainCfg, ChainConfig]
	chains.NodeConfigs[string, db.Node]

	EnsureChains([]string, ...pg.QOpt) error
}

type ChainConfig = chains.ChainConfig[string, *db.ChainCfg]
