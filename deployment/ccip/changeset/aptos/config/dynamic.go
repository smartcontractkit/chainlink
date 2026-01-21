package config

import (
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

type DynamicConfig struct {
	Defs          []operations.Definition
	Inputs        []any // Each element should be the corresponding input type for its operation
	ChainSelector uint64
	Description   string
	MCMSConfig    *proposalutils.TimelockConfig
}
