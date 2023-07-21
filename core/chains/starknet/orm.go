package starknet

import (
	starknetdb "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
	"github.com/smartcontractkit/chainlink/v2/core/chains/starknet/types"
)

func NewConfigs(cfgs chains.ConfigsV2[string, starknetdb.Node]) types.Configs {
	return chains.NewConfigs(cfgs)
}
