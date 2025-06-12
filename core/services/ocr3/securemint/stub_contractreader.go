package securemint

import (
	"context"
	"fmt"
	"time"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/por_mock_ocr3plugin/por"
	sm_plugin "github.com/smartcontractkit/por_mock_ocr3plugin/por"
)

var _ sm_plugin.ContractReader = &stubContractReader{}

// stubContractReader is a mock implementation of the ContractReader interface.
// It retrieves the latest config digest from the config contract and then uses that to return mocked report details.
type stubContractReader struct {
	contractConfigTracker ocrtypes.ContractConfigTracker
}

func newStubContractReader(contractConfigTracker ocrtypes.ContractConfigTracker) *stubContractReader {
	return &stubContractReader{
		contractConfigTracker: contractConfigTracker,
	}
}

func (m *stubContractReader) GetLatestTransmittedReportDetails(ctx context.Context, _ por.ChainSelector) (sm_plugin.TransmittedReportDetails, error) {
	_, configDigest, err := m.contractConfigTracker.LatestConfigDetails(ctx)
	if err != nil {
		return sm_plugin.TransmittedReportDetails{}, fmt.Errorf("failed to get config digest: %w", err)
	}

	return sm_plugin.TransmittedReportDetails{
		ConfigDigest:    [32]byte(configDigest),
		SeqNr:           1,          // Mock sequence number
		LatestTimestamp: time.Now(), // Mock timestamp
	}, nil
}
