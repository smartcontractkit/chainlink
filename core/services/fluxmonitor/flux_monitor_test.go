package fluxmonitor_test

import (
	"fmt"
	"math"
	"math/big"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/eth/contracts"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitor"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	updateAnswerHash     = utils.MustHash("updateAnswer(uint256,int256)")
	updateAnswerSelector = updateAnswerHash[:4]
)

func ensureAccount(t *testing.T, store *store.Store) common.Address {
	t.Helper()
	auth := cmd.TerminalKeyStoreAuthenticator{Prompter: &cltest.MockCountingPrompter{T: t}}
	_, err := auth.Authenticate(store, "somepassword")
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

			fm := fluxmonitor.New(store, runManager)
			logBroadcaster := fm.(fluxmonitor.MockableLogBroadcaster).MockLogBroadcaster()

			err := fm.Start()
			assert.NoError(t, err)
			assert.Equal(t, test.wantStarted, logBroadcaster.Started)
		})
	}
}

func TestConcreteFluxMonitor_AddJobRemoveJob(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	txm := new(mocks.TxManager)
	store.TxManager = txm
	txm.On("GetBlockHeight").Return(uint64(123), nil)

	t.Run("starts and stops DeviationCheckers when jobs are added and removed", func(t *testing.T) {
		job := cltest.NewJobWithFluxMonitorInitiator()
		runManager := new(mocks.RunManager)
		started := make(chan struct{}, 1)

		dc := new(mocks.DeviationChecker)
		dc.On("Start", mock.Anything, mock.Anything).Return(nil).Run(func(mock.Arguments) {
			started <- struct{}{}
		})

		checkerFactory := new(mocks.DeviationCheckerFactory)
		checkerFactory.On("New", job.Initiators[0], runManager, store.ORM, store.Config.DefaultHTTPTimeout()).Return(dc, nil)
		fm := fluxmonitor.New(store, runManager)
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
		fm := fluxmonitor.New(store, runManager)
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
	tests := []struct {
		name             string
		eligible         bool
		connected        bool
		funded           bool
		threshold        float64
		latestAnswer     int64
		polledAnswer     int64
		expectedToPoll   bool
		expectedToSubmit bool
	}{
		{"eligible, connected, funded, threshold > 0, answers deviate", true, true, true, 0.1, 1, 100, true, true},
		{"eligible, connected, funded, threshold > 0, answers do not deviate", true, true, true, 0.1, 100, 100, true, false},
		{"eligible, connected, funded, threshold == 0, answers deviate", true, true, true, 0, 1, 100, true, true},
		{"eligible, connected, funded, threshold == 0, answers do not deviate", true, true, true, 0, 1, 100, true, true},

		{"eligible, disconnected, funded, threshold > 0, answers deviate", true, false, true, 0.1, 1, 100, false, false},
		{"eligible, disconnected, funded, threshold > 0, answers do not deviate", true, false, true, 0.1, 100, 100, false, false},
		{"eligible, disconnected, funded, threshold == 0, answers deviate", true, false, true, 0, 1, 100, false, false},
		{"eligible, disconnected, funded, threshold == 0, answers do not deviate", true, false, true, 0, 1, 100, false, false},

		{"ineligible, connected, funded, threshold > 0, answers deviate", false, true, true, 0.1, 1, 100, false, false},
		{"ineligible, connected, funded, threshold > 0, answers do not deviate", false, true, true, 0.1, 100, 100, false, false},
		{"ineligible, connected, funded, threshold == 0, answers deviate", false, true, true, 0, 1, 100, false, false},
		{"ineligible, connected, funded, threshold == 0, answers do not deviate", false, true, true, 0, 1, 100, false, false},

		{"ineligible, disconnected, funded, threshold > 0, answers deviate", false, false, true, 0.1, 1, 100, false, false},
		{"ineligible, disconnected, funded, threshold > 0, answers do not deviate", false, false, true, 0.1, 100, 100, false, false},
		{"ineligible, disconnected, funded, threshold == 0, answers deviate", false, false, true, 0, 1, 100, false, false},
		{"ineligible, disconnected, funded, threshold == 0, answers do not deviate", false, false, true, 0, 1, 100, false, false},

		{"eligible, connected, underfunded, threshold > 0, answers deviate", true, true, false, 0.1, 1, 100, false, false},
		{"eligible, connected, underfunded, threshold > 0, answers do not deviate", true, true, false, 0.1, 100, 100, false, false},
		{"eligible, connected, underfunded, threshold == 0, answers deviate", true, true, false, 0, 1, 100, false, false},
		{"eligible, connected, underfunded, threshold == 0, answers do not deviate", true, true, false, 0, 1, 100, false, false},

		{"eligible, disconnected, underfunded, threshold > 0, answers deviate", true, false, false, 0.1, 1, 100, false, false},
		{"eligible, disconnected, underfunded, threshold > 0, answers do not deviate", true, false, false, 0.1, 100, 100, false, false},
		{"eligible, disconnected, underfunded, threshold == 0, answers deviate", true, false, false, 0, 1, 100, false, false},
		{"eligible, disconnected, underfunded, threshold == 0, answers do not deviate", true, false, false, 0, 1, 100, false, false},

		{"ineligible, connected, underfunded, threshold > 0, answers deviate", false, true, false, 0.1, 1, 100, false, false},
		{"ineligible, connected, underfunded, threshold > 0, answers do not deviate", false, true, false, 0.1, 100, 100, false, false},
		{"ineligible, connected, underfunded, threshold == 0, answers deviate", false, true, false, 0, 1, 100, false, false},
		{"ineligible, connected, underfunded, threshold == 0, answers do not deviate", false, true, false, 0, 1, 100, false, false},

		{"ineligible, disconnected, underfunded, threshold > 0, answers deviate", false, false, false, 0.1, 1, 100, false, false},
		{"ineligible, disconnected, underfunded, threshold > 0, answers do not deviate", false, false, false, 0.1, 100, 100, false, false},
		{"ineligible, disconnected, underfunded, threshold == 0, answers deviate", false, false, false, 0, 1, 100, false, false},
		{"ineligible, disconnected, underfunded, threshold == 0, answers do not deviate", false, false, false, 0, 1, 100, false, false},
	}

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	nodeAddr := ensureAccount(t, store)

	for _, test := range tests {
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
				availableFunds = minPayment
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
			}
			fluxAggregator.On("RoundState", nodeAddr).Return(roundState, nil).Maybe()

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
				}`, test.polledAnswer, initr.InitiatorParams.Address.Hex(), updateAnswerSelector, reportableRoundID)))
				require.NoError(t, err)

				rm.On("Create", job.ID, &initr, mock.Anything, mock.MatchedBy(func(runRequest *models.RunRequest) bool {
					return reflect.DeepEqual(runRequest.RequestParams.Result.Value(), data.Result.Value())
				})).Return(&run, nil)

				fluxAggregator.On("GetMethodID", "updateAnswer").Return(updateAnswerSelector, nil)
			}

			checker, err := fluxmonitor.NewPollingDeviationChecker(store, fluxAggregator, initr, rm, fetcher, time.Second)
			require.NoError(t, err)

			if test.connected {
				checker.OnConnect()
			}

			checker.ExportedPollIfEligible(test.threshold)

			fluxAggregator.AssertExpectations(t)
			fetcher.AssertExpectations(t)
			rm.AssertExpectations(t)
		})
	}
}

func TestPollingDeviationChecker_TriggerIdleTimeThreshold(t *testing.T) {

	tests := []struct {
		name             string
		idleThreshold    time.Duration
		expectedToSubmit bool
	}{
		{"no idleThreshold", 0, false},
		{"idleThreshold > 0", 10 * time.Millisecond, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			nodeAddr := ensureAccount(t, store)

			fetcher := new(mocks.Fetcher)
			runManager := new(mocks.RunManager)
			fluxAggregator := new(mocks.FluxAggregator)

			job := cltest.NewJobWithFluxMonitorInitiator()
			initr := job.Initiators[0]
			initr.ID = 1
			initr.PollingInterval = models.Duration(math.MaxInt64)
			initr.IdleThreshold = models.Duration(test.idleThreshold)

			const fetchedAnswer = 100
			answerBigInt := big.NewInt(fetchedAnswer * int64(math.Pow10(int(initr.InitiatorParams.Precision))))

			fluxAggregator.On("SubscribeToLogs", mock.Anything).Return(true, eth.UnsubscribeFunc(func() {}), nil)

			roundState1 := contracts.FluxAggregatorRoundState{ReportableRoundID: 1, EligibleToSubmit: false, LatestAnswer: answerBigInt} // Initial poll
			roundState2 := contracts.FluxAggregatorRoundState{ReportableRoundID: 2, EligibleToSubmit: false, LatestAnswer: answerBigInt} // idleThreshold 1
			roundState3 := contracts.FluxAggregatorRoundState{ReportableRoundID: 3, EligibleToSubmit: false, LatestAnswer: answerBigInt} // NewRound
			roundState4 := contracts.FluxAggregatorRoundState{ReportableRoundID: 4, EligibleToSubmit: false, LatestAnswer: answerBigInt} // idleThreshold 2

			idleThresholdOccured := make(chan struct{}, 3)

			fluxAggregator.On("RoundState", nodeAddr).Return(roundState1, nil).Once() // Initial poll
			if test.expectedToSubmit {
				// idleThreshold 1
				fluxAggregator.On("RoundState", nodeAddr).Return(roundState2, nil).Once().Run(func(args mock.Arguments) { idleThresholdOccured <- struct{}{} })
				// NewRound
				fluxAggregator.On("RoundState", nodeAddr).Return(roundState3, nil).Once()
				// idleThreshold 2
				fluxAggregator.On("RoundState", nodeAddr).Return(roundState4, nil).Once().Run(func(args mock.Arguments) { idleThresholdOccured <- struct{}{} })
			}

			deviationChecker, err := fluxmonitor.NewPollingDeviationChecker(
				store,
				fluxAggregator,
				initr,
				runManager,
				fetcher,
				time.Duration(math.MaxInt64),
			)
			require.NoError(t, err)

			deviationChecker.OnConnect()
			deviationChecker.Start()
			require.Len(t, idleThresholdOccured, 0, "no Job Runs created")

			if test.expectedToSubmit {
				require.Eventually(t, func() bool { return len(idleThresholdOccured) == 1 }, 3*time.Second, 10*time.Millisecond)
				deviationChecker.HandleLog(&contracts.LogNewRound{RoundId: big.NewInt(int64(roundState1.ReportableRoundID))}, nil)
				require.Eventually(t, func() bool { return len(idleThresholdOccured) == 2 }, 3*time.Second, 10*time.Millisecond)
			}

			deviationChecker.Stop()

			if !test.expectedToSubmit {
				require.Len(t, idleThresholdOccured, 0)
			}

			fetcher.AssertExpectations(t)
			runManager.AssertExpectations(t)
			fluxAggregator.AssertExpectations(t)
		})
	}
}

func TestPollingDeviationChecker_RoundTimeoutCausesPoll(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	nodeAddr := ensureAccount(t, store)

	tests := []struct {
		name              string
		timesOutAt        func() int64
		expectedToTrigger bool
	}{
		{"timesOutAt == 0", func() int64 { return 0 }, false},
		{"timesOutAt != 0", func() int64 { return time.Now().Add(1 * time.Second).Unix() }, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fetcher := new(mocks.Fetcher)
			runManager := new(mocks.RunManager)
			fluxAggregator := new(mocks.FluxAggregator)

			job := cltest.NewJobWithFluxMonitorInitiator()
			initr := job.Initiators[0]
			initr.ID = 1
			initr.PollingInterval = models.Duration(math.MaxInt64)
			initr.IdleThreshold = models.Duration(0)

			const fetchedAnswer = 100
			answerBigInt := big.NewInt(fetchedAnswer * int64(math.Pow10(int(initr.InitiatorParams.Precision))))

			fluxAggregator.On("SubscribeToLogs", mock.Anything).Return(true, eth.UnsubscribeFunc(func() {}), nil)

			if test.expectedToTrigger {
				fluxAggregator.On("RoundState", nodeAddr).Return(contracts.FluxAggregatorRoundState{
					ReportableRoundID: 1,
					EligibleToSubmit:  false,
					LatestAnswer:      answerBigInt,
					TimesOutAt:        uint64(test.timesOutAt()),
				}, nil).Once()
				fluxAggregator.On("RoundState", nodeAddr).Return(contracts.FluxAggregatorRoundState{
					ReportableRoundID: 1,
					EligibleToSubmit:  false,
					LatestAnswer:      answerBigInt,
					TimesOutAt:        0,
				}, nil).Once()
			} else {
				fluxAggregator.On("RoundState", nodeAddr).Return(contracts.FluxAggregatorRoundState{
					ReportableRoundID: 1,
					EligibleToSubmit:  false,
					LatestAnswer:      answerBigInt,
					TimesOutAt:        uint64(test.timesOutAt()),
				}, nil).Once()
			}

			deviationChecker, err := fluxmonitor.NewPollingDeviationChecker(
				store,
				fluxAggregator,
				initr,
				runManager,
				fetcher,
				time.Duration(math.MaxInt64),
			)
			require.NoError(t, err)

			deviationChecker.Start()
			deviationChecker.OnConnect()
			time.Sleep(5 * time.Second)
			deviationChecker.Stop()

			fetcher.AssertExpectations(t)
			runManager.AssertExpectations(t)
			fluxAggregator.AssertExpectations(t)
		})
	}
}

func TestPollingDeviationChecker_RespondToNewRound(t *testing.T) {

	type roundIDCase struct {
		name                     string
		storedReportableRoundID  *big.Int
		fetchedReportableRoundID uint32
		logRoundID               int64
	}
	var (
		stored_lt_fetched_lt_log = roundIDCase{"stored < fetched < log", big.NewInt(5), 10, 15}
		stored_lt_log_lt_fetched = roundIDCase{"stored < log < fetched", big.NewInt(5), 15, 10}
		fetched_lt_stored_lt_log = roundIDCase{"fetched < stored < log", big.NewInt(10), 5, 15}
		fetched_lt_log_lt_stored = roundIDCase{"fetched < log < stored", big.NewInt(15), 5, 10}
		log_lt_fetched_lt_stored = roundIDCase{"log < fetched < stored", big.NewInt(15), 10, 5}
		log_lt_stored_lt_fetched = roundIDCase{"log < stored < fetched", big.NewInt(10), 15, 5}
		stored_lt_fetched_eq_log = roundIDCase{"stored < fetched = log", big.NewInt(5), 10, 10}
		stored_eq_fetched_lt_log = roundIDCase{"stored = fetched < log", big.NewInt(5), 5, 10}
		stored_eq_log_lt_fetched = roundIDCase{"stored = log < fetched", big.NewInt(5), 10, 5}
		fetched_lt_stored_eq_log = roundIDCase{"fetched < stored = log", big.NewInt(10), 5, 10}
		fetched_eq_log_lt_stored = roundIDCase{"fetched = log < stored", big.NewInt(10), 5, 5}
		log_lt_fetched_eq_stored = roundIDCase{"log < fetched = stored", big.NewInt(10), 10, 5}
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
		{true, true, true, stored_lt_fetched_lt_log, deviationThresholdExceeded},
		{true, true, true, stored_lt_log_lt_fetched, deviationThresholdExceeded},
		{true, true, true, fetched_lt_stored_lt_log, deviationThresholdExceeded},
		{true, true, true, fetched_lt_log_lt_stored, deviationThresholdExceeded},
		{true, true, true, log_lt_fetched_lt_stored, deviationThresholdExceeded},
		{true, true, true, log_lt_stored_lt_fetched, deviationThresholdExceeded},
		{true, true, true, stored_lt_fetched_eq_log, deviationThresholdExceeded},
		{true, true, true, stored_eq_fetched_lt_log, deviationThresholdExceeded},
		{true, true, true, stored_eq_log_lt_fetched, deviationThresholdExceeded},
		{true, true, true, fetched_lt_stored_eq_log, deviationThresholdExceeded},
		{true, true, true, fetched_eq_log_lt_stored, deviationThresholdExceeded},
		{true, true, true, log_lt_fetched_eq_stored, deviationThresholdExceeded},
		{true, true, true, stored_lt_fetched_lt_log, deviationThresholdNotExceeded},
		{true, true, true, stored_lt_log_lt_fetched, deviationThresholdNotExceeded},
		{true, true, true, fetched_lt_stored_lt_log, deviationThresholdNotExceeded},
		{true, true, true, fetched_lt_log_lt_stored, deviationThresholdNotExceeded},
		{true, true, true, log_lt_fetched_lt_stored, deviationThresholdNotExceeded},
		{true, true, true, log_lt_stored_lt_fetched, deviationThresholdNotExceeded},
		{true, true, true, stored_lt_fetched_eq_log, deviationThresholdNotExceeded},
		{true, true, true, stored_eq_fetched_lt_log, deviationThresholdNotExceeded},
		{true, true, true, stored_eq_log_lt_fetched, deviationThresholdNotExceeded},
		{true, true, true, fetched_lt_stored_eq_log, deviationThresholdNotExceeded},
		{true, true, true, fetched_eq_log_lt_stored, deviationThresholdNotExceeded},
		{true, true, true, log_lt_fetched_eq_stored, deviationThresholdNotExceeded},
		{true, true, false, stored_lt_fetched_lt_log, deviationThresholdExceeded},
		{true, true, false, stored_lt_log_lt_fetched, deviationThresholdExceeded},
		{true, true, false, fetched_lt_stored_lt_log, deviationThresholdExceeded},
		{true, true, false, fetched_lt_log_lt_stored, deviationThresholdExceeded},
		{true, true, false, log_lt_fetched_lt_stored, deviationThresholdExceeded},
		{true, true, false, log_lt_stored_lt_fetched, deviationThresholdExceeded},
		{true, true, false, stored_lt_fetched_eq_log, deviationThresholdExceeded},
		{true, true, false, stored_eq_fetched_lt_log, deviationThresholdExceeded},
		{true, true, false, stored_eq_log_lt_fetched, deviationThresholdExceeded},
		{true, true, false, fetched_lt_stored_eq_log, deviationThresholdExceeded},
		{true, true, false, fetched_eq_log_lt_stored, deviationThresholdExceeded},
		{true, true, false, log_lt_fetched_eq_stored, deviationThresholdExceeded},
		{true, true, false, stored_lt_fetched_lt_log, deviationThresholdNotExceeded},
		{true, true, false, stored_lt_log_lt_fetched, deviationThresholdNotExceeded},
		{true, true, false, fetched_lt_stored_lt_log, deviationThresholdNotExceeded},
		{true, true, false, fetched_lt_log_lt_stored, deviationThresholdNotExceeded},
		{true, true, false, log_lt_fetched_lt_stored, deviationThresholdNotExceeded},
		{true, true, false, log_lt_stored_lt_fetched, deviationThresholdNotExceeded},
		{true, true, false, stored_lt_fetched_eq_log, deviationThresholdNotExceeded},
		{true, true, false, stored_eq_fetched_lt_log, deviationThresholdNotExceeded},
		{true, true, false, stored_eq_log_lt_fetched, deviationThresholdNotExceeded},
		{true, true, false, fetched_lt_stored_eq_log, deviationThresholdNotExceeded},
		{true, true, false, fetched_eq_log_lt_stored, deviationThresholdNotExceeded},
		{true, true, false, log_lt_fetched_eq_stored, deviationThresholdNotExceeded},
		{true, false, true, stored_lt_fetched_lt_log, deviationThresholdExceeded},
		{true, false, true, stored_lt_log_lt_fetched, deviationThresholdExceeded},
		{true, false, true, fetched_lt_stored_lt_log, deviationThresholdExceeded},
		{true, false, true, fetched_lt_log_lt_stored, deviationThresholdExceeded},
		{true, false, true, log_lt_fetched_lt_stored, deviationThresholdExceeded},
		{true, false, true, log_lt_stored_lt_fetched, deviationThresholdExceeded},
		{true, false, true, stored_lt_fetched_eq_log, deviationThresholdExceeded},
		{true, false, true, stored_eq_fetched_lt_log, deviationThresholdExceeded},
		{true, false, true, stored_eq_log_lt_fetched, deviationThresholdExceeded},
		{true, false, true, fetched_lt_stored_eq_log, deviationThresholdExceeded},
		{true, false, true, fetched_eq_log_lt_stored, deviationThresholdExceeded},
		{true, false, true, log_lt_fetched_eq_stored, deviationThresholdExceeded},
		{true, false, true, stored_lt_fetched_lt_log, deviationThresholdNotExceeded},
		{true, false, true, stored_lt_log_lt_fetched, deviationThresholdNotExceeded},
		{true, false, true, fetched_lt_stored_lt_log, deviationThresholdNotExceeded},
		{true, false, true, fetched_lt_log_lt_stored, deviationThresholdNotExceeded},
		{true, false, true, log_lt_fetched_lt_stored, deviationThresholdNotExceeded},
		{true, false, true, log_lt_stored_lt_fetched, deviationThresholdNotExceeded},
		{true, false, true, stored_lt_fetched_eq_log, deviationThresholdNotExceeded},
		{true, false, true, stored_eq_fetched_lt_log, deviationThresholdNotExceeded},
		{true, false, true, stored_eq_log_lt_fetched, deviationThresholdNotExceeded},
		{true, false, true, fetched_lt_stored_eq_log, deviationThresholdNotExceeded},
		{true, false, true, fetched_eq_log_lt_stored, deviationThresholdNotExceeded},
		{true, false, true, log_lt_fetched_eq_stored, deviationThresholdNotExceeded},
		{true, false, false, stored_lt_fetched_lt_log, deviationThresholdExceeded},
		{true, false, false, stored_lt_log_lt_fetched, deviationThresholdExceeded},
		{true, false, false, fetched_lt_stored_lt_log, deviationThresholdExceeded},
		{true, false, false, fetched_lt_log_lt_stored, deviationThresholdExceeded},
		{true, false, false, log_lt_fetched_lt_stored, deviationThresholdExceeded},
		{true, false, false, log_lt_stored_lt_fetched, deviationThresholdExceeded},
		{true, false, false, stored_lt_fetched_eq_log, deviationThresholdExceeded},
		{true, false, false, stored_eq_fetched_lt_log, deviationThresholdExceeded},
		{true, false, false, stored_eq_log_lt_fetched, deviationThresholdExceeded},
		{true, false, false, fetched_lt_stored_eq_log, deviationThresholdExceeded},
		{true, false, false, fetched_eq_log_lt_stored, deviationThresholdExceeded},
		{true, false, false, log_lt_fetched_eq_stored, deviationThresholdExceeded},
		{true, false, false, stored_lt_fetched_lt_log, deviationThresholdNotExceeded},
		{true, false, false, stored_lt_log_lt_fetched, deviationThresholdNotExceeded},
		{true, false, false, fetched_lt_stored_lt_log, deviationThresholdNotExceeded},
		{true, false, false, fetched_lt_log_lt_stored, deviationThresholdNotExceeded},
		{true, false, false, log_lt_fetched_lt_stored, deviationThresholdNotExceeded},
		{true, false, false, log_lt_stored_lt_fetched, deviationThresholdNotExceeded},
		{true, false, false, stored_lt_fetched_eq_log, deviationThresholdNotExceeded},
		{true, false, false, stored_eq_fetched_lt_log, deviationThresholdNotExceeded},
		{true, false, false, stored_eq_log_lt_fetched, deviationThresholdNotExceeded},
		{true, false, false, fetched_lt_stored_eq_log, deviationThresholdNotExceeded},
		{true, false, false, fetched_eq_log_lt_stored, deviationThresholdNotExceeded},
		{true, false, false, log_lt_fetched_eq_stored, deviationThresholdNotExceeded},
		{false, true, true, stored_lt_fetched_lt_log, deviationThresholdExceeded},
		{false, true, true, stored_lt_log_lt_fetched, deviationThresholdExceeded},
		{false, true, true, fetched_lt_stored_lt_log, deviationThresholdExceeded},
		{false, true, true, fetched_lt_log_lt_stored, deviationThresholdExceeded},
		{false, true, true, log_lt_fetched_lt_stored, deviationThresholdExceeded},
		{false, true, true, log_lt_stored_lt_fetched, deviationThresholdExceeded},
		{false, true, true, stored_lt_fetched_eq_log, deviationThresholdExceeded},
		{false, true, true, stored_eq_fetched_lt_log, deviationThresholdExceeded},
		{false, true, true, stored_eq_log_lt_fetched, deviationThresholdExceeded},
		{false, true, true, fetched_lt_stored_eq_log, deviationThresholdExceeded},
		{false, true, true, fetched_eq_log_lt_stored, deviationThresholdExceeded},
		{false, true, true, log_lt_fetched_eq_stored, deviationThresholdExceeded},
		{false, true, true, stored_lt_fetched_lt_log, deviationThresholdNotExceeded},
		{false, true, true, stored_lt_log_lt_fetched, deviationThresholdNotExceeded},
		{false, true, true, fetched_lt_stored_lt_log, deviationThresholdNotExceeded},
		{false, true, true, fetched_lt_log_lt_stored, deviationThresholdNotExceeded},
		{false, true, true, log_lt_fetched_lt_stored, deviationThresholdNotExceeded},
		{false, true, true, log_lt_stored_lt_fetched, deviationThresholdNotExceeded},
		{false, true, true, stored_lt_fetched_eq_log, deviationThresholdNotExceeded},
		{false, true, true, stored_eq_fetched_lt_log, deviationThresholdNotExceeded},
		{false, true, true, stored_eq_log_lt_fetched, deviationThresholdNotExceeded},
		{false, true, true, fetched_lt_stored_eq_log, deviationThresholdNotExceeded},
		{false, true, true, fetched_eq_log_lt_stored, deviationThresholdNotExceeded},
		{false, true, true, log_lt_fetched_eq_stored, deviationThresholdNotExceeded},
		{false, true, false, stored_lt_fetched_lt_log, deviationThresholdExceeded},
		{false, true, false, stored_lt_log_lt_fetched, deviationThresholdExceeded},
		{false, true, false, fetched_lt_stored_lt_log, deviationThresholdExceeded},
		{false, true, false, fetched_lt_log_lt_stored, deviationThresholdExceeded},
		{false, true, false, log_lt_fetched_lt_stored, deviationThresholdExceeded},
		{false, true, false, log_lt_stored_lt_fetched, deviationThresholdExceeded},
		{false, true, false, stored_lt_fetched_eq_log, deviationThresholdExceeded},
		{false, true, false, stored_eq_fetched_lt_log, deviationThresholdExceeded},
		{false, true, false, stored_eq_log_lt_fetched, deviationThresholdExceeded},
		{false, true, false, fetched_lt_stored_eq_log, deviationThresholdExceeded},
		{false, true, false, fetched_eq_log_lt_stored, deviationThresholdExceeded},
		{false, true, false, log_lt_fetched_eq_stored, deviationThresholdExceeded},
		{false, true, false, stored_lt_fetched_lt_log, deviationThresholdNotExceeded},
		{false, true, false, stored_lt_log_lt_fetched, deviationThresholdNotExceeded},
		{false, true, false, fetched_lt_stored_lt_log, deviationThresholdNotExceeded},
		{false, true, false, fetched_lt_log_lt_stored, deviationThresholdNotExceeded},
		{false, true, false, log_lt_fetched_lt_stored, deviationThresholdNotExceeded},
		{false, true, false, log_lt_stored_lt_fetched, deviationThresholdNotExceeded},
		{false, true, false, stored_lt_fetched_eq_log, deviationThresholdNotExceeded},
		{false, true, false, stored_eq_fetched_lt_log, deviationThresholdNotExceeded},
		{false, true, false, stored_eq_log_lt_fetched, deviationThresholdNotExceeded},
		{false, true, false, fetched_lt_stored_eq_log, deviationThresholdNotExceeded},
		{false, true, false, fetched_eq_log_lt_stored, deviationThresholdNotExceeded},
		{false, true, false, log_lt_fetched_eq_stored, deviationThresholdNotExceeded},
		{false, false, true, stored_lt_fetched_lt_log, deviationThresholdExceeded},
		{false, false, true, stored_lt_log_lt_fetched, deviationThresholdExceeded},
		{false, false, true, fetched_lt_stored_lt_log, deviationThresholdExceeded},
		{false, false, true, fetched_lt_log_lt_stored, deviationThresholdExceeded},
		{false, false, true, log_lt_fetched_lt_stored, deviationThresholdExceeded},
		{false, false, true, log_lt_stored_lt_fetched, deviationThresholdExceeded},
		{false, false, true, stored_lt_fetched_eq_log, deviationThresholdExceeded},
		{false, false, true, stored_eq_fetched_lt_log, deviationThresholdExceeded},
		{false, false, true, stored_eq_log_lt_fetched, deviationThresholdExceeded},
		{false, false, true, fetched_lt_stored_eq_log, deviationThresholdExceeded},
		{false, false, true, fetched_eq_log_lt_stored, deviationThresholdExceeded},
		{false, false, true, log_lt_fetched_eq_stored, deviationThresholdExceeded},
		{false, false, true, stored_lt_fetched_lt_log, deviationThresholdNotExceeded},
		{false, false, true, stored_lt_log_lt_fetched, deviationThresholdNotExceeded},
		{false, false, true, fetched_lt_stored_lt_log, deviationThresholdNotExceeded},
		{false, false, true, fetched_lt_log_lt_stored, deviationThresholdNotExceeded},
		{false, false, true, log_lt_fetched_lt_stored, deviationThresholdNotExceeded},
		{false, false, true, log_lt_stored_lt_fetched, deviationThresholdNotExceeded},
		{false, false, true, stored_lt_fetched_eq_log, deviationThresholdNotExceeded},
		{false, false, true, stored_eq_fetched_lt_log, deviationThresholdNotExceeded},
		{false, false, true, stored_eq_log_lt_fetched, deviationThresholdNotExceeded},
		{false, false, true, fetched_lt_stored_eq_log, deviationThresholdNotExceeded},
		{false, false, true, fetched_eq_log_lt_stored, deviationThresholdNotExceeded},
		{false, false, true, log_lt_fetched_eq_stored, deviationThresholdNotExceeded},
		{false, false, false, stored_lt_fetched_lt_log, deviationThresholdExceeded},
		{false, false, false, stored_lt_log_lt_fetched, deviationThresholdExceeded},
		{false, false, false, fetched_lt_stored_lt_log, deviationThresholdExceeded},
		{false, false, false, fetched_lt_log_lt_stored, deviationThresholdExceeded},
		{false, false, false, log_lt_fetched_lt_stored, deviationThresholdExceeded},
		{false, false, false, log_lt_stored_lt_fetched, deviationThresholdExceeded},
		{false, false, false, stored_lt_fetched_eq_log, deviationThresholdExceeded},
		{false, false, false, stored_eq_fetched_lt_log, deviationThresholdExceeded},
		{false, false, false, stored_eq_log_lt_fetched, deviationThresholdExceeded},
		{false, false, false, fetched_lt_stored_eq_log, deviationThresholdExceeded},
		{false, false, false, fetched_eq_log_lt_stored, deviationThresholdExceeded},
		{false, false, false, log_lt_fetched_eq_stored, deviationThresholdExceeded},
		{false, false, false, stored_lt_fetched_lt_log, deviationThresholdNotExceeded},
		{false, false, false, stored_lt_log_lt_fetched, deviationThresholdNotExceeded},
		{false, false, false, fetched_lt_stored_lt_log, deviationThresholdNotExceeded},
		{false, false, false, fetched_lt_log_lt_stored, deviationThresholdNotExceeded},
		{false, false, false, log_lt_fetched_lt_stored, deviationThresholdNotExceeded},
		{false, false, false, log_lt_stored_lt_fetched, deviationThresholdNotExceeded},
		{false, false, false, stored_lt_fetched_eq_log, deviationThresholdNotExceeded},
		{false, false, false, stored_eq_fetched_lt_log, deviationThresholdNotExceeded},
		{false, false, false, stored_eq_log_lt_fetched, deviationThresholdNotExceeded},
		{false, false, false, fetched_lt_stored_eq_log, deviationThresholdNotExceeded},
		{false, false, false, fetched_eq_log_lt_stored, deviationThresholdNotExceeded},
		{false, false, false, log_lt_fetched_eq_stored, deviationThresholdNotExceeded},
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
			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			nodeAddr := ensureAccount(t, store)

			expectedToFetchRoundState := !test.startedBySelf
			expectedToPoll := expectedToFetchRoundState && test.eligible && test.funded && test.logRoundID >= int64(test.fetchedReportableRoundID)
			expectedToSubmit := expectedToPoll

			job := cltest.NewJobWithFluxMonitorInitiator()
			initr := job.Initiators[0]
			initr.ID = 1
			initr.InitiatorParams.PollingInterval = models.Duration(1 * time.Hour)

			rm := new(mocks.RunManager)
			fetcher := new(mocks.Fetcher)
			fluxAggregator := new(mocks.FluxAggregator)

			var availableFunds *big.Int
			var paymentAmount *big.Int
			minPayment := store.Config.MinimumContractPayment().ToInt()
			if test.funded {
				availableFunds = minPayment
				paymentAmount = minPayment
			} else {
				availableFunds = big.NewInt(1)
				paymentAmount = minPayment
			}

			if expectedToFetchRoundState {
				fluxAggregator.On("RoundState", nodeAddr).Return(contracts.FluxAggregatorRoundState{
					ReportableRoundID: test.fetchedReportableRoundID,
					LatestAnswer:      big.NewInt(test.latestAnswer * int64(math.Pow10(int(initr.InitiatorParams.Precision)))),
					EligibleToSubmit:  test.eligible,
					AvailableFunds:    availableFunds,
					PaymentAmount:     paymentAmount,
				}, nil).Once()
			}

			if expectedToPoll {
				fetcher.On("Fetch").Return(decimal.NewFromInt(test.polledAnswer), nil).Once()
			}

			if expectedToSubmit {
				fluxAggregator.On("GetMethodID", "updateAnswer").Return(updateAnswerSelector, nil)

				data, err := models.ParseJSON([]byte(fmt.Sprintf(`{
					"result": "%d",
					"address": "%s",
					"functionSelector": "0xe6330cf7",
					"dataPrefix": "0x%0x"
				}`, test.polledAnswer, initr.InitiatorParams.Address.Hex(), utils.EVMWordUint64(uint64(test.fetchedReportableRoundID)))))
				require.NoError(t, err)

				rm.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.MatchedBy(func(runRequest *models.RunRequest) bool {
					return reflect.DeepEqual(runRequest.RequestParams.Result.Value(), data.Result.Value())
				})).Return(nil, nil)
			}

			checker, err := fluxmonitor.NewPollingDeviationChecker(store, fluxAggregator, initr, rm, fetcher, time.Hour)
			require.NoError(t, err)

			checker.ExportedSetStoredReportableRoundID(test.storedReportableRoundID)

			checker.OnConnect()

			var startedBy common.Address
			if test.startedBySelf {
				startedBy = nodeAddr
			}
			checker.ExportedRespondToLog(&contracts.LogNewRound{RoundId: big.NewInt(test.logRoundID), StartedBy: startedBy})

			fluxAggregator.AssertExpectations(t)
			fetcher.AssertExpectations(t)
			rm.AssertExpectations(t)
		})
	}
}

func TestOutsideDeviation(t *testing.T) {
	tests := []struct {
		name                string
		curPrice, nextPrice decimal.Decimal
		threshold           float64 // in percentage
		expectation         bool
	}{
		{"0 current price, outside deviation", decimal.NewFromInt(0), decimal.NewFromInt(100), 2, true},
		{"0 current price, inside deviation", decimal.NewFromInt(0), decimal.NewFromInt(1), 2, true},
		{"0 current and next price", decimal.NewFromInt(0), decimal.NewFromInt(0), 2, false},

		{"inside deviation", decimal.NewFromInt(100), decimal.NewFromInt(101), 2, false},
		{"equal to deviation", decimal.NewFromInt(100), decimal.NewFromInt(102), 2, true},
		{"outside deviation", decimal.NewFromInt(100), decimal.NewFromInt(103), 2, true},
		{"outside deviation zero", decimal.NewFromInt(100), decimal.NewFromInt(0), 2, true},

		{"inside deviation, crosses 0 backwards", decimal.NewFromFloat(0.1), decimal.NewFromFloat(-0.1), 201, false},
		{"equal to deviation, crosses 0 backwards", decimal.NewFromFloat(0.1), decimal.NewFromFloat(-0.1), 200, true},
		{"outside deviation, crosses 0 backwards", decimal.NewFromFloat(0.1), decimal.NewFromFloat(-0.1), 199, true},

		{"inside deviation, crosses 0 forwards", decimal.NewFromFloat(-0.1), decimal.NewFromFloat(0.1), 201, false},
		{"equal to deviation, crosses 0 forwards", decimal.NewFromFloat(-0.1), decimal.NewFromFloat(0.1), 200, true},
		{"outside deviation, crosses 0 forwards", decimal.NewFromFloat(-0.1), decimal.NewFromFloat(0.1), 199, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := fluxmonitor.OutsideDeviation(test.curPrice, test.nextPrice, test.threshold)
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
			val, err := fluxmonitor.ExtractFeedURLs(initiatorParams.Feeds, store.ORM)
			require.NoError(t, err)
			assert.Equal(t, val, expectation)
		})
	}
}
