package securemint

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/por_mock_ocr3plugin/por"
	sm_plugin "github.com/smartcontractkit/por_mock_ocr3plugin/por"
)

var _ sm_plugin.ContractReader = &stubContractReader{}

// stubContractReader is a mock implementation of the ContractReader interface.
// It retrieves the latest config digest from the config contract and then uses that to return a mocked report.
// This is needed so that sm_plugin.ShouldTransmitAcceptedReport() does not fail (it checks the config digest).
type stubContractReader struct {
	getConfigDigestFunc func() ([32]byte, error)
}

// TODO(gg): make it use the ContractConfigTracker (and call LatestConfigDetails()) instead of a function.
func newStubContractReader(getConfigDigestFunc func() ([32]byte, error)) *stubContractReader {
	return &stubContractReader{
		getConfigDigestFunc: getConfigDigestFunc,
	}
}

func (m *stubContractReader) GetLatestTransmittedReportDetails(ctx context.Context, chainId por.ChainSelector) (sm_plugin.TransmittedReportDetails, error) {
	configDigest, err := m.getConfigDigestFunc()
	if err != nil {
		return sm_plugin.TransmittedReportDetails{}, fmt.Errorf("failed to get config digest: %w", err)
	}

	return sm_plugin.TransmittedReportDetails{
		ConfigDigest:    configDigest,
		SeqNr:           1,          // Mock sequence number
		LatestTimestamp: time.Now(), // Mock timestamp
	}, nil
}
