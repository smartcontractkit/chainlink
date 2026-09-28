package vault_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"

	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestNewOCR2KeySigner_SignsDigestWithOnchainKey(t *testing.T) {
	t.Parallel()

	kb, err := ocr2key.New(corekeys.EVM)
	require.NoError(t, err)

	signer, err := vaultcap.NewOCR2KeySigner(kb)
	require.NoError(t, err)

	digest := "6acc256f1ba5bf28f834040ecd80ff28d316aa1806ef21ace295f45226b0a66d"
	sig, err := signer.Sign(t.Context(), []byte(digest))
	require.NoError(t, err)

	recovered, err := utils.GetSignersEthAddress([]byte(digest), sig)
	require.NoError(t, err)
	require.Equal(t, common.HexToAddress(kb.OnChainPublicKey()), recovered)
}

func TestNewOCR2KeySigner_NonEVMBundleRejected(t *testing.T) {
	t.Parallel()

	kb, err := ocr2key.New(corekeys.Solana)
	require.NoError(t, err)

	_, err = vaultcap.NewOCR2KeySigner(kb)
	require.ErrorContains(t, err, "requires an EVM OCR2 key bundle")
}

func TestNewOCR2KeySigner_NilBundleRejected(t *testing.T) {
	t.Parallel()

	_, err := vaultcap.NewOCR2KeySigner(nil)
	require.ErrorContains(t, err, "OCR2 key bundle is nil")
}
