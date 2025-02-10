package chainlink

import (
	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"

	"github.com/smartcontractkit/chainlink-integrations/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/config"
)

type GeneralConfig interface {
	config.AppConfig
	toml.HasEVMConfigs
	CosmosConfigs() RawConfigs
	SolanaConfigs() solcfg.TOMLConfigs
	StarknetConfigs() RawConfigs
	AptosConfigs() RawConfigs
	TronConfigs() RawConfigs
	// ConfigTOML returns both the user provided and effective configuration as TOML.
	ConfigTOML() (user, effective string)
}
