package evm

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-automation/pkg/v3/plugin"
)

func TestNewOnchainKeyringV3Wrapper(t *testing.T) {
	t.Run("the on chain keyring wrapper gets the public key and max signature length from the wrapped keyring", func(t *testing.T) {
		onchainKeyring := &mockOnchainKeyring{
			MaxSignatureLengthFn: func() int {
				return 123
			},
			PublicKeyFn: func() types.OnchainPublicKey {
				return types.OnchainPublicKey([]byte("abcdef"))
			},
		}
		keyring := NewOnchainKeyringV3Wrapper(onchainKeyring)
		assert.Equal(t, 123, keyring.MaxSignatureLength())
		assert.Equal(t, types.OnchainPublicKey([]byte("abcdef")), keyring.PublicKey())
	})

	t.Run("a report context is created and the wrapped keyring signs the report", func(t *testing.T) {
		onchainKeyring := &mockOnchainKeyring{
			SignFn: func(context types.ReportContext, report types.Report) (signature []byte, err error) {
				assert.Equal(t, types.ReportContext{
					ReportTimestamp: types.ReportTimestamp{
						ConfigDigest: types.ConfigDigest([32]byte{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4}),
						Epoch:        101,
						Round:        0,
					},
					ExtraHash: [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				}, context)
				assert.Equal(t, types.Report([]byte("a report to sign")), report)
				return []byte("signature"), err
			},
		}
		keyring := NewOnchainKeyringV3Wrapper(onchainKeyring)
		signed, err := keyring.Sign(
			types.ConfigDigest([32]byte{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4}),
			101,
			ocr3types.ReportWithInfo[plugin.AutomationReportInfo]{
				Report: []byte("a report to sign"),
				Info:   plugin.AutomationReportInfo{},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, []byte("signature"), signed)
	})

	t.Run("a report context is created and the wrapped keyring verifies the report", func(t *testing.T) {
		onchainKeyring := &mockOnchainKeyring{
			VerifyFn: func(pk types.OnchainPublicKey, rc types.ReportContext, r types.Report, signature []byte) bool {
				assert.Equal(t, types.OnchainPublicKey([]byte("key")), pk)
				assert.Equal(t, types.ReportContext{
					ReportTimestamp: types.ReportTimestamp{
						ConfigDigest: types.ConfigDigest([32]byte{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4}),
						Epoch:        999,
						Round:        0,
					},
					ExtraHash: [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				}, rc)
				assert.Equal(t, types.Report([]byte("a report to sign")), r)
				assert.Equal(t, []byte("signature"), signature)
				return true
			},
		}
		keyring := NewOnchainKeyringV3Wrapper(onchainKeyring)
		valid := keyring.Verify(
			types.OnchainPublicKey([]byte("key")),
			types.ConfigDigest([32]byte{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4}),
			999,
			ocr3types.ReportWithInfo[plugin.AutomationReportInfo]{
				Report: []byte("a report to sign"),
				Info:   plugin.AutomationReportInfo{},
			},
			[]byte("signature"),
		)
		assert.True(t, valid)
	})
}

type mockOnchainKeyring struct {
	PublicKeyFn          func() types.OnchainPublicKey
	SignFn               func(types.ReportContext, types.Report) (signature []byte, err error)
	VerifyFn             func(_ types.OnchainPublicKey, _ types.ReportContext, _ types.Report, signature []byte) bool
	MaxSignatureLengthFn func() int
}

func (k *mockOnchainKeyring) PublicKey() types.OnchainPublicKey {
	return k.PublicKeyFn()
}

func (k *mockOnchainKeyring) Sign(c types.ReportContext, r types.Report) (signature []byte, err error) {
	return k.SignFn(c, r)
}

func (k *mockOnchainKeyring) Verify(pk types.OnchainPublicKey, c types.ReportContext, r types.Report, signature []byte) bool {
	return k.VerifyFn(pk, c, r, signature)
}

func (k *mockOnchainKeyring) MaxSignatureLength() int {
	return k.MaxSignatureLengthFn()
}
