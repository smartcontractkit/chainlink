package ocrcommon_test

import (
	"sync"
	"testing"
	"time"

	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func Test_SingletonPeerWrapper_Start(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)

	require.NoError(t, utils.JustError(db.Exec(`DELETE FROM encrypted_key_rings`)))

	peerID, err := p2ppeer.Decode("12D3KooWPjceQrSwdWXPyLLeABRXmuqt69Rg3sBYbU1Nft9HyQ6X")
	require.NoError(t, err)

	t.Run("with no p2p keys returns error", func(t *testing.T) {
		cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.P2P.V1.Enabled = ptr(true)
		})
		keyStore := cltest.NewKeyStore(t, db, cfg)
		noopMetrics := ocrcommon.NewMetricVecFactory(ocrcommon.NewNoopMetricVec)
		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg, db, logger.TestLogger(t), noopMetrics)
		require.Contains(t, pw.Start(testutils.Context(t)).Error(), "No P2P keys found in keystore. Peer wrapper will not be fully initialized")
	})

	t.Run("with one p2p key and matching P2P_PEER_ID returns nil", func(t *testing.T) {
		cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.P2P.V1.Enabled = ptr(true)
		})
		keyStore := cltest.NewKeyStore(t, db, cfg)
		k, err := keyStore.P2P().Create()
		require.NoError(t, err)

		cfg = configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.P2P.V1.Enabled = ptr(true)
			c.P2P.PeerID = ptr(k.PeerID())
		})
		keyStore = cltest.NewKeyStore(t, db, cfg)
		noopMetrics := ocrcommon.NewMetricVecFactory(ocrcommon.NewNoopMetricVec)
		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg, db, logger.TestLogger(t), noopMetrics)

		require.NoError(t, pw.Start(testutils.Context(t)), "foo")
		require.Equal(t, k.PeerID(), pw.PeerID)
	})

	t.Run("with one p2p key and mismatching P2P_PEER_ID returns error", func(t *testing.T) {
		cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.P2P.V1.Enabled = ptr(true)
			c.P2P.PeerID = ptr(p2pkey.PeerID(peerID))
		})
		keyStore := cltest.NewKeyStore(t, db, cfg)
		noopMetrics := ocrcommon.NewMetricVecFactory(ocrcommon.NewNoopMetricVec)
		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg, db, logger.TestLogger(t), noopMetrics)

		require.Contains(t, pw.Start(testutils.Context(t)).Error(), "unable to find P2P key with id")
	})

	t.Run("with multiple p2p keys and valid P2P_PEER_ID returns nil", func(t *testing.T) {
		cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.P2P.V1.Enabled = ptr(true)
		})
		keyStore := cltest.NewKeyStore(t, db, cfg)
		k2, err := keyStore.P2P().Create()
		require.NoError(t, err)

		cfg = configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.P2P.V1.Enabled = ptr(true)
			c.P2P.PeerID = ptr(k2.PeerID())
		})
		keyStore = cltest.NewKeyStore(t, db, cfg)
		noopMetrics := ocrcommon.NewMetricVecFactory(ocrcommon.NewNoopMetricVec)
		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg, db, logger.TestLogger(t), noopMetrics)

		require.NoError(t, pw.Start(testutils.Context(t)), "foo")
		require.Equal(t, k2.PeerID(), pw.PeerID)
	})

	t.Run("with multiple p2p keys and mismatching P2P_PEER_ID returns error", func(t *testing.T) {
		cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.P2P.V1.Enabled = ptr(true)
			c.P2P.PeerID = ptr(p2pkey.PeerID(peerID))
		})
		keyStore := cltest.NewKeyStore(t, db, cfg)
		noopMetrics := ocrcommon.NewMetricVecFactory(ocrcommon.NewNoopMetricVec)
		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg, db, logger.TestLogger(t), noopMetrics)

		require.Contains(t, pw.Start(testutils.Context(t)).Error(), "unable to find P2P key with id")
	})

	t.Run("ocr metric instationation", func(t *testing.T) {
		cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.P2P.V2.Enabled = ptr(true)
			p2paddresses := []string{
				"127.0.0.1:17193",
			}
			c.P2P.V2.ListenAddresses = ptr(p2paddresses)
		})
		keyStore := cltest.NewKeyStore(t, db, cfg)
		k2, err := keyStore.P2P().Create()
		require.NoError(t, err)

		cfg = configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.P2P.V2.Enabled = ptr(true)
			p2paddresses := []string{
				"127.0.0.1:17195",
			}
			c.P2P.PeerID = ptr(k2.PeerID())
			c.P2P.V2.ListenAddresses = ptr(p2paddresses)
		})
		keyStore = cltest.NewKeyStore(t, db, cfg)

		// setup wrapper around noop metrics to ensure that libocr instantiates the metric handler
		var mu sync.RWMutex
		mck := &mock.Mock{}
		type mockGenerator struct {
			commontypes.MetricVec
		}

		// will be called multiple times in libocr, once per metric.
		newMockGenerator := func(name string, help string, labelNames ...string) (commontypes.MetricVec, error) {
			n, err := ocrcommon.NewNoopMetricVec(name, help, labelNames...)
			if err != nil {
				return nil, err
			}
			m := &mockGenerator{
				MetricVec: n,
			}
			mu.Lock()
			mck.On("generator").Return(nil)
			mck.MethodCalled("generator")
			mu.Unlock()
			return m, nil
		}

		mockMetrics := ocrcommon.NewMetricVecFactory(newMockGenerator)
		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg, db, logger.TestLogger(t), mockMetrics)

		require.NoError(t, pw.Start(testutils.Context(t)), "foo")

		testutils.AssertEventually(t, func() bool {
			mu.RLock()
			defer mu.RUnlock()
			return len(mck.Calls) > 0
		})
	})
}

func Test_SingletonPeerWrapper_Close(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)

	require.NoError(t, utils.JustError(db.Exec(`DELETE FROM encrypted_key_rings`)))

	cfg := configtest.NewGeneralConfig(t, nil)
	keyStore := cltest.NewKeyStore(t, db, cfg)
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
	keyStore = cltest.NewKeyStore(t, db, cfg)
	noopMetrics := ocrcommon.NewMetricVecFactory(ocrcommon.NewNoopMetricVec)
	pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg, db, logger.TestLogger(t), noopMetrics)

	require.NoError(t, pw.Start(testutils.Context(t)))
	require.True(t, pw.IsStarted(), "Should have started successfully")
	pw.Close()

	/* If peer is still stuck in listenLoop, we will get a bind error trying to start on the same port */
	require.False(t, pw.IsStarted())
	pw = ocrcommon.NewSingletonPeerWrapper(keyStore, cfg, db, logger.TestLogger(t), noopMetrics)
	require.NoError(t, pw.Start(testutils.Context(t)), "Should have shut down gracefully, and be able to re-use same port")
	require.True(t, pw.IsStarted(), "Should have started successfully")
	pw.Close()
}

func ptr[T any](t T) *T { return &t }
