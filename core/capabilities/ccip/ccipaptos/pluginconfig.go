package ccipaptos

import (
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ocrimpls"
)

// initializePluginConfig returns a PluginConfig for Aptos chains.
func initializePluginConfig(lggr logger.Logger, extraDataCodec ccipcommon.ExtraDataCodec) ccipcommon.PluginConfig {
	return ccipcommon.PluginConfig{
		CommitPluginCodec:          NewCommitPluginCodecV1(),
		ExecutePluginCodec:         NewExecutePluginCodecV1(extraDataCodec),
		MessageHasher:              NewMessageHasherV1(logger.Sugared(lggr).Named(chainsel.FamilyAptos).Named("MessageHasherV1"), extraDataCodec),
		TokenDataEncoder:           NewAptosTokenDataEncoder(),
		GasEstimateProvider:        NewGasEstimateProvider(),
		RMNCrypto:                  nil,
		ContractTransmitterFactory: ocrimpls.NewAptosContractTransmitterFactory(extraDataCodec),
		ChainAccessorFactory:       AptosChainAccessorFactory{},
		ChainRW:                    ChainCWProvider{},
		ExtraDataCodec:             ExtraDataDecoder{},
		AddressCodec:               AddressCodec{},
	}
}

func init() {
	// Register the Aptos plugin config factory
	ccipcommon.RegisterPluginConfig(chainsel.FamilyAptos, initializePluginConfig)
}
