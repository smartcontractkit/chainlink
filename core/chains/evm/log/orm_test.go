package log_test

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
)

func TestORM_broadcasts(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := cltest.NewTestGeneralConfig(t)
	lggr := logger.TestLogger(t)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	orm := log.NewORM(db, lggr, cfg, cltest.FixtureChainID)

	_, addr := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)
	specV2 := cltest.MustInsertV2JobSpec(t, db, addr)

	const selectQuery = `SELECT consumed FROM log_broadcasts
		WHERE block_hash = $1 AND block_number = $2 AND log_index = $3 AND job_id = $4 AND evm_chain_id = $5`

	listener := &mockListener{specV2.ID}

	rawLog := cltest.RandomLog(t)
	queryArgs := []interface{}{rawLog.BlockHash, rawLog.BlockNumber, rawLog.Index, listener.JobID(), cltest.FixtureChainID.String()}

	// No rows
	res, err := db.Exec(selectQuery, queryArgs...)
	require.NoError(t, err)
	rowsAffected, err := res.RowsAffected()
	require.NoError(t, err)
	require.Zero(t, rowsAffected)

	t.Run("WasBroadcastConsumed_DNE", func(t *testing.T) {
		_, err := orm.WasBroadcastConsumed(rawLog.BlockHash, rawLog.Index, listener.JobID())
		require.NoError(t, err)
	})

	require.True(t, t.Run("CreateBroadcast", func(t *testing.T) {
		err := orm.CreateBroadcast(rawLog.BlockHash, rawLog.BlockNumber, rawLog.Index, listener.JobID())
		require.NoError(t, err)

		var consumed null.Bool
		err = db.Get(&consumed, selectQuery, queryArgs...)
		require.NoError(t, err)
		require.Equal(t, null.BoolFrom(false), consumed)
	}))

	t.Run("WasBroadcastConsumed_false", func(t *testing.T) {
		was, err := orm.WasBroadcastConsumed(rawLog.BlockHash, rawLog.Index, listener.JobID())
		require.NoError(t, err)
		require.False(t, was)
	})

	require.True(t, t.Run("MarkBroadcastConsumed", func(t *testing.T) {
		err := orm.MarkBroadcastConsumed(rawLog.BlockHash, rawLog.BlockNumber, rawLog.Index, listener.JobID())
		require.NoError(t, err)

		var consumed null.Bool
		err = db.Get(&consumed, selectQuery, queryArgs...)
		require.NoError(t, err)
		require.Equal(t, null.BoolFrom(true), consumed)
	}))

	t.Run("MarkBroadcastsConsumed Success", func(t *testing.T) {
		var (
			err          error
			blockHashes  []common.Hash
			blockNumbers []uint64
			logIndexes   []uint
			jobIDs       []int32
		)
		for i := 0; i < 3; i++ {
			l := cltest.RandomLog(t)
			err = orm.CreateBroadcast(l.BlockHash, l.BlockNumber, l.Index, listener.JobID())
			require.NoError(t, err)
			blockHashes = append(blockHashes, l.BlockHash)
			blockNumbers = append(blockNumbers, l.BlockNumber)
			logIndexes = append(logIndexes, l.Index)
			jobIDs = append(jobIDs, listener.JobID())

		}
		err = orm.MarkBroadcastsConsumed(blockHashes, blockNumbers, logIndexes, jobIDs)
		require.NoError(t, err)

		for i := range blockHashes {
			was, err := orm.WasBroadcastConsumed(blockHashes[i], logIndexes[i], jobIDs[i])
			require.NoError(t, err)
			require.True(t, was)
		}
	})

	t.Run("MarkBroadcastsConsumed Failure", func(t *testing.T) {
		var (
			err          error
			blockHashes  []common.Hash
			blockNumbers []uint64
			logIndexes   []uint
			jobIDs       []int32
		)
		for i := 0; i < 5; i++ {
			l := cltest.RandomLog(t)
			err = orm.CreateBroadcast(l.BlockHash, l.BlockNumber, l.Index, listener.JobID())
			require.NoError(t, err)
			blockHashes = append(blockHashes, l.BlockHash)
			blockNumbers = append(blockNumbers, l.BlockNumber)
			logIndexes = append(logIndexes, l.Index)
			jobIDs = append(jobIDs, listener.JobID())

		}
		err = orm.MarkBroadcastsConsumed(blockHashes[:len(blockHashes)-2], blockNumbers, logIndexes, jobIDs)
		require.Error(t, err)
	})

	t.Run("WasBroadcastConsumed_true", func(t *testing.T) {
		was, err := orm.WasBroadcastConsumed(rawLog.BlockHash, rawLog.Index, listener.JobID())
		require.NoError(t, err)
		require.True(t, was)
	})
}

func TestORM_pending(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := cltest.NewTestGeneralConfig(t)
	lggr := logger.TestLogger(t)
	orm := log.NewORM(db, lggr, cfg, cltest.FixtureChainID)

	num, err := orm.GetPendingMinBlock()
	require.NoError(t, err)
	require.Nil(t, num)

	var num10 int64 = 10
	err = orm.SetPendingMinBlock(&num10)
	require.NoError(t, err)

	num, err = orm.GetPendingMinBlock()
	require.NoError(t, err)
	require.Equal(t, num10, *num)

	err = orm.SetPendingMinBlock(nil)
	require.NoError(t, err)

	num, err = orm.GetPendingMinBlock()
	require.NoError(t, err)
	require.Nil(t, num)
}

func TestORM_MarkUnconsumed(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := cltest.NewTestGeneralConfig(t)
	lggr := logger.TestLogger(t)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	orm := log.NewORM(db, lggr, cfg, cltest.FixtureChainID)

	_, addr1 := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)
	job1 := cltest.MustInsertV2JobSpec(t, db, addr1)

	_, addr2 := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)
	job2 := cltest.MustInsertV2JobSpec(t, db, addr2)

	logBefore := cltest.RandomLog(t)
	logBefore.BlockNumber = 34
	require.NoError(t,
		orm.CreateBroadcast(logBefore.BlockHash, logBefore.BlockNumber, logBefore.Index, job1.ID))
	require.NoError(t,
		orm.MarkBroadcastConsumed(logBefore.BlockHash, logBefore.BlockNumber, logBefore.Index, job1.ID))

	logAt := cltest.RandomLog(t)
	logAt.BlockNumber = 38
	require.NoError(t,
		orm.CreateBroadcast(logAt.BlockHash, logAt.BlockNumber, logAt.Index, job1.ID))
	require.NoError(t,
		orm.MarkBroadcastConsumed(logAt.BlockHash, logAt.BlockNumber, logAt.Index, job1.ID))

	logAfter := cltest.RandomLog(t)
	logAfter.BlockNumber = 40
	require.NoError(t,
		orm.CreateBroadcast(logAfter.BlockHash, logAfter.BlockNumber, logAfter.Index, job2.ID))
	require.NoError(t,
		orm.MarkBroadcastConsumed(logAfter.BlockHash, logAfter.BlockNumber, logAfter.Index, job2.ID))

	// logAt and logAfter should now be marked unconsumed. logBefore is still consumed.
	require.NoError(t, orm.MarkBroadcastsUnconsumed(38))

	consumed, err := orm.WasBroadcastConsumed(logBefore.BlockHash, logBefore.Index, job1.ID)
	require.NoError(t, err)
	require.True(t, consumed)

	consumed, err = orm.WasBroadcastConsumed(logAt.BlockHash, logAt.Index, job1.ID)
	require.NoError(t, err)
	require.False(t, consumed)

	consumed, err = orm.WasBroadcastConsumed(logAfter.BlockHash, logAfter.Index, job2.ID)
	require.NoError(t, err)
	require.False(t, consumed)
}

func TestORM_Reinitialize(t *testing.T) {
	type TestLogBroadcast struct {
		BlockNumber big.Int
		log.LogBroadcast
	}
	var unconsumed = func(blockNum int64) TestLogBroadcast {
		hash := common.BigToHash(big.NewInt(rand.Int63()))
		return TestLogBroadcast{*big.NewInt(blockNum),
			log.LogBroadcast{hash, false, uint(rand.Uint32()), 0},
		}
	}
	var consumed = func(blockNum int64) TestLogBroadcast {
		hash := common.BigToHash(big.NewInt(rand.Int63()))
		return TestLogBroadcast{*big.NewInt(blockNum),
			log.LogBroadcast{hash, true, uint(rand.Uint32()), 0},
		}
	}

	tests := []struct {
		name               string
		pendingBlockNum    *int64
		expPendingBlockNum *int64
		broadcasts         []TestLogBroadcast
	}{
		{name: "empty", expPendingBlockNum: nil},
		{name: "both-delete", expPendingBlockNum: null.IntFrom(10).Ptr(),
			pendingBlockNum: null.IntFrom(10).Ptr(), broadcasts: []TestLogBroadcast{
				unconsumed(11), unconsumed(12),
				consumed(9),
			}},
		{name: "both-update", expPendingBlockNum: null.IntFrom(9).Ptr(),
			pendingBlockNum: null.IntFrom(10).Ptr(), broadcasts: []TestLogBroadcast{
				unconsumed(9), unconsumed(10),
				consumed(8),
			}},
		{name: "broadcasts-update", expPendingBlockNum: null.IntFrom(9).Ptr(),
			pendingBlockNum: nil, broadcasts: []TestLogBroadcast{
				unconsumed(9), unconsumed(10),
				consumed(8),
			}},
		{name: "pending-noop", expPendingBlockNum: null.IntFrom(10).Ptr(),
			pendingBlockNum: null.IntFrom(10).Ptr(), broadcasts: []TestLogBroadcast{
				consumed(8), consumed(9),
			}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db := pgtest.NewSqlxDB(t)
			cfg := cltest.NewTestGeneralConfig(t)
			lggr := logger.TestLogger(t)
			orm := log.NewORM(db, lggr, cfg, cltest.FixtureChainID)

			jobID := cltest.MustInsertV2JobSpec(t, db, common.BigToAddress(big.NewInt(rand.Int63()))).ID

			for _, b := range tt.broadcasts {
				if b.Consumed {
					err := orm.MarkBroadcastConsumed(b.BlockHash, b.BlockNumber.Uint64(), b.LogIndex, jobID)
					require.NoError(t, err)
				} else {
					err := orm.CreateBroadcast(b.BlockHash, b.BlockNumber.Uint64(), b.LogIndex, jobID)
					require.NoError(t, err)
				}
			}
			if tt.pendingBlockNum != nil {
				require.NoError(t, orm.SetPendingMinBlock(tt.pendingBlockNum))
			}

			pendingBlockNum, err := orm.Reinitialize()
			require.NoError(t, err)
			assert.Equal(t, tt.expPendingBlockNum, pendingBlockNum)

			pendingBlockNum, err = orm.GetPendingMinBlock()
			if assert.NoError(t, err) {
				assert.Equal(t, tt.expPendingBlockNum, pendingBlockNum)
			}

			bs, err := orm.FindBroadcasts(0, 20)
			if assert.NoError(t, err) {
				for _, b := range bs {
					assert.True(t, b.Consumed)
				}
			}
		})
	}
}
