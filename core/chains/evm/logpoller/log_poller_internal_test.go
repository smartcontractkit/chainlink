package logpoller

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
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgconn"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_emitter"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	EmitterABI, _ = abi.JSON(strings.NewReader(log_emitter.LogEmitterABI))
)

// Validate that filters stored in log_filters_table match the filters stored in memory
func validateFiltersTable(t *testing.T, lp *logPoller, orm *ORM) {
	filters, err := orm.LoadFilters()
	require.NoError(t, err)
	require.Equal(t, len(filters), len(lp.filters))
	for name, dbFilter := range filters {
		dbFilter := dbFilter
		memFilter, ok := lp.filters[name]
		require.True(t, ok)
		assert.True(t, memFilter.Contains(&dbFilter),
			fmt.Sprintf("in-memory Filter %s is missing some addresses or events from db Filter table", name))
		assert.True(t, dbFilter.Contains(&memFilter),
			fmt.Sprintf("db Filter table %s is missing some addresses or events from in-memory Filter", name))
	}
}

func TestLogPoller_RegisterFilter(t *testing.T) {
	t.Parallel()
	a1 := common.HexToAddress("0x2ab9a2dc53736b361b72d900cdf9f78f9406fbbb")
	a2 := common.HexToAddress("0x2ab9a2dc53736b361b72d900cdf9f78f9406fbbc")

	lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.ErrorLevel)
	chainID := testutils.NewRandomEVMChainID()
	db := pgtest.NewSqlxDB(t)

	orm := NewORM(chainID, db, lggr, pgtest.NewQConfig(true))
	lp := NewLogPoller(orm, nil, lggr, 15*time.Second, 1, 1, 2, 1000)

	filter := Filter{"test Filter", []common.Hash{EmitterABI.Events["Log1"].ID}, []common.Address{a1}}
	err := lp.RegisterFilter(filter)
	require.Error(t, err, "RegisterFilter failed to save Filter to db")
	require.Equal(t, 1, observedLogs.Len())
	assertForeignConstraintError(t, observedLogs.TakeAll()[0], "evm_log_poller_filters", "evm_log_poller_filters_evm_chain_id_fkey")

	db.Close()
	db = pgtest.NewSqlxDB(t)
	orm = NewORM(chainID, db, lggr, pgtest.NewQConfig(true))

	// disable check that chain id exists for rest of tests
	require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS evm_log_poller_filters_evm_chain_id_fkey DEFERRED`)))
	// Set up a test chain with a log emitting contract deployed.

	lp = NewLogPoller(orm, nil, lggr, 15*time.Second, 1, 1, 2, 1000)

	// We expect a zero Filter if nothing registered yet.
	f := lp.Filter(nil, nil, nil)
	require.Equal(t, 1, len(f.Addresses))
	assert.Equal(t, common.HexToAddress("0x0000000000000000000000000000000000000000"), f.Addresses[0])

	err = lp.RegisterFilter(Filter{"Emitter Log 1", []common.Hash{EmitterABI.Events["Log1"].ID}, []common.Address{a1}})
	require.NoError(t, err)
	assert.Equal(t, []common.Address{a1}, lp.Filter(nil, nil, nil).Addresses)
	assert.Equal(t, [][]common.Hash{{EmitterABI.Events["Log1"].ID}}, lp.Filter(nil, nil, nil).Topics)
	validateFiltersTable(t, lp, orm)

	// Should de-dupe EventSigs
	err = lp.RegisterFilter(Filter{"Emitter Log 1 + 2", []common.Hash{EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}, []common.Address{a2}})
	require.NoError(t, err)
	assert.Equal(t, []common.Address{a1, a2}, lp.Filter(nil, nil, nil).Addresses)
	assert.Equal(t, [][]common.Hash{{EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}}, lp.Filter(nil, nil, nil).Topics)
	validateFiltersTable(t, lp, orm)

	// Should de-dupe Addresses
	err = lp.RegisterFilter(Filter{"Emitter Log 1 + 2 dupe", []common.Hash{EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}, []common.Address{a2}})
	require.NoError(t, err)
	assert.Equal(t, []common.Address{a1, a2}, lp.Filter(nil, nil, nil).Addresses)
	assert.Equal(t, [][]common.Hash{{EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}}, lp.Filter(nil, nil, nil).Topics)
	validateFiltersTable(t, lp, orm)

	// Address required.
	err = lp.RegisterFilter(Filter{"no address", []common.Hash{EmitterABI.Events["Log1"].ID}, []common.Address{}})
	require.Error(t, err)
	// Event required
	err = lp.RegisterFilter(Filter{"No event", []common.Hash{}, []common.Address{a1}})
	require.Error(t, err)
	validateFiltersTable(t, lp, orm)

	// Removing non-existence Filter should log error but return nil
	err = lp.UnregisterFilter("Filter doesn't exist", nil)
	require.NoError(t, err)
	require.Equal(t, observedLogs.Len(), 1)
	require.Contains(t, observedLogs.TakeAll()[0].Entry.Message, "not found")

	// Check that all filters are still there
	_, ok := lp.filters["Emitter Log 1"]
	require.True(t, ok, "'Emitter Log 1 Filter' missing")
	_, ok = lp.filters["Emitter Log 1 + 2"]
	require.True(t, ok, "'Emitter Log 1 + 2' Filter missing")
	_, ok = lp.filters["Emitter Log 1 + 2 dupe"]
	require.True(t, ok, "'Emitter Log 1 + 2 dupe' Filter missing")

	// Removing an existing Filter should remove it from both memory and db
	err = lp.UnregisterFilter("Emitter Log 1 + 2", nil)
	require.NoError(t, err)
	_, ok = lp.filters["Emitter Log 1 + 2"]
	require.False(t, ok, "'Emitter Log 1 Filter' should have been removed by UnregisterFilter()")
	require.Len(t, lp.filters, 2)
	validateFiltersTable(t, lp, orm)

	err = lp.UnregisterFilter("Emitter Log 1 + 2 dupe", nil)
	require.NoError(t, err)
	err = lp.UnregisterFilter("Emitter Log 1", nil)
	require.NoError(t, err)
	assert.Len(t, lp.filters, 0)
	filters, err := lp.orm.LoadFilters()
	require.NoError(t, err)
	assert.Len(t, filters, 0)

	// Make sure cache was invalidated
	assert.Len(t, lp.Filter(nil, nil, nil).Addresses, 1)
	assert.Equal(t, lp.Filter(nil, nil, nil).Addresses[0], common.HexToAddress("0x0000000000000000000000000000000000000000"))
	assert.Len(t, lp.Filter(nil, nil, nil).Topics, 1)
	assert.Len(t, lp.Filter(nil, nil, nil).Topics[0], 0)
}

func assertForeignConstraintError(t *testing.T, observedLog observer.LoggedEntry,
	table string, constraint string) {

	assert.Equal(t, "SQL ERROR", observedLog.Entry.Message)

	field := observedLog.Context[0]
	require.Equal(t, zapcore.ErrorType, field.Type)
	err, ok := field.Interface.(error)
	var pgErr *pgconn.PgError
	require.True(t, errors.As(err, &pgErr))
	require.True(t, ok)
	assert.Equal(t, "23503", pgErr.SQLState()) // foreign key constraint violation code
	assert.Equal(t, table, pgErr.TableName)
	assert.Equal(t, constraint, pgErr.ConstraintName)
}

func TestLogPoller_DBErrorHandling(t *testing.T) {
	t.Parallel()
	lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.WarnLevel)
	chainID1 := testutils.NewRandomEVMChainID()
	chainID2 := testutils.NewRandomEVMChainID()
	db := pgtest.NewSqlxDB(t)
	o := NewORM(chainID1, db, lggr, pgtest.NewQConfig(true))

	owner := testutils.MustNewSimTransactor(t)
	ethDB := rawdb.NewMemoryDatabase()
	ec := backends.NewSimulatedBackendWithDatabase(ethDB, map[common.Address]core.GenesisAccount{
		owner.From: {
			Balance: big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)),
		},
	}, 10e6)
	_, _, emitter, err := log_emitter.DeployLogEmitter(owner, ec)
	require.NoError(t, err)
	_, err = emitter.EmitLog1(owner, []*big.Int{big.NewInt(9)})
	require.NoError(t, err)
	_, err = emitter.EmitLog1(owner, []*big.Int{big.NewInt(7)})
	require.NoError(t, err)
	ec.Commit()
	ec.Commit()
	ec.Commit()

	lp := NewLogPoller(o, client.NewSimulatedBackendClient(t, ec, chainID2), lggr, 1*time.Hour, 2, 3, 2, 1000)
	ctx, cancelReplay := context.WithCancel(testutils.Context(t))
	lp.ctx, lp.cancel = context.WithCancel(testutils.Context(t))
	defer cancelReplay()
	defer lp.cancel()

	err = lp.Replay(ctx, 5) // block number too high
	require.ErrorContains(t, err, "Invalid replay block number")

	// Force a db error while loading the filters (tx aborted, already rolled back)
	require.Error(t, utils.JustError(db.Exec(`invalid query`)))
	go func() {
		err = lp.Replay(ctx, 2)
		assert.Error(t, err, ErrReplayAbortedByClient)
	}()

	time.Sleep(100 * time.Millisecond)
	go lp.run()
	require.Eventually(t, func() bool {
		return observedLogs.Len() >= 5
	}, 2*time.Second, 20*time.Millisecond)
	lp.cancel()
	lp.Close()
	<-lp.done

	logMsgs := make(map[string]int)
	for _, obs := range observedLogs.All() {
		_, ok := logMsgs[obs.Entry.Message]
		if ok {
			logMsgs[(obs.Entry.Message)] = 1
		} else {
			logMsgs[(obs.Entry.Message)]++
		}
	}

	assert.Contains(t, logMsgs, "SQL ERROR")
	assert.Contains(t, logMsgs, "Failed loading filters in main logpoller loop, retrying later")
	assert.Contains(t, logMsgs, "Error executing replay, could not get fromBlock")
	assert.Contains(t, logMsgs, "backup log poller ran before filters loaded, skipping")
}

func TestLogPoller_ConvertLogs(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)

	topics := []common.Hash{EmitterABI.Events["Log1"].ID}

	var cases = []struct {
		name     string
		logs     []types.Log
		blocks   []LogPollerBlock
		expected int
	}{
		{"SingleBlock",
			[]types.Log{{Topics: topics}, {Topics: topics}},
			[]LogPollerBlock{{BlockTimestamp: time.Now()}},
			2},
		{"BlockList",
			[]types.Log{{Topics: topics}, {Topics: topics}, {Topics: topics}},
			[]LogPollerBlock{{BlockTimestamp: time.Now()}},
			3},
		{"EmptyList",
			[]types.Log{},
			[]LogPollerBlock{},
			0},
		{"TooManyBlocks",
			[]types.Log{{}},
			[]LogPollerBlock{{}, {}},
			0},
		{"TooFewBlocks",
			[]types.Log{{}, {}, {}},
			[]LogPollerBlock{{}, {}},
			0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			logs := convertLogs(c.logs, c.blocks, lggr, big.NewInt(53))
			require.Len(t, logs, c.expected)
			for i := 0; i < c.expected; i++ {
				if len(c.blocks) == 1 {
					assert.Equal(t, c.blocks[0].BlockTimestamp, logs[i].BlockTimestamp)
				} else {
					assert.Equal(t, logs[i].BlockTimestamp, c.blocks[i].BlockTimestamp)
				}
			}
		})
	}
}

func TestFilterName(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "a - b:c:d", FilterName("a", "b", "c", "d"))
	assert.Equal(t, "empty args test", FilterName("empty args test"))
}

func benchmarkFilter(b *testing.B, nFilters, nAddresses, nEvents int) {
	lggr := logger.TestLogger(b)
	lp := NewLogPoller(nil, nil, lggr, 1*time.Hour, 2, 3, 2, 1000)
	for i := 0; i < nFilters; i++ {
		var addresses []common.Address
		var events []common.Hash
		for j := 0; j < nAddresses; j++ {
			addresses = append(addresses, common.BigToAddress(big.NewInt(int64(j+1))))
		}
		for j := 0; j < nEvents; j++ {
			events = append(events, common.BigToHash(big.NewInt(int64(j+1))))
		}
		err := lp.RegisterFilter(Filter{Name: "my Filter", EventSigs: events, Addresses: addresses})
		require.NoError(b, err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		lp.Filter(nil, nil, nil)
	}
}

func BenchmarkFilter10_1(b *testing.B) {
	benchmarkFilter(b, 10, 1, 1)
}
func BenchmarkFilter100_10(b *testing.B) {
	benchmarkFilter(b, 100, 10, 10)
}
func BenchmarkFilter1000_100(b *testing.B) {
	benchmarkFilter(b, 1000, 100, 100)
}
