package offchainreporting_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

func Test_DB_ReadWriteState(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	sqldb, _ := store.DB.DB()
	configDigest := cltest.MakeConfigDigest(t)
	key := cltest.MustInsertRandomKey(t, store.DB)
	spec := cltest.MustInsertOffchainreportingOracleSpec(t, store, key.Address)

	t.Run("reads and writes state", func(t *testing.T) {
		db := offchainreporting.NewDB(sqldb, spec.ID)
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
		db := offchainreporting.NewDB(sqldb, spec.ID)
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
		db := offchainreporting.NewDB(sqldb, spec.ID)
		state := ocrtypes.PersistentState{
			Epoch:                3,
			HighestSentEpoch:     4,
			HighestReceivedEpoch: []uint32{5, 6},
		}

		err := db.WriteState(ctx, configDigest, state)
		require.NoError(t, err)

		// db with different spec
		db = offchainreporting.NewDB(sqldb, -1)

		readState, err := db.ReadState(ctx, configDigest)
		require.NoError(t, err)

		require.Nil(t, readState)
	})

	t.Run("does not return result for wrong config digest", func(t *testing.T) {
		db := offchainreporting.NewDB(sqldb, spec.ID)
		state := ocrtypes.PersistentState{
			Epoch:                4,
			HighestSentEpoch:     5,
			HighestReceivedEpoch: []uint32{6, 7},
		}

		err := db.WriteState(ctx, configDigest, state)
		require.NoError(t, err)

		readState, err := db.ReadState(ctx, cltest.MakeConfigDigest(t))
		require.NoError(t, err)

		require.Nil(t, readState)
	})
}

func Test_DB_ReadWriteConfig(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	sqldb, _ := store.DB.DB()
	config := ocrtypes.ContractConfig{
		ConfigDigest:         cltest.MakeConfigDigest(t),
		Signers:              []common.Address{cltest.NewAddress()},
		Transmitters:         []common.Address{cltest.NewAddress()},
		Threshold:            uint8(35),
		EncodedConfigVersion: uint64(987654),
		Encoded:              []byte{1, 2, 3, 4, 5},
	}
	key := cltest.MustInsertRandomKey(t, store.DB)
	spec := cltest.MustInsertOffchainreportingOracleSpec(t, store, key.Address)
	transmitterAddress := key.Address.Address()

	t.Run("reads and writes config", func(t *testing.T) {
		db := offchainreporting.NewDB(sqldb, spec.ID)

		err := db.WriteConfig(ctx, config)
		require.NoError(t, err)

		readConfig, err := db.ReadConfig(ctx)
		require.NoError(t, err)

		require.Equal(t, &config, readConfig)
	})

	t.Run("updates config", func(t *testing.T) {
		db := offchainreporting.NewDB(sqldb, spec.ID)

		newConfig := ocrtypes.ContractConfig{
			ConfigDigest:         cltest.MakeConfigDigest(t),
			Signers:              []common.Address{utils.ZeroAddress, transmitterAddress, cltest.NewAddress()},
			Transmitters:         []common.Address{utils.ZeroAddress, transmitterAddress, cltest.NewAddress()},
			Threshold:            uint8(36),
			EncodedConfigVersion: uint64(987655),
			Encoded:              []byte{2, 3, 4, 5, 6},
		}

		err := db.WriteConfig(ctx, newConfig)
		require.NoError(t, err)

		readConfig, err := db.ReadConfig(ctx)
		require.NoError(t, err)

		require.Equal(t, &newConfig, readConfig)
	})

	t.Run("does not return result for wrong spec", func(t *testing.T) {
		db := offchainreporting.NewDB(sqldb, spec.ID)

		err := db.WriteConfig(ctx, config)
		require.NoError(t, err)

		db = offchainreporting.NewDB(sqldb, -1)

		readConfig, err := db.ReadConfig(ctx)
		require.NoError(t, err)

		require.Nil(t, readConfig)
	})
}

func Test_DB_PendingTransmissions(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	sqldb, _ := store.DB.DB()
	key := cltest.MustInsertRandomKey(t, store.DB)

	spec := cltest.MustInsertOffchainreportingOracleSpec(t, store, key.Address)
	spec2 := cltest.MustInsertOffchainreportingOracleSpec(t, store, key.Address)
	db := offchainreporting.NewDB(sqldb, spec.ID)
	db2 := offchainreporting.NewDB(sqldb, spec2.ID)
	configDigest := cltest.MakeConfigDigest(t)

	k := ocrtypes.PendingTransmissionKey{
		ConfigDigest: configDigest,
		Epoch:        0,
		Round:        1,
	}
	k2 := ocrtypes.PendingTransmissionKey{
		ConfigDigest: configDigest,
		Epoch:        1,
		Round:        2,
	}

	t.Run("stores and retrieves pending transmissions", func(t *testing.T) {
		p := ocrtypes.PendingTransmission{
			Time:             time.Now(),
			Median:           ocrtypes.Observation(big.NewInt(41)),
			SerializedReport: []byte{0, 2, 3},
			Rs:               [][32]byte{cltest.Random32Byte()},
			Ss:               [][32]byte{cltest.Random32Byte()},
			Vs:               cltest.Random32Byte(),
		}

		err := db.StorePendingTransmission(ctx, k, p)
		require.NoError(t, err)

		// Now overwrite value for k to prove that updating works
		p = ocrtypes.PendingTransmission{
			Time:             time.Now(),
			Median:           ocrtypes.Observation(big.NewInt(42)),
			SerializedReport: []byte{1, 2, 3},
			Rs:               [][32]byte{cltest.Random32Byte()},
			Ss:               [][32]byte{cltest.Random32Byte()},
			Vs:               cltest.Random32Byte(),
		}
		err = db.StorePendingTransmission(ctx, k, p)
		require.NoError(t, err)

		p2 := ocrtypes.PendingTransmission{
			Time:             time.Now(),
			Median:           ocrtypes.Observation(big.NewInt(43)),
			SerializedReport: []byte{2, 2, 3},
			Rs:               [][32]byte{cltest.Random32Byte()},
			Ss:               [][32]byte{cltest.Random32Byte()},
			Vs:               cltest.Random32Byte(),
		}

		err = db.StorePendingTransmission(ctx, k2, p2)
		require.NoError(t, err)

		kRedHerring := ocrtypes.PendingTransmissionKey{
			ConfigDigest: ocrtypes.ConfigDigest{43},
			Epoch:        1,
			Round:        2,
		}
		pRedHerring := ocrtypes.PendingTransmission{
			Time:             time.Now(),
			Median:           ocrtypes.Observation(big.NewInt(43)),
			SerializedReport: []byte{3, 2, 3},
			Rs:               [][32]byte{cltest.Random32Byte()},
			Ss:               [][32]byte{cltest.Random32Byte()},
			Vs:               cltest.Random32Byte(),
		}

		err = db.StorePendingTransmission(ctx, kRedHerring, pRedHerring)
		require.NoError(t, err)

		m, err := db.PendingTransmissionsWithConfigDigest(ctx, configDigest)
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
			Time:             time.Unix(100, 0),
			Median:           ocrtypes.Observation(big.NewInt(44)),
			SerializedReport: []byte{1, 4, 3},
			Rs:               [][32]byte{cltest.Random32Byte()},
			Ss:               [][32]byte{cltest.Random32Byte()},
			Vs:               cltest.Random32Byte(),
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
			Time:             time.Unix(100, 0),
			Median:           ocrtypes.Observation(big.NewInt(44)),
			SerializedReport: []byte{1, 4, 3},
			Rs:               [][32]byte{cltest.Random32Byte()},
			Ss:               [][32]byte{cltest.Random32Byte()},
			Vs:               cltest.Random32Byte(),
		}
		err := db.StorePendingTransmission(ctx, k2, p)
		require.NoError(t, err)

		m, err := db.PendingTransmissionsWithConfigDigest(ctx, configDigest)
		require.NoError(t, err)
		require.Len(t, m, 1)
		require.Equal(t, p.Median, m[k2].Median)
	})

	t.Run("deletes pending transmission older than", func(t *testing.T) {
		p := ocrtypes.PendingTransmission{
			Time:             time.Unix(100, 0),
			Median:           ocrtypes.Observation(big.NewInt(41)),
			SerializedReport: []byte{0, 2, 3},
			Rs:               [][32]byte{cltest.Random32Byte()},
			Ss:               [][32]byte{cltest.Random32Byte()},
			Vs:               cltest.Random32Byte(),
		}

		err := db.StorePendingTransmission(ctx, k, p)
		require.NoError(t, err)

		p2 := ocrtypes.PendingTransmission{
			Time:             time.Unix(1000, 0),
			Median:           ocrtypes.Observation(big.NewInt(42)),
			SerializedReport: []byte{1, 2, 3},
			Rs:               [][32]byte{cltest.Random32Byte()},
			Ss:               [][32]byte{cltest.Random32Byte()},
			Vs:               cltest.Random32Byte(),
		}
		err = db.StorePendingTransmission(ctx, k2, p2)
		require.NoError(t, err)

		p2 = ocrtypes.PendingTransmission{
			Time:             time.Now(),
			Median:           ocrtypes.Observation(big.NewInt(43)),
			SerializedReport: []byte{2, 2, 3},
			Rs:               [][32]byte{cltest.Random32Byte()},
			Ss:               [][32]byte{cltest.Random32Byte()},
			Vs:               cltest.Random32Byte(),
		}

		err = db.StorePendingTransmission(ctx, k2, p2)
		require.NoError(t, err)

		err = db.DeletePendingTransmissionsOlderThan(ctx, time.Unix(900, 0))
		require.NoError(t, err)

		m, err := db.PendingTransmissionsWithConfigDigest(ctx, configDigest)
		require.NoError(t, err)
		require.Len(t, m, 1)

		// Didn't affect other oracleSpecIDs
		db = offchainreporting.NewDB(sqldb, spec2.ID)
		m, err = db.PendingTransmissionsWithConfigDigest(ctx, configDigest)
		require.NoError(t, err)
		require.Len(t, m, 1)
	})
}

func Test_DB_LatestRoundRequested(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	require.NoError(t, store.DB.Exec(`SET CONSTRAINTS offchainreporting_latest_roun_offchainreporting_oracle_spe_fkey DEFERRED`).Error)
	sqldb, _ := store.DB.DB()

	db := offchainreporting.NewDB(sqldb, 1)
	db2 := offchainreporting.NewDB(sqldb, 2)

	rawLog := cltest.LogFromFixture(t, "./testdata/round_requested_log_1_1.json")

	rr := offchainaggregator.OffchainAggregatorRoundRequested{
		Requester:    cltest.NewAddress(),
		ConfigDigest: cltest.MakeConfigDigest(t),
		Epoch:        42,
		Round:        9,
		Raw:          rawLog,
	}

	t.Run("saves latest round requested", func(t *testing.T) {
		err := db.SaveLatestRoundRequested(rr)
		require.NoError(t, err)

		rawLog.Index = 42

		// Now overwrite to prove that updating works
		rr = offchainaggregator.OffchainAggregatorRoundRequested{
			Requester:    cltest.NewAddress(),
			ConfigDigest: cltest.MakeConfigDigest(t),
			Epoch:        43,
			Round:        8,
			Raw:          rawLog,
		}

		err = db.SaveLatestRoundRequested(rr)
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
}
