package ccipaptos

import (
	"context"
	"encoding/json"
	"fmt"

	aptosloop "github.com/smartcontractkit/chainlink-aptos/relayer/chainreader/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	aptosconfig "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/configs/aptos"
)

// ChainCWProvider is a struct that implements the ChainRWProvider interface for Aptos chains.
type ChainCWProvider struct{}

// GetChainReader returns a new ContractReader for Aptos chains.
func (g ChainCWProvider) GetChainReader(ctx context.Context, params ccipcommon.ChainReaderProviderOpts) (types.ContractReader, error) {
	cfg, err := aptosconfig.GetChainReaderConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get Aptos chain reader config: %w", err)
	}
	marshaledConfig, err := json.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Aptos chain reader config: %w", err)
	}

	cr, err := params.Relayer.NewContractReader(ctx, marshaledConfig)
	if err != nil {
		return nil, err
	}

	cr = aptosloop.NewLoopChainReader(params.Lggr, cr)

	return cr, nil
}

// GetChainWriter returns a new ContractWriter for Aptos chains.
func (g ChainCWProvider) GetChainWriter(ctx context.Context, params ccipcommon.ChainWriterProviderOpts) (types.ContractWriter, error) {
	transmitter := params.Transmitters[types.NewRelayID(params.ChainFamily, params.ChainID)]
	cfg, err := aptosconfig.GetChainWriterConfig(transmitter[0])
	if err != nil {
		return nil, fmt.Errorf("failed to get Aptos chain writer config: %w", err)
	}
	chainWriterConfig, err := json.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Aptos chain writer config: %w", err)
	}

	cw, err := params.Relayer.NewContractWriter(ctx, chainWriterConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create chain writer for chain %s: %w", params.ChainID, err)
	}

	return cw, nil
}
