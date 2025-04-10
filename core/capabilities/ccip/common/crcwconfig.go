package common

import (
	"context"
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// GetChainReaderParams is a struct that contains the parameters for GetChainReader.
type GetChainReaderParams struct {
	Lggr             logger.Logger
	Relayer          loop.Relayer
	ChainID          string
	DestChainID      string
	HomeChainID      string
	Ofc              OffChainConfig
	ChainSelector    cciptypes.ChainSelector
	RelayChainFamily string
}

// GetChainWriterParams is a struct that contains the parameters for GetChainWriter.
type GetChainWriterParams struct {
	ChainID               string
	Relayer               loop.Relayer
	Transmitters          map[types.RelayID][]string
	ExecBatchGasLimit     uint64
	ChainFamily           string
	OfframpProgramAddress []byte
}

// GetChainReaderWriter is an interface that defines the methods to get a ContractReader and a ContractWriter.
type GetChainReaderWriter interface {
	GetChainReader(ctx context.Context, params GetChainReaderParams) (types.ContractReader, error)
	GetChainWriter(ctx context.Context, params GetChainWriterParams) (types.ContractWriter, error)
}

// CRCW is a struct that implements the GetChainReaderWriter interface for all chains.
type CRCW struct {
	EVMCRCW    GetChainReaderWriter
	SolanaCRCW GetChainReaderWriter
	lggr       logger.Logger
}

// NewCRCW is a constructor for CRCW.
func NewCRCW(evmCRCW, solanaCRCW GetChainReaderWriter, lggr logger.Logger) *CRCW {
	return &CRCW{
		EVMCRCW:    evmCRCW,
		SolanaCRCW: solanaCRCW,
		lggr:       lggr,
	}
}

// GetChainReader returns a new ContractReader base on relay chain family.
func (c *CRCW) GetChainReader(ctx context.Context, params GetChainReaderParams) (types.ContractReader, error) {
	switch params.RelayChainFamily {
	case chainsel.FamilyEVM:
		return c.EVMCRCW.GetChainReader(ctx, params)
	case chainsel.FamilySolana:
		return c.SolanaCRCW.GetChainReader(ctx, params)
	default:
		return nil, fmt.Errorf("unsupported chain family %s", params.RelayChainFamily)
	}
}

// GetChainWriter returns a new ContractWriter based on relay chain family.
func (c *CRCW) GetChainWriter(ctx context.Context, params GetChainWriterParams) (types.ContractWriter, error) {
	switch params.ChainFamily {
	case chainsel.FamilyEVM:
		return c.EVMCRCW.GetChainWriter(ctx, params)
	case chainsel.FamilySolana:
		return c.SolanaCRCW.GetChainWriter(ctx, params)
	default:
		return nil, fmt.Errorf("unsupported chain family %s", params.ChainFamily)
	}
}
