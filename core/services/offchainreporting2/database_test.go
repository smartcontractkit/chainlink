package offchainreporting2_test

import (
	"context"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting2/testhelpers"

	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	offchainreporting "github.com/smartcontractkit/chainlink/core/services/offchainreporting2"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

func MustInsertOffchainreportingOracleSpec(t *testing.T, db *sqlx.DB, transmitterAddress ethkey.EIP55Address) job.OffchainReporting2OracleSpec {
	t.Helper()

	spec := job.OffchainReporting2OracleSpec{}
	mockJuelsPerFeeCoinSource := `ds1          [type=bridge name=voter_turnout];
	ds1_parse    [type=jsonparse path="one,two"];
	ds1_multiply [type=multiply times=1.23];
	ds1 -> ds1_parse -> ds1_multiply -> answer1;
	answer1      [type=median index=0];`
	require.NoError(t, db.Get(&spec, `INSERT INTO offchainreporting2_oracle_specs (created_at, updated_at, relay, relay_config, contract_id, p2p_bootstrap_peers, is_bootstrap_peer, ocr_key_bundle_id, monitoring_endpoint, transmitter_id, blockchain_timeout, contract_config_tracker_subscribe_interval, contract_config_tracker_poll_interval, contract_config_confirmations, juels_per_fee_coin_pipeline) VALUES (
NOW(),NOW(), 'ethereum', '{}', $1,'{}',false,$2,$3,$4,0,0,0,0,$5
) RETURNING *`, cltest.NewEIP55Address().String(), cltest.DefaultOCR2KeyBundleID, "chain.link:1234", transmitterAddress.String(), mockJuelsPerFeeCoinSource))
	return spec
}

func setupDB(t *testing.T) *sqlx.DB {
	t.Helper()

	sqlx := pgtest.NewSqlxDB(t)

	return sqlx
}

func Test_DB_ReadWriteState(t *testing.T) {
	sqlDB := setupDB(t)

	configDigest := testhelpers.MakeConfigDigest(t)
	cfg := configtest.NewTestGeneralConfig(t)
	ethKeyStore := cltest.NewKeyStore(t, sqlDB, cfg).Eth()
	key, _ := cltest.MustInsertRandomKey(t, ethKeyStore)
	spec := MustInsertOffchainreportingOracleSpec(t, sqlDB, key.Address)
	lggr := logger.TestLogger(t)

	t.Run("reads and writes state", func(t *testing.T) {
		db := offchainreporting.NewDB(sqlDB.DB, spec.ID, lggr)
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
		db := offchainreporting.NewDB(sqlDB.DB, spec.ID, lggr)
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
		db := offchainreporting.NewDB(sqlDB.DB, spec.ID, lggr)
		state := ocrtypes.PersistentState{
			Epoch:                3,
			HighestSentEpoch:     4,
			HighestReceivedEpoch: []uint32{5, 6},
		}

		err := db.WriteState(ctx, configDigest, state)
		require.NoError(t, err)

		// odb with different spec
		db = offchainreporting.NewDB(sqlDB.DB, -1, lggr)

		readState, err := db.ReadState(ctx, configDigest)
		require.NoError(t, err)

		require.Nil(t, readState)
	})

	t.Run("does not return result for wrong config digest", func(t *testing.T) {
		db := offchainreporting.NewDB(sqlDB.DB, spec.ID, lggr)
		state := ocrtypes.PersistentState{
			Epoch:                4,
			HighestSentEpoch:     5,
			HighestReceivedEpoch: []uint32{6, 7},
		}

		err := db.WriteState(ctx, configDigest, state)
		require.NoError(t, err)

		readState, err := db.ReadState(ctx, testhelpers.MakeConfigDigest(t))
		require.NoError(t, err)

		require.Nil(t, readState)
	})
}

func Test_DB_ReadWriteConfig(t *testing.T) {
	sqlDB := setupDB(t)

	config := ocrtypes.ContractConfig{
		ConfigDigest:          testhelpers.MakeConfigDigest(t),
		ConfigCount:           1,
		Signers:               []ocrtypes.OnchainPublicKey{},
		Transmitters:          []ocrtypes.Account{"account1"},
		F:                     79,
		OnchainConfig:         []byte{},
		OffchainConfigVersion: 111,
		OffchainConfig:        []byte{},
	}
	cfg := configtest.NewTestGeneralConfig(t)
	ethKeyStore := cltest.NewKeyStore(t, sqlDB, cfg).Eth()
	key, _ := cltest.MustInsertRandomKey(t, ethKeyStore)
	spec := MustInsertOffchainreportingOracleSpec(t, sqlDB, key.Address)
	lggr := logger.TestLogger(t)

	t.Run("reads and writes config", func(t *testing.T) {
		db := offchainreporting.NewDB(sqlDB.DB, spec.ID, lggr)

		err := db.WriteConfig(ctx, config)
		require.NoError(t, err)

		readConfig, err := db.ReadConfig(ctx)
		require.NoError(t, err)

		require.Equal(t, &config, readConfig)
	})

	t.Run("updates config", func(t *testing.T) {
		db := offchainreporting.NewDB(sqlDB.DB, spec.ID, lggr)

		newConfig := ocrtypes.ContractConfig{
			ConfigDigest: testhelpers.MakeConfigDigest(t),
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
		db := offchainreporting.NewDB(sqlDB.DB, spec.ID, lggr)

		err := db.WriteConfig(ctx, config)
		require.NoError(t, err)

		db = offchainreporting.NewDB(sqlDB.DB, -1, lggr)

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
	sqlDB := setupDB(t)

	cfg := configtest.NewTestGeneralConfig(t)
	ethKeyStore := cltest.NewKeyStore(t, sqlDB, cfg).Eth()
	key, _ := cltest.MustInsertRandomKey(t, ethKeyStore)

	lggr := logger.TestLogger(t)
	spec := MustInsertOffchainreportingOracleSpec(t, sqlDB, key.Address)
	spec2 := MustInsertOffchainreportingOracleSpec(t, sqlDB, key.Address)
	db := offchainreporting.NewDB(sqlDB.DB, spec.ID, lggr)
	db2 := offchainreporting.NewDB(sqlDB.DB, spec2.ID, lggr)
	configDigest := testhelpers.MakeConfigDigest(t)

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
			AttributedSignatures: []ocrtypes.AttributedOnchainSignature{
				{Signature: cltest.MustRandomBytes(t, 7), Signer: 248},
				{Signature: cltest.MustRandomBytes(t, 17), Signer: 31},
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
			AttributedSignatures: []ocrtypes.AttributedOnchainSignature{
				{Signature: cltest.MustRandomBytes(t, 7), Signer: 248},
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
			AttributedSignatures: []ocrtypes.AttributedOnchainSignature{
				{Signature: cltest.MustRandomBytes(t, 7), Signer: 248},
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
			AttributedSignatures: []ocrtypes.AttributedOnchainSignature{
				{Signature: cltest.MustRandomBytes(t, 7), Signer: 248},
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
			AttributedSignatures: []ocrtypes.AttributedOnchainSignature{
				{Signature: cltest.MustRandomBytes(t, 7), Signer: 248},
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
			Time:      time.Unix(100, 0),
			ExtraHash: cltest.Random32Byte(),
			Report:    []byte{2, 2, 3},
			AttributedSignatures: []ocrtypes.AttributedOnchainSignature{
				{Signature: cltest.MustRandomBytes(t, 7), Signer: 248},
			},
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
			Time:      time.Unix(100, 0),
			ExtraHash: cltest.Random32Byte(),
			Report:    []byte{2, 2, 3},
			AttributedSignatures: []ocrtypes.AttributedOnchainSignature{
				{Signature: cltest.MustRandomBytes(t, 7), Signer: 248},
			},
		}

		err := db.StorePendingTransmission(ctx, k, p)
		require.NoError(t, err)

		p2 := ocrtypes.PendingTransmission{
			Time:      time.Unix(1000, 0),
			ExtraHash: cltest.Random32Byte(),
			Report:    []byte{2, 2, 3},
			AttributedSignatures: []ocrtypes.AttributedOnchainSignature{
				{Signature: cltest.MustRandomBytes(t, 7), Signer: 248},
			},
		}
		err = db.StorePendingTransmission(ctx, k2, p2)
		require.NoError(t, err)

		p2 = ocrtypes.PendingTransmission{
			Time:      time.Now(),
			ExtraHash: cltest.Random32Byte(),
			Report:    []byte{2, 2, 3},
			AttributedSignatures: []ocrtypes.AttributedOnchainSignature{
				{Signature: cltest.MustRandomBytes(t, 7), Signer: 248},
			},
		}

		err = db.StorePendingTransmission(ctx, k2, p2)
		require.NoError(t, err)

		err = db.DeletePendingTransmissionsOlderThan(ctx, time.Unix(900, 0))
		require.NoError(t, err)

		m, err := db.PendingTransmissionsWithConfigDigest(ctx, configDigest)
		require.NoError(t, err)
		require.Len(t, m, 1)

		// Didn't affect other oracleSpecIDs
		db = offchainreporting.NewDB(sqlDB.DB, spec2.ID, lggr)
		m, err = db.PendingTransmissionsWithConfigDigest(ctx, configDigest)
		require.NoError(t, err)
		require.Len(t, m, 1)
	})
}
