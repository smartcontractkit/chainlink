package keyringadapter

import (
	"github.com/smartcontractkit/chainlink-common/pkg/types/core/securemint"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// SecureMintOCR3OnchainKeyringAdapter adapts an OCR2 OnchainKeyring to implement ocr3types.OnchainKeyring[ChainSelector]
// This adapter enables the use of existing OCR2 keyrings with the OCR3 PoR plugin.
// Copied and adapted from core/services/ocrcommon/adapters.go
// Ideally we use ocrcommon.OCR3OnchainKeyringMultiChainAdapter instead? Problem is that that one is not typed, it assumes []byte as the Report type, while we use por.ChainSelector.
type SecureMintOCR3OnchainKeyringAdapter struct {
	ocr2Keyring types.OnchainKeyring
}

// Ensure OnchainKeyringAdapter implements the OCR3 OnchainKeyring interface for PoR ChainSelector
var _ ocr3types.OnchainKeyring[securemint.ChainSelector] = &SecureMintOCR3OnchainKeyringAdapter{}

// NewSecureMintOCR3OnchainKeyringAdapter creates a new adapter that wraps an OCR2 OnchainKeyring
// to implement the OCR3 OnchainKeyring interface for PoR ChainSelector.
func NewSecureMintOCR3OnchainKeyringAdapter(keyring types.OnchainKeyring) *SecureMintOCR3OnchainKeyringAdapter {
	return &SecureMintOCR3OnchainKeyringAdapter{
		ocr2Keyring: keyring,
	}
}

// PublicKey returns the public key of the underlying OCR2 keyring.
func (adapter *SecureMintOCR3OnchainKeyringAdapter) PublicKey() types.OnchainPublicKey {
	return adapter.ocr2Keyring.PublicKey()
}

// Sign creates a signature over the given report using the OCR2 keyring.
// It converts the OCR3 parameters (config digest, sequence number, and report with info)
// into the OCR2 ReportContext format expected by the underlying keyring.
func (adapter *SecureMintOCR3OnchainKeyringAdapter) Sign(
	configDigest types.ConfigDigest,
	seqNr uint64,
	reportWithInfo ocr3types.ReportWithInfo[securemint.ChainSelector],
) (signature []byte, err error) {
	// Convert OCR3 parameters to OCR2 ReportContext
	// Note: seqNr is converted to uint32 for Epoch field, which may truncate for very large values
	reportContext := types.ReportContext{
		ReportTimestamp: types.ReportTimestamp{
			ConfigDigest: configDigest,
			Epoch:        uint32(seqNr), //nolint:gosec // Intentional conversion, matches OCR protocol
			Round:        0,             // OCR3 doesn't use rounds in the same way as OCR2
		},
		ExtraHash: [32]byte{}, // Initialize with empty hash
	}

	return adapter.ocr2Keyring.Sign(reportContext, reportWithInfo.Report)
}

// Verify verifies a signature over the given report using the OCR2 keyring.
// It converts the OCR3 parameters into the OCR2 ReportContext format for verification.
func (adapter *SecureMintOCR3OnchainKeyringAdapter) Verify(
	publicKey types.OnchainPublicKey,
	configDigest types.ConfigDigest,
	seqNr uint64,
	reportWithInfo ocr3types.ReportWithInfo[securemint.ChainSelector],
	signature []byte,
) bool {
	// Convert OCR3 parameters to OCR2 ReportContext
	// Note: seqNr is converted to uint32 for Epoch field, which may truncate for very large values
	reportContext := types.ReportContext{
		ReportTimestamp: types.ReportTimestamp{
			ConfigDigest: configDigest,
			Epoch:        uint32(seqNr), //nolint:gosec // Intentional conversion, matches OCR protocol
			Round:        0,             // OCR3 doesn't use rounds in the same way as OCR2
		},
		ExtraHash: [32]byte{}, // Initialize with empty hash
	}

	return adapter.ocr2Keyring.Verify(publicKey, reportContext, reportWithInfo.Report, signature)
}

// MaxSignatureLength returns the maximum signature length from the underlying OCR2 keyring.
func (adapter *SecureMintOCR3OnchainKeyringAdapter) MaxSignatureLength() int {
	return adapter.ocr2Keyring.MaxSignatureLength()
}
