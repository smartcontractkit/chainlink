package vault

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/smdkg/dkgocr/dkgocrtypes"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

func setupORM(t *testing.T) (*sqlx.DB, dkgocrtypes.ResultPackageDatabase) {
	t.Helper()

	db := pgtest.NewSqlxDB(t)
	orm := NewVaultORM(db)

	return db, orm
}

func createTestResultPackage() dkgocrtypes.ResultPackageDatabaseValue {
	var configDigest types.ConfigDigest
	copy(configDigest[:], common.Hex2Bytes("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"))

	signatures := []types.AttributedOnchainSignature{
		{
			Signature: common.Hex2Bytes("deadbeef"),
			Signer:    commontypes.OracleID(1),
		},
		{
			Signature: common.Hex2Bytes("cafebabe"),
			Signer:    commontypes.OracleID(2),
		},
	}

	return dkgocrtypes.ResultPackageDatabaseValue{
		ConfigDigest:            configDigest,
		SeqNr:                   42,
		ReportWithResultPackage: common.Hex2Bytes("feedface"),
		Signatures:              signatures,
	}
}

func createTestInstanceID() dkgocrtypes.InstanceID {
	configContract := testutils.NewAddress()
	var configDigest types.ConfigDigest
	copy(configDigest[:], common.Hex2Bytes("abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"))
	return dkgocrtypes.MakeInstanceID(configContract, configDigest)
}

func TestORM_WriteAndReadResultPackage(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	_, orm := setupORM(t)
	instanceID := createTestInstanceID()
	value := createTestResultPackage()

	err := orm.WriteResultPackage(ctx, instanceID, value)
	require.NoError(t, err)

	readValue, err := orm.ReadResultPackage(ctx, instanceID)
	require.NoError(t, err)
	require.NotNil(t, readValue)

	assert.Equal(t, value.ConfigDigest, readValue.ConfigDigest)
	assert.Equal(t, value.SeqNr, readValue.SeqNr)
	assert.Equal(t, value.ReportWithResultPackage, readValue.ReportWithResultPackage)
	require.Len(t, readValue.Signatures, len(value.Signatures))

	for i, expectedSig := range value.Signatures {
		assert.Equal(t, expectedSig.Signature, readValue.Signatures[i].Signature)
		assert.Equal(t, expectedSig.Signer, readValue.Signatures[i].Signer)
	}
}

func TestORM_ReadResultPackage_NotFound(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	_, orm := setupORM(t)
	instanceID := createTestInstanceID()

	readValue, err := orm.ReadResultPackage(ctx, instanceID)
	require.NoError(t, err)
	assert.Nil(t, readValue)
}

func TestORM_WriteResultPackage_Upsert(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	db, orm := setupORM(t)
	instanceID := createTestInstanceID()
	value := createTestResultPackage()

	err := orm.WriteResultPackage(ctx, instanceID, value)
	require.NoError(t, err)

	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM dkg_results WHERE instance_id = $1", instanceID)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	updatedValue := value
	updatedValue.SeqNr = 99
	updatedValue.ReportWithResultPackage = common.Hex2Bytes("deadfeedface")
	updatedValue.Signatures = []types.AttributedOnchainSignature{
		{
			Signature: common.Hex2Bytes("beefcafe"),
			Signer:    commontypes.OracleID(3),
		},
	}

	err = orm.WriteResultPackage(ctx, instanceID, updatedValue)
	require.NoError(t, err)

	err = db.Get(&count, "SELECT COUNT(*) FROM dkg_results WHERE instance_id = $1", instanceID)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	readValue, err := orm.ReadResultPackage(ctx, instanceID)
	require.NoError(t, err)
	require.NotNil(t, readValue)

	assert.Equal(t, uint64(99), readValue.SeqNr)
	assert.Equal(t, common.Hex2Bytes("deadfeedface"), readValue.ReportWithResultPackage)
	require.Len(t, readValue.Signatures, 1)
	assert.Equal(t, common.Hex2Bytes("beefcafe"), readValue.Signatures[0].Signature)
	assert.Equal(t, commontypes.OracleID(3), readValue.Signatures[0].Signer)
}

func TestORM_WriteResultPackage_MultipleSignatures(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	_, orm := setupORM(t)
	instanceID := createTestInstanceID()
	value := createTestResultPackage()

	value.Signatures = []types.AttributedOnchainSignature{
		{Signature: common.Hex2Bytes("sig1"), Signer: commontypes.OracleID(1)},
		{Signature: common.Hex2Bytes("sig2"), Signer: commontypes.OracleID(2)},
		{Signature: common.Hex2Bytes("sig3"), Signer: commontypes.OracleID(3)},
		{Signature: common.Hex2Bytes("sig4"), Signer: commontypes.OracleID(4)},
		{Signature: common.Hex2Bytes("sig5"), Signer: commontypes.OracleID(5)},
	}

	err := orm.WriteResultPackage(ctx, instanceID, value)
	require.NoError(t, err)

	readValue, err := orm.ReadResultPackage(ctx, instanceID)
	require.NoError(t, err)
	require.NotNil(t, readValue)

	require.Len(t, readValue.Signatures, 5)
	for i, expectedSig := range value.Signatures {
		assert.Equal(t, expectedSig.Signature, readValue.Signatures[i].Signature)
		assert.Equal(t, expectedSig.Signer, readValue.Signatures[i].Signer)
	}
}

func TestORM_WriteResultPackage_DifferentInstanceIDs(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	db, orm := setupORM(t)

	instanceID1 := createTestInstanceID()
	instanceID2 := dkgocrtypes.MakeInstanceID(testutils.NewAddress(), types.ConfigDigest{})

	value1 := createTestResultPackage()
	value1.SeqNr = 1

	value2 := createTestResultPackage()
	value2.SeqNr = 2

	err := orm.WriteResultPackage(ctx, instanceID1, value1)
	require.NoError(t, err)

	err = orm.WriteResultPackage(ctx, instanceID2, value2)
	require.NoError(t, err)

	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM dkg_results")
	require.NoError(t, err)
	assert.Equal(t, 2, count)

	readValue1, err := orm.ReadResultPackage(ctx, instanceID1)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), readValue1.SeqNr)

	readValue2, err := orm.ReadResultPackage(ctx, instanceID2)
	require.NoError(t, err)
	assert.Equal(t, uint64(2), readValue2.SeqNr)
}

func TestORM_WriteResultPackage_TimestampsUpdated(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	db, orm := setupORM(t)
	instanceID := createTestInstanceID()
	value := createTestResultPackage()

	err := orm.WriteResultPackage(ctx, instanceID, value)
	require.NoError(t, err)

	var initialResult struct {
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}
	err = db.Get(&initialResult, "SELECT created_at, updated_at FROM dkg_results WHERE instance_id = $1", instanceID)
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	value.SeqNr = 999
	err = orm.WriteResultPackage(ctx, instanceID, value)
	require.NoError(t, err)

	var finalResult struct {
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}
	err = db.Get(&finalResult, "SELECT created_at, updated_at FROM dkg_results WHERE instance_id = $1", instanceID)
	require.NoError(t, err)

	assert.Equal(t, initialResult.CreatedAt, finalResult.CreatedAt)
	assert.True(t, finalResult.UpdatedAt.Equal(initialResult.UpdatedAt) || finalResult.UpdatedAt.After(initialResult.UpdatedAt),
		"updated_at should be >= initial timestamp")
}

func TestORM_ConfigDigestHandling(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	_, orm := setupORM(t)
	instanceID := createTestInstanceID()
	value := createTestResultPackage()

	expectedBytes := []byte{
		0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0,
		0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0,
		0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0,
		0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0,
	}
	var testDigest types.ConfigDigest
	copy(testDigest[:], expectedBytes)
	value.ConfigDigest = testDigest

	err := orm.WriteResultPackage(ctx, instanceID, value)
	require.NoError(t, err)

	readValue, err := orm.ReadResultPackage(ctx, instanceID)
	require.NoError(t, err)
	require.NotNil(t, readValue)

	assert.Equal(t, testDigest, readValue.ConfigDigest)
}

func TestORM_WriteResultPackage_ValidationErrors(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	_, orm := setupORM(t)
	instanceID := createTestInstanceID()

	t.Run("zero config digest", func(t *testing.T) {
		value := createTestResultPackage()
		value.ConfigDigest = types.ConfigDigest{}

		err := orm.WriteResultPackage(ctx, instanceID, value)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "config digest cannot be zero")
	})

	t.Run("zero sequence number", func(t *testing.T) {
		value := createTestResultPackage()
		value.SeqNr = 0

		err := orm.WriteResultPackage(ctx, instanceID, value)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "sequence number cannot be zero")
	})

	t.Run("empty report", func(t *testing.T) {
		value := createTestResultPackage()
		value.ReportWithResultPackage = []byte{}

		err := orm.WriteResultPackage(ctx, instanceID, value)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "report with result package cannot be empty")
	})

	t.Run("empty signatures", func(t *testing.T) {
		value := createTestResultPackage()
		value.Signatures = []types.AttributedOnchainSignature{}

		err := orm.WriteResultPackage(ctx, instanceID, value)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "signatures cannot be empty")
	})

	t.Run("valid package passes validation", func(t *testing.T) {
		value := createTestResultPackage()

		err := orm.WriteResultPackage(ctx, instanceID, value)
		require.NoError(t, err)
	})
}
