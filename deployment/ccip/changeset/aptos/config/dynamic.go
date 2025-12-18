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

// RMNRemoteCurseInput config for cursing/uncursing multiple subjects on Aptos RMN Remote.
type RMNRemoteCurseInput struct {
	ChainSelector uint64
	Subjects      [][]byte
	MCMSConfig    *proposalutils.TimelockConfig
}
