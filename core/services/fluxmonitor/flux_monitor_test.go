package fluxmonitor_test

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/eth/contracts"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitor"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/onsi/gomega"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const oracleCount uint8 = 17

var (
	submitHash     = utils.MustHash("submit(uint256,int256)")
	submitSelector = submitHash[:4]
)

func ensureAccount(t *testing.T, store *store.Store) common.Address {
	t.Helper()
	auth := cmd.TerminalKeyStoreAuthenticator{Prompter: &cltest.MockCountingPrompter{T: t}}
	_, err := auth.Authenticate(store, cltest.Password)
	assert.NoError(t, err)
	assert.True(t, store.KeyStore.HasAccounts())
	acct, err := store.KeyStore.GetFirstAccount()
	assert.NoError(t, err)
	return acct.Address
}

func TestConcreteFluxMonitor_Start_withEthereumDisabled(t *testing.T) {
	tests := []struct {
		name        string
		enabled     bool
		wantStarted bool
	}{
		{"enabled", true, false},
		{"disabled", false, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, cleanup := cltest.NewConfig(t)
			defer cleanup()
			config.Config.Set("ETH_DISABLED", test.enabled)
			store, cleanup := cltest.NewStoreWithConfig(config)
			defer cleanup()
			runManager := new(mocks.RunManager)

			lb := eth.NewLogBroadcaster(store.TxManager, store.ORM, store.Config.BlockBackfillDepth())
			fm := fluxmonitor.New(store, runManager, lb)

			err := fm.Start()
			require.NoError(t, err)
			defer fm.Stop()
		})
	}
}

func TestConcreteFluxMonitor_AddJobRemoveJob(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	txm := new(mocks.TxManager)
	store.TxManager = txm
	txm.On("GetLatestBlock").Return(models.Block{Number: hexutil.Uint64(123)}, nil)
	txm.On("GetLogs", mock.Anything).Return([]models.Log{}, nil)

	t.Run("starts and stops DeviationCheckers when jobs are added and removed", func(t *testing.T) {
		job := cltest.NewJobWithFluxMonitorInitiator()
		runManager := new(mocks.RunManager)
		started := make(chan struct{}, 1)

		dc := new(mocks.DeviationChecker)
		dc.On("Start", mock.Anything, mock.Anything).Return(nil).Run(func(mock.Arguments) {
			started <- struct{}{}
		})

		checkerFactory := new(mocks.DeviationCheckerFactory)
		checkerFactory.On("New", job.Initiators[0], mock.Anything, runManager, store.ORM, store.Config.DefaultHTTPTimeout()).Return(dc, nil)
		lb := eth.NewLogBroadcaster(store.TxManager, store.ORM, store.Config.BlockBackfillDepth())
		require.NoError(t, lb.Start())
		fm := fluxmonitor.New(store, runManager, lb)
		fluxmonitor.ExportedSetCheckerFactory(fm, checkerFactory)
		require.NoError(t, fm.Start())

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

		fm.Stop()

		dc.AssertExpectations(t)
	})

	t.Run("does not error or attempt to start a DeviationChecker when receiving a non-Flux Monitor job", func(t *testing.T) {
		job := cltest.NewJobWithRunLogInitiator()
		runManager := new(mocks.RunManager)
		checkerFactory := new(mocks.DeviationCheckerFactory)
		lb := eth.NewLogBroadcaster(store.TxManager, store.ORM, store.Config.BlockBackfillDepth())
		require.NoError(t, lb.Start())
		fm := fluxmonitor.New(store, runManager, lb)
		fluxmonitor.ExportedSetCheckerFactory(fm, checkerFactory)

		err := fm.Start()
		require.NoError(t, err)
		defer fm.Stop()

		err = fm.AddJob(job)
		require.NoError(t, err)

		checkerFactory.AssertNotCalled(t, "New", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestPollingDeviationChecker_PollIfEligible(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name              string
		eligible          bool
		connected         bool
		funded            bool
		threshold         float64
		absoluteThreshold float64
		latestAnswer      int64
		polledAnswer      int64
		expectedToPoll    bool
		expectedToSubmit  bool
	}{
		{name: "eligible, connected, funded, threshold > 0, answers deviate",
			eligible: true, connected: true, funded: true, threshold: 0.1,
			absoluteThreshold: 200, latestAnswer: 1, polledAnswer: 100,
			expectedToPoll: true, expectedToSubmit: true},
		{name: "eligible, connected, funded, threshold > 0, answers do not deviate",
			eligible: true, connected: true, funded: true, threshold: 0.1,
			absoluteThreshold: 200, latestAnswer: 100, polledAnswer: 100,
			expectedToPoll: true, expectedToSubmit: false},

		{name: "eligible, disconnected, funded, threshold > 0, answers deviate",
			eligible: true, connected: false, funded: true, threshold: 0.1,
			absoluteThreshold: 200, latestAnswer: 1, polledAnswer: 100,
			expectedToPoll: false, expectedToSubmit: false},
		{name: "eligible, disconnected, funded, threshold > 0, answers do not deviate",
			eligible: true, connected: false, funded: true, threshold: 0.1,
			absoluteThreshold: 200, latestAnswer: 100, polledAnswer: 100,
			expectedToPoll: false, expectedToSubmit: false},

		{name: "ineligible, connected, funded, threshold > 0, answers deviate",
			eligible: false, connected: true, funded: true, threshold: 0.1,
			absoluteThreshold: 200, latestAnswer: 1, polledAnswer: 100,
			expectedToPoll: false, expectedToSubmit: false},
		{name: "ineligible, connected, funded, threshold > 0, answers do not deviate",
			eligible: false, connected: true, funded: true, threshold: 0.1,
			absoluteThreshold: 200, latestAnswer: 100, polledAnswer: 100,
			expectedToPoll: false, expectedToSubmit: false},

		{name: "ineligible, disconnected, funded, threshold > 0, answers deviate",
			eligible: false, connected: false, funded: true, threshold: 0.1,
			absoluteThreshold: 200, latestAnswer: 1, polledAnswer: 100,
			expectedToPoll: false, expectedToSubmit: false},
		{name: "ineligible, disconnected, funded, threshold > 0, answers do not deviate",
			eligible: false, connected: false, funded: true, threshold: 0.1,
			absoluteThreshold: 200, latestAnswer: 100, polledAnswer: 100,
			expectedToPoll: false, expectedToSubmit: false},

		{name: "eligible, connected, underfunded, threshold > 0, answers deviate",
			eligible: true, connected: true, funded: false, threshold: 0.1,
			absoluteThreshold: 200, latestAnswer: 1, polledAnswer: 100,
			expectedToPoll: false, expectedToSubmit: false},
		{name: "eligible, connected, underfunded, threshold > 0, answers do not deviate",
			eligible: true, connected: true, funded: false, threshold: 0.1,
			absoluteThreshold: 200, latestAnswer: 100, polledAnswer: 100,
			expectedToPoll: false, expectedToSubmit: false},

		{name: "eligible, disconnected, underfunded, threshold > 0, answers deviate",
			eligible: true, connected: false, funded: false, threshold: 0.1,
			absoluteThreshold: 1, latestAnswer: 200, polledAnswer: 100,
			expectedToPoll: false, expectedToSubmit: false},
		{name: "eligible, disconnected, underfunded, threshold > 0, answers do not deviate",
			eligible: true, connected: false, funded: false, threshold: 0.1,
			absoluteThreshold: 200, latestAnswer: 100, polledAnswer: 100,
			expectedToPoll: false, expectedToSubmit: false},

		{name: "ineligible, connected, underfunded, threshold > 0, answers deviate",
			eligible: false, connected: true, funded: false, threshold: 0.1,
			absoluteThreshold: 200, latestAnswer: 1, polledAnswer: 100,
			expectedToPoll: false, expectedToSubmit: false},
		{name: "ineligible, connected, underfunded, threshold > 0, answers do not deviate",
			eligible: false, connected: true, funded: false, threshold: 0.1,
			absoluteThreshold: 200, latestAnswer: 100, polledAnswer: 100,
			expectedToPoll: false, expectedToSubmit: false},

		{name: "ineligible, disconnected, underfunded, threshold > 0, answers deviate",
			eligible: false, connected: false, funded: false, threshold: 0.1,
			absoluteThreshold: 200, latestAnswer: 1, polledAnswer: 100,
			expectedToPoll: false, expectedToSubmit: false},
		{name: "ineligible, disconnected, underfunded, threshold > 0, answers do not deviate",
			eligible: false, connected: false, funded: false, threshold: 0.1,
			absoluteThreshold: 200, latestAnswer: 100, polledAnswer: 100,
			expectedToPoll: false, expectedToSubmit: false},
	}

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	nodeAddr := ensureAccount(t, store)

	for _, test := range tests {

		// Run one test for relative thresholds, one for absolute thresholds
		for _, thresholds := range []struct{ abs, rel float64 }{{0.1, 200}, {1, 10}} {
			test := test // Copy test so that for loop can overwrite test during asynchronous operation (t.Parallel())
			test.threshold = thresholds.rel
			test.absoluteThreshold = thresholds.abs
			t.Run(test.name, func(t *testing.T) {
				rm := new(mocks.RunManager)
				fetcher := new(mocks.Fetcher)
				fluxAggregator := new(mocks.FluxAggregator)

				job := cltest.NewJobWithFluxMonitorInitiator()
				initr := job.Initiators[0]
				initr.ID = 1

				const reportableRoundID = 2
				latestAnswerNoPrecision := test.latestAnswer * int64(math.Pow10(int(initr.InitiatorParams.Precision)))

				var availableFunds *big.Int
				var paymentAmount *big.Int
				minPayment := store.Config.MinimumContractPayment().ToInt()
				if test.funded {
					availableFunds = big.NewInt(1).Mul(big.NewInt(10000), minPayment)
					paymentAmount = minPayment
				} else {
					availableFunds = big.NewInt(1)
					paymentAmount = minPayment
				}

				roundState := contracts.FluxAggregatorRoundState{
					ReportableRoundID: reportableRoundID,
					EligibleToSubmit:  test.eligible,
					LatestAnswer:      big.NewInt(latestAnswerNoPrecision),
					AvailableFunds:    availableFunds,
					PaymentAmount:     paymentAmount,
					OracleCount:       oracleCount,
				}
				fluxAggregator.On("RoundState", nodeAddr, uint32(0)).Return(roundState, nil).Maybe()

				if test.expectedToPoll {
					fetcher.On("Fetch").Return(decimal.NewFromInt(test.polledAnswer), nil)
				}

				if test.expectedToSubmit {
					run := cltest.NewJobRun(job)
					data, err := models.ParseJSON([]byte(fmt.Sprintf(`{
					"result": "%d",
					"address": "%s",
					"functionSelector": "0x%x",
					"dataPrefix": "0x000000000000000000000000000000000000000000000000000000000000000%d"
				}`, test.polledAnswer, initr.InitiatorParams.Address.Hex(), submitSelector, reportableRoundID)))
					require.NoError(t, err)

					rm.On("Create", job.ID, &initr, mock.Anything, mock.MatchedBy(func(runRequest *models.RunRequest) bool {
						return reflect.DeepEqual(runRequest.RequestParams.Result.Value(), data.Result.Value())
					})).Return(&run, nil)

					fluxAggregator.On("GetMethodID", "submit").Return(submitSelector, nil)
				}

				checker, err := fluxmonitor.NewPollingDeviationChecker(
					store,
					fluxAggregator,
					initr,
					nil,
					rm,
					fetcher,
					func() {},
				)
				require.NoError(t, err)

				if test.connected {
					checker.OnConnect()
				}

				checker.ExportedPollIfEligible(test.threshold, test.absoluteThreshold)

				fluxAggregator.AssertExpectations(t)
				fetcher.AssertExpectations(t)
				rm.AssertExpectations(t)
			})
		}
	}
}

// If the roundState method is unable to communicate with the contract (possibly due to
// incorrect address) then the pollIfEligible method should create a JobSpecErr record
func TestPollingDeviationChecker_PollIfEligible_Creates_JobSpecErr(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	nodeAddr := ensureAccount(t, store)

	rm := new(mocks.RunManager)
	fetcher := new(mocks.Fetcher)
	fluxAggregator := new(mocks.FluxAggregator)

	job := cltest.NewJobWithFluxMonitorInitiator()
	initr := job.Initiators[0]
	roundState := contracts.FluxAggregatorRoundState{}
	require.Len(t, job.Errors, 0)
	err := store.CreateJob(&job)
	require.NoError(t, err)

	fluxAggregator.On("RoundState", nodeAddr, mock.Anything).Return(roundState, errors.New("err")).Once()
	checker, err := fluxmonitor.NewPollingDeviationChecker(
		store,
		fluxAggregator,
		initr,
		nil,
		rm,
		fetcher,
		func() {},
	)
	require.NoError(t, err)
	checker.OnConnect()

	checker.ExportedPollIfEligible(1, 1)

	job, err = store.FindJobWithErrors(job.ID)
	require.NoError(t, err)
	require.Len(t, job.Errors, 1)

	fluxAggregator.AssertExpectations(t)
	fetcher.AssertExpectations(t)
	rm.AssertExpectations(t)
}

func TestPollingDeviationChecker_BuffersLogs(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	nodeAddr := ensureAccount(t, store)

	const (
		fetchedValue = 100
	)

	job := cltest.NewJobWithFluxMonitorInitiator()
	initr := job.Initiators[0]
	initr.ID = 1
	initr.PollTimer.Disabled = true
	initr.IdleTimer.Disabled = true

	// Test helpers
	var (
		makeRoundStateForRoundID = func(roundID uint32) contracts.FluxAggregatorRoundState {
			return contracts.FluxAggregatorRoundState{
				ReportableRoundID: roundID,
				EligibleToSubmit:  true,
				LatestAnswer:      big.NewInt(100 * int64(math.Pow10(int(initr.InitiatorParams.Precision)))),
				AvailableFunds:    store.Config.MinimumContractPayment().ToInt(),
				PaymentAmount:     store.Config.MinimumContractPayment().ToInt(),
			}
		}

		matchRunRequestForRoundID = func(roundID uint32) interface{} {
			data, err := models.ParseJSON([]byte(fmt.Sprintf(`{
                "result": "%d",
                "address": "%s",
                "functionSelector": "0x%x",
                "dataPrefix": "0x000000000000000000000000000000000000000000000000000000000000000%d"
            }`, fetchedValue, initr.InitiatorParams.Address.Hex(), submitSelector, roundID)))
			require.NoError(t, err)

			return mock.MatchedBy(func(runRequest *models.RunRequest) bool {
				return reflect.DeepEqual(runRequest.RequestParams.Result.Value(), data.Result.Value())
			})
		}
	)

	chBlock := make(chan struct{})
	chSafeToAssert := make(chan struct{})
	chSafeToFillQueue := make(chan struct{})

	fluxAggregator := new(mocks.FluxAggregator)
	fluxAggregator.On("SubscribeToLogs", mock.Anything).Return(true, eth.UnsubscribeFunc(func() {}), nil)
	fluxAggregator.On("GetMethodID", "submit").Return(submitSelector, nil)
	fluxAggregator.On("RoundState", nodeAddr, uint32(1)).
		Return(makeRoundStateForRoundID(1), nil).
		Run(func(mock.Arguments) {
			close(chSafeToFillQueue)
			<-chBlock
		}).
		Once()
	fluxAggregator.On("RoundState", nodeAddr, uint32(3)).Return(makeRoundStateForRoundID(3), nil).Once()
	fluxAggregator.On("RoundState", nodeAddr, uint32(4)).Return(makeRoundStateForRoundID(4), nil).Once()

	fetcher := new(mocks.Fetcher)
	fetcher.On("Fetch").Return(decimal.NewFromInt(fetchedValue), nil)

	rm := new(mocks.RunManager)
	run := cltest.NewJobRun(job)

	rm.On("Create", job.ID, &initr, mock.Anything, matchRunRequestForRoundID(1)).Return(&run, nil).Once()
	rm.On("Create", job.ID, &initr, mock.Anything, matchRunRequestForRoundID(3)).Return(&run, nil).Once()
	rm.On("Create", job.ID, &initr, mock.Anything, matchRunRequestForRoundID(4)).Return(&run, nil).Once().
		Run(func(mock.Arguments) { close(chSafeToAssert) })

	checker, err := fluxmonitor.NewPollingDeviationChecker(
		store,
		fluxAggregator,
		initr,
		nil,
		rm,
		fetcher,
		func() {},
	)
	require.NoError(t, err)

	checker.OnConnect()
	checker.Start()

	var logBroadcasts []*mocks.LogBroadcast

	for i := 1; i <= 4; i++ {
		logBroadcast := new(mocks.LogBroadcast)
		logBroadcast.On("Log").Return(&contracts.LogNewRound{RoundId: big.NewInt(int64(i)), StartedAt: big.NewInt(0)})
		logBroadcast.On("WasAlreadyConsumed").Return(false, nil)
		logBroadcast.On("MarkConsumed").Return(nil)
		logBroadcasts = append(logBroadcasts, logBroadcast)
	}

	checker.HandleLog(logBroadcasts[0], nil) // Get the checker to start processing a log so we can freeze it
	<-chSafeToFillQueue
	checker.HandleLog(logBroadcasts[1], nil) // This log is evicted from the priority queue
	checker.HandleLog(logBroadcasts[2], nil)
	checker.HandleLog(logBroadcasts[3], nil)

	close(chBlock)
	<-chSafeToAssert

	fluxAggregator.AssertExpectations(t)
	fetcher.AssertExpectations(t)
	rm.AssertExpectations(t)
}

func TestPollingDeviationChecker_TriggerIdleTimeThreshold(t *testing.T) {

	tests := []struct {
		name              string
		idleTimerDisabled bool
		idleDuration      models.Duration
		expectedToSubmit  bool
	}{
		{"no idleDuration", true, models.MustMakeDuration(0), false},
		{"idleDuration > 0", false, models.MustMakeDuration(2 * time.Second), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			nodeAddr := ensureAccount(t, store)

			fetcher := new(mocks.Fetcher)
			runManager := new(mocks.RunManager)
			fluxAggregator := new(mocks.FluxAggregator)
			logBroadcast := new(mocks.LogBroadcast)

			job := cltest.NewJobWithFluxMonitorInitiator()
			initr := job.Initiators[0]
			initr.ID = 1
			initr.PollTimer.Disabled = true
			initr.IdleTimer.Disabled = test.idleTimerDisabled
			initr.IdleTimer.Duration = test.idleDuration

			const fetchedAnswer = 100
			answerBigInt := big.NewInt(fetchedAnswer * int64(math.Pow10(int(initr.InitiatorParams.Precision))))

			fluxAggregator.On("SubscribeToLogs", mock.Anything).Return(true, eth.UnsubscribeFunc(func() {}), nil)

			idleDurationOccured := make(chan struct{}, 3)

			now := func() uint64 { return uint64(time.Now().UTC().Unix()) }

			if test.expectedToSubmit {
				// idleDuration 1
				roundState1 := contracts.FluxAggregatorRoundState{ReportableRoundID: 1, EligibleToSubmit: false, LatestAnswer: answerBigInt, StartedAt: now()}
				fluxAggregator.On("RoundState", nodeAddr, uint32(0)).Return(roundState1, nil).Once().Run(func(args mock.Arguments) {
					idleDurationOccured <- struct{}{}
				})
			}

			deviationChecker, err := fluxmonitor.NewPollingDeviationChecker(
				store,
				fluxAggregator,
				initr,
				nil,
				runManager,
				fetcher,
				func() {},
			)
			require.NoError(t, err)

			deviationChecker.OnConnect()
			deviationChecker.Start()
			require.Len(t, idleDurationOccured, 0, "no Job Runs created")

			if test.expectedToSubmit {
				require.Eventually(t, func() bool { return len(idleDurationOccured) == 1 }, 3*time.Second, 10*time.Millisecond)

				chBlock := make(chan struct{})
				// NewRound resets the idle timer
				roundState2 := contracts.FluxAggregatorRoundState{ReportableRoundID: 2, EligibleToSubmit: false, LatestAnswer: answerBigInt, StartedAt: now()}
				fluxAggregator.On("RoundState", nodeAddr, uint32(2)).Return(roundState2, nil).Once().Run(func(args mock.Arguments) {
					close(chBlock)
				})

				decodedLog := contracts.LogNewRound{RoundId: big.NewInt(2), StartedAt: big.NewInt(0)}
				logBroadcast.On("Log").Return(&decodedLog)
				logBroadcast.On("WasAlreadyConsumed").Return(false, nil).Once()
				logBroadcast.On("MarkConsumed").Return(nil).Once()
				deviationChecker.HandleLog(logBroadcast, nil)

				gomega.NewGomegaWithT(t).Eventually(chBlock).Should(gomega.BeClosed())

				// idleDuration 2
				roundState3 := contracts.FluxAggregatorRoundState{ReportableRoundID: 3, EligibleToSubmit: false, LatestAnswer: answerBigInt, StartedAt: now()}
				fluxAggregator.On("RoundState", nodeAddr, uint32(0)).Return(roundState3, nil).Once().Run(func(args mock.Arguments) {
					idleDurationOccured <- struct{}{}
				})
				require.Eventually(t, func() bool { return len(idleDurationOccured) == 2 }, 3*time.Second, 10*time.Millisecond)
			}

			deviationChecker.Stop()

			if !test.expectedToSubmit {
				require.Len(t, idleDurationOccured, 0)
			}

			fetcher.AssertExpectations(t)
			runManager.AssertExpectations(t)
			fluxAggregator.AssertExpectations(t)
		})
	}
}

func TestPollingDeviationChecker_RoundTimeoutCausesPoll_timesOutAtZero(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	nodeAddr := ensureAccount(t, store)
	fetcher := new(mocks.Fetcher)
	runManager := new(mocks.RunManager)
	fluxAggregator := new(mocks.FluxAggregator)

	job := cltest.NewJobWithFluxMonitorInitiator()
	initr := job.Initiators[0]
	initr.ID = 1
	initr.PollTimer.Disabled = true
	initr.IdleTimer.Disabled = true

	ch := make(chan struct{})

	const fetchedAnswer = 100
	answerBigInt := big.NewInt(fetchedAnswer * int64(math.Pow10(int(initr.InitiatorParams.Precision))))
	fluxAggregator.On("SubscribeToLogs", mock.Anything).Return(true, eth.UnsubscribeFunc(func() {}), nil)
	fluxAggregator.On("RoundState", nodeAddr, uint32(0)).Return(contracts.FluxAggregatorRoundState{
		ReportableRoundID: 1,
		EligibleToSubmit:  false,
		LatestAnswer:      answerBigInt,
		StartedAt:         0,
		Timeout:           0,
	}, nil).
		Run(func(mock.Arguments) { close(ch) }).
		Once()

	deviationChecker, err := fluxmonitor.NewPollingDeviationChecker(
		store,
		fluxAggregator,
		initr,
		nil,
		runManager,
		fetcher,
		func() {},
	)
	require.NoError(t, err)

	deviationChecker.ExportedRoundState()
	deviationChecker.Start()
	deviationChecker.OnConnect()

	gomega.NewGomegaWithT(t).Eventually(ch).Should(gomega.BeClosed())

	deviationChecker.Stop()

	fetcher.AssertExpectations(t)
	runManager.AssertExpectations(t)
	fluxAggregator.AssertExpectations(t)
}

func TestPollingDeviationChecker_RoundTimeoutCausesPoll_timesOutNotZero(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	nodeAddr := ensureAccount(t, store)

	fetcher := new(mocks.Fetcher)
	runManager := new(mocks.RunManager)
	fluxAggregator := new(mocks.FluxAggregator)

	job := cltest.NewJobWithFluxMonitorInitiator()
	initr := job.Initiators[0]
	initr.ID = 1
	initr.PollTimer.Disabled = true
	initr.IdleTimer.Disabled = true

	const fetchedAnswer = 100
	answerBigInt := big.NewInt(fetchedAnswer * int64(math.Pow10(int(initr.InitiatorParams.Precision))))

	chRoundState1 := make(chan struct{})
	chRoundState2 := make(chan struct{})

	fluxAggregator.On("SubscribeToLogs", mock.Anything).Return(true, eth.UnsubscribeFunc(func() {}), nil)

	startedAt := uint64(time.Now().Unix())
	timeout := uint64(3)
	fluxAggregator.On("RoundState", nodeAddr, uint32(0)).Return(contracts.FluxAggregatorRoundState{
		ReportableRoundID: 1,
		EligibleToSubmit:  false,
		LatestAnswer:      answerBigInt,
		StartedAt:         startedAt,
		Timeout:           timeout,
	}, nil).Once().
		Run(func(mock.Arguments) { close(chRoundState1) }).
		Once()
	fluxAggregator.On("RoundState", nodeAddr, uint32(0)).Return(contracts.FluxAggregatorRoundState{
		ReportableRoundID: 1,
		EligibleToSubmit:  false,
		LatestAnswer:      answerBigInt,
		StartedAt:         startedAt,
		Timeout:           timeout,
	}, nil).Once().
		Run(func(mock.Arguments) { close(chRoundState2) }).
		Once()

	deviationChecker, err := fluxmonitor.NewPollingDeviationChecker(
		store,
		fluxAggregator,
		initr,
		nil,
		runManager,
		fetcher,
		func() {},
	)
	require.NoError(t, err)

	deviationChecker.ExportedRoundState()
	deviationChecker.Start()
	deviationChecker.OnConnect()

	gomega.NewGomegaWithT(t).Eventually(chRoundState1).Should(gomega.BeClosed())
	gomega.NewGomegaWithT(t).Eventually(chRoundState2).Should(gomega.BeClosed())

	time.Sleep(time.Duration(2*timeout) * time.Second)
	deviationChecker.Stop()

	fetcher.AssertExpectations(t)
	runManager.AssertExpectations(t)
	fluxAggregator.AssertExpectations(t)
}

func TestPollingDeviationChecker_RespondToNewRound(t *testing.T) {

	type roundIDCase struct {
		name                     string
		fetchedReportableRoundID uint32
		logRoundID               int64
	}
	var (
		fetched_lt_log = roundIDCase{"fetched < log", 10, 15}
		fetched_gt_log = roundIDCase{"fetched > log", 15, 10}
		fetched_eq_log = roundIDCase{"fetched = log", 10, 10}
	)

	type answerCase struct {
		name         string
		latestAnswer int64
		polledAnswer int64
	}
	var (
		deviationThresholdExceeded    = answerCase{"deviation", 10, 100}
		deviationThresholdNotExceeded = answerCase{"no deviation", 10, 10}
	)

	tests := []struct {
		funded        bool
		eligible      bool
		startedBySelf bool
		roundIDCase
		answerCase
	}{
		{true, true, true, fetched_lt_log, deviationThresholdExceeded},
		{true, true, true, fetched_gt_log, deviationThresholdExceeded},
		{true, true, true, fetched_eq_log, deviationThresholdExceeded},
		{true, true, true, fetched_lt_log, deviationThresholdNotExceeded},
		{true, true, true, fetched_gt_log, deviationThresholdNotExceeded},
		{true, true, true, fetched_eq_log, deviationThresholdNotExceeded},
		{true, true, false, fetched_lt_log, deviationThresholdExceeded},
		{true, true, false, fetched_gt_log, deviationThresholdExceeded},
		{true, true, false, fetched_eq_log, deviationThresholdExceeded},
		{true, true, false, fetched_lt_log, deviationThresholdNotExceeded},
		{true, true, false, fetched_gt_log, deviationThresholdNotExceeded},
		{true, true, false, fetched_eq_log, deviationThresholdNotExceeded},
		{true, false, true, fetched_lt_log, deviationThresholdExceeded},
		{true, false, true, fetched_gt_log, deviationThresholdExceeded},
		{true, false, true, fetched_eq_log, deviationThresholdExceeded},
		{true, false, true, fetched_lt_log, deviationThresholdNotExceeded},
		{true, false, true, fetched_gt_log, deviationThresholdNotExceeded},
		{true, false, true, fetched_eq_log, deviationThresholdNotExceeded},
		{true, false, false, fetched_lt_log, deviationThresholdExceeded},
		{true, false, false, fetched_gt_log, deviationThresholdExceeded},
		{true, false, false, fetched_eq_log, deviationThresholdExceeded},
		{true, false, false, fetched_lt_log, deviationThresholdNotExceeded},
		{true, false, false, fetched_gt_log, deviationThresholdNotExceeded},
		{true, false, false, fetched_eq_log, deviationThresholdNotExceeded},
		{false, true, true, fetched_lt_log, deviationThresholdExceeded},
		{false, true, true, fetched_gt_log, deviationThresholdExceeded},
		{false, true, true, fetched_eq_log, deviationThresholdExceeded},
		{false, true, true, fetched_lt_log, deviationThresholdNotExceeded},
		{false, true, true, fetched_gt_log, deviationThresholdNotExceeded},
		{false, true, true, fetched_eq_log, deviationThresholdNotExceeded},
		{false, true, false, fetched_lt_log, deviationThresholdExceeded},
		{false, true, false, fetched_gt_log, deviationThresholdExceeded},
		{false, true, false, fetched_eq_log, deviationThresholdExceeded},
		{false, true, false, fetched_lt_log, deviationThresholdNotExceeded},
		{false, true, false, fetched_gt_log, deviationThresholdNotExceeded},
		{false, true, false, fetched_eq_log, deviationThresholdNotExceeded},
		{false, false, true, fetched_lt_log, deviationThresholdExceeded},
		{false, false, true, fetched_gt_log, deviationThresholdExceeded},
		{false, false, true, fetched_eq_log, deviationThresholdExceeded},
		{false, false, true, fetched_lt_log, deviationThresholdNotExceeded},
		{false, false, true, fetched_gt_log, deviationThresholdNotExceeded},
		{false, false, true, fetched_eq_log, deviationThresholdNotExceeded},
		{false, false, false, fetched_lt_log, deviationThresholdExceeded},
		{false, false, false, fetched_gt_log, deviationThresholdExceeded},
		{false, false, false, fetched_eq_log, deviationThresholdExceeded},
		{false, false, false, fetched_lt_log, deviationThresholdNotExceeded},
		{false, false, false, fetched_gt_log, deviationThresholdNotExceeded},
		{false, false, false, fetched_eq_log, deviationThresholdNotExceeded},
	}

	for _, test := range tests {
		name := test.answerCase.name + ", " + test.roundIDCase.name
		if test.eligible {
			name += ", eligible"
		} else {
			name += ", ineligible"
		}
		if test.startedBySelf {
			name += ", started by self"
		} else {
			name += ", started by other"
		}
		if test.funded {
			name += ", funded"
		} else {
			name += ", underfunded"
		}

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			nodeAddr := ensureAccount(t, store)

			expectedToFetchRoundState := !test.startedBySelf
			expectedToPoll := expectedToFetchRoundState && test.eligible && test.funded && test.logRoundID >= int64(test.fetchedReportableRoundID)
			expectedToSubmit := expectedToPoll

			job := cltest.NewJobWithFluxMonitorInitiator()
			initr := job.Initiators[0]
			initr.ID = 1
			initr.PollTimer.Disabled = true
			initr.IdleTimer.Disabled = true

			rm := new(mocks.RunManager)
			fetcher := new(mocks.Fetcher)
			fluxAggregator := new(mocks.FluxAggregator)

			paymentAmount := store.Config.MinimumContractPayment().ToInt()
			var availableFunds *big.Int
			if test.funded {
				availableFunds = big.NewInt(1).Mul(paymentAmount, big.NewInt(1000))
			} else {
				availableFunds = big.NewInt(1)
			}

			if expectedToFetchRoundState {
				fluxAggregator.On("RoundState", nodeAddr, uint32(test.logRoundID)).Return(contracts.FluxAggregatorRoundState{
					ReportableRoundID: test.fetchedReportableRoundID,
					LatestAnswer:      big.NewInt(test.latestAnswer * int64(math.Pow10(int(initr.InitiatorParams.Precision)))),
					EligibleToSubmit:  test.eligible,
					AvailableFunds:    availableFunds,
					PaymentAmount:     paymentAmount,
					OracleCount:       oracleCount,
				}, nil).Once()
			}

			if expectedToPoll {
				fetcher.On("Fetch").Return(decimal.NewFromInt(test.polledAnswer), nil).Once()
			}

			if expectedToSubmit {
				fluxAggregator.On("GetMethodID", "submit").Return(submitSelector, nil)

				data, err := models.ParseJSON([]byte(fmt.Sprintf(`{
					"result": "%d",
					"address": "%s",
					"functionSelector": "0x202ee0ed",
					"dataPrefix": "0x%0x"
				}`, test.polledAnswer, initr.InitiatorParams.Address.Hex(), utils.EVMWordUint64(uint64(test.fetchedReportableRoundID)))))
				require.NoError(t, err)

				rm.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.MatchedBy(func(runRequest *models.RunRequest) bool {
					return reflect.DeepEqual(runRequest.RequestParams.Result.Value(), data.Result.Value())
				})).Return(nil, nil)
			}

			checker, err := fluxmonitor.NewPollingDeviationChecker(
				store,
				fluxAggregator,
				initr,
				nil,
				rm,
				fetcher,
				func() {},
			)
			require.NoError(t, err)

			checker.OnConnect()

			var startedBy common.Address
			if test.startedBySelf {
				startedBy = nodeAddr
			}
			checker.ExportedRespondToNewRoundLog(&contracts.LogNewRound{RoundId: big.NewInt(test.logRoundID), StartedBy: startedBy, StartedAt: big.NewInt(0)})

			fluxAggregator.AssertExpectations(t)
			fetcher.AssertExpectations(t)
			rm.AssertExpectations(t)
		})
	}
}

type outsideDeviationRow struct {
	name                string
	curPrice, nextPrice decimal.Decimal
	threshold           float64 // in percentage
	absoluteThreshold   float64
	expectation         bool
}

func (o outsideDeviationRow) String() string {
	return fmt.Sprintf(
		`{name: "%s", curPrice: %s, nextPrice: %s, threshold: %.2f, `+
			"absoluteThreshold: %f, expectation: %v}", o.name, o.curPrice, o.nextPrice,
		o.threshold, o.absoluteThreshold, o.expectation)
}

func TestOutsideDeviation(t *testing.T) {
	t.Parallel()
	f, i := decimal.NewFromFloat, decimal.NewFromInt
	tests := []outsideDeviationRow{
		// Start with a huge absoluteThreshold, to test relative threshold behavior
		{"0 current price, outside deviation", i(0), i(100), 2, 0, true},
		{"0 current and next price", i(0), i(0), 2, 0, false},

		{"inside deviation", i(100), i(101), 2, 0, false},
		{"equal to deviation", i(100), i(102), 2, 0, true},
		{"outside deviation", i(100), i(103), 2, 0, true},
		{"outside deviation zero", i(100), i(0), 2, 0, true},

		{"inside deviation, crosses 0 backwards", f(0.1), f(-0.1), 201, 0, false},
		{"equal to deviation, crosses 0 backwards", f(0.1), f(-0.1), 200, 0, true},
		{"outside deviation, crosses 0 backwards", f(0.1), f(-0.1), 199, 0, true},

		{"inside deviation, crosses 0 forwards", f(-0.1), f(0.1), 201, 0, false},
		{"equal to deviation, crosses 0 forwards", f(-0.1), f(0.1), 200, 0, true},
		{"outside deviation, crosses 0 forwards", f(-0.1), f(0.1), 199, 0, true},

		{"thresholds=0, deviation", i(0), i(100), 0, 0, true},
		{"thresholds=0, no deviation", i(100), i(100), 0, 0, true},
		{"thresholds=0, all zeros", i(0), i(0), 0, 0, true},
	}

	c := func(test outsideDeviationRow) {
		actual := fluxmonitor.OutsideDeviation(test.curPrice, test.nextPrice,
			fluxmonitor.DeviationThresholds{Rel: test.threshold,
				Abs: test.absoluteThreshold})
		assert.Equal(t, test.expectation, actual,
			"check on OutsideDeviation failed for %s", test)
	}

	for _, test := range tests {
		test := test
		// Checks on relative threshold
		t.Run(test.name, func(t *testing.T) { c(test) })
		// Check corresponding absolute threshold tests; make relative threshold
		// always pass (as long as curPrice and nextPrice aren't both 0.)
		test2 := test
		test2.threshold = 0
		// absoluteThreshold is initially zero, so any change will trigger
		test2.expectation = test2.curPrice.Sub(test.nextPrice).Abs().GreaterThan(i(0)) ||
			test2.absoluteThreshold == 0
		t.Run(test.name+" threshold zeroed", func(t *testing.T) { c(test2) })
		// Huge absoluteThreshold means trigger always fails
		test3 := test
		test3.absoluteThreshold = 1e307
		test3.expectation = false
		t.Run(test.name+" max absolute threshold", func(t *testing.T) { c(test3) })
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
			val, err := fluxmonitor.ExtractFeedURLs(initiatorParams.Feeds, store.ORM)
			require.NoError(t, err)
			assert.Equal(t, val, expectation)
		})
	}
}

func TestPollingDeviationChecker_SufficientPayment(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job := cltest.NewJobWithFluxMonitorInitiator()
	initr := job.Initiators[0]
	rm := new(mocks.RunManager)
	fetcher := new(mocks.Fetcher)
	fluxAggregator := new(mocks.FluxAggregator)

	var payment int64 = 10
	var eq = payment
	var gt int64 = payment + 1
	var lt int64 = payment - 1

	tests := []struct {
		name               string
		minContractPayment int64
		minJobPayment      interface{} // nil or int64
		want               bool
	}{
		{"payment above min contract payment, no min job payment", lt, nil, true},
		{"payment equal to min contract payment, no min job payment", eq, nil, true},
		{"payment below min contract payment, no min job payment", gt, nil, false},

		{"payment above min contract payment, above min job payment", lt, lt, true},
		{"payment equal to min contract payment, above min job payment", eq, lt, true},
		{"payment below min contract payment, above min job payment", gt, lt, false},

		{"payment above min contract payment, equal to min job payment", lt, eq, true},
		{"payment equal to min contract payment, equal to min job payment", eq, eq, true},
		{"payment below min contract payment, equal to min job payment", gt, eq, false},

		{"payment above minimum contract payment, below min job payment", lt, gt, false},
		{"payment equal to minimum contract payment, below min job payment", eq, gt, false},
		{"payment below minimum contract payment, below min job payment", gt, gt, false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			store.Config.Set(orm.EnvVarName("MinimumContractPayment"), test.minContractPayment)
			var minJobPayment *assets.Link

			if test.minJobPayment != nil {
				mjb := assets.Link(*big.NewInt(test.minJobPayment.(int64)))
				minJobPayment = &mjb
			}

			checker, err := fluxmonitor.NewPollingDeviationChecker(
				store,
				fluxAggregator,
				initr,
				minJobPayment,
				rm,
				fetcher,
				func() {},
			)
			require.NoError(t, err)

			assert.Equal(t, test.want, checker.ExportedSufficientPayment(big.NewInt(payment)))
		})
	}
}

func TestPollingDeviationChecker_SufficientFunds(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	checker := cltest.NewPollingDeviationChecker(t, store)

	payment := 100
	rounds := 3
	oracleCount := 21
	min := payment * rounds * oracleCount

	tests := []struct {
		name  string
		funds int
		want  bool
	}{
		{"above minimum", min + 1, true},
		{"equal to minimum", min, true},
		{"below minimum", min - 1, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			state := contracts.FluxAggregatorRoundState{
				AvailableFunds: big.NewInt(int64(test.funds)),
				PaymentAmount:  big.NewInt(int64(payment)),
				OracleCount:    uint8(oracleCount),
			}
			assert.Equal(t, test.want, checker.ExportedSufficientFunds(state))
		})
	}
}

func TestFluxMonitor_MakeIdleTimer_RoundStartedAtIsNil(t *testing.T) {
	t.Parallel()

	log := contracts.LogNewRound{}
	idleThreshold, err := models.MakeDuration(5 * time.Second)
	require.NoError(t, err)
	clock := new(mocks.AfterNower)

	clock.On("Now").Return(time.Unix(11, 0))

	timerChannel := make(<-chan time.Time)
	clock.On("After", idleThreshold.Duration()).Return(timerChannel)

	idleTimer := fluxmonitor.MakeIdleTimer(log, idleThreshold, clock)

	assert.Equal(t, timerChannel, idleTimer)

	clock.AssertExpectations(t)
}

func TestFluxMonitor_MakeIdleTimer_RoundStartedAtIsInPast(t *testing.T) {
	// We want to err on the side of the shorter idle timeout, so if round started at is in the past
	// we trust the local clock and adjust the idle timeout down to assume it started counting from
	// round startedAt in terms of our local clock
	t.Parallel()

	log := contracts.LogNewRound{StartedAt: big.NewInt(10)}
	idleThreshold, err := models.MakeDuration(5 * time.Second)
	require.NoError(t, err)
	clock := new(mocks.AfterNower)

	clock.On("Now").Return(time.Unix(11, 0))

	timerChannel := make(<-chan time.Time)
	clock.On("After", 4*time.Second).Return(timerChannel)

	idleTimer := fluxmonitor.MakeIdleTimer(log, idleThreshold, clock)

	assert.Equal(t, timerChannel, idleTimer)

	clock.AssertExpectations(t)
}

func TestFluxMonitor_MakeIdleTimer_IdleThresholdAlreadyPassed(t *testing.T) {
	// If idle threshold is already passed, node should trigger a new round immediately
	t.Parallel()

	log := contracts.LogNewRound{StartedAt: big.NewInt(10)}
	idleThreshold, err := models.MakeDuration(5 * time.Second)
	require.NoError(t, err)
	clock := new(mocks.AfterNower)

	clock.On("Now").Return(time.Unix(42, 0))
	timerChannel := make(<-chan time.Time)
	clock.On("After", mock.MatchedBy(func(d time.Duration) bool {
		// Anything 0 or less is fine since this will expire immediately
		return d <= 0
	})).Return(timerChannel)

	idleTimer := fluxmonitor.MakeIdleTimer(log, idleThreshold, clock)

	assert.Equal(t, timerChannel, idleTimer)

	clock.AssertExpectations(t)
}

func TestFluxMonitor_MakeIdleTimer_OutOfBoundsStartedAt(t *testing.T) {
	// If idle threshold is out of bounds (should never happen!) simply ignore
	// it and wait exactly the idle threshold from now
	t.Parallel()

	var startedAt big.Int
	startedAt.SetUint64(math.MaxUint64)
	log := contracts.LogNewRound{StartedAt: &startedAt}
	idleThreshold, err := models.MakeDuration(5 * time.Second)
	require.NoError(t, err)
	clock := new(mocks.AfterNower)

	clock.On("Now").Return(time.Unix(11, 0))
	timerChannel := make(<-chan time.Time)
	clock.On("After", idleThreshold.Duration()).Return(timerChannel)

	idleTimer := fluxmonitor.MakeIdleTimer(log, idleThreshold, clock)

	assert.Equal(t, timerChannel, idleTimer)

	clock.AssertExpectations(t)
}

func TestFluxMonitor_MakeIdleTimer_RoundStartedAtIsInFuture(t *testing.T) {
	// If the round started at is somehow in the future, this machine probably has a slow clock.
	// Since local time is skewed backwards, we should not attempt to use it for
	// calculating expiry time and instead start counting down the idle timer from now.
	t.Parallel()

	log := contracts.LogNewRound{StartedAt: big.NewInt(40)}
	idleThreshold, err := models.MakeDuration(42 * time.Second)
	require.NoError(t, err)
	clock := new(mocks.AfterNower)

	clock.On("Now").Return(time.Unix(9, 0))
	timerChannel := make(<-chan time.Time)
	clock.On("After", idleThreshold.Duration()).Return(timerChannel)

	idleTimer := fluxmonitor.MakeIdleTimer(log, idleThreshold, clock)

	assert.Equal(t, timerChannel, idleTimer)

	clock.AssertExpectations(t)
}

func TestFluxMonitor_PollingDeviationChecker_HandlesNilLogs(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	p := cltest.NewPollingDeviationChecker(t, store)

	logBroadcast := new(mocks.LogBroadcast)
	var logNewRound *contracts.LogNewRound
	var logAnswerUpdated *contracts.LogAnswerUpdated
	var randomType interface{}

	logBroadcast.On("Log").Return(logNewRound).Once()
	assert.NotPanics(t, func() {
		p.HandleLog(logBroadcast, nil)
	})

	logBroadcast.On("Log").Return(logAnswerUpdated).Once()
	assert.NotPanics(t, func() {
		p.HandleLog(logBroadcast, nil)
	})

	logBroadcast.On("Log").Return(randomType).Once()
	assert.NotPanics(t, func() {
		p.HandleLog(logBroadcast, nil)
	})
}

func TestFluxMonitor_ConsumeLogBroadcast_Happy(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	p := cltest.NewPollingDeviationChecker(t, store)
	p.ExportedFluxAggregator().(*mocks.FluxAggregator).
		On("RoundState", mock.Anything, mock.Anything).
		Return(contracts.FluxAggregatorRoundState{ReportableRoundID: 123}, nil)

	logBroadcast := new(mocks.LogBroadcast)
	logBroadcast.On("WasAlreadyConsumed").Return(false, nil).Once()
	logBroadcast.On("Log").Return(&contracts.LogAnswerUpdated{})
	logBroadcast.On("MarkConsumed").Return(nil).Once()

	p.ExportedBacklog().Add(fluxmonitor.PriorityNewRoundLog, logBroadcast)
	p.ExportedProcessLogs()

	logBroadcast.AssertExpectations(t)
}

func TestFluxMonitor_ConsumeLogBroadcast_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		consumed bool
		err      error
	}{
		{"already consumed", true, nil},
		{"error determining already consumed", false, errors.New("err")},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			p := cltest.NewPollingDeviationChecker(t, store)

			logBroadcast := new(mocks.LogBroadcast)
			logBroadcast.On("WasAlreadyConsumed").Return(test.consumed, test.err).Once()

			p.ExportedBacklog().Add(fluxmonitor.PriorityNewRoundLog, logBroadcast)
			p.ExportedProcessLogs()

			logBroadcast.AssertExpectations(t)
		})
	}
}

func TestPollingDeviationChecker_DoesNotDoubleSubmit(t *testing.T) {
	t.Run("when NewRound log arrives, then poll ticker fires", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()

		nodeAddr := ensureAccount(t, store)

		job := cltest.NewJobWithFluxMonitorInitiator()
		initr := job.Initiators[0]
		initr.ID = 1
		initr.PollTimer.Disabled = true
		initr.IdleTimer.Disabled = true
		run := cltest.NewJobRun(job)

		rm := new(mocks.RunManager)
		fetcher := new(mocks.Fetcher)
		fluxAggregator := new(mocks.FluxAggregator)

		paymentAmount := store.Config.MinimumContractPayment().ToInt()
		availableFunds := big.NewInt(1).Mul(paymentAmount, big.NewInt(1000))

		const (
			roundID = 3
			answer  = 100
		)

		checker, err := fluxmonitor.NewPollingDeviationChecker(
			store,
			fluxAggregator,
			initr,
			nil,
			rm,
			fetcher,
			func() {},
		)
		require.NoError(t, err)

		checker.OnConnect()

		// Fire off the NewRound log, which the node should respond to
		fluxAggregator.On("RoundState", nodeAddr, uint32(roundID)).
			Return(contracts.FluxAggregatorRoundState{
				ReportableRoundID: roundID,
				LatestAnswer:      big.NewInt(answer),
				EligibleToSubmit:  true,
				AvailableFunds:    availableFunds,
				PaymentAmount:     paymentAmount,
				OracleCount:       1,
			}, nil).
			Once()
		fetcher.On("Fetch").
			Return(decimal.NewFromInt(answer), nil).
			Once()
		fluxAggregator.On("GetMethodID", "submit").
			Return(submitSelector, nil).
			Once()
		rm.On("Create", job.ID, &initr, mock.Anything, mock.Anything).
			Return(&run, nil).
			Once()
		checker.ExportedRespondToNewRoundLog(&contracts.LogNewRound{
			RoundId:   big.NewInt(roundID),
			StartedAt: big.NewInt(0),
		})

		// Now force the node to try to poll and ensure it does not respond this time
		fluxAggregator.On("RoundState", nodeAddr, uint32(0)).
			Return(contracts.FluxAggregatorRoundState{
				ReportableRoundID: roundID,
				LatestAnswer:      big.NewInt(answer),
				EligibleToSubmit:  true,
				AvailableFunds:    availableFunds,
				PaymentAmount:     paymentAmount,
				OracleCount:       1,
			}, nil).
			Once()
		checker.ExportedPollIfEligible(0, 0)

		rm.AssertExpectations(t)
		fetcher.AssertExpectations(t)
		fluxAggregator.AssertExpectations(t)
	})

	t.Run("when poll ticker fires, then NewRound log arrives", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()

		nodeAddr := ensureAccount(t, store)

		job := cltest.NewJobWithFluxMonitorInitiator()
		initr := job.Initiators[0]
		initr.ID = 1
		initr.PollTimer.Disabled = true
		initr.IdleTimer.Disabled = true
		run := cltest.NewJobRun(job)

		rm := new(mocks.RunManager)
		fetcher := new(mocks.Fetcher)
		fluxAggregator := new(mocks.FluxAggregator)

		paymentAmount := store.Config.MinimumContractPayment().ToInt()
		availableFunds := big.NewInt(1).Mul(paymentAmount, big.NewInt(1000))

		const (
			roundID = 3
			answer  = 100
		)

		checker, err := fluxmonitor.NewPollingDeviationChecker(
			store,
			fluxAggregator,
			initr,
			nil,
			rm,
			fetcher,
			func() {},
		)
		require.NoError(t, err)

		checker.OnConnect()

		// First, force the node to try to poll, which should result in a submission
		fluxAggregator.On("RoundState", nodeAddr, uint32(0)).
			Return(contracts.FluxAggregatorRoundState{
				ReportableRoundID: roundID,
				LatestAnswer:      big.NewInt(answer),
				EligibleToSubmit:  true,
				AvailableFunds:    availableFunds,
				PaymentAmount:     paymentAmount,
				OracleCount:       1,
			}, nil).
			Once()
		fetcher.On("Fetch").
			Return(decimal.NewFromInt(answer), nil).
			Once()
		fluxAggregator.On("GetMethodID", "submit").
			Return(submitSelector, nil).
			Once()
		rm.On("Create", job.ID, &initr, mock.Anything, mock.Anything).
			Return(&run, nil).
			Once()
		checker.ExportedPollIfEligible(0, 0)

		// Now fire off the NewRound log and ensure it does not respond this time
		checker.ExportedRespondToNewRoundLog(&contracts.LogNewRound{
			RoundId:   big.NewInt(roundID),
			StartedAt: big.NewInt(0),
		})

		rm.AssertExpectations(t)
		fetcher.AssertExpectations(t)
		fluxAggregator.AssertExpectations(t)
	})
}
