package evm

import (
	"context"
	"math/big"
	"testing"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestUpkeepStateStore_SelectByWorkIDs(t *testing.T) {
	workIDs := []string{"a", "b", "c", "d"}
	inserts := []ocr2keepers.CheckResult{
		{
			UpkeepID: createUpkeepIDForTest(1),
			WorkID:   workIDs[0],
			Eligible: false,
			Trigger: ocr2keepers.Trigger{
				BlockNumber: ocr2keepers.BlockNumber(1),
			},
		},
		{
			UpkeepID: createUpkeepIDForTest(2),
			WorkID:   workIDs[1],
			Eligible: false,
			Trigger: ocr2keepers.Trigger{
				BlockNumber: ocr2keepers.BlockNumber(2),
			},
		},
		{
			UpkeepID: createUpkeepIDForTest(3),
			WorkID:   workIDs[2],
			Eligible: false,
			Trigger: ocr2keepers.Trigger{
				BlockNumber: ocr2keepers.BlockNumber(3),
			},
		},
	}

	expected := map[string]ocr2keepers.UpkeepState{
		workIDs[0]: ocr2keepers.Ineligible,
		workIDs[1]: ocr2keepers.Ineligible,
		workIDs[2]: ocr2keepers.Ineligible,
	}

	store := NewUpkeepStateStore(logger.TestLogger(t))

	for _, insert := range inserts {
		assert.NoError(t, store.SetUpkeepState(context.Background(), insert, ocr2keepers.Performed), "storing states should not produce an error")
	}

	states, err := store.SelectByWorkIDs(workIDs...)
	assert.NoError(t, err, "no error expected from selecting states")

	assert.Equal(t, expected, states, "upkeep state values should match expected")
}

func TestUpkeepStateStore_SetUpkeepState(t *testing.T) {
	t.Run("should not save state for upkeep eligible", func(t *testing.T) {
		uid := &ocr2keepers.UpkeepIdentifier{}
		_ = uid.FromBigInt(big.NewInt(1))

		store := NewUpkeepStateStore(logger.TestLogger(t))

		assert.NoError(t, store.SetUpkeepState(context.Background(), ocr2keepers.CheckResult{
			UpkeepID: *uid,
			WorkID:   "test",
			Eligible: true,
		}, ocr2keepers.Ineligible), "setting state should not return an error")

		assert.Len(t, store.workIDIdx, 0, "should not add to upkeep states")
	})

	t.Run("should insert new state when ineligible and state does not exist in store and ignore state input", func(t *testing.T) {
		uid := &ocr2keepers.UpkeepIdentifier{}
		_ = uid.FromBigInt(big.NewInt(1))

		store := NewUpkeepStateStore(logger.TestLogger(t))
		input := ocr2keepers.CheckResult{
			UpkeepID: *uid,
			WorkID:   "test",
			Trigger: ocr2keepers.Trigger{
				BlockNumber: ocr2keepers.BlockNumber(1),
			},
			Eligible: false,
		}

		assert.NoError(t, store.SetUpkeepState(context.Background(), input, ocr2keepers.Performed))

		require.Len(t, store.workIDIdx, 1, "should add to upkeep states")

		assert.Equal(t, ocr2keepers.Ineligible, store.workIDIdx["test"].state, "stored state should be ineligible")
		assert.Equal(t, input.UpkeepID, store.workIDIdx["test"].upkeepID, "stored upkeepID should match input")
		assert.Equal(t, input.WorkID, store.workIDIdx["test"].workID, "stored workID should match input")
		assert.Equal(t, uint64(input.Trigger.BlockNumber), store.workIDIdx["test"].block, "stored block should match input")
	})

	// when eligible and state exists in store, override state, ignore state input
	t.Run("should override block when ineligible and state exists in store and ignore state input", func(t *testing.T) {
		store := NewUpkeepStateStore(logger.TestLogger(t))
		input := ocr2keepers.CheckResult{
			UpkeepID: createUpkeepIDForTest(1),
			WorkID:   "test",
			Trigger: ocr2keepers.Trigger{
				BlockNumber: ocr2keepers.BlockNumber(1),
			},
			Eligible: false,
		}

		assert.NoError(t, store.SetUpkeepState(context.Background(), input, ocr2keepers.Performed), "setting state should not return an error")

		require.Len(t, store.workIDIdx, 1, "should add to upkeep states")

		// update the block number for the input to indicate a state data change
		input.Trigger.BlockNumber = ocr2keepers.BlockNumber(5)

		assert.NoError(t, store.SetUpkeepState(context.Background(), input, ocr2keepers.Performed), "setting state should not return an error")

		require.Len(t, store.workIDIdx, 1, "should update existing upkeep state")

		assert.Equal(t, ocr2keepers.Ineligible, store.workIDIdx["test"].state, "stored state should be ineligible")
		assert.Equal(t, input.UpkeepID, store.workIDIdx["test"].upkeepID, "stored upkeepID should match input")
		assert.Equal(t, input.WorkID, store.workIDIdx["test"].workID, "stored workID should match input")
		assert.Equal(t, uint64(input.Trigger.BlockNumber), store.workIDIdx["test"].block, "stored block should match input")
	})
}

func createUpkeepIDForTest(v int64) ocr2keepers.UpkeepIdentifier {
	uid := &ocr2keepers.UpkeepIdentifier{}
	_ = uid.FromBigInt(big.NewInt(v))

	return *uid
}
