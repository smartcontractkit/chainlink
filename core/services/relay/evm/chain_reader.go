package evm

import (
	"context"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
)

type ChainReaderService interface {
	services.ServiceCtx
	types.ChainReader
}

type chainReader struct {
	lggr logger.Logger
	lp   logpoller.LogPoller
}

// chainReader constructor
func NewChainReaderService(lggr logger.Logger, lp logpoller.LogPoller) (*chainReader, error) {
	return &chainReader{lggr, lp}, nil
}

func (cr *chainReader) GetLatestValue(ctx context.Context, bc types.BoundContract, method string, params any, returnVal any) ([]byte, error) {

	// TODO: implement GetLatestValue

	return nil, nil
}

func (cr *chainReader) Start(ctx context.Context) error { return nil }
func (cr *chainReader) Close() error                    { return nil }

func (cr *chainReader) Ready() error { return nil }
func (cr *chainReader) HealthReport() map[string]error {
	return map[string]error{cr.Name(): nil}
}
func (cr *chainReader) Name() string { return cr.lggr.Name() }
