package ocr3impls_test

import (
	"encoding/hex"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	chainsel "github.com/smartcontractkit/chain-selectors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/no_op_ocr3"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	discoverermocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/discoverer/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/graph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/ocr3impls"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
)

func setupLogPoller[RI ocr3impls.MultichainMeta](t *testing.T, db *sqlx.DB, bs *keyringsAndSigners[RI]) (logpoller.LogPoller, testUniverse[RI]) {
	lggr := logger.TestLogger(t)

	o := logpoller.NewORM(testutils.SimulatedChainID, db, lggr)

	// create the universe which will deploy the OCR contract and set config
	// we will replay on the log poller to get the appropriate ConfigSet log
	uni := newTestUniverse[RI](t, bs)
	lpOpts := logpoller.Opts{
		PollPeriod:               1 * time.Second,
		FinalityDepth:            100,
		BackfillBatchSize:        100,
		RpcBatchSize:             100,
		KeepFinalizedBlocksDepth: 200,
	}
	headTracker := headtracker.NewSimulatedHeadTracker(uni.simClient, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
	if lpOpts.PollPeriod == 0 {
		lpOpts.PollPeriod = 1 * time.Hour
	}
	lp := logpoller.NewLogPoller(o, uni.simClient, lggr, headTracker, lpOpts)
	return lp, uni
}

func TestConfigSet(t *testing.T) {
	require.Equal(t, no_op_ocr3.NoOpOCR3ConfigSet{}.Topic().Hex(), ocr3impls.ConfigSet.Hex())
}

func TestMultichainConfigTracker_New(t *testing.T) {
	t.Run("master chain not in log pollers", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		_, uni := setupLogPoller[multichainMeta](t, db, nil)

		masterChain := commontypes.RelayID{
			Network: relay.NetworkEVM,
			ChainID: testutils.SimulatedChainID.String(),
		}
		mockDiscovererFactory := discoverermocks.NewFactory(t)
		_, err := ocr3impls.NewMultichainConfigTracker(
			masterChain,
			logger.TestLogger(t),
			map[commontypes.RelayID]logpoller.LogPoller{},
			uni.simClient,
			uni.wrapper.Address(),
			mockDiscovererFactory,
			ocr3impls.TransmitterCombiner,
			nil,
		)
		require.Error(t, err, "expected error creating multichain config tracker")
	})

	t.Run("combiner is nil", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		lp, uni := setupLogPoller[multichainMeta](t, db, nil)

		masterChain := commontypes.RelayID{
			Network: relay.NetworkEVM,
			ChainID: testutils.SimulatedChainID.String(),
		}
		mockDiscovererFactory := discoverermocks.NewFactory(t)
		_, err := ocr3impls.NewMultichainConfigTracker(
			masterChain,
			logger.TestLogger(t),
			map[commontypes.RelayID]logpoller.LogPoller{masterChain: lp},
			uni.simClient,
			uni.wrapper.Address(),
			mockDiscovererFactory,
			nil,
			nil,
		)
		require.Error(t, err, "expected error creating multichain config tracker")
	})

	t.Run("factory is nil", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		lp, uni := setupLogPoller[multichainMeta](t, db, nil)

		masterChain := commontypes.RelayID{
			Network: relay.NetworkEVM,
			ChainID: testutils.SimulatedChainID.String(),
		}
		_, err := ocr3impls.NewMultichainConfigTracker(
			masterChain,
			logger.TestLogger(t),
			map[commontypes.RelayID]logpoller.LogPoller{masterChain: lp},
			uni.simClient,
			uni.wrapper.Address(),
			nil,
			ocr3impls.TransmitterCombiner,
			nil,
		)
		require.Error(t, err, "expected error creating multichain config tracker")
	})
}

func TestMultichainConfigTracker_SingleChain(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	lp, uni := setupLogPoller[multichainMeta](t, db, nil)
	require.NoError(t, lp.Start(testutils.Context(t)))
	t.Cleanup(func() { require.NoError(t, lp.Close()) })

	masterChain := commontypes.RelayID{
		Network: relay.NetworkEVM,
		ChainID: testutils.SimulatedChainID.String(),
	}

	ch, exists := chainsel.ChainByEvmChainID(uint64(mustStrToI64(t, masterChain.ChainID)))
	assert.True(t, exists)

	// for this test only one LM is "deployed"
	// so the discovery will return a single LM which is the master LM
	g := graph.NewGraph()
	g.(graph.GraphTest).AddNetwork(models.NetworkSelector(ch.Selector), graph.Data{
		Liquidity:               big.NewInt(1234), // liquidity doesn't matter for this test
		LiquidityManagerAddress: models.Address(uni.wrapper.Address()),
	})
	mockDiscoverer := discoverermocks.NewDiscoverer(t)
	mockDiscoverer.On("Discover", mock.Anything).Return(g, nil)
	defer mockDiscoverer.AssertExpectations(t)
	mockDiscovererFactory := discoverermocks.NewFactory(t)
	mockDiscovererFactory.On("NewDiscoverer", models.NetworkSelector(ch.Selector), models.Address(uni.wrapper.Address())).
		Return(mockDiscoverer, nil)
	defer mockDiscovererFactory.AssertExpectations(t)
	tracker, err := ocr3impls.NewMultichainConfigTracker(
		masterChain,
		logger.TestLogger(t),
		map[commontypes.RelayID]logpoller.LogPoller{masterChain: lp},
		uni.simClient,
		uni.wrapper.Address(),
		mockDiscovererFactory,
		ocr3impls.TransmitterCombiner,
		nil,
	)
	require.NoError(t, err, "failed to create multichain config tracker")

	// Replay the log poller to get the ConfigSet log
	err = tracker.ReplayChain(testutils.Context(t), masterChain, 1)
	require.NoError(t, err, "failed to replay log poller")

	// fetch config digest from the tracker
	changedInBlock, configDigest, err := tracker.LatestConfigDetails(testutils.Context(t))
	require.NoError(t, err, "failed to get latest config details")
	c, err := uni.wrapper.LatestConfigDetails(nil)
	require.NoError(t, err, "failed to get latest config digest and epoch")
	require.Equal(t, hex.EncodeToString(c.ConfigDigest[:]), configDigest.Hex(), "expected latest config digest to match")

	// fetch config details from the tracker
	config, err := tracker.LatestConfig(testutils.Context(t), changedInBlock)
	require.NoError(t, err, "failed to get latest config")
	require.Equal(t, uint64(1), config.ConfigCount, "expected config count to match")
	require.Equal(t, configDigest, config.ConfigDigest, "expected config digest to match")
	require.Equal(t, uint8(1), config.F, "expected f to match")
	require.Equal(t, []byte{}, config.OnchainConfig, "expected onchain config to match")
	require.Equal(t, []byte{}, config.OffchainConfig, "expected offchain config to match")
	require.Equal(t, uint64(3), config.OffchainConfigVersion, "expected offchain config version to match")
	expectedSigners := func() []ocrtypes.OnchainPublicKey {
		var signers []ocrtypes.OnchainPublicKey
		for _, b := range uni.keyrings {
			signers = append(signers, b.PublicKey())
		}
		return signers
	}()
	expectedTransmitters := func() []ocrtypes.Account {
		var accounts []ocrtypes.Account
		for _, tm := range uni.transmitters {
			accounts = append(accounts, ocrtypes.Account(ocr3impls.EncodeTransmitter(masterChain, ocrtypes.Account(tm.From.Hex()))))
		}
		return accounts
	}()
	require.Equal(t, expectedSigners, config.Signers, "expected signers to match")
	require.Equal(t, expectedTransmitters, config.Transmitters, "expected transmitters to match")
}

func TestMultichainConfigTracker_Multichain(t *testing.T) {
	// create heavyweight db's because the log pollers need to have separate
	// databases to avoid conflicts.
	_, db1 := heavyweight.FullTestDBV2(t, nil)
	_, db2 := heavyweight.FullTestDBV2(t, nil)

	lp1, uni1 := setupLogPoller[multichainMeta](t, db1, nil)
	lp2, uni2 := setupLogPoller[multichainMeta](t, db2, &keyringsAndSigners[multichainMeta]{
		keyrings: uni1.keyrings,
		signers:  uni1.signers,
	})
	t.Cleanup(func() {
		require.NoError(t, lp1.Close())
		require.NoError(t, lp2.Close())
	})

	// finality depth
	uni2.backend.Commit()
	uni2.backend.Commit()

	// start the log pollers
	require.NoError(t, lp1.Start(testutils.Context(t)))
	require.NoError(t, lp2.Start(testutils.Context(t)))

	// create the multichain config tracker
	// the chain id's we're using in the mappings are different from the
	// simulated chain id but that should be fine for this test.
	masterChain := commontypes.RelayID{
		Network: relay.NetworkEVM,
		ChainID: strconv.FormatUint(chainsel.TEST_90000001.EvmChainID, 10),
	}
	secondChain := commontypes.RelayID{
		Network: relay.NetworkEVM,
		ChainID: strconv.FormatUint(chainsel.TEST_90000002.EvmChainID, 10),
	}

	chain1, exists := chainsel.ChainByEvmChainID(uint64(mustStrToI64(t, masterChain.ChainID)))
	assert.True(t, exists)

	chain2, exists := chainsel.ChainByEvmChainID(uint64(mustStrToI64(t, secondChain.ChainID)))
	assert.True(t, exists)

	// this test doesn't care about the connections, just the vertices themselves
	g := graph.NewGraph()
	g.(graph.GraphTest).AddNetwork(models.NetworkSelector(chain1.Selector), graph.Data{
		Liquidity:               big.NewInt(1234), // liquidity doesn't matter for this test
		LiquidityManagerAddress: models.Address(uni1.wrapper.Address()),
	})
	g.(graph.GraphTest).AddNetwork(models.NetworkSelector(chain2.Selector), graph.Data{
		Liquidity:               big.NewInt(1234), // liquidity doesn't matter for this test
		LiquidityManagerAddress: models.Address(uni2.wrapper.Address()),
	})
	mockDiscoverer := discoverermocks.NewDiscoverer(t)
	mockDiscoverer.On("Discover", mock.Anything).Return(g, nil)
	defer mockDiscoverer.AssertExpectations(t)
	mockDiscovererFactory := discoverermocks.NewFactory(t)
	mockDiscovererFactory.On("NewDiscoverer", models.NetworkSelector(chain1.Selector), models.Address(uni1.wrapper.Address())).
		Return(mockDiscoverer, nil)
	defer mockDiscovererFactory.AssertExpectations(t)
	tracker, err := ocr3impls.NewMultichainConfigTracker(
		masterChain,
		logger.TestLogger(t),
		map[commontypes.RelayID]logpoller.LogPoller{
			masterChain: lp1,
			secondChain: lp2,
		},
		uni1.simClient,
		uni1.wrapper.Address(),
		mockDiscovererFactory,
		ocr3impls.TransmitterCombiner,
		nil, // we call replay explicitly below
	)
	require.NoError(t, err, "failed to create multichain config tracker")

	// Replay the log pollers to get the ConfigSet log
	// on each respective chain
	require.NoError(t, tracker.ReplayChain(testutils.Context(t), masterChain, 1), "failed to replay log poller on master chain")
	require.NoError(t, tracker.ReplayChain(testutils.Context(t), secondChain, 1), "failed to replay log poller on second chain")

	// fetch config digest from the tracker
	changedInBlock, configDigest, err := tracker.LatestConfigDetails(testutils.Context(t))
	require.NoError(t, err, "failed to get latest config details")
	c, err := uni1.wrapper.LatestConfigDetails(nil)
	require.NoError(t, err, "failed to get latest config digest and epoch")
	require.Equal(t, hex.EncodeToString(c.ConfigDigest[:]), configDigest.Hex(), "expected latest config digest to match")

	// fetch config details from the tracker
	config, err := tracker.LatestConfig(testutils.Context(t), changedInBlock)
	require.NoError(t, err, "failed to get latest config")
	require.Equal(t, uint64(1), config.ConfigCount, "expected config count to match")
	require.Equal(t, configDigest, config.ConfigDigest, "expected config digest to match")
	require.Equal(t, uint8(1), config.F, "expected f to match")
	require.Equal(t, []byte{}, config.OnchainConfig, "expected onchain config to match")
	require.Equal(t, []byte{}, config.OffchainConfig, "expected offchain config to match")
	require.Equal(t, uint64(3), config.OffchainConfigVersion, "expected offchain config version to match")
	expectedSigners := func() []ocrtypes.OnchainPublicKey {
		var signers []ocrtypes.OnchainPublicKey
		for _, b := range uni1.keyrings {
			signers = append(signers, b.PublicKey())
		}
		return signers
	}()
	require.Equal(t, expectedSigners, config.Signers, "expected signers to match")
	expectedTransmitters := func() []ocrtypes.Account {
		var accounts []ocrtypes.Account
		for i := range uni1.transmitters {
			t1 := ocr3impls.EncodeTransmitter(masterChain, ocrtypes.Account(uni1.transmitters[i].From.Hex()))
			t2 := ocr3impls.EncodeTransmitter(secondChain, ocrtypes.Account(uni2.transmitters[i].From.Hex()))
			accounts = append(accounts, ocrtypes.Account(ocr3impls.JoinTransmitters([]string{t1, t2})))
		}
		return accounts
	}()
	require.Equal(t, expectedTransmitters, config.Transmitters, "expected transmitters to match")
}

func mustStrToI64(t *testing.T, s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	require.NoError(t, err)
	return i
}
