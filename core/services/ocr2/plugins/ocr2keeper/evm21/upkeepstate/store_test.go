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
)

func TestUpkeepStateStore(t *testing.T) {
	tests := []struct {
		name               string
		inserts            []ocr2keepers.CheckResult
		workIDsSelect      []string
		workIDsFromScanner []string
		errScanner         error
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
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testutils.Context(t)

			scanner := &mockScanner{}
			scanner.addWorkID(tc.workIDsFromScanner...)
			scanner.setErr(tc.errScanner)
			store := NewUpkeepStateStore(nil, logger.TestLogger(t), scanner)

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
		expected     map[string]ocr2keepers.UpkeepState
	}{
		{
			name:     "querying non-stored workIDs on empty db returns empty results",
			queryIDs: []string{"a", "b", "c", "d"},
			expected: make(map[string]ocr2keepers.UpkeepState),
		},
		{
			name:     "querying non-stored workIDs on db with values returns empty results",
			queryIDs: []string{"a", "b", "c", "d"},
			storedValues: []storedValue{
				{result: makeTestResult(1, "e", false, 1), state: ocr2keepers.Ineligible},
				{result: makeTestResult(2, "f", false, 1), state: ocr2keepers.Ineligible},
				{result: makeTestResult(3, "g", false, 1), state: ocr2keepers.Ineligible},
				{result: makeTestResult(4, "h", false, 1), state: ocr2keepers.Ineligible},
			},
			expected: make(map[string]ocr2keepers.UpkeepState),
		},
		{
			name:     "querying workIDs with non-stored values returns valid results",
			queryIDs: []string{"a", "b", "c", "d"},
			storedValues: []storedValue{
				{result: makeTestResult(5, "a", false, 1), state: ocr2keepers.Ineligible},
				{result: makeTestResult(6, "b", false, 1), state: ocr2keepers.Ineligible},
				{result: makeTestResult(7, "c", false, 1), state: ocr2keepers.Ineligible},
				{result: makeTestResult(8, "h", false, 1), state: ocr2keepers.Ineligible},
			},
			expected: map[string]ocr2keepers.UpkeepState{
				"a": ocr2keepers.Ineligible,
				"b": ocr2keepers.Ineligible,
				"c": ocr2keepers.Ineligible,
			},
		},
		{
			name:     "storing eligible values is a noop",
			queryIDs: []string{"a", "b", "c", "d"},
			storedValues: []storedValue{
				{result: makeTestResult(9, "a", false, 1), state: ocr2keepers.Ineligible},
				{result: makeTestResult(10, "b", false, 1), state: ocr2keepers.Ineligible},
				{result: makeTestResult(11, "c", false, 1), state: ocr2keepers.Ineligible},
				{result: makeTestResult(12, "d", true, 1), state: ocr2keepers.Performed},
			},
			expected: map[string]ocr2keepers.UpkeepState{
				"a": ocr2keepers.Ineligible,
				"b": ocr2keepers.Ineligible,
				"c": ocr2keepers.Ineligible,
			},
		},
		{
			name:     "provided state on setupkeepstate is currently ignored for eligible check results",
			queryIDs: []string{"a", "b"},
			storedValues: []storedValue{
				{result: makeTestResult(13, "a", true, 1), state: ocr2keepers.Ineligible},
				{result: makeTestResult(14, "b", false, 1), state: ocr2keepers.Performed},
			},
			expected: map[string]ocr2keepers.UpkeepState{
				"b": ocr2keepers.Ineligible,
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

			states, err := store.SelectByWorkIDsInRange(ctx, 1, 100, test.queryIDs...)

			require.NoError(t, err, "no error expected from selecting states")

			assert.Equal(t, test.expected, states, "upkeep state values should match expected")

			observedLogs.TakeAll()

			require.Equal(t, 0, observedLogs.Len())
		})
	}
}

func TestUpkeepStateStore_Upsert(t *testing.T) {
	ctx := testutils.Context(t)
	store := NewUpkeepStateStore(nil, logger.TestLogger(t), &mockScanner{})

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
