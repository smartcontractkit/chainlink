package coordinator_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm/headtracker"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ocr2vrf/coordinator"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestOrm_HeadsByNumbers(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	lg := logger.TestLogger(t)
	cfg := cltest.NewTestGeneralConfig(t)
	orm := coordinator.NewORM(db, *utils.NewBigI(cltest.FixtureChainID.Int64()), lg)
	htORM := headtracker.NewORM(db, lg, cfg, cltest.FixtureChainID)

	var expectedHashes []common.Hash
	for i := 0; i < 10; i++ {
		head := cltest.Head(i + 1)
		require.NoError(t, htORM.IdempotentInsertHead(testutils.Context(t), head))
		expectedHashes = append(expectedHashes, head.Hash)
	}

	dbHeads, err := orm.HeadsByNumbers(testutils.Context(t), []uint64{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
	})
	require.NoError(t, err)
	require.Len(t, dbHeads, len(expectedHashes))
	var dbHashes []common.Hash

	for _, dbh := range dbHeads {
		dbHashes = append(dbHashes, dbh.Hash)
	}
	require.ElementsMatch(t, expectedHashes, dbHashes)
}
