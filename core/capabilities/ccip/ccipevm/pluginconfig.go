package ccipevm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipsolana"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/configs/evm"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	evmrelaytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type PluginConfig struct {
	extraDataCodec cctypes.ExtraDataCodec
}

func NewPluginConfig(extraDataCodec cctypes.ExtraDataCodec) *PluginConfig {
	return &PluginConfig{
		extraDataCodec: extraDataCodec,
	}
}

func (p PluginConfig) InitializePluginConfig() ccipcommon.PluginConfig {
	extraDataCodec := ccipcommon.NewExtraDataCodec(ExtraDataDecoder{}, ccipsolana.ExtraDataDecoder{})
	return ccipcommon.PluginConfig{
		CommitPluginCodec:  NewCommitPluginCodecV1(),
		ExecutePluginCodec: NewExecutePluginCodecV1(extraDataCodec),
		MessageHasher: func(lggr logger.Logger) cciptypes.MessageHasher {
			return NewMessageHasherV1(lggr, extraDataCodec)
		},
		TokenDataEncoder:     NewEVMTokenDataEncoder(),
		GasEstimateProvider:  NewGasEstimateProvider(),
		RMNCrypto:            func(lggr logger.Logger) cciptypes.RMNCrypto { return NewEVMRMNCrypto(lggr) },
		GetChainReaderConfig: getEVMChainReaderConfig,
		GetChainWriter:       getEVMChainWriter,
	}
}

func getEVMChainReaderConfig(
	lggr logger.Logger,
	chainID string,
	destChainID string,
	homeChainID string,
	ofc ccipcommon.OffChainConfig,
	chainSelector cciptypes.ChainSelector,
) ([]byte, error) {
	var chainReaderConfig evmrelaytypes.ChainReaderConfig
	if chainID == destChainID {
		chainReaderConfig = evmconfig.DestReaderConfig
	} else {
		chainReaderConfig = evmconfig.SourceReaderConfig
	}

	if !ofc.CommitEmpty() && ofc.Commit().PriceFeedChainSelector == chainSelector {
		lggr.Debugw("Adding feed reader config", "chainID", chainID)
		chainReaderConfig = evmconfig.MergeReaderConfigs(chainReaderConfig, evmconfig.FeedReaderConfig)
	}

	if isUSDCEnabled(ofc) {
		lggr.Debugw("Adding USDC reader config", "chainID", chainID)
		chainReaderConfig = evmconfig.MergeReaderConfigs(chainReaderConfig, evmconfig.USDCReaderConfig)
	}

	if chainID == homeChainID {
		lggr.Debugw("Adding home chain reader config", "chainID", chainID)
		chainReaderConfig = evmconfig.MergeReaderConfigs(chainReaderConfig, evmconfig.HomeChainReaderConfigRaw)
	}

	marshaledConfig, err := json.Marshal(chainReaderConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chain reader config: %w", err)
	}

	return marshaledConfig, nil
}

func getEVMChainWriter(
	ctx context.Context,
	chainID string,
	relayer loop.Relayer,
	transmitters map[types.RelayID][]string,
	execBatchGasLimit uint64,
	chainFamily string,
	offrampProgramAddress []byte,
) (types.ContractWriter, error) {
	var fromAddress common.Address
	transmitter, ok := transmitters[types.NewRelayID(chainFamily, chainID)]
	if ok {
		fromAddress = common.HexToAddress(transmitter[0])
	}

	evmConfig, err := evmconfig.ChainWriterConfigRaw(
		fromAddress,
		defaultCommitGasLimit,
		execBatchGasLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to create EVM chain writer config: %w", err)
	}

	chainWriterConfig, err := json.Marshal(evmConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal EVM chain writer config: %w", err)
	}

	cw, err := relayer.NewContractWriter(ctx, chainWriterConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create chain writer for chain %s: %w", chainID, err)
	}

	return cw, nil
}

func isUSDCEnabled(ofc ccipcommon.OffChainConfig) bool {
	if ofc.ExecEmpty() {
		return false
	}

	return ofc.Exec().IsUSDCEnabled()
}
