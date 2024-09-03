package chainlink

import (
	coscfg "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/config"
	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	stkcfg "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/config"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/config"
)

type GeneralConfig interface {
	config.AppConfig
	toml.HasEVMConfigs
	CosmosConfigs() coscfg.TOMLConfigs
	SolanaConfigs() solcfg.TOMLConfigs
	StarknetConfigs() stkcfg.TOMLConfigs
	AptosConfigs() RawConfigs
	// ConfigTOML returns both the user provided and effective configuration as TOML.
	ConfigTOML() (user, effective string)
}
