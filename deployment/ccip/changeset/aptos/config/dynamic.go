package config

import (
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type DynamicConfig struct {
	Defs          []operations.Definition
	Inputs        []any // Each element should be the corresponding input type for its operation
	ChainSelector uint64
	Description   string
	MCMSConfig    *cldfproposalutils.TimelockConfig
}
