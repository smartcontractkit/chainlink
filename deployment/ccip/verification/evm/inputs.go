package evm

import (
	"encoding/json"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_0_0/rmn_proxy_contract"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/token_admin_registry"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_from_mint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_with_from_mint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/factory_burn_mint_erc20"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/token_pool_factory"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/ccip_home"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/nonce_manager"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/registry_module_owner_custom"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_home"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_remote"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
)

type rawContractInfo struct {
	solidityStandardJSONInput string
	bytecode                  string
	name                      string
}

// contracts maps type & version of a contract to its corresponding standard JSON input, name, and bytecode.
// TODO: USDCTokenPool, HybridLockReleaseUSDCTokenPool, CapabilitiesRegistry, WETH9? (not urgent as these are infrequently deployed)
var contracts map[deployment.ContractType]map[semver.Version]rawContractInfo = map[deployment.ContractType]map[semver.Version]rawContractInfo{
	changeset.RMNRemote: {
		deployment.Version1_6_0: rawContractInfo{
			solidityStandardJSONInput: rmn_remote.SolidityStandardInput,
			bytecode:                  rmn_remote.RMNRemoteBin,
			name:                      "contracts/rmn/RMNRemote.sol:RMNRemote",
		},
	},
	changeset.TokenPoolFactory: {
		deployment.Version1_5_1: rawContractInfo{
			solidityStandardJSONInput: token_pool_factory.SolidityStandardInput,
			bytecode:                  token_pool_factory.TokenPoolFactoryBin,
			name:                      "contracts/tokenAdminRegistry/TokenPoolFactory/TokenPoolFactory.sol:TokenPoolFactory",
		},
	},
	changeset.RegistryModule: {
		deployment.Version1_6_0: rawContractInfo{
			solidityStandardJSONInput: registry_module_owner_custom.SolidityStandardInput,
			bytecode:                  registry_module_owner_custom.RegistryModuleOwnerCustomBin,
			name:                      "contracts/tokenAdminRegistry/RegistryModuleOwnerCustom.sol:RegistryModuleOwnerCustom",
		},
	},
	changeset.NonceManager: {
		deployment.Version1_6_0: rawContractInfo{
			solidityStandardJSONInput: nonce_manager.SolidityStandardInput,
			bytecode:                  nonce_manager.NonceManagerBin,
			name:                      "contracts/NonceManager.sol:NonceManager",
		},
	},
	changeset.FeeQuoter: {
		deployment.Version1_6_0: rawContractInfo{
			solidityStandardJSONInput: fee_quoter.SolidityStandardInput,
			bytecode:                  fee_quoter.FeeQuoterBin,
			name:                      "contracts/FeeQuoter.sol:FeeQuoter",
		},
	},
	changeset.CCIPHome: {
		deployment.Version1_6_0: rawContractInfo{
			solidityStandardJSONInput: ccip_home.SolidityStandardInput,
			bytecode:                  ccip_home.CCIPHomeBin,
			name:                      "contracts/capability/CCIPHome.sol:CCIPHome",
		},
	},
	changeset.RMNHome: {
		deployment.Version1_6_0: rawContractInfo{
			solidityStandardJSONInput: rmn_home.SolidityStandardInput,
			bytecode:                  rmn_home.RMNHomeBin,
			name:                      "contracts/rmn/RMNHome.sol:RMNHome",
		},
	},
	changeset.OnRamp: {
		deployment.Version1_6_0: rawContractInfo{
			solidityStandardJSONInput: onramp.SolidityStandardInput,
			bytecode:                  onramp.OnRampBin,
			name:                      "contracts/onRamp/OnRamp.sol:OnRamp",
		},
	},
	changeset.OffRamp: {
		deployment.Version1_6_0: rawContractInfo{
			solidityStandardJSONInput: offramp.SolidityStandardInput,
			bytecode:                  offramp.OffRampBin,
			name:                      "contracts/offRamp/OffRamp.sol:OffRamp",
		},
	},
	changeset.FactoryBurnMintERC20Token: {
		deployment.Version1_5_1: rawContractInfo{
			solidityStandardJSONInput: factory_burn_mint_erc20.SolidityStandardInput,
			bytecode:                  factory_burn_mint_erc20.FactoryBurnMintERC20Bin,
			name:                      "contracts/tokenAdminRegistry/TokenPoolFactory/FactoryBurnMintERC20.sol:FactoryBurnMintERC20",
		},
	},
	changeset.BurnMintTokenPool: {
		deployment.Version1_5_1: rawContractInfo{
			solidityStandardJSONInput: burn_mint_token_pool.SolidityStandardInput,
			bytecode:                  burn_mint_token_pool.BurnMintTokenPoolBin,
			name:                      "contracts/pools/BurnMintTokenPool.sol:BurnMintTokenPool",
		},
	},
	changeset.BurnWithFromMintTokenPool: {
		deployment.Version1_5_1: rawContractInfo{
			solidityStandardJSONInput: burn_with_from_mint_token_pool.SolidityStandardInput,
			bytecode:                  burn_with_from_mint_token_pool.BurnWithFromMintTokenPoolBin,
			name:                      "contracts/pools/BurnWithFromMintTokenPool.sol:BurnWithFromMintTokenPool",
		},
	},
	changeset.BurnFromMintTokenPool: {
		deployment.Version1_5_1: rawContractInfo{
			solidityStandardJSONInput: burn_from_mint_token_pool.SolidityStandardInput,
			bytecode:                  burn_from_mint_token_pool.BurnFromMintTokenPoolBin,
			name:                      "contracts/pools/BurnFromMintTokenPool.sol:BurnFromMintTokenPool",
		},
	},
	changeset.LockReleaseTokenPool: {
		deployment.Version1_5_1: rawContractInfo{
			solidityStandardJSONInput: lock_release_token_pool.SolidityStandardInput,
			bytecode:                  lock_release_token_pool.LockReleaseTokenPoolBin,
			name:                      "contracts/pools/LockReleaseTokenPool.sol:LockReleaseTokenPool",
		},
	},
	changeset.TokenAdminRegistry: {
		deployment.Version1_5_0: rawContractInfo{
			solidityStandardJSONInput: token_admin_registry.SolidityStandardInput,
			bytecode:                  token_admin_registry.TokenAdminRegistryBin,
			name:                      "contracts/tokenAdminRegistry/TokenAdminRegistry.sol:TokenAdminRegistry",
		},
	},
	changeset.Router: {
		deployment.Version1_2_0: rawContractInfo{
			solidityStandardJSONInput: router.SolidityStandardInput,
			bytecode:                  router.RouterBin,
			name:                      "contracts/Router.sol:Router",
		},
	},
	changeset.ARMProxy: {
		deployment.Version1_0_0: rawContractInfo{
			solidityStandardJSONInput: rmn_proxy_contract.SolidityStandardInput,
			bytecode:                  rmn_proxy_contract.RMNProxyBin,
			name:                      "contracts/rmn/RMNProxy.sol:RMNProxy",
		},
	},
}

// loadSolidityContractMetadata loads the metadata for a contract type and version, including the standard JSON input, bytecode, and name.
func loadSolidityContractMetadata(contractType deployment.ContractType, version semver.Version) (solidityContractMetadata, error) {
	contract, ok := contracts[contractType]
	if !ok {
		return solidityContractMetadata{}, fmt.Errorf("no contract found for type %s", contractType)
	}
	contractWithVersion, ok := contract[version]
	if !ok {
		return solidityContractMetadata{}, fmt.Errorf("no contract found for type %s with version %s", contractType, version)
	}

	var input solidityContractMetadata
	err := json.Unmarshal([]byte(contractWithVersion.solidityStandardJSONInput), &input)
	if err != nil {
		return solidityContractMetadata{}, fmt.Errorf("failed to unmarshal solidity standard JSON input for contract type %s: %w", contractType, err)
	}
	// Add remaining fields that don't exist in the standard JSON input
	input.Bytecode = contractWithVersion.bytecode
	input.Name = contractWithVersion.name

	return input, nil
}

// solidityContractMetadata defines the metadata for a Solidity contract, including the standard JSON input, bytecode, and contract name.
type solidityContractMetadata struct {
	Version  string         `json:"version"`
	Language string         `json:"language"`
	Settings map[string]any `json:"settings"`
	Sources  map[string]any `json:"sources"`
	Bytecode string         `json:"bytecode"`
	Name     string         `json:"name"`
}

// SourceCode returns the source code of the contract as a string.
func (s solidityContractMetadata) SourceCode() (string, error) {
	sourceCodeMap := map[string]any{
		"language": s.Language,
		"settings": s.Settings,
		"sources":  s.Sources,
	}
	jsonBytes, err := json.Marshal(sourceCodeMap)
	if err != nil {
		return "", fmt.Errorf("failed to marshal source code: %w", err)
	}
	return string(jsonBytes), nil
}
