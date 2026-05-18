package vault

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
)

// Wire-size allowances for top-level protobuf wrappers. These are intentionally
// small conservative paddings (not exact protobuf encodings) so byte budgets
// stay slightly below the OCR3.1 hard caps.
const (
	// observationsOuterProtoOverhead accounts for tags/length prefixes on the
	// Observations message fields (observations, pending_queue_items, sort_nonce).
	observationsOuterProtoOverhead = 10
	// outcomesOuterProtoOverhead accounts for the Outcomes wrapper around the
	// repeated outcomes field. It is a separate named constant from
	// observationsOuterProtoOverhead because the two messages can diverge on the wire.
	outcomesOuterProtoOverhead = 10
	// sortNonceWireBytes is the 32-byte random SortNonce plus typical proto3 bytes-field overhead.
	sortNonceWireBytes = 35
)

// computePluginByteBudgets derives how many bytes are available for the
// repeated Observation messages and repeated Outcome messages respectively,
// after subtracting fixed costs that Observation() always pays next round
// (blob handles for pending-queue blobs + sort nonce + outer framing).
//
// Used by NewReportingPlugin and by applyTestByteBudgets in tests.
func computePluginByteBudgets(ctx context.Context, cfg *ReportingPluginConfig, maxObservationBytes, maxReportsPlusPrecursorBytes, n, f int) (obsBudget int, precursorBudget int, err error) {
	if n <= 0 || f < 0 {
		return 0, 0, fmt.Errorf("invalid OCR config for byte budgets: N=%d F=%d", n, f)
	}
	batchSizeLimit, err := cfg.MaxBatchSize.Limit(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("could not fetch max batch size limit: %w", err)
	}
	maxBlobHandleCount := 2 * batchSizeLimit
	blobHandleBytes := ocr3_1types.BlobHandleMarshalledBytesUpperBound(n, f)
	obsBudget = maxObservationBytes -
		maxBlobHandleCount*blobHandleBytes -
		sortNonceWireBytes -
		observationsOuterProtoOverhead
	precursorBudget = maxReportsPlusPrecursorBytes - outcomesOuterProtoOverhead
	if obsBudget <= 0 {
		return 0, 0, fmt.Errorf(
			"VaultMaxObservationSizeLimit leaves no room for observations after pending-queue blob handle overhead (effective obs budget=%d bytes; N=%d F=%d batchSize=%d blobHandleBytes=%d)",
			obsBudget, n, f, batchSizeLimit, blobHandleBytes,
		)
	}
	if precursorBudget <= 0 {
		return 0, 0, fmt.Errorf(
			"VaultMaxReportsPlusPrecursorSizeLimit is too small (effective precursor budget=%d bytes)",
			precursorBudget,
		)
	}
	return obsBudget, precursorBudget, nil
}
