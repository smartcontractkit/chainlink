package ccipsolana

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gagliardetto/solana-go"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	solanaconfig "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/configs/solana"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/oraclecreator"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func CreatePluginConfig() oraclecreator.Plugin {
	return oraclecreator.Plugin{
		CommitPluginCodec:    NewCommitPluginCodecV1(),
		ExecutePluginCodec:   NewExecutePluginCodecV1(),
		ExtraArgsCodec:       ccipcommon.NewExtraDataCodec(),
		MessageHasher:        func(lggr logger.Logger) cciptypes.MessageHasher { return NewMessageHasherV1(lggr) },
		TokenDataEncoder:     NewSolanaTokenDataEncoder(),
		GasEstimateProvider:  NewGasEstimateProvider(),
		RMNCrypto:            func(lggr logger.Logger) cciptypes.RMNCrypto { return nil },
		AddressToString:      func(addr []byte, checkSum bool) string { return solana.PublicKeyFromBytes(addr).String() },
		GetChainReaderConfig: getSolanaChainReaderConfig,
		GetChainWriter:       getSolanaChainWriter,
	}
}

func getSolanaChainReaderConfig(lggr logger.Logger,
	chainID string,
	destChainID string,
	homeChainID string,
	ofc cctypes.OffChainConfig,
	chainSelector cciptypes.ChainSelector,
) ([]byte, error) {
	// TODO update chain reader config in contract_reader.go
	var cfg config.ContractReader
	if chainID == destChainID {
		cfg = solanaconfig.DestReaderConfig
	} else {
		cfg = solanaconfig.SourceReaderConfig
	}

	marshaledConfig, err := json.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chain reader config: %w", err)
	}

	return marshaledConfig, nil
}

func getSolanaChainWriter(
	ctx context.Context,
	chainID string,
	relayer loop.Relayer,
	transmitters map[types.RelayID][]string,
	execBatchGasLimit uint64,
	chainFamily string,
	offrampProgramAddress []byte,
	destChainSelector uint64,
) (types.ContractWriter, error) {
	transmitter := transmitters[types.NewRelayID(chainFamily, chainID)]
	offrampAddress := solana.PublicKeyFromBytes(offrampProgramAddress)
	solConfig, err := solanaconfig.GetSolanaChainWriterConfig(offrampAddress.String(), transmitter[0], destChainSelector)
	if err == nil {
		return nil, fmt.Errorf("failed to get Solana chain writer config: %w", err)
	}
	chainWriterConfig, err := json.Marshal(solConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Solana chain writer config: %w", err)
	}

	cw, err := relayer.NewContractWriter(ctx, chainWriterConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create chain writer for chain %s: %w", chainID, err)
	}

	return cw, nil
}
