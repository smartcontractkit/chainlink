package securemint

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/por_mock_ocr3plugin/por"
)

// Ensure StubContractTransmitter implements the ContractTransmitter interface
var _ ocr3types.ContractTransmitter[por.ChainSelector] = (*stubContractTransmitter)(nil)

// stubContractTransmitter is a stub implementation of the ContractTransmitter interface
// that logs messages when its functions are invoked instead of performing actual operations.
type stubContractTransmitter struct {
	services.Service

	logger      logger.Logger
	fromAccount types.Account
}

// StubTransmissionCounter is a global counter to track the number of transmissions, used for testing purposes.
// Since this is a stub implementation, we can get away with it.
var StubTransmissionCounter atomic.Int32

// newStubContractTransmitter creates a new StubContractTransmitter instance
func NewStubContractTransmitter(logger logger.Logger, fromAccount types.Account) *stubContractTransmitter {
	t := &stubContractTransmitter{
		logger:      logger,
		fromAccount: fromAccount,
	}
	t.Service = services.Config{
		Name:  "StubContractTransmitter",
		Start: t.Start,
		Close: t.Close,
	}.NewService(logger)

	return t
}

// Transmit logs the transmission details instead of actually transmitting
func (s *stubContractTransmitter) Transmit(
	_ context.Context,
	configDigest types.ConfigDigest,
	seqNr uint64,
	reportWithInfo ocr3types.ReportWithInfo[por.ChainSelector],
	aos []types.AttributedOnchainSignature,
) error {
	s.logger.Info("Transmit called ", map[string]any{
		"configDigest":    fmt.Sprintf("%x", configDigest),
		"sequenceNumber":  seqNr,
		"reportLength":    len(reportWithInfo.Report),
		"reportInfo":      reportWithInfo.Info,
		"signaturesCount": len(aos),
	})

	// Log report details if available
	if len(reportWithInfo.Report) > 0 {
		s.logger.Debug("Report data ", map[string]any{
			"reportHex": fmt.Sprintf("%x", reportWithInfo.Report),
		})
	}

	s.logger.Info("Transmit completed successfully (stub implementation)", nil)
	StubTransmissionCounter.Add(1)
	return nil
}

// FromAccount returns the configured account and logs the call
func (s *stubContractTransmitter) FromAccount(_ context.Context) (types.Account, error) {
	s.logger.Debug("FromAccount called ", map[string]any{
		"account": string(s.fromAccount),
	})

	return s.fromAccount, nil
}
