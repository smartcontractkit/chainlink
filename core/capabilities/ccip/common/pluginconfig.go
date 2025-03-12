package common

import (
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipsolana"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// CreatePluginConfig creates a PluginConfig for the given chain family.
func CreatePluginConfig(chainFamily string) (cctypes.PluginConfig, error) {
	extraDataCodec := cctypes.NewExtraDataCodec(
		cctypes.NewExtraDataCodecParams(
			ccipevm.ExtraDataDecoder{},
			ccipsolana.ExtraDataDecoder{},
		),
	)
	switch chainFamily {
	case chainsel.FamilyEVM:
		return cctypes.PluginConfig{
			CommitPluginCodec:  ccipevm.NewCommitPluginCodecV1(),
			ExecutePluginCodec: ccipevm.NewExecutePluginCodecV1(extraDataCodec),
			MessageHasher: func(lggr logger.Logger) cciptypes.MessageHasher {
				return ccipevm.NewMessageHasherV1(lggr, extraDataCodec)
			},
			TokenDataEncoder:     ccipevm.NewEVMTokenDataEncoder(),
			GasEstimateProvider:  ccipevm.NewGasEstimateProvider(),
			RMNCrypto:            func(lggr logger.Logger) cciptypes.RMNCrypto { return ccipevm.NewEVMRMNCrypto(lggr) },
			GetChainReaderConfig: ccipevm.GetEVMChainReaderConfig,
			GetChainWriter:       ccipevm.GetEVMChainWriter,
		}, nil
	case chainsel.FamilySolana:
		return cctypes.PluginConfig{
			CommitPluginCodec:  ccipsolana.NewCommitPluginCodecV1(),
			ExecutePluginCodec: ccipsolana.NewExecutePluginCodecV1(extraDataCodec),
			MessageHasher: func(lggr logger.Logger) cciptypes.MessageHasher {
				return ccipsolana.NewMessageHasherV1(lggr, extraDataCodec)
			},
			TokenDataEncoder:     ccipsolana.NewSolanaTokenDataEncoder(),
			GasEstimateProvider:  ccipsolana.NewGasEstimateProvider(),
			RMNCrypto:            func(lggr logger.Logger) cciptypes.RMNCrypto { return nil },
			GetChainReaderConfig: ccipsolana.GetSolanaChainReaderConfig,
			GetChainWriter:       ccipsolana.GetSolanaChainWriter,
			PriceOnlyCommitFn:    consts.MethodCommitPriceOnly,
		}, nil
	}

	return cctypes.PluginConfig{}, fmt.Errorf("unsupported chain family: %s", chainFamily)
}
