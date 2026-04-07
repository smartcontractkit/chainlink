package cre

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/secrets"
	"github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
)

func TestNewNodeKeys_IgnoresEmptyImportedAptosSecretWhenAptosDisabled(t *testing.T) {
	t.Parallel()

	p2pKey, err := crypto.NewP2PKey("dev-password")
	require.NoError(t, err)
	dkgKey, err := crypto.NewDKGRecipientKey("dev-password")
	require.NoError(t, err)

	baseSecrets, err := (&secrets.NodeKeys{
		P2PKey: p2pKey,
		DKGKey: dkgKey,
		EVM:    map[uint64]*crypto.EVMKey{},
		Solana: map[string]*crypto.SolKey{},
	}).ToNodeSecretsTOML()
	require.NoError(t, err)

	keys, err := NewNodeKeys(NodeKeyInput{
		ImportedSecrets: baseSecrets,
		AptosChainIDs:   nil,
	})
	require.NoError(t, err)
	require.Nil(t, keys.Aptos)
	require.Equal(t, p2pKey.PeerID, keys.P2PKey.PeerID)
}

func TestNewNodeKeys_RejectsMissingImportedAptosSecretWhenAptosEnabled(t *testing.T) {
	t.Parallel()

	p2pKey, err := crypto.NewP2PKey("dev-password")
	require.NoError(t, err)
	dkgKey, err := crypto.NewDKGRecipientKey("dev-password")
	require.NoError(t, err)

	baseSecrets, err := (&secrets.NodeKeys{
		P2PKey: p2pKey,
		DKGKey: dkgKey,
		EVM:    map[uint64]*crypto.EVMKey{},
		Solana: map[string]*crypto.SolKey{},
	}).ToNodeSecretsTOML()
	require.NoError(t, err)

	_, err = NewNodeKeys(NodeKeyInput{
		ImportedSecrets: baseSecrets,
		AptosChainIDs:   []uint64{4},
	})
	require.ErrorContains(t, err, "missing an Aptos key")
}

func TestNewNodeKeys_PreservesAptosChainIDs(t *testing.T) {
	t.Parallel()

	keys, err := NewNodeKeys(NodeKeyInput{
		AptosChainIDs: []uint64{4, 5},
		Password:      "dev-password",
	})
	require.NoError(t, err)
	require.NotNil(t, keys.Aptos)
	require.ElementsMatch(t, []uint64{4, 5}, keys.AptosChainIDs)
}
