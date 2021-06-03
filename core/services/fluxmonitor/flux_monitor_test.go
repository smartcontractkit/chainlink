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
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitor"
	"github.com/smartcontractkit/chainlink/core/services/log"
	logmocks "github.com/smartcontractkit/chainlink/core/services/log/mocks"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const oracleCount uint8 = 17

type answerSet struct{ latestAnswer, polledAnswer int64 }

var (
	submitHash     = utils.MustHash("submit(uint256,int256)")
	submitSelector = submitHash[:4]
	now            = func() uint64 { return uint64(time.Now().UTC().Unix()) }
	nilOpts        *bind.CallOpts

	makeRoundDataForRoundID = func(roundID uint32) flux_aggregator_wrapper.LatestRoundData {
		return flux_aggregator_wrapper.LatestRoundData{
			RoundId: big.NewInt(int64(roundID)),
		}
	}
	freshContractRoundDataResponse = func() (flux_aggregator_wrapper.LatestRoundData, error) {
		return flux_aggregator_wrapper.LatestRoundData{}, errors.New("No data present")
	}
)

func TestConcreteFluxMonitor_AddJobRemoveJob(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

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
		lb := log.NewBroadcaster(log.NewORM(store.DB), store.EthClient, store.Config, nil)
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

		fm.Close()

		dc.AssertExpectations(t)
	})

	t.Run("does not error or attempt to start a DeviationChecker when receiving a non-Flux Monitor job", func(t *testing.T) {
		job := cltest.NewJobWithRunLogInitiator()
		runManager := new(mocks.RunManager)
		checkerFactory := new(mocks.DeviationCheckerFactory)
		lb := log.NewBroadcaster(log.NewORM(store.DB), store.EthClient, store.Config, nil)
		require.NoError(t, lb.Start())
		fm := fluxmonitor.New(store, runManager, lb)
		fluxmonitor.ExportedSetCheckerFactory(fm, checkerFactory)

		err := fm.Start()
		require.NoError(t, err)
		defer fm.Close()

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
		answersDeviate    bool
		hasPreviousRun    bool
		previousRunStatus models.RunStatus
		expectedToPoll    bool
		expectedToSubmit  bool
	}{
		{
			name:     "eligible",
			eligible: true, connected: true, funded: true, answersDeviate: true,
			expectedToPoll: true, expectedToSubmit: true,
		}, {
			name:     "ineligible",
			eligible: false, connected: true, funded: true, answersDeviate: true,
			expectedToPoll: false, expectedToSubmit: false,
		}, {
			name:     "disconnected",
			eligible: true, connected: false, funded: true, answersDeviate: true,
			expectedToPoll: false, expectedToSubmit: false,
		}, {
			name:     "under funded",
			eligible: true, connected: true, funded: false, answersDeviate: true,
			expectedToPoll: false, expectedToSubmit: false,
		}, {
			name:     "answer undeviated",
			eligible: true, connected: true, funded: true, answersDeviate: false,
			expectedToPoll: true, expectedToSubmit: false,
		}, {
			name:     "previous job run completed",
			eligible: true, connected: true, funded: true, answersDeviate: true,
			hasPreviousRun: true, previousRunStatus: models.RunStatusCompleted,
			expectedToPoll: false, expectedToSubmit: false,
		}, {
			name:     "previous job run in progress",
			eligible: true, connected: true, funded: true, answersDeviate: true,
			hasPreviousRun: true, previousRunStatus: models.RunStatusInProgress,
			expectedToPoll: false, expectedToSubmit: false,
		}, {
			name:     "previous job run cancelled",
			eligible: true, connected: true, funded: true, answersDeviate: true,
			hasPreviousRun: true, previousRunStatus: models.RunStatusCancelled,
			expectedToPoll: false, expectedToSubmit: false,
		}, {
			name:     "previous job run errored",
			eligible: true, connected: true, funded: true, answersDeviate: true,
			hasPreviousRun: true, previousRunStatus: models.RunStatusErrored,
			expectedToPoll: true, expectedToSubmit: true,
		},
	}

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	_, nodeAddr := cltest.MustAddRandomKeyToKeystore(t, store)

	const reportableRoundID = 2
	thresholds := struct{ abs, rel float64 }{0.1, 200}
	deviatedAnswers := answerSet{1, 100}
	undeviatedAnswers := answerSet{100, 101}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {

			answers := undeviatedAnswers
			if test.answersDeviate {
				answers = deviatedAnswers
			}

			rm := new(mocks.RunManager)
			fetcher := new(mocks.Fetcher)
			fluxAggregator := new(mocks.FluxAggregator)
			logBroadcaster := new(logmocks.Broadcaster)
			logBroadcaster.On("IsConnected").Return(test.connected).Once()

			job := cltest.NewJobWithFluxMonitorInitiator()
			initr := job.Initiators[0]
			initr.ID = 1
			require.NoError(t, store.CreateJob(&job))

			if test.hasPreviousRun {
				run := cltest.NewJobRun(job)
				run.Status = test.previousRunStatus
				require.NoError(t, store.CreateJobRun(&run))
				_, err := store.FindOrCreateFluxMonitorRoundStats(initr.Address, reportableRoundID)
				require.NoError(t, err)
				store.UpdateFluxMonitorRoundStats(initr.Address, reportableRoundID, run.ID)
			}

			latestAnswerNoPrecision := answers.latestAnswer * int64(math.Pow10(int(initr.InitiatorParams.Precision)))

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

			roundState := flux_aggregator_wrapper.OracleRoundState{
				RoundId:          reportableRoundID,
				EligibleToSubmit: test.eligible,
				LatestSubmission: big.NewInt(latestAnswerNoPrecision),
				AvailableFunds:   availableFunds,
				PaymentAmount:    paymentAmount,
				OracleCount:      oracleCount,
			}
			fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(0)).Return(roundState, nil).Maybe()
			fluxAggregator.On("LatestRoundData", nilOpts).Return(freshContractRoundDataResponse()).Maybe()

			if test.expectedToPoll {
				fetcher.On("Fetch", mock.Anything, mock.Anything, mock.Anything).Return(decimal.NewFromInt(answers.polledAnswer), nil)
			}

			if test.expectedToSubmit {
				run := cltest.NewJobRun(job)
				require.NoError(t, store.CreateJobRun(&run))

				data, err := models.ParseJSON([]byte(fmt.Sprintf(`{
					"result": "%d",
					"address": "%s",
					"functionSelector": "0x%x",
					"dataPrefix": "0x000000000000000000000000000000000000000000000000000000000000000%d"
				}`, answers.polledAnswer, initr.InitiatorParams.Address.Hex(), submitSelector, reportableRoundID)))
				require.NoError(t, err)

				rm.On("Create", job.ID, &initr, mock.Anything, mock.MatchedBy(func(runRequest *models.RunRequest) bool {
					return reflect.DeepEqual(runRequest.RequestParams.Result.Value(), data.Result.Value())
				})).Return(&run, nil)
			}

			checker, err := fluxmonitor.NewPollingDeviationChecker(
				store,
				fluxAggregator,
				nil,
				logBroadcaster,
				initr,
				nil,
				rm,
				fetcher,
				big.NewInt(0),
				big.NewInt(100000000000),
			)
			require.NoError(t, err)

			oracles := []common.Address{nodeAddr, cltest.NewAddress()}
			fluxAggregator.On("GetOracles", nilOpts).Return(oracles, nil)
			checker.SetOracleAddress()

			checker.ExportedPollIfEligible(thresholds.rel, thresholds.abs)

			logBroadcaster.AssertExpectations(t)
			fluxAggregator.AssertExpectations(t)
			fetcher.AssertExpectations(t)
			rm.AssertExpectations(t)
		})
	}
}

// If the roundState method is unable to communicate with the contract (possibly due to
// incorrect address) then the pollIfEligible method should create a JobSpecErr record
func TestPollingDeviationChecker_PollIfEligible_Creates_JobSpecErr(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	_, nodeAddr := cltest.MustAddRandomKeyToKeystore(t, store)
	oracles := []common.Address{nodeAddr, cltest.NewAddress()}

	rm := new(mocks.RunManager)
	fetcher := new(mocks.Fetcher)
	fluxAggregator := new(mocks.FluxAggregator)
	logBroadcaster := new(logmocks.Broadcaster)
	logBroadcaster.On("IsConnected").Return(true).Once()

	job := cltest.NewJobWithFluxMonitorInitiator()
	initr := job.Initiators[0]
	roundState := flux_aggregator_wrapper.OracleRoundState{}
	require.Len(t, job.Errors, 0)
	err := store.CreateJob(&job)
	require.NoError(t, err)

	fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, mock.Anything).Return(roundState, errors.New("err")).Once()
	checker, err := fluxmonitor.NewPollingDeviationChecker(
		store,
		fluxAggregator,
		nil,
		logBroadcaster,
		initr,
		nil,
		rm,
		fetcher,
		big.NewInt(0),
		big.NewInt(100000000000),
	)
	require.NoError(t, err)

	fluxAggregator.On("GetOracles", nilOpts).Return(oracles, nil)
	require.NoError(t, checker.SetOracleAddress())

	checker.ExportedPollIfEligible(1, 1)

	job, err = store.FindJobWithErrors(job.ID)
	require.NoError(t, err)
	require.Len(t, job.Errors, 1)

	logBroadcaster.AssertExpectations(t)
	fluxAggregator.AssertExpectations(t)
	fetcher.AssertExpectations(t)
	rm.AssertExpectations(t)
}

func TestPollingDeviationChecker_BuffersLogs(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	_, nodeAddr := cltest.MustAddRandomKeyToKeystore(t, store)
	oracles := []common.Address{nodeAddr, cltest.NewAddress()}

	const (
		fetchedValue = 100
	)

	job := cltest.NewJobWithFluxMonitorInitiator()
	initr := job.Initiators[0]
	initr.ID = 1
	initr.PollTimer.Disabled = true
	initr.IdleTimer.Disabled = true
	require.NoError(t, store.CreateJob(&job))

	// Test helpers
	var (
		makeRoundStateForRoundID = func(roundID uint32) flux_aggregator_wrapper.OracleRoundState {
			return flux_aggregator_wrapper.OracleRoundState{
				RoundId:          roundID,
				EligibleToSubmit: true,
				LatestSubmission: big.NewInt(100 * int64(math.Pow10(int(initr.InitiatorParams.Precision)))),
				AvailableFunds:   store.Config.MinimumContractPayment().ToInt(),
				PaymentAmount:    store.Config.MinimumContractPayment().ToInt(),
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
	fluxAggregator.On("LatestRoundData", nilOpts).Return(freshContractRoundDataResponse()).Times(4)
	fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(1)).
		Return(makeRoundStateForRoundID(1), nil).
		Run(func(mock.Arguments) {
			close(chSafeToFillQueue)
			<-chBlock
		}).
		Once()
	fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(3)).Return(makeRoundStateForRoundID(3), nil).Once()
	fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(4)).Return(makeRoundStateForRoundID(4), nil).Once()
	fluxAggregator.On("GetOracles", nilOpts).Return(oracles, nil)
	fluxAggregator.On("Address").Return(initr.Address, nil)

	fetcher := new(mocks.Fetcher)
	fetcher.On("Fetch", mock.Anything, mock.Anything, mock.Anything).Return(decimal.NewFromInt(fetchedValue), nil)

	logBroadcaster := new(logmocks.Broadcaster)
	logBroadcaster.On("Register", mock.Anything, mock.MatchedBy(func(opts log.ListenerOpts) bool {
		return opts.Contract.Address() == initr.Address
	})).Return(func() {})

	rm := new(mocks.RunManager)
	run := cltest.NewJobRun(job)
	require.NoError(t, store.CreateJobRun(&run))

	rm.On("Create", job.ID, &initr, mock.Anything, matchRunRequestForRoundID(1)).Return(&run, nil).Once()
	rm.On("Create", job.ID, &initr, mock.Anything, matchRunRequestForRoundID(3)).Return(&run, nil).Once()
	rm.On("Create", job.ID, &initr, mock.Anything, matchRunRequestForRoundID(4)).Return(&run, nil).Once().
		Run(func(mock.Arguments) { close(chSafeToAssert) })

	checker, err := fluxmonitor.NewPollingDeviationChecker(
		store,
		fluxAggregator,
		nil,
		logBroadcaster,
		initr,
		nil,
		rm,
		fetcher,
		big.NewInt(0),
		big.NewInt(100000000000),
	)
	require.NoError(t, err)

	checker.Start()

	var logBroadcasts []*logmocks.Broadcast

	for i := 1; i <= 4; i++ {
		logBroadcast := new(logmocks.Broadcast)
		logBroadcast.On("DecodedLog").Return(&flux_aggregator_wrapper.FluxAggregatorNewRound{RoundId: big.NewInt(int64(i)), StartedAt: big.NewInt(0)})
		logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
		logBroadcasts = append(logBroadcasts, logBroadcast)
	}

	checker.HandleLog(logBroadcasts[0]) // Get the checker to start processing a log so we can freeze it
	<-chSafeToFillQueue
	checker.HandleLog(logBroadcasts[1]) // This log is evicted from the priority queue
	checker.HandleLog(logBroadcasts[2])
	checker.HandleLog(logBroadcasts[3])

	close(chBlock)
	<-chSafeToAssert

	logBroadcaster.AssertExpectations(t)
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

			_, nodeAddr := cltest.MustAddRandomKeyToKeystore(t, store)
			oracles := []common.Address{nodeAddr, cltest.NewAddress()}

			var (
				fetcher        = new(mocks.Fetcher)
				runManager     = new(mocks.RunManager)
				fluxAggregator = new(mocks.FluxAggregator)
				logBroadcast   = new(logmocks.Broadcast)
				logBroadcaster = new(logmocks.Broadcaster)
			)

			job := cltest.NewJobWithFluxMonitorInitiator()
			initr := job.Initiators[0]
			initr.ID = 1
			initr.PollTimer.Disabled = true
			initr.IdleTimer.Disabled = test.idleTimerDisabled
			initr.IdleTimer.Duration = test.idleDuration

			const fetchedAnswer = 100
			answerBigInt := big.NewInt(fetchedAnswer * int64(math.Pow10(int(initr.InitiatorParams.Precision))))

			fluxAggregator.On("GetOracles", nilOpts).Return(oracles, nil)
			fluxAggregator.On("Address").Return(initr.Address).Maybe()
			logBroadcaster.On("Register", mock.Anything, mock.Anything).Return(func() {})
			logBroadcaster.On("IsConnected").Return(true).Maybe()

			idleDurationOccured := make(chan struct{}, 3)

			fluxAggregator.On("LatestRoundData", nilOpts).Return(freshContractRoundDataResponse()).Maybe()
			if test.expectedToSubmit {
				// performInitialPoll()
				roundState1 := flux_aggregator_wrapper.OracleRoundState{RoundId: 1, EligibleToSubmit: false, LatestSubmission: answerBigInt, StartedAt: now()}
				fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(0)).Return(roundState1, nil).Once()
				// idleDuration 1
				roundState2 := flux_aggregator_wrapper.OracleRoundState{RoundId: 1, EligibleToSubmit: false, LatestSubmission: answerBigInt, StartedAt: now()}
				fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(0)).Return(roundState2, nil).Once().Run(func(args mock.Arguments) {
					idleDurationOccured <- struct{}{}
				})
			}

			deviationChecker, err := fluxmonitor.NewPollingDeviationChecker(
				store,
				fluxAggregator,
				nil,
				logBroadcaster,
				initr,
				nil,
				runManager,
				fetcher,
				big.NewInt(0),
				big.NewInt(100000000000),
			)
			require.NoError(t, err)

			deviationChecker.Start()
			require.Len(t, idleDurationOccured, 0, "no Job Runs created")

			if test.expectedToSubmit {
				require.Eventually(t, func() bool { return len(idleDurationOccured) == 1 }, 3*time.Second, 10*time.Millisecond)

				chBlock := make(chan struct{})
				// NewRound resets the idle timer
				roundState2 := flux_aggregator_wrapper.OracleRoundState{RoundId: 2, EligibleToSubmit: false, LatestSubmission: answerBigInt, StartedAt: now()}
				fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(2)).Return(roundState2, nil).Once().Run(func(args mock.Arguments) {
					close(chBlock)
				})

				decodedLog := flux_aggregator_wrapper.FluxAggregatorNewRound{RoundId: big.NewInt(2), StartedAt: big.NewInt(0)}
				logBroadcast.On("DecodedLog").Return(&decodedLog)
				logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
				logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil).Once()
				deviationChecker.HandleLog(logBroadcast)

				gomega.NewGomegaWithT(t).Eventually(chBlock).Should(gomega.BeClosed())

				// idleDuration 2
				roundState3 := flux_aggregator_wrapper.OracleRoundState{RoundId: 3, EligibleToSubmit: false, LatestSubmission: answerBigInt, StartedAt: now()}
				fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(0)).Return(roundState3, nil).Once().Run(func(args mock.Arguments) {
					idleDurationOccured <- struct{}{}
				})
				require.Eventually(t, func() bool { return len(idleDurationOccured) == 2 }, 3*time.Second, 10*time.Millisecond)
			}

			deviationChecker.Stop()

			if !test.expectedToSubmit {
				require.Len(t, idleDurationOccured, 0)
			}

			logBroadcaster.AssertExpectations(t)
			fetcher.AssertExpectations(t)
			runManager.AssertExpectations(t)
			fluxAggregator.AssertExpectations(t)
		})
	}
}

func TestPollingDeviationChecker_RoundTimeoutCausesPoll_timesOutAtZero(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	_, nodeAddr := cltest.MustAddRandomKeyToKeystore(t, store)
	oracles := []common.Address{nodeAddr, cltest.NewAddress()}

	fetcher := new(mocks.Fetcher)
	runManager := new(mocks.RunManager)
	fluxAggregator := new(mocks.FluxAggregator)
	logBroadcaster := new(logmocks.Broadcaster)

	job := cltest.NewJobWithFluxMonitorInitiator()
	initr := job.Initiators[0]
	initr.ID = 1
	initr.PollTimer.Disabled = true
	initr.IdleTimer.Disabled = true

	ch := make(chan struct{})

	const fetchedAnswer = 100
	answerBigInt := big.NewInt(fetchedAnswer * int64(math.Pow10(int(initr.InitiatorParams.Precision))))
	logBroadcaster.On("Register", mock.Anything, mock.Anything).Return(func() {})

	fluxAggregator.On("LatestRoundData", nilOpts).Return(makeRoundDataForRoundID(1), nil).Maybe()
	fluxAggregator.On("Address").Return(initr.Address).Maybe()
	roundState0 := flux_aggregator_wrapper.OracleRoundState{RoundId: 1, EligibleToSubmit: false, LatestSubmission: answerBigInt, StartedAt: now()}
	fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(1)).Return(roundState0, nil).Once() // initialRoundState()
	fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(0)).Return(flux_aggregator_wrapper.OracleRoundState{
		RoundId:          1,
		EligibleToSubmit: false,
		LatestSubmission: answerBigInt,
		StartedAt:        0,
		Timeout:          0,
	}, nil).
		Run(func(mock.Arguments) { close(ch) }).
		Once()

	deviationChecker, err := fluxmonitor.NewPollingDeviationChecker(
		store,
		fluxAggregator,
		nil,
		logBroadcaster,
		initr,
		nil,
		runManager,
		fetcher,
		big.NewInt(0),
		big.NewInt(100000000000),
	)
	require.NoError(t, err)

	fluxAggregator.On("GetOracles", nilOpts).Return(oracles, nil)

	deviationChecker.SetOracleAddress()
	deviationChecker.ExportedRoundState()
	deviationChecker.Start()

	gomega.NewGomegaWithT(t).Eventually(ch).Should(gomega.BeClosed())

	deviationChecker.Stop()

	logBroadcaster.AssertExpectations(t)
	fetcher.AssertExpectations(t)
	runManager.AssertExpectations(t)
	fluxAggregator.AssertExpectations(t)
}

func TestPollingDeviationChecker_UsesPreviousRoundStateOnStartup_RoundTimeout(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	_, nodeAddr := cltest.MustAddRandomKeyToKeystore(t, store)
	oracles := []common.Address{nodeAddr, cltest.NewAddress()}

	fetcher := new(mocks.Fetcher)
	runManager := new(mocks.RunManager)
	logBroadcaster := new(logmocks.Broadcaster)

	job := cltest.NewJobWithFluxMonitorInitiator()
	initr := job.Initiators[0]
	initr.PollTimer.Disabled = true
	initr.IdleTimer.Disabled = true

	tests := []struct {
		name             string
		timeout          uint64
		expectedToSubmit bool
	}{
		{"active round exists - round will time out", 2, true},
		{"active round exists - round will not time out", 100, false},
		{"no active round", 0, false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			fluxAggregator := new(mocks.FluxAggregator)

			logBroadcaster.On("Register", mock.Anything, mock.Anything).Return(func() {})
			logBroadcaster.On("IsConnected").Return(true).Maybe()

			fluxAggregator.On("Address").Return(initr.Address).Maybe()
			fluxAggregator.On("GetOracles", nilOpts).Return(oracles, nil)

			fluxAggregator.On("LatestRoundData", nilOpts).Return(makeRoundDataForRoundID(1), nil).Maybe()
			fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(1)).Return(flux_aggregator_wrapper.OracleRoundState{
				RoundId:          1,
				EligibleToSubmit: false,
				StartedAt:        now(),
				Timeout:          test.timeout,
			}, nil).Once()

			// 2nd roundstate call means round timer triggered
			chRoundState := make(chan struct{})
			fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(0)).Return(flux_aggregator_wrapper.OracleRoundState{
				RoundId:          1,
				EligibleToSubmit: false,
			}, nil).
				Run(func(mock.Arguments) { close(chRoundState) }).
				Maybe()

			deviationChecker, err := fluxmonitor.NewPollingDeviationChecker(
				store,
				fluxAggregator,
				nil,
				logBroadcaster,
				initr,
				nil,
				runManager,
				fetcher,
				big.NewInt(0),
				big.NewInt(100000000000),
			)
			require.NoError(t, err)

			deviationChecker.Start()

			if test.expectedToSubmit {
				gomega.NewGomegaWithT(t).Eventually(chRoundState).Should(gomega.BeClosed())
			} else {
				gomega.NewGomegaWithT(t).Consistently(chRoundState).ShouldNot(gomega.BeClosed())
			}

			deviationChecker.Stop()
			logBroadcaster.AssertExpectations(t)
			fluxAggregator.AssertExpectations(t)
		})
	}
}

func TestPollingDeviationChecker_UsesPreviousRoundStateOnStartup_IdleTimer(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	_, nodeAddr := cltest.MustAddRandomKeyToKeystore(t, store)
	oracles := []common.Address{nodeAddr, cltest.NewAddress()}

	fetcher := new(mocks.Fetcher)
	runManager := new(mocks.RunManager)
	logBroadcaster := new(logmocks.Broadcaster)

	job := cltest.NewJobWithFluxMonitorInitiator()
	initr := job.Initiators[0]
	initr.PollTimer.Disabled = true
	initr.IdleTimer.Disabled = false

	almostExpired := time.Now().
		Add(initr.IdleTimer.Duration.Duration() * -1).
		Add(2 * time.Second).
		Unix()

	tests := []struct {
		name             string
		startedAt        uint64
		expectedToSubmit bool
	}{
		{"active round exists - idleTimer about to expired", uint64(almostExpired), true},
		{"active round exists - idleTimer will not expire", 100, false},
		{"no active round", 0, false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			fluxAggregator := new(mocks.FluxAggregator)

			logBroadcaster.On("Register", mock.Anything, mock.Anything).Return(func() {})
			logBroadcaster.On("IsConnected").Return(true).Maybe()

			fluxAggregator.On("Address").Return(initr.Address).Maybe()
			fluxAggregator.On("GetOracles", nilOpts).Return(oracles, nil)
			fluxAggregator.On("LatestRoundData", nilOpts).Return(makeRoundDataForRoundID(1), nil).Maybe()
			// first roundstate in setInitialTickers()
			fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(1)).Return(flux_aggregator_wrapper.OracleRoundState{
				RoundId:          1,
				EligibleToSubmit: false,
				StartedAt:        test.startedAt,
				Timeout:          10000, // round won't time out
			}, nil).Once()

			// 2nd roundstate in performInitialPoll()
			roundState := flux_aggregator_wrapper.OracleRoundState{RoundId: 1, EligibleToSubmit: false}
			fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(0)).Return(roundState, nil).Once()

			// 3rd roundState call means idleTimer triggered
			chRoundState := make(chan struct{})
			fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(0)).Return(roundState, nil).
				Run(func(mock.Arguments) { close(chRoundState) }).
				Maybe()

			deviationChecker, err := fluxmonitor.NewPollingDeviationChecker(
				store,
				fluxAggregator,
				nil,
				logBroadcaster,
				initr,
				nil,
				runManager,
				fetcher,
				big.NewInt(0),
				big.NewInt(100000000000),
			)
			require.NoError(t, err)

			deviationChecker.Start()

			if test.expectedToSubmit {
				gomega.NewGomegaWithT(t).Eventually(chRoundState).Should(gomega.BeClosed())
			} else {
				gomega.NewGomegaWithT(t).Consistently(chRoundState).ShouldNot(gomega.BeClosed())
			}

			deviationChecker.Stop()
			fluxAggregator.AssertExpectations(t)
		})
	}
}

func TestPollingDeviationChecker_RoundTimeoutCausesPoll_timesOutNotZero(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	_, nodeAddr := cltest.MustAddRandomKeyToKeystore(t, store)
	oracles := []common.Address{nodeAddr, cltest.NewAddress()}

	fetcher := new(mocks.Fetcher)
	runManager := new(mocks.RunManager)
	fluxAggregator := new(mocks.FluxAggregator)
	logBroadcaster := new(logmocks.Broadcaster)

	job := cltest.NewJobWithFluxMonitorInitiator()
	initr := job.Initiators[0]
	initr.ID = 1
	initr.PollTimer.Disabled = true
	initr.IdleTimer.Disabled = true

	const fetchedAnswer = 100
	answerBigInt := big.NewInt(fetchedAnswer * int64(math.Pow10(int(initr.InitiatorParams.Precision))))

	chRoundState1 := make(chan struct{})
	chRoundState2 := make(chan struct{})

	logBroadcaster.On("Register", mock.Anything, mock.Anything).Return(func() {})
	logBroadcaster.On("IsConnected").Return(true).Maybe()

	fluxAggregator.On("Address").Return(initr.Address).Maybe()
	fluxAggregator.On("GetOracles", nilOpts).Return(oracles, nil)
	fluxAggregator.On("LatestRoundData", nilOpts).Return(makeRoundDataForRoundID(1), nil).Maybe()
	fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(1)).Return(flux_aggregator_wrapper.OracleRoundState{
		RoundId:          1,
		EligibleToSubmit: false,
		LatestSubmission: answerBigInt,
		StartedAt:        now(),
		Timeout:          uint64(1000000),
	}, nil).Once()

	startedAt := uint64(time.Now().Unix())
	timeout := uint64(3)
	fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(0)).Return(flux_aggregator_wrapper.OracleRoundState{
		RoundId:          1,
		EligibleToSubmit: false,
		LatestSubmission: answerBigInt,
		StartedAt:        startedAt,
		Timeout:          timeout,
	}, nil).Once().
		Run(func(mock.Arguments) { close(chRoundState1) }).
		Once()
	fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(0)).Return(flux_aggregator_wrapper.OracleRoundState{
		RoundId:          1,
		EligibleToSubmit: false,
		LatestSubmission: answerBigInt,
		StartedAt:        startedAt,
		Timeout:          timeout,
	}, nil).Once().
		Run(func(mock.Arguments) { close(chRoundState2) }).
		Once()

	deviationChecker, err := fluxmonitor.NewPollingDeviationChecker(
		store,
		fluxAggregator,
		nil,
		logBroadcaster,
		initr,
		nil,
		runManager,
		fetcher,
		big.NewInt(0),
		big.NewInt(100000000000),
	)
	require.NoError(t, err)
	deviationChecker.Start()

	logBroadcast := new(logmocks.Broadcast)
	logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
	logBroadcast.On("DecodedLog").Return(&flux_aggregator_wrapper.FluxAggregatorNewRound{RoundId: big.NewInt(0), StartedAt: big.NewInt(time.Now().UTC().Unix())})
	logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	deviationChecker.HandleLog(logBroadcast)

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

	type testCase struct {
		funded        bool
		eligible      bool
		startedBySelf bool
		duplicateLog  bool
		runStatus     models.RunStatus
		roundIDCase
		answerCase
	}

	// generate all permutations of test cases
	tests := []testCase{}
	duplicateLogOptions := []bool{true, false}
	runStatusOptions := []models.RunStatus{
		models.RunStatusCompleted, models.RunStatusCancelled, models.RunStatusErrored,
		models.RunStatusInProgress, models.RunStatusUnstarted, models.RunStatusPendingOutgoingConfirmations,
	}
	fundedOptions := []bool{true, false}
	eligibleOptions := []bool{true, false}
	startedBySelfOptions := []bool{true, false}
	roundIDCaseOptions := []roundIDCase{fetched_lt_log, fetched_gt_log, fetched_eq_log}
	answerCaseOptions := []answerCase{deviationThresholdExceeded, deviationThresholdNotExceeded}
	for _, funded := range fundedOptions {
		for _, eligible := range eligibleOptions {
			for _, startedBySelf := range startedBySelfOptions {
				for _, duplicateLog := range duplicateLogOptions {
					for _, runStatus := range runStatusOptions {
						for _, roundIDCase := range roundIDCaseOptions {
							for _, answerCase := range answerCaseOptions {
								newTestCase := testCase{funded, eligible, startedBySelf, duplicateLog, runStatus, roundIDCase, answerCase}
								tests = append(tests, newTestCase)
							}
						}
					}
				}
			}
		}
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

			_, nodeAddr := cltest.MustAddRandomKeyToKeystore(t, store)
			oracles := []common.Address{nodeAddr, cltest.NewAddress()}

			previousSubmissionReorged := test.duplicateLog &&
				(test.runStatus == models.RunStatusCompleted || test.runStatus == models.RunStatusErrored)
			expectedToFetchRoundState := previousSubmissionReorged || !test.startedBySelf
			expectedToPoll := expectedToFetchRoundState && test.eligible && test.funded && test.logRoundID >= int64(test.fetchedReportableRoundID)
			expectedToSubmit := expectedToPoll

			job := cltest.NewJobWithFluxMonitorInitiator()
			initr := job.Initiators[0]
			initr.ID = 1
			initr.PollTimer.Disabled = true
			initr.IdleTimer.Disabled = true
			require.NoError(t, store.CreateJob(&job))

			if test.duplicateLog {
				jobRun := cltest.NewJobRun(job)
				jobRun.Status = test.runStatus
				require.NoError(t, store.CreateJobRun(&jobRun))
				err := store.UpdateFluxMonitorRoundStats(initr.Address, uint32(test.logRoundID), jobRun.ID)
				require.NoError(t, err)
			}

			rm := new(mocks.RunManager)
			fetcher := new(mocks.Fetcher)
			fluxAggregator := new(mocks.FluxAggregator)
			logBroadcaster := new(logmocks.Broadcaster)
			logBroadcaster.On("IsConnected").Return(true).Maybe()

			paymentAmount := store.Config.MinimumContractPayment().ToInt()
			var availableFunds *big.Int
			if test.funded {
				availableFunds = big.NewInt(1).Mul(paymentAmount, big.NewInt(1000))
			} else {
				availableFunds = big.NewInt(1)
			}

			if expectedToFetchRoundState {
				fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(test.logRoundID)).Return(flux_aggregator_wrapper.OracleRoundState{
					RoundId:          test.fetchedReportableRoundID,
					LatestSubmission: big.NewInt(test.latestAnswer * int64(math.Pow10(int(initr.InitiatorParams.Precision)))),
					EligibleToSubmit: test.eligible,
					AvailableFunds:   availableFunds,
					PaymentAmount:    paymentAmount,
					OracleCount:      oracleCount,
				}, nil).Once()
			}

			if expectedToPoll {
				fetcher.On("Fetch", mock.Anything, mock.Anything, mock.Anything).Return(decimal.NewFromInt(test.polledAnswer), nil).Once()
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
				nil,
				logBroadcaster,
				initr,
				nil,
				rm,
				fetcher,
				big.NewInt(0),
				big.NewInt(0),
			)
			require.NoError(t, err)

			fluxAggregator.On("Address").Return(initr.Address).Maybe()
			fluxAggregator.On("GetOracles", nilOpts).Return(oracles, nil)
			checker.SetOracleAddress()

			var startedBy common.Address
			if test.startedBySelf {
				startedBy = nodeAddr
			}
			checker.ExportedRespondToNewRoundLog(&flux_aggregator_wrapper.FluxAggregatorNewRound{
				RoundId:   big.NewInt(test.logRoundID),
				StartedBy: startedBy,
				StartedAt: big.NewInt(0),
			})

			fluxAggregator.AssertExpectations(t)
			logBroadcaster.AssertExpectations(t)
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
	logBroadcaster := new(logmocks.Broadcaster)
	logBroadcaster.On("IsConnected").Return(true).Maybe()

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
				nil,
				logBroadcaster,
				initr,
				minJobPayment,
				rm,
				fetcher,
				big.NewInt(0),
				big.NewInt(100000000000),
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

			state := flux_aggregator_wrapper.OracleRoundState{
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

	log := flux_aggregator_wrapper.FluxAggregatorNewRound{}
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

	log := flux_aggregator_wrapper.FluxAggregatorNewRound{StartedAt: big.NewInt(10)}
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

	log := flux_aggregator_wrapper.FluxAggregatorNewRound{StartedAt: big.NewInt(10)}
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
	log := flux_aggregator_wrapper.FluxAggregatorNewRound{StartedAt: &startedAt}
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

	log := flux_aggregator_wrapper.FluxAggregatorNewRound{StartedAt: big.NewInt(40)}
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

	var (
		p                = cltest.NewPollingDeviationChecker(t, store)
		logBroadcast     = new(logmocks.Broadcast)
		logNewRound      *flux_aggregator_wrapper.FluxAggregatorNewRound
		logAnswerUpdated *flux_aggregator_wrapper.FluxAggregatorAnswerUpdated
		randomType       interface{}
	)

	logBroadcast.On("DecodedLog").Return(logNewRound).Once()
	assert.NotPanics(t, func() {
		p.HandleLog(logBroadcast)
	})

	logBroadcast.On("DecodedLog").Return(logAnswerUpdated).Once()
	assert.NotPanics(t, func() {
		p.HandleLog(logBroadcast)
	})

	logBroadcast.On("DecodedLog").Return(randomType).Once()
	assert.NotPanics(t, func() {
		p.HandleLog(logBroadcast)
	})
}

func TestFluxMonitor_IdleTimer(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	fluxAggregator := new(mocks.FluxAggregator)
	runManager := new(mocks.RunManager)
	fetcher := new(mocks.Fetcher)
	initr := models.Initiator{
		JobSpecID: models.NewJobID(),
		InitiatorParams: models.InitiatorParams{
			IdleTimer: models.IdleTimerConfig{
				Disabled: false,
				Duration: models.MustMakeDuration(10 * time.Millisecond),
			},
			PollTimer: models.PollTimerConfig{
				Disabled: true,
			},
		},
	}
	lb := new(logmocks.Broadcaster)
	lb.On("Register", mock.Anything, mock.Anything).Return(func() {})
	lb.On("IsConnected").Return(true)
	fluxAggregator.On("GetOracles", mock.Anything).Return([]common.Address{}, nil)
	fluxAggregator.On("LatestRoundData", mock.Anything).Return(
		flux_aggregator_wrapper.LatestRoundData{RoundId: big.NewInt(10), StartedAt: nil}, nil)

	// By returning this old round state started at, we stop the idle timer from getting reset.
	startedAtTs := big.NewInt(time.Now().Unix() - 10)
	// Normally there are 2 oracle round state calls upon startup.
	fluxAggregator.On("OracleRoundState", mock.Anything, mock.Anything, mock.Anything).Return(
		flux_aggregator_wrapper.OracleRoundState{EligibleToSubmit: false, RoundId: 10, StartedAt: startedAtTs.Uint64()}, nil).Times(2)
	done := make(chan struct{})
	// To get a 3rd call we need the idle timer to fire
	fluxAggregator.On("OracleRoundState", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		done <- struct{}{}
	}).Return(
		flux_aggregator_wrapper.OracleRoundState{EligibleToSubmit: false, RoundId: 10, StartedAt: startedAtTs.Uint64()}, nil)

	checker, err := fluxmonitor.NewPollingDeviationChecker(store, fluxAggregator, nil, lb, initr, nil, runManager, fetcher, big.NewInt(0), big.NewInt(100000000000))
	require.NoError(t, err)
	checker.Start()
	defer checker.Stop()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Error("idle timer did not fire as expected")
	}
}

func TestFluxMonitor_ConsumeLogBroadcast_Happy(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	p := cltest.NewPollingDeviationChecker(t, store)
	p.ExportedFluxAggregator().(*mocks.FluxAggregator).
		On("OracleRoundState", nilOpts, mock.Anything, mock.Anything).
		Return(flux_aggregator_wrapper.OracleRoundState{RoundId: 123}, nil)
	p.ExportedFluxAggregator().(*mocks.FluxAggregator).
		On("Address").
		Return(cltest.NewAddress())

	logBroadcast := new(logmocks.Broadcast)
	p.ExportedLogBroadcaster().On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
	logBroadcast.On("DecodedLog").Return(&flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{})
	p.ExportedLogBroadcaster().On("MarkConsumed", mock.Anything, mock.Anything).Return(nil).Once()

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

			logBroadcast := new(logmocks.Broadcast)
			p.ExportedLogBroadcaster().On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(test.consumed, test.err).Once()

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

		_, nodeAddr := cltest.MustAddRandomKeyToKeystore(t, store)
		oracles := []common.Address{nodeAddr, cltest.NewAddress()}

		job := cltest.NewJobWithFluxMonitorInitiator()
		initr := job.Initiators[0]
		initr.ID = 1
		initr.PollTimer.Disabled = true
		initr.IdleTimer.Disabled = true
		run := cltest.NewJobRun(job)

		var (
			rm             = new(mocks.RunManager)
			fetcher        = new(mocks.Fetcher)
			fluxAggregator = new(mocks.FluxAggregator)
			logBroadcaster = new(logmocks.Broadcaster)

			paymentAmount  = store.Config.MinimumContractPayment().ToInt()
			availableFunds = big.NewInt(1).Mul(paymentAmount, big.NewInt(1000))
		)
		logBroadcaster.On("IsConnected").Return(true).Maybe()

		const (
			roundID = 3
			answer  = 100
		)

		checker, err := fluxmonitor.NewPollingDeviationChecker(
			store,
			fluxAggregator,
			nil,
			logBroadcaster,
			initr,
			nil,
			rm,
			fetcher,
			big.NewInt(0),
			big.NewInt(100000000000),
		)
		require.NoError(t, err)

		fluxAggregator.On("Address").Return(initr.Address)
		fluxAggregator.On("GetOracles", nilOpts).Return(oracles, nil)
		checker.SetOracleAddress()

		fluxAggregator.On("LatestRoundData", nilOpts).Return(freshContractRoundDataResponse()).Maybe()
		// Fire off the NewRound log, which the node should respond to
		fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(roundID)).
			Return(flux_aggregator_wrapper.OracleRoundState{
				RoundId:          roundID,
				LatestSubmission: big.NewInt(answer),
				EligibleToSubmit: true,
				AvailableFunds:   availableFunds,
				PaymentAmount:    paymentAmount,
				OracleCount:      1,
			}, nil).
			Once()
		fetcher.On("Fetch", mock.Anything, mock.Anything, mock.Anything).
			Return(decimal.NewFromInt(answer), nil).
			Once()
		rm.On("Create", job.ID, &initr, mock.Anything, mock.Anything).
			Return(&run, nil).
			Once()
		checker.ExportedRespondToNewRoundLog(&flux_aggregator_wrapper.FluxAggregatorNewRound{
			RoundId:   big.NewInt(roundID),
			StartedAt: big.NewInt(0),
		})

		g := gomega.NewGomegaWithT(t)
		expectation := func() bool {
			jrs, err := store.JobRunsFor(job.ID)
			require.NoError(t, err)
			return len(jrs) == 1
		}
		g.Eventually(expectation, cltest.DBWaitTimeout, cltest.DBPollingInterval)
		g.Consistently(expectation, cltest.DBWaitTimeout, cltest.DBPollingInterval)

		// Now force the node to try to poll and ensure it does not respond this time
		fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(0)).
			Return(flux_aggregator_wrapper.OracleRoundState{
				RoundId:          roundID,
				LatestSubmission: big.NewInt(answer),
				EligibleToSubmit: true,
				AvailableFunds:   availableFunds,
				PaymentAmount:    paymentAmount,
				OracleCount:      1,
			}, nil).
			Once()
		checker.ExportedPollIfEligible(0, 0)

		rm.AssertExpectations(t)
		fetcher.AssertExpectations(t)
		fluxAggregator.AssertExpectations(t)
		logBroadcaster.AssertExpectations(t)
	})

	t.Run("when poll ticker fires, then NewRound log arrives", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()

		_, nodeAddr := cltest.MustAddRandomKeyToKeystore(t, store)
		oracles := []common.Address{nodeAddr, cltest.NewAddress()}

		job := cltest.NewJobWithFluxMonitorInitiator()
		initr := job.Initiators[0]
		initr.ID = 1
		initr.PollTimer.Disabled = true
		initr.IdleTimer.Disabled = true
		run := cltest.NewJobRun(job)

		rm := new(mocks.RunManager)
		fetcher := new(mocks.Fetcher)
		fluxAggregator := new(mocks.FluxAggregator)
		logBroadcaster := new(logmocks.Broadcaster)
		logBroadcaster.On("IsConnected").Return(true).Maybe()

		paymentAmount := store.Config.MinimumContractPayment().ToInt()
		availableFunds := big.NewInt(1).Mul(paymentAmount, big.NewInt(1000))

		const (
			roundID = 3
			answer  = 100
		)

		checker, err := fluxmonitor.NewPollingDeviationChecker(
			store,
			fluxAggregator,
			nil,
			logBroadcaster,
			initr,
			nil,
			rm,
			fetcher,
			big.NewInt(0),
			big.NewInt(100000000000),
		)
		require.NoError(t, err)

		fluxAggregator.On("LatestRoundData", nilOpts).Return(flux_aggregator_wrapper.LatestRoundData{Answer: big.NewInt(100), UpdatedAt: big.NewInt(1616447984)}, nil).Maybe()
		// First, force the node to try to poll, which should result in a submission
		fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(0)).
			Return(flux_aggregator_wrapper.OracleRoundState{
				RoundId:          roundID,
				LatestSubmission: big.NewInt(answer),
				EligibleToSubmit: true,
				AvailableFunds:   availableFunds,
				PaymentAmount:    paymentAmount,
				OracleCount:      1,
			}, nil).
			Once()
		md, _ := models.MarshalBridgeMetaData(big.NewInt(100), big.NewInt(1616447984))
		fetcher.On("Fetch", mock.Anything, md, mock.Anything).
			Return(decimal.NewFromInt(answer), nil).
			Once()
		rm.On("Create", job.ID, &initr, mock.Anything, mock.Anything).
			Return(&run, nil).
			Once()
		fluxAggregator.On("Address").Return(initr.Address)
		fluxAggregator.On("GetOracles", nilOpts).Return(oracles, nil)
		checker.SetOracleAddress()
		checker.ExportedPollIfEligible(0, 0)

		// Now fire off the NewRound log and ensure it does not respond this time
		checker.ExportedRespondToNewRoundLog(&flux_aggregator_wrapper.FluxAggregatorNewRound{
			RoundId:   big.NewInt(roundID),
			StartedAt: big.NewInt(0),
		})

		rm.AssertExpectations(t)
		fetcher.AssertExpectations(t)
		fluxAggregator.AssertExpectations(t)
		logBroadcaster.AssertExpectations(t)
	})
}

func TestFluxMonitor_PollingDeviationChecker_IsFlagLowered(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		getFlagsResult []bool
		expected       bool
	}{
		{"both lowered", []bool{false, false}, true},
		{"global lowered", []bool{false, true}, true},
		{"contract lowered", []bool{true, false}, true},
		{"both raised", []bool{true, true}, false},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			config, configCleanup := cltest.NewConfig(t)
			defer configCleanup()
			store, storeCleanup := cltest.NewStoreWithConfig(t, config)
			defer storeCleanup()

			var (
				fluxAggregator = new(mocks.FluxAggregator)
				rm             = new(mocks.RunManager)
				fetcher        = new(mocks.Fetcher)
				logBroadcaster = new(logmocks.Broadcaster)
				flagsContract  = new(mocks.Flags)

				job   = cltest.NewJobWithFluxMonitorInitiator()
				initr = job.Initiators[0]
			)
			logBroadcaster.On("IsConnected").Return(true).Maybe()

			flagsContract.On("GetFlags", mock.Anything, mock.Anything).
				Run(func(args mock.Arguments) {
					require.Equal(t, []common.Address{utils.ZeroAddress, initr.Address}, args.Get(1).([]common.Address))
				}).
				Return(test.getFlagsResult, nil)

			checker, err := fluxmonitor.NewPollingDeviationChecker(
				store,
				fluxAggregator,
				flagsContract,
				logBroadcaster,
				initr,
				nil,
				rm,
				fetcher,
				big.NewInt(0),
				big.NewInt(100000000000),
			)
			require.NoError(t, err)

			result, err := checker.ExportedIsFlagLowered()
			require.NoError(t, err)
			require.Equal(t, test.expected, result)

			flagsContract.AssertExpectations(t)
		})
	}
}
