package upkeepstate

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

func TestInsertSelectDelete(t *testing.T) {
	ctx := testutils.Context(t)
	chainID := testutils.FixtureChainID
	db := pgtest.NewSqlxDB(t)
	orm := NewORM(chainID, db)

	inserted := []persistedStateRecord{
		{
			UpkeepID:            ubig.New(big.NewInt(2)),
			WorkID:              "0x1",
			CompletionState:     100,
			BlockNumber:         2,
			IneligibilityReason: 2,
			InsertedAt:          time.Now(),
		},
	}

	err := orm.BatchInsertRecords(ctx, inserted)

	require.NoError(t, err, "no error expected from insert")

	states, err := orm.SelectStatesByWorkIDs(ctx, []string{"0x1"})

	require.NoError(t, err, "no error expected from select")
	require.Len(t, states, 1, "records return should equal records inserted")

	err = orm.DeleteExpired(ctx, time.Now())

	assert.NoError(t, err, "no error expected from delete")

	states, err = orm.SelectStatesByWorkIDs(ctx, []string{"0x1"})

	require.NoError(t, err, "no error expected from select")
	require.Len(t, states, 0, "records return should be empty since records were deleted")
}
