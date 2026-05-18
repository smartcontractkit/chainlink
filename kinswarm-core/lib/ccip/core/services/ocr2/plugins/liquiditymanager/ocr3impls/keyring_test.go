package ocr3impls

import (
	"testing"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

func TestKeyring(t *testing.T) {
	t.Run("PublicKey", func(t *testing.T) {
		t.Parallel()
		bundle, err := ocr2key.New(chaintype.EVM)
		require.NoError(t, err, "failed to create key bundle")
		keyring := NewOnchainKeyring[models.Report](bundle, logger.TestLogger(t))
		require.Equal(t, bundle.PublicKey(), keyring.PublicKey())
	})

	t.Run("MaxSignatureLength", func(t *testing.T) {
		t.Parallel()
		bundle, err := ocr2key.New(chaintype.EVM)
		require.NoError(t, err, "failed to create key bundle")
		keyring := NewOnchainKeyring[models.Report](bundle, logger.TestLogger(t))
		require.Equal(t, 65, keyring.MaxSignatureLength())
	})

	t.Run("Sign/Verify", func(t *testing.T) {
		t.Parallel()
		bundle, err := ocr2key.New(chaintype.EVM)
		require.NoError(t, err, "failed to create key bundle")
		keyring := NewOnchainKeyring[models.Report](bundle, logger.TestLogger(t))
		digest := testutils.Random32Byte()
		seqNr := uint64(1)
		report := testutils.Random32Byte()
		sig, err := keyring.Sign(digest, seqNr, ocr3types.ReportWithInfo[models.Report]{
			Info: models.Report{
				ConfigDigest: models.ConfigDigest{
					ConfigDigest: digest,
				},
			},
			Report: report[:],
		})
		require.NoError(t, err, "failed to sign")
		require.True(t, keyring.Verify(
			keyring.PublicKey(),
			digest,
			seqNr,
			ocr3types.ReportWithInfo[models.Report]{
				Info: models.Report{
					ConfigDigest: models.ConfigDigest{
						ConfigDigest: digest,
					},
				},
				Report: report[:],
			},
			sig,
		))

		// bork sig, verify should fail
		old := sig[0]
		sig[0] = sig[0] ^ 0xFF
		require.False(t, keyring.Verify(
			keyring.PublicKey(),
			digest,
			seqNr,
			ocr3types.ReportWithInfo[models.Report]{
				Info: models.Report{
					ConfigDigest: models.ConfigDigest{
						ConfigDigest: digest,
					},
				},
				Report: report[:],
			},
			sig,
		))

		sig[0] = old
		// bork report, verify should fail
		report[0] = report[0] ^ 0xFF
		require.False(t, keyring.Verify(
			keyring.PublicKey(),
			digest,
			seqNr,
			ocr3types.ReportWithInfo[models.Report]{
				Info: models.Report{
					ConfigDigest: models.ConfigDigest{
						ConfigDigest: digest,
					},
				},
				Report: report[:],
			},
			sig,
		))

		// bork config digest, verify should fail
		digest[0] = digest[0] ^ 0xFF
		require.False(t, keyring.Verify(
			keyring.PublicKey(),
			digest,
			seqNr,
			ocr3types.ReportWithInfo[models.Report]{
				Info: models.Report{
					ConfigDigest: models.ConfigDigest{
						ConfigDigest: digest,
					},
				},
				Report: report[:],
			},
			sig,
		))
	})
}
