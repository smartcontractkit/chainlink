package offchainreporting2_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	offchainreporting "github.com/smartcontractkit/chainlink/core/services/offchainreporting2"
	"github.com/stretchr/testify/require"
)

func Test_SingletonPeerWrapper_Start(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	cfg := configtest.NewTestGeneralConfig(t)
	cfg.Overrides.OCR2P2PV2ListenAddresses = []string{"127.0.0.1:9999"}
	cfg.Overrides.OCR2P2PV2AnnounceAddresses = []string{"127.0.0.1:9999"}
	db := store.DB

	// Clear out fixture
	require.NoError(t, db.Exec(`DELETE FROM encrypted_key_rings`).Error)

	t.Run("with no p2p keys returns nil", func(t *testing.T) {
		keyStore := cltest.NewKeyStore(t, db)
		pw := offchainreporting.NewSingletonPeerWrapper(keyStore, cfg, store.DB)

		require.NoError(t, pw.Start())
	})

	var k p2pkey.KeyV2
	var err error

	t.Run("with one p2p key and matching OCR2_P2P_PEER_ID returns nil", func(t *testing.T) {
		keyStore := cltest.NewKeyStore(t, db)
		k, err = keyStore.P2P().Create()
		require.NoError(t, err)

		peerID := k.PeerID()
		cfg.Overrides.OCR2P2PPeerID = peerID

		require.NoError(t, err)

		pw := offchainreporting.NewSingletonPeerWrapper(keyStore, cfg, store.DB)

		require.NoError(t, pw.Start(), "foo")
		require.Equal(t, k.PeerID(), pw.PeerID)
	})

	t.Run("with multiple p2p keys and mismatching OCR2_P2P_PEER_ID returns error", func(t *testing.T) {
		keyStore := cltest.NewKeyStore(t, db)
		cfg.Overrides.OCR2P2PPeerID = cltest.DefaultP2PPeerID

		pw := offchainreporting.NewSingletonPeerWrapper(keyStore, cfg, store.DB)

		require.Contains(t, pw.Start().Error(), "multiple p2p keys found but none matched the given OCR2_P2P_PEER_ID of")
	})
}
