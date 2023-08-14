package upkeepstate

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestInsertSelectDelete(t *testing.T) {
	lggr, _ := logger.TestLoggerObserved(t, zapcore.ErrorLevel)
	chainID := testutils.FixtureChainID
	db := pgtest.NewSqlxDB(t)
	orm := NewORM(chainID, db, lggr, pgtest.NewQConfig(true))

	inserted := upkeepStateRecord{
		WorkID:          "0x1",
		CompletionState: ocr2keepers.UpkeepState(100),
		BlockNumber:     2,
		AddedAt:         time.Now(),
	}

	err := orm.InsertUpkeepState(inserted)

	require.NoError(t, err, "no error expected from insert")

	states, err := orm.SelectStatesByWorkIDs([]string{"0x1"})

	require.NoError(t, err, "no error expected from select")
	require.Len(t, states, 1, "records return should equal records inserted")

	err = orm.DeleteBeforeTime(time.Now())

	assert.NoError(t, err, "no error expected from delete")

	states, err = orm.SelectStatesByWorkIDs([]string{"0x1"})

	require.NoError(t, err, "no error expected from select")
	require.Len(t, states, 0, "records return should be empty since records were deleted")
}
