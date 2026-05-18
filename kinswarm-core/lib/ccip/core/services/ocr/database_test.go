package ocr_test

import (
	"bytes"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr"
)

func Test_DB_ReadWriteState(t *testing.T) {
	db := pgtest.NewSqlxDB(t)

	configDigest := cltest.MakeConfigDigest(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	key, _ := cltest.MustInsertRandomKey(t, ethKeyStore)
	spec := cltest.MustInsertOffchainreportingOracleSpec(t, db, key.EIP55Address)

	t.Run("reads and writes state", func(t *testing.T) {
		t.Log("creating DB")
		odb := ocr.NewTestDB(t, db, spec.ID)
		state := ocrtypes.PersistentState{
			Epoch:                1,
			HighestSentEpoch:     2,
			HighestReceivedEpoch: []uint32{3},
		}

		err := odb.WriteState(testutils.Context(t), configDigest, state)
		require.NoError(t, err)

		readState, err := odb.ReadState(testutils.Context(t), configDigest)
		require.NoError(t, err)

		require.Equal(t, state, *readState)
	})

	t.Run("updates state", func(t *testing.T) {
		odb := ocr.NewTestDB(t, db, spec.ID)
		newState := ocrtypes.PersistentState{
			Epoch:                2,
			HighestSentEpoch:     3,
			HighestReceivedEpoch: []uint32{4, 5},
		}

		err := odb.WriteState(testutils.Context(t), configDigest, newState)
		require.NoError(t, err)

		readState, err := odb.ReadState(testutils.Context(t), configDigest)
		require.NoError(t, err)

		require.Equal(t, newState, *readState)
	})

	t.Run("does not return result for wrong spec", func(t *testing.T) {
		odb := ocr.NewTestDB(t, db, spec.ID)
		state := ocrtypes.PersistentState{
			Epoch:                3,
			HighestSentEpoch:     4,
			HighestReceivedEpoch: []uint32{5, 6},
		}

		err := odb.WriteState(testutils.Context(t), configDigest, state)
		require.NoError(t, err)

		// db with different spec
		odb = ocr.NewTestDB(t, db, -1)

		readState, err := odb.ReadState(testutils.Context(t), configDigest)
		require.NoError(t, err)

		require.Nil(t, readState)
	})

	t.Run("does not return result for wrong config digest", func(t *testing.T) {
		odb := ocr.NewTestDB(t, db, spec.ID)
		state := ocrtypes.PersistentState{
			Epoch:                4,
			HighestSentEpoch:     5,
			HighestReceivedEpoch: []uint32{6, 7},
		}

		err := odb.WriteState(testutils.Context(t), configDigest, state)
		require.NoError(t, err)

		readState, err := odb.ReadState(testutils.Context(t), cltest.MakeConfigDigest(t))
		require.NoError(t, err)

		require.Nil(t, readState)
	})
}

func Test_DB_ReadWriteConfig(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	sqlDB := db

	config := ocrtypes.ContractConfig{
		ConfigDigest:         cltest.MakeConfigDigest(t),
		Signers:              []common.Address{testutils.NewAddress(), testutils.NewAddress()},
		Transmitters:         []common.Address{testutils.NewAddress(), testutils.NewAddress()},
		Threshold:            uint8(35),
		EncodedConfigVersion: uint64(987654),
		Encoded:              []byte{1, 2, 3, 4, 5},
	}
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	key, _ := cltest.MustInsertRandomKey(t, ethKeyStore)
	spec := cltest.MustInsertOffchainreportingOracleSpec(t, db, key.EIP55Address)
	transmitterAddress := key.Address

	t.Run("reads and writes config", func(t *testing.T) {
		db := ocr.NewTestDB(t, sqlDB, spec.ID)

		err := db.WriteConfig(testutils.Context(t), config)
		require.NoError(t, err)

		readConfig, err := db.ReadConfig(testutils.Context(t))
		require.NoError(t, err)

		require.Equal(t, &config, readConfig)
	})

	t.Run("updates config", func(t *testing.T) {
		db := ocr.NewTestDB(t, sqlDB, spec.ID)

		newConfig := ocrtypes.ContractConfig{
			ConfigDigest:         cltest.MakeConfigDigest(t),
			Signers:              []common.Address{utils.ZeroAddress, transmitterAddress, testutils.NewAddress()},
			Transmitters:         []common.Address{utils.ZeroAddress, transmitterAddress, testutils.NewAddress()},
			Threshold:            uint8(36),
			EncodedConfigVersion: uint64(987655),
			Encoded:              []byte{2, 3, 4, 5, 6},
		}

		err := db.WriteConfig(testutils.Context(t), newConfig)
		require.NoError(t, err)

		readConfig, err := db.ReadConfig(testutils.Context(t))
		require.NoError(t, err)

		require.Equal(t, &newConfig, readConfig)
	})

	t.Run("does not return result for wrong spec", func(t *testing.T) {
		db := ocr.NewTestDB(t, sqlDB, spec.ID)

		err := db.WriteConfig(testutils.Context(t), config)
		require.NoError(t, err)

		db = ocr.NewTestDB(t, sqlDB, -1)

		readConfig, err := db.ReadConfig(testutils.Context(t))
		require.NoError(t, err)

		require.Nil(t, readConfig)
	})
}

func assertPendingTransmissionEqual(t *testing.T, pt1, pt2 ocrtypes.PendingTransmission) {
	t.Helper()

	require.Equal(t, pt1.Rs, pt2.Rs)
	require.Equal(t, pt1.Ss, pt2.Ss)
	assert.True(t, bytes.Equal(pt1.Vs[:], pt2.Vs[:]))
	assert.True(t, bytes.Equal(pt1.SerializedReport[:], pt2.SerializedReport[:]))
	assert.Equal(t, pt1.Median, pt2.Median)
	for i := range pt1.Ss {
		assert.True(t, bytes.Equal(pt1.Ss[i][:], pt2.Ss[i][:]))
	}
	for i := range pt1.Rs {
		assert.True(t, bytes.Equal(pt1.Rs[i][:], pt2.Rs[i][:]))
	}
}

func Test_DB_PendingTransmissions(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	sqlDB := db
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	key, _ := cltest.MustInsertRandomKey(t, ethKeyStore)

	spec := cltest.MustInsertOffchainreportingOracleSpec(t, db, key.EIP55Address)
	spec2 := cltest.MustInsertOffchainreportingOracleSpec(t, db, key.EIP55Address)
	odb := ocr.NewTestDB(t, sqlDB, spec.ID)
	odb2 := ocr.NewTestDB(t, sqlDB, spec2.ID)
	configDigest := cltest.MakeConfigDigest(t)

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
			Time:             time.Now(),
			Median:           ocrtypes.Observation(big.NewInt(41)),
			SerializedReport: []byte{0, 2, 3},
			Rs:               [][32]byte{testutils.Random32Byte(), testutils.Random32Byte()},
			Ss:               [][32]byte{testutils.Random32Byte(), testutils.Random32Byte()},
			Vs:               testutils.Random32Byte(),
		}

		err := odb.StorePendingTransmission(testutils.Context(t), k, p)
		require.NoError(t, err)
		m, err := odb.PendingTransmissionsWithConfigDigest(testutils.Context(t), configDigest)
		require.NoError(t, err)
		assertPendingTransmissionEqual(t, m[k], p)

		// Now overwrite value for k to prove that updating works
		p = ocrtypes.PendingTransmission{
			Time:             time.Now(),
			Median:           ocrtypes.Observation(big.NewInt(42)),
			SerializedReport: []byte{1, 2, 3},
			Rs:               [][32]byte{testutils.Random32Byte()},
			Ss:               [][32]byte{testutils.Random32Byte()},
			Vs:               testutils.Random32Byte(),
		}
		err = odb.StorePendingTransmission(testutils.Context(t), k, p)
		require.NoError(t, err)
		m, err = odb.PendingTransmissionsWithConfigDigest(testutils.Context(t), configDigest)
		require.NoError(t, err)
		assertPendingTransmissionEqual(t, m[k], p)

		p2 := ocrtypes.PendingTransmission{
			Time:             time.Now(),
			Median:           ocrtypes.Observation(big.NewInt(43)),
			SerializedReport: []byte{2, 2, 3},
			Rs:               [][32]byte{testutils.Random32Byte()},
			Ss:               [][32]byte{testutils.Random32Byte()},
			Vs:               testutils.Random32Byte(),
		}

		err = odb.StorePendingTransmission(testutils.Context(t), k2, p2)
		require.NoError(t, err)

		kRedHerring := ocrtypes.ReportTimestamp{
			ConfigDigest: ocrtypes.ConfigDigest{43},
			Epoch:        1,
			Round:        2,
		}
		pRedHerring := ocrtypes.PendingTransmission{
			Time:             time.Now(),
			Median:           ocrtypes.Observation(big.NewInt(43)),
			SerializedReport: []byte{3, 2, 3},
			Rs:               [][32]byte{testutils.Random32Byte()},
			Ss:               [][32]byte{testutils.Random32Byte()},
			Vs:               testutils.Random32Byte(),
		}

		err = odb.StorePendingTransmission(testutils.Context(t), kRedHerring, pRedHerring)
		require.NoError(t, err)

		m, err = odb.PendingTransmissionsWithConfigDigest(testutils.Context(t), configDigest)
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
		m, err = odb2.PendingTransmissionsWithConfigDigest(testutils.Context(t), configDigest)
		require.NoError(t, err)
		require.Len(t, m, 0)
	})

	t.Run("deletes pending transmission by key", func(t *testing.T) {
		p := ocrtypes.PendingTransmission{
			Time:             time.Unix(100, 0),
			Median:           ocrtypes.Observation(big.NewInt(44)),
			SerializedReport: []byte{1, 4, 3},
			Rs:               [][32]byte{testutils.Random32Byte()},
			Ss:               [][32]byte{testutils.Random32Byte()},
			Vs:               testutils.Random32Byte(),
		}
		err := odb.StorePendingTransmission(testutils.Context(t), k, p)
		require.NoError(t, err)
		err = odb2.StorePendingTransmission(testutils.Context(t), k, p)
		require.NoError(t, err)

		err = odb.DeletePendingTransmission(testutils.Context(t), k)
		require.NoError(t, err)

		m, err := odb.PendingTransmissionsWithConfigDigest(testutils.Context(t), configDigest)
		require.NoError(t, err)
		require.Len(t, m, 1)

		// Did not affect other oracleSpecID
		m, err = odb2.PendingTransmissionsWithConfigDigest(testutils.Context(t), configDigest)
		require.NoError(t, err)
		require.Len(t, m, 1)
	})

	t.Run("allows multiple duplicate keys for different spec ID", func(t *testing.T) {
		p := ocrtypes.PendingTransmission{
			Time:             time.Unix(100, 0),
			Median:           ocrtypes.Observation(big.NewInt(44)),
			SerializedReport: []byte{1, 4, 3},
			Rs:               [][32]byte{testutils.Random32Byte()},
			Ss:               [][32]byte{testutils.Random32Byte()},
			Vs:               testutils.Random32Byte(),
		}
		err := odb.StorePendingTransmission(testutils.Context(t), k2, p)
		require.NoError(t, err)

		m, err := odb.PendingTransmissionsWithConfigDigest(testutils.Context(t), configDigest)
		require.NoError(t, err)
		require.Len(t, m, 1)
		require.Equal(t, p.Median, m[k2].Median)
	})

	t.Run("deletes pending transmission older than", func(t *testing.T) {
		p := ocrtypes.PendingTransmission{
			Time:             time.Unix(100, 0),
			Median:           ocrtypes.Observation(big.NewInt(41)),
			SerializedReport: []byte{0, 2, 3},
			Rs:               [][32]byte{testutils.Random32Byte()},
			Ss:               [][32]byte{testutils.Random32Byte()},
			Vs:               testutils.Random32Byte(),
		}

		err := odb.StorePendingTransmission(testutils.Context(t), k, p)
		require.NoError(t, err)

		p2 := ocrtypes.PendingTransmission{
			Time:             time.Unix(1000, 0),
			Median:           ocrtypes.Observation(big.NewInt(42)),
			SerializedReport: []byte{1, 2, 3},
			Rs:               [][32]byte{testutils.Random32Byte()},
			Ss:               [][32]byte{testutils.Random32Byte()},
			Vs:               testutils.Random32Byte(),
		}
		err = odb.StorePendingTransmission(testutils.Context(t), k2, p2)
		require.NoError(t, err)

		p2 = ocrtypes.PendingTransmission{
			Time:             time.Now(),
			Median:           ocrtypes.Observation(big.NewInt(43)),
			SerializedReport: []byte{2, 2, 3},
			Rs:               [][32]byte{testutils.Random32Byte()},
			Ss:               [][32]byte{testutils.Random32Byte()},
			Vs:               testutils.Random32Byte(),
		}

		err = odb.StorePendingTransmission(testutils.Context(t), k2, p2)
		require.NoError(t, err)

		err = odb.DeletePendingTransmissionsOlderThan(testutils.Context(t), time.Unix(900, 0))
		require.NoError(t, err)

		m, err := odb.PendingTransmissionsWithConfigDigest(testutils.Context(t), configDigest)
		require.NoError(t, err)
		require.Len(t, m, 1)

		// Didn't affect other oracleSpecIDs
		odb = ocr.NewTestDB(t, sqlDB, spec2.ID)
		m, err = odb.PendingTransmissionsWithConfigDigest(testutils.Context(t), configDigest)
		require.NoError(t, err)
		require.Len(t, m, 1)
	})
}

func Test_DB_LatestRoundRequested(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	sqlDB := db

	pgtest.MustExec(t, db, `SET CONSTRAINTS offchainreporting_latest_roun_offchainreporting_oracle_spe_fkey DEFERRED`)

	odb := ocr.NewTestDB(t, sqlDB, 1)
	odb2 := ocr.NewTestDB(t, sqlDB, 2)

	rawLog := cltest.LogFromFixture(t, "../../testdata/jsonrpc/round_requested_log_1_1.json")

	rr := offchainaggregator.OffchainAggregatorRoundRequested{
		Requester:    testutils.NewAddress(),
		ConfigDigest: cltest.MakeConfigDigest(t),
		Epoch:        42,
		Round:        9,
		Raw:          rawLog,
	}

	t.Run("saves latest round requested", func(t *testing.T) {
		ctx := testutils.Context(t)
		err := odb.SaveLatestRoundRequested(ctx, rr)
		require.NoError(t, err)

		rawLog.Index = 42

		// Now overwrite to prove that updating works
		rr = offchainaggregator.OffchainAggregatorRoundRequested{
			Requester:    testutils.NewAddress(),
			ConfigDigest: cltest.MakeConfigDigest(t),
			Epoch:        43,
			Round:        8,
			Raw:          rawLog,
		}

		err = odb.SaveLatestRoundRequested(ctx, rr)
		require.NoError(t, err)
	})

	t.Run("loads latest round requested", func(t *testing.T) {
		ctx := testutils.Context(t)
		// There is no round for db2
		lrr, err := odb2.LoadLatestRoundRequested(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, int(lrr.Epoch))

		lrr, err = odb.LoadLatestRoundRequested(ctx)
		require.NoError(t, err)

		assert.Equal(t, rr, lrr)
	})

	t.Run("spec with latest round requested can be deleted", func(t *testing.T) {
		_, err := sqlDB.Exec(`DELETE FROM ocr_oracle_specs`)
		assert.NoError(t, err)
	})
}
