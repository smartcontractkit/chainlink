package keystore

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/cosmoskey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocrkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/solkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/tronkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
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
	ocr2Raw := make([][]byte, 0, len(chaintype.SupportedChainTypes))
	for _, chain := range chaintype.SupportedChainTypes {
		key := ocr2key.MustNewInsecure(rand.Reader, chain)
		ocr2 = append(ocr2, key)
		ocr2Raw = append(ocr2Raw, internal.RawBytes(key))
	}
	p2p1, p2p2 := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1)), p2pkey.MustNewV2XXXTestingOnly(big.NewInt(2))
	sol1, sol2 := solkey.MustNewInsecure(rand.Reader), solkey.MustNewInsecure(rand.Reader)
	vrf1, vrf2 := vrfkey.MustNewV2XXXTestingOnly(big.NewInt(1)), vrfkey.MustNewV2XXXTestingOnly(big.NewInt(2))
	tk1, tk2 := cosmoskey.MustNewInsecure(rand.Reader), cosmoskey.MustNewInsecure(rand.Reader)
	uk1, uk2 := tronkey.MustNewInsecure(rand.Reader), tronkey.MustNewInsecure(rand.Reader)
	originalKeyRingRaw := rawKeyRing{
		CSA:    [][]byte{internal.RawBytes(csa1), internal.RawBytes(csa2)},
		Eth:    [][]byte{internal.RawBytes(eth1), internal.RawBytes(eth2)},
		OCR:    [][]byte{internal.RawBytes(ocr[0]), internal.RawBytes(ocr[1])},
		OCR2:   ocr2Raw,
		P2P:    [][]byte{internal.RawBytes(p2p1), internal.RawBytes(p2p2)},
		Solana: [][]byte{internal.RawBytes(sol1), internal.RawBytes(sol2)},
		VRF:    [][]byte{internal.RawBytes(vrf1), internal.RawBytes(vrf2)},
		Cosmos: [][]byte{internal.RawBytes(tk1), internal.RawBytes(tk2)},
		Tron:   [][]byte{internal.RawBytes(uk1), internal.RawBytes(uk2)},
	}
	originalKeyRing, kerr := originalKeyRingRaw.keys()
	require.NoError(t, kerr)

	t.Run("test encrypt/decrypt", func(t *testing.T) {
		encryptedKr, err := originalKeyRing.Encrypt(password, utils.FastScryptParams)
		require.NoError(t, err)
		decryptedKeyRing, err := encryptedKr.Decrypt(password)
		require.NoError(t, err)
		// compare cosmos keys
		require.Len(t, decryptedKeyRing.Cosmos, 2)
		require.Equal(t, originalKeyRing.Cosmos[tk1.ID()].PublicKey(), decryptedKeyRing.Cosmos[tk1.ID()].PublicKey())
		require.Equal(t, originalKeyRing.Cosmos[tk2.ID()].PublicKey(), decryptedKeyRing.Cosmos[tk2.ID()].PublicKey())
		// compare tron keys
		require.Len(t, decryptedKeyRing.Tron, 2)
		require.Equal(t, originalKeyRing.Tron[uk1.ID()].Base58Address(), decryptedKeyRing.Tron[uk1.ID()].Base58Address())
		require.Equal(t, originalKeyRing.Tron[uk2.ID()].Base58Address(), decryptedKeyRing.Tron[uk2.ID()].Base58Address())
		// compare csa keys
		require.Len(t, decryptedKeyRing.CSA, 2)
		require.Equal(t, originalKeyRing.CSA[csa1.ID()].PublicKey, decryptedKeyRing.CSA[csa1.ID()].PublicKey)
		require.Equal(t, originalKeyRing.CSA[csa2.ID()].PublicKey, decryptedKeyRing.CSA[csa2.ID()].PublicKey)
		// compare eth keys
		require.Len(t, decryptedKeyRing.Eth, 2)
		require.Equal(t, originalKeyRing.Eth[eth1.ID()].Address, decryptedKeyRing.Eth[eth1.ID()].Address)
		require.Equal(t, originalKeyRing.Eth[eth2.ID()].Address, decryptedKeyRing.Eth[eth2.ID()].Address)
		// compare ocr keys
		require.Len(t, decryptedKeyRing.OCR, 2)
		require.Equal(t, internal.RawBytes(originalKeyRing.OCR[ocr[0].ID()]), internal.RawBytes(decryptedKeyRing.OCR[ocr[0].ID()]))
		require.Equal(t, originalKeyRing.OCR[ocr[0].ID()].OffChainEncryption, decryptedKeyRing.OCR[ocr[0].ID()].OffChainEncryption)
		require.Equal(t, internal.RawBytes(originalKeyRing.OCR[ocr[1].ID()]), internal.RawBytes(decryptedKeyRing.OCR[ocr[1].ID()]))
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
		require.Len(t, decryptedKeyRing.P2P, 2)
		require.Equal(t, originalKeyRing.P2P[p2p1.ID()].PublicKeyHex(), decryptedKeyRing.P2P[p2p1.ID()].PublicKeyHex())
		require.Equal(t, originalKeyRing.P2P[p2p1.ID()].PeerID(), decryptedKeyRing.P2P[p2p1.ID()].PeerID())
		require.Equal(t, originalKeyRing.P2P[p2p2.ID()].PublicKeyHex(), decryptedKeyRing.P2P[p2p2.ID()].PublicKeyHex())
		require.Equal(t, originalKeyRing.P2P[p2p2.ID()].PeerID(), decryptedKeyRing.P2P[p2p2.ID()].PeerID())
		// compare solana keys
		require.Len(t, decryptedKeyRing.Solana, 2)
		require.Equal(t, originalKeyRing.Solana[sol1.ID()].GetPublic(), decryptedKeyRing.Solana[sol1.ID()].GetPublic())
		// compare vrf keys
		require.Len(t, decryptedKeyRing.VRF, 2)
		require.Equal(t, originalKeyRing.VRF[vrf1.ID()].PublicKey, decryptedKeyRing.VRF[vrf1.ID()].PublicKey)
		require.Equal(t, originalKeyRing.VRF[vrf2.ID()].PublicKey, decryptedKeyRing.VRF[vrf2.ID()].PublicKey)
	})

	t.Run("test legacy system", func(t *testing.T) {
		// Add unsupported keys to raw json
		rawJson, _ := json.Marshal(originalKeyRing.raw())
		var allKeys = map[string][]string{
			"foo": {
				"bar", "biz",
			},
		}
		err := json.Unmarshal(rawJson, &allKeys)
		require.NoError(t, err)
		// Add more ocr2 keys
		newOCR2Key1 := ocrkey.MustNewV2XXXTestingOnly(big.NewInt(5))
		newOCR2Key2 := ocrkey.MustNewV2XXXTestingOnly(big.NewInt(6))
		allKeys["OCR2"] = append(allKeys["OCR2"], newOCR2Key1.Raw().String())
		allKeys["OCR2"] = append(allKeys["OCR2"], newOCR2Key2.Raw().String())

		// Add more p2p keys
		newP2PKey1 := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(5))
		newP2PKey2 := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(7))
		allKeys["P2P"] = append(allKeys["P2P"], newP2PKey1.Raw().String())
		allKeys["P2P"] = append(allKeys["P2P"], newP2PKey2.Raw().String())

		// Run legacy system
		newRawJson, _ := json.Marshal(allKeys)
		err = originalKeyRing.LegacyKeys.StoreUnsupported(newRawJson, originalKeyRing)
		require.NoError(t, err)
		require.Equal(t, 6, originalKeyRing.LegacyKeys.legacyRawKeys.len())
		marshalledRawKeyRingJson, err := json.Marshal(originalKeyRing.raw())
		require.NoError(t, err)
		unloadedKeysJson, err := originalKeyRing.LegacyKeys.UnloadUnsupported(marshalledRawKeyRingJson)
		require.NoError(t, err)
		var shouldHaveAllKeys = map[string][]string{}
		err = json.Unmarshal(unloadedKeysJson, &shouldHaveAllKeys)
		require.NoError(t, err)

		// Check if keys where added to the raw json
		require.Equal(t, []string{"bar", "biz"}, shouldHaveAllKeys["foo"])
		require.Contains(t, shouldHaveAllKeys["OCR2"], newOCR2Key1.Raw().String())
		require.Contains(t, shouldHaveAllKeys["OCR2"], newOCR2Key2.Raw().String())
		require.Contains(t, shouldHaveAllKeys["P2P"], newP2PKey1.Raw().String())
		require.Contains(t, shouldHaveAllKeys["P2P"], newP2PKey2.Raw().String())

		// Check error
		err = originalKeyRing.LegacyKeys.StoreUnsupported(newRawJson, nil)
		require.Error(t, err)
		_, err = originalKeyRing.LegacyKeys.UnloadUnsupported(nil)
		require.Error(t, err)
	})
}
