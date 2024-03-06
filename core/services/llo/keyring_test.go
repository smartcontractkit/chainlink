package llo

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ocr3types "github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
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
}

func (m *mockKey) Sign3(digest ocrtypes.ConfigDigest, seqNr uint64, r ocrtypes.Report) (signature []byte, err error) {
	return []byte(fmt.Sprintf("sig-%d", m.format)), nil
}

func (m *mockKey) Verify3(publicKey ocrtypes.OnchainPublicKey, cd ocrtypes.ConfigDigest, seqNr uint64, r ocrtypes.Report, signature []byte) bool {
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
		llotypes.ReportFormatEVM:      &mockKey{format: llotypes.ReportFormatEVM, maxSignatureLen: 1},
		llotypes.ReportFormatSolana:   &mockKey{format: llotypes.ReportFormatSolana, maxSignatureLen: 2},
		llotypes.ReportFormatCosmos:   &mockKey{format: llotypes.ReportFormatCosmos, maxSignatureLen: 4},
		llotypes.ReportFormatStarknet: &mockKey{format: llotypes.ReportFormatStarknet, maxSignatureLen: 8},
	}

	kr := NewOnchainKeyring(lggr, ks)

	cases := []struct {
		format llotypes.ReportFormat
	}{
		{
			llotypes.ReportFormatEVM,
		},
		{
			llotypes.ReportFormatSolana,
		},
		{
			llotypes.ReportFormatCosmos,
		},
		{
			llotypes.ReportFormatStarknet,
		},
	}

	cd, err := ocrtypes.BytesToConfigDigest(testutils.MustRandBytes(32))
	require.NoError(t, err)
	seqNr := rand.Uint64()
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

				for _, tc2 := range cases {
					verified := kr.Verify(nil, cd, seqNr, ocr3types.ReportWithInfo[llotypes.ReportInfo]{Info: llotypes.ReportInfo{ReportFormat: tc2.format}}, sig)
					if tc.format == tc2.format {
						assert.True(t, verified, "expected true for %s", tc2.format)
					} else {
						assert.False(t, verified, "expected false for %s", tc2.format)
					}
				}
			})
		}
	})

	t.Run("MaxSignatureLength", func(t *testing.T) {
		assert.Equal(t, 8+4+2+1, kr.MaxSignatureLength())
	})
	t.Run("PublicKey", func(t *testing.T) {
		b := make([]byte, 8+4+2+1)
		for i := 0; i < len(b); i++ {
			b[i] = byte(255)
		}
		assert.Equal(t, types.OnchainPublicKey(b), kr.PublicKey())
	})
}
