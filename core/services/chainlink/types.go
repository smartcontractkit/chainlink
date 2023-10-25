package chainlink

import (
	"github.com/smartcontractkit/chainlink-solana/pkg/solana"
	stkcfg "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/config"

	"github.com/smartcontractkit/chainlink/v2/core/chains/cosmos"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/config"
)

//go:generate mockery --quiet --name GeneralConfig --output ./mocks/ --case=underscore

type GeneralConfig interface {
	config.AppConfig
	toml.HasEVMConfigs
	CosmosConfigs() cosmos.CosmosConfigs
	SolanaConfigs() solana.TOMLConfigs
	StarknetConfigs() stkcfg.TOMLConfigs
	// ConfigTOML returns both the user provided and effective configuration as TOML.
	ConfigTOML() (user, effective string)
}
