package fluxmonitor_test

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/eth"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	ethsvc "github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/eth/contracts"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitor"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const oracleCount uint32 = 17

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
			require.NoError(t, err)
			defer fm.Stop()
			assert.Equal(t, test.wantStarted, logBroadcaster.Started)
		})
	}
}

func TestConcreteFluxMonitor_AddJobRemoveJob(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	txm := new(mocks.TxManager)
	store.TxManager = txm
	txm.On("GetLatestBlock").Return(eth.Block{Number: hexutil.Uint64(123)}, nil)
	txm.On("GetLogs", mock.Anything).Return([]eth.Log{}, nil)

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
		name              string
		eligible          bool
		connected         bool
		funded            bool
		relativeThreshold float64
		absoluteThreshold float64
		latestAnswer      int64
		polledAnswer      int64
		expectedToPoll    bool
		expectedToSubmit  bool
		reportRegardless  bool
	}{
		{"eligible, connected, funded, answers deviate", true, true, true, 0.1, 1e-10, 1, 100, true, true, false},
		{"eligible, connected, funded, answers do not deviate", true, true, true, 0.1, 1e-10, 100, 100, true, false, false},
		{"eligible, connected, funded, report regardless of deviation, answers deviate", true, true, true, 0.1, 1e-10, 1, 100, true, true, true},
		{"eligible, connected, funded, report regardless of deviation, answers do not deviate", true, true, true, 0.1, 1e-10, 1, 100, true, true, true},
		{"eligible, connected, funded, absolute but not relative deviation", true, true, true, 5, 1, 100, 102, true, false, false},

		{"eligible, disconnected, funded, answers deviate", true, false, true, 0.1, 1e-10, 1, 100, false, false, false},
		{"eligible, disconnected, funded, answers do not deviate", true, false, true, 0.1, 1e-10, 100, 100, false, false, false},
		{"eligible, disconnected, funded, report regardless of deviation, answers deviate", true, false, true, 0.1, 1e-10, 1, 100, false, false, true},
		{"eligible, disconnected, funded, report regardless of deviation, answers do not deviate", true, false, true, 0.1, 1e-10, 1, 100, false, false, true},

		{"ineligible, connected, funded, answers deviate", false, true, true, 0.1, 1e-10, 1, 100, false, false, false},
		{"ineligible, connected, funded, answers do not deviate", false, true, true, 0.1, 1e-10, 100, 100, false, false, false},
		{"ineligible, connected, funded, report regardless of deviation, answers deviate", false, true, true, 0.1, 1e-10, 1, 100, false, false, true},
		{"ineligible, connected, funded, report regardless of deviation, answers do not deviate", false, true, true, 0.1, 1e-10, 1, 100, false, false, true},

		{"ineligible, disconnected, funded, answers deviate", false, false, true, 0.1, 1e-10, 1, 100, false, false, false},
		{"ineligible, disconnected, funded, answers do not deviate", false, false, true, 0.1, 1e-10, 100, 100, false, false, false},
		{"ineligible, disconnected, funded, report regardless of deviation, answers deviate", false, false, true, 0.1, 1e-10, 1, 100, false, false, true},
		{"ineligible, disconnected, funded, report regardless of deviation, answers do not deviate", false, false, true, 0.1, 1e-10, 1, 100, false, false, true},

		{"eligible, connected, underfunded, answers deviate", true, true, false, 0.1, 1e-10, 1, 100, false, false, false},
		{"eligible, connected, underfunded, answers do not deviate", true, true, false, 0.1, 1e-10, 100, 100, false, false, false},
		{"eligible, connected, underfunded, report regardless of deviation, answers deviate", true, true, false, 0.1, 1e-10, 1, 100, false, false, true},
		{"eligible, connected, underfunded, report regardless of deviation, answers do not deviate", true, true, false, 0.1, 1e-10, 1, 100, false, false, true},

		{"eligible, disconnected, underfunded, answers deviate", true, false, false, 0.1, 1e-10, 1, 100, false, false, false},
		{"eligible, disconnected, underfunded, answers do not deviate", true, false, false, 0.1, 1e-10, 100, 100, false, false, false},
		{"eligible, disconnected, underfunded, report regardless of deviation, answers deviate", true, false, false, 0.1, 1e-10, 1, 100, false, false, true},
		{"eligible, disconnected, underfunded, report regardless of deviation, answers do not deviate", true, false, false, 0.1, 1e-10, 1, 100, false, false, true},

		{"ineligible, connected, underfunded, answers deviate", false, true, false, 0.1, 1e-10, 1, 100, false, false, false},
		{"ineligible, connected, underfunded, answers do not deviate", false, true, false, 0.1, 1e-10, 100, 100, false, false, false},
		{"ineligible, connected, underfunded, report regardless of deviation, answers deviate", false, true, false, 0.1, 1e-10, 1, 100, false, false, true},
		{"ineligible, connected, underfunded, report regardless of deviation, answers do not deviate", false, true, false, 0.1, 1e-10, 1, 100, false, false, true},

		{"ineligible, disconnected, underfunded, answers deviate", false, false, false, 0.1, 1e-10, 1, 100, false, false, false},
		{"ineligible, disconnected, underfunded, answers do not deviate", false, false, false, 0.1, 1e-10, 100, 100, false, false, false},
		{"ineligible, disconnected, underfunded, report regardless of deviation, answers deviate", false, false, false, 0.1, 1e-10, 1, 100, false, false, true},
		{"ineligible, disconnected, underfunded, report regardless of deviation, answers do not deviate", false, false, false, 0.1, 1e-10, 1, 100, false, false, true},
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
			initr.ValueTriggers = initr.ValueTriggers[:0]
			triggerJSON := fmt.Sprintf(`{"absoluteThreshold": %e, "relativeThreshold": %e}`,
				test.absoluteThreshold, test.relativeThreshold)
			require.NoError(t, json.Unmarshal([]byte(triggerJSON), &initr.ValueTriggers),
				"failed to construct trigger functions from %s", triggerJSON)

			const reportableRoundID = 2
			latestAnswerNoPrecision := test.latestAnswer * int64(math.Pow10(int(initr.InitiatorParams.Precision)))

			paymentAmount := store.Config.MinimumContractPayment().ToInt()
			var availableFunds *big.Int
			if test.funded {
				availableFunds = big.NewInt(1).Mul(paymentAmount, big.NewInt(1000))
			} else {
				availableFunds = big.NewInt(1)
			}

			roundState := contracts.FluxAggregatorRoundState{
				ReportableRoundID: reportableRoundID,
				EligibleToSubmit:  test.eligible,
				LatestAnswer:      big.NewInt(latestAnswerNoPrecision),
				AvailableFunds:    availableFunds,
				PaymentAmount:     paymentAmount,
				OracleCount:       oracleCount,
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

			checker, err := fluxmonitor.NewPollingDeviationChecker(store,
				fluxAggregator, initr, rm, fetcher, models.MustMakeDuration(time.Second), func() {})
			require.NoError(t, err)

			if test.connected {
				checker.OnConnect()
			}

			// If you see panics with complaints about unregistered mocks here, it may
			// be because something in the table has expectedToSubmit false, but the
			// fluxmonitor has decided to report the polledAnswer anyway. (TODO:
			// Recover those panics, and give a more informative error message.)
			checker.ExportedPollIfEligible()

			fluxAggregator.AssertExpectations(t)
			fetcher.AssertExpectations(t)
			rm.AssertExpectations(t)
		})
	}
}

func TestPollingDeviationChecker_TriggerIdleTimeThreshold(t *testing.T) {

	tests := []struct {
		name             string
		idleThreshold    models.Duration
		expectedToSubmit bool
	}{
		{"no idleThreshold", models.MustMakeDuration(0), false},
		{"idleThreshold > 0", models.MustMakeDuration(10 * time.Millisecond), true},
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
			initr.PollingInterval = models.MustMakeDuration(math.MaxInt64)
			initr.IdleThreshold = test.idleThreshold

			const fetchedAnswer = 100
			answerBigInt := big.NewInt(fetchedAnswer * int64(math.Pow10(int(initr.InitiatorParams.Precision))))

			fluxAggregator.On("SubscribeToLogs", mock.Anything).Return(true, ethsvc.UnsubscribeFunc(func() {}), nil)

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
				models.MustMakeDuration(time.Duration(math.MaxInt64)),
				func() {},
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
			initr.PollingInterval = models.MustMakeDuration(math.MaxInt64)
			initr.IdleThreshold = models.MustMakeDuration(0)

			const fetchedAnswer = 100
			answerBigInt := big.NewInt(fetchedAnswer * int64(math.Pow10(int(initr.InitiatorParams.Precision))))

			fluxAggregator.On("SubscribeToLogs", mock.Anything).Return(true, ethsvc.UnsubscribeFunc(func() {}), nil)

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
				models.MustMakeDuration(time.Duration(math.MaxInt64)),
				func() {},
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
			initr.InitiatorParams.PollingInterval = models.MustMakeDuration(1 * time.Hour)

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
				fluxAggregator.On("RoundState", nodeAddr).Return(contracts.FluxAggregatorRoundState{
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

			checker, err := fluxmonitor.NewPollingDeviationChecker(store,
				fluxAggregator, initr, rm, fetcher, models.MustMakeDuration(time.Hour), func() {})
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
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	checker := cltest.NewPollingDeviationChecker(t, store)

	min := store.Config.MinimumContractPayment().ToInt().Int64()

	tests := []struct {
		name    string
		payment int64
		want    bool
	}{
		{"above minimum", min + 1, true},
		{"equal to minimum", min, true},
		{"below minimum", min - 1, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, checker.SufficientPayment(big.NewInt(test.payment)))
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
				OracleCount:    uint32(oracleCount),
			}
			assert.Equal(t, test.want, checker.SufficientFunds(state))
		})
	}
}
