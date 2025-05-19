package upkeepstate

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestUpkeepStateStore(t *testing.T) {
	tests := []struct {
		name               string
		inserts            []ocr2keepers.CheckResult
		workIDsSelect      []string
		workIDsFromScanner []string
		errScanner         error
		recordsFromDB      []persistedStateRecord
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
				ocr2keepers.UnknownState,
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
			recordsFromDB: []persistedStateRecord{
				{
					WorkID:          "0x3",
					CompletionState: 2,
					BlockNumber:     2,
					InsertedAt:      time.Now(),
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
				ocr2keepers.UnknownState,
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
			errScanner:         errors.New("test error"),
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
			errDB:              errors.New("test error"),
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

			states, err := store.SelectByWorkIDs(ctx, tc.workIDsSelect...)
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
		name           string
		flushSize      int
		expectedWrites int
		queryIDs       []string
		storedValues   []storedValue
		expected       []ocr2keepers.UpkeepState
	}{
		{
			name:           "querying non-stored workIDs on db with values returns unknown state results",
			queryIDs:       []string{"0x1", "0x2", "0x3", "0x4"},
			flushSize:      10,
			expectedWrites: 1,
			storedValues: []storedValue{
				{result: makeTestResult(1, "0x11", false, 1), state: ocr2keepers.Performed},
				{result: makeTestResult(2, "0x22", false, 1), state: ocr2keepers.Performed},
				{result: makeTestResult(3, "0x33", false, 1), state: ocr2keepers.Performed},
				{result: makeTestResult(4, "0x44", false, 1), state: ocr2keepers.Performed},
			},
			expected: []ocr2keepers.UpkeepState{
				ocr2keepers.UnknownState,
				ocr2keepers.UnknownState,
				ocr2keepers.UnknownState,
				ocr2keepers.UnknownState,
			},
		},
		{
			name:           "storing eligible values is a noop",
			queryIDs:       []string{"0x1", "0x2", "0x3", "0x4"},
			flushSize:      4,
			expectedWrites: 1,
			storedValues: []storedValue{
				{result: makeTestResult(9, "0x1", false, 1), state: ocr2keepers.Ineligible},
				{result: makeTestResult(10, "0x2", false, 1), state: ocr2keepers.Ineligible},
				{result: makeTestResult(11, "0x3", false, 1), state: ocr2keepers.Ineligible},
				{result: makeTestResult(12, "0x4", true, 1), state: ocr2keepers.Performed}, // gets inserted
			},
			expected: []ocr2keepers.UpkeepState{
				ocr2keepers.Ineligible,
				ocr2keepers.Ineligible,
				ocr2keepers.Ineligible,
				ocr2keepers.UnknownState,
			},
		},
		{
			name:           "provided state on setupkeepstate is currently ignored for eligible check results",
			queryIDs:       []string{"0x1", "0x2"},
			flushSize:      1,
			expectedWrites: 1,
			storedValues: []storedValue{
				{result: makeTestResult(13, "0x1", true, 1), state: ocr2keepers.Ineligible},
				{result: makeTestResult(14, "0x2", false, 1), state: ocr2keepers.Performed}, // gets inserted
			},
			expected: []ocr2keepers.UpkeepState{
				ocr2keepers.UnknownState,
				ocr2keepers.Ineligible,
			},
		},
		{
			name:           "provided state outside the flush batch isn't registered in the db",
			queryIDs:       []string{"0x1", "0x2", "0x3", "0x4", "0x5", "0x6", "0x7", "0x8"},
			flushSize:      3,
			expectedWrites: 2,
			storedValues: []storedValue{
				{result: makeTestResult(13, "0x1", true, 1), state: ocr2keepers.Ineligible},
				{result: makeTestResult(14, "0x2", false, 1), state: ocr2keepers.Performed}, // gets inserted
				{result: makeTestResult(15, "0x3", true, 1), state: ocr2keepers.Ineligible},
				{result: makeTestResult(16, "0x4", false, 1), state: ocr2keepers.Performed}, // gets inserted
				{result: makeTestResult(17, "0x5", true, 1), state: ocr2keepers.Ineligible},
				{result: makeTestResult(18, "0x6", false, 1), state: ocr2keepers.Performed}, // gets inserted
				{result: makeTestResult(19, "0x7", true, 1), state: ocr2keepers.Ineligible},
				{result: makeTestResult(20, "0x8", false, 1), state: ocr2keepers.Performed}, // gets inserted
			},
			expected: []ocr2keepers.UpkeepState{
				ocr2keepers.UnknownState,
				ocr2keepers.Ineligible,
				ocr2keepers.UnknownState,
				ocr2keepers.Ineligible,
				ocr2keepers.UnknownState,
				ocr2keepers.Ineligible,
				ocr2keepers.UnknownState,
				ocr2keepers.Ineligible,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := testutils.Context(t)

			tickerCh := make(chan time.Time)

			oldNewTickerFn := newTickerFn
			oldFlushSize := batchSize
			newTickerFn = func(d time.Duration) *time.Ticker {
				t := time.NewTicker(d)
				t.C = tickerCh
				return t
			}
			batchSize = test.flushSize
			defer func() {
				newTickerFn = oldNewTickerFn
				batchSize = oldFlushSize
			}()

			lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.ErrorLevel)
			chainID := testutils.FixtureChainID
			db := pgtest.NewSqlxDB(t)
			realORM := NewORM(chainID, db)
			insertFinished := make(chan struct{}, 1)
			orm := &wrappedORM{
				BatchInsertRecordsFn: func(ctx context.Context, records []persistedStateRecord) error {
					err := realORM.BatchInsertRecords(ctx, records)
					insertFinished <- struct{}{}
					return err
				},
				SelectStatesByWorkIDsFn: realORM.SelectStatesByWorkIDs,
				DeleteExpiredFn:         realORM.DeleteExpired,
			}
			scanner := &mockScanner{}
			store := NewUpkeepStateStore(orm, lggr, scanner)

			servicetest.Run(t, store)

			t.Cleanup(func() {
				t.Log("cleaning up database")

				if _, err := db.Exec(`DELETE FROM evm.upkeep_states`); err != nil {
					t.Logf("error in cleanup: %s", err)
				}
			})

			for _, insert := range test.storedValues {
				require.NoError(t, store.SetUpkeepState(ctx, insert.result, insert.state), "storing states should not produce an error")
			}

			tickerCh <- time.Now()

			// if this test inserts data, wait for the insert to complete before proceeding
			for i := 0; i < test.expectedWrites; i++ {
				<-insertFinished
			}

			// empty the cache before doing selects to force a db lookup
			store.cache = make(map[string]*upkeepStateRecord)

			states, err := store.SelectByWorkIDs(ctx, test.queryIDs...)

			require.NoError(t, err, "no error expected from selecting states")

			assert.Equal(t, test.expected, states, "upkeep state values should match expected")

			observedLogs.TakeAll()

			require.Equal(t, 0, observedLogs.Len())
		})
	}
}

func TestUpkeepStateStore_emptyDB(t *testing.T) {
	t.Run("querying non-stored workIDs on empty db returns unknown state results", func(t *testing.T) {
		lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.ErrorLevel)
		chainID := testutils.FixtureChainID
		db := pgtest.NewSqlxDB(t)
		realORM := NewORM(chainID, db)
		insertFinished := make(chan struct{}, 1)
		orm := &wrappedORM{
			BatchInsertRecordsFn: func(ctx context.Context, records []persistedStateRecord) error {
				err := realORM.BatchInsertRecords(ctx, records)
				insertFinished <- struct{}{}
				return err
			},
			SelectStatesByWorkIDsFn: realORM.SelectStatesByWorkIDs,
			DeleteExpiredFn:         realORM.DeleteExpired,
		}
		scanner := &mockScanner{}
		store := NewUpkeepStateStore(orm, lggr, scanner)

		states, err := store.SelectByWorkIDs(testutils.Context(t), []string{"0x1", "0x2", "0x3", "0x4"}...)
		assert.NoError(t, err)
		assert.Equal(t, []ocr2keepers.UpkeepState{
			ocr2keepers.UnknownState,
			ocr2keepers.UnknownState,
			ocr2keepers.UnknownState,
			ocr2keepers.UnknownState,
		}, states)

		observedLogs.TakeAll()

		require.Equal(t, 0, observedLogs.Len())
	})
}

func TestUpkeepStateStore_Upsert(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)
	chainID := testutils.FixtureChainID
	orm := NewORM(chainID, db)

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
	addedAt := store.cache["0x1"].addedAt
	store.mu.Unlock()

	require.True(t, now.After(addedAt))
}

func TestUpkeepStateStore_Service(t *testing.T) {
	ctx := testutils.Context(t)
	orm := &mockORM{
		onDelete: func(tm time.Time) {

		},
	}
	scanner := &mockScanner{}

	store := NewUpkeepStateStore(orm, logger.TestLogger(t), scanner)

	store.retention = 500 * time.Millisecond
	store.cleanCadence = 100 * time.Millisecond

	servicetest.Run(t, store)

	// add a value to set up the test
	require.NoError(t, store.SetUpkeepState(ctx, ocr2keepers.CheckResult{
		Eligible: false,
		WorkID:   "0x2",
		Trigger: ocr2keepers.Trigger{
			BlockNumber: ocr2keepers.BlockNumber(1),
		},
	}, ocr2keepers.Ineligible))

	// allow one cycle of cleaning the cache
	time.Sleep(110 * time.Millisecond)

	// select from store to ensure values still exist
	values, err := store.SelectByWorkIDs(ctx, "0x2")
	require.NoError(t, err, "no error from selecting states")
	require.Equal(t, []ocr2keepers.UpkeepState{ocr2keepers.Ineligible}, values, "selected values should match expected")

	// wait longer than cache timeout
	time.Sleep(700 * time.Millisecond)

	// select from store to ensure cached values were removed
	values, err = store.SelectByWorkIDs(ctx, "0x2")
	require.NoError(t, err, "no error from selecting states")
	require.Equal(t, []ocr2keepers.UpkeepState{ocr2keepers.UnknownState}, values, "selected values should match expected")
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

func (s *mockScanner) ScanWorkIDs(context.Context, ...string) ([]string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	res := s.workIDs
	s.workIDs = nil
	return res, s.err
}

func (s *mockScanner) Start(context.Context) error {
	return nil
}

func (s *mockScanner) Close() error {
	return nil
}

type mockORM struct {
	lock           sync.Mutex
	records        []persistedStateRecord
	lastPruneDepth time.Time
	onDelete       func(tm time.Time)
	err            error
}

func (_m *mockORM) addRecords(records ...persistedStateRecord) {
	_m.lock.Lock()
	defer _m.lock.Unlock()

	_m.records = append(_m.records, records...)
}

func (_m *mockORM) setErr(err error) {
	_m.lock.Lock()
	defer _m.lock.Unlock()

	_m.err = err
}

func (_m *mockORM) BatchInsertRecords(ctx context.Context, state []persistedStateRecord) error {
	return nil
}

func (_m *mockORM) SelectStatesByWorkIDs(ctx context.Context, workIDs []string) ([]persistedStateRecord, error) {
	_m.lock.Lock()
	defer _m.lock.Unlock()

	res := _m.records
	_m.records = nil

	return res, _m.err
}

func (_m *mockORM) DeleteExpired(ctx context.Context, tm time.Time) error {
	_m.lock.Lock()
	defer _m.lock.Unlock()

	_m.lastPruneDepth = tm
	_m.onDelete(tm)

	return _m.err
}

type wrappedORM struct {
	BatchInsertRecordsFn    func(context.Context, []persistedStateRecord) error
	SelectStatesByWorkIDsFn func(context.Context, []string) ([]persistedStateRecord, error)
	DeleteExpiredFn         func(context.Context, time.Time) error
}

func (o *wrappedORM) BatchInsertRecords(ctx context.Context, r []persistedStateRecord) error {
	return o.BatchInsertRecordsFn(ctx, r)
}

func (o *wrappedORM) SelectStatesByWorkIDs(ctx context.Context, ids []string) ([]persistedStateRecord, error) {
	return o.SelectStatesByWorkIDsFn(ctx, ids)
}

func (o *wrappedORM) DeleteExpired(ctx context.Context, t time.Time) error {
	return o.DeleteExpiredFn(ctx, t)
}
