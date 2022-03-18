package ocrcommon_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	"gopkg.in/guregu/null.v4"

	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func Test_SingletonPeerWrapper_Start(t *testing.T) {
	t.Parallel()

	cfg := configtest.NewTestGeneralConfigWithOverrides(t, configtest.GeneralConfigOverrides{
		P2PEnabled: null.BoolFrom(true),
	})
	db := pgtest.NewSqlxDB(t)

	require.NoError(t, utils.JustError(db.Exec(`DELETE FROM encrypted_key_rings`)))

	peerID, err := p2ppeer.Decode("12D3KooWPjceQrSwdWXPyLLeABRXmuqt69Rg3sBYbU1Nft9HyQ6X")
	require.NoError(t, err)

	t.Run("with no p2p keys returns error", func(t *testing.T) {
		keyStore := cltest.NewKeyStore(t, db, cfg)
		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg, db, logger.TestLogger(t))
		require.Contains(t, pw.Start(testutils.Context(t)).Error(), "No P2P keys found in keystore. Peer wrapper will not be fully initialized")
	})

	var k p2pkey.KeyV2

	t.Run("with one p2p key and matching P2P_PEER_ID returns nil", func(t *testing.T) {
		keyStore := cltest.NewKeyStore(t, db, cfg)
		k, err = keyStore.P2P().Create()
		require.NoError(t, err)

		cfg.Overrides.P2PPeerID = k.PeerID()

		require.NoError(t, err)

		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg, db, logger.TestLogger(t))

		require.NoError(t, pw.Start(testutils.Context(t)), "foo")
		require.Equal(t, k.PeerID(), pw.PeerID)
	})

	t.Run("with one p2p key and mismatching P2P_PEER_ID returns error", func(t *testing.T) {
		keyStore := cltest.NewKeyStore(t, db, cfg)

		cfg.Overrides.P2PPeerID = p2pkey.PeerID(peerID)

		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg, db, logger.TestLogger(t))

		require.Contains(t, pw.Start(testutils.Context(t)).Error(), "unable to find P2P key with id")
	})

	var k2 p2pkey.KeyV2

	t.Run("with multiple p2p keys and valid P2P_PEER_ID returns nil", func(t *testing.T) {
		keyStore := cltest.NewKeyStore(t, db, cfg)
		k2, err = keyStore.P2P().Create()
		require.NoError(t, err)

		cfg.Overrides.P2PPeerID = k2.PeerID()

		require.NoError(t, err)

		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg, db, logger.TestLogger(t))

		require.NoError(t, pw.Start(testutils.Context(t)), "foo")
		require.Equal(t, k2.PeerID(), pw.PeerID)
	})

	t.Run("with multiple p2p keys and mismatching P2P_PEER_ID returns error", func(t *testing.T) {
		keyStore := cltest.NewKeyStore(t, db, cfg)

		cfg.Overrides.P2PPeerID = p2pkey.PeerID(peerID)

		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg, db, logger.TestLogger(t))

		require.Contains(t, pw.Start(testutils.Context(t)).Error(), "unable to find P2P key with id")
	})
}
