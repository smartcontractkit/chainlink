package generic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

func TestKeyBundleKeyring(t *testing.T) {
	t.Parallel()
	inner := &recordingKeyring{}
	k := keyBundleKeyring{keyring: inner}
	report := ocr3types.ReportWithInfo[string]{Report: []byte("report"), Info: "aptos"}

	sig, err := k.Sign(ocrtypes.ConfigDigest{1}, 2, report)
	require.NoError(t, err)
	assert.Equal(t, []byte("sig"), sig)
	assert.True(t, k.Verify(nil, ocrtypes.ConfigDigest{1}, 2, report, sig))

	for _, got := range inner.reports {
		assert.Equal(t, ocrtypes.Report("report"), got.Report)
		info := &structpb.Struct{}
		require.NoError(t, proto.Unmarshal(got.Info, info))
		assert.Equal(t, "aptos", info.Fields["keyBundleName"].GetStringValue())
	}
}

type recordingKeyring struct {
	reports []ocr3types.ReportWithInfo[[]byte]
}

func (r *recordingKeyring) PublicKey() ocrtypes.OnchainPublicKey { return nil }

func (r *recordingKeyring) Sign(_ ocrtypes.ConfigDigest, _ uint64, report ocr3types.ReportWithInfo[[]byte]) ([]byte, error) {
	r.reports = append(r.reports, report)
	return []byte("sig"), nil
}

func (r *recordingKeyring) Verify(_ ocrtypes.OnchainPublicKey, _ ocrtypes.ConfigDigest, _ uint64,
	report ocr3types.ReportWithInfo[[]byte], _ []byte) bool {
	r.reports = append(r.reports, report)
	return true
}

func (r *recordingKeyring) MaxSignatureLength() int { return 3 }
