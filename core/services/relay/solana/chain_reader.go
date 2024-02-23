package solana

import (
	"context"
	"errors"
	"fmt"

	commonservices "github.com/smartcontractkit/chainlink-common/pkg/services"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
)

const ServiceName = "SolanaChainReader"

type SolanaChainReaderService struct {
	lggr logger.Logger
	commonservices.StateMachine
}

var (
	ErrUnimplemented                         = errors.New("function unimplemented")
	_                services.ServiceCtx     = &SolanaChainReaderService{}
	_                commontypes.ChainReader = &SolanaChainReaderService{}
)

// NewChainReaderService is a constructor for a new ChainReaderService for Solana. Returns a nil service on error.
func NewChainReaderService(lggr logger.Logger) (*SolanaChainReaderService, error) {
	return &SolanaChainReaderService{
		lggr: lggr.Named(ServiceName),
	}, nil
}

// Name implements the services.ServiceCtx interface and returns the logger service name.
func (s *SolanaChainReaderService) Name() string {
	return s.lggr.Name()
}

// Start implements the services.ServiceCtx interface and starts necessary background services.
// An error is returned if starting any internal services fails. Subsequent calls to Start return
// and error.
func (s *SolanaChainReaderService) Start(_ context.Context) error {
	return s.StartOnce(ServiceName, func() error {
		return nil
	})
}

// Close implements the services.ServiceCtx interface and stops all background services and cleans
// up used resources. Subsequent calls to Close return an error.
func (s *SolanaChainReaderService) Close() error {
	return s.StopOnce(ServiceName, func() error {
		return nil
	})
}

// Ready implements the services.ServiceCtx interface and returns an error if starting the service
// encountered any errors or if the service is not ready to serve requests.
func (s *SolanaChainReaderService) Ready() error {
	return s.StateMachine.Ready()
}

// HealthReport implements the services.ServiceCtx interface and returns errors for any internal
// function or service that may have failed.
func (s *SolanaChainReaderService) HealthReport() map[string]error {
	return map[string]error{s.Name(): s.Healthy()}
}

// GetLatestValue implements the types.ChainReader interface and requests and parses on-chain
// data named by the provided contract, method, and params.
func (s *SolanaChainReaderService) GetLatestValue(_ context.Context, contractName, method string, params any, returnVal any) error {
	return fmt.Errorf("%w: GetLatestValue not available", ErrUnimplemented)
}

// Bind implements the types.ChainReader interface and allows new contract bindings to be added
// to the service.
func (s *SolanaChainReaderService) Bind(_ context.Context, bindings []commontypes.BoundContract) error {
	return fmt.Errorf("%w: Bind not available", ErrUnimplemented)
}
