package types

import (
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"

	"github.com/smartcontractkit/chainlink/core/chains"
)

type Configs interface {
	chains.ChainConfigs[string]
	chains.NodeConfigs[string, db.Node]
}
