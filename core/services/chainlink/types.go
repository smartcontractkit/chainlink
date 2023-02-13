package chainlink

import (
	v2 "github.com/smartcontractkit/chainlink/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/core/chains/solana"
	"github.com/smartcontractkit/chainlink/core/chains/starknet"
	"github.com/smartcontractkit/chainlink/core/chains/terra"
	"github.com/smartcontractkit/chainlink/core/config"
)

//go:generate mockery --quiet --name GeneralConfig --output ./mocks/ --case=underscore

type GeneralConfig interface {
	config.GeneralConfig
	v2.HasEVMConfigs
	SolanaConfigs() solana.SolanaConfigs
	StarknetConfigs() starknet.StarknetConfigs
	TerraConfigs() terra.TerraConfigs
	// ConfigTOML returns both the user provided and effective configuration as TOML.
	ConfigTOML() (user, effective string)
}
