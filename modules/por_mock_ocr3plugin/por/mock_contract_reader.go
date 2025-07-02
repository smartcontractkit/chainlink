package por

import (
	"context"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// A mock implementation of the ContractReader interface for testing purposes.
type MockContractReader struct {
	digest types.ConfigDigest
}

// NewMockContractReader creates a new instance of MockContractReader.
func NewMockContractReader(configDigest types.ConfigDigest) *MockContractReader {
	return &MockContractReader{configDigest}
}

// GetLatestTransmittedReportDetails simulates the retrieval of the latest transmission details from the contract on-chain.
func (m *MockContractReader) GetLatestTransmittedReportDetails(ctx context.Context, chainId ChainSelector) (TransmittedReportDetails, error) {
	return TransmittedReportDetails{m.digest, 0, time.Now()}, nil
}
