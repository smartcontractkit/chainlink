package keyringadapter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/types/core/securemint"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// mockOCR2OnchainKeyring is a mock implementation of types.OnchainKeyring for testing
type mockOCR2OnchainKeyring struct {
	publicKey          types.OnchainPublicKey
	maxSignatureLength int
	signFunc           func(types.ReportContext, types.Report) ([]byte, error)
	verifyFunc         func(types.OnchainPublicKey, types.ReportContext, types.Report, []byte) bool
}

func (m *mockOCR2OnchainKeyring) PublicKey() types.OnchainPublicKey {
	return m.publicKey
}

func (m *mockOCR2OnchainKeyring) Sign(ctx types.ReportContext, report types.Report) ([]byte, error) {
	if m.signFunc != nil {
		return m.signFunc(ctx, report)
	}
	return []byte("mock-signature"), nil
}

func (m *mockOCR2OnchainKeyring) Verify(
	pubKey types.OnchainPublicKey,
	ctx types.ReportContext,
	report types.Report,
	signature []byte,
) bool {
	if m.verifyFunc != nil {
		return m.verifyFunc(pubKey, ctx, report, signature)
	}
	return true
}

func (m *mockOCR2OnchainKeyring) MaxSignatureLength() int {
	return m.maxSignatureLength
}

func TestPorOnchainKeyringAdapter(t *testing.T) {
	// Setup test data
	testPublicKey := types.OnchainPublicKey("test-public-key")
	testConfigDigest := types.ConfigDigest([32]byte{1, 2, 3, 4, 5})
	testSeqNr := uint64(42)
	testReport := types.Report([]byte("test-report"))
	testChainSelector := securemint.ChainSelector(1234)
	testSignature := []byte("test-signature")
	testMaxSigLen := 65

	reportWithInfo := ocr3types.ReportWithInfo[securemint.ChainSelector]{
		Report: testReport,
		Info:   testChainSelector,
	}

	t.Run("adapter implements the correct interface", func(t *testing.T) {
		mockKeyring := &mockOCR2OnchainKeyring{
			publicKey:          testPublicKey,
			maxSignatureLength: testMaxSigLen,
		}

		adapter := NewSecureMintOCR3OnchainKeyringAdapter(mockKeyring)

		// Verify that the adapter implements the OCR3 OnchainKeyring interface
		var _ ocr3types.OnchainKeyring[securemint.ChainSelector] = adapter
	})

	t.Run("PublicKey returns the underlying keyring's public key", func(t *testing.T) {
		mockKeyring := &mockOCR2OnchainKeyring{
			publicKey:          testPublicKey,
			maxSignatureLength: testMaxSigLen,
		}

		adapter := NewSecureMintOCR3OnchainKeyringAdapter(mockKeyring)
		assert.Equal(t, testPublicKey, adapter.PublicKey())
	})

	t.Run("MaxSignatureLength returns the underlying keyring's max signature length", func(t *testing.T) {
		mockKeyring := &mockOCR2OnchainKeyring{
			publicKey:          testPublicKey,
			maxSignatureLength: testMaxSigLen,
		}

		adapter := NewSecureMintOCR3OnchainKeyringAdapter(mockKeyring)
		assert.Equal(t, testMaxSigLen, adapter.MaxSignatureLength())
	})

	t.Run("Sign correctly converts OCR3 parameters to OCR2 format", func(t *testing.T) {
		var capturedReportContext types.ReportContext
		var capturedReport types.Report

		mockKeyring := &mockOCR2OnchainKeyring{
			publicKey:          testPublicKey,
			maxSignatureLength: testMaxSigLen,
			signFunc: func(ctx types.ReportContext, report types.Report) ([]byte, error) {
				capturedReportContext = ctx
				capturedReport = report
				return testSignature, nil
			},
		}

		adapter := NewSecureMintOCR3OnchainKeyringAdapter(mockKeyring)
		signature, err := adapter.Sign(testConfigDigest, testSeqNr, reportWithInfo)

		require.NoError(t, err)
		assert.Equal(t, testSignature, signature)

		// Verify the conversion from OCR3 to OCR2 format
		assert.Equal(t, testConfigDigest, capturedReportContext.ReportTimestamp.ConfigDigest)
		assert.Equal(t, uint32(testSeqNr), capturedReportContext.ReportTimestamp.Epoch)
		assert.Equal(t, uint8(0), capturedReportContext.ReportTimestamp.Round)
		assert.Equal(t, [32]byte{}, capturedReportContext.ExtraHash)
		assert.Equal(t, testReport, capturedReport)
	})

	t.Run("Verify correctly converts OCR3 parameters to OCR2 format", func(t *testing.T) {
		var capturedPublicKey types.OnchainPublicKey
		var capturedReportContext types.ReportContext
		var capturedReport types.Report
		var capturedSignature []byte

		mockKeyring := &mockOCR2OnchainKeyring{
			publicKey:          testPublicKey,
			maxSignatureLength: testMaxSigLen,
			verifyFunc: func(
				pubKey types.OnchainPublicKey,
				ctx types.ReportContext,
				report types.Report,
				signature []byte,
			) bool {
				capturedPublicKey = pubKey
				capturedReportContext = ctx
				capturedReport = report
				capturedSignature = signature
				return true
			},
		}

		adapter := NewSecureMintOCR3OnchainKeyringAdapter(mockKeyring)
		result := adapter.Verify(testPublicKey, testConfigDigest, testSeqNr, reportWithInfo, testSignature)

		assert.True(t, result)

		// Verify the conversion from OCR3 to OCR2 format
		assert.Equal(t, testPublicKey, capturedPublicKey)
		assert.Equal(t, testConfigDigest, capturedReportContext.ReportTimestamp.ConfigDigest)
		assert.Equal(t, uint32(testSeqNr), capturedReportContext.ReportTimestamp.Epoch)
		assert.Equal(t, uint8(0), capturedReportContext.ReportTimestamp.Round)
		assert.Equal(t, [32]byte{}, capturedReportContext.ExtraHash)
		assert.Equal(t, testReport, capturedReport)
		assert.Equal(t, testSignature, capturedSignature)
	})

	t.Run("Sign and Verify work together", func(t *testing.T) {
		mockKeyring := &mockOCR2OnchainKeyring{
			publicKey:          testPublicKey,
			maxSignatureLength: testMaxSigLen,
		}

		adapter := NewSecureMintOCR3OnchainKeyringAdapter(mockKeyring)

		// Sign a report
		signature, err := adapter.Sign(testConfigDigest, testSeqNr, reportWithInfo)
		require.NoError(t, err)

		// Verify the signature
		isValid := adapter.Verify(testPublicKey, testConfigDigest, testSeqNr, reportWithInfo, signature)
		assert.True(t, isValid)
	})
}
