package evm

import (
	"context"
	"math/big"
	"testing"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

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
			lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.ErrorLevel)
			chainID := testutils.FixtureChainID
			db := pgtest.NewSqlxDB(t)
			orm := NewORM(chainID, db, lggr, pgtest.NewQConfig(true))
			store := NewUpkeepStateStore(orm, lggr)

			t.Cleanup(func() {
				t.Log("cleaning up database")

				if _, err := db.Exec(`DELETE FROM evm_upkeep_state`); err != nil {
					t.Logf("error in cleanup: %s", err)
				}
			})

			for _, insert := range test.storedValues {
				require.NoError(t, store.SetUpkeepState(context.Background(), insert.result, insert.state), "storing states should not produce an error")
			}

			states, err := store.SelectByWorkIDs(test.queryIDs...)

			require.NoError(t, err, "no error expected from selecting states")

			assert.Equal(t, test.expected, states, "upkeep state values should match expected")

			observedLogs.TakeAll()

			require.Equal(t, 0, observedLogs.Len())
		})
	}
}
