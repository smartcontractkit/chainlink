package ocr2key

import (
	"testing"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/require"
)

func TestTerraKeyRing_Sign_Verify(t *testing.T) {
	kr1 := newTerraKeyring()
	kr2 := newTerraKeyring()
	ctx := ocrtypes.ReportContext{}
	report := ocrtypes.Report{}
	sig, err := kr1.Sign(ctx, report)
	require.NoError(t, err)
	result := kr2.Verify(kr1.PublicKey(), ctx, report, sig)
	require.True(t, result)
}
