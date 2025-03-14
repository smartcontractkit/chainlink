package ccipevm

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/configs/evm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	evmrelaytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type PluginConfig struct {
	extraDataCodec ccipcommon.ExtraDataCodec
}

func NewPluginConfig(extraDataCodec ccipcommon.ExtraDataCodec) *PluginConfig {
	return &PluginConfig{
		extraDataCodec: extraDataCodec,
	}
}

func (p PluginConfig) InitializePluginConfig() ccipcommon.PluginConfig {
	return ccipcommon.PluginConfig{
		CommitPluginCodec:  NewCommitPluginCodecV1(),
		ExecutePluginCodec: NewExecutePluginCodecV1(p.extraDataCodec),
		MessageHasher: func(lggr logger.Logger) cciptypes.MessageHasher {
			return NewMessageHasherV1(lggr, p.extraDataCodec)
		},
		TokenDataEncoder:     NewEVMTokenDataEncoder(),
		GasEstimateProvider:  NewGasEstimateProvider(),
		RMNCrypto:            func(lggr logger.Logger) cciptypes.RMNCrypto { return NewEVMRMNCrypto(lggr) },
		GetChainReaderWriter: GetCRCW{},
	}
}

type GetCRCW struct{}

func (g GetCRCW) GetChainReader(params ccipcommon.GetChainReaderParams) (types.ContractReader, error) {
	var chainReaderConfig evmrelaytypes.ChainReaderConfig
	if params.ChainID == params.DestChainID {
		chainReaderConfig = evmconfig.DestReaderConfig
	} else {
		chainReaderConfig = evmconfig.SourceReaderConfig
	}

	if !params.Ofc.CommitEmpty() && params.Ofc.Commit().PriceFeedChainSelector == params.ChainSelector {
		params.Lggr.Debugw("Adding feed reader config", "chainID", params.ChainID)
		chainReaderConfig = evmconfig.MergeReaderConfigs(chainReaderConfig, evmconfig.FeedReaderConfig)
	}

	if isUSDCEnabled(params.Ofc) {
		params.Lggr.Debugw("Adding USDC reader config", "chainID", params.ChainID)
		chainReaderConfig = evmconfig.MergeReaderConfigs(chainReaderConfig, evmconfig.USDCReaderConfig)
	}

	if params.ChainID == params.HomeChainID {
		params.Lggr.Debugw("Adding home chain reader config", "chainID", params.ChainID)
		chainReaderConfig = evmconfig.MergeReaderConfigs(chainReaderConfig, evmconfig.HomeChainReaderConfigRaw)
	}

	marshaledConfig, err := json.Marshal(chainReaderConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chain reader config: %w", err)
	}

	cr, err := params.Relayer.NewContractReader(params.Ctx, marshaledConfig)
	if err != nil {
		return nil, err
	}

	return cr, nil
}

func (g GetCRCW) GetChainWriter(params ccipcommon.GetChainWriterParams) (types.ContractWriter, error) {
	var fromAddress common.Address
	transmitter, ok := params.Transmitters[types.NewRelayID(params.ChainFamily, params.ChainID)]
	if ok {
		fromAddress = common.HexToAddress(transmitter[0])
	}

	evmConfig, err := evmconfig.ChainWriterConfigRaw(
		fromAddress,
		defaultCommitGasLimit,
		params.ExecBatchGasLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to create EVM chain writer config: %w", err)
	}

	chainWriterConfig, err := json.Marshal(evmConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal EVM chain writer config: %w", err)
	}

	cw, err := params.Relayer.NewContractWriter(params.Ctx, chainWriterConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create chain writer for chain %s: %w", params.ChainID, err)
	}

	return cw, nil
}

func isUSDCEnabled(ofc ccipcommon.OffChainConfig) bool {
	if ofc.ExecEmpty() {
		return false
	}

	return ofc.Exec().IsUSDCEnabled()
}
