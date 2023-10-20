package evm

import (
	"context"
	"fmt"

	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

// constructor for ChainReader, returns nil if there is any error
func newChainReader(lggr logger.Logger, chain evm.Chain, ropts *types.RelayOpts) *chainReader {
	relayConfig, err := ropts.RelayConfig()
	if err != nil {
		lggr.Errorf("Failed parsing RelayConfig: %w", err.Error())
		return nil
	}

	if relayConfig.ChainReader == nil {
		return nil
	}

	if err = validateChainReaderConfig(*relayConfig.ChainReader); err != nil {
		lggr.Errorf("Invalid ChainReader configuration: %w", err.Error())
		return nil
	}

	return &chainReader{lggr, chain.LogPoller()}
}

func validateChainReaderConfig(cfg types.ChainReaderConfig) error {
	// Validate config (check ABI from job spec against imported gethwrappers, etc.)
	return nil
}

func (cr *chainReader) initialize() error {
	// Initialize chain reader, start cache polling loop, etc.
	return nil
}

type ChainReaderService interface {
	services.ServiceCtx
	relaytypes.ChainReader
}

type chainReader struct {
	lggr logger.Logger
	lp   logpoller.LogPoller
}

// chainReader constructor
func NewChainReaderService(lggr logger.Logger, lp logpoller.LogPoller) (*chainReader, error) {
	return &chainReader{lggr, lp}, nil
}

func (cr *chainReader) GetLatestValue(ctx context.Context, bc relaytypes.BoundContract, method string, params any, returnVal any) error {

	// TODO: implement GetLatestValue

	return fmt.Errorf("Unimplemented method")
}

func (cr *chainReader) Start(ctx context.Context) error {
	if err := cr.initialize(); err != nil {
		return fmt.Errorf("Failed to initialize ChainReader: %w", err)
	}
	return nil
}
func (cr *chainReader) Close() error { return nil }

func (cr *chainReader) Ready() error { return nil }
func (cr *chainReader) HealthReport() map[string]error {
	return map[string]error{cr.Name(): nil}
}
func (cr *chainReader) Name() string { return cr.lggr.Name() }
