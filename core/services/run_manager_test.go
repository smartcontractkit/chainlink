package services_test

import (
	"fmt"

	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	clnull "github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"
)

func makeJobRunWithInitiator(t *testing.T, store *strpkg.Store, job models.JobSpec) models.JobRun {
	require.NoError(t, store.CreateJob(&job))

	initr := models.Initiator{
		JobSpecID: job.ID,
	}

	err := store.CreateInitiator(&initr)
	require.NoError(t, err)

	return models.MakeJobRun(&job, time.Now(), &initr, big.NewInt(0), &models.RunRequest{})
}

func TestRunManager_ResumePendingBridge(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	pusher := new(mocks.StatsPusher)
	pusher.On("PushNow").Return(nil)

	runQueue := new(mocks.RunQueue)
	runQueue.On("Run", mock.Anything).Maybe().Return(nil)

	runManager := services.NewRunManager(runQueue, store.Config, store.ORM, pusher, store.TxManager, store.Clock)

	input := cltest.JSONFromString(t, `{"address":"0xdfcfc2b9200dbb10952c2b7cce60fc7260e03c6f"}`)

	t.Run("reject a run with an invalid state", func(t *testing.T) {
		run := makeJobRunWithInitiator(t, store, cltest.NewJob())
		require.NoError(t, store.CreateJobRun(&run))
		err := runManager.ResumePendingBridge(run.ID, models.BridgeRunResult{})
		assert.Error(t, err)
	})

	t.Run("reject a run with no tasks", func(t *testing.T) {
		run := makeJobRunWithInitiator(t, store, models.NewJob())
		run.SetStatus(models.RunStatusPendingBridge)
		require.NoError(t, store.CreateJobRun(&run))
		err := runManager.ResumePendingBridge(run.ID, models.BridgeRunResult{})
		assert.NoError(t, err)

		run, err = store.FindJobRun(run.ID)
		require.NoError(t, err)
		assert.Equal(t, models.RunStatusErrored, run.GetStatus())
	})

	t.Run("input with error errors run", func(t *testing.T) {
		run := makeJobRunWithInitiator(t, store, cltest.NewJob())
		run.SetStatus(models.RunStatusPendingBridge)
		require.NoError(t, store.CreateJobRun(&run))

		err := runManager.ResumePendingBridge(run.ID, models.BridgeRunResult{Status: models.RunStatusErrored})
		assert.NoError(t, err)

		run, err = store.FindJobRun(run.ID)
		require.NoError(t, err)
		assert.Equal(t, models.RunStatusErrored, run.GetStatus())
		assert.True(t, run.FinishedAt.Valid)
		assert.Len(t, run.TaskRuns, 1)
		assert.Equal(t, models.RunStatusErrored, run.TaskRuns[0].Status)
	})

	t.Run("completed input with remaining tasks should put task into in-progress", func(t *testing.T) {
		job := cltest.NewJob()
		job.Tasks = []models.TaskSpec{{Type: adapters.TaskTypeNoOp}, {Type: adapters.TaskTypeNoOp}}
		run := makeJobRunWithInitiator(t, store, job)
		run.SetStatus(models.RunStatusPendingBridge)
		require.NoError(t, store.CreateJobRun(&run))

		err := runManager.ResumePendingBridge(run.ID, models.BridgeRunResult{Data: input, Status: models.RunStatusCompleted})
		assert.NoError(t, err)

		run, err = store.FindJobRun(run.ID)
		require.NoError(t, err)
		assert.NoError(t, err)
		assert.Equal(t, string(models.RunStatusInProgress), string(run.GetStatus()))
		assert.Len(t, run.TaskRuns, 2)
		assert.Equal(t, string(models.RunStatusCompleted), string(run.TaskRuns[0].Status))
	})

	t.Run("completed input with no remaining tasks should get marked as complete", func(t *testing.T) {
		run := makeJobRunWithInitiator(t, store, cltest.NewJob())
		run.SetStatus(models.RunStatusPendingBridge)
		require.NoError(t, store.CreateJobRun(&run))

		err := runManager.ResumePendingBridge(run.ID, models.BridgeRunResult{Data: input, Status: models.RunStatusCompleted})
		assert.NoError(t, err)

		run, err = store.FindJobRun(run.ID)
		require.NoError(t, err)
		assert.Equal(t, string(models.RunStatusCompleted), string(run.GetStatus()))
		assert.True(t, run.FinishedAt.Valid)
		assert.Len(t, run.TaskRuns, 1)
		assert.Equal(t, string(models.RunStatusCompleted), string(run.TaskRuns[0].Status))
	})

	runQueue.AssertExpectations(t)
}

func TestRunManager_ResumeAllPendingNextBlock(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	pusher := new(mocks.StatsPusher)
	pusher.On("PushNow").Return(nil)

	runQueue := new(mocks.RunQueue)
	runQueue.On("Run", mock.Anything).Maybe().Return(nil)

	runManager := services.NewRunManager(runQueue, store.Config, store.ORM, pusher, store.TxManager, store.Clock)

	t.Run("reject a run with no tasks", func(t *testing.T) {
		run := makeJobRunWithInitiator(t, store, models.NewJob())
		run.SetStatus(models.RunStatusPendingIncomingConfirmations)
		require.NoError(t, store.CreateJobRun(&run))

		err := runManager.ResumeAllPendingNextBlock(nil)
		assert.NoError(t, err)

		run, err = store.FindJobRun(run.ID)
		require.NoError(t, err)
		assert.Equal(t, models.RunStatusErrored, run.GetStatus())
	})

	t.Run("leave in pending if not enough incoming confirmations have been met yet", func(t *testing.T) {
		run := makeJobRunWithInitiator(t, store, cltest.NewJob())
		run.SetStatus(models.RunStatusPendingIncomingConfirmations)
		run.TaskRuns[0].MinRequiredIncomingConfirmations = clnull.Uint32From(2)
		require.NoError(t, store.CreateJobRun(&run))

		err := runManager.ResumeAllPendingNextBlock(big.NewInt(0))
		require.NoError(t, err)

		run, err = store.FindJobRun(run.ID)
		require.NoError(t, err)
		assert.Equal(t, models.RunStatusPendingIncomingConfirmations, run.GetStatus())
		assert.Equal(t, uint32(1), run.TaskRuns[0].ObservedIncomingConfirmations.Uint32)
	})

	t.Run("input, should go from pending_incoming_confirmations -> in_progress and save the input", func(t *testing.T) {
		run := makeJobRunWithInitiator(t, store, cltest.NewJob())
		run.SetStatus(models.RunStatusPendingIncomingConfirmations)
		run.TaskRuns[0].MinRequiredIncomingConfirmations = clnull.Uint32From(2)
		require.NoError(t, store.CreateJobRun(&run))

		observedHeight := big.NewInt(1)
		err := runManager.ResumeAllPendingNextBlock(observedHeight)
		require.NoError(t, err)

		run, err = store.FindJobRun(run.ID)
		require.NoError(t, err)
		assert.Equal(t, string(models.RunStatusInProgress), string(run.GetStatus()))
	})

	runQueue.AssertExpectations(t)
}

func TestRunManager_ResumeAllPendingConnection(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	pusher := new(mocks.StatsPusher)
	pusher.On("PushNow").Return(nil)

	runQueue := new(mocks.RunQueue)
	runQueue.On("Run", mock.Anything).Maybe().Return(nil)

	runManager := services.NewRunManager(runQueue, store.Config, store.ORM, pusher, store.TxManager, store.Clock)

	t.Run("reject a run with no tasks", func(t *testing.T) {
		run := makeJobRunWithInitiator(t, store, models.NewJob())
		run.SetStatus(models.RunStatusPendingConnection)
		require.NoError(t, store.CreateJobRun(&run))

		err := runManager.ResumeAllPendingConnection()
		assert.NoError(t, err)

		run, err = store.FindJobRun(run.ID)
		require.NoError(t, err)
		assert.Equal(t, models.RunStatusErrored, run.GetStatus())
	})

	t.Run("input, should go from pending -> in progress", func(t *testing.T) {
		run := makeJobRunWithInitiator(t, store, cltest.NewJob())
		run.SetStatus(models.RunStatusPendingConnection)

		job, err := store.FindJob(run.JobSpecID)
		require.NoError(t, err)
		run.TaskRuns = []models.TaskRun{models.TaskRun{ID: models.NewID(), TaskSpecID: job.Tasks[0].ID, Status: models.RunStatusUnstarted}}

		require.NoError(t, store.CreateJobRun(&run))
		err = runManager.ResumeAllPendingConnection()
		assert.NoError(t, err)

		run, err = store.FindJobRun(run.ID)
		require.NoError(t, err)
		assert.Equal(t, models.RunStatusInProgress, run.GetStatus())
	})
}

func TestRunManager_ResumeAllPendingConnection_NotEnoughConfirmations(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()

	store := app.Store

	app.StartAndConnect()

	job := cltest.NewJobWithRunLogInitiator()
	job.Tasks = []models.TaskSpec{cltest.NewTask(t, "NoOp")}
	require.NoError(t, store.CreateJob(&job))

	run := cltest.NewJobRun(job)
	run.SetStatus(models.RunStatusPendingConnection)
	run.CreationHeight = utils.NewBig(big.NewInt(0))
	run.ObservedHeight = run.CreationHeight
	run.TaskRuns[0].MinRequiredIncomingConfirmations = clnull.Uint32From(807)
	run.TaskRuns[0].Status = models.RunStatusPendingConnection
	require.NoError(t, store.CreateJobRun(&run))

	app.RunManager.ResumeAllPendingConnection()

	cltest.WaitForJobRunToPendIncomingConfirmations(t, store, run)
}

func TestRunManager_Create(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()

	store := app.Store

	app.StartAndConnect()

	job := cltest.NewJobWithRunLogInitiator()
	job.Tasks = []models.TaskSpec{cltest.NewTask(t, "NoOp")} // empty params
	require.NoError(t, store.CreateJob(&job))

	requestID := common.HexToHash("0xcafe")
	initiator := job.Initiators[0]
	rr := models.NewRunRequest(models.JSON{})
	rr.RequestID = &requestID
	rr.RequestParams = cltest.JSONFromString(t, `{"random": "input"}`)
	jr, err := app.RunManager.Create(job.ID, &initiator, nil, rr)
	require.NoError(t, err)
	updatedJR := cltest.WaitForJobRunToComplete(t, store, *jr)
	assert.Equal(t, rr.RequestID, updatedJR.RunRequest.RequestID)
}

func TestRunManager_Create_DoesNotSaveToTaskSpec(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()

	store := app.Store

	app.StartAndConnect()

	job := cltest.NewJobWithWebInitiator()
	job.Tasks = []models.TaskSpec{cltest.NewTask(t, "NoOp")} // empty params
	require.NoError(t, store.CreateJob(&job))

	initiator := job.Initiators[0]
	data := cltest.JSONFromString(t, `{"random": "input"}`)
	rr := &models.RunRequest{RequestParams: data}
	jr, err := app.RunManager.Create(job.ID, &initiator, nil, rr)
	require.NoError(t, err)
	cltest.WaitForJobRunToComplete(t, store, *jr)

	retrievedJob, err := store.FindJob(job.ID)
	require.NoError(t, err)
	require.Len(t, job.Tasks, 1)
	require.Len(t, retrievedJob.Tasks, 1)
	assert.JSONEq(t, job.Tasks[0].Params.String(), retrievedJob.Tasks[0].Params.String())
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

			gethClient := new(mocks.GethClient)
			sub := new(mocks.Subscription)
			app, cleanup := cltest.NewApplicationWithConfig(t, config,
				eth.NewClientWith(nil, gethClient),
			)
			gethClient.On("ChainID", mock.Anything).Return(app.Config.ChainID(), nil)
			gethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(1), nil)
			gethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).Return(sub, nil)
			gethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).Return(sub, nil)
			sub.On("Err").Return(nil)
			sub.On("Unsubscribe").Return()

			kst := new(mocks.KeyStoreInterface)
			kst.On("Accounts").Return([]accounts.Account{})
			app.Store.KeyStore = kst
			defer cleanup()

			app.StartAndConnect()

			job := cltest.NewJobWithRunLogInitiator()
			job.Tasks = []models.TaskSpec{cltest.NewTask(t, "NoOp")}
			require.NoError(t, app.Store.CreateJob(&job))

			creationHeight := big.NewInt(1)
			requestID := common.HexToHash("0xcafe")
			initiator := job.Initiators[0]
			rr := models.NewRunRequest(models.JSON{})
			rr.RequestID = &requestID
			rr.TxHash = &initiatingTxHash
			rr.BlockHash = &test.logBlockHash
			rr.RequestParams = cltest.JSONFromString(t, `{"random": "input"}`)
			jr, err := app.RunManager.Create(job.ID, &initiator, creationHeight, rr)
			require.NoError(t, err)

			run := cltest.WaitForJobRunToPendIncomingConfirmations(t, app.Store, *jr)
			assert.Equal(t, models.RunStatusPendingIncomingConfirmations, run.TaskRuns[0].Status)
			assert.Equal(t, models.RunStatusPendingIncomingConfirmations, run.GetStatus())

			gethClient.On("TransactionReceipt", mock.Anything, mock.Anything).Return(&types.Receipt{
				TxHash:      initiatingTxHash,
				BlockHash:   test.receiptBlockHash,
				BlockNumber: big.NewInt(3),
			}, nil)

			err = app.RunManager.ResumeAllPendingNextBlock(big.NewInt(2))
			require.NoError(t, err)
			run = cltest.WaitForJobRunStatus(t, app.Store, *jr, test.wantStatus)
			assert.Equal(t, rr.RequestID, run.RunRequest.RequestID)
			assert.Equal(t, minimumConfirmations, run.TaskRuns[0].MinRequiredIncomingConfirmations.Uint32)
			assert.True(t, run.TaskRuns[0].MinRequiredIncomingConfirmations.Valid)
			assert.Equal(t, minimumConfirmations, run.TaskRuns[0].ObservedIncomingConfirmations.Uint32, "task run should track its current confirmations")
			assert.True(t, run.TaskRuns[0].ObservedIncomingConfirmations.Valid)

			assert.True(t, app.EthMock.AllCalled(), app.EthMock.Remaining())

			kst.AssertExpectations(t)
		})
	}
}

func TestRunManager_Create_fromRunLogPayments(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		inputPayment         *assets.Link
		jobMinimumPayment    *assets.Link
		configMinimumPayment string
		bridgePayment        *assets.Link
		jobStatus            models.RunStatus
	}{
		// no payments required
		{
			name:                 "no payment required and none given",
			inputPayment:         assets.NewLink(0),
			jobMinimumPayment:    nil,
			configMinimumPayment: "0",
			bridgePayment:        assets.NewLink(0),
			jobStatus:            models.RunStatusInProgress,
		},
		{
			name:                 "no payment required and some given",
			inputPayment:         assets.NewLink(13),
			jobMinimumPayment:    nil,
			configMinimumPayment: "0",
			bridgePayment:        assets.NewLink(0),
			jobStatus:            models.RunStatusInProgress,
		},

		// configuration payments only
		{
			name:                 "configuration payment required and none given",
			inputPayment:         assets.NewLink(0),
			jobMinimumPayment:    nil,
			configMinimumPayment: "13",
			bridgePayment:        assets.NewLink(0),
			jobStatus:            models.RunStatusErrored,
		},
		{
			name:                 "configuration payment required and insufficient given",
			inputPayment:         assets.NewLink(7),
			jobMinimumPayment:    nil,
			configMinimumPayment: "13",
			bridgePayment:        assets.NewLink(0),
			jobStatus:            models.RunStatusErrored,
		},
		{
			name:                 "configuration payment required and exact amount given",
			inputPayment:         assets.NewLink(13),
			jobMinimumPayment:    nil,
			configMinimumPayment: "13",
			bridgePayment:        assets.NewLink(0),
			jobStatus:            models.RunStatusInProgress,
		},
		{
			name:                 "configuration payment required and excess amount given",
			inputPayment:         assets.NewLink(17),
			jobMinimumPayment:    nil,
			configMinimumPayment: "13",
			bridgePayment:        assets.NewLink(0),
			jobStatus:            models.RunStatusInProgress,
		},

		// job payments only
		{
			name:                 "job payment required and none given",
			inputPayment:         assets.NewLink(0),
			jobMinimumPayment:    assets.NewLink(13),
			configMinimumPayment: "0",
			bridgePayment:        assets.NewLink(0),
			jobStatus:            models.RunStatusErrored,
		},
		{
			name:                 "job payment required and insufficient given",
			inputPayment:         assets.NewLink(7),
			jobMinimumPayment:    assets.NewLink(13),
			configMinimumPayment: "0",
			bridgePayment:        assets.NewLink(0),
			jobStatus:            models.RunStatusErrored,
		},
		{
			name:                 "job payment required and exact amount given",
			inputPayment:         assets.NewLink(13),
			jobMinimumPayment:    assets.NewLink(13),
			configMinimumPayment: "0",
			bridgePayment:        assets.NewLink(0),
			jobStatus:            models.RunStatusInProgress,
		},
		{
			name:                 "job payment required and excess amount given",
			inputPayment:         assets.NewLink(17),
			jobMinimumPayment:    assets.NewLink(13),
			configMinimumPayment: "0",
			bridgePayment:        assets.NewLink(0),
			jobStatus:            models.RunStatusInProgress,
		},

		// bridge payments only
		{
			name:                 "bridge payment required and none given",
			inputPayment:         assets.NewLink(0),
			jobMinimumPayment:    assets.NewLink(0),
			configMinimumPayment: "0",
			bridgePayment:        assets.NewLink(13),
			jobStatus:            models.RunStatusErrored,
		},
		{
			name:                 "bridge payment required and insufficient given",
			inputPayment:         assets.NewLink(7),
			jobMinimumPayment:    assets.NewLink(0),
			configMinimumPayment: "0",
			bridgePayment:        assets.NewLink(13),
			jobStatus:            models.RunStatusErrored,
		},
		{
			name:                 "bridge payment required and exact amount given",
			inputPayment:         assets.NewLink(13),
			jobMinimumPayment:    assets.NewLink(0),
			configMinimumPayment: "0",
			bridgePayment:        assets.NewLink(13),
			jobStatus:            models.RunStatusInProgress,
		},
		{
			name:                 "bridge payment required and excess amount given",
			inputPayment:         assets.NewLink(17),
			jobMinimumPayment:    assets.NewLink(0),
			configMinimumPayment: "0",
			bridgePayment:        assets.NewLink(13),
			jobStatus:            models.RunStatusInProgress,
		},

		// job and bridge payments
		{
			name:                 "job and bridge payment required and none given",
			inputPayment:         assets.NewLink(0),
			jobMinimumPayment:    assets.NewLink(11),
			configMinimumPayment: "0",
			bridgePayment:        assets.NewLink(13),
			jobStatus:            models.RunStatusErrored,
		},
		{
			name:                 "job and bridge payment required and insufficient given",
			inputPayment:         assets.NewLink(11),
			jobMinimumPayment:    assets.NewLink(11),
			configMinimumPayment: "0",
			bridgePayment:        assets.NewLink(13),
			jobStatus:            models.RunStatusErrored,
		},
		{
			name:                 "job and bridge payment required and exact amount given",
			inputPayment:         assets.NewLink(24),
			jobMinimumPayment:    assets.NewLink(11),
			configMinimumPayment: "0",
			bridgePayment:        assets.NewLink(13),
			jobStatus:            models.RunStatusInProgress,
		},
		{
			name:                 "job and bridge payment required and excess amount given",
			inputPayment:         assets.NewLink(25),
			jobMinimumPayment:    assets.NewLink(11),
			configMinimumPayment: "0",
			bridgePayment:        assets.NewLink(13),
			jobStatus:            models.RunStatusInProgress,
		},

		// config and job payments (uses job minimum payment)
		{
			name:                 "both payments required and no payment given",
			inputPayment:         assets.NewLink(0),
			jobMinimumPayment:    assets.NewLink(11),
			configMinimumPayment: "13",
			bridgePayment:        assets.NewLink(0),
			jobStatus:            models.RunStatusErrored,
		},
		{
			name:                 "both payments required and job payment amount given",
			inputPayment:         assets.NewLink(11),
			jobMinimumPayment:    assets.NewLink(11),
			configMinimumPayment: "13",
			bridgePayment:        assets.NewLink(0),
			jobStatus:            models.RunStatusInProgress,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			config, configCleanup := cltest.NewConfig(t)
			defer configCleanup()
			config.Set("DATABASE_TIMEOUT", "10s") // Lots of parallelized tests
			config.Set("MINIMUM_CONTRACT_PAYMENT", test.configMinimumPayment)

			store, storeCleanup := cltest.NewStoreWithConfig(config)
			defer storeCleanup()

			bt := &models.BridgeType{
				Name:                   models.MustNewTaskType("expensiveBridge"),
				URL:                    cltest.WebURL(t, "https://localhost:80"),
				Confirmations:          0,
				MinimumContractPayment: test.bridgePayment,
			}
			require.NoError(t, store.CreateBridgeType(bt))

			job := cltest.NewJobWithRunLogInitiator()
			job.MinPayment = test.jobMinimumPayment
			job.Tasks = []models.TaskSpec{
				cltest.NewTask(t, "NoOp"),
				cltest.NewTask(t, bt.Name.String()),
			}
			require.NoError(t, store.CreateJob(&job))
			initiator := job.Initiators[0]

			creationHeight := big.NewInt(1)

			runRequest := models.NewRunRequest(models.JSON{})
			runRequest.Payment = test.inputPayment
			runRequest.RequestParams = cltest.JSONFromString(t, `{"random": "input"}`)

			pusher := new(mocks.StatsPusher)
			pusher.On("PushNow").Return(nil)

			runQueue := new(mocks.RunQueue)
			runQueue.On("Run", mock.Anything).Return(nil)

			runManager := services.NewRunManager(runQueue, store.Config, store.ORM, pusher, store.TxManager, store.Clock)
			run, err := runManager.Create(job.ID, &initiator, creationHeight, runRequest)
			require.NoError(t, err)

			assert.Equal(t, test.jobStatus, run.GetStatus())
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
	app, cleanup := cltest.NewApplicationWithConfig(t, config,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	kst := new(mocks.KeyStoreInterface)
	kst.On("Accounts").Return([]accounts.Account{})
	app.Store.KeyStore = kst
	defer cleanup()

	require.NoError(t, app.StartAndConnect())

	store := app.GetStore()
	job := cltest.NewJobWithRunLogInitiator()
	job.Tasks = []models.TaskSpec{cltest.NewTask(t, "NoOp")}
	require.NoError(t, store.CreateJob(&job))

	requestID := common.HexToHash("0xcafe")
	initiator := job.Initiators[0]
	rr := models.NewRunRequest(models.JSON{})
	rr.RequestID = &requestID
	rr.TxHash = &initiatingTxHash
	rr.BlockHash = &triggeringBlockHash

	futureCreationHeight := big.NewInt(9)
	pastCurrentHeight := big.NewInt(1)

	rr.RequestParams = cltest.JSONFromString(t, `{"random": "input"}`)
	jr, err := app.RunManager.Create(job.ID, &initiator, futureCreationHeight, rr)
	require.NoError(t, err)
	cltest.WaitForJobRunToPendIncomingConfirmations(t, app.Store, *jr)

	err = app.RunManager.ResumeAllPendingNextBlock(pastCurrentHeight)
	require.NoError(t, err)

	updatedJR := cltest.WaitForJobRunToPendIncomingConfirmations(t, app.Store, *jr)
	assert.True(t, updatedJR.TaskRuns[0].ObservedIncomingConfirmations.Valid)
	assert.Equal(t, uint32(0), updatedJR.TaskRuns[0].ObservedIncomingConfirmations.Uint32)

	kst.AssertExpectations(t)
}

func TestRunManager_ResumeConfirmingTasks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status models.RunStatus
	}{
		{models.RunStatusPendingConnection},
		{models.RunStatusPendingIncomingConfirmations},
		{models.RunStatusPendingOutgoingConfirmations},
	}

	for _, test := range tests {
		t.Run(string(test.status), func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			job := cltest.NewJobWithWebInitiator()
			require.NoError(t, store.CreateJob(&job))
			run := cltest.NewJobRun(job)
			run.SetStatus(test.status)
			require.NoError(t, store.CreateJobRun(&run))

			pusher := new(mocks.StatsPusher)
			pusher.On("PushNow").Return(nil)

			runQueue := new(mocks.RunQueue)
			runQueue.On("Run", mock.Anything).Return(nil)

			runManager := services.NewRunManager(runQueue, store.Config, store.ORM, pusher, store.TxManager, store.Clock)
			runManager.ResumeAllPendingNextBlock(big.NewInt(3821))

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
			run := cltest.NewJobRun(job)
			run.SetStatus(test.status)
			require.NoError(t, store.CreateJobRun(&run))

			pusher := new(mocks.StatsPusher)

			runQueue := new(mocks.RunQueue)
			runQueue.On("Run", mock.Anything).Return(nil)

			runManager := services.NewRunManager(runQueue, store.Config, store.ORM, pusher, store.TxManager, store.Clock)
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
			run := cltest.NewJobRun(job)
			run.SetStatus(test.status)
			run.DeletedAt = null.TimeFrom(time.Now())
			require.NoError(t, store.CreateJobRun(&run))

			pusher := new(mocks.StatsPusher)

			runQueue := new(mocks.RunQueue)
			runQueue.On("Run", mock.Anything).Return(nil)

			runManager := services.NewRunManager(runQueue, store.Config, store.ORM, pusher, store.TxManager, store.Clock)
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
		{models.RunStatusPendingIncomingConfirmations},
		{models.RunStatusPendingOutgoingConfirmations},
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
			run := cltest.NewJobRun(job)
			run.SetStatus(test.status)
			require.NoError(t, store.CreateJobRun(&run))

			pusher := new(mocks.StatsPusher)

			runQueue := new(mocks.RunQueue)
			runQueue.On("Run", mock.Anything).Maybe().Return(nil)

			runManager := services.NewRunManager(runQueue, store.Config, store.ORM, pusher, store.TxManager, store.Clock)
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
		{models.RunStatusPendingIncomingConfirmations},
		{models.RunStatusPendingOutgoingConfirmations},
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
			run := cltest.NewJobRun(job)
			run.SetStatus(test.status)
			run.DeletedAt = null.TimeFrom(time.Now())
			require.NoError(t, store.CreateJobRun(&run))

			pusher := new(mocks.StatsPusher)

			runQueue := new(mocks.RunQueue)
			runQueue.On("Run", mock.Anything).Maybe().Return(nil)

			runManager := services.NewRunManager(runQueue, store.Config, store.ORM, pusher, store.TxManager, store.Clock)
			runManager.ResumeAllInProgress()

			runQueue.AssertExpectations(t)
		})
	}
}

func TestRunManager_ValidateRun_PaymentAboveThreshold(t *testing.T) {
	jobSpecID := cltest.NewJob().ID
	run := &models.JobRun{ID: models.NewID(), JobSpecID: jobSpecID, Payment: assets.NewLink(2)}
	contractCost := assets.NewLink(1)

	services.ValidateRun(run, contractCost)

	assert.Equal(t, models.RunStatus(""), run.GetStatus())
}

func TestRunManager_ValidateRun_PaymentBelowThreshold(t *testing.T) {
	jobSpecID := cltest.NewJob().ID
	run := &models.JobRun{ID: models.NewID(), JobSpecID: jobSpecID, Payment: assets.NewLink(1)}
	contractCost := assets.NewLink(2)

	services.ValidateRun(run, contractCost)

	assert.Equal(t, models.RunStatusErrored, run.GetStatus())

	expectedErrorMsg := fmt.Sprintf("rejecting job %s with payment 1 below minimum threshold (2)", jobSpecID)
	assert.Equal(t, expectedErrorMsg, run.Result.ErrorMessage.String)
}

func TestRunManager_NewRun(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	now := time.Now()
	job := cltest.NewJobWithWebInitiator()
	assert.Len(t, job.Tasks, 1)

	t.Run("creates a run with a block height and all adapters", func(t *testing.T) {
		run, adapters := services.NewRun(&job, &job.Initiators[0], big.NewInt(0), &models.RunRequest{}, store.Config, store.ORM, now)
		assert.NotNil(t, run.ID)
		assert.NotNil(t, run.JobSpecID)
		assert.Equal(t, run.GetStatus(), models.RunStatusInProgress)
		assert.Equal(t, utils.NewBig(big.NewInt(0)), run.CreationHeight)
		assert.Equal(t, utils.NewBig(big.NewInt(0)), run.ObservedHeight)
		require.Len(t, run.TaskRuns, 1)
		assert.NotNil(t, run.TaskRuns[0].ID)
		assert.Len(t, adapters, 1)
	})

	t.Run("with no block height creates a run with all adapters", func(t *testing.T) {
		run, adapters := services.NewRun(&job, &job.Initiators[0], nil, &models.RunRequest{}, store.Config, store.ORM, now)
		assert.NotNil(t, run.ID)
		assert.NotNil(t, run.JobSpecID)
		assert.Equal(t, run.GetStatus(), models.RunStatusInProgress)
		assert.Nil(t, run.CreationHeight)
		assert.Nil(t, run.ObservedHeight)
		require.Len(t, run.TaskRuns, 1)
		assert.NotNil(t, run.TaskRuns[0].ID)
		assert.Len(t, adapters, 1)
	})
}
