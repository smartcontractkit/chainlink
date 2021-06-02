package offchainreporting_test

import (
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/stretchr/testify/require"
)

func Test_SingletonPeerWrapper_Start(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	db := store.DB

	t.Run("with locked KeyStore returns nil", func(t *testing.T) {
		keyStore := cltest.NewKeyStore(t, db).OCR
		pw := offchainreporting.NewSingletonPeerWrapper(keyStore, store.Config, store.DB)

		require.NoError(t, pw.Start())
	})

	// Clear out fixture
	require.NoError(t, db.Exec(`DELETE FROM encrypted_p2p_keys`).Error)

	t.Run("with no p2p keys returns nil", func(t *testing.T) {
		keyStore := cltest.NewKeyStore(t, db).OCR
		require.NoError(t, keyStore.Unlock(cltest.Password))
		pw := offchainreporting.NewSingletonPeerWrapper(keyStore, store.Config, store.DB)

		require.NoError(t, pw.Start())
	})

	var k p2pkey.Key
	var err error

	t.Run("with one p2p key and matching P2P_PEER_ID returns nil", func(t *testing.T) {
		keyStore := cltest.NewKeyStore(t, db).OCR
		require.NoError(t, keyStore.Unlock(cltest.Password))
		k, _, err = keyStore.GenerateEncryptedP2PKey()
		require.NoError(t, err)

		store.Config.Set("P2P_PEER_ID", k.MustGetPeerID())

		require.NoError(t, err)

		pw := offchainreporting.NewSingletonPeerWrapper(keyStore, store.Config, store.DB)

		require.NoError(t, pw.Start(), "foo")
		require.Equal(t, k.MustGetPeerID(), pw.PeerID)
	})

	t.Run("with one p2p key and no P2P_PEER_ID returns error", func(t *testing.T) {
		keyStore := cltest.NewKeyStore(t, db).OCR
		require.NoError(t, keyStore.Unlock(cltest.Password))

		store.Config.Set("P2P_PEER_ID", "")

		pw := offchainreporting.NewSingletonPeerWrapper(keyStore, store.Config, store.DB)

		require.EqualError(t, pw.Start(), "failed to start peer wrapper: P2P_PEER_ID: env var unset")
	})

	t.Run("with one p2p key and mismatching P2P_PEER_ID returns error", func(t *testing.T) {
		keyStore := cltest.NewKeyStore(t, db).OCR
		require.NoError(t, keyStore.Unlock(cltest.Password))

		store.Config.Set("P2P_PEER_ID", cltest.DefaultP2PPeerID)

		pw := offchainreporting.NewSingletonPeerWrapper(keyStore, store.Config, store.DB)

		require.EqualError(t, pw.Start(), fmt.Sprintf("multiple p2p keys found but none matched the given P2P_PEER_ID of '%v'. Keys available: %s", cltest.DefaultP2PPeerID, k.MustGetPeerID()))
	})

	var k2 p2pkey.Key

	t.Run("with multiple p2p keys and valid P2P_PEER_ID returns nil", func(t *testing.T) {
		keyStore := cltest.NewKeyStore(t, db).OCR
		require.NoError(t, keyStore.Unlock(cltest.Password))
		k2, _, err = keyStore.GenerateEncryptedP2PKey()
		require.NoError(t, err)

		store.Config.Set("P2P_PEER_ID", k2.MustGetPeerID())

		require.NoError(t, err)

		pw := offchainreporting.NewSingletonPeerWrapper(keyStore, store.Config, store.DB)

		require.NoError(t, pw.Start(), "foo")
		require.Equal(t, k2.MustGetPeerID(), pw.PeerID)
	})

	t.Run("with multiple p2p keys and mismatching P2P_PEER_ID returns error", func(t *testing.T) {
		keyStore := cltest.NewKeyStore(t, db).OCR
		require.NoError(t, keyStore.Unlock(cltest.Password))

		store.Config.Set("P2P_PEER_ID", cltest.DefaultP2PPeerID.String())

		pw := offchainreporting.NewSingletonPeerWrapper(keyStore, store.Config, store.DB)

		require.Contains(t, pw.Start().Error(), "multiple p2p keys found but none matched the given P2P_PEER_ID of")
	})
}
