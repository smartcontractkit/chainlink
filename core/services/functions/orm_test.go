package functions_test

import (
	"errors"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/functions"
)

var (
	defaultFlags               = []byte{0x1, 0x2, 0x3}
	defaultAggregationMethod   = functions.AggregationMethod(65)
	defaultGasLimit            = uint32(100_000)
	defaultCoordinatorContract = common.HexToAddress("0x0000000000000000000000000000000000000000")
	defaultMetadata            = []byte{0xbb}
)

func setupORM(t *testing.T) functions.ORM {
	t.Helper()

	var (
		db       = pgtest.NewSqlxDB(t)
		contract = testutils.NewAddress()
		orm      = functions.NewORM(db, contract)
	)

	return orm
}

func newRequestID() functions.RequestID {
	return testutils.Random32Byte()
}

func createRequest(t *testing.T, orm functions.ORM) (functions.RequestID, common.Hash, time.Time) {
	ts := time.Now().Round(time.Second)
	id, hash := createRequestWithTimestamp(t, orm, ts)
	return id, hash, ts
}

func createRequestWithTimestamp(t *testing.T, orm functions.ORM, ts time.Time) (functions.RequestID, common.Hash) {
	ctx := testutils.Context(t)
	id := newRequestID()
	txHash := utils.RandomHash()
	newReq := &functions.Request{
		RequestID:                  id,
		RequestTxHash:              &txHash,
		ReceivedAt:                 ts,
		Flags:                      defaultFlags,
		AggregationMethod:          &defaultAggregationMethod,
		CallbackGasLimit:           &defaultGasLimit,
		CoordinatorContractAddress: &defaultCoordinatorContract,
		OnchainMetadata:            defaultMetadata,
	}
	err := orm.CreateRequest(ctx, newReq)
	require.NoError(t, err)
	return id, txHash
}

func TestORM_CreateRequestsAndFindByID(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	orm := setupORM(t)
	id1, txHash1, ts1 := createRequest(t, orm)
	id2, txHash2, ts2 := createRequest(t, orm)

	req1, err := orm.FindById(ctx, id1)
	require.NoError(t, err)
	require.Equal(t, id1, req1.RequestID)
	require.Equal(t, &txHash1, req1.RequestTxHash)
	require.Equal(t, ts1, req1.ReceivedAt)
	require.Equal(t, functions.IN_PROGRESS, req1.State)
	require.Equal(t, defaultFlags, req1.Flags)
	require.Equal(t, defaultAggregationMethod, *req1.AggregationMethod)
	require.Equal(t, defaultGasLimit, *req1.CallbackGasLimit)
	require.Equal(t, defaultCoordinatorContract, *req1.CoordinatorContractAddress)
	require.Equal(t, defaultMetadata, req1.OnchainMetadata)

	req2, err := orm.FindById(ctx, id2)
	require.NoError(t, err)
	require.Equal(t, id2, req2.RequestID)
	require.Equal(t, &txHash2, req2.RequestTxHash)
	require.Equal(t, ts2, req2.ReceivedAt)
	require.Equal(t, functions.IN_PROGRESS, req2.State)

	t.Run("missing ID", func(t *testing.T) {
		req, err := orm.FindById(testutils.Context(t), newRequestID())
		require.Error(t, err)
		require.Nil(t, req)
	})

	t.Run("duplicated", func(t *testing.T) {
		newReq := &functions.Request{RequestID: id1, RequestTxHash: &txHash1, ReceivedAt: ts1}
		err := orm.CreateRequest(testutils.Context(t), newReq)
		require.Error(t, err)
		require.True(t, errors.Is(err, functions.ErrDuplicateRequestID))
	})
}

func TestORM_SetResult(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	orm := setupORM(t)
	id, _, ts := createRequest(t, orm)

	rdts := time.Now().Round(time.Second)
	err := orm.SetResult(ctx, id, []byte("result"), rdts)
	require.NoError(t, err)

	req, err := orm.FindById(ctx, id)
	require.NoError(t, err)
	require.Equal(t, id, req.RequestID)
	require.Equal(t, ts, req.ReceivedAt)
	require.NotNil(t, req.ResultReadyAt)
	require.Equal(t, rdts, *req.ResultReadyAt)
	require.Equal(t, functions.RESULT_READY, req.State)
	require.Equal(t, []byte("result"), req.Result)
}

func TestORM_SetError(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	orm := setupORM(t)
	id, _, ts := createRequest(t, orm)

	rdts := time.Now().Round(time.Second)
	err := orm.SetError(ctx, id, functions.USER_ERROR, []byte("error"), rdts, true)
	require.NoError(t, err)

	req, err := orm.FindById(ctx, id)
	require.NoError(t, err)
	require.Equal(t, id, req.RequestID)
	require.Equal(t, ts, req.ReceivedAt)
	require.NotNil(t, req.ResultReadyAt)
	require.Equal(t, rdts, *req.ResultReadyAt)
	require.NotNil(t, req.ErrorType)
	require.Equal(t, functions.USER_ERROR, *req.ErrorType)
	require.Equal(t, functions.RESULT_READY, req.State)
	require.Equal(t, []byte("error"), req.Error)
}

func TestORM_SetError_Internal(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	orm := setupORM(t)
	id, _, ts := createRequest(t, orm)

	rdts := time.Now().Round(time.Second)
	err := orm.SetError(ctx, id, functions.INTERNAL_ERROR, []byte("error"), rdts, false)
	require.NoError(t, err)

	req, err := orm.FindById(ctx, id)
	require.NoError(t, err)
	require.Equal(t, id, req.RequestID)
	require.Equal(t, ts, req.ReceivedAt)
	require.Equal(t, functions.INTERNAL_ERROR, *req.ErrorType)
	require.Equal(t, functions.IN_PROGRESS, req.State)
	require.Equal(t, []byte("error"), req.Error)
}

func TestORM_SetFinalized(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	orm := setupORM(t)
	id, _, _ := createRequest(t, orm)

	err := orm.SetFinalized(ctx, id, []byte("result"), []byte("error"))
	require.NoError(t, err)

	req, err := orm.FindById(ctx, id)
	require.NoError(t, err)
	require.Equal(t, []byte("result"), req.TransmittedResult)
	require.Equal(t, []byte("error"), req.TransmittedError)
	require.Equal(t, functions.FINALIZED, req.State)
}

func TestORM_SetConfirmed(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	orm := setupORM(t)
	id, _, _ := createRequest(t, orm)

	err := orm.SetConfirmed(ctx, id)
	require.NoError(t, err)

	req, err := orm.FindById(ctx, id)
	require.NoError(t, err)
	require.Equal(t, functions.CONFIRMED, req.State)
}

func TestORM_StateTransitions(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	orm := setupORM(t)
	now := time.Now()
	id, _ := createRequestWithTimestamp(t, orm, now)
	req, err := orm.FindById(ctx, id)
	require.NoError(t, err)
	require.Equal(t, functions.IN_PROGRESS, req.State)

	err = orm.SetResult(ctx, id, []byte{}, now)
	require.NoError(t, err)
	req, err = orm.FindById(ctx, id)
	require.NoError(t, err)
	require.Equal(t, functions.RESULT_READY, req.State)

	_, err = orm.TimeoutExpiredResults(ctx, now.Add(time.Minute), 1)
	require.NoError(t, err)
	req, err = orm.FindById(ctx, id)
	require.NoError(t, err)
	require.Equal(t, functions.TIMED_OUT, req.State)

	err = orm.SetFinalized(ctx, id, nil, nil)
	require.Error(t, err)
	req, err = orm.FindById(ctx, id)
	require.NoError(t, err)
	require.Equal(t, functions.TIMED_OUT, req.State)

	err = orm.SetConfirmed(ctx, id)
	require.NoError(t, err)
	req, err = orm.FindById(ctx, id)
	require.NoError(t, err)
	require.Equal(t, functions.CONFIRMED, req.State)
}

func TestORM_FindOldestEntriesByState(t *testing.T) {
	t.Parallel()

	orm := setupORM(t)
	now := time.Now()
	id2, _ := createRequestWithTimestamp(t, orm, now.Add(2*time.Minute))
	createRequestWithTimestamp(t, orm, now.Add(3*time.Minute))
	id1, _ := createRequestWithTimestamp(t, orm, now.Add(1*time.Minute))

	t.Run("with limit", func(t *testing.T) {
		ctx := testutils.Context(t)
		result, err := orm.FindOldestEntriesByState(ctx, functions.IN_PROGRESS, 2)
		require.NoError(t, err)
		require.Equal(t, 2, len(result), "incorrect results length")
		require.Equal(t, id1, result[0].RequestID, "incorrect results order")
		require.Equal(t, id2, result[1].RequestID, "incorrect results order")

		require.Equal(t, defaultFlags, result[0].Flags)
		require.Equal(t, defaultAggregationMethod, *result[0].AggregationMethod)
		require.Equal(t, defaultGasLimit, *result[0].CallbackGasLimit)
		require.Equal(t, defaultCoordinatorContract, *result[0].CoordinatorContractAddress)
		require.Equal(t, defaultMetadata, result[0].OnchainMetadata)
	})

	t.Run("with no limit", func(t *testing.T) {
		ctx := testutils.Context(t)
		result, err := orm.FindOldestEntriesByState(ctx, functions.IN_PROGRESS, 20)
		require.NoError(t, err)
		require.Equal(t, 3, len(result), "incorrect results length")
	})

	t.Run("no matching entries", func(t *testing.T) {
		ctx := testutils.Context(t)
		result, err := orm.FindOldestEntriesByState(ctx, functions.RESULT_READY, 10)
		require.NoError(t, err)
		require.Equal(t, 0, len(result), "incorrect results length")
	})
}

func TestORM_TimeoutExpiredResults(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	orm := setupORM(t)
	now := time.Now()
	var ids []functions.RequestID
	for offset := -50; offset <= -10; offset += 10 {
		id, _ := createRequestWithTimestamp(t, orm, now.Add(time.Duration(offset)*time.Minute))
		ids = append(ids, id)
	}
	// can time out IN_PROGRESS, RESULT_READY or FINALIZED
	err := orm.SetResult(ctx, ids[0], []byte("result"), now)
	require.NoError(t, err)
	err = orm.SetFinalized(ctx, ids[1], []byte("result"), []byte(""))
	require.NoError(t, err)
	// can't time out CONFIRMED
	err = orm.SetConfirmed(ctx, ids[2])
	require.NoError(t, err)

	results, err := orm.TimeoutExpiredResults(ctx, now.Add(-35*time.Minute), 1)
	require.NoError(t, err)
	require.Equal(t, 1, len(results), "not respecting limit")
	require.Equal(t, ids[0], results[0], "incorrect results order")

	results, err = orm.TimeoutExpiredResults(ctx, now.Add(-15*time.Minute), 10)
	require.NoError(t, err)
	require.Equal(t, 2, len(results), "incorrect results length")
	require.Equal(t, ids[1], results[0], "incorrect results order")
	require.Equal(t, ids[3], results[1], "incorrect results order")

	results, err = orm.TimeoutExpiredResults(ctx, now.Add(-15*time.Minute), 10)
	require.NoError(t, err)
	require.Equal(t, 0, len(results), "not idempotent")

	expectedFinalStates := []functions.RequestState{
		functions.TIMED_OUT,
		functions.TIMED_OUT,
		functions.CONFIRMED,
		functions.TIMED_OUT,
		functions.IN_PROGRESS,
	}
	for i, expectedState := range expectedFinalStates {
		req, err := orm.FindById(ctx, ids[i])
		require.NoError(t, err)
		require.Equal(t, req.State, expectedState, "incorrect state")
	}
}

func TestORM_PruneOldestRequests(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	orm := setupORM(t)
	now := time.Now()
	var ids []functions.RequestID
	// store 5 requests
	for offset := -50; offset <= -10; offset += 10 {
		id, _ := createRequestWithTimestamp(t, orm, now.Add(time.Duration(offset)*time.Minute))
		ids = append(ids, id)
	}

	// don't prune if max not hit
	total, pruned, err := orm.PruneOldestRequests(ctx, 6, 3)
	require.NoError(t, err)
	require.Equal(t, uint32(5), total)
	require.Equal(t, uint32(0), pruned)

	// prune up to max batch size
	total, pruned, err = orm.PruneOldestRequests(ctx, 1, 2)
	require.NoError(t, err)
	require.Equal(t, uint32(5), total)
	require.Equal(t, uint32(2), pruned)

	// prune all above the limit
	total, pruned, err = orm.PruneOldestRequests(ctx, 1, 20)
	require.NoError(t, err)
	require.Equal(t, uint32(3), total)
	require.Equal(t, uint32(2), pruned)

	// no pruning needed any more
	total, pruned, err = orm.PruneOldestRequests(ctx, 1, 20)
	require.NoError(t, err)
	require.Equal(t, uint32(1), total)
	require.Equal(t, uint32(0), pruned)

	// verify only the newest one is left after pruning
	result, err := orm.FindOldestEntriesByState(ctx, functions.IN_PROGRESS, 20)
	require.NoError(t, err)
	require.Equal(t, 1, len(result), "incorrect results length")
	require.Equal(t, ids[4], result[0].RequestID, "incorrect results order")
}

func TestORM_PruneOldestRequests_Large(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	orm := setupORM(t)
	now := time.Now()
	// store 1000 requests
	for offset := -1000; offset <= -1; offset++ {
		_, _ = createRequestWithTimestamp(t, orm, now.Add(time.Duration(offset)*time.Minute))
	}

	// prune 900/1000
	total, pruned, err := orm.PruneOldestRequests(ctx, 100, 1000)
	require.NoError(t, err)
	require.Equal(t, uint32(1000), total)
	require.Equal(t, uint32(900), pruned)

	// verify there's 100 left
	result, err := orm.FindOldestEntriesByState(ctx, functions.IN_PROGRESS, 200)
	require.NoError(t, err)
	require.Equal(t, 100, len(result), "incorrect results length")
}
