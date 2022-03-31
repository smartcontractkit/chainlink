package logpoller_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/log_emitter"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestPopulateLoadedDB(t *testing.T) {
	t.Skip()
	lggr := logger.TestLogger(t)
	_, db := heavyweight.FullTestDB(t, "logs_scale", true, false)
	chainID := big.NewInt(137)
	_, err := db.Exec(`INSERT INTO evm_chains (id, created_at, updated_at) VALUES ($1, NOW(), NOW())`, utils.NewBig(chainID))
	require.NoError(t, err)
	o := logpoller.NewORM(big.NewInt(137), db, lggr, pgtest.NewPGCfg(true))
	event1 := logpoller.EmitterABI.Events["Log1"].ID
	address1 := common.HexToAddress("0x2ab9a2Dc53736b361b72d900CdF9F78F9406fbbb")
	address2 := common.HexToAddress("0x6E225058950f237371261C985Db6bDe26df2200E")

	for j := 0; j < 1000; j++ {
		var logs []logpoller.Log
		// Max we can insert per batch
		for i := 0; i < 1000; i++ {
			addr := address1
			if (i+(1000*j))%2 == 0 {
				addr = address2
			}
			logs = append(logs, logpoller.GenLog(chainID, 1, int64(i+(1000*j)), fmt.Sprintf("0x%d", i+(1000*j)), event1[:], addr))
		}
		require.NoError(t, o.InsertLogs(logs))
	}
	s := time.Now()
	lgs, err := o.SelectLogsByBlockRangeFilter(750000, 800000, address1, event1[:])
	require.NoError(t, err)
	t.Log(time.Since(s), len(lgs))
}

func TestLogPollerIntegration(t *testing.T) {
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	//_, db := heavyweight.FullTestDB(t, "log", true, false)
	chainID := big.NewInt(42)
	_, err := db.Exec(`INSERT INTO evm_chains (id, created_at, updated_at) VALUES ($1, NOW(), NOW())`, utils.NewBig(chainID))
	require.NoError(t, err)

	// Set up a test chain with a log emitting contract deployed.
	owner := testutils.MustNewSimTransactor()
	ec := backends.NewSimulatedBackend(map[common.Address]core.GenesisAccount{
		owner.From: {
			Balance: big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)),
		},
	}, 10e6)
	t.Cleanup(func() { ec.Close() })
	emitterAddress1, _, emitter1, err := log_emitter.DeployLogEmitter(owner, ec)
	require.NoError(t, err)
	ec.Commit() // Block 2.

	// Set up a log poller listening for log emitter logs.
	lp := logpoller.NewLogPoller(logpoller.NewORM(chainID, db, lggr, pgtest.NewPGCfg(true)),
		client.NewSimulatedBackendClient(t, ec, chainID), lggr, 100*time.Millisecond, 2, 3)
	// Only filter for log1 events.
	lp.MergeFilter([]common.Hash{logpoller.EmitterABI.Events["Log1"].ID}, emitterAddress1)
	require.NoError(t, lp.Start(context.Background()))

	// Emit some logs in blocks 2->6.
	for i := 0; i < 5; i++ {
		emitter1.EmitLog1(owner, []*big.Int{big.NewInt(int64(i))})
		emitter1.EmitLog2(owner, []*big.Int{big.NewInt(int64(i))})
		ec.Commit()
	}
	// We should eventually receive all those Log1 logs.
	testutils.AssertEventually(t, func() bool {
		logs, err := lp.Logs(1, 6, logpoller.EmitterABI.Events["Log1"].ID, emitterAddress1)
		require.NoError(t, err)
		t.Logf("Received %d/%d logs\n", len(logs), 5)
		return len(logs) == 5
	})
	// Now let's update the filter and replay to get Log2 logs.
	lp.MergeFilter([]common.Hash{logpoller.EmitterABI.Events["Log2"].ID}, emitterAddress1)
	// Replay only from block 3, so we should see logs in block 3,4,5,6 (4 logs)
	lp.Replay(3)

	// We should eventually see 4 logs2 logs.
	testutils.AssertEventually(t, func() bool {
		logs, err := lp.Logs(1, 6, logpoller.EmitterABI.Events["Log2"].ID, emitterAddress1)
		require.NoError(t, err)
		t.Logf("Received %d/%d logs\n", len(logs), 4)
		return len(logs) == 4
	})

	require.NoError(t, lp.Close())
}
