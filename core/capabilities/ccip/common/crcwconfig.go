package common

import (
	"context"
	"fmt"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// ChainReaderProviderOpts is a struct that contains the parameters for GetChainReader.
type ChainReaderProviderOpts struct {
	Lggr          logger.Logger
	Relayer       loop.Relayer
	ChainID       string
	DestChainID   string
	HomeChainID   string
	Ofc           OffChainConfig
	ChainSelector cciptypes.ChainSelector
	ChainFamily   string
}

// ChainWriterProviderOpts is a struct that contains the parameters for GetChainWriter.
type ChainWriterProviderOpts struct {
	ChainID               string
	Relayer               loop.Relayer
	Transmitters          map[types.RelayID][]string
	ExecBatchGasLimit     uint64
	ChainFamily           string
	OfframpProgramAddress []byte
}

// ChainRWProvider is an interface that defines the methods to get a ContractReader and a ContractWriter.
type ChainRWProvider interface {
	GetChainReader(ctx context.Context, params ChainReaderProviderOpts) (types.ContractReader, error)
	GetChainWriter(ctx context.Context, params ChainWriterProviderOpts) (types.ContractWriter, error)
}

// MultiChainRW is a struct that implements the ChainRWProvider interface for all chains.
type MultiChainRW struct {
	cwProviderMap map[string]ChainRWProvider
}

// NewCRCW is a constructor for MultiChainRW.
func NewCRCW(cwProviderMap map[string]ChainRWProvider) MultiChainRW {
	return MultiChainRW{
		cwProviderMap: cwProviderMap,
	}
}

// GetChainReader returns a new ContractReader base on relay chain family.
func (c *MultiChainRW) GetChainReader(ctx context.Context, params ChainReaderProviderOpts) (types.ContractReader, error) {
	provider, exist := c.cwProviderMap[params.ChainFamily]
	if !exist {
		return nil, fmt.Errorf("unsupported chain family %s", params.ChainFamily)
	}

	return provider.GetChainReader(ctx, params)
}

// GetChainWriter returns a new ContractWriter based on relay chain family.
func (c *MultiChainRW) GetChainWriter(ctx context.Context, params ChainWriterProviderOpts) (types.ContractWriter, error) {
	provider, exist := c.cwProviderMap[params.ChainFamily]
	if !exist {
		return nil, fmt.Errorf("unsupported chain family %s", params.ChainFamily)
	}

	return provider.GetChainWriter(ctx, params)
}
