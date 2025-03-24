package llo

import (
	"fmt"
	"math"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var _ Key = &mockKey{}

type mockKey struct {
	format          llotypes.ReportFormat
	verify          bool
	maxSignatureLen int
	sig             []byte
}

func (m *mockKey) Sign(reportCtx ocrtypes.ReportContext, report ocrtypes.Report) ([]byte, error) {
	return m.sig, nil
}
func (m *mockKey) Verify(publicKey ocrtypes.OnchainPublicKey, reportCtx ocrtypes.ReportContext, report ocrtypes.Report, signature []byte) bool {
	return m.verify
}
func (m *mockKey) Sign3(digest ocrtypes.ConfigDigest, seqNr uint64, r ocrtypes.Report) (signature []byte, err error) {
	return m.sig, nil
}
func (m *mockKey) Verify3(publicKey ocrtypes.OnchainPublicKey, cd ocrtypes.ConfigDigest, seqNr uint64, r ocrtypes.Report, signature []byte) bool {
	return m.verify
}
func (m *mockKey) SignBlob(b []byte) (sig []byte, err error) { return m.sig, nil }
func (m *mockKey) VerifyBlob(publicKey ocrtypes.OnchainPublicKey, b []byte, sig []byte) bool {
	return m.verify
}

func (m *mockKey) PublicKey() ocrtypes.OnchainPublicKey {
	b := make([]byte, m.maxSignatureLen)
	for i := 0; i < m.maxSignatureLen; i++ {
		b[i] = byte(255)
	}
	return ocrtypes.OnchainPublicKey(b)
}

func (m *mockKey) MaxSignatureLength() int {
	return m.maxSignatureLen
}

func (m *mockKey) reset(format llotypes.ReportFormat) {
	m.format = format
	m.verify = false
}

func Test_Keyring(t *testing.T) {
	lggr := logger.TestLogger(t)

	ks := map[llotypes.ReportFormat]Key{
		llotypes.ReportFormatEVMPremiumLegacy: &mockKey{format: llotypes.ReportFormatEVMPremiumLegacy, maxSignatureLen: 1, sig: []byte("sig-1")},
		llotypes.ReportFormatJSON:             &mockKey{format: llotypes.ReportFormatJSON, maxSignatureLen: 2, sig: []byte("sig-2")},
		llotypes.ReportFormatEVMStreamlined:   &mockKey{format: llotypes.ReportFormatEVMStreamlined, maxSignatureLen: 6, sig: []byte("sig-6")},
	}

	kr := NewOnchainKeyring(lggr, ks, 2)

	cases := []struct {
		format llotypes.ReportFormat
	}{
		{
			llotypes.ReportFormatEVMPremiumLegacy,
		},
		{
			llotypes.ReportFormatJSON,
		},
		{
			llotypes.ReportFormatEVMStreamlined,
		},
	}

	cd, err := ocrtypes.BytesToConfigDigest(testutils.MustRandBytes(32))
	require.NoError(t, err)
	seqNr := rand.Uint64N(math.MaxUint32 << 8)
	t.Run("Sign+Verify", func(t *testing.T) {
		for _, tc := range cases {
			t.Run(tc.format.String(), func(t *testing.T) {
				k := ks[tc.format]
				defer k.(*mockKey).reset(tc.format)

				sig, err := kr.Sign(cd, seqNr, ocr3types.ReportWithInfo[llotypes.ReportInfo]{Info: llotypes.ReportInfo{ReportFormat: tc.format}})
				require.NoError(t, err)

				assert.Equal(t, []byte(fmt.Sprintf("sig-%d", tc.format)), sig)

				assert.False(t, kr.Verify(nil, cd, seqNr, ocr3types.ReportWithInfo[llotypes.ReportInfo]{Info: llotypes.ReportInfo{ReportFormat: tc.format}}, sig))

				k.(*mockKey).verify = true
			})
		}
	})

	t.Run("MaxSignatureLength", func(t *testing.T) {
		assert.Equal(t, 6+2+1, kr.MaxSignatureLength())
	})
	t.Run("PublicKey", func(t *testing.T) {
		b := make([]byte, 6+2+1)
		for i := 0; i < len(b); i++ {
			b[i] = byte(255)
		}
		assert.Equal(t, types.OnchainPublicKey(b), kr.PublicKey())
	})
}
