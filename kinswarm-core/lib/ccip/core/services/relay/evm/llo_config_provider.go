package evm

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

func DonIDToBytes32(donID uint32) [32]byte {
	var b [32]byte
	copy(b[:], common.LeftPadBytes(big.NewInt(int64(donID)).Bytes(), 32))
	return b
}

func newLLOConfigProvider(ctx context.Context, lggr logger.Logger, chain legacyevm.Chain, opts *types.RelayOpts) (*configWatcher, error) {
	if !common.IsHexAddress(opts.ContractID) {
		return nil, errors.New("invalid contractID, expected hex address")
	}

	configuratorAddress := common.HexToAddress(opts.ContractID)

	relayConfig, err := opts.RelayConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get relay config: %w", err)
	}
	if relayConfig.LLODONID == 0 {
		return nil, errors.New("donID must be specified in relayConfig for LLO jobs")
	}
	donIDPadded := DonIDToBytes32(relayConfig.LLODONID)
	cp, err := mercury.NewConfigPoller(
		ctx,
		lggr.Named(fmt.Sprintf("LLO-%d", relayConfig.LLODONID)),
		chain.LogPoller(),
		configuratorAddress,
		donIDPadded,
		// TODO: Does LLO need to support config contract? MERC-1827
	)
	if err != nil {
		return nil, err
	}

	configDigester := mercury.NewOffchainConfigDigester(donIDPadded, chain.Config().EVM().ChainID(), configuratorAddress, ocrtypes.ConfigDigestPrefixLLO)
	return newConfigWatcher(lggr, configuratorAddress, configDigester, cp, chain, relayConfig.FromBlock, opts.New), nil
}
