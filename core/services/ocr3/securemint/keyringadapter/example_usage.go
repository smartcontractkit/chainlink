package keyringadapter

import (
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/por_mock_ocr3plugin/por"
)

// ExampleUsage demonstrates how to use the OnchainKeyringAdapter
// to wrap an existing OCR2 keyring for use with OCR3 PoR plugin.
func ExampleUsage() {
	// This is a simplified example showing the adapter usage pattern.
	// In real usage, you would obtain the OCR2 keyring from your keystore.

	// Step 1: Get an existing OCR2 OnchainKeyring
	// This could be from your keystore, e.g.:
	// ocr2Bundle, err := ocr2key.New(chaintype.EVM)
	// if err != nil { ... }
	// ocr2Keyring := ocr2Bundle

	// For this example, we'll use a mock keyring
	mockOCR2Keyring := &mockExampleKeyring{
		publicKey: types.OnchainPublicKey("example-public-key"),
		maxSigLen: 65, // typical for ECDSA signatures
	}

	// Step 2: Wrap the OCR2 keyring with the PoR adapter
	porKeyring := NewSecureMintOCR3OnchainKeyringAdapter(mockOCR2Keyring)

	// Step 3: Use the adapter as an OCR3 OnchainKeyring for PoR
	configDigest := types.ConfigDigest([32]byte{1, 2, 3, 4, 5}) // example digest
	seqNr := uint64(42)
	chainSelector := por.ChainSelector(1234) // example chain selector

	reportWithInfo := ocr3types.ReportWithInfo[por.ChainSelector]{
		Report: []byte("example-por-report"),
		Info:   chainSelector,
	}

	// Sign a report
	signature, err := porKeyring.Sign(configDigest, seqNr, reportWithInfo)
	if err != nil {
		fmt.Printf("Error signing report: %v\n", err)
		return
	}

	// Verify the signature
	isValid := porKeyring.Verify(
		porKeyring.PublicKey(),
		configDigest,
		seqNr,
		reportWithInfo,
		signature,
	)

	fmt.Printf("Report signed successfully\n")
	fmt.Printf("Signature length: %d bytes\n", len(signature))
	fmt.Printf("Max signature length: %d bytes\n", porKeyring.MaxSignatureLength())
	fmt.Printf("Signature valid: %t\n", isValid)
	fmt.Printf("Public key: %x\n", porKeyring.PublicKey())
}

// mockExampleKeyring is a simple mock implementation for demonstration purposes
type mockExampleKeyring struct {
	publicKey types.OnchainPublicKey
	maxSigLen int
}

func (m *mockExampleKeyring) PublicKey() types.OnchainPublicKey {
	return m.publicKey
}

func (m *mockExampleKeyring) Sign(ctx types.ReportContext, report types.Report) ([]byte, error) {
	// In a real implementation, this would use cryptographic signing
	return []byte("example-signature"), nil
}

func (m *mockExampleKeyring) Verify(
	pubKey types.OnchainPublicKey,
	ctx types.ReportContext,
	report types.Report,
	signature []byte,
) bool {
	// In a real implementation, this would verify the cryptographic signature
	return true
}

func (m *mockExampleKeyring) MaxSignatureLength() int {
	return m.maxSigLen
}
