package ocr3_test

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/no_op_ocr3"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ocr3"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/require"
)

func setupLogPoller[RI any](t *testing.T, db *sqlx.DB) (logpoller.LogPoller, testUniverse[RI]) {
	lggr := logger.TestLogger(t)

	o := logpoller.NewORM(testutils.SimulatedChainID, db, lggr, pgtest.NewQConfig(false))

	// create the universe which will deploy the OCR contract and set config
	// we will replay on the log poller to get the appropriate ConfigSet log
	uni := newTestUniverse[RI](t)

	lp := logpoller.NewLogPoller(o, uni.simClient, lggr, 1*time.Second, false, 100, 100, 100, 200)
	return lp, uni
}

func TestConfigSet(t *testing.T) {
	require.Equal(t, no_op_ocr3.NoOpOCR3ConfigSet{}.Topic().Hex(), ocr3.ConfigSet.Hex())
}

func TestMultichainConfigTracker_New(t *testing.T) {
	t.Run("master chain not in log pollers", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		_, uni := setupLogPoller[struct{}](t, db)

		masterChain := relay.ID{
			Network: relay.EVM,
			ChainID: testutils.SimulatedChainID.String(),
		}
		_, err := ocr3.NewMultichainConfigTracker(
			masterChain,
			logger.TestLogger(t),
			map[relay.ID]logpoller.LogPoller{},
			map[relay.ID]evmclient.Client{masterChain: uni.simClient},
			map[relay.ID]common.Address{masterChain: uni.wrapper.Address()},
			ocr3.TransmitterCombiner,
		)
		require.Error(t, err, "expected error creating multichain config tracker")
	})

	t.Run("master chain not in clients", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		lp, uni := setupLogPoller[struct{}](t, db)

		masterChain := relay.ID{
			Network: relay.EVM,
			ChainID: testutils.SimulatedChainID.String(),
		}
		_, err := ocr3.NewMultichainConfigTracker(
			masterChain,
			logger.TestLogger(t),
			map[relay.ID]logpoller.LogPoller{masterChain: lp},
			map[relay.ID]evmclient.Client{},
			map[relay.ID]common.Address{masterChain: uni.wrapper.Address()},
			ocr3.TransmitterCombiner,
		)
		require.Error(t, err, "expected error creating multichain config tracker")
	})

	t.Run("master chain not in contract addresses", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		lp, uni := setupLogPoller[struct{}](t, db)

		masterChain := relay.ID{
			Network: relay.EVM,
			ChainID: testutils.SimulatedChainID.String(),
		}
		_, err := ocr3.NewMultichainConfigTracker(
			masterChain,
			logger.TestLogger(t),
			map[relay.ID]logpoller.LogPoller{masterChain: lp},
			map[relay.ID]evmclient.Client{masterChain: uni.simClient},
			map[relay.ID]common.Address{},
			ocr3.TransmitterCombiner,
		)
		require.Error(t, err, "expected error creating multichain config tracker")
	})

	t.Run("combiner is nil", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		lp, uni := setupLogPoller[struct{}](t, db)

		masterChain := relay.ID{
			Network: relay.EVM,
			ChainID: testutils.SimulatedChainID.String(),
		}
		_, err := ocr3.NewMultichainConfigTracker(
			masterChain,
			logger.TestLogger(t),
			map[relay.ID]logpoller.LogPoller{masterChain: lp},
			map[relay.ID]evmclient.Client{masterChain: uni.simClient},
			map[relay.ID]common.Address{masterChain: uni.wrapper.Address()},
			nil,
		)
		require.Error(t, err, "expected error creating multichain config tracker")
	})
}

func TestMultichainConfigTracker_SingleChain(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	lp, uni := setupLogPoller[struct{}](t, db)
	lp.Start(testutils.Context(t))

	masterChain := relay.ID{
		Network: relay.EVM,
		ChainID: testutils.SimulatedChainID.String(),
	}
	tracker, err := ocr3.NewMultichainConfigTracker(
		masterChain,
		logger.TestLogger(t),
		map[relay.ID]logpoller.LogPoller{masterChain: lp},
		map[relay.ID]evmclient.Client{masterChain: uni.simClient},
		map[relay.ID]common.Address{masterChain: uni.wrapper.Address()},
		ocr3.TransmitterCombiner,
	)
	require.NoError(t, err, "failed to create multichain config tracker")

	// Replay the log poller to get the ConfigSet log
	err = tracker.Replay(testutils.Context(t), masterChain, 1)
	require.NoError(t, err, "failed to replay log poller")

	// fetch config digest from the tracker
	changedInBlock, configDigest, err := tracker.LatestConfigDetails(testutils.Context(t))
	require.NoError(t, err, "failed to get latest config details")
	c, err := uni.wrapper.LatestConfigDigestAndEpoch(nil)
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
		for _, b := range uni.bundles {
			signers = append(signers, b.PublicKey())
		}
		return signers
	}()
	expectedTransmitters := func() []ocrtypes.Account {
		var accounts []ocrtypes.Account
		for _, tm := range uni.transmitters {
			accounts = append(accounts, ocrtypes.Account(tm.From.Hex()))
		}
		return accounts
	}()
	require.Equal(t, expectedSigners, config.Signers, "expected signers to match")
	require.Equal(t, expectedTransmitters, config.Transmitters, "expected transmitters to match")
}
