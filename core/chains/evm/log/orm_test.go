package log_test

import (
	"context"
	crand "crypto/rand"
	"math/big"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	evmutils "github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
)

func TestORM_broadcasts(t *testing.T) {
	db := testutils.NewSqlxDB(t)

	orm := log.NewORM(db, *testutils.FixtureChainID)

	jobID := mustInsertV2JobSpec(t, db, testutils.NewAddress())

	const selectQuery = `SELECT consumed FROM log_broadcasts
		WHERE block_hash = $1 AND block_number = $2 AND log_index = $3 AND job_id = $4 AND evm_chain_id = $5`

	listener := &mockListener{jobID}

	rawLog := randomLog(t)
	queryArgs := []interface{}{rawLog.BlockHash, rawLog.BlockNumber, rawLog.Index, listener.JobID(), testutils.FixtureChainID.String()}

	// No rows
	res, err := db.Exec(selectQuery, queryArgs...)
	require.NoError(t, err)
	rowsAffected, err := res.RowsAffected()
	require.NoError(t, err)
	require.Zero(t, rowsAffected)

	t.Run("WasBroadcastConsumed_DNE", func(t *testing.T) {
		_, err := orm.WasBroadcastConsumed(testutils.Context(t), rawLog.BlockHash, rawLog.Index, listener.JobID())
		require.NoError(t, err)
	})

	require.True(t, t.Run("CreateBroadcast", func(t *testing.T) {
		err := orm.CreateBroadcast(testutils.Context(t), rawLog.BlockHash, rawLog.BlockNumber, rawLog.Index, listener.JobID())
		require.NoError(t, err)

		var consumed null.Bool
		err = db.Get(&consumed, selectQuery, queryArgs...)
		require.NoError(t, err)
		require.Equal(t, null.BoolFrom(false), consumed)
	}))

	t.Run("WasBroadcastConsumed_false", func(t *testing.T) {
		was, err := orm.WasBroadcastConsumed(testutils.Context(t), rawLog.BlockHash, rawLog.Index, listener.JobID())
		require.NoError(t, err)
		require.False(t, was)
	})

	require.True(t, t.Run("MarkBroadcastConsumed", func(t *testing.T) {
		err := orm.MarkBroadcastConsumed(testutils.Context(t), rawLog.BlockHash, rawLog.BlockNumber, rawLog.Index, listener.JobID())
		require.NoError(t, err)

		var consumed null.Bool
		err = db.Get(&consumed, selectQuery, queryArgs...)
		require.NoError(t, err)
		require.Equal(t, null.BoolFrom(true), consumed)
	}))

	t.Run("WasBroadcastConsumed_true", func(t *testing.T) {
		was, err := orm.WasBroadcastConsumed(testutils.Context(t), rawLog.BlockHash, rawLog.Index, listener.JobID())
		require.NoError(t, err)
		require.True(t, was)
	})
}

func TestORM_pending(t *testing.T) {
	ctx := testutils.Context(t)
	db := testutils.NewSqlxDB(t)
	orm := log.NewORM(db, *testutils.FixtureChainID)

	num, err := orm.GetPendingMinBlock(ctx)
	require.NoError(t, err)
	require.Nil(t, num)

	var num10 int64 = 10
	err = orm.SetPendingMinBlock(ctx, &num10)
	require.NoError(t, err)

	num, err = orm.GetPendingMinBlock(ctx)
	require.NoError(t, err)
	require.Equal(t, num10, *num)

	err = orm.SetPendingMinBlock(ctx, nil)
	require.NoError(t, err)

	num, err = orm.GetPendingMinBlock(ctx)
	require.NoError(t, err)
	require.Nil(t, num)
}

func TestORM_MarkUnconsumed(t *testing.T) {
	ctx := testutils.Context(t)
	db := testutils.NewSqlxDB(t)

	orm := log.NewORM(db, *testutils.FixtureChainID)

	job1 := mustInsertV2JobSpec(t, db, testutils.NewAddress())

	job2 := mustInsertV2JobSpec(t, db, testutils.NewAddress())

	logBefore := randomLog(t)
	logBefore.BlockNumber = 34
	require.NoError(t,
		orm.CreateBroadcast(ctx, logBefore.BlockHash, logBefore.BlockNumber, logBefore.Index, job1))
	require.NoError(t,
		orm.MarkBroadcastConsumed(ctx, logBefore.BlockHash, logBefore.BlockNumber, logBefore.Index, job1))

	logAt := randomLog(t)
	logAt.BlockNumber = 38
	require.NoError(t,
		orm.CreateBroadcast(ctx, logAt.BlockHash, logAt.BlockNumber, logAt.Index, job1))
	require.NoError(t,
		orm.MarkBroadcastConsumed(ctx, logAt.BlockHash, logAt.BlockNumber, logAt.Index, job1))

	logAfter := randomLog(t)
	logAfter.BlockNumber = 40
	require.NoError(t,
		orm.CreateBroadcast(ctx, logAfter.BlockHash, logAfter.BlockNumber, logAfter.Index, job2))
	require.NoError(t,
		orm.MarkBroadcastConsumed(ctx, logAfter.BlockHash, logAfter.BlockNumber, logAfter.Index, job2))

	// logAt and logAfter should now be marked unconsumed. logBefore is still consumed.
	require.NoError(t, orm.MarkBroadcastsUnconsumed(ctx, 38))

	consumed, err := orm.WasBroadcastConsumed(ctx, logBefore.BlockHash, logBefore.Index, job1)
	require.NoError(t, err)
	require.True(t, consumed)

	consumed, err = orm.WasBroadcastConsumed(ctx, logAt.BlockHash, logAt.Index, job1)
	require.NoError(t, err)
	require.False(t, consumed)

	consumed, err = orm.WasBroadcastConsumed(ctx, logAfter.BlockHash, logAfter.Index, job2)
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
			db := testutils.NewSqlxDB(t)
			ctx := testutils.Context(t)
			orm := log.NewORM(db, *testutils.FixtureChainID)

			jobID := mustInsertV2JobSpec(t, db, common.BigToAddress(big.NewInt(rand.Int63())))

			for _, b := range tt.broadcasts {
				if b.Consumed {
					err := orm.MarkBroadcastConsumed(ctx, b.BlockHash, b.BlockNumber.Uint64(), b.LogIndex, jobID)
					require.NoError(t, err)
				} else {
					err := orm.CreateBroadcast(ctx, b.BlockHash, b.BlockNumber.Uint64(), b.LogIndex, jobID)
					require.NoError(t, err)
				}
			}
			if tt.pendingBlockNum != nil {
				require.NoError(t, orm.SetPendingMinBlock(ctx, tt.pendingBlockNum))
			}

			pendingBlockNum, err := orm.Reinitialize(ctx)
			require.NoError(t, err)
			assert.Equal(t, tt.expPendingBlockNum, pendingBlockNum)

			pendingBlockNum, err = orm.GetPendingMinBlock(ctx)
			if assert.NoError(t, err) {
				assert.Equal(t, tt.expPendingBlockNum, pendingBlockNum)
			}

			bs, err := orm.FindBroadcasts(ctx, 0, 20)
			if assert.NoError(t, err) {
				for _, b := range bs {
					assert.True(t, b.Consumed)
				}
			}
		})
	}
}

type mockListener struct {
	jobID int32
}

func (l *mockListener) JobID() int32                             { return l.jobID }
func (l *mockListener) HandleLog(context.Context, log.Broadcast) {}

func randomLog(t *testing.T) types.Log {
	t.Helper()

	topics := make([]common.Hash, 4)
	for i := range topics {
		topics[i] = evmutils.NewHash()
	}

	return types.Log{
		Address:     testutils.NewAddress(),
		BlockHash:   evmutils.NewHash(),
		BlockNumber: uint64(rand.Intn(9999999)),
		Index:       uint(rand.Intn(9999999)),
		Data:        randomBytes(t, 512),
		Topics:      []common.Hash{evmutils.NewHash(), evmutils.NewHash(), evmutils.NewHash(), evmutils.NewHash()},
	}
}

func randomBytes(t *testing.T, n int) []byte {
	b := make([]byte, n)
	_, err := crand.Read(b)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func mustInsertV2JobSpec(t *testing.T, db *sqlx.DB, transmitterAddress common.Address) int32 {
	t.Helper()

	addr := evmtypes.MustEIP55Address(transmitterAddress.Hex())

	var oracleSpecId int32
	require.NoError(t, db.Get(&oracleSpecId, `INSERT INTO ocr_oracle_specs (created_at, updated_at, contract_address, p2pv2_bootstrappers, is_bootstrap_peer, transmitter_address, evm_chain_id, contract_config_confirmations) VALUES (
	NOW(),NOW(),$1,'{}',false,$2,0,0
	) RETURNING id`, evmtypes.MustEIP55Address(testutils.NewAddress().Hex()), &addr))

	var jobSpecId int32
	require.NoError(t, db.Get(&jobSpecId, `INSERT INTO jobs (schema_version, created_at, type, ocr_oracle_spec_id, external_job_id) VALUES (
    1,NOW(),0,$1,$2) RETURNING id`, oracleSpecId, uuid.New()))

	return jobSpecId
}
