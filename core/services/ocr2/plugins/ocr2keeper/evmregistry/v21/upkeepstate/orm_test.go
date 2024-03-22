package upkeepstate

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestInsertSelectDelete(t *testing.T) {
	lggr, _ := logger.TestLoggerObserved(t, zapcore.ErrorLevel)
	chainID := testutils.FixtureChainID
	db := pgtest.NewSqlxDB(t)
	orm := NewORM(chainID, db, lggr, pgtest.NewQConfig(true))

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

	err := orm.BatchInsertRecords(inserted)

	require.NoError(t, err, "no error expected from insert")

	states, err := orm.SelectStatesByWorkIDs([]string{"0x1"})

	require.NoError(t, err, "no error expected from select")
	require.Len(t, states, 1, "records return should equal records inserted")

	err = orm.DeleteExpired(time.Now())

	assert.NoError(t, err, "no error expected from delete")

	states, err = orm.SelectStatesByWorkIDs([]string{"0x1"})

	require.NoError(t, err, "no error expected from select")
	require.Len(t, states, 0, "records return should be empty since records were deleted")
}
