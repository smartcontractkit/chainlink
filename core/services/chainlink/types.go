package chainlink

import (
	"github.com/smartcontractkit/chainlink/v2/core/chains/cosmos"
	v2 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/chains/solana"
	"github.com/smartcontractkit/chainlink/v2/core/chains/starknet"
	"github.com/smartcontractkit/chainlink/v2/core/config"
)

//go:generate mockery --quiet --name GeneralConfig --output ./mocks/ --case=underscore

type GeneralConfig interface {
	config.GeneralConfig
	v2.HasEVMConfigs
	CosmosConfigs() cosmos.CosmosConfigs
	SolanaConfigs() solana.SolanaConfigs
	StarknetConfigs() starknet.StarknetConfigs
	// ConfigTOML returns both the user provided and effective configuration as TOML.
	ConfigTOML() (user, effective string)
}
