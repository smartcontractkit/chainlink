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

func newSecureMintConfigProvider(ctx context.Context, lggr logger.Logger, chain legacyevm.Chain, opts *types.RelayOpts) (*configWatcher, error) {
	if !common.IsHexAddress(opts.ContractID) {
		return nil, errors.New("invalid contractID, expected hex address")
	}

	configStoreAddress := common.HexToAddress(opts.ContractID)
	offchainConfigDigester := evmutil.EVMOffchainConfigDigester{
		ChainID:         chain.Config().EVM().ChainID().Uint64(),
		ContractAddress: configStoreAddress,
	}
	lggr.Infof("TRACE - Creating SecureMintConfigProvider with contract address: %s", configStoreAddress.Hex())

	// Create a log decoder for OCRConfigurationStoreEVMSimple contract
	logDecoder, err := newOCRConfigurationStoreEVMSimpleLogDecoder(chain, configStoreAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create log decoder: %w", err)
	}

	// Use the new ConfigPollerEVMSimple implementation with logpoller
	cp, err := NewConfigPollerEVMSimple(ctx, lggr, ConfigPollerEVMSimpleConfig{
		LogPoller:  chain.LogPoller(),
		Address:    configStoreAddress,
		LogDecoder: logDecoder,
		Client:     chain.Client(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create ConfigPollerEVMSimple: %w", err)
	}

	relayConfig, err := opts.RelayConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get relay config: %w", err)
	}

	return newConfigWatcher(lggr, configStoreAddress, offchainConfigDigester, cp, chain, relayConfig.FromBlock, opts.New), nil
}
