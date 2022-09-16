package logpoller_test

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/log_emitter"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	EmitterABI, _ = abi.JSON(strings.NewReader(log_emitter.LogEmitterABI))
)

func logRuntime(t *testing.T, start time.Time) {
	t.Log("runtime", time.Since(start))
}

func TestPopulateLoadedDB(t *testing.T) {
	t.Skip("only for local load testing and query analysis")
	lggr := logger.TestLogger(t)
	_, db := heavyweight.FullTestDB(t, "logs_scale")
	chainID := big.NewInt(137)
	_, err := db.Exec(`INSERT INTO evm_chains (id, created_at, updated_at) VALUES ($1, NOW(), NOW())`, utils.NewBig(chainID))
	require.NoError(t, err)
	o := logpoller.NewORM(big.NewInt(137), db, lggr, pgtest.NewPGCfg(true))
	event1 := EmitterABI.Events["Log1"].ID
	address1 := common.HexToAddress("0x2ab9a2Dc53736b361b72d900CdF9F78F9406fbbb")
	address2 := common.HexToAddress("0x6E225058950f237371261C985Db6bDe26df2200E")

	// We start at 1 just so block number > 0
	for j := 1; j < 1000; j++ {
		var logs []logpoller.Log
		// Max we can insert per batch
		for i := 0; i < 1000; i++ {
			addr := address1
			if (i+(1000*j))%2 == 0 {
				addr = address2
			}
			logs = append(logs, logpoller.Log{
				EvmChainId:  utils.NewBig(chainID),
				LogIndex:    1,
				BlockHash:   common.HexToHash(fmt.Sprintf("0x%d", i+(1000*j))),
				BlockNumber: int64(i + (1000 * j)),
				EventSig:    event1[:],
				Topics:      [][]byte{event1[:], logpoller.EvmWord(uint64(i + 1000*j)).Bytes()},
				Address:     addr,
				TxHash:      common.HexToHash("0x1234"),
				Data:        logpoller.EvmWord(uint64(i + 1000*j)).Bytes(),
			})
		}
		require.NoError(t, o.InsertLogs(logs))
	}
	func() {
		defer logRuntime(t, time.Now())
		_, err := o.SelectLogsByBlockRangeFilter(750000, 800000, address1, event1[:])
		require.NoError(t, err)
	}()
	func() {
		defer logRuntime(t, time.Now())
		_, err = o.SelectLatestLogEventSigsAddrsWithConfs(0, []common.Address{address1}, []common.Hash{event1}, 0)
		require.NoError(t, err)
	}()

	// Confirm all the logs.
	require.NoError(t, o.InsertBlock(common.HexToHash("0x10"), 1000000))
	func() {
		defer logRuntime(t, time.Now())
		lgs, err := o.SelectDataWordRange(address1, event1[:], 0, logpoller.EvmWord(500000), logpoller.EvmWord(500020), 0)
		require.NoError(t, err)
		// 10 since every other log is for address1
		assert.Equal(t, 10, len(lgs))
	}()

	func() {
		defer logRuntime(t, time.Now())
		lgs, err := o.SelectIndexedLogs(address2, event1[:], 1, []common.Hash{logpoller.EvmWord(500000), logpoller.EvmWord(500020)}, 0)
		require.NoError(t, err)
		assert.Equal(t, 2, len(lgs))
	}()

	func() {
		defer logRuntime(t, time.Now())
		lgs, err := o.SelectIndexLogsTopicRange(address1, event1[:], 1, logpoller.EvmWord(500000), logpoller.EvmWord(500020), 0)
		require.NoError(t, err)
		assert.Equal(t, 10, len(lgs))
	}()
}

func TestLogPoller_Integration(t *testing.T) {
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	chainID := testutils.NewRandomEVMChainID()
	_, err := db.Exec(`INSERT INTO evm_chains (id, created_at, updated_at) VALUES ($1, NOW(), NOW())`, utils.NewBig(chainID))
	require.NoError(t, err)

	// Set up a test chain with a log emitting contract deployed.
	owner := testutils.MustNewSimTransactor(t)
	ec := backends.NewSimulatedBackend(map[common.Address]core.GenesisAccount{
		owner.From: {
			Balance: big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)),
		},
	}, 10e6)
	t.Cleanup(func() { ec.Close() })
	emitterAddress1, _, emitter1, err := log_emitter.DeployLogEmitter(owner, ec)
	require.NoError(t, err)
	ec.Commit()
	ec.Commit() // Block 2. Ensure we have finality number of blocks

	// Set up a log poller listening for log emitter logs.
	lp := logpoller.NewLogPoller(logpoller.NewORM(chainID, db, lggr, pgtest.NewPGCfg(true)),
		client.NewSimulatedBackendClient(t, ec, chainID), lggr, 100*time.Millisecond, 2, 3, 2)
	// Only filter for log1 events.
	require.NoError(t, lp.MergeFilter([]common.Hash{EmitterABI.Events["Log1"].ID}, []common.Address{emitterAddress1}))
	require.NoError(t, lp.Start(testutils.Context(t)))

	// Emit some logs in blocks 3->7.
	for i := 0; i < 5; i++ {
		emitter1.EmitLog1(owner, []*big.Int{big.NewInt(int64(i))})
		emitter1.EmitLog2(owner, []*big.Int{big.NewInt(int64(i))})
		ec.Commit()
	}
	// The poller starts on a new chain at latest-finality (5 in this case),
	// replay to ensure we get all the logs.
	require.NoError(t, lp.Replay(testutils.Context(t), 1))

	// We should immediately have all those Log1 logs.
	logs, err := lp.Logs(2, 7, EmitterABI.Events["Log1"].ID, emitterAddress1)
	require.NoError(t, err)
	assert.Equal(t, 5, len(logs))
	// Now let's update the filter and replay to get Log2 logs.
	require.NoError(t, lp.MergeFilter([]common.Hash{EmitterABI.Events["Log2"].ID}, []common.Address{emitterAddress1}))
	// Replay an invalid block should error
	assert.Error(t, lp.Replay(testutils.Context(t), 0))
	assert.Error(t, lp.Replay(testutils.Context(t), 20))
	// Replay only from block 4, so we should see logs in block 4,5,6,7 (4 logs)
	require.NoError(t, lp.Replay(testutils.Context(t), 4))

	// We should immediately see 4 logs2 logs.
	logs, err = lp.Logs(2, 7, EmitterABI.Events["Log2"].ID, emitterAddress1)
	require.NoError(t, err)
	assert.Equal(t, 4, len(logs))

	// Cancelling a replay should return an error synchronously.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	assert.True(t, errors.Is(lp.Replay(ctx, 4), logpoller.ErrReplayAbortedByClient))
	require.NoError(t, lp.Close())
}
