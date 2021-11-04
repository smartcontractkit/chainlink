package offchainreporting_test

import (
	"fmt"
	"testing"

	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/stretchr/testify/require"
)

func Test_SingletonPeerWrapper_Start(t *testing.T) {
	t.Parallel()

	cfg := configtest.NewTestGeneralConfig(t)
	gdb := pgtest.NewGormDB(t)
	db := postgres.UnwrapGormDB(gdb)

	require.NoError(t, gdb.Exec(`DELETE FROM encrypted_key_rings`).Error)

	peerID, err := p2ppeer.Decode("12D3KooWPjceQrSwdWXPyLLeABRXmuqt69Rg3sBYbU1Nft9HyQ6X")
	require.NoError(t, err)

	t.Run("with no p2p keys returns nil", func(t *testing.T) {
		keyStore := cltest.NewKeyStore(t, db)
		pw := offchainreporting.NewSingletonPeerWrapper(keyStore, cfg, gdb, logger.TestLogger(t))

		require.NoError(t, pw.Start())
	})

	var k p2pkey.KeyV2

	t.Run("with one p2p key and matching P2P_PEER_ID returns nil", func(t *testing.T) {
		keyStore := cltest.NewKeyStore(t, db)
		k, err = keyStore.P2P().Create()
		require.NoError(t, err)

		cfg.Overrides.P2PPeerID = k.PeerID()

		require.NoError(t, err)

		pw := offchainreporting.NewSingletonPeerWrapper(keyStore, cfg, gdb, logger.TestLogger(t))

		require.NoError(t, pw.Start(), "foo")
		require.Equal(t, k.PeerID(), pw.PeerID)
	})

	t.Run("with one p2p key and mismatching P2P_PEER_ID returns error", func(t *testing.T) {
		keyStore := cltest.NewKeyStore(t, db)

		cfg.Overrides.P2PPeerID = p2pkey.PeerID(peerID)

		pw := offchainreporting.NewSingletonPeerWrapper(keyStore, cfg, gdb, logger.TestLogger(t))

		require.Contains(t, pw.Start().Error(), fmt.Sprintf("unable to find P2P key with id %s", p2pkey.PeerID(peerID)))
	})

	var k2 p2pkey.KeyV2

	t.Run("with multiple p2p keys and valid P2P_PEER_ID returns nil", func(t *testing.T) {
		keyStore := cltest.NewKeyStore(t, db)
		k2, err = keyStore.P2P().Create()
		require.NoError(t, err)

		cfg.Overrides.P2PPeerID = k2.PeerID()

		require.NoError(t, err)

		pw := offchainreporting.NewSingletonPeerWrapper(keyStore, cfg, gdb, logger.TestLogger(t))

		require.NoError(t, pw.Start(), "foo")
		require.Equal(t, k2.PeerID(), pw.PeerID)
	})

	t.Run("with multiple p2p keys and mismatching P2P_PEER_ID returns error", func(t *testing.T) {
		keyStore := cltest.NewKeyStore(t, db)

		cfg.Overrides.P2PPeerID = p2pkey.PeerID(peerID)

		pw := offchainreporting.NewSingletonPeerWrapper(keyStore, cfg, gdb, logger.TestLogger(t))

		require.Contains(t, pw.Start().Error(), fmt.Sprintf("unable to find P2P key with id %s", p2pkey.PeerID(peerID)))
	})
}
