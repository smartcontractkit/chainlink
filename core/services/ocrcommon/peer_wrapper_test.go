package ocrcommon_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/networking"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"

	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func setupV2Networking(
	t *testing.T,
	port int64,
	dbName string,
	nodeKey ethkey.KeyV2,
	backend *backends.SimulatedBackend,
) (chainlink.Application, string, common.Address, ocr2key.KeyBundle, *configtest.TestGeneralConfig) {
	p2paddresses := []string{
		fmt.Sprintf("127.0.0.1:%d", port),
	}
	config, _ := heavyweight.FullTestDB(t, fmt.Sprintf("%s%d", dbName, port))
	config.Overrides.FeatureOffchainReporting = null.BoolFrom(false)
	config.Overrides.FeatureOffchainReporting2 = null.BoolFrom(true)
	config.Overrides.FeatureLogPoller = null.BoolFrom(true)
	config.Overrides.GlobalGasEstimatorMode = null.NewString("FixedPrice", true)
	config.Overrides.P2PEnabled = null.BoolFrom(true)
	config.Overrides.SetP2PV2DeltaDial(500 * time.Millisecond)
	config.Overrides.SetP2PV2DeltaReconcile(5 * time.Second)
	config.Overrides.P2PListenPort = null.NewInt(0, true)
	config.Overrides.P2PV2ListenAddresses = p2paddresses
	config.Overrides.P2PV2AnnounceAddresses = p2paddresses
	config.Overrides.P2PNetworkingStack = networking.NetworkingStackV2
	config.Overrides.GlobalEvmGasLimitOCRJobType = null.IntFrom(5300000)

	app := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, backend, nodeKey)

	require.NoError(t, app.GetKeyStore().Unlock(testutils.Password))
	_, err := app.GetKeyStore().P2P().Create()
	require.NoError(t, err)
	p2pIDs, err := app.GetKeyStore().P2P().GetAll()
	require.NoError(t, err)
	require.Len(t, p2pIDs, 1)
	peerID := p2pIDs[0].PeerID()
	config.Overrides.P2PPeerID = peerID

	kb, err := app.GetKeyStore().OCR2().Create(chaintype.EVM)
	require.NoError(t, err)

	err = app.Start(testutils.Context(t))
	require.NoError(t, err)

	t.Cleanup(func() {
		err := app.Stop()
		require.NoError(t, err)
	})
	return app, peerID.Raw(), nodeKey.Address, kb, config
}

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

func Test_SingletonPeerWrapper_Close(t *testing.T) {
	t.Parallel()

	cfg := configtest.NewTestGeneralConfigWithOverrides(t, configtest.GeneralConfigOverrides{
		P2PEnabled: null.BoolFrom(true),
	})
	db := pgtest.NewSqlxDB(t)

	require.NoError(t, utils.JustError(db.Exec(`DELETE FROM encrypted_key_rings`)))

	keyStore := cltest.NewKeyStore(t, db, cfg)
	k, err := keyStore.P2P().Create()
	require.NoError(t, err)

	p2paddresses := []string{
		"127.0.0.1:17193",
	}

	cfg.Overrides.P2PPeerID = k.PeerID()
	cfg.Overrides.P2PNetworkingStack = networking.NetworkingStackV2
	cfg.Overrides.P2PEnabled = null.BoolFrom(true)
	cfg.Overrides.SetP2PV2DeltaDial(100 * time.Millisecond)
	cfg.Overrides.SetP2PV2DeltaReconcile(1 * time.Second)
	cfg.Overrides.P2PListenPort = null.NewInt(0, true)
	cfg.Overrides.P2PV2ListenAddresses = p2paddresses
	cfg.Overrides.P2PV2AnnounceAddresses = p2paddresses

	pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg, db, logger.TestLogger(t))

	require.NoError(t, pw.Start(testutils.Context(t)))
	require.True(t, pw.IsStarted(), "Should have started successfully")
	pw.Close()

	/* If peer is still stuck in listenLoop, we will get a bind error trying to start on the same port */
	require.False(t, pw.IsStarted())
	pw = ocrcommon.NewSingletonPeerWrapper(keyStore, cfg, db, logger.TestLogger(t))
	require.NoError(t, pw.Start(testutils.Context(t)), "Should have shut down gracefully, and be able to re-use same port")
	require.True(t, pw.IsStarted(), "Should have started successfully")
	pw.Close()
}
