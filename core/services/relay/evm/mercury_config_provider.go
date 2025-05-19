package evm

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

func newMercuryConfigProvider(ctx context.Context, lggr logger.Logger, chain legacyevm.Chain, opts *types.RelayOpts) (commontypes.ConfigProvider, error) {
	if !common.IsHexAddress(opts.ContractID) {
		return nil, errors.New("invalid contractID, expected hex address")
	}

	aggregatorAddress := common.HexToAddress(opts.ContractID)

	relayConfig, err := opts.RelayConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get relay config: %w", err)
	}
	if relayConfig.FeedID == nil {
		return nil, errors.New("feed ID is required for tracking config on mercury contracts")
	}
	cp, err := mercury.NewConfigPoller(
		ctx,
		logger.Named(lggr, relayConfig.FeedID.String()),
		chain.LogPoller(),
		aggregatorAddress,
		*relayConfig.FeedID,
		// TODO: Does mercury need to support config contract? DF-19182
	)
	if err != nil {
		return nil, err
	}

	offchainConfigDigester := mercury.NewOffchainConfigDigester(*relayConfig.FeedID, chain.Config().EVM().ChainID(), aggregatorAddress, ocrtypes.ConfigDigestPrefixMercuryV02)
	return newConfigWatcher(lggr, aggregatorAddress, offchainConfigDigester, cp, chain, relayConfig.FromBlock, opts.New), nil
}
