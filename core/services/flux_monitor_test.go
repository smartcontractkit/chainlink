package services_test

import (
	"chainlink/core/cmd"
	"chainlink/core/internal/cltest"
	"chainlink/core/internal/mocks"
	"chainlink/core/services"
	"chainlink/core/store/models"
	"chainlink/core/utils"
	"context"
	"fmt"
	"math/big"
	"net/url"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type successFetcher decimal.Decimal

func (f *successFetcher) Fetch() (decimal.Decimal, error) {
	return decimal.Decimal(*f), nil
}

func fakeSubscription() *mocks.Subscription {
	sub := new(mocks.Subscription)
	sub.On("Unsubscribe").Return()
	sub.On("Err").Return(nil)
	return sub
}

func TestConcreteFluxMonitor_AddJobRemoveJobHappy(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job := cltest.NewJobWithFluxMonitorInitiator()
	runManager := new(mocks.RunManager)
	started := make(chan struct{}, 1)

	dc := new(mocks.DeviationChecker)
	dc.On("Start", mock.Anything, mock.Anything).Return(nil).Run(func(mock.Arguments) {
		started <- struct{}{}
	})

	checkerFactory := new(mocks.DeviationCheckerFactory)
	checkerFactory.On("New", job.Initiators[0], runManager, store.ORM).Return(dc, nil)
	fm := services.NewFluxMonitor(store, runManager)
	services.ExportedSetCheckerFactory(fm, checkerFactory)
	require.NoError(t, fm.Start())
	defer fm.Stop()
	require.NoError(t, fm.Connect(nil))
	defer fm.Disconnect()

	// Add Job
	require.NoError(t, fm.AddJob(job))

	cltest.CallbackOrTimeout(t, "deviation checker started", func() {
		<-started
	})
	checkerFactory.AssertExpectations(t)
	dc.AssertExpectations(t)

	// Remove Job
	removed := make(chan struct{})
	dc.On("Stop").Return().Run(func(mock.Arguments) {
		removed <- struct{}{}
	})
	fm.RemoveJob(job.ID)
	cltest.CallbackOrTimeout(t, "deviation checker stopped", func() {
		<-removed
	})
	dc.AssertExpectations(t)
}

func TestConcreteFluxMonitor_AddJobError(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job := cltest.NewJobWithFluxMonitorInitiator()
	runManager := new(mocks.RunManager)
	dc := new(mocks.DeviationChecker)
	dc.On("Start", mock.Anything, mock.Anything).Return(errors.New("deliberate test error"))
	checkerFactory := new(mocks.DeviationCheckerFactory)
	checkerFactory.On("New", job.Initiators[0], runManager, store.ORM).Return(dc, nil)
	fm := services.NewFluxMonitor(store, runManager)
	services.ExportedSetCheckerFactory(fm, checkerFactory)
	require.NoError(t, fm.Start())
	defer fm.Stop()
	require.NoError(t, fm.Connect(nil))
	defer fm.Disconnect()

	require.Error(t, fm.AddJob(job))
	checkerFactory.AssertExpectations(t)
	dc.AssertExpectations(t)
}

func TestConcreteFluxMonitor_AddJobDisconnected(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job := cltest.NewJobWithFluxMonitorInitiator()
	runManager := new(mocks.RunManager)
	checkerFactory := new(mocks.DeviationCheckerFactory)
	dc := new(mocks.DeviationChecker)
	checkerFactory.On("New", job.Initiators[0], runManager, store.ORM).Return(dc, nil)
	fm := services.NewFluxMonitor(store, runManager)
	services.ExportedSetCheckerFactory(fm, checkerFactory)
	require.NoError(t, fm.Start())
	defer fm.Stop()

	require.NoError(t, fm.AddJob(job))
}

func TestConcreteFluxMonitor_AddJobNonFluxMonitor(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job := cltest.NewJobWithRunLogInitiator()
	runManager := new(mocks.RunManager)
	checkerFactory := new(mocks.DeviationCheckerFactory)
	fm := services.NewFluxMonitor(store, runManager)
	services.ExportedSetCheckerFactory(fm, checkerFactory)
	require.NoError(t, fm.Start())
	defer fm.Stop()

	require.NoError(t, fm.AddJob(job))
}

func TestConcreteFluxMonitor_ConnectStartsExistingJobs(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	runManager := new(mocks.RunManager)
	started := make(chan struct{})

	dc := new(mocks.DeviationChecker)
	dc.On("Start", mock.Anything, mock.Anything).Return(nil).Run(func(mock.Arguments) {
		started <- struct{}{}
	})

	checkerFactory := new(mocks.DeviationCheckerFactory)

	for i := 0; i < 3; i++ {
		job := cltest.NewJobWithFluxMonitorInitiator()
		require.NoError(t, store.CreateJob(&job))
		job, err := store.FindJob(job.ID)
		require.NoError(t, err)
		checkerFactory.On("New", job.Initiators[0], runManager, store.ORM).Return(dc, nil)
	}

	fm := services.NewFluxMonitor(store, runManager)
	services.ExportedSetCheckerFactory(fm, checkerFactory)
	err := fm.Start()
	require.NoError(t, err)
	defer fm.Stop()

	require.NoError(t, fm.Connect(nil))
	cltest.CallbackOrTimeout(t, "deviation checker started", func() {
		<-started
	})
	checkerFactory.AssertExpectations(t)
	dc.AssertExpectations(t)
}

func TestConcreteFluxMonitor_StopWithoutStart(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	runManager := new(mocks.RunManager)

	fm := services.NewFluxMonitor(store, runManager)
	fm.Stop()
}

func TestPollingDeviationChecker_PollHappy(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	fetcher := new(mocks.Fetcher)
	fetcher.On("Fetch").Return(decimal.NewFromInt(102), nil)

	job := cltest.NewJobWithFluxMonitorInitiator()
	initr := job.Initiators[0]
	initr.ID = 1

	rm := new(mocks.RunManager)
	run := cltest.NewJobRun(job)
	data, err := models.ParseJSON([]byte(fmt.Sprintf(`{
			"result": "102",
			"address": "%s",
			"functionSelector": "0xe6330cf7",
			"dataPrefix": "0x0000000000000000000000000000000000000000000000000000000000000002"
	}`, initr.InitiatorParams.Address.Hex())))
	require.NoError(t, err)
	rm.On("Create", job.ID, &initr, mock.Anything, mock.MatchedBy(func(runRequest *models.RunRequest) bool {
		return runRequest.RequestParams == data
	})).Return(&run, nil)

	checker, err := services.NewPollingDeviationChecker(store, initr, rm, fetcher, time.Second)
	require.NoError(t, err)

	ethClient := new(mocks.Client)
	ethClient.On("GetAggregatorPrice", initr.InitiatorParams.Address, initr.InitiatorParams.Precision).
		Return(decimal.NewFromInt(100), nil)
	ethClient.On("GetAggregatorRound", initr.InitiatorParams.Address).
		Return(big.NewInt(1), nil)

	require.NoError(t, checker.ExportedFetchAggregatorData(ethClient)) // setup
	ethClient.AssertExpectations(t)
	assert.Equal(t, decimal.NewFromInt(100), checker.ExportedCurrentPrice())
	assert.Equal(t, big.NewInt(1), checker.ExportedCurrentRound())

	_, err = checker.ExportedPoll()
	require.NoError(t, err) // main entry point

	fetcher.AssertExpectations(t)
	rm.AssertExpectations(t)
	assert.Equal(t, decimal.NewFromInt(102), checker.ExportedCurrentPrice())
	assert.Equal(t, big.NewInt(2), checker.ExportedCurrentRound())
}

func TestPollingDeviationChecker_TriggerIdleTimeThreshold(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	auth := cmd.TerminalKeyStoreAuthenticator{Prompter: &cltest.MockCountingPrompter{T: t}}
	_, err := auth.Authenticate(store, "somepassword")
	assert.NoError(t, err)
	assert.True(t, store.KeyStore.HasAccounts())

	job := cltest.NewJobWithFluxMonitorInitiator()
	initr := job.Initiators[0]
	initr.ID = 1
	initr.PollingInterval = models.Duration(5 * time.Millisecond)
	initr.IdleThreshold = models.Duration(10 * time.Millisecond)

	jobRun := cltest.NewJobRun(job)

	runManager := new(mocks.RunManager)

	randomLargeNumber := 100
	jobRunCreated := make(chan struct{}, randomLargeNumber)
	runManager.On("Create", job.ID, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&jobRun, nil).
		Run(func(args mock.Arguments) {
			jobRunCreated <- struct{}{}
		})

	fetcher := successFetcher(decimal.NewFromInt(100))
	deviationChecker, err := services.NewPollingDeviationChecker(
		store,
		initr,
		runManager,
		&fetcher,
		time.Second,
	)
	require.NoError(t, err)

	ethClient := new(mocks.Client)
	ethClient.On("GetAggregatorPrice", initr.InitiatorParams.Address, initr.InitiatorParams.Precision).
		Return(decimal.NewFromInt(100), nil)
	ethClient.On("GetAggregatorRound", initr.InitiatorParams.Address).
		Return(big.NewInt(1), nil)
	ethClient.On("SubscribeToLogs", mock.Anything, mock.Anything).
		Return(fakeSubscription(), nil)
	ethClient.On("GetLatestSubmission", mock.Anything, mock.Anything).
		Return(big.NewInt(0), big.NewInt(0), nil)

	err = deviationChecker.Start(context.Background(), ethClient)
	require.NoError(t, err)
	require.Len(t, jobRunCreated, 0, "no Job Runs created")

	require.Eventually(t, func() bool { return len(jobRunCreated) >= 1 }, time.Second, time.Millisecond, "idleThreshold triggers Job Run")
	require.Eventually(t, func() bool { return len(jobRunCreated) >= 5 }, time.Second, time.Millisecond, "idleThreshold triggers succeeding Job Runs")

	deviationChecker.Stop()

	assert.Equal(t, decimal.NewFromInt(100).String(), deviationChecker.ExportedCurrentPrice().String())
}

func TestPollingDeviationChecker_StartError(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	rm := new(mocks.RunManager)
	job := cltest.NewJobWithFluxMonitorInitiator()
	initr := job.Initiators[0]
	initr.ID = 1

	ethClient := new(mocks.Client)
	ethClient.On("GetAggregatorPrice", initr.InitiatorParams.Address, initr.InitiatorParams.Precision).
		Return(decimal.NewFromInt(0), errors.New("deliberate test error"))

	checker, err := services.NewPollingDeviationChecker(store, initr, rm, nil, time.Second)
	require.NoError(t, err)
	require.Error(t, checker.Start(context.Background(), ethClient))
}

func TestPollingDeviationChecker_StartStop(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	// Prepare initialization to 100, which matches external adapter, so no deviation
	job := cltest.NewJobWithFluxMonitorInitiator()
	initr := job.Initiators[0]
	initr.ID = 1

	ethClient := new(mocks.Client)
	ethClient.On("GetAggregatorPrice", initr.InitiatorParams.Address, initr.InitiatorParams.Precision).
		Return(decimal.NewFromInt(100), nil)
	ethClient.On("GetAggregatorRound", initr.InitiatorParams.Address).
		Return(big.NewInt(1), nil)
	ethClient.On("SubscribeToLogs", mock.Anything, mock.Anything).
		Return(fakeSubscription(), nil)

	rm := new(mocks.RunManager)
	fetcher := new(mocks.Fetcher)
	checker, err := services.NewPollingDeviationChecker(store, initr, rm, fetcher, time.Millisecond)
	require.NoError(t, err)

	// Set up fetcher to mark when polled
	started := make(chan struct{})
	fetcher.On("Fetch").Return(decimal.NewFromFloat(100.0), nil).Maybe().Run(func(mock.Arguments) {
		started <- struct{}{}
	})

	// Start() with no delay to speed up test and polling.
	done := make(chan struct{})
	go func() {
		checker.Start(context.Background(), ethClient) // Start() polling
		done <- struct{}{}
	}()

	cltest.CallbackOrTimeout(t, "Start() starts", func() {
		<-started
	})
	fetcher.AssertExpectations(t)

	checker.Stop()
	cltest.CallbackOrTimeout(t, "Stop() unblocks Start()", func() {
		<-done
	})
}

func TestPollingDeviationChecker_NoDeviation_CanBeCanceled(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	auth := cmd.TerminalKeyStoreAuthenticator{Prompter: &cltest.MockCountingPrompter{T: t}}
	_, err := auth.Authenticate(store, "somepassword")
	assert.NoError(t, err)
	assert.True(t, store.KeyStore.HasAccounts())

	// Set up fetcher to mark when polled
	fetcher := new(mocks.Fetcher)
	polled := make(chan struct{})
	fetcher.On("Fetch").Return(decimal.NewFromFloat(100.0), nil).Run(func(mock.Arguments) {
		polled <- struct{}{}
	})

	// Prepare initialization to 100, which matches external adapter, so no deviation
	job := cltest.NewJobWithFluxMonitorInitiator()
	initr := job.Initiators[0]
	initr.ID = 1

	ethClient := new(mocks.Client)
	ethClient.On("GetAggregatorPrice", initr.InitiatorParams.Address, initr.InitiatorParams.Precision).
		Return(decimal.NewFromInt(100), nil)
	ethClient.On("GetAggregatorRound", initr.InitiatorParams.Address).
		Return(big.NewInt(1), nil)
	ethClient.On("SubscribeToLogs", mock.Anything, mock.Anything).
		Return(fakeSubscription(), nil)
	ethClient.On("GetLatestSubmission", mock.Anything, mock.Anything).
		Return(big.NewInt(0), big.NewInt(0), nil)

	// Start() with no delay to speed up test and polling.
	rm := new(mocks.RunManager) // No mocks assert no runs are created
	checker, err := services.NewPollingDeviationChecker(store, initr, rm, fetcher, time.Millisecond)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		checker.Start(ctx, ethClient) // Start() polling until cancel()
		done <- struct{}{}
	}()

	// Check if Polled
	cltest.CallbackOrTimeout(t, "start repeatedly polls external adapter", func() {
		<-polled // launched at the beginning of Start
		<-polled // launched after time.After
	})
	fetcher.AssertExpectations(t)

	// Cancel parent context and ensure Start() stops.
	cancel()
	cltest.CallbackOrTimeout(t, "Start() unblocks and is done", func() {
		<-done
	})
}

func TestPollingDeviationChecker_StopWithoutStart(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	rm := new(mocks.RunManager)
	job := cltest.NewJobWithFluxMonitorInitiator()
	initr := job.Initiators[0]
	initr.ID = 1

	checker, err := services.NewPollingDeviationChecker(store, initr, rm, nil, time.Second)
	require.NoError(t, err)
	checker.Stop()
}

func TestPollingDeviationChecker_RespondToNewRound_Ignore(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	currentRound := int64(5)

	// Prepare on-chain initialization to 100
	job := cltest.NewJobWithFluxMonitorInitiator()
	initr := job.Initiators[0]
	initr.ID = 1

	ethClient := new(mocks.Client)
	ethClient.On("GetAggregatorPrice", initr.InitiatorParams.Address, initr.InitiatorParams.Precision).
		Return(decimal.NewFromInt(100), nil)
	ethClient.On("GetAggregatorRound", initr.InitiatorParams.Address).
		Return(big.NewInt(currentRound), nil)

	// Initialize
	rm := new(mocks.RunManager)
	fetcher := new(mocks.Fetcher)
	checker, err := services.NewPollingDeviationChecker(store, initr, rm, fetcher, time.Minute)
	require.NoError(t, err)
	require.NoError(t, checker.ExportedFetchAggregatorData(ethClient))
	ethClient.AssertExpectations(t)

	// Send rounds less than or equal to current, sequentially
	tests := []struct {
		name  string
		round uint64
	}{
		{"less than", 4},
		{"equal", 5},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			log := cltest.LogFromFixture(t, "testdata/new_round_log.json")
			log.Topics[models.NewRoundTopicRoundID] = common.BytesToHash(utils.EVMWordUint64(test.round))
			require.NoError(t, checker.ExportedRespondToNewRound(log))
			rm.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		})
	}
}

func TestPollingDeviationChecker_RespondToNewRound_Respond(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	currentRound := int64(5)

	// Prepare on-chain initialization to 100, which matches external adapter,
	// so no deviation
	job := cltest.NewJobWithFluxMonitorInitiator()
	initr := job.Initiators[0]
	initr.ID = 1

	ethClient := new(mocks.Client)
	ethClient.On("GetAggregatorPrice", initr.InitiatorParams.Address, initr.InitiatorParams.Precision).
		Return(decimal.NewFromInt(100), nil)
	ethClient.On("GetAggregatorRound", initr.InitiatorParams.Address).
		Return(big.NewInt(currentRound), nil)

	// Initialize
	rm := new(mocks.RunManager)
	fetcher := new(mocks.Fetcher)
	checker, err := services.NewPollingDeviationChecker(store, initr, rm, fetcher, time.Minute)
	require.NoError(t, err)
	require.NoError(t, checker.ExportedFetchAggregatorData(ethClient))
	ethClient.AssertExpectations(t)

	// Send log greater than current
	data, err := models.ParseJSON([]byte(fmt.Sprintf(`{
			"result": "100",
			"address": "%s",
			"functionSelector": "0xe6330cf7",
			"dataPrefix": "0x0000000000000000000000000000000000000000000000000000000000000006"
	}`, initr.InitiatorParams.Address.Hex()))) // dataPrefix has currentRound + 1
	require.NoError(t, err)

	// Set up fetcher for 100; even if within deviation, forces the creation of run.
	fetcher.On("Fetch").Return(decimal.NewFromFloat(100.0), nil).Maybe()

	rm.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.MatchedBy(func(runRequest *models.RunRequest) bool {
		return runRequest.RequestParams == data
	})).Return(nil, nil) // only round 6 triggers run.

	log := cltest.LogFromFixture(t, "testdata/new_round_log.json")
	log.Topics[models.NewRoundTopicRoundID] = common.BytesToHash(utils.EVMWordUint64(6))
	require.NoError(t, checker.ExportedRespondToNewRound(log))
	fetcher.AssertExpectations(t)
	rm.AssertExpectations(t)
}

func TestOutsideDeviation(t *testing.T) {
	tests := []struct {
		name                string
		curPrice, nextPrice decimal.Decimal
		threshold           float64 // in percentage
		expectation         bool
	}{
		{"0 current price", decimal.NewFromInt(0), decimal.NewFromInt(100), 2, true},
		{"inside deviation", decimal.NewFromInt(100), decimal.NewFromInt(101), 2, false},
		{"equal to deviation", decimal.NewFromInt(100), decimal.NewFromInt(102), 2, true},
		{"outside deviation", decimal.NewFromInt(100), decimal.NewFromInt(103), 2, true},
		{"outside deviation zero", decimal.NewFromInt(100), decimal.NewFromInt(0), 2, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := services.OutsideDeviation(test.curPrice, test.nextPrice, test.threshold)
			assert.Equal(t, test.expectation, actual)
		})
	}
}

func TestExtractFeedURLs(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	bridge := &models.BridgeType{
		Name: models.MustNewTaskType("testbridge"),
		URL:  cltest.WebURL(t, "https://testing.com/bridges"),
	}
	require.NoError(t, store.CreateBridgeType(bridge))

	tests := []struct {
		name        string
		in          string
		expectation []string
	}{
		{
			"single",
			`["https://lambda.staging.devnet.tools/bnc/call"]`,
			[]string{"https://lambda.staging.devnet.tools/bnc/call"},
		},
		{
			"double",
			`["https://lambda.staging.devnet.tools/bnc/call", "https://lambda.staging.devnet.tools/cc/call"]`,
			[]string{"https://lambda.staging.devnet.tools/bnc/call", "https://lambda.staging.devnet.tools/cc/call"},
		},
		{
			"bridge",
			`[{"bridge":"testbridge"}]`,
			[]string{"https://testing.com/bridges"},
		},
		{
			"mixed",
			`["https://lambda.staging.devnet.tools/bnc/call", {"bridge": "testbridge"}]`,
			[]string{"https://lambda.staging.devnet.tools/bnc/call", "https://testing.com/bridges"},
		},
		{
			"empty",
			`[]`,
			[]string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			initiatorParams := models.InitiatorParams{
				Feeds: cltest.JSONFromString(t, test.in),
			}
			var expectation []*url.URL
			for _, urlString := range test.expectation {
				expectation = append(expectation, cltest.MustParseURL(urlString))
			}
			val, err := services.ExtractFeedURLs(initiatorParams.Feeds, store.ORM)
			require.NoError(t, err)
			assert.Equal(t, val, expectation)
		})
	}
}
