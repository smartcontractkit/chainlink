package ccipsui

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	suiloop "github.com/smartcontractkit/chainlink-sui/relayer/chainreader/loop"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	suiconfig "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/configs/sui"
)

// ChainCWProvider is a struct that implements the ChainRWProvider interface for EVM chains.
type ChainCWProvider struct{}

// GetChainReader returns a new ContractReader for EVM chains.
func (g ChainCWProvider) GetChainReader(ctx context.Context, params ccipcommon.ChainReaderProviderOpts) (types.ContractReader, error) {
	transmitter := params.Transmitters[types.NewRelayID(params.ChainFamily, params.ChainID)]

	cfg, err := suiconfig.GetChainReaderConfig(transmitter[0])
	if err != nil {
		return nil, fmt.Errorf("failed to get SUI config: %w", err)
	}

	marshaledConfig, err := json.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal SUI chain reader config: %w", err)
	}

	cr, err := params.Relayer.NewContractReader(ctx, marshaledConfig)
	if err != nil {
		return nil, err
	}

	cr = suiloop.NewLoopChainReader(params.Lggr, cr)

	return cr, nil
}

// GetChainWriter returns a new ContractWriter for EVM chains.
func (g ChainCWProvider) GetChainWriter(ctx context.Context, params ccipcommon.ChainWriterProviderOpts) (types.ContractWriter, error) {
	transmitter := params.Transmitters[types.NewRelayID(params.ChainFamily, params.ChainID)]
	cfg, err := suiconfig.GetChainWriterConfig(transmitter[0])
	if err != nil {
		return nil, fmt.Errorf("failed to get Sui chain writer config: %w", err)
	}
	chainWriterConfig, err := json.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Sui chain writer config: %w", err)
	}

	cw, err := params.Relayer.NewContractWriter(ctx, chainWriterConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create chain writer for chain %s: %w", params.ChainID, err)
	}
	return cw, nil
}
