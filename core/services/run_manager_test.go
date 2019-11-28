package services_test

import (
	"math/big"
	"testing"
	"time"

	"chainlink/core/adapters"
	"chainlink/core/eth"
	"chainlink/core/internal/cltest"
	"chainlink/core/internal/mocks"
	clnull "chainlink/core/null"
	"chainlink/core/services"
	"chainlink/core/store/models"
	"chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"
)

func TestRunManager_ResumePending(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	runQueue := new(mocks.RunQueue)
	runQueue.On("Run", mock.Anything).Maybe().Return(nil)

	runManager := services.NewRunManager(runQueue, store.Config, store.ORM, store.TxManager, store.Clock)

	job := cltest.NewJob()
	require.NoError(t, store.CreateJob(&job))
	input := cltest.JSONFromString(t, `{"address":"0xdfcfc2b9200dbb10952c2b7cce60fc7260e03c6f"}`)

	t.Run("reject a run with an invalid state", func(t *testing.T) {
		run := &models.JobRun{ID: models.NewID(), JobSpecID: job.ID}
		require.NoError(t, store.CreateJobRun(run))
		err := runManager.ResumePending(run.ID, models.BridgeRunResult{})
		assert.Error(t, err)
	})

	t.Run("reject a run with no tasks", func(t *testing.T) {
		run := models.JobRun{ID: models.NewID(), JobSpecID: job.ID, Status: models.RunStatusPendingBridge}
		require.NoError(t, store.CreateJobRun(&run))
		err := runManager.ResumePending(run.ID, models.BridgeRunResult{})
		assert.NoError(t, err)

		run, err = store.FindJobRun(run.ID)
		require.NoError(t, err)
		assert.Equal(t, models.RunStatusErrored, run.Status)
	})

	t.Run("input with error errors run", func(t *testing.T) {
		runID := models.NewID()
		run := models.JobRun{
			ID:        runID,
			JobSpecID: job.ID,
			Status:    models.RunStatusPendingBridge,
			TaskRuns:  []models.TaskRun{models.TaskRun{ID: models.NewID(), JobRunID: runID}},
		}
		require.NoError(t, store.CreateJobRun(&run))

		err := runManager.ResumePending(run.ID, models.BridgeRunResult{Status: models.RunStatusErrored})
		assert.NoError(t, err)

		run, err = store.FindJobRun(run.ID)
		require.NoError(t, err)
		assert.Equal(t, models.RunStatusErrored, run.Status)
		assert.True(t, run.FinishedAt.Valid)
		assert.Len(t, run.TaskRuns, 1)
		assert.Equal(t, models.RunStatusErrored, run.TaskRuns[0].Status)
	})

	t.Run("completed input with remaining tasks should put task into in-progress", func(t *testing.T) {
		runID := models.NewID()
		run := models.JobRun{
			ID:        runID,
			JobSpecID: job.ID,
			Status:    models.RunStatusPendingBridge,
			TaskRuns:  []models.TaskRun{models.TaskRun{ID: models.NewID(), JobRunID: runID}, models.TaskRun{ID: models.NewID(), JobRunID: runID}},
		}
		require.NoError(t, store.CreateJobRun(&run))

		err := runManager.ResumePending(run.ID, models.BridgeRunResult{Data: input, Status: models.RunStatusCompleted})
		assert.NoError(t, err)

		run, err = store.FindJobRun(run.ID)
		require.NoError(t, err)
		assert.NoError(t, err)
		assert.Equal(t, string(models.RunStatusInProgress), string(run.Status))
		assert.Len(t, run.TaskRuns, 2)
		assert.Equal(t, string(models.RunStatusCompleted), string(run.TaskRuns[0].Status))
	})

	t.Run("completed input with no remaining tasks should get marked as complete", func(t *testing.T) {
		runID := models.NewID()
		run := models.JobRun{
			ID:        runID,
			JobSpecID: job.ID,
			Status:    models.RunStatusPendingBridge,
			TaskRuns:  []models.TaskRun{models.TaskRun{ID: models.NewID(), JobRunID: runID}},
		}
		require.NoError(t, store.CreateJobRun(&run))

		err := runManager.ResumePending(run.ID, models.BridgeRunResult{Data: input, Status: models.RunStatusCompleted})
		assert.NoError(t, err)

		run, err = store.FindJobRun(run.ID)
		require.NoError(t, err)
		assert.Equal(t, string(models.RunStatusCompleted), string(run.Status))
		assert.True(t, run.FinishedAt.Valid)
		assert.Len(t, run.TaskRuns, 1)
		assert.Equal(t, string(models.RunStatusCompleted), string(run.TaskRuns[0].Status))
	})

	runQueue.AssertExpectations(t)
}

func TestRunManager_ResumeAllConfirming(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	runQueue := new(mocks.RunQueue)
	runQueue.On("Run", mock.Anything).Maybe().Return(nil)

	runManager := services.NewRunManager(runQueue, store.Config, store.ORM, store.TxManager, store.Clock)

	job := cltest.NewJob()
	require.NoError(t, store.CreateJob(&job))

	t.Run("reject a run with no tasks", func(t *testing.T) {
		run := models.JobRun{
			ID:        models.NewID(),
			JobSpecID: job.ID,
			Status:    models.RunStatusPendingConfirmations,
		}
		require.NoError(t, store.CreateJobRun(&run))

		err := runManager.ResumeAllConfirming(nil)
		assert.NoError(t, err)

		run, err = store.FindJobRun(run.ID)
		require.NoError(t, err)
		assert.Equal(t, models.RunStatusErrored, run.Status)
	})

	creationHeight := utils.NewBig(big.NewInt(0))

	t.Run("leave in pending if not enough confirmations have been met yet", func(t *testing.T) {
		run := models.JobRun{
			ID:             models.NewID(),
			JobSpecID:      job.ID,
			CreationHeight: creationHeight,
			Status:         models.RunStatusPendingConfirmations,
			TaskRuns: []models.TaskRun{models.TaskRun{
				ID:                   models.NewID(),
				MinimumConfirmations: clnull.Uint32From(2),
				TaskSpec: models.TaskSpec{
					Type: adapters.TaskTypeNoOp,
				},
			}},
		}
		require.NoError(t, store.CreateJobRun(&run))

		err := runManager.ResumeAllConfirming(creationHeight.ToInt())
		require.NoError(t, err)

		run, err = store.FindJobRun(run.ID)
		require.NoError(t, err)
		assert.Equal(t, models.RunStatusPendingConfirmations, run.Status)
		assert.Equal(t, uint32(1), run.TaskRuns[0].Confirmations.Uint32)
	})

	t.Run("input, should go from pending -> in progress and save the input", func(t *testing.T) {
		run := models.JobRun{
			ID:             models.NewID(),
			JobSpecID:      job.ID,
			CreationHeight: creationHeight,
			Status:         models.RunStatusPendingConfirmations,
			TaskRuns: []models.TaskRun{models.TaskRun{
				ID:                   models.NewID(),
				MinimumConfirmations: clnull.Uint32From(1),
				TaskSpec: models.TaskSpec{
					Type: adapters.TaskTypeNoOp,
				},
			}},
		}
		require.NoError(t, store.CreateJobRun(&run))

		observedHeight := big.NewInt(1)
		err := runManager.ResumeAllConfirming(observedHeight)
		require.NoError(t, err)

		run, err = store.FindJobRun(run.ID)
		require.NoError(t, err)
		assert.Equal(t, string(models.RunStatusInProgress), string(run.Status))
	})

	runQueue.AssertExpectations(t)
}

func TestRunManager_ResumeAllConnecting(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	runQueue := new(mocks.RunQueue)
	runQueue.On("Run", mock.Anything).Maybe().Return(nil)

	runManager := services.NewRunManager(runQueue, store.Config, store.ORM, store.TxManager, store.Clock)

	job := cltest.NewJob()
	require.NoError(t, store.CreateJob(&job))

	t.Run("reject a run with no tasks", func(t *testing.T) {
		run := models.JobRun{
			ID:        models.NewID(),
			JobSpecID: job.ID,
			Status:    models.RunStatusPendingConnection,
		}
		require.NoError(t, store.CreateJobRun(&run))

		err := runManager.ResumeAllConnecting()
		assert.NoError(t, err)

		run, err = store.FindJobRun(run.ID)
		require.NoError(t, err)
		assert.Equal(t, models.RunStatusErrored, run.Status)
	})

	t.Run("input, should go from pending -> in progress", func(t *testing.T) {
		run := models.JobRun{
			ID:        models.NewID(),
			JobSpecID: job.ID,
			Status:    models.RunStatusPendingConnection,
			TaskRuns: []models.TaskRun{models.TaskRun{
				ID: models.NewID(),
			}},
		}
		require.NoError(t, store.CreateJobRun(&run))
		err := runManager.ResumeAllConnecting()
		assert.NoError(t, err)

		run, err = store.FindJobRun(run.ID)
		require.NoError(t, err)
		assert.Equal(t, models.RunStatusInProgress, run.Status)
	})
}

func TestRunManager_ResumeAllConnecting_NotEnoughConfirmations(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t)
	defer cleanup()

	store := app.Store
	eth := cltest.MockEthOnStore(t, store)
	eth.Register("eth_chainId", store.Config.ChainID())

	app.Start()

	job := cltest.NewJobWithRunLogInitiator()
	job.Tasks = []models.TaskSpec{cltest.NewTask(t, "NoOp")}
	require.NoError(t, store.CreateJob(&job))

	initiator := job.Initiators[0]
	run := job.NewRun(initiator)
	run.Status = models.RunStatusPendingConnection
	run.CreationHeight = utils.NewBig(big.NewInt(0))
	run.ObservedHeight = run.CreationHeight
	run.TaskRuns[0].MinimumConfirmations = clnull.Uint32From(807)
	run.TaskRuns[0].Status = models.RunStatusPendingConnection
	require.NoError(t, store.CreateJobRun(&run))

	app.RunManager.ResumeAllConnecting()

	cltest.WaitForJobRunToPendConfirmations(t, store, run)
}

func TestRunManager_Create(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t)
	defer cleanup()

	store := app.Store
	eth := cltest.MockEthOnStore(t, store)
	eth.Register("eth_chainId", store.Config.ChainID())

	app.Start()

	job := cltest.NewJobWithRunLogInitiator()
	job.Tasks = []models.TaskSpec{cltest.NewTask(t, "NoOp")} // empty params
	require.NoError(t, store.CreateJob(&job))

	requestID := "RequestID"
	initr := job.Initiators[0]
	rr := models.NewRunRequest()
	rr.RequestID = &requestID
	data := cltest.JSONFromString(t, `{"random": "input"}`)
	jr, err := app.RunManager.Create(job.ID, &initr, &data, nil, rr)
	require.NoError(t, err)
	updatedJR := cltest.WaitForJobRunToComplete(t, store, *jr)
	assert.Equal(t, rr.RequestID, updatedJR.RunRequest.RequestID)
}

func TestRunManager_Create_DoesNotSaveToTaskSpec(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t)
	defer cleanup()

	store := app.Store
	mocketh := cltest.MockEthOnStore(t, store)
	mocketh.Register("eth_chainId", store.Config.ChainID())

	app.Start()

	job := cltest.NewJobWithWebInitiator()
	job.Tasks = []models.TaskSpec{cltest.NewTask(t, "NoOp")} // empty params
	require.NoError(t, store.CreateJob(&job))

	initr := job.Initiators[0]
	data := cltest.JSONFromString(t, `{"random": "input"}`)
	jr, err := app.RunManager.Create(job.ID, &initr, &data, nil, &models.RunRequest{})
	require.NoError(t, err)
	cltest.WaitForJobRunToComplete(t, store, *jr)

	retrievedJob, err := store.FindJob(job.ID)
	require.NoError(t, err)
	require.Len(t, job.Tasks, 1)
	require.Len(t, retrievedJob.Tasks, 1)
	assert.Equal(t, job.Tasks[0].Params, retrievedJob.Tasks[0].Params)
}

func TestRunManager_Create_fromRunLog_Happy(t *testing.T) {
	t.Parallel()

	initiatingTxHash := cltest.NewHash()
	triggeringBlockHash := cltest.NewHash()
	otherBlockHash := cltest.NewHash()

	tests := []struct {
		name             string
		logBlockHash     common.Hash
		receiptBlockHash common.Hash
		wantStatus       models.RunStatus
	}{
		{
			name:             "main chain",
			logBlockHash:     triggeringBlockHash,
			receiptBlockHash: triggeringBlockHash,
			wantStatus:       models.RunStatusCompleted,
		},
		{
			name:             "ommered chain",
			logBlockHash:     triggeringBlockHash,
			receiptBlockHash: otherBlockHash,
			wantStatus:       models.RunStatusErrored,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			config, cfgCleanup := cltest.NewConfig(t)
			defer cfgCleanup()
			minimumConfirmations := uint32(2)
			config.Set("MIN_INCOMING_CONFIRMATIONS", minimumConfirmations)
			app, cleanup := cltest.NewApplicationWithConfig(t, config)
			defer cleanup()

			mocketh := app.MockEthCallerSubscriber()
			store := app.GetStore()
			mocketh.Context("app.Start()", func(meth *cltest.EthMock) {
				meth.Register("eth_chainId", store.Config.ChainID())
			})
			app.Start()

			job := cltest.NewJobWithRunLogInitiator()
			job.Tasks = []models.TaskSpec{cltest.NewTask(t, "NoOp")}
			require.NoError(t, store.CreateJob(&job))

			creationHeight := big.NewInt(1)
			requestID := "RequestID"
			initr := job.Initiators[0]
			rr := models.NewRunRequest()
			rr.RequestID = &requestID
			rr.TxHash = &initiatingTxHash
			rr.BlockHash = &test.logBlockHash
			data := cltest.JSONFromString(t, `{"random": "input"}`)
			jr, err := app.RunManager.Create(job.ID, &initr, &data, creationHeight, rr)
			require.NoError(t, err)

			run := cltest.WaitForJobRunToPendConfirmations(t, app.Store, *jr)
			assert.Equal(t, models.RunStatusPendingConfirmations, run.TaskRuns[0].Status)
			assert.Equal(t, models.RunStatusPendingConfirmations, run.Status)

			confirmedReceipt := eth.TxReceipt{
				Hash:        initiatingTxHash,
				BlockHash:   &test.receiptBlockHash,
				BlockNumber: cltest.Int(3),
			}
			mocketh.Context("validateOnMainChain", func(meth *cltest.EthMock) {
				meth.Register("eth_getTransactionReceipt", confirmedReceipt)
			})

			err = app.RunManager.ResumeAllConfirming(big.NewInt(2))
			require.NoError(t, err)
			run = cltest.WaitForJobRunStatus(t, store, *jr, test.wantStatus)
			assert.Equal(t, rr.RequestID, run.RunRequest.RequestID)
			assert.Equal(t, minimumConfirmations, run.TaskRuns[0].MinimumConfirmations.Uint32)
			assert.True(t, run.TaskRuns[0].MinimumConfirmations.Valid)
			assert.Equal(t, minimumConfirmations, run.TaskRuns[0].Confirmations.Uint32, "task run should track its current confirmations")
			assert.True(t, run.TaskRuns[0].Confirmations.Valid)

			assert.True(t, mocketh.AllCalled(), mocketh.Remaining())
		})
	}
}

func TestRunManager_Create_fromRunLog_ConnectToLaggingEthNode(t *testing.T) {
	t.Parallel()

	initiatingTxHash := cltest.NewHash()
	triggeringBlockHash := cltest.NewHash()

	config, cfgCleanup := cltest.NewConfig(t)
	defer cfgCleanup()
	minimumConfirmations := uint32(2)
	config.Set("MIN_INCOMING_CONFIRMATIONS", minimumConfirmations)
	app, cleanup := cltest.NewApplicationWithConfig(t, config)
	defer cleanup()

	eth := app.MockEthCallerSubscriber()
	app.MockStartAndConnect()

	store := app.GetStore()
	job := cltest.NewJobWithRunLogInitiator()
	job.Tasks = []models.TaskSpec{cltest.NewTask(t, "NoOp")}
	require.NoError(t, store.CreateJob(&job))

	requestID := "RequestID"
	initr := job.Initiators[0]
	rr := models.NewRunRequest()
	rr.RequestID = &requestID
	rr.TxHash = &initiatingTxHash
	rr.BlockHash = &triggeringBlockHash

	futureCreationHeight := big.NewInt(9)
	pastCurrentHeight := big.NewInt(1)

	data := cltest.JSONFromString(t, `{"random": "input"}`)
	jr, err := app.RunManager.Create(job.ID, &initr, &data, futureCreationHeight, rr)
	require.NoError(t, err)
	cltest.WaitForJobRunToPendConfirmations(t, app.Store, *jr)

	err = app.RunManager.ResumeAllConfirming(pastCurrentHeight)
	require.NoError(t, err)
	updatedJR := cltest.WaitForJobRunToPendConfirmations(t, app.Store, *jr)
	assert.True(t, updatedJR.TaskRuns[0].Confirmations.Valid)
	assert.Equal(t, uint32(0), updatedJR.TaskRuns[0].Confirmations.Uint32)
	assert.True(t, eth.AllCalled(), eth.Remaining())
}

func TestRunManager_ResumeConfirmingTasks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status models.RunStatus
	}{
		{models.RunStatusPendingConnection},
		{models.RunStatusPendingConfirmations},
	}

	for _, test := range tests {
		t.Run(string(test.status), func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			job := cltest.NewJobWithWebInitiator()
			require.NoError(t, store.CreateJob(&job))
			initr := job.Initiators[0]
			run := job.NewRun(initr)
			run.Status = test.status
			require.NoError(t, store.CreateJobRun(&run))

			runQueue := new(mocks.RunQueue)
			runQueue.On("Run", mock.Anything).Return(nil)

			runManager := services.NewRunManager(runQueue, store.Config, store.ORM, store.TxManager, store.Clock)
			runManager.ResumeAllConfirming(big.NewInt(3821))

			runQueue.AssertExpectations(t)
		})
	}
}

func TestRunManager_ResumeAllInProgress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status models.RunStatus
	}{
		{models.RunStatusInProgress},
		{models.RunStatusPendingSleep},
	}

	for _, test := range tests {
		t.Run(string(test.status), func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			job := cltest.NewJobWithWebInitiator()
			require.NoError(t, store.CreateJob(&job))
			initr := job.Initiators[0]
			run := job.NewRun(initr)
			run.Status = test.status
			require.NoError(t, store.CreateJobRun(&run))

			runQueue := new(mocks.RunQueue)
			runQueue.On("Run", mock.Anything).Return(nil)

			runManager := services.NewRunManager(runQueue, store.Config, store.ORM, store.TxManager, store.Clock)
			runManager.ResumeAllInProgress()

			runQueue.AssertExpectations(t)
		})
	}
}

// XXX: In progress tasks that are archived should still be run as they have been paid for
func TestRunManager_ResumeAllInProgress_Archived(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status models.RunStatus
	}{
		{models.RunStatusInProgress},
		{models.RunStatusPendingSleep},
	}

	for _, test := range tests {
		t.Run(string(test.status), func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			job := cltest.NewJobWithWebInitiator()
			require.NoError(t, store.CreateJob(&job))
			initr := job.Initiators[0]
			run := job.NewRun(initr)
			run.Status = test.status
			run.DeletedAt = null.TimeFrom(time.Now())
			require.NoError(t, store.CreateJobRun(&run))

			runQueue := new(mocks.RunQueue)
			runQueue.On("Run", mock.Anything).Return(nil)

			runManager := services.NewRunManager(runQueue, store.Config, store.ORM, store.TxManager, store.Clock)
			runManager.ResumeAllInProgress()

			runQueue.AssertExpectations(t)
		})
	}
}

func TestRunManager_ResumeAllInProgress_NotInProgress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status models.RunStatus
	}{
		{models.RunStatusPendingConnection},
		{models.RunStatusPendingConfirmations},
		{models.RunStatusPendingBridge},
		{models.RunStatusCompleted},
		{models.RunStatusCancelled},
	}

	for _, test := range tests {
		t.Run(string(test.status), func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			job := cltest.NewJobWithWebInitiator()
			require.NoError(t, store.CreateJob(&job))
			initr := job.Initiators[0]
			run := job.NewRun(initr)
			run.Status = test.status
			require.NoError(t, store.CreateJobRun(&run))

			runQueue := new(mocks.RunQueue)
			runQueue.On("Run", mock.Anything).Maybe().Return(nil)

			runManager := services.NewRunManager(runQueue, store.Config, store.ORM, store.TxManager, store.Clock)
			runManager.ResumeAllInProgress()

			runQueue.AssertExpectations(t)
		})
	}
}

func TestRunManager_ResumeAllInProgress_NotInProgressAndArchived(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status models.RunStatus
	}{
		{models.RunStatusPendingConnection},
		{models.RunStatusPendingConfirmations},
		{models.RunStatusPendingBridge},
		{models.RunStatusCompleted},
		{models.RunStatusCancelled},
	}

	for _, test := range tests {
		t.Run(string(test.status), func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			job := cltest.NewJobWithWebInitiator()
			require.NoError(t, store.CreateJob(&job))
			initr := job.Initiators[0]
			run := job.NewRun(initr)
			run.Status = test.status
			run.DeletedAt = null.TimeFrom(time.Now())
			require.NoError(t, store.CreateJobRun(&run))

			runQueue := new(mocks.RunQueue)
			runQueue.On("Run", mock.Anything).Maybe().Return(nil)

			runManager := services.NewRunManager(runQueue, store.Config, store.ORM, store.TxManager, store.Clock)
			runManager.ResumeAllInProgress()

			runQueue.AssertExpectations(t)
		})
	}
}
