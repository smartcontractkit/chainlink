package keystore

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/dkgencryptkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/dkgsignkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocrkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/solkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/terrakey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const password = "password"

func TestKeyRing_Encrypt_Decrypt(t *testing.T) {
	csa1, csa2 := csakey.MustNewV2XXXTestingOnly(big.NewInt(1)), csakey.MustNewV2XXXTestingOnly(big.NewInt(2))
	eth1, eth2 := mustNewEthKey(t), mustNewEthKey(t)
	ocr := []ocrkey.KeyV2{
		ocrkey.MustNewV2XXXTestingOnly(big.NewInt(1)),
		ocrkey.MustNewV2XXXTestingOnly(big.NewInt(2)),
	}
	var ocr2 []ocr2key.KeyBundle
	var ocr2_raw []ocr2key.Raw
	for _, chain := range chaintype.SupportedChainTypes {
		key := ocr2key.MustNewInsecure(rand.Reader, chain)
		ocr2 = append(ocr2, key)
		ocr2_raw = append(ocr2_raw, key.Raw())
	}
	p2p1, p2p2 := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1)), p2pkey.MustNewV2XXXTestingOnly(big.NewInt(2))
	sol1, sol2 := solkey.MustNewInsecure(rand.Reader), solkey.MustNewInsecure(rand.Reader)
	vrf1, vrf2 := vrfkey.MustNewV2XXXTestingOnly(big.NewInt(1)), vrfkey.MustNewV2XXXTestingOnly(big.NewInt(2))
	tk1, tk2 := terrakey.MustNewInsecure(rand.Reader), terrakey.MustNewInsecure(rand.Reader)
	dkgsign1, dkgsign2 := dkgsignkey.MustNewXXXTestingOnly(big.NewInt(1)), dkgsignkey.MustNewXXXTestingOnly(big.NewInt(2))
	dkgencrypt1, dkgencrypt2 := dkgencryptkey.MustNewXXXTestingOnly(big.NewInt(1)), dkgencryptkey.MustNewXXXTestingOnly(big.NewInt(2))
	originalKeyRingRaw := rawKeyRing{
		CSA:        []csakey.Raw{csa1.Raw(), csa2.Raw()},
		Eth:        []ethkey.Raw{eth1.Raw(), eth2.Raw()},
		OCR:        []ocrkey.Raw{ocr[0].Raw(), ocr[1].Raw()},
		OCR2:       ocr2_raw,
		P2P:        []p2pkey.Raw{p2p1.Raw(), p2p2.Raw()},
		Solana:     []solkey.Raw{sol1.Raw(), sol2.Raw()},
		VRF:        []vrfkey.Raw{vrf1.Raw(), vrf2.Raw()},
		Terra:      []terrakey.Raw{tk1.Raw(), tk2.Raw()},
		DKGSign:    []dkgsignkey.Raw{dkgsign1.Raw(), dkgsign2.Raw()},
		DKGEncrypt: []dkgencryptkey.Raw{dkgencrypt1.Raw(), dkgencrypt2.Raw()},
	}
	originalKeyRing, err := originalKeyRingRaw.keys()
	require.NoError(t, err)

	encryptedKeyRing, err := originalKeyRing.Encrypt(password, utils.FastScryptParams)
	require.NoError(t, err)
	decryptedKeyRing, err := encryptedKeyRing.Decrypt(password)
	require.NoError(t, err)
	// compare csa keys
	require.Equal(t, 2, len(decryptedKeyRing.CSA))
	require.Equal(t, originalKeyRing.CSA[csa1.ID()].PublicKey, decryptedKeyRing.CSA[csa1.ID()].PublicKey)
	require.Equal(t, originalKeyRing.CSA[csa2.ID()].PublicKey, decryptedKeyRing.CSA[csa2.ID()].PublicKey)
	// compare eth keys
	require.Equal(t, 2, len(decryptedKeyRing.Eth))
	require.Equal(t, originalKeyRing.Eth[eth1.ID()].Address, decryptedKeyRing.Eth[eth1.ID()].Address)
	require.Equal(t, originalKeyRing.Eth[eth2.ID()].Address, decryptedKeyRing.Eth[eth2.ID()].Address)
	// compare ocr keys
	require.Equal(t, 2, len(decryptedKeyRing.OCR))
	require.Equal(t, originalKeyRing.OCR[ocr[0].ID()].OnChainSigning.X, decryptedKeyRing.OCR[ocr[0].ID()].OnChainSigning.X)
	require.Equal(t, originalKeyRing.OCR[ocr[0].ID()].OnChainSigning.Y, decryptedKeyRing.OCR[ocr[0].ID()].OnChainSigning.Y)
	require.Equal(t, originalKeyRing.OCR[ocr[0].ID()].OnChainSigning.D, decryptedKeyRing.OCR[ocr[0].ID()].OnChainSigning.D)
	require.Equal(t, originalKeyRing.OCR[ocr[0].ID()].OffChainSigning, decryptedKeyRing.OCR[ocr[0].ID()].OffChainSigning)
	require.Equal(t, originalKeyRing.OCR[ocr[0].ID()].OffChainEncryption, decryptedKeyRing.OCR[ocr[0].ID()].OffChainEncryption)
	require.Equal(t, originalKeyRing.OCR[ocr[1].ID()].OnChainSigning.X, decryptedKeyRing.OCR[ocr[1].ID()].OnChainSigning.X)
	require.Equal(t, originalKeyRing.OCR[ocr[1].ID()].OnChainSigning.Y, decryptedKeyRing.OCR[ocr[1].ID()].OnChainSigning.Y)
	require.Equal(t, originalKeyRing.OCR[ocr[1].ID()].OnChainSigning.D, decryptedKeyRing.OCR[ocr[1].ID()].OnChainSigning.D)
	require.Equal(t, originalKeyRing.OCR[ocr[1].ID()].OffChainSigning, decryptedKeyRing.OCR[ocr[1].ID()].OffChainSigning)
	require.Equal(t, originalKeyRing.OCR[ocr[1].ID()].OffChainEncryption, decryptedKeyRing.OCR[ocr[1].ID()].OffChainEncryption)
	// compare ocr2 keys
	require.Equal(t, len(chaintype.SupportedChainTypes), len(decryptedKeyRing.OCR2))
	for i := range ocr2 {
		id := ocr2[i].ID()
		require.Equal(t, originalKeyRing.OCR2[id].ID(), decryptedKeyRing.OCR2[id].ID())
		require.Equal(t, ocr2[i].OnChainPublicKey(), decryptedKeyRing.OCR2[id].OnChainPublicKey())
		require.Equal(t, originalKeyRing.OCR2[id].ChainType(), decryptedKeyRing.OCR2[id].ChainType())
	}
	// compare p2p keys
	require.Equal(t, 2, len(decryptedKeyRing.P2P))
	require.Equal(t, originalKeyRing.P2P[p2p1.ID()].GetPublic(), decryptedKeyRing.P2P[p2p1.ID()].GetPublic())
	require.Equal(t, originalKeyRing.P2P[p2p1.ID()].PeerID(), decryptedKeyRing.P2P[p2p1.ID()].PeerID())
	require.Equal(t, originalKeyRing.P2P[p2p2.ID()].GetPublic(), decryptedKeyRing.P2P[p2p2.ID()].GetPublic())
	require.Equal(t, originalKeyRing.P2P[p2p2.ID()].PeerID(), decryptedKeyRing.P2P[p2p2.ID()].PeerID())
	// compare solana keys
	require.Equal(t, 2, len(decryptedKeyRing.Solana))
	require.Equal(t, originalKeyRing.Solana[sol1.ID()].GetPublic(), decryptedKeyRing.Solana[sol1.ID()].GetPublic())
	// compare vrf keys
	require.Equal(t, 2, len(decryptedKeyRing.VRF))
	require.Equal(t, originalKeyRing.VRF[vrf1.ID()].PublicKey, decryptedKeyRing.VRF[vrf1.ID()].PublicKey)
	require.Equal(t, originalKeyRing.VRF[vrf2.ID()].PublicKey, decryptedKeyRing.VRF[vrf2.ID()].PublicKey)
	// compare terra keys
	require.Equal(t, 2, len(decryptedKeyRing.Terra))
	require.Equal(t, originalKeyRing.Terra[tk1.ID()].PublicKey(), decryptedKeyRing.Terra[tk1.ID()].PublicKey())
	require.Equal(t, originalKeyRing.Terra[tk2.ID()].PublicKey(), decryptedKeyRing.Terra[tk2.ID()].PublicKey())
	// compare dkgsign keys
	require.Equal(t, 2, len(decryptedKeyRing.DKGSign))
	require.Equal(t, originalKeyRing.DKGSign[dkgsign1.ID()].PublicKey, decryptedKeyRing.DKGSign[dkgsign1.ID()].PublicKey)
	require.Equal(t, originalKeyRing.DKGSign[dkgsign2.ID()].PublicKey, decryptedKeyRing.DKGSign[dkgsign2.ID()].PublicKey)
	// compare dkgencrypt keys
	require.Equal(t, 2, len(decryptedKeyRing.DKGEncrypt))
	require.Equal(t, originalKeyRing.DKGEncrypt[dkgencrypt1.ID()].PublicKey, decryptedKeyRing.DKGEncrypt[dkgencrypt1.ID()].PublicKey)
	require.Equal(t, originalKeyRing.DKGEncrypt[dkgencrypt2.ID()].PublicKey, decryptedKeyRing.DKGEncrypt[dkgencrypt2.ID()].PublicKey)
}
