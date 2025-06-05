package evm

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

func newStandardConfigProvider(ctx context.Context, lggr logger.Logger, chain legacyevm.Chain, opts *types.RelayOpts) (*configWatcher, error) {
	if !common.IsHexAddress(opts.ContractID) {
		return nil, errors.New("invalid contractID, expected hex address")
	}

	aggregatorAddress := common.HexToAddress(opts.ContractID)
	offchainConfigDigester := evmutil.EVMOffchainConfigDigester{
		ChainID:         chain.Config().EVM().ChainID().Uint64(),
		ContractAddress: aggregatorAddress,
	}
	return newContractConfigProvider(ctx, lggr, chain, opts, aggregatorAddress, OCR2AggregatorLogDecoder, offchainConfigDigester)
}

func newContractConfigProvider(ctx context.Context, lggr logger.Logger, chain legacyevm.Chain, opts *types.RelayOpts, aggregatorAddress common.Address, ld LogDecoder, digester ocrtypes.OffchainConfigDigester) (*configWatcher, error) {
	var cp types.ConfigPoller

	relayConfig, err := opts.RelayConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get relay config: %w", err)
	}
	cp, err = NewConfigPoller(
		ctx,
		lggr,
		CPConfig{
			chain.Client(),
			chain.LogPoller(),
			aggregatorAddress,
			relayConfig.ConfigContractAddress,
			ld,
		},
	)
	if err != nil {
		return nil, err
	}

	return newConfigWatcher(lggr, aggregatorAddress, digester, cp, chain, relayConfig.FromBlock, opts.New), nil
}
