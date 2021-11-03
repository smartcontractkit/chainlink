package log_test

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
)

func TestORM_broadcasts(t *testing.T) {
	gdb := pgtest.NewGormDB(t)
	db := postgres.UnwrapGormDB(gdb)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	orm := log.NewORM(db, cltest.FixtureChainID)

	_, addr := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)
	specV2 := cltest.MustInsertV2JobSpec(t, gdb, addr)

	const selectQuery = `SELECT consumed FROM log_broadcasts
		WHERE block_hash = ? AND block_number = ? AND log_index = ? AND job_id = ? AND evm_chain_id = ?`

	listener := &mockListener{specV2.ID}

	rawLog := cltest.RandomLog(t)
	queryArgs := []interface{}{rawLog.BlockHash, rawLog.BlockNumber, rawLog.Index, listener.JobID(), cltest.FixtureChainID.String()}

	// No rows
	q := gdb.Exec(selectQuery, queryArgs...)
	require.NoError(t, q.Error)
	require.Zero(t, q.RowsAffected)

	t.Run("WasBroadcastConsumed_DNE", func(t *testing.T) {
		_, err := orm.WasBroadcastConsumed(rawLog.BlockHash, rawLog.Index, listener.JobID())
		require.NoError(t, err)
	})

	require.True(t, t.Run("CreateBroadcast", func(t *testing.T) {
		err := orm.CreateBroadcast(rawLog.BlockHash, rawLog.BlockNumber, rawLog.Index, listener.JobID())
		require.NoError(t, err)

		var consumed null.Bool
		err = gdb.Raw(selectQuery, queryArgs...).Row().Scan(&consumed)
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
		err = gdb.Raw(selectQuery, queryArgs...).Row().Scan(&consumed)
		require.NoError(t, err)
		require.Equal(t, null.BoolFrom(true), consumed)
	}))

	t.Run("WasBroadcastConsumed_true", func(t *testing.T) {
		was, err := orm.WasBroadcastConsumed(rawLog.BlockHash, rawLog.Index, listener.JobID())
		require.NoError(t, err)
		require.True(t, was)
	})
}

func TestORM_pending(t *testing.T) {
	gdb := pgtest.NewGormDB(t)
	db := postgres.UnwrapGormDB(gdb)
	orm := log.NewORM(db, cltest.FixtureChainID)

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
			gdb := pgtest.NewGormDB(t)
			db := postgres.UnwrapGormDB(gdb)
			orm := log.NewORM(db, cltest.FixtureChainID)

			jobID := cltest.MustInsertV2JobSpec(t, gdb, common.BigToAddress(big.NewInt(rand.Int63()))).ID

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
