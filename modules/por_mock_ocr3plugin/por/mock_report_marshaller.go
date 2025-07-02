package por

import (
	"context"
	"encoding/json"
)

type MockReportMarshaler struct {
	maxReportSize int
}

// NewMockReportMarshaler creates a new instance of MockReportMarshaler with the specified maximum report size.
func NewMockReportMarshaler() *MockReportMarshaler {
	return &MockReportMarshaler{1024}
}

// Serialize simulates the serialization of a PorReport into bytes.
func (m *MockReportMarshaler) Serialize(ctx context.Context, chain ChainSelector, report PorReport) ([]byte, error) {
	// Simulate a serialization process by encoding the report into JSON.
	// In a real implementation, this would involve more complex logic to handle the report's structure and the target chain.
	// E.g., using an encoding logic compatible with the DataFeedsCache contract in Chainlink.
	// (See https://github.com/smartcontractkit/chainlink/blob/develop/contracts/src/v0.8/data-feeds/DataFeedsCache.sol#L53)
	encodedReport, err := json.Marshal(report)
	if err != nil {
		return nil, err
	}

	return encodedReport, nil
}

func (m *MockReportMarshaler) MaxReportSize(ctx context.Context) int {
	// Return the maximum report size set during initialization.
	return m.maxReportSize
}
