package ocrcommon_test

import (
	"fmt"
	"testing"
	"time"

	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/freeport"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
)

func Test_SingletonPeerWrapper_Start(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)

	var peerID ragep2ptypes.PeerID
	require.NoError(t, peerID.UnmarshalText([]byte(configtest.DefaultPeerID)))

	t.Run("with no p2p keys returns error", func(t *testing.T) {
		cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.P2P.V2.Enabled = ptr(true)
		})
		keyStore := cltest.NewKeyStore(t, db)
		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg.P2P(), cfg.OCR(), db, logger.TestLogger(t))
		require.Contains(t, pw.Start(testutils.Context(t)).Error(), "No P2P keys found in keystore. Peer wrapper will not be fully initialized")
	})

	t.Run("with one p2p key and matching P2P.PeerID returns nil", func(t *testing.T) {
		ctx := testutils.Context(t)
		keyStore := cltest.NewKeyStore(t, db)
		k, err := keyStore.P2P().Create(ctx)
		require.NoError(t, err)

		cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.P2P.V2.Enabled = ptr(true)
			c.P2P.V2.ListenAddresses = &[]string{fmt.Sprintf("127.0.0.1:%d", freeport.GetOne(t))}
			c.P2P.PeerID = ptr(k.PeerID())
		})
		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg.P2P(), cfg.OCR(), db, logger.TestLogger(t))

		servicetest.Run(t, pw)
		require.Equal(t, k.PeerID(), pw.PeerID)
	})

	t.Run("with one p2p key and mismatching P2P.PeerID returns error", func(t *testing.T) {
		ctx := testutils.Context(t)
		cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.P2P.V2.Enabled = ptr(true)
			c.P2P.PeerID = ptr(p2pkey.PeerID(peerID))
		})
		keyStore := cltest.NewKeyStore(t, db)

		_, err := keyStore.P2P().Create(ctx)
		require.NoError(t, err)

		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg.P2P(), cfg.OCR(), db, logger.TestLogger(t))

		require.Contains(t, pw.Start(testutils.Context(t)).Error(), "unable to find P2P key with id")
	})

	t.Run("with multiple p2p keys and valid P2P.PeerID returns nil", func(t *testing.T) {
		ctx := testutils.Context(t)
		keyStore := cltest.NewKeyStore(t, db)
		k2, err := keyStore.P2P().Create(ctx)
		require.NoError(t, err)

		cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.P2P.V2.Enabled = ptr(true)
			c.P2P.V2.ListenAddresses = &[]string{fmt.Sprintf("127.0.0.1:%d", freeport.GetOne(t))}
			c.P2P.PeerID = ptr(k2.PeerID())
		})

		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg.P2P(), cfg.OCR(), db, logger.TestLogger(t))

		servicetest.Run(t, pw)
		require.Equal(t, k2.PeerID(), pw.PeerID)
	})

	t.Run("with multiple p2p keys and mismatching P2P.PeerID returns error", func(t *testing.T) {
		ctx := testutils.Context(t)
		cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.P2P.V2.Enabled = ptr(true)
			c.P2P.V2.ListenAddresses = &[]string{fmt.Sprintf("127.0.0.1:%d", freeport.GetOne(t))}
			c.P2P.PeerID = ptr(p2pkey.PeerID(peerID))
		})
		keyStore := cltest.NewKeyStore(t, db)

		_, err := keyStore.P2P().Create(ctx)
		require.NoError(t, err)

		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg.P2P(), cfg.OCR(), db, logger.TestLogger(t))

		require.Contains(t, pw.Start(testutils.Context(t)).Error(), "unable to find P2P key with id")
	})
}

func Test_SingletonPeerWrapper_Close(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	db := pgtest.NewSqlxDB(t)

	keyStore := cltest.NewKeyStore(t, db)
	k, err := keyStore.P2P().Create(ctx)
	require.NoError(t, err)

	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.P2P.V2.Enabled = ptr(true)
		c.P2P.PeerID = ptr(k.PeerID())
		c.P2P.V2.DeltaDial = commonconfig.MustNewDuration(100 * time.Millisecond)
		c.P2P.V2.DeltaReconcile = commonconfig.MustNewDuration(1 * time.Second)

		p2paddresses := []string{
			"127.0.0.1:17193",
		}
		c.P2P.V2.ListenAddresses = ptr(p2paddresses)
		c.P2P.V2.AnnounceAddresses = ptr(p2paddresses)
	})

	pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg.P2P(), cfg.OCR(), db, logger.TestLogger(t))

	require.NoError(t, pw.Start(testutils.Context(t)))
	require.True(t, pw.IsStarted(), "Should have started successfully")
	require.NoError(t, pw.Close())

	/* If peer is still stuck in listenLoop, we will get a bind error trying to start on the same port */
	require.False(t, pw.IsStarted())
	pw = ocrcommon.NewSingletonPeerWrapper(keyStore, cfg.P2P(), cfg.OCR(), db, logger.TestLogger(t))
	require.NoError(t, pw.Start(testutils.Context(t)), "Should have shut down gracefully, and be able to re-use same port")
	require.True(t, pw.IsStarted(), "Should have started successfully")
	require.NoError(t, pw.Close())
}

func ptr[T any](t T) *T { return &t }
