package ocrcommon_test

import (
	"testing"
	"time"

	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func Test_SingletonPeerWrapper_Start(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)

	peerID, err := p2ppeer.Decode("12D3KooWPjceQrSwdWXPyLLeABRXmuqt69Rg3sBYbU1Nft9HyQ6X")
	require.NoError(t, err)

	t.Run("with no p2p keys returns error", func(t *testing.T) {
		cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.P2P.V1.Enabled = ptr(true)
		})
		keyStore := cltest.NewKeyStore(t, db, cfg.Database())
		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg.P2P(), cfg.OCR(), cfg.Database(), db, logger.TestLogger(t))
		require.Contains(t, pw.Start(testutils.Context(t)).Error(), "No P2P keys found in keystore. Peer wrapper will not be fully initialized")
	})

	t.Run("with one p2p key and matching P2P.PeerID returns nil", func(t *testing.T) {
		cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.P2P.V1.Enabled = ptr(true)
		})
		keyStore := cltest.NewKeyStore(t, db, cfg.Database())
		k, err := keyStore.P2P().Create()
		require.NoError(t, err)

		cfg = configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.P2P.V1.Enabled = ptr(true)
			c.P2P.PeerID = ptr(k.PeerID())
		})
		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg.P2P(), cfg.OCR(), cfg.Database(), db, logger.TestLogger(t))

		require.NoError(t, pw.Start(testutils.Context(t)), "foo")
		require.Equal(t, k.PeerID(), pw.PeerID)
		require.NoError(t, pw.Close())
	})

	t.Run("with one p2p key and mismatching P2P.PeerID returns error", func(t *testing.T) {
		cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.P2P.V1.Enabled = ptr(true)
			c.P2P.PeerID = ptr(p2pkey.PeerID(peerID))
		})
		keyStore := cltest.NewKeyStore(t, db, cfg.Database())

		_, err := keyStore.P2P().Create()
		require.NoError(t, err)

		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg.P2P(), cfg.OCR(), cfg.Database(), db, logger.TestLogger(t))

		require.Contains(t, pw.Start(testutils.Context(t)).Error(), "unable to find P2P key with id")
	})

	t.Run("with multiple p2p keys and valid P2P.PeerID returns nil", func(t *testing.T) {
		cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.P2P.V1.Enabled = ptr(true)
		})
		keyStore := cltest.NewKeyStore(t, db, cfg.Database())
		k2, err := keyStore.P2P().Create()
		require.NoError(t, err)

		cfg = configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.P2P.V1.Enabled = ptr(true)
			c.P2P.PeerID = ptr(k2.PeerID())
		})

		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg.P2P(), cfg.OCR(), cfg.Database(), db, logger.TestLogger(t))

		require.NoError(t, pw.Start(testutils.Context(t)), "foo")
		require.Equal(t, k2.PeerID(), pw.PeerID)
		require.NoError(t, pw.Close())
	})

	t.Run("with multiple p2p keys and mismatching P2P.PeerID returns error", func(t *testing.T) {
		cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.P2P.V1.Enabled = ptr(true)
			c.P2P.PeerID = ptr(p2pkey.PeerID(peerID))
		})
		keyStore := cltest.NewKeyStore(t, db, cfg.Database())

		_, err := keyStore.P2P().Create()
		require.NoError(t, err)

		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg.P2P(), cfg.OCR(), cfg.Database(), db, logger.TestLogger(t))

		require.Contains(t, pw.Start(testutils.Context(t)).Error(), "unable to find P2P key with id")
	})
}

func Test_SingletonPeerWrapper_Close(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)

	cfg := configtest.NewGeneralConfig(t, nil)
	keyStore := cltest.NewKeyStore(t, db, cfg.Database())
	k, err := keyStore.P2P().Create()
	require.NoError(t, err)

	cfg = configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.P2P.V2.Enabled = ptr(true)
		c.P2P.PeerID = ptr(k.PeerID())
		c.P2P.V2.DeltaDial = models.MustNewDuration(100 * time.Millisecond)
		c.P2P.V2.DeltaReconcile = models.MustNewDuration(1 * time.Second)
		c.P2P.V1.ListenPort = ptr[uint16](0)

		p2paddresses := []string{
			"127.0.0.1:17193",
		}
		c.P2P.V2.ListenAddresses = ptr(p2paddresses)
		c.P2P.V2.AnnounceAddresses = ptr(p2paddresses)

	})

	pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg.P2P(), cfg.OCR(), cfg.Database(), db, logger.TestLogger(t))

	require.NoError(t, pw.Start(testutils.Context(t)))
	require.True(t, pw.IsStarted(), "Should have started successfully")
	require.NoError(t, pw.Close())

	/* If peer is still stuck in listenLoop, we will get a bind error trying to start on the same port */
	require.False(t, pw.IsStarted())
	pw = ocrcommon.NewSingletonPeerWrapper(keyStore, cfg.P2P(), cfg.OCR(), cfg.Database(), db, logger.TestLogger(t))
	require.NoError(t, pw.Start(testutils.Context(t)), "Should have shut down gracefully, and be able to re-use same port")
	require.True(t, pw.IsStarted(), "Should have started successfully")
	require.NoError(t, pw.Close())
}

func TestSingletonPeerWrapper_PeerConfig(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)

	require.NoError(t, utils.JustError(db.Exec(`DELETE FROM encrypted_key_rings`)))

	cfg := configtest.NewGeneralConfig(t, nil)
	keyStore := cltest.NewKeyStore(t, db, cfg.Database())
	k, err := keyStore.P2P().Create()
	require.NoError(t, err)

	t.Run("generates a random port if v1 enabled and listen port isn't set", func(t *testing.T) {
		cfg = configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.P2P.V1.Enabled = ptr(true)
			c.P2P.V1.ListenPort = ptr(uint16(0))
			c.P2P.PeerID = ptr(k.PeerID())
		})

		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg.P2P(), cfg.OCR(), cfg.Database(), db, logger.TestLogger(t))
		peerConfig, err := pw.PeerConfig()
		require.NoError(t, err)

		assert.NotEqual(t, peerConfig.V1ListenPort, 0)
	})

	t.Run("generates a random port if v1v2 enabled and listen port isn't set", func(t *testing.T) {
		cfg = configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.P2P.V1.Enabled = ptr(true)
			c.P2P.V2.Enabled = ptr(true)
			c.P2P.V1.ListenPort = ptr(uint16(0))
			c.P2P.PeerID = ptr(k.PeerID())
		})

		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg.P2P(), cfg.OCR(), cfg.Database(), db, logger.TestLogger(t))
		peerConfig, err := pw.PeerConfig()
		require.NoError(t, err)

		assert.NotEqual(t, peerConfig.V1ListenPort, 0)
	})

	t.Run("doesnt generate a port if v2 is enabled", func(t *testing.T) {
		cfg = configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.P2P.V2.Enabled = ptr(true)
			c.P2P.V1.ListenPort = ptr(uint16(0))
			c.P2P.PeerID = ptr(k.PeerID())
		})

		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg.P2P(), cfg.OCR(), cfg.Database(), db, logger.TestLogger(t))
		peerConfig, err := pw.PeerConfig()
		require.NoError(t, err)

		assert.NotEqual(t, peerConfig.V1ListenPort, 0)
	})

	t.Run("doesnt override a port if listenport is set", func(t *testing.T) {
		portNo := uint16(33247)
		cfg = configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.P2P.V1.Enabled = ptr(true)
			c.P2P.V1.ListenPort = ptr(uint16(33247))
			c.P2P.PeerID = ptr(k.PeerID())
		})

		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg.P2P(), cfg.OCR(), cfg.Database(), db, logger.TestLogger(t))
		peerConfig, err := pw.PeerConfig()
		require.NoError(t, err)

		assert.Equal(t, peerConfig.V1ListenPort, portNo)
	})
}

func ptr[T any](t T) *T { return &t }
