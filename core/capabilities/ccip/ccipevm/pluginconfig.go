package ccipevm

import (
	chainsel "github.com/smartcontractkit/chain-selectors"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ocrimpls"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const defaultCommitGasLimit = 500_000

// InitializePluginConfig returns a PluginConfig for EVM chains.
func InitializePluginConfig(lggr logger.Logger, extraDataCodec ccipcommon.ExtraDataCodec) ccipcommon.PluginConfig {
	return ccipcommon.PluginConfig{
		CommitPluginCodec:          NewCommitPluginCodecV1(),
		ExecutePluginCodec:         NewExecutePluginCodecV1(extraDataCodec),
		MessageHasher:              NewMessageHasherV1(lggr.Named(chainsel.FamilyEVM).Named("MessageHasherV1"), extraDataCodec),
		TokenDataEncoder:           NewEVMTokenDataEncoder(),
		GasEstimateProvider:        NewGasEstimateProvider(extraDataCodec),
		RMNCrypto:                  NewEVMRMNCrypto(lggr.Named(chainsel.FamilyEVM).Named("RMNCrypto")),
		ContractTransmitterFactory: ocrimpls.NewEVMContractTransmitterFactory(extraDataCodec),
		ChainRW:                    ChainCWProvider{},
		ExtraDataCodec:             ExtraDataDecoder{},
		AddressCodec:               AddressCodec{},
	}
}

func init() {
	// Register the EVM plugin config factory
	ccipcommon.RegisterPluginConfig(chainsel.FamilyEVM, InitializePluginConfig)
}
