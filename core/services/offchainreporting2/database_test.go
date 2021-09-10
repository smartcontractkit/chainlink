package offchainreporting2_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	offchainreporting "github.com/smartcontractkit/chainlink/core/services/offchainreporting2"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var ctx = context.Background()

func setupDB(t *testing.T) (*gorm.DB, *sql.DB) {
	t.Helper()

	gormDB := pgtest.NewGormDB(t)
	sqlDB, err := gormDB.DB()
	require.NoError(t, err)

	return gormDB, sqlDB
}

func Test_DB_ReadWriteState(t *testing.T) {
	gormDB, sqlDB := setupDB(t)

	configDigest := MakeConfigDigest(t)
	ethKeyStore := cltest.NewKeyStore(t, gormDB).Eth()
	key, _ := cltest.MustInsertRandomKey(t, ethKeyStore)
	spec := MustInsertOffchainreportingOracleSpec(t, gormDB, key.Address)

	t.Run("reads and writes state", func(t *testing.T) {
		db := offchainreporting.NewDB(sqlDB, spec.ID)
		state := ocrtypes.PersistentState{
			Epoch:                1,
			HighestSentEpoch:     2,
			HighestReceivedEpoch: []uint32{3},
		}

		err := db.WriteState(ctx, configDigest, state)
		require.NoError(t, err)

		readState, err := db.ReadState(ctx, configDigest)
		require.NoError(t, err)

		require.Equal(t, state, *readState)
	})

	t.Run("updates state", func(t *testing.T) {
		db := offchainreporting.NewDB(sqlDB, spec.ID)
		newState := ocrtypes.PersistentState{
			Epoch:                2,
			HighestSentEpoch:     3,
			HighestReceivedEpoch: []uint32{4, 5},
		}

		err := db.WriteState(ctx, configDigest, newState)
		require.NoError(t, err)

		readState, err := db.ReadState(ctx, configDigest)
		require.NoError(t, err)

		require.Equal(t, newState, *readState)
	})

	t.Run("does not return result for wrong spec", func(t *testing.T) {
		db := offchainreporting.NewDB(sqlDB, spec.ID)
		state := ocrtypes.PersistentState{
			Epoch:                3,
			HighestSentEpoch:     4,
			HighestReceivedEpoch: []uint32{5, 6},
		}

		err := db.WriteState(ctx, configDigest, state)
		require.NoError(t, err)

		// db with different spec
		db = offchainreporting.NewDB(sqlDB, -1)

		readState, err := db.ReadState(ctx, configDigest)
		require.NoError(t, err)

		require.Nil(t, readState)
	})

	t.Run("does not return result for wrong config digest", func(t *testing.T) {
		db := offchainreporting.NewDB(sqlDB, spec.ID)
		state := ocrtypes.PersistentState{
			Epoch:                4,
			HighestSentEpoch:     5,
			HighestReceivedEpoch: []uint32{6, 7},
		}

		err := db.WriteState(ctx, configDigest, state)
		require.NoError(t, err)

		readState, err := db.ReadState(ctx, MakeConfigDigest(t))
		require.NoError(t, err)

		require.Nil(t, readState)
	})
}

func Test_DB_ReadWriteConfig(t *testing.T) {
	gormDB, sqlDB := setupDB(t)

	config := ocrtypes.ContractConfig{
		ConfigDigest:          MakeConfigDigest(t),
		ConfigCount:           1,
		Signers:               []ocrtypes.OnchainPublicKey{},
		Transmitters:          []ocrtypes.Account{"account1"},
		F:                     79,
		OnchainConfig:         []byte{},
		OffchainConfigVersion: 111,
		OffchainConfig:        []byte{},
	}
	ethKeyStore := cltest.NewKeyStore(t, gormDB).Eth()
	key, _ := cltest.MustInsertRandomKey(t, ethKeyStore)
	spec := MustInsertOffchainreportingOracleSpec(t, gormDB, key.Address)

	t.Run("reads and writes config", func(t *testing.T) {
		db := offchainreporting.NewDB(sqlDB, spec.ID)

		err := db.WriteConfig(ctx, config)
		require.NoError(t, err)

		readConfig, err := db.ReadConfig(ctx)
		require.NoError(t, err)

		require.Equal(t, &config, readConfig)
	})

	t.Run("updates config", func(t *testing.T) {
		db := offchainreporting.NewDB(sqlDB, spec.ID)

		newConfig := ocrtypes.ContractConfig{
			ConfigDigest: MakeConfigDigest(t),
			Signers:      []ocrtypes.OnchainPublicKey{},
			Transmitters: []ocrtypes.Account{},
		}

		err := db.WriteConfig(ctx, newConfig)
		require.NoError(t, err)

		readConfig, err := db.ReadConfig(ctx)
		require.NoError(t, err)

		require.Equal(t, &newConfig, readConfig)
	})

	t.Run("does not return result for wrong spec", func(t *testing.T) {
		db := offchainreporting.NewDB(sqlDB, spec.ID)

		err := db.WriteConfig(ctx, config)
		require.NoError(t, err)

		db = offchainreporting.NewDB(sqlDB, -1)

		readConfig, err := db.ReadConfig(ctx)
		require.NoError(t, err)

		require.Nil(t, readConfig)
	})
}

func assertPendingTransmissionEqual(t *testing.T, pt1, pt2 ocrtypes.PendingTransmission) {
	t.Helper()

	require.Equal(t, pt1.Time.Unix(), pt2.Time.Unix())
	require.Equal(t, pt1.ExtraHash, pt2.ExtraHash)
	require.Equal(t, pt1.Report, pt2.Report)
	require.Equal(t, pt1.AttributedSignatures, pt2.AttributedSignatures)
}

func Test_DB_PendingTransmissions(t *testing.T) {
	gormDB, sqlDB := setupDB(t)

	ethKeyStore := cltest.NewKeyStore(t, gormDB).Eth()
	key, _ := cltest.MustInsertRandomKey(t, ethKeyStore)

	spec := MustInsertOffchainreportingOracleSpec(t, gormDB, key.Address)
	spec2 := MustInsertOffchainreportingOracleSpec(t, gormDB, key.Address)
	db := offchainreporting.NewDB(sqlDB, spec.ID)
	db2 := offchainreporting.NewDB(sqlDB, spec2.ID)
	configDigest := MakeConfigDigest(t)

	k := ocrtypes.ReportTimestamp{
		ConfigDigest: configDigest,
		Epoch:        0,
		Round:        1,
	}
	k2 := ocrtypes.ReportTimestamp{
		ConfigDigest: configDigest,
		Epoch:        1,
		Round:        2,
	}

	t.Run("stores and retrieves pending transmissions", func(t *testing.T) {
		p := ocrtypes.PendingTransmission{
			Time:      time.Now(),
			ExtraHash: cltest.Random32Byte(),
			Report:    []byte{0, 2, 3},
			AttributedSignatures: []ocrtypes.AttributedOnChainSignature{
				{Signature: cltest.RandomBytes(7), Signer: 2487},
				{Signature: cltest.RandomBytes(17), Signer: 312},
			},
		}

		err := db.StorePendingTransmission(ctx, k, p)
		require.NoError(t, err)
		m, err := db.PendingTransmissionsWithConfigDigest(ctx, configDigest)
		require.NoError(t, err)
		assertPendingTransmissionEqual(t, p, m[k])

		// Now overwrite value for k to prove that updating works
		p = ocrtypes.PendingTransmission{
			Time:      time.Now(),
			ExtraHash: cltest.Random32Byte(),
			Report:    []byte{1, 2, 3},
			AttributedSignatures: []ocrtypes.AttributedOnChainSignature{
				{Signature: cltest.RandomBytes(7), Signer: 2487},
			},
		}
		err = db.StorePendingTransmission(ctx, k, p)
		require.NoError(t, err)
		m, err = db.PendingTransmissionsWithConfigDigest(ctx, configDigest)
		require.NoError(t, err)
		assertPendingTransmissionEqual(t, p, m[k])

		p2 := ocrtypes.PendingTransmission{
			Time:      time.Now(),
			ExtraHash: cltest.Random32Byte(),
			Report:    []byte{2, 2, 3},
			AttributedSignatures: []ocrtypes.AttributedOnChainSignature{
				{Signature: cltest.RandomBytes(7), Signer: 2487},
			},
		}

		err = db.StorePendingTransmission(ctx, k2, p2)
		require.NoError(t, err)

		kRedHerring := ocrtypes.ReportTimestamp{
			ConfigDigest: ocrtypes.ConfigDigest{43},
			Epoch:        1,
			Round:        2,
		}
		pRedHerring := ocrtypes.PendingTransmission{
			Time:      time.Now(),
			ExtraHash: cltest.Random32Byte(),
			Report:    []byte{3, 2, 3},
			AttributedSignatures: []ocrtypes.AttributedOnChainSignature{
				{Signature: cltest.RandomBytes(7), Signer: 2487},
			},
		}

		err = db.StorePendingTransmission(ctx, kRedHerring, pRedHerring)
		require.NoError(t, err)

		m, err = db.PendingTransmissionsWithConfigDigest(ctx, configDigest)
		require.NoError(t, err)

		require.Len(t, m, 2)

		// HACK to get around time equality because otherwise its annoying (time storage into postgres is mildly lossy)
		require.Equal(t, p.Time.Unix(), m[k].Time.Unix())
		require.Equal(t, p2.Time.Unix(), m[k2].Time.Unix())

		var zt time.Time
		p.Time, p2.Time = zt, zt
		for k, v := range m {
			v.Time = zt
			m[k] = v
		}

		require.Equal(t, p, m[k])
		require.Equal(t, p2, m[k2])

		// No keys for this oracleSpecID yet
		m, err = db2.PendingTransmissionsWithConfigDigest(ctx, configDigest)
		require.NoError(t, err)
		require.Len(t, m, 0)
	})

	t.Run("deletes pending transmission by key", func(t *testing.T) {
		p := ocrtypes.PendingTransmission{
			Time:      time.Unix(100, 0),
			ExtraHash: cltest.Random32Byte(),
			Report:    []byte{1, 4, 3},
			AttributedSignatures: []ocrtypes.AttributedOnChainSignature{
				{Signature: cltest.RandomBytes(7), Signer: 2487},
			},
		}
		err := db.StorePendingTransmission(ctx, k, p)
		require.NoError(t, err)
		err = db2.StorePendingTransmission(ctx, k, p)
		require.NoError(t, err)

		err = db.DeletePendingTransmission(ctx, k)
		require.NoError(t, err)

		m, err := db.PendingTransmissionsWithConfigDigest(ctx, configDigest)
		require.NoError(t, err)
		require.Len(t, m, 1)

		// Did not affect other oracleSpecID
		m, err = db2.PendingTransmissionsWithConfigDigest(ctx, configDigest)
		require.NoError(t, err)
		require.Len(t, m, 1)
	})

	t.Run("allows multiple duplicate keys for different spec ID", func(t *testing.T) {
		p := ocrtypes.PendingTransmission{
			Time: time.Unix(100, 0),
		}
		err := db.StorePendingTransmission(ctx, k2, p)
		require.NoError(t, err)

		m, err := db.PendingTransmissionsWithConfigDigest(ctx, configDigest)
		require.NoError(t, err)
		require.Len(t, m, 1)
		// FIXME: don't understand how the median is being used as a key or what the replacement is yet
		// require.Equal(t, p.Median, m[k2].Median)
	})

	t.Run("deletes pending transmission older than", func(t *testing.T) {
		p := ocrtypes.PendingTransmission{
			Time: time.Unix(100, 0),
		}

		err := db.StorePendingTransmission(ctx, k, p)
		require.NoError(t, err)

		p2 := ocrtypes.PendingTransmission{
			Time: time.Unix(1000, 0),
		}
		err = db.StorePendingTransmission(ctx, k2, p2)
		require.NoError(t, err)

		p2 = ocrtypes.PendingTransmission{
			Time: time.Now(),
		}

		err = db.StorePendingTransmission(ctx, k2, p2)
		require.NoError(t, err)

		err = db.DeletePendingTransmissionsOlderThan(ctx, time.Unix(900, 0))
		require.NoError(t, err)

		m, err := db.PendingTransmissionsWithConfigDigest(ctx, configDigest)
		require.NoError(t, err)
		require.Len(t, m, 1)

		// Didn't affect other oracleSpecIDs
		db = offchainreporting.NewDB(sqlDB, spec2.ID)
		m, err = db.PendingTransmissionsWithConfigDigest(ctx, configDigest)
		require.NoError(t, err)
		require.Len(t, m, 1)
	})
}

func Test_DB_LatestRoundRequested(t *testing.T) {
	gormDB, sqlDB := setupDB(t)

	require.NoError(t, gormDB.Exec(`SET CONSTRAINTS offchainreporting2_latest_roun_offchainreporting2_oracle_spe_fkey DEFERRED`).Error)

	db := offchainreporting.NewDB(sqlDB, 1)
	db2 := offchainreporting.NewDB(sqlDB, 2)

	rawLog := cltest.LogFromFixture(t, "../../testdata/jsonrpc/round_requested_log_1_1.json")

	rr := ocr2aggregator.OCR2AggregatorRoundRequested{
		Requester:    cltest.NewAddress(),
		ConfigDigest: MakeConfigDigest(t),
		Epoch:        42,
		Round:        9,
		Raw:          rawLog,
	}

	t.Run("saves latest round requested", func(t *testing.T) {
		err := postgres.GormTransactionWithDefaultContext(gormDB, func(tx *gorm.DB) error {
			return db.SaveLatestRoundRequested(postgres.MustSQLTx(tx), rr)
		})
		require.NoError(t, err)

		rawLog.Index = 42

		// Now overwrite to prove that updating works
		rr = ocr2aggregator.OCR2AggregatorRoundRequested{
			Requester:    cltest.NewAddress(),
			ConfigDigest: MakeConfigDigest(t),
			Epoch:        43,
			Round:        8,
			Raw:          rawLog,
		}

		err = postgres.GormTransactionWithDefaultContext(gormDB, func(tx *gorm.DB) error {
			return db.SaveLatestRoundRequested(postgres.MustSQLTx(tx), rr)
		})
		require.NoError(t, err)
	})

	t.Run("loads latest round requested", func(t *testing.T) {
		// There is no round for db2
		lrr, err := db2.LoadLatestRoundRequested()
		require.NoError(t, err)
		require.Equal(t, 0, int(lrr.Epoch))

		lrr, err = db.LoadLatestRoundRequested()
		require.NoError(t, err)

		assert.Equal(t, rr, lrr)
	})

	t.Run("spec with latest round requested can be deleted", func(t *testing.T) {
		_, err := sqlDB.Exec(`DELETE FROM offchainreporting_oracle_specs`)
		assert.NoError(t, err)
	})
}
