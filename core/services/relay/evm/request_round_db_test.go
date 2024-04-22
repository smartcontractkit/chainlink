package evm_test

import (
	"testing"

	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/testhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

func Test_DB_LatestRoundRequested(t *testing.T) {
	sqlDB := pgtest.NewSqlxDB(t)

	_, err := sqlDB.Exec(`SET CONSTRAINTS offchainreporting2_latest_round_oracle_spec_fkey DEFERRED`)
	require.NoError(t, err)

	lggr := logger.TestLogger(t)
	db := evm.NewRoundRequestedDB(sqlDB, 1, lggr)
	db2 := evm.NewRoundRequestedDB(sqlDB, 2, lggr)

	rawLog := cltest.LogFromFixture(t, "../../../testdata/jsonrpc/round_requested_log_1_1.json")

	rr := ocr2aggregator.OCR2AggregatorRoundRequested{
		Requester:    testutils.NewAddress(),
		ConfigDigest: testhelpers.MakeConfigDigest(t),
		Epoch:        42,
		Round:        9,
		Raw:          rawLog,
	}

	t.Run("saves latest round requested", func(t *testing.T) {
		ctx := testutils.Context(t)
		err := db.SaveLatestRoundRequested(ctx, rr)
		require.NoError(t, err)

		rawLog.Index = 42

		// Now overwrite to prove that updating works
		rr = ocr2aggregator.OCR2AggregatorRoundRequested{
			Requester:    testutils.NewAddress(),
			ConfigDigest: testhelpers.MakeConfigDigest(t),
			Epoch:        43,
			Round:        8,
			Raw:          rawLog,
		}

		err = db.SaveLatestRoundRequested(ctx, rr)
		require.NoError(t, err)
	})

	t.Run("loads latest round requested", func(t *testing.T) {
		ctx := testutils.Context(t)
		// There is no round for db2
		lrr, err := db2.LoadLatestRoundRequested(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, int(lrr.Epoch))

		lrr, err = db.LoadLatestRoundRequested(ctx)
		require.NoError(t, err)

		assert.Equal(t, rr, lrr)
	})

	t.Run("spec with latest round requested can be deleted", func(t *testing.T) {
		_, err := sqlDB.Exec(`DELETE FROM ocr2_oracle_specs`)
		assert.NoError(t, err)
	})
}
