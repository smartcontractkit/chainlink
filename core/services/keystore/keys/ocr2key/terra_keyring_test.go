package ocr2key

import (
	cryptorand "crypto/rand"
	"testing"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/require"
)

func TestTerraKeyRing_Sign_Verify(t *testing.T) {
	kr1, err := newTerraKeyring(cryptorand.Reader)
	require.NoError(t, err)
	kr2, err := newTerraKeyring(cryptorand.Reader)
	require.NoError(t, err)
	ctx := ocrtypes.ReportContext{}
	report := ocrtypes.Report{}
	sig, err := kr1.Sign(ctx, report)
	require.NoError(t, err)
	result := kr2.Verify(kr1.PublicKey(), ctx, report, sig)
	require.True(t, result)
}
