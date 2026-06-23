package chainlink

import (
	"github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	coreconfig "github.com/smartcontractkit/chainlink/v2/core/config"
)

type GeneralConfig interface {
	coreconfig.AppConfig
	toml.HasEVMConfigs
	CosmosConfigs() RawConfigs
	SolanaConfigs() RawConfigs
	StarknetConfigs() RawConfigs
	AptosConfigs() RawConfigs
	TronConfigs() RawConfigs
	TONConfigs() RawConfigs
	SuiConfigs() RawConfigs
	StellarConfigs() RawConfigs
	// ConfigTOML returns both the user provided and effective configuration as TOML.
	ConfigTOML() (user, effective string)
	ImportedSecretConfig
}

// ImportedSecretConfig is a configuration for imported secrets
// to be imported into the keystore upon startup.
type ImportedSecretConfig interface {
	ImportedP2PKey() coreconfig.ImportableKey
	ImportedEthKeys() coreconfig.ImportableChainKeyLister
	ImportedSolKeys() coreconfig.ImportableChainKeyLister
	ImportedAptosKeys() coreconfig.ImportableChainKeyLister
	ImportedDKGRecipientKey() coreconfig.ImportableKey
}
