package ocr2_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	medianconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/median/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jmoiron/sqlx"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/testhelpers"
)

const defaultPluginID = 0

func MustInsertOCROracleSpec(t *testing.T, db *sqlx.DB, transmitterAddress types.EIP55Address) job.OCR2OracleSpec {
	t.Helper()

	spec := job.OCR2OracleSpec{}
	mockJuelsPerFeeCoinSource := `ds1          [type=bridge name=voter_turnout];
	ds1_parse    [type=jsonparse path="one,two"];
	ds1_multiply [type=multiply times=1.23];
	ds1 -> ds1_parse -> ds1_multiply -> answer1;
	answer1      [type=median index=0];`
	mockGasPriceSubunitsSource := `ds1          [type=bridge name=voter_turnout];
	ds1_parse    [type=jsonparse path="one,two"];
	ds1_multiply [type=multiply times=1.23];
	ds1 -> ds1_parse -> ds1_multiply -> answer1;
	answer1      [type=median index=0];`
	config := medianconfig.PluginConfig{JuelsPerFeeCoinPipeline: mockJuelsPerFeeCoinSource, GasPriceSubunitsPipeline: mockGasPriceSubunitsSource}
	jsonConfig, err := json.Marshal(config)
	require.NoError(t, err)

	require.NoError(t, db.Get(&spec, `INSERT INTO ocr2_oracle_specs (
relay, relay_config, contract_id, p2pv2_bootstrappers, ocr_key_bundle_id, monitoring_endpoint, transmitter_id, 
blockchain_timeout, contract_config_tracker_poll_interval, contract_config_confirmations, plugin_type, plugin_config, onchain_signing_strategy, created_at, updated_at) VALUES (
'ethereum', '{}', $1, '{}', $2, $3, $4,
0, 0, 0, 'median', $5, '{}', NOW(), NOW()
) RETURNING *`, cltest.NewEIP55Address().String(), cltest.DefaultOCR2KeyBundleID, "chain.link:1234", transmitterAddress.String(), jsonConfig))
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
	ethKeyStore := cltest.NewKeyStore(t, sqlDB).Eth()
	key, _ := cltest.MustInsertRandomKey(t, ethKeyStore)
	spec := MustInsertOCROracleSpec(t, sqlDB, key.EIP55Address)
	lggr := logger.TestLogger(t)

	t.Run("reads and writes state", func(t *testing.T) {
		db := ocr2.NewDB(sqlDB, spec.ID, defaultPluginID, lggr)
		state := ocrtypes.PersistentState{
			Epoch:                1,
			HighestSentEpoch:     2,
			HighestReceivedEpoch: []uint32{3},
		}

		err := db.WriteState(testutils.Context(t), configDigest, state)
		require.NoError(t, err)

		readState, err := db.ReadState(testutils.Context(t), configDigest)
		require.NoError(t, err)

		require.Equal(t, state, *readState)
	})

	t.Run("updates state", func(t *testing.T) {
		db := ocr2.NewDB(sqlDB, spec.ID, defaultPluginID, lggr)
		newState := ocrtypes.PersistentState{
			Epoch:                2,
			HighestSentEpoch:     3,
			HighestReceivedEpoch: []uint32{4, 5},
		}

		err := db.WriteState(testutils.Context(t), configDigest, newState)
		require.NoError(t, err)

		readState, err := db.ReadState(testutils.Context(t), configDigest)
		require.NoError(t, err)

		require.Equal(t, newState, *readState)
	})

	t.Run("does not return result for wrong spec", func(t *testing.T) {
		db := ocr2.NewDB(sqlDB, spec.ID, defaultPluginID, lggr)
		state := ocrtypes.PersistentState{
			Epoch:                3,
			HighestSentEpoch:     4,
			HighestReceivedEpoch: []uint32{5, 6},
		}

		err := db.WriteState(testutils.Context(t), configDigest, state)
		require.NoError(t, err)

		// odb with different spec
		db = ocr2.NewDB(sqlDB, -1, defaultPluginID, lggr)

		readState, err := db.ReadState(testutils.Context(t), configDigest)
		require.NoError(t, err)

		require.Nil(t, readState)
	})

	t.Run("does not return result for wrong config digest", func(t *testing.T) {
		db := ocr2.NewDB(sqlDB, spec.ID, defaultPluginID, lggr)
		state := ocrtypes.PersistentState{
			Epoch:                4,
			HighestSentEpoch:     5,
			HighestReceivedEpoch: []uint32{6, 7},
		}

		err := db.WriteState(testutils.Context(t), configDigest, state)
		require.NoError(t, err)

		readState, err := db.ReadState(testutils.Context(t), testhelpers.MakeConfigDigest(t))
		require.NoError(t, err)

		require.Nil(t, readState)
	})
}

func Test_DB_ReadWriteConfig(t *testing.T) {
	sqlDB := setupDB(t)

	config := ocrtypes.ContractConfig{
		ConfigDigest:          testhelpers.MakeConfigDigest(t),
		ConfigCount:           1,
		Signers:               []ocrtypes.OnchainPublicKey{{0x01}, {0x02}},
		Transmitters:          []ocrtypes.Account{"account1", "account2"},
		F:                     79,
		OnchainConfig:         []byte{0x01, 0x02},
		OffchainConfigVersion: 111,
		OffchainConfig:        []byte{0x03, 0x04},
	}
	ethKeyStore := cltest.NewKeyStore(t, sqlDB).Eth()
	key, _ := cltest.MustInsertRandomKey(t, ethKeyStore)
	spec := MustInsertOCROracleSpec(t, sqlDB, key.EIP55Address)
	lggr := logger.TestLogger(t)

	t.Run("reads and writes config", func(t *testing.T) {
		db := ocr2.NewDB(sqlDB, spec.ID, defaultPluginID, lggr)

		err := db.WriteConfig(testutils.Context(t), config)
		require.NoError(t, err)

		readConfig, err := db.ReadConfig(testutils.Context(t))
		require.NoError(t, err)

		require.Equal(t, &config, readConfig)
	})

	t.Run("updates config", func(t *testing.T) {
		db := ocr2.NewDB(sqlDB, spec.ID, defaultPluginID, lggr)

		newConfig := ocrtypes.ContractConfig{
			ConfigDigest: testhelpers.MakeConfigDigest(t),
			Signers:      []ocrtypes.OnchainPublicKey{},
			Transmitters: []ocrtypes.Account{},
		}

		err := db.WriteConfig(testutils.Context(t), newConfig)
		require.NoError(t, err)

		readConfig, err := db.ReadConfig(testutils.Context(t))
		require.NoError(t, err)

		require.Equal(t, &newConfig, readConfig)
	})

	t.Run("does not return result for wrong spec", func(t *testing.T) {
		db := ocr2.NewDB(sqlDB, spec.ID, defaultPluginID, lggr)

		err := db.WriteConfig(testutils.Context(t), config)
		require.NoError(t, err)

		db = ocr2.NewDB(sqlDB, -1, defaultPluginID, lggr)

		readConfig, err := db.ReadConfig(testutils.Context(t))
		require.NoError(t, err)

		require.Nil(t, readConfig)
	})

	t.Run("reads and writes config for multiple plugins", func(t *testing.T) {
		otherPluginID := int32(2)
		db1 := ocr2.NewDB(sqlDB, spec.ID, defaultPluginID, lggr)
		db2 := ocr2.NewDB(sqlDB, spec.ID, otherPluginID, lggr)

		otherConfig := ocrtypes.ContractConfig{
			ConfigDigest: testhelpers.MakeConfigDigest(t),
			Signers:      []ocrtypes.OnchainPublicKey{},
			Transmitters: []ocrtypes.Account{},
		}
		err := db1.WriteConfig(testutils.Context(t), config)
		require.NoError(t, err)
		err = db2.WriteConfig(testutils.Context(t), otherConfig)
		require.NoError(t, err)

		readConfig, err := db1.ReadConfig(testutils.Context(t))
		require.NoError(t, err)
		require.Equal(t, &config, readConfig)

		readConfig, err = db2.ReadConfig(testutils.Context(t))
		require.NoError(t, err)
		require.Equal(t, &otherConfig, readConfig)
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

	ethKeyStore := cltest.NewKeyStore(t, sqlDB).Eth()
	key, _ := cltest.MustInsertRandomKey(t, ethKeyStore)

	lggr := logger.TestLogger(t)
	spec := MustInsertOCROracleSpec(t, sqlDB, key.EIP55Address)
	spec2 := MustInsertOCROracleSpec(t, sqlDB, key.EIP55Address)
	db := ocr2.NewDB(sqlDB, spec.ID, defaultPluginID, lggr)
	db2 := ocr2.NewDB(sqlDB, spec2.ID, defaultPluginID, lggr)
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
			ExtraHash: testutils.Random32Byte(),
			Report:    []byte{0, 2, 3},
			AttributedSignatures: []ocrtypes.AttributedOnchainSignature{
				{Signature: cltest.MustRandomBytes(t, 7), Signer: 248},
				{Signature: cltest.MustRandomBytes(t, 17), Signer: 31},
			},
		}

		err := db.StorePendingTransmission(testutils.Context(t), k, p)
		require.NoError(t, err)
		m, err := db.PendingTransmissionsWithConfigDigest(testutils.Context(t), configDigest)
		require.NoError(t, err)
		assertPendingTransmissionEqual(t, p, m[k])

		// Now overwrite value for k to prove that updating works
		p = ocrtypes.PendingTransmission{
			Time:      time.Now(),
			ExtraHash: testutils.Random32Byte(),
			Report:    []byte{1, 2, 3},
			AttributedSignatures: []ocrtypes.AttributedOnchainSignature{
				{Signature: cltest.MustRandomBytes(t, 7), Signer: 248},
			},
		}
		err = db.StorePendingTransmission(testutils.Context(t), k, p)
		require.NoError(t, err)
		m, err = db.PendingTransmissionsWithConfigDigest(testutils.Context(t), configDigest)
		require.NoError(t, err)
		assertPendingTransmissionEqual(t, p, m[k])

		p2 := ocrtypes.PendingTransmission{
			Time:      time.Now(),
			ExtraHash: testutils.Random32Byte(),
			Report:    []byte{2, 2, 3},
			AttributedSignatures: []ocrtypes.AttributedOnchainSignature{
				{Signature: cltest.MustRandomBytes(t, 7), Signer: 248},
			},
		}

		err = db.StorePendingTransmission(testutils.Context(t), k2, p2)
		require.NoError(t, err)

		kRedHerring := ocrtypes.ReportTimestamp{
			ConfigDigest: ocrtypes.ConfigDigest{43},
			Epoch:        1,
			Round:        2,
		}
		pRedHerring := ocrtypes.PendingTransmission{
			Time:      time.Now(),
			ExtraHash: testutils.Random32Byte(),
			Report:    []byte{3, 2, 3},
			AttributedSignatures: []ocrtypes.AttributedOnchainSignature{
				{Signature: cltest.MustRandomBytes(t, 7), Signer: 248},
			},
		}

		err = db.StorePendingTransmission(testutils.Context(t), kRedHerring, pRedHerring)
		require.NoError(t, err)

		m, err = db.PendingTransmissionsWithConfigDigest(testutils.Context(t), configDigest)
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
		m, err = db2.PendingTransmissionsWithConfigDigest(testutils.Context(t), configDigest)
		require.NoError(t, err)
		require.Len(t, m, 0)
	})

	t.Run("deletes pending transmission by key", func(t *testing.T) {
		p := ocrtypes.PendingTransmission{
			Time:      time.Unix(100, 0),
			ExtraHash: testutils.Random32Byte(),
			Report:    []byte{1, 4, 3},
			AttributedSignatures: []ocrtypes.AttributedOnchainSignature{
				{Signature: cltest.MustRandomBytes(t, 7), Signer: 248},
			},
		}
		err := db.StorePendingTransmission(testutils.Context(t), k, p)
		require.NoError(t, err)
		err = db2.StorePendingTransmission(testutils.Context(t), k, p)
		require.NoError(t, err)

		err = db.DeletePendingTransmission(testutils.Context(t), k)
		require.NoError(t, err)

		m, err := db.PendingTransmissionsWithConfigDigest(testutils.Context(t), configDigest)
		require.NoError(t, err)
		require.Len(t, m, 1)

		// Did not affect other oracleSpecID
		m, err = db2.PendingTransmissionsWithConfigDigest(testutils.Context(t), configDigest)
		require.NoError(t, err)
		require.Len(t, m, 1)
	})

	t.Run("allows multiple duplicate keys for different spec ID", func(t *testing.T) {
		p := ocrtypes.PendingTransmission{
			Time:      time.Unix(100, 0),
			ExtraHash: testutils.Random32Byte(),
			Report:    []byte{2, 2, 3},
			AttributedSignatures: []ocrtypes.AttributedOnchainSignature{
				{Signature: cltest.MustRandomBytes(t, 7), Signer: 248},
			},
		}
		err := db.StorePendingTransmission(testutils.Context(t), k2, p)
		require.NoError(t, err)

		m, err := db.PendingTransmissionsWithConfigDigest(testutils.Context(t), configDigest)
		require.NoError(t, err)
		require.Len(t, m, 1)
		// FIXME: don't understand how the median is being used as a key or what the replacement is yet
		// require.Equal(t, p.Median, m[k2].Median)
	})

	t.Run("deletes pending transmission older than", func(t *testing.T) {
		p := ocrtypes.PendingTransmission{
			Time:      time.Unix(100, 0),
			ExtraHash: testutils.Random32Byte(),
			Report:    []byte{2, 2, 3},
			AttributedSignatures: []ocrtypes.AttributedOnchainSignature{
				{Signature: cltest.MustRandomBytes(t, 7), Signer: 248},
			},
		}

		err := db.StorePendingTransmission(testutils.Context(t), k, p)
		require.NoError(t, err)

		p2 := ocrtypes.PendingTransmission{
			Time:      time.Unix(1000, 0),
			ExtraHash: testutils.Random32Byte(),
			Report:    []byte{2, 2, 3},
			AttributedSignatures: []ocrtypes.AttributedOnchainSignature{
				{Signature: cltest.MustRandomBytes(t, 7), Signer: 248},
			},
		}
		err = db.StorePendingTransmission(testutils.Context(t), k2, p2)
		require.NoError(t, err)

		p2 = ocrtypes.PendingTransmission{
			Time:      time.Now(),
			ExtraHash: testutils.Random32Byte(),
			Report:    []byte{2, 2, 3},
			AttributedSignatures: []ocrtypes.AttributedOnchainSignature{
				{Signature: cltest.MustRandomBytes(t, 7), Signer: 248},
			},
		}

		err = db.StorePendingTransmission(testutils.Context(t), k2, p2)
		require.NoError(t, err)

		err = db.DeletePendingTransmissionsOlderThan(testutils.Context(t), time.Unix(900, 0))
		require.NoError(t, err)

		m, err := db.PendingTransmissionsWithConfigDigest(testutils.Context(t), configDigest)
		require.NoError(t, err)
		require.Len(t, m, 1)

		// Didn't affect other oracleSpecIDs
		db = ocr2.NewDB(sqlDB, spec2.ID, defaultPluginID, lggr)
		m, err = db.PendingTransmissionsWithConfigDigest(testutils.Context(t), configDigest)
		require.NoError(t, err)
		require.Len(t, m, 1)
	})
}

func Test_DB_ReadWriteProtocolState(t *testing.T) {
	sqlDB := setupDB(t)

	lggr := logger.TestLogger(t)
	db := ocr2.NewDB(sqlDB, 0, defaultPluginID, lggr)
	cd1 := testhelpers.MakeConfigDigest(t)
	cd2 := testhelpers.MakeConfigDigest(t)
	ctx := testutils.Context(t)

	assertCount := func(expected int64) {
		testutils.AssertCount(t, sqlDB, "ocr_protocol_states", expected)
	}

	t.Run("stores and retrieves protocol state", func(t *testing.T) {
		assertCount(0)

		err := db.WriteProtocolState(ctx, cd1, "key1", []byte{1})
		assert.NoError(t, err)

		assertCount(1)

		err = db.WriteProtocolState(ctx, cd2, "key1", []byte{2})
		assert.NoError(t, err)

		assertCount(2)

		err = db.WriteProtocolState(ctx, cd2, "key2", []byte{3})
		assert.NoError(t, err)

		assertCount(3)

		// should overwrite
		err = db.WriteProtocolState(ctx, cd2, "key2", []byte{4})
		assert.NoError(t, err)

		val, err := db.ReadProtocolState(ctx, cd1, "key1")
		assert.NoError(t, err)
		assert.Equal(t, []byte{1}, val)

		val, err = db.ReadProtocolState(ctx, cd2, "key1")
		assert.NoError(t, err)
		assert.Equal(t, []byte{2}, val)

		val, err = db.ReadProtocolState(ctx, cd2, "key2")
		assert.NoError(t, err)
		assert.Equal(t, []byte{4}, val)

		// should write empty value
		err = db.WriteProtocolState(ctx, cd1, "key1", []byte{})
		assert.NoError(t, err)

		val, err = db.ReadProtocolState(ctx, cd1, "key1")
		assert.NoError(t, err)
		assert.Equal(t, []byte{}, val)

		assertCount(3)

		// should delete value
		err = db.WriteProtocolState(ctx, cd1, "key1", nil)
		assert.NoError(t, err)

		assertCount(2)

		// trying to read missing value yields nil
		val, err = db.ReadProtocolState(ctx, cd1, "key1")
		assert.NoError(t, err)
		assert.Nil(t, val)
	})
}
