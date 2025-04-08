package chainlink

import (
	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"

	"github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	coreconfig "github.com/smartcontractkit/chainlink/v2/core/config"
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
	ImportedSecretConfig
}

// ImportedSecretConfig is a configuration for imported secrets
// to be imported into the keystore upon startup.
type ImportedSecretConfig interface {
	ImportedP2PKey() coreconfig.ImportableKey
	ImportedEthKeys() coreconfig.ImportableEthKeyLister
}
