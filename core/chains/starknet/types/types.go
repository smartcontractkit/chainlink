package types

import (
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
)

type Configs interface {
	chains.ChainConfigs
	chains.NodeConfigs[string, db.Node]
}
