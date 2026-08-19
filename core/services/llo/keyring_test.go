package llo

import (
	crand "crypto/rand"
	"fmt"
	"math"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var _ Key = &mockKey{}

type mockKey struct {
	format          llotypes.ReportFormat
	verify          bool
	maxSignatureLen int
	sig             []byte
}

func (m *mockKey) Sign(reportCtx types.ReportContext, report types.Report) ([]byte, error) {
	return m.sig, nil
}
func (m *mockKey) Verify(publicKey types.OnchainPublicKey, reportCtx types.ReportContext, report types.Report, signature []byte) bool {
	return m.verify
}
func (m *mockKey) Sign3(digest types.ConfigDigest, seqNr uint64, r types.Report) (signature []byte, err error) {
	return m.sig, nil
}
func (m *mockKey) Verify3(publicKey types.OnchainPublicKey, cd types.ConfigDigest, seqNr uint64, r types.Report, signature []byte) bool {
	return m.verify
}
func (m *mockKey) SignBlob(b []byte) (sig []byte, err error) { return m.sig, nil }
func (m *mockKey) VerifyBlob(publicKey types.OnchainPublicKey, b []byte, sig []byte) bool {
	return m.verify
}

func (m *mockKey) PublicKey() types.OnchainPublicKey {
	b := make([]byte, m.maxSignatureLen)
	for i := range m.maxSignatureLen {
		b[i] = byte(255)
	}
	return types.OnchainPublicKey(b)
}

func (m *mockKey) MaxSignatureLength() int {
	return m.maxSignatureLen
}

func Test_Keyring(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

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

	cd, err := types.BytesToConfigDigest(mustRandBytes(32))
	require.NoError(t, err)
	seqNr := rand.Uint64N(math.MaxUint32 << 8)
	t.Run("Sign+Verify", func(t *testing.T) {
		t.Parallel()
		for _, tc := range cases {
			t.Run(tc.format.String(), func(t *testing.T) {
				t.Parallel()
				ksLocal := map[llotypes.ReportFormat]Key{
					llotypes.ReportFormatEVMPremiumLegacy: &mockKey{format: llotypes.ReportFormatEVMPremiumLegacy, maxSignatureLen: 1, sig: []byte("sig-1")},
					llotypes.ReportFormatJSON:             &mockKey{format: llotypes.ReportFormatJSON, maxSignatureLen: 2, sig: []byte("sig-2")},
					llotypes.ReportFormatEVMStreamlined:   &mockKey{format: llotypes.ReportFormatEVMStreamlined, maxSignatureLen: 6, sig: []byte("sig-6")},
				}
				krLocal := NewOnchainKeyring(lggr, ksLocal, 2)
				k := ksLocal[tc.format]

				sig, err := krLocal.Sign(cd, seqNr, ocr3types.ReportWithInfo[llotypes.ReportInfo]{Info: llotypes.ReportInfo{ReportFormat: tc.format}})
				require.NoError(t, err)

				assert.Equal(t, fmt.Appendf(nil, "sig-%d", tc.format), sig)

				assert.False(t, krLocal.Verify(nil, cd, seqNr, ocr3types.ReportWithInfo[llotypes.ReportInfo]{Info: llotypes.ReportInfo{ReportFormat: tc.format}}, sig))

				k.(*mockKey).verify = true
			})
		}
	})

	t.Run("MaxSignatureLength", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, 6+2+1, kr.MaxSignatureLength())
	})
	t.Run("PublicKey", func(t *testing.T) {
		t.Parallel()
		b := make([]byte, 6+2+1)
		for i := range b {
			b[i] = byte(255)
		}
		assert.Equal(t, types.OnchainPublicKey(b), kr.PublicKey())
	})
}

func mustRandBytes(n int) (b []byte) {
	b = make([]byte, n)
	_, err := crand.Read(b)
	if err != nil {
		panic(err)
	}
	return
}
