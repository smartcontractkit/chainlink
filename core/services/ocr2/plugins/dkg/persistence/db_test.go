package persistence

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	ocr2vrftypes "github.com/smartcontractkit/ocr2vrf/types"
	"github.com/smartcontractkit/ocr2vrf/types/hash"
	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/relay"
)

func setup(t testing.TB) (ocr2vrftypes.DKGSharePersistence, *sqlx.DB) {
	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	return NewShareDB(db, lggr, pgtest.NewQConfig(true), big.NewInt(1337), relay.EVM), db
}

func TestShareDB_WriteShareRecords(t *testing.T) {
	configDigest := testutils.Random32Byte()
	keyID := testutils.Random32Byte()

	t.Run("valid input", func(tt *testing.T) {
		shareDB, db := setup(tt)
		var expectedRecords []ocr2vrftypes.PersistentShareSetRecord

		// Starting from 1 because player indexes must not be 0
		for i := 1; i < 4; i++ {
			b := ocr2vrftypes.RawMarshalPlayerIdxInt(ocr2vrftypes.PlayerIdxInt(i))
			playerIdx, _, err := ocr2vrftypes.UnmarshalPlayerIdx(b)
			require.NoError(t, err)
			shareRecord := crypto.Keccak256Hash([]byte(fmt.Sprintf("%d", i)))
			shareRecordHash := hash.GetHash(shareRecord[:])
			var h hash.Hash
			copy(h[:], shareRecordHash[:])
			rec := ocr2vrftypes.PersistentShareSetRecord{
				Dealer:               *playerIdx,
				MarshaledShareRecord: shareRecord[:],
				Hash:                 h,
			}
			expectedRecords = append(expectedRecords, rec)
		}

		err := shareDB.WriteShareRecords(context.TODO(), configDigest, keyID, expectedRecords)
		require.NoError(tt, err)

		rows, err := db.Query(`SELECT COUNT(*) AS count FROM dkg_shares`)
		require.NoError(tt, err)

		var count int
		for rows.Next() {
			require.NoError(tt, rows.Scan(&count))
		}

		require.Equal(tt, 3, count)
	})

	t.Run("bad input, zero hash", func(tt *testing.T) {
		shareDB, db := setup(tt)
		b := ocr2vrftypes.RawMarshalPlayerIdxInt(ocr2vrftypes.PlayerIdxInt(1))
		dealer, _, err := ocr2vrftypes.UnmarshalPlayerIdx(b)
		require.NoError(tt, err)
		records := []ocr2vrftypes.PersistentShareSetRecord{
			{
				Dealer:               *dealer,
				MarshaledShareRecord: []byte{1},
				Hash:                 hash.Hash{}, // There's a problem here
			},
		}

		// no error, but there will be no rows inserted in the db
		err = shareDB.WriteShareRecords(context.TODO(), configDigest, keyID, records)
		require.NoError(tt, err)

		rows, err := db.Query(`SELECT COUNT(*) AS count FROM dkg_shares`)
		require.NoError(tt, err)

		var count int
		for rows.Next() {
			require.NoError(tt, rows.Scan(&count))
		}

		require.Equal(tt, 0, count)
	})

	t.Run("bad input, nonmatching hash", func(tt *testing.T) {
		shareDB, db := setup(tt)
		var records []ocr2vrftypes.PersistentShareSetRecord

		// Starting from 1 because player indexes must not be 0
		for i := 1; i < 4; i++ {
			b := ocr2vrftypes.RawMarshalPlayerIdxInt(ocr2vrftypes.PlayerIdxInt(i))
			playerIdx, _, err := ocr2vrftypes.UnmarshalPlayerIdx(b)
			require.NoError(t, err)
			shareRecord := crypto.Keccak256Hash([]byte(fmt.Sprintf("%d", i)))
			// Expected hash is SHA256, not Keccak256.
			shareRecordHash := crypto.Keccak256Hash(shareRecord[:])
			var h hash.Hash
			copy(h[:], shareRecordHash[:])
			rec := ocr2vrftypes.PersistentShareSetRecord{
				Dealer:               *playerIdx,
				MarshaledShareRecord: shareRecord[:],
				Hash:                 h,
			}
			records = append(records, rec)
		}

		err := shareDB.WriteShareRecords(context.TODO(), configDigest, keyID, records)
		require.Error(tt, err)

		// no rows should have been inserted
		rows, err := db.Query(`SELECT COUNT(*) AS count FROM dkg_shares`)
		require.NoError(tt, err)

		var count int
		for rows.Next() {
			require.NoError(tt, rows.Scan(&count))
		}

		require.Equal(tt, 0, count)
	})
}

func TestShareDBE2E(t *testing.T) {
	shareDB, _ := setup(t)

	// create some fake data to insert and retrieve
	configDigest := testutils.Random32Byte()
	keyID := testutils.Random32Byte()
	var expectedRecords []ocr2vrftypes.PersistentShareSetRecord
	expectedRecordsMap := make(map[ocr2vrftypes.PlayerIdx]ocr2vrftypes.PersistentShareSetRecord)

	// Starting from 1 because player indexes must not be 0
	for i := 1; i < 4; i++ {
		b := ocr2vrftypes.RawMarshalPlayerIdxInt(ocr2vrftypes.PlayerIdxInt(i))
		playerIdx, _, err := ocr2vrftypes.UnmarshalPlayerIdx(b)
		require.NoError(t, err)
		shareRecord := crypto.Keccak256Hash([]byte(fmt.Sprintf("%d", i)))
		shareRecordHash := hash.GetHash(shareRecord[:])
		var h hash.Hash
		copy(h[:], shareRecordHash[:])
		rec := ocr2vrftypes.PersistentShareSetRecord{
			Dealer:               *playerIdx,
			MarshaledShareRecord: shareRecord[:],
			Hash:                 h,
		}
		expectedRecords = append(expectedRecords, rec)
		expectedRecordsMap[*playerIdx] = rec
	}

	err := shareDB.WriteShareRecords(context.TODO(), configDigest, keyID, expectedRecords)
	require.NoError(t, err)

	actualRecords, err := shareDB.ReadShareRecords(configDigest, keyID)
	require.NoError(t, err)

	assert.Equal(t, len(expectedRecords), len(actualRecords))
	numAssertions := 0
	for _, actualRecord := range actualRecords {
		expectedRecord, ok := expectedRecordsMap[actualRecord.Dealer]
		require.True(t, ok)
		require.Equal(t, expectedRecord.MarshaledShareRecord, actualRecord.MarshaledShareRecord)
		require.Equal(t, expectedRecord.Hash[:], actualRecord.Hash[:])
		numAssertions++
	}

	require.Equal(t, len(expectedRecords), numAssertions)
}
