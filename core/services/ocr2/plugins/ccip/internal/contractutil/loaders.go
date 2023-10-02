package contractutil

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp_1_1_0"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/observability"
)

func LoadOnRamp(onRampAddress common.Address, pluginName string, client client.Client) (evm_2_evm_onramp.EVM2EVMOnRampInterface, semver.Version, error) {
	version, err := ccipconfig.VerifyTypeAndVersion(onRampAddress, client, ccipconfig.EVM2EVMOnRamp)
	if err != nil {
		return nil, semver.Version{}, errors.Wrap(err, "Invalid onRamp contract")
	}

	onRamp, err := observability.NewObservedEvm2EvmOnRamp(onRampAddress, pluginName, client)
	return onRamp, version, err
}

func LoadOnRampDynamicConfig(onRamp evm_2_evm_onramp.EVM2EVMOnRampInterface, version semver.Version, client client.Client) (evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig, error) {
	opts := &bind.CallOpts{}

	switch version.String() {
	case "1.0.0":
		legacyOnramp, err := evm_2_evm_onramp_1_0_0.NewEVM2EVMOnRamp(onRamp.Address(), client)
		if err != nil {
			return evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{}, err
		}
		legacyDynamicConfig, err := legacyOnramp.GetDynamicConfig(opts)
		if err != nil {
			return evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{}, err
		}
		return evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{
			Router:                          legacyDynamicConfig.Router,
			MaxTokensLength:                 legacyDynamicConfig.MaxTokensLength,
			DestGasOverhead:                 0,
			DestGasPerPayloadByte:           0,
			DestDataAvailabilityOverheadGas: 0,
			DestGasPerDataAvailabilityByte:  0,
			DestDataAvailabilityMultiplier:  0,
			PriceRegistry:                   legacyDynamicConfig.PriceRegistry,
			MaxDataSize:                     legacyDynamicConfig.MaxDataSize,
			MaxGasLimit:                     uint32(legacyDynamicConfig.MaxGasLimit),
		}, nil
	case "1.1.0":
		legacyOnramp, err := evm_2_evm_onramp_1_1_0.NewEVM2EVMOnRamp(onRamp.Address(), client)
		if err != nil {
			return evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{}, err
		}
		legacyDynamicConfig, err := legacyOnramp.GetDynamicConfig(opts)
		if err != nil {
			return evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{}, err
		}
		return evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{
			Router:                          legacyDynamicConfig.Router,
			MaxTokensLength:                 legacyDynamicConfig.MaxTokensLength,
			DestGasOverhead:                 legacyDynamicConfig.DestGasOverhead,
			DestGasPerPayloadByte:           legacyDynamicConfig.DestGasPerPayloadByte,
			DestDataAvailabilityOverheadGas: 0,
			DestGasPerDataAvailabilityByte:  0,
			DestDataAvailabilityMultiplier:  0,
			PriceRegistry:                   legacyDynamicConfig.PriceRegistry,
			MaxDataSize:                     legacyDynamicConfig.MaxDataSize,
			MaxGasLimit:                     uint32(legacyDynamicConfig.MaxGasLimit),
		}, nil
	case "1.2.0":
		return onRamp.GetDynamicConfig(opts)
	default:
		return evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{}, errors.Errorf("Invalid onramp version: %s", version)
	}
}

func LoadOffRamp(offRampAddress common.Address, pluginName string, client client.Client) (evm_2_evm_offramp.EVM2EVMOffRampInterface, semver.Version, error) {
	version, err := ccipconfig.VerifyTypeAndVersion(offRampAddress, client, ccipconfig.EVM2EVMOffRamp)
	if err != nil {
		return nil, semver.Version{}, errors.Wrap(err, "Invalid offRamp contract")
	}

	offRamp, err := observability.NewObservedEvm2EvmOffRamp(offRampAddress, pluginName, client)
	return offRamp, version, err
}

func LoadCommitStore(commitStoreAddress common.Address, pluginName string, client client.Client) (commit_store.CommitStoreInterface, semver.Version, error) {
	version, err := ccipconfig.VerifyTypeAndVersion(commitStoreAddress, client, ccipconfig.CommitStore)
	if err != nil {
		return nil, semver.Version{}, errors.Wrap(err, "Invalid commitStore contract")
	}

	commitStore, err := observability.NewObservedCommitStore(commitStoreAddress, pluginName, client)
	return commitStore, version, err
}

func DecodeCommitStoreOffchainConfig(version semver.Version, offchainConfig []byte) (ccipconfig.CommitOffchainConfig, error) {
	switch version.String() {
	case "1.0.0", "1.1.0":
		offchainConfigV1, err := ccipconfig.DecodeOffchainConfig[ccipconfig.CommitOffchainConfigV1](offchainConfig)
		if err != nil {
			return ccipconfig.CommitOffchainConfig{}, err
		}

		return ccipconfig.CommitOffchainConfig{
			SourceFinalityDepth:      offchainConfigV1.SourceFinalityDepth,
			DestFinalityDepth:        offchainConfigV1.DestFinalityDepth,
			GasPriceHeartBeat:        offchainConfigV1.FeeUpdateHeartBeat,
			DAGasPriceDeviationPPB:   offchainConfigV1.FeeUpdateDeviationPPB,
			ExecGasPriceDeviationPPB: offchainConfigV1.FeeUpdateDeviationPPB,
			TokenPriceHeartBeat:      offchainConfigV1.FeeUpdateHeartBeat,
			TokenPriceDeviationPPB:   offchainConfigV1.FeeUpdateDeviationPPB,
			MaxGasPrice:              offchainConfigV1.MaxGasPrice,
			InflightCacheExpiry:      offchainConfigV1.InflightCacheExpiry,
		}, nil
	case "1.2.0":
		offchainConfig, err := ccipconfig.DecodeOffchainConfig[ccipconfig.CommitOffchainConfig](offchainConfig)
		return offchainConfig, err
	default:
		return ccipconfig.CommitOffchainConfig{}, errors.Errorf("Invalid commitStore version: %s", version)
	}
}
