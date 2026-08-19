package feeds

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"

	ksmocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
)

func Test_service_newOCR2ConfigMsg_OnchainSigningPubKey(t *testing.T) {
	t.Parallel()

	t.Run("EVM key bundle includes raw onchain signing pub key", func(t *testing.T) {
		t.Parallel()

		evmKb, err := ocr2key.New(corekeys.EVM)
		require.NoError(t, err)

		rawPubKey, ok := ocr2key.RawEVMOnChainPublicKey(evmKb)
		require.True(t, ok)
		require.Len(t, rawPubKey, 130) // 65-byte uncompressed pubkey, hex-encoded

		ocr2Store := ksmocks.NewOCR2(t)
		ocr2Store.On("Get", evmKb.ID()).Return(evmKb, nil)

		s := &service{ocr2KeyStore: ocr2Store}
		msg, err := s.newOCR2ConfigMsg(OCR2ConfigModel{
			Enabled:     true,
			KeyBundleID: null.StringFrom(evmKb.ID()),
		})
		require.NoError(t, err)
		require.NotNil(t, msg.OcrKeyBundle)
		assert.Equal(t, evmKb.OnChainPublicKey(), msg.OcrKeyBundle.OnchainSigningAddress)
		assert.Equal(t, rawPubKey, msg.OcrKeyBundle.OnchainSigningPubKey)
	})

	t.Run("non-EVM key bundle omits raw onchain signing pub key", func(t *testing.T) {
		t.Parallel()

		solKb, err := ocr2key.New(corekeys.Solana)
		require.NoError(t, err)

		ocr2Store := ksmocks.NewOCR2(t)
		ocr2Store.On("Get", solKb.ID()).Return(solKb, nil)

		s := &service{ocr2KeyStore: ocr2Store}
		msg, err := s.newOCR2ConfigMsg(OCR2ConfigModel{
			Enabled:     true,
			KeyBundleID: null.StringFrom(solKb.ID()),
		})
		require.NoError(t, err)
		require.NotNil(t, msg.OcrKeyBundle)
		assert.Empty(t, msg.OcrKeyBundle.OnchainSigningPubKey)
	})

	t.Run("disabled config does not fetch key bundle", func(t *testing.T) {
		t.Parallel()

		s := &service{}
		msg, err := s.newOCR2ConfigMsg(OCR2ConfigModel{Enabled: false})
		require.NoError(t, err)
		assert.False(t, msg.Enabled)
		assert.Nil(t, msg.OcrKeyBundle)
	})
}
