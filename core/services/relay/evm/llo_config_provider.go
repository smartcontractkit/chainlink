package evm

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

func newLLOConfigProvider(ctx context.Context, lggr logger.Logger, chain legacyevm.Chain, opts *types.RelayOpts) (commontypes.ConfigProvider, error) {
	// ContractID should be the address of the capabilities registry contract
	if !common.IsHexAddress(opts.ContractID) {
		return nil, errors.New("invalid contractID, expected hex address")
	}

	capabilitiesRegistryAddr := common.HexToAddress(opts.ContractID)
	digester := evmutil.EVMOffchainConfigDigester{
		ChainID:         chain.Config().EVM().ChainID().Uint64(),
		ContractAddress: capabilitiesRegistryAddr,
	}

	var cp types.ConfigPoller

	relayConfig, err := opts.RelayConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get relay config: %w", err)
	}
	cp, err = NewLLOConfigPoller(
		ctx,
		lggr,
		chain.Client(),
		chain.LogPoller(),
		capabilitiesRegistryAddr,
	)
	if err != nil {
		return nil, err
	}

	return newConfigWatcher(lggr, capabilitiesRegistryAddr, digester, cp, chain, relayConfig.FromBlock, opts.New), nil
}

// type configWrapper struct {
//     ocr2types.ConfigTracker
// }

// func (c *configWrapper) Start()
