package upkeepstate

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

func TestUpkeepStateStore(t *testing.T) {
	tests := []struct {
		name               string
		inserts            []ocr2keepers.CheckResult
		workIDsSelect      []string
		workIDsFromScanner []string
		errScanner         error
		recordsFromDB      []PersistedStateRecord
		errDB              error
		expected           []ocr2keepers.UpkeepState
		errored            bool
	}{
		{
			name: "empty store",
		},
		{
			name: "save only ineligible states",
			inserts: []ocr2keepers.CheckResult{
				{
					UpkeepID: createUpkeepIDForTest(1),
					WorkID:   "0x1",
					Eligible: false,
					Trigger: ocr2keepers.Trigger{
						BlockNumber: ocr2keepers.BlockNumber(1),
					},
				},
				{
					UpkeepID: createUpkeepIDForTest(2),
					WorkID:   "ox2",
					Eligible: true,
					Trigger: ocr2keepers.Trigger{
						BlockNumber: ocr2keepers.BlockNumber(2),
					},
				},
			},
			workIDsSelect: []string{"0x1", "0x2"},
			expected: []ocr2keepers.UpkeepState{
				ocr2keepers.Ineligible,
				StateUnknown,
			},
		},
		{
			name: "fetch results from scanner",
			inserts: []ocr2keepers.CheckResult{
				{
					UpkeepID: createUpkeepIDForTest(1),
					WorkID:   "0x1",
					Eligible: false,
					Trigger: ocr2keepers.Trigger{
						BlockNumber: ocr2keepers.BlockNumber(1),
					},
				},
			},
			workIDsSelect:      []string{"0x1", "0x2"},
			workIDsFromScanner: []string{"0x2", "0x222"},
			expected: []ocr2keepers.UpkeepState{
				ocr2keepers.Ineligible,
				ocr2keepers.Performed,
			},
		},
		{
			name: "fetch results from db",
			inserts: []ocr2keepers.CheckResult{
				{
					UpkeepID: createUpkeepIDForTest(1),
					WorkID:   "0x1",
					Eligible: false,
					Trigger: ocr2keepers.Trigger{
						BlockNumber: ocr2keepers.BlockNumber(1),
					},
				},
			},
			workIDsSelect:      []string{"0x1", "0x2", "0x3"},
			workIDsFromScanner: []string{"0x2", "0x222"},
			recordsFromDB: []PersistedStateRecord{
				{
					WorkID:          "0x3",
					CompletionState: 1,
					BlockNumber:     2,
					AddedAt:         time.Now(),
				},
			},
			expected: []ocr2keepers.UpkeepState{
				ocr2keepers.Ineligible,
				ocr2keepers.Performed,
				ocr2keepers.Ineligible,
			},
		},
		{
			name: "unknown states",
			inserts: []ocr2keepers.CheckResult{
				{
					UpkeepID: createUpkeepIDForTest(1),
					WorkID:   "0x1",
					Eligible: false,
					Trigger: ocr2keepers.Trigger{
						BlockNumber: ocr2keepers.BlockNumber(1),
					},
				},
			},
			workIDsSelect:      []string{"0x2"},
			workIDsFromScanner: []string{},
			expected: []ocr2keepers.UpkeepState{
				StateUnknown,
			},
		},
		{
			name: "scanner error",
			inserts: []ocr2keepers.CheckResult{
				{
					UpkeepID: createUpkeepIDForTest(1),
					WorkID:   "0x1",
					Eligible: false,
					Trigger: ocr2keepers.Trigger{
						BlockNumber: ocr2keepers.BlockNumber(1),
					},
				},
			},
			workIDsSelect:      []string{"0x1", "0x2"},
			workIDsFromScanner: []string{"0x2", "0x222"},
			errScanner:         fmt.Errorf("test error"),
			errored:            true,
		},
		{
			name: "db error",
			inserts: []ocr2keepers.CheckResult{
				{
					UpkeepID: createUpkeepIDForTest(1),
					WorkID:   "0x1",
					Eligible: false,
					Trigger: ocr2keepers.Trigger{
						BlockNumber: ocr2keepers.BlockNumber(1),
					},
				},
			},
			workIDsSelect:      []string{"0x1", "0x2"},
			workIDsFromScanner: []string{"0x2", "0x222"},
			errDB:              fmt.Errorf("test error"),
			errored:            true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testutils.Context(t)
			lggr := logger.TestLogger(t)

			scanner := &mockScanner{}
			scanner.addWorkID(tc.workIDsFromScanner...)
			scanner.setErr(tc.errScanner)

			orm := &mockORM{}
			orm.addRecords(tc.recordsFromDB...)
			orm.setErr(tc.errDB)

			store := NewUpkeepStateStore(orm, lggr, scanner)

			for _, insert := range tc.inserts {
				assert.NoError(t, store.SetUpkeepState(ctx, insert, ocr2keepers.Performed))
			}

			states, err := store.SelectByWorkIDsInRange(ctx, 1, 100, tc.workIDsSelect...)
			if tc.errored {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			assert.Equal(t, len(tc.expected), len(states))
			for i, state := range states {
				assert.Equal(t, tc.expected[i], state)
			}
		})
	}
}

func TestUpkeepStateStore_SetSelectIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("database required for upkeep state store integration test")
	}

	makeTestResult := func(id int64, workID string, eligible bool, block uint64) ocr2keepers.CheckResult {
		uid := &ocr2keepers.UpkeepIdentifier{}
		_ = uid.FromBigInt(big.NewInt(id))

		return ocr2keepers.CheckResult{
			UpkeepID: *uid,
			WorkID:   workID,
			Eligible: eligible,
			Trigger: ocr2keepers.Trigger{
				BlockNumber: ocr2keepers.BlockNumber(block),
			},
		}
	}

	type storedValue struct {
		result ocr2keepers.CheckResult
		state  ocr2keepers.UpkeepState
	}

	tests := []struct {
		name         string
		queryIDs     []string
		storedValues []storedValue
		expected     []ocr2keepers.UpkeepState
	}{
		{
			name:     "querying non-stored workIDs on empty db returns unknown state results",
			queryIDs: []string{"0x1", "0x2", "0x3", "0x4"},
			expected: []ocr2keepers.UpkeepState{
				StateUnknown,
				StateUnknown,
				StateUnknown,
				StateUnknown,
			},
		},
		{
			name:     "querying non-stored workIDs on db with values returns unknown state results",
			queryIDs: []string{"0x1", "0x2", "0x3", "0x4"},
			storedValues: []storedValue{
				{result: makeTestResult(1, "0x11", false, 1), state: ocr2keepers.Ineligible},
				{result: makeTestResult(2, "0x22", false, 1), state: ocr2keepers.Ineligible},
				{result: makeTestResult(3, "0x33", false, 1), state: ocr2keepers.Ineligible},
				{result: makeTestResult(4, "0x44", false, 1), state: ocr2keepers.Ineligible},
			},
			expected: []ocr2keepers.UpkeepState{
				StateUnknown,
				StateUnknown,
				StateUnknown,
				StateUnknown,
			},
		},
		{
			name:     "querying workIDs with non-stored values returns valid results",
			queryIDs: []string{"0x1", "0x2", "0x3", "0x4"},
			storedValues: []storedValue{
				{result: makeTestResult(5, "0x1", false, 1), state: ocr2keepers.Ineligible},
				{result: makeTestResult(6, "0x2", false, 1), state: ocr2keepers.Ineligible},
				{result: makeTestResult(7, "0x3", false, 1), state: ocr2keepers.Ineligible},
				{result: makeTestResult(8, "0x44", false, 1), state: ocr2keepers.Ineligible},
			},
			expected: []ocr2keepers.UpkeepState{
				ocr2keepers.Ineligible,
				ocr2keepers.Ineligible,
				ocr2keepers.Ineligible,
				StateUnknown,
			},
		},
		{
			name:     "storing eligible values is a noop",
			queryIDs: []string{"0x1", "0x2", "0x3", "0x4"},
			storedValues: []storedValue{
				{result: makeTestResult(9, "0x1", false, 1), state: ocr2keepers.Ineligible},
				{result: makeTestResult(10, "0x2", false, 1), state: ocr2keepers.Ineligible},
				{result: makeTestResult(11, "0x3", false, 1), state: ocr2keepers.Ineligible},
				{result: makeTestResult(12, "0x4", true, 1), state: ocr2keepers.Performed},
			},
			expected: []ocr2keepers.UpkeepState{
				ocr2keepers.Ineligible,
				ocr2keepers.Ineligible,
				ocr2keepers.Ineligible,
				StateUnknown,
			},
		},
		{
			name:     "provided state on setupkeepstate is currently ignored for eligible check results",
			queryIDs: []string{"0x1", "0x2"},
			storedValues: []storedValue{
				{result: makeTestResult(13, "0x1", true, 1), state: ocr2keepers.Ineligible},
				{result: makeTestResult(14, "0x2", false, 1), state: ocr2keepers.Performed},
			},
			expected: []ocr2keepers.UpkeepState{
				StateUnknown,
				ocr2keepers.Ineligible,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := testutils.Context(t)

			lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.ErrorLevel)
			chainID := testutils.FixtureChainID
			db := pgtest.NewSqlxDB(t)
			orm := NewORM(chainID, db, lggr, pgtest.NewQConfig(true))
			scanner := &mockScanner{}
			store := NewUpkeepStateStore(orm, lggr, scanner)

			t.Cleanup(func() {
				t.Log("cleaning up database")

				if _, err := db.Exec(`DELETE FROM evm_upkeep_state`); err != nil {
					t.Logf("error in cleanup: %s", err)
				}
			})

			for _, insert := range test.storedValues {
				require.NoError(t, store.SetUpkeepState(context.Background(), insert.result, insert.state), "storing states should not produce an error")
			}

			// empty the cache before doing selects to force a db lookup
			store.cache = make(map[string]*upkeepStateRecord)

			states, err := store.SelectByWorkIDsInRange(ctx, 1, 100, test.queryIDs...)

			require.NoError(t, err, "no error expected from selecting states")

			assert.Equal(t, test.expected, states, "upkeep state values should match expected")

			observedLogs.TakeAll()

			require.Equal(t, 0, observedLogs.Len())
		})
	}
}

func TestUpkeepStateStore_Upsert(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)
	chainID := testutils.FixtureChainID
	orm := NewORM(chainID, db, lggr, pgtest.NewQConfig(true))

	store := NewUpkeepStateStore(orm, lggr, &mockScanner{})

	res := ocr2keepers.CheckResult{
		UpkeepID: createUpkeepIDForTest(1),
		WorkID:   "0x1",
		Eligible: false,
		Trigger: ocr2keepers.Trigger{
			BlockNumber: ocr2keepers.BlockNumber(1),
		},
	}
	require.NoError(t, store.SetUpkeepState(ctx, res, ocr2keepers.Performed))
	<-time.After(10 * time.Millisecond)
	res.Trigger.BlockNumber = ocr2keepers.BlockNumber(2)
	now := time.Now()
	require.NoError(t, store.SetUpkeepState(ctx, res, ocr2keepers.Performed))

	store.mu.Lock()
	addedAt := store.cache["0x1"].AddedAt
	block := store.cache["0x1"].BlockNumber
	store.mu.Unlock()

	require.True(t, now.After(addedAt))
	require.Equal(t, uint64(2), block)
}

func TestUpkeepStateStore_Service(t *testing.T) {
	orm := &mockORM{
		onDelete: func(tm time.Time) {

		},
	}
	scanner := &mockScanner{}

	store := NewUpkeepStateStore(orm, logger.TestLogger(t), scanner)

	store.retention = 500 * time.Millisecond
	store.cleanCadence = 100 * time.Millisecond

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		assert.NoError(t, store.Start(context.Background()), "no error from starting service")

		wg.Done()
	}()

	// add a value to set up the test
	require.NoError(t, store.SetUpkeepState(context.Background(), ocr2keepers.CheckResult{
		Eligible: false,
		WorkID:   "0x2",
		Trigger: ocr2keepers.Trigger{
			BlockNumber: ocr2keepers.BlockNumber(1),
		},
	}, ocr2keepers.Ineligible))

	// allow one cycle of cleaning the cache
	time.Sleep(110 * time.Millisecond)

	// select from store to ensure values still exist
	values, err := store.SelectByWorkIDsInRange(context.Background(), 1, 100, "0x2")
	require.NoError(t, err, "no error from selecting states")
	require.Equal(t, []ocr2keepers.UpkeepState{ocr2keepers.Ineligible}, values, "selected values should match expected")

	// wait longer than cache timeout
	time.Sleep(700 * time.Millisecond)

	// select from store to ensure cached values were removed
	values, err = store.SelectByWorkIDsInRange(context.Background(), 1, 100, "0x2")
	require.NoError(t, err, "no error from selecting states")
	require.Equal(t, []ocr2keepers.UpkeepState{StateUnknown}, values, "selected values should match expected")

	assert.NoError(t, store.Close(), "no error from closing service")

	wg.Wait()
}

func createUpkeepIDForTest(v int64) ocr2keepers.UpkeepIdentifier {
	uid := &ocr2keepers.UpkeepIdentifier{}
	_ = uid.FromBigInt(big.NewInt(v))

	return *uid
}

type mockScanner struct {
	lock    sync.Mutex
	workIDs []string
	err     error
}

func (s *mockScanner) addWorkID(workIDs ...string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.workIDs = append(s.workIDs, workIDs...)
}

func (s *mockScanner) setErr(err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.err = err
}

func (s *mockScanner) WorkIDsInRange(ctx context.Context, start, end int64) ([]string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	res := s.workIDs[:]
	s.workIDs = nil
	return res, s.err
}

type mockORM struct {
	lock           sync.Mutex
	records        []PersistedStateRecord
	lastPruneDepth time.Time
	onDelete       func(tm time.Time)
	err            error
}

func (_m *mockORM) addRecords(records ...PersistedStateRecord) {
	_m.lock.Lock()
	defer _m.lock.Unlock()

	_m.records = append(_m.records, records...)
}

func (_m *mockORM) setErr(err error) {
	_m.lock.Lock()
	defer _m.lock.Unlock()

	_m.err = err
}

func (_m *mockORM) InsertUpkeepState(state PersistedStateRecord, _ ...pg.QOpt) error {
	return nil
}

func (_m *mockORM) SelectStatesByWorkIDs(workIDs []string, _ ...pg.QOpt) ([]PersistedStateRecord, error) {
	_m.lock.Lock()
	defer _m.lock.Unlock()

	res := _m.records[:]
	_m.records = nil

	return res, _m.err
}

func (_m *mockORM) DeleteExpired(tm time.Time, _ ...pg.QOpt) error {
	_m.lock.Lock()
	defer _m.lock.Unlock()

	_m.lastPruneDepth = tm
	_m.onDelete(tm)

	return _m.err
}
