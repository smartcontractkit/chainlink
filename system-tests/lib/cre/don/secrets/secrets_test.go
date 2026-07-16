package secrets

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
)

func TestNodeKeysAptosSecretsRoundTrip(t *testing.T) {
	t.Parallel()

	aptosKey, err := crypto.NewAptosKey("dev-password")
	require.NoError(t, err)
	p2pKey, err := crypto.NewP2PKey("dev-password")
	require.NoError(t, err)
	dkgKey, err := crypto.NewDKGRecipientKey("dev-password")
	require.NoError(t, err)

	keys := &NodeKeys{
		Aptos:         aptosKey,
		AptosChainIDs: []uint64{4, 5},
		P2PKey:        p2pKey,
		DKGKey:        dkgKey,
		EVM:           map[uint64]*crypto.EVMKey{},
		Solana:        map[string]*crypto.SolKey{},
	}

	secretsTOML, err := keys.ToNodeSecretsTOML()
	require.NoError(t, err)
	require.Equal(t, 2, strings.Count(secretsTOML, "[[Aptos.Keys]]"))

	imported, err := ImportNodeKeys(secretsTOML)
	require.NoError(t, err)
	require.NotNil(t, imported.Aptos)
	require.ElementsMatch(t, []uint64{4, 5}, imported.AptosChainIDs)
	require.Equal(t, aptosKey.Account, imported.Aptos.Account)
	require.Equal(t, aptosKey.PublicKey, imported.Aptos.PublicKey)
	require.Equal(t, aptosKey.Password, imported.Aptos.Password)
	require.Equal(t, p2pKey.PeerID, imported.P2PKey.PeerID)
	require.Equal(t, dkgKey.PubKey, imported.DKGKey.PubKey)
}

func TestImportNodeKeys_IgnoresEmptyAptosSecret(t *testing.T) {
	t.Parallel()

	p2pKey, err := crypto.NewP2PKey("dev-password")
	require.NoError(t, err)
	dkgKey, err := crypto.NewDKGRecipientKey("dev-password")
	require.NoError(t, err)

	keys := &NodeKeys{
		P2PKey: p2pKey,
		DKGKey: dkgKey,
		EVM:    map[uint64]*crypto.EVMKey{},
		Solana: map[string]*crypto.SolKey{},
	}

	secretsTOML, err := keys.ToNodeSecretsTOML()
	require.NoError(t, err)

	imported, err := ImportNodeKeys(secretsTOML)
	require.NoError(t, err)
	require.Nil(t, imported.Aptos)
	require.Equal(t, p2pKey.PeerID, imported.P2PKey.PeerID)
	require.Equal(t, dkgKey.PubKey, imported.DKGKey.PubKey)
}
