package loader

import (
	"database/sql"
	"testing"

	"github.com/graph-gophers/dataloader"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	evmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	txmgrMocks "github.com/smartcontractkit/chainlink/core/chains/evm/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	coremocks "github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/feeds"
	feedsMocks "github.com/smartcontractkit/chainlink/core/services/feeds/mocks"
	"github.com/smartcontractkit/chainlink/core/services/job"
	jobORMMocks "github.com/smartcontractkit/chainlink/core/services/job/mocks"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestLoader_Chains(t *testing.T) {
	t.Parallel()

	app := &coremocks.Application{}
	ctx := InjectDataloader(testutils.Context(t), app)

	defer t.Cleanup(func() {
		mock.AssertExpectationsForObjects(t, app)
	})

	id := utils.Big{}
	err := id.UnmarshalText([]byte("1"))
	require.NoError(t, err)

	id2 := utils.Big{}
	err = id2.UnmarshalText([]byte("2"))
	require.NoError(t, err)

	chainId3 := utils.Big{}
	err = chainId3.UnmarshalText([]byte("3"))
	require.NoError(t, err)

	chain := types.DBChain{
		ID:      id,
		Enabled: true,
	}
	chain2 := types.DBChain{
		ID:      id2,
		Enabled: true,
	}
	evmORM := evmtest.NewMockORM([]types.DBChain{chain, chain2}, nil)
	app.On("EVMORM").Return(evmORM)

	batcher := chainBatcher{app}

	keys := dataloader.NewKeysFromStrings([]string{"2", "1", "3"})
	results := batcher.loadByIDs(ctx, keys)

	assert.Len(t, results, 3)
	assert.Equal(t, chain2, results[0].Data.(types.DBChain))
	assert.Equal(t, chain, results[1].Data.(types.DBChain))
	assert.Nil(t, results[2].Data)
	assert.Error(t, results[2].Error)
	assert.Equal(t, "chain not found", results[2].Error.Error())
}

func TestLoader_Nodes(t *testing.T) {
	t.Parallel()

	evmChainSet := evmmocks.NewChainSet(t)
	app := coremocks.NewApplication(t)
	ctx := InjectDataloader(testutils.Context(t), app)

	defer t.Cleanup(func() {
		mock.AssertExpectationsForObjects(t, app, evmChainSet)
	})

	chainId1 := utils.Big{}
	err := chainId1.UnmarshalText([]byte("1"))
	require.NoError(t, err)

	chainId2 := utils.Big{}
	err = chainId2.UnmarshalText([]byte("2"))
	require.NoError(t, err)

	chainId3 := utils.Big{}
	err = chainId3.UnmarshalText([]byte("3"))
	require.NoError(t, err)

	node1 := types.Node{
		ID:         int32(1),
		Name:       "test-node-1",
		EVMChainID: chainId1,
	}
	node2 := types.Node{
		ID:         int32(2),
		Name:       "test-node-1",
		EVMChainID: chainId2,
	}

	evmChainSet.On("GetNodesByChainIDs", mock.Anything, []utils.Big{chainId2, chainId1, chainId3}).Return([]types.Node{
		node1, node2,
	}, nil)
	app.On("GetChains").Return(chainlink.Chains{EVM: evmChainSet})

	batcher := nodeBatcher{app}

	keys := dataloader.NewKeysFromStrings([]string{"2", "1", "3"})
	found := batcher.loadByChainIDs(ctx, keys)

	require.Len(t, found, 3)
	assert.Equal(t, []types.Node{node2}, found[0].Data)
	assert.Equal(t, []types.Node{node1}, found[1].Data)
	assert.Equal(t, []types.Node{}, found[2].Data)
}

func TestLoader_FeedsManagers(t *testing.T) {
	t.Parallel()

	fsvc := &feedsMocks.Service{}
	app := &coremocks.Application{}
	ctx := InjectDataloader(testutils.Context(t), app)

	defer t.Cleanup(func() {
		mock.AssertExpectationsForObjects(t, app, fsvc)
	})

	mgr1 := feeds.FeedsManager{
		ID:   int64(1),
		Name: "manager 1",
	}
	mgr2 := feeds.FeedsManager{
		ID:   int64(2),
		Name: "manager 2",
	}
	mgr3 := feeds.FeedsManager{
		ID:   int64(3),
		Name: "manager 3",
	}

	fsvc.On("ListManagersByIDs", []int64{3, 1, 2, 5}).Return([]feeds.FeedsManager{
		mgr1, mgr2, mgr3,
	}, nil)
	app.On("GetFeedsService").Return(fsvc)

	batcher := feedsBatcher{app}

	keys := dataloader.NewKeysFromStrings([]string{"3", "1", "2", "5"})
	found := batcher.loadByIDs(ctx, keys)

	require.Len(t, found, 4)
	assert.Equal(t, mgr3, found[0].Data)
	assert.Equal(t, mgr1, found[1].Data)
	assert.Equal(t, mgr2, found[2].Data)
	assert.Nil(t, found[3].Data)
	assert.Error(t, found[3].Error)
	assert.Equal(t, "feeds manager not found", found[3].Error.Error())
}

func TestLoader_JobProposals(t *testing.T) {
	t.Parallel()

	fsvc := &feedsMocks.Service{}
	app := &coremocks.Application{}
	ctx := InjectDataloader(testutils.Context(t), app)

	defer t.Cleanup(func() {
		mock.AssertExpectationsForObjects(t, app, fsvc)
	})

	jp1 := feeds.JobProposal{
		ID:             int64(1),
		FeedsManagerID: int64(3),
		Status:         feeds.JobProposalStatusPending,
	}
	jp2 := feeds.JobProposal{
		ID:             int64(2),
		FeedsManagerID: int64(1),
		Status:         feeds.JobProposalStatusApproved,
	}
	jp3 := feeds.JobProposal{
		ID:             int64(3),
		FeedsManagerID: int64(1),
		Status:         feeds.JobProposalStatusRejected,
	}

	fsvc.On("ListJobProposalsByManagersIDs", []int64{3, 1, 2}).Return([]feeds.JobProposal{
		jp1, jp3, jp2,
	}, nil)
	app.On("GetFeedsService").Return(fsvc)

	batcher := jobProposalBatcher{app}

	keys := dataloader.NewKeysFromStrings([]string{"3", "1", "2"})
	found := batcher.loadByManagersIDs(ctx, keys)

	require.Len(t, found, 3)
	assert.Equal(t, []feeds.JobProposal{jp1}, found[0].Data)
	assert.Equal(t, []feeds.JobProposal{jp3, jp2}, found[1].Data)
	assert.Equal(t, []feeds.JobProposal{}, found[2].Data)
}

func TestLoader_JobRuns(t *testing.T) {
	t.Parallel()

	jobsORM := &jobORMMocks.ORM{}
	app := &coremocks.Application{}
	ctx := InjectDataloader(testutils.Context(t), app)

	defer t.Cleanup(func() {
		mock.AssertExpectationsForObjects(t, app, jobsORM)
	})

	run1 := pipeline.Run{ID: int64(1)}
	run2 := pipeline.Run{ID: int64(2)}
	run3 := pipeline.Run{ID: int64(3)}

	jobsORM.On("FindPipelineRunsByIDs", []int64{3, 1, 2}).Return([]pipeline.Run{
		run3, run1, run2,
	}, nil)
	app.On("JobORM").Return(jobsORM)

	batcher := jobRunBatcher{app}

	keys := dataloader.NewKeysFromStrings([]string{"3", "1", "2"})
	found := batcher.loadByIDs(ctx, keys)

	require.Len(t, found, 3)
	assert.Equal(t, run3, found[0].Data)
	assert.Equal(t, run1, found[1].Data)
	assert.Equal(t, run2, found[2].Data)
}

func TestLoader_JobsByPipelineSpecIDs(t *testing.T) {
	t.Parallel()

	t.Run("with out errors", func(t *testing.T) {
		t.Parallel()

		jobsORM := &jobORMMocks.ORM{}
		app := &coremocks.Application{}
		ctx := InjectDataloader(testutils.Context(t), app)

		defer t.Cleanup(func() {
			mock.AssertExpectationsForObjects(t, app, jobsORM)
		})

		job1 := job.Job{ID: int32(2), PipelineSpecID: int32(1)}
		job2 := job.Job{ID: int32(3), PipelineSpecID: int32(2)}
		job3 := job.Job{ID: int32(4), PipelineSpecID: int32(3)}

		jobsORM.On("FindJobsByPipelineSpecIDs", []int32{3, 1, 2}).Return([]job.Job{
			job1, job2, job3,
		}, nil)
		app.On("JobORM").Return(jobsORM)

		batcher := jobBatcher{app}

		keys := dataloader.NewKeysFromStrings([]string{"3", "1", "2"})
		found := batcher.loadByPipelineSpecIDs(ctx, keys)

		require.Len(t, found, 3)
		assert.Equal(t, job3, found[0].Data)
		assert.Equal(t, job1, found[1].Data)
		assert.Equal(t, job2, found[2].Data)
	})

	t.Run("with errors", func(t *testing.T) {
		t.Parallel()

		jobsORM := &jobORMMocks.ORM{}
		app := &coremocks.Application{}
		ctx := InjectDataloader(testutils.Context(t), app)

		defer t.Cleanup(func() {
			mock.AssertExpectationsForObjects(t, app, jobsORM)
		})

		jobsORM.On("FindJobsByPipelineSpecIDs", []int32{3, 1, 2}).Return([]job.Job{}, sql.ErrNoRows)
		app.On("JobORM").Return(jobsORM)

		batcher := jobBatcher{app}

		keys := dataloader.NewKeysFromStrings([]string{"3", "1", "2"})
		found := batcher.loadByPipelineSpecIDs(ctx, keys)

		require.Len(t, found, 1)
		assert.Nil(t, found[0].Data)
		assert.ErrorIs(t, found[0].Error, sql.ErrNoRows)
	})
}

func TestLoader_JobsByExternalJobIDs(t *testing.T) {
	t.Parallel()

	t.Run("with out errors", func(t *testing.T) {
		t.Parallel()

		jobsORM := &jobORMMocks.ORM{}
		app := &coremocks.Application{}
		ctx := InjectDataloader(testutils.Context(t), app)

		defer t.Cleanup(func() {
			mock.AssertExpectationsForObjects(t, app, jobsORM)
		})

		ejID := uuid.NewV4()
		job := job.Job{ID: int32(2), ExternalJobID: ejID}

		jobsORM.On("FindJobByExternalJobID", ejID).Return(job, nil)
		app.On("JobORM").Return(jobsORM)

		batcher := jobBatcher{app}

		keys := dataloader.NewKeysFromStrings([]string{ejID.String()})
		found := batcher.loadByExternalJobIDs(ctx, keys)

		require.Len(t, found, 1)
		assert.Equal(t, job, found[0].Data)
	})
}

func TestLoader_EthTransactionsAttempts(t *testing.T) {
	t.Parallel()

	txmORM := &txmgrMocks.ORM{}
	app := &coremocks.Application{}
	ctx := InjectDataloader(testutils.Context(t), app)

	defer t.Cleanup(func() {
		mock.AssertExpectationsForObjects(t, app, txmORM)
	})

	ethTxIDs := []int64{1, 2, 3}

	attempt1 := txmgr.EthTxAttempt{
		ID:      int64(1),
		EthTxID: ethTxIDs[0],
	}
	attempt2 := txmgr.EthTxAttempt{
		ID:      int64(1),
		EthTxID: ethTxIDs[1],
	}

	txmORM.On("FindEthTxAttemptConfirmedByEthTxIDs", []int64{ethTxIDs[2], ethTxIDs[1], ethTxIDs[0]}).Return([]txmgr.EthTxAttempt{
		attempt1, attempt2,
	}, nil)
	app.On("TxmORM").Return(txmORM)

	batcher := ethTransactionAttemptBatcher{app}

	keys := dataloader.NewKeysFromStrings([]string{"3", "2", "1"})
	found := batcher.loadByEthTransactionIDs(ctx, keys)

	require.Len(t, found, 3)
	assert.Equal(t, []txmgr.EthTxAttempt{}, found[0].Data)
	assert.Equal(t, []txmgr.EthTxAttempt{attempt2}, found[1].Data)
	assert.Equal(t, []txmgr.EthTxAttempt{attempt1}, found[2].Data)
}

func TestLoader_SpecErrorsByJobID(t *testing.T) {
	t.Parallel()

	t.Run("without errors", func(t *testing.T) {
		t.Parallel()

		jobsORM := &jobORMMocks.ORM{}
		app := &coremocks.Application{}
		ctx := InjectDataloader(testutils.Context(t), app)

		defer t.Cleanup(func() {
			mock.AssertExpectationsForObjects(t, app, jobsORM)
		})

		specErr1 := job.SpecError{ID: int64(2), JobID: int32(1)}
		specErr2 := job.SpecError{ID: int64(3), JobID: int32(2)}
		specErr3 := job.SpecError{ID: int64(4), JobID: int32(3)}

		jobsORM.On("FindSpecErrorsByJobIDs", []int32{3, 1, 2}, mock.Anything).Return([]job.SpecError{
			specErr1, specErr2, specErr3,
		}, nil)
		app.On("JobORM").Return(jobsORM)

		batcher := jobSpecErrorsBatcher{app}

		keys := dataloader.NewKeysFromStrings([]string{"3", "1", "2"})
		found := batcher.loadByJobIDs(ctx, keys)

		require.Len(t, found, 3)
		assert.Equal(t, []job.SpecError{specErr3}, found[0].Data)
		assert.Equal(t, []job.SpecError{specErr1}, found[1].Data)
		assert.Equal(t, []job.SpecError{specErr2}, found[2].Data)
	})

	t.Run("with errors", func(t *testing.T) {
		t.Parallel()

		jobsORM := &jobORMMocks.ORM{}
		app := &coremocks.Application{}
		ctx := InjectDataloader(testutils.Context(t), app)

		defer t.Cleanup(func() {
			mock.AssertExpectationsForObjects(t, app, jobsORM)
		})

		jobsORM.On("FindSpecErrorsByJobIDs", []int32{3, 1, 2}, mock.Anything).Return([]job.SpecError{}, sql.ErrNoRows)
		app.On("JobORM").Return(jobsORM)

		batcher := jobSpecErrorsBatcher{app}

		keys := dataloader.NewKeysFromStrings([]string{"3", "1", "2"})
		found := batcher.loadByJobIDs(ctx, keys)

		require.Len(t, found, 1)
		assert.Nil(t, found[0].Data)
		assert.ErrorIs(t, found[0].Error, sql.ErrNoRows)
	})
}

func TestLoader_loadByEthTransactionID(t *testing.T) {
	t.Parallel()

	txmORM := txmgrMocks.NewORM(t)
	app := coremocks.NewApplication(t)
	ctx := InjectDataloader(testutils.Context(t), app)

	ethTxID := int64(3)
	ethTxHash := utils.NewHash()

	receipt := txmgr.EthReceipt{
		ID:     int64(1),
		TxHash: ethTxHash,
	}

	attempt1 := txmgr.EthTxAttempt{
		ID:          int64(1),
		EthTxID:     ethTxID,
		Hash:        ethTxHash,
		EthReceipts: []txmgr.EthReceipt{receipt},
	}

	txmORM.On("FindEthTxAttemptConfirmedByEthTxIDs", []int64{ethTxID}).Return([]txmgr.EthTxAttempt{
		attempt1,
	}, nil)

	app.On("TxmORM").Return(txmORM)

	batcher := ethTransactionAttemptBatcher{app}

	keys := dataloader.NewKeysFromStrings([]string{"3"})
	found := batcher.loadByEthTransactionIDs(ctx, keys)

	require.Len(t, found, 1)
	assert.Equal(t, []txmgr.EthTxAttempt{attempt1}, found[0].Data)
}
