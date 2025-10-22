package round_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	"github.com/smartcontractkit/chainlink/v2/round"
)

func Test_DB_LatestRoundRequested(t *testing.T) {
	sqlDB := testutils.NewSqlxDB(t)

	_, err := sqlDB.Exec(`SET CONSTRAINTS offchainreporting2_latest_round_oracle_spec_fkey DEFERRED`)
	require.NoError(t, err)

	lggr := logger.Test(t)
	db := round.NewRoundRequestedDB(sqlDB, 1, lggr)
	db2 := round.NewRoundRequestedDB(sqlDB, 2, lggr)

	rawLog := logFromFixture(t, "./testdata/round_requested_log_1_1.json")

	rr := ocr2aggregator.OCR2AggregatorRoundRequested{
		Requester:    testutils.NewAddress(),
		ConfigDigest: ocrtypes.ConfigDigest{},
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
			ConfigDigest: ocrtypes.ConfigDigest{},
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

func logFromFixture(t *testing.T, path string) types.Log {
	value := gjson.Get(string(mustReadFile(t, path)), "params.result")
	var el types.Log
	require.NoError(t, json.Unmarshal([]byte(value.String()), &el))

	return el
}

func mustReadFile(t testing.TB, file string) []byte {
	t.Helper()

	content, err := os.ReadFile(file)
	require.NoError(t, err)
	return content
}
