package persistence

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	ocr2vrftypes "github.com/smartcontractkit/ocr2vrf/types"
	"github.com/smartcontractkit/ocr2vrf/types/hash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
)

func setup(t testing.TB) ocr2vrftypes.DKGSharePersistence {
	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	return NewShareDB(db, lggr, pgtest.NewQConfig(true))
}

func TestShareDB(t *testing.T) {
	shareDB := setup(t)

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
