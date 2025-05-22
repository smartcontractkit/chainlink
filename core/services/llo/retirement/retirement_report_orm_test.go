package retirement

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

func Test_RetirementReportCache_ORM(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	orm := &retirementReportCacheORM{db}
	ctx := t.Context()

	cd := ocr2types.ConfigDigest{1}
	attestedRetirementReport := []byte("report1")
	cd2 := ocr2types.ConfigDigest{2}
	attestedRetirementReport2 := []byte("report2")

	t.Run("StoreAttestedRetirementReport", func(t *testing.T) {
		err := orm.StoreAttestedRetirementReport(ctx, cd, attestedRetirementReport)
		require.NoError(t, err)
		err = orm.StoreAttestedRetirementReport(ctx, cd2, attestedRetirementReport2)
		require.NoError(t, err)
	})
	t.Run("LoadAttestedRetirementReports", func(t *testing.T) {
		arrs, err := orm.LoadAttestedRetirementReports(ctx)
		require.NoError(t, err)

		require.Len(t, arrs, 2)
		assert.Equal(t, attestedRetirementReport, arrs[cd])
		assert.Equal(t, attestedRetirementReport2, arrs[cd2])
	})
	t.Run("StoreConfig", func(t *testing.T) {
		signers := [][]byte{[]byte("signer1"), []byte("signer2")}
		err := orm.StoreConfig(ctx, cd, signers, 1)
		require.NoError(t, err)

		err = orm.StoreConfig(ctx, cd2, signers, 2)
		require.NoError(t, err)

		// overwriting does nothing, the 255 is ignored
		err = orm.StoreConfig(ctx, cd2, signers, 255)
		require.NoError(t, err)
	})
	t.Run("LoadConfigs", func(t *testing.T) {
		configs, err := orm.LoadConfigs(ctx)
		require.NoError(t, err)

		require.Len(t, configs, 2)
		assert.Equal(t, Config{
			Digest:  cd,
			Signers: [][]byte{[]byte("signer1"), []byte("signer2")},
			F:       1,
		}, configs[0])
		assert.Equal(t, Config{
			Digest:  cd2,
			Signers: [][]byte{[]byte("signer1"), []byte("signer2")},
			F:       2,
		}, configs[1])
	})
}
