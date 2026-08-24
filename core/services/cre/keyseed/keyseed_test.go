package keyseed

import (
	"strings"
	"testing"

	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonkeystore "github.com/smartcontractkit/chainlink-common/keystore"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ethkey"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"
	"github.com/smartcontractkit/chainlink-common/keystore/ocr2offchain"
	"github.com/smartcontractkit/chainlink-common/keystore/ragep2p"
)

const password = "a keystore password"

// store is an empty keystore in memory: what this tests is the conversion, and the
// database it normally lives in adds nothing to that.
func store(t *testing.T) commonkeystore.Keystore {
	t.Helper()
	ks, err := commonkeystore.LoadKeystore(t.Context(), commonkeystore.NewMemoryStorage(), password)
	require.NoError(t, err)
	return ks
}

// TestPeerKey is the test that matters for the peer: the copy has to be the same
// identity, because the address other DON members dial is this peer ID and the
// registry lists the node under it. A copy that produced a different key would
// look like a node that had quietly left the DON.
func TestPeerKey(t *testing.T) {
	legacy, err := p2pkey.NewV2()
	require.NoError(t, err)

	copied, err := peerKey(legacy, password, Names{})
	require.NoError(t, err)
	assert.Equal(t, "ragep2p_peer/"+DefaultNames.Peer, copied.name)
	assert.Equal(t, commonkeystore.Ed25519, copied.keyType)

	ks := store(t)
	n, err := imported(t.Context(), ks, password, copied)
	require.NoError(t, err)
	assert.Equal(t, 1, n)

	keyrings, err := ragep2p.GetPeerKeyrings(t.Context(), ks, []string{DefaultNames.Peer})
	require.NoError(t, err)
	require.Len(t, keyrings, 1)

	peerID, err := keyrings[0].PeerID()
	require.NoError(t, err)
	// The legacy form prints with a p2p_ prefix; the identity is what is after it.
	assert.Equal(t, strings.TrimPrefix(legacy.PeerID().String(), "p2p_"), peerID, "the copied peer must be the same peer")
}

// TestOCR2Keys checks the bundle arrives as the three keys the reading side looks
// for, and that each is the same key: the DON's configuration lists these public
// halves, so a copy that changed any of them would sign rounds nobody accepts.
func TestOCR2Keys(t *testing.T) {
	bundle, err := ocr2key.New(corekeys.EVM)
	require.NoError(t, err)

	keys, err := ocr2Keys(bundle, password, Names{})
	require.NoError(t, err)
	require.Len(t, keys, 3)

	ks := store(t)
	n, err := imported(t.Context(), ks, password, keys...)
	require.NoError(t, err)
	assert.Equal(t, 3, n)

	t.Run("the offchain keys are the bundle's", func(t *testing.T) {
		keyrings, err := ocr2offchain.GetOCR2OffchainKeyrings(t.Context(), ks, []string{DefaultNames.OCR})
		require.NoError(t, err)
		require.Len(t, keyrings, 1)

		assert.Equal(t, bundle.OffchainPublicKey(), keyrings[0].OffchainPublicKey())
		assert.Equal(t, bundle.ConfigEncryptionPublicKey(), keyrings[0].ConfigEncryptionPublicKey())
	})

	t.Run("the onchain key is the bundle's", func(t *testing.T) {
		onchain, err := ks.GetKeys(t.Context(), commonkeystore.GetKeysRequest{
			KeyNames: []string{path(prefixOCR2Onchain, DefaultNames.OCR, onchainSigning)},
		})
		require.NoError(t, err)
		require.Len(t, onchain.Keys, 1)
		// Pinned: the reading side spells this name too, in another repository (the
		// capabilities repo's libs/standalone/nodekeys).
		require.Equal(t, "ocr2_onchain/ocr2/ocr2_onchain_signing", path(prefixOCR2Onchain, DefaultNames.OCR, onchainSigning))

		// The bundle reports its onchain public key as an address, so compare addresses.
		address := gethcrypto.Keccak256(onchain.Keys[0].KeyInfo.PublicKey[1:])[12:]
		assert.Equal(t, []byte(bundle.PublicKey()), address, "the copied key must sign as the same account")
	})
}

// TestEVMKey checks an EVM key is copied under the name a chain asks for it by -
// its address - and that it is the same key.
func TestEVMKey(t *testing.T) {
	legacy, err := ethkey.NewV2()
	require.NoError(t, err)

	copied, err := evmKey(legacy, password)
	require.NoError(t, err)
	assert.Equal(t, legacy.Address.Hex(), copied.name, "a chain asks for an account by address")

	ks := store(t)
	n, err := imported(t.Context(), ks, password, copied)
	require.NoError(t, err)
	assert.Equal(t, 1, n)

	// Signed through the interface a chain uses, and recovered to the address the
	// node knows this key by.
	digest := gethcrypto.Keccak256([]byte("a transaction"))
	signature, err := commonkeystore.NewCoreKeystore(ks).Sign(t.Context(), legacy.Address.Hex(), digest)
	require.NoError(t, err)

	recovered, err := gethcrypto.SigToPub(digest, signature)
	require.NoError(t, err)
	assert.Equal(t, legacy.Address, gethcrypto.PubkeyToAddress(*recovered))
}

// TestImportedIsIdempotent covers the property that lets this run on every boot: a
// key already in the store is left as it is.
func TestImportedIsIdempotent(t *testing.T) {
	legacy, err := p2pkey.NewV2()
	require.NoError(t, err)
	copied, err := peerKey(legacy, password, Names{})
	require.NoError(t, err)

	ks := store(t)

	n, err := imported(t.Context(), ks, password, copied)
	require.NoError(t, err)
	require.Equal(t, 1, n)

	n, err = imported(t.Context(), ks, password, copied)
	require.NoError(t, err)
	assert.Zero(t, n, "a second boot must copy nothing")
}
