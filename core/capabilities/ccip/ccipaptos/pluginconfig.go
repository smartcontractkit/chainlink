package ccipaptos

import (
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipsui"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ocrimpls"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
)

// initializePluginConfig returns a PluginConfig for Aptos chains.
func initializePluginConfigFunc(chainselFamily string) ccipcommon.InitFunction {
	return func(lggr logger.Logger, extraDataCodec ccipocr3.ExtraDataCodecBundle) ccipcommon.PluginConfig {
		var cwProvider ccipcommon.ChainRWProvider
		var transmitterFactory types.ContractTransmitterFactory
		var msgHasher ccipocr3.MessageHasher
		var executeCodec ccipocr3.ExecutePluginCodec

		if chainselFamily == chainsel.FamilyAptos {
			cwProvider = ChainCWProvider{}
			transmitterFactory = ocrimpls.NewAptosContractTransmitterFactory(extraDataCodec)
			msgHasher = NewMessageHasherV1(logger.Sugared(lggr).Named(chainselFamily).Named("MessageHasherV1"), extraDataCodec)
			executeCodec = NewExecutePluginCodecV1(extraDataCodec)
		} else {
			cwProvider = ccipsui.ChainCWProvider{}
			transmitterFactory = ocrimpls.NewSuiContractTransmitterFactory(extraDataCodec)
			msgHasher = ccipsui.NewMessageHasherV1(logger.Sugared(lggr).Named(chainselFamily).Named("MessageHasherV1"), extraDataCodec)
			executeCodec = ccipsui.NewExecutePluginCodecV1(extraDataCodec)
		}

		return ccipcommon.PluginConfig{
			CommitPluginCodec:          NewCommitPluginCodecV1(),
			ExecutePluginCodec:         executeCodec,
			MessageHasher:              msgHasher,
			TokenDataEncoder:           NewAptosTokenDataEncoder(),
			GasEstimateProvider:        NewGasEstimateProvider(),
			RMNCrypto:                  nil,
			ContractTransmitterFactory: transmitterFactory,
			ChainRW:                    cwProvider,
			ExtraDataCodec:             ExtraDataDecoder{},
			AddressCodec:               AddressCodec{},
		}
	}
}

func init() {
	// Register the Aptos and Sui plugin config factory
	ccipcommon.RegisterPluginConfig(chainsel.FamilyAptos, initializePluginConfigFunc(chainsel.FamilyAptos))
	ccipcommon.RegisterPluginConfig(chainsel.FamilySui, initializePluginConfigFunc(chainsel.FamilySui))
}
