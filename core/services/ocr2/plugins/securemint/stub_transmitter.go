package securemint

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/por_mock_ocr3plugin/por"
)

// Ensure StubContractTransmitter implements the ContractTransmitter interface
var _ ocr3types.ContractTransmitter[por.ChainSelector] = (*stubContractTransmitter)(nil)

// stubContractTransmitter is a stub implementation of the ContractTransmitter interface
// that logs messages when its functions are invoked instead of performing actual operations.
type stubContractTransmitter struct {
	logger      logger.Logger
	fromAccount types.Account
}

// newStubContractTransmitter creates a new StubContractTransmitter instance
func newStubContractTransmitter(logger logger.Logger, fromAccount types.Account) *stubContractTransmitter {
	return &stubContractTransmitter{
		logger:      logger,
		fromAccount: fromAccount,
	}
}

// Transmit logs the transmission details instead of actually transmitting
func (s *stubContractTransmitter) Transmit(
	ctx context.Context,
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

	// Log signature details
	for i, sig := range aos {
		s.logger.Debug("Signature details ", map[string]any{
			"signatureIndex": i,
			"signer":         fmt.Sprintf("%x", sig.Signer),
			"signatureHex":   fmt.Sprintf("%x", sig.Signature),
		})
	}

	s.logger.Info("Transmit completed successfully (stub implementation)", nil)
	return nil
}

// FromAccount returns the configured account and logs the call
func (s *stubContractTransmitter) FromAccount(ctx context.Context) (types.Account, error) {
	s.logger.Debug("FromAccount called ", map[string]any{
		"account": string(s.fromAccount),
	})

	return s.fromAccount, nil
}
