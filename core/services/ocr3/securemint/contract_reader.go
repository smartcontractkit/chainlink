package securemint

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/por_mock_ocr3plugin/por"
	sm_plugin "github.com/smartcontractkit/por_mock_ocr3plugin/por"
)

// secureMintContractReader implements the sm_plugin.ContractReader interface
// by wrapping a chain reader and providing secure mint specific functionality
type secureMintContractReader struct {
	chainReader   types.ContractReader
	chainSelector por.ChainSelector
	logger        logger.Logger
}

// Ensure secureMintContractReader implements the ContractReader interface
var _ sm_plugin.ContractReader = (*secureMintContractReader)(nil)

// GetLatestTransmittedReportDetails retrieves the latest transmitted report details from the chain
func (s *secureMintContractReader) GetLatestTransmittedReportDetails(ctx context.Context, chainSelector por.ChainSelector) (sm_plugin.TransmittedReportDetails, error) {
	if chainSelector != s.chainSelector {
		return sm_plugin.TransmittedReportDetails{}, fmt.Errorf("chain selector mismatch: expected %d, got %d", s.chainSelector, chainSelector)
	}

	// Call the contract to get the latest transmitted details using the correct method name
	var result sm_plugin.TransmittedReportDetails
	err := s.chainReader.GetLatestValue(ctx, "SecureMint.latestTransmittedDetails", primitives.Unconfirmed, nil, &result)
	if err != nil {
		return sm_plugin.TransmittedReportDetails{}, fmt.Errorf("failed to get latest transmitted details: %w", err)
	}

	return result, nil
}

// Close implements the Close method for services that need cleanup
func (s *secureMintContractReader) Close() error {
	// The underlying chain reader should handle its own cleanup
	// We don't need to do anything special here
	return nil
}
