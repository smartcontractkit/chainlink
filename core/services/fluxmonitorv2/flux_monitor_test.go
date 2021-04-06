package fluxmonitorv2_test

import (
	"context"
	"math"
	"math/big"
	"testing"
	"time"

	"gopkg.in/guregu/null.v4"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/logger"
	corenull "github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitorv2"
	fmmocks "github.com/smartcontractkit/chainlink/core/services/fluxmonitorv2/mocks"
	jobmocks "github.com/smartcontractkit/chainlink/core/services/job/mocks"
	logmocks "github.com/smartcontractkit/chainlink/core/services/log/mocks"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	pipelinemocks "github.com/smartcontractkit/chainlink/core/services/pipeline/mocks"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const oracleCount uint8 = 17

type answerSet struct{ latestAnswer, polledAnswer int64 }

var (
	now     = func() uint64 { return uint64(time.Now().UTC().Unix()) }
	nilOpts *bind.CallOpts

	makeRoundDataForRoundID = func(roundID uint32) flux_aggregator_wrapper.LatestRoundData {
		return flux_aggregator_wrapper.LatestRoundData{
			RoundId: big.NewInt(int64(roundID)),
		}
	}
	freshContractRoundDataResponse = func() (flux_aggregator_wrapper.LatestRoundData, error) {
		return flux_aggregator_wrapper.LatestRoundData{}, errors.New("No data present")
	}

	contractAddress   = cltest.NewAddress()
	threshold         = float64(0.5)
	absoluteThreshold = float64(0.01)
	idleTimerPeriod   = time.Minute
	precision         = int32(2)
	defaultLogger     = *logger.Default
	pipelineSpec      = pipeline.Spec{
		ID: 1,
		DotDagSource: `
// data source 1
ds1 [type=http method=GET url="https://pricesource1.com" requestData="{\\"coin\\": \\"ETH\\", \\"market\\": \\"USD\\"}"];
ds1_parse [type=jsonparse path="latest"];

// data source 2
ds2 [type=http method=GET url="https://pricesource1.com" requestData="{\\"coin\\": \\"ETH\\", \\"market\\": \\"USD\\"}"];
ds2_parse [type=jsonparse path="latest"];

ds1 -> ds1_parse -> answer1;
ds2 -> ds2_parse -> answer1;

answer1 [type=median index=0];
`,
		JobID: 1,
	}
)

// testMocks defines all the mock interfaces used by the Flux Monitor
type testMocks struct {
	fluxAggregator    *mocks.FluxAggregator
	logBroadcast      *logmocks.Broadcast
	logBroadcaster    *logmocks.Broadcaster
	orm               *fmmocks.ORM
	jobORM            *jobmocks.ORM
	pipelineORM       *pipelinemocks.ORM
	pipelineRunner    *pipelinemocks.Runner
	keyStore          *fmmocks.KeyStoreInterface
	contractSubmitter *fmmocks.ContractSubmitter
}

func newTestMocks() *testMocks {
	return &testMocks{
		fluxAggregator:    new(mocks.FluxAggregator),
		logBroadcast:      new(logmocks.Broadcast),
		logBroadcaster:    new(logmocks.Broadcaster),
		orm:               new(fmmocks.ORM),
		jobORM:            new(jobmocks.ORM),
		pipelineORM:       new(pipelinemocks.ORM),
		pipelineRunner:    new(pipelinemocks.Runner),
		keyStore:          new(fmmocks.KeyStoreInterface),
		contractSubmitter: new(fmmocks.ContractSubmitter),
	}
}

// AssertExpectations asserts expectations of all the mocks
func (tm *testMocks) AssertExpectations(t *testing.T) {
	tm.fluxAggregator.AssertExpectations(t)
	tm.logBroadcast.AssertExpectations(t)
	tm.logBroadcaster.AssertExpectations(t)
	tm.orm.AssertExpectations(t)
	tm.jobORM.AssertExpectations(t)
	tm.pipelineORM.AssertExpectations(t)
	tm.pipelineRunner.AssertExpectations(t)
	tm.keyStore.AssertExpectations(t)
	tm.contractSubmitter.AssertExpectations(t)
}

func setupMocks(t *testing.T) *testMocks {
	t.Helper()

	tm := newTestMocks()

	t.Cleanup(func() {
		tm.AssertExpectations(t)
	})

	return tm
}

type setupOptions struct {
	pollTickerDisabled bool
	idleTimerDisabled  bool
	idleTimerPeriod    time.Duration
	orm                fluxmonitorv2.ORM
}

// setup sets up a Flux Monitor for testing, allowing the test to provide
// functional options to configure the setup
func setup(t *testing.T, optionFns ...func(*setupOptions)) (*fluxmonitorv2.FluxMonitor, *testMocks) {
	t.Helper()

	tm := setupMocks(t)

	pipelineRun := fluxmonitorv2.NewPipelineRun(
		tm.pipelineRunner,
		pipelineSpec,
		defaultLogger,
	)

	options := setupOptions{
		idleTimerPeriod: time.Minute,
		orm:             tm.orm,
	}

	for _, optionFn := range optionFns {
		optionFn(&options)
	}

	fm, err := fluxmonitorv2.NewFluxMonitor(
		pipelineRun,
		options.orm,
		tm.jobORM,
		tm.pipelineORM,
		tm.keyStore,
		fluxmonitorv2.NewPollManager(
			fluxmonitorv2.PollManagerConfig{
				PollTickerInterval:      time.Minute,
				PollTickerDisabled:      options.pollTickerDisabled,
				IdleTimerPeriod:         options.idleTimerPeriod,
				IdleTimerDisabled:       options.idleTimerDisabled,
				HibernationPollPeriod:   24 * time.Hour,
				MinRetryBackoffDuration: 1 * time.Minute,
				MaxRetryBackoffDuration: 1 * time.Hour,
			},
			logger.Default,
		),
		fluxmonitorv2.NewPaymentChecker(assets.NewLink(1), nil),
		contractAddress,
		tm.contractSubmitter,
		fluxmonitorv2.NewDeviationChecker(threshold, absoluteThreshold),
		fluxmonitorv2.NewSubmissionChecker(big.NewInt(0), big.NewInt(100000000000), precision),
		fluxmonitorv2.Flags{},
		tm.fluxAggregator,
		tm.logBroadcaster,
		precision,
		func() {},
		logger.Default,
	)
	require.NoError(t, err)

	return fm, tm
}

// disablePollTicker is an option to disable the poll ticker during setup
func disablePollTicker(disabled bool) func(*setupOptions) {
	return func(opts *setupOptions) {
		opts.pollTickerDisabled = disabled
	}
}

// disableIdleTimer is an option to disable the idle timer during setup
func disableIdleTimer(disabled bool) func(*setupOptions) {
	return func(opts *setupOptions) {
		opts.idleTimerDisabled = disabled
	}
}

// setIdleTimerPeriod is an option to set the idle timer period during setup
func setIdleTimerPeriod(period time.Duration) func(*setupOptions) {
	return func(opts *setupOptions) {
		opts.idleTimerPeriod = period
	}
}

// withORM is an option to switch out the ORM during set up. Useful when you
// want to use a database backed ORM
func withORM(orm fluxmonitorv2.ORM) func(*setupOptions) {
	return func(opts *setupOptions) {
		opts.orm = orm
	}
}

// setupStoreWithKey setups a new store and adds a key to the keystore
func setupStoreWithKey(t *testing.T) (*store.Store, common.Address) {
	store, cleanup := cltest.NewStore(t)
	_, nodeAddr := cltest.MustAddRandomKeyToKeystore(t, store)

	t.Cleanup(cleanup)

	return store, nodeAddr
}

func TestFluxMonitor_PollIfEligible(t *testing.T) {
	testCases := []struct {
		name              string
		eligible          bool
		connected         bool
		funded            bool
		answersDeviate    bool
		hasPreviousRun    bool
		previousRunStatus pipeline.RunStatus
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
			hasPreviousRun: true, previousRunStatus: pipeline.RunStatusCompleted,
			expectedToPoll: false, expectedToSubmit: false,
		}, {
			name:     "previous job run in progress",
			eligible: true, connected: true, funded: true, answersDeviate: true,
			hasPreviousRun: true, previousRunStatus: pipeline.RunStatusInProgress,
			expectedToPoll: false, expectedToSubmit: false,
		}, {
			name:     "previous job run errored",
			eligible: true, connected: true, funded: true, answersDeviate: true,
			hasPreviousRun: true, previousRunStatus: pipeline.RunStatusErrored,
			expectedToPoll: true, expectedToSubmit: true,
		},
	}

	store, nodeAddr := setupStoreWithKey(t)

	const reportableRoundID = 2
	var (
		thresholds        = struct{ abs, rel float64 }{0.1, 200}
		deviatedAnswers   = answerSet{1, 100}
		undeviatedAnswers = answerSet{100, 101}
	)

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fm, tm := setup(t)

			tm.keyStore.On("Accounts").Return([]accounts.Account{{Address: nodeAddr}}).Once()

			// Setup Answers
			answers := undeviatedAnswers
			if tc.answersDeviate {
				answers = deviatedAnswers
			}
			latestAnswerNoPrecision := answers.latestAnswer * int64(
				math.Pow10(int(precision)),
			)

			// Setup Run
			run := pipeline.Run{
				ID:             1,
				PipelineSpecID: 1,
			}
			if tc.hasPreviousRun {
				switch tc.previousRunStatus {
				case pipeline.RunStatusCompleted:
					now := time.Now()
					run.FinishedAt = &now
				case pipeline.RunStatusErrored:
					run.Errors = []null.String{
						null.StringFrom("Random: String, foo"),
					}
				default:
				}

				tm.orm.
					On("FindOrCreateFluxMonitorRoundStats", contractAddress, uint32(reportableRoundID)).
					Return(fluxmonitorv2.FluxMonitorRoundStatsV2{
						Aggregator:     contractAddress,
						RoundID:        reportableRoundID,
						PipelineRunID:  corenull.Int64From(run.ID),
						NumSubmissions: 1,
					}, nil)

				tm.pipelineORM.
					On("FindRun", run.ID).
					Return(run, nil)
			} else {
				if tc.connected {
					tm.orm.
						On("FindOrCreateFluxMonitorRoundStats", contractAddress, uint32(reportableRoundID)).
						Return(fluxmonitorv2.FluxMonitorRoundStatsV2{
							Aggregator: contractAddress,
							RoundID:    reportableRoundID,
						}, nil)
				}
			}

			// Set up funds
			var availableFunds *big.Int
			var paymentAmount *big.Int
			minPayment := store.Config.MinimumContractPayment().ToInt()
			if tc.funded {
				availableFunds = big.NewInt(1).Mul(big.NewInt(10000), minPayment)
				paymentAmount = minPayment
			} else {
				availableFunds = big.NewInt(1)
				paymentAmount = minPayment
			}

			roundState := flux_aggregator_wrapper.OracleRoundState{
				RoundId:          reportableRoundID,
				EligibleToSubmit: tc.eligible,
				LatestSubmission: big.NewInt(latestAnswerNoPrecision),
				AvailableFunds:   availableFunds,
				PaymentAmount:    paymentAmount,
				OracleCount:      oracleCount,
			}
			tm.fluxAggregator.
				On("OracleRoundState", nilOpts, nodeAddr, uint32(0)).
				Return(roundState, nil).Maybe()

			if tc.expectedToPoll {
				tm.fluxAggregator.On("LatestRoundData", nilOpts).Return(flux_aggregator_wrapper.LatestRoundData{
					Answer:    big.NewInt(10),
					UpdatedAt: big.NewInt(100),
				}, nil)
				tm.pipelineRunner.
					On("ExecuteAndInsertNewRun", context.Background(), pipelineSpec, pipeline.JSONSerializable{
						Val: map[string]interface{}{
							"latestAnswer": float64(10),
							"updatedAt":    float64(100),
						},
					}, defaultLogger).
					Return(int64(1), pipeline.FinalResult{
						Values: []interface{}{decimal.NewFromInt(answers.polledAnswer)},
						Errors: []error{nil},
					}, nil)
			}

			if tc.expectedToSubmit {
				tm.contractSubmitter.
					On("Submit", big.NewInt(reportableRoundID), big.NewInt(answers.polledAnswer)).
					Return(nil).
					Once()

				tm.orm.
					On("UpdateFluxMonitorRoundStats",
						contractAddress,
						uint32(reportableRoundID),
						int64(1),
					).
					Return(nil)
			}

			if tc.connected {
				fm.OnConnect()
			}

			oracles := []common.Address{nodeAddr, cltest.NewAddress()}
			tm.fluxAggregator.On("GetOracles", nilOpts).Return(oracles, nil)
			fm.SetOracleAddress()

			fm.ExportedPollIfEligible(thresholds.rel, thresholds.abs)
		})
	}
}

// If the roundState method is unable to communicate with the contract (possibly due to
// incorrect address) then the pollIfEligible method should create a JobErr record
func TestFluxMonitor_PollIfEligible_Creates_JobErr(t *testing.T) {
	_, nodeAddr := setupStoreWithKey(t)
	oracles := []common.Address{nodeAddr, cltest.NewAddress()}

	var (
		roundState = flux_aggregator_wrapper.OracleRoundState{}
	)

	fm, tm := setup(t)

	tm.keyStore.On("Accounts").Return([]accounts.Account{{Address: nodeAddr}}).Once()

	tm.jobORM.
		On("RecordError",
			context.Background(),
			pipelineSpec.JobID,
			"Unable to call roundState method on provided contract. Check contract address.",
		).Once()

	tm.fluxAggregator.
		On("OracleRoundState", nilOpts, nodeAddr, mock.Anything).
		Return(roundState, errors.New("err")).
		Once()

	fm.OnConnect()

	tm.fluxAggregator.On("GetOracles", nilOpts).Return(oracles, nil)
	require.NoError(t, fm.SetOracleAddress())

	fm.ExportedPollIfEligible(1, 1)
}

func TestPollingDeviationChecker_BuffersLogs(t *testing.T) {
	store, nodeAddr := setupStoreWithKey(t)
	oracles := []common.Address{nodeAddr, cltest.NewAddress()}

	fm, tm := setup(t,
		disableIdleTimer(true),
		disablePollTicker(true),
	)

	const (
		fetchedValue = 100
	)

	// Test helpers
	var (
		makeRoundStateForRoundID = func(roundID uint32) flux_aggregator_wrapper.OracleRoundState {
			return flux_aggregator_wrapper.OracleRoundState{
				RoundId:          roundID,
				EligibleToSubmit: true,
				LatestSubmission: big.NewInt(100 * int64(math.Pow10(int(precision)))),
				AvailableFunds:   store.Config.MinimumContractPayment().ToInt(),
				PaymentAmount:    store.Config.MinimumContractPayment().ToInt(),
			}
		}
	)

	chBlock := make(chan struct{})
	chSafeToAssert := make(chan struct{})
	chSafeToFillQueue := make(chan struct{})

	tm.keyStore.On("Accounts").Return([]accounts.Account{{Address: nodeAddr}}).Once()

	tm.fluxAggregator.On("LatestRoundData", nilOpts).Return(freshContractRoundDataResponse()).Maybe()
	tm.fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(1)).
		Return(makeRoundStateForRoundID(1), nil).
		Run(func(mock.Arguments) {
			close(chSafeToFillQueue)
			<-chBlock
		}).
		Once()
	tm.fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(3)).Return(makeRoundStateForRoundID(3), nil).Once()
	tm.fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(4)).Return(makeRoundStateForRoundID(4), nil).Once()
	tm.fluxAggregator.On("GetOracles", nilOpts).Return(oracles, nil)
	// tm.fluxAggregator.On("Address").Return(contractAddress, nil)

	tm.logBroadcaster.On("Register", fm, mock.Anything).Return(true, func() {})

	tm.orm.On("MostRecentFluxMonitorRoundID", contractAddress).Return(uint32(1), nil)
	tm.orm.On("MostRecentFluxMonitorRoundID", contractAddress).Return(uint32(3), nil)
	tm.orm.On("MostRecentFluxMonitorRoundID", contractAddress).Return(uint32(4), nil)

	// Round 1
	tm.orm.
		On("FindOrCreateFluxMonitorRoundStats", contractAddress, uint32(1)).
		Return(fluxmonitorv2.FluxMonitorRoundStatsV2{
			Aggregator: contractAddress,
			RoundID:    1,
		}, nil)
	tm.pipelineRunner.
		On("ExecuteAndInsertNewRun", context.Background(), pipelineSpec, pipeline.JSONSerializable{Val: map[string]interface{}(nil), Null: false}, defaultLogger).
		Return(int64(1), pipeline.FinalResult{
			Values: []interface{}{decimal.NewFromInt(fetchedValue)},
			Errors: []error{nil},
		}, nil)
	tm.contractSubmitter.
		On("Submit", big.NewInt(1), big.NewInt(fetchedValue)).
		Return(nil).
		Once()

	tm.orm.
		On("UpdateFluxMonitorRoundStats",
			contractAddress,
			uint32(1),
			mock.AnythingOfType("int64"), //int64(1),
		).
		Return(nil).Once()

	// Round 3
	tm.orm.
		On("FindOrCreateFluxMonitorRoundStats", contractAddress, uint32(3)).
		Return(fluxmonitorv2.FluxMonitorRoundStatsV2{
			Aggregator: contractAddress,
			RoundID:    3,
		}, nil)
	tm.pipelineRunner.
		On("ExecuteAndInsertNewRun", context.Background(), pipelineSpec, pipeline.JSONSerializable{Val: map[string]interface{}(nil), Null: false}, defaultLogger).
		Return(int64(2), pipeline.FinalResult{
			Values: []interface{}{decimal.NewFromInt(fetchedValue)},
			Errors: []error{nil},
		}, nil)
	tm.contractSubmitter.
		On("Submit", big.NewInt(3), big.NewInt(fetchedValue)).
		Return(nil).
		Once()
	tm.orm.
		On("UpdateFluxMonitorRoundStats",
			contractAddress,
			uint32(3),
			mock.AnythingOfType("int64"), //int64(2),

		).
		Return(nil).Once()

	// Round 4
	tm.orm.
		On("FindOrCreateFluxMonitorRoundStats", contractAddress, uint32(4)).
		Return(fluxmonitorv2.FluxMonitorRoundStatsV2{
			Aggregator: contractAddress,
			RoundID:    3,
		}, nil)
	tm.pipelineRunner.
		On("ExecuteAndInsertNewRun", context.Background(), pipelineSpec, pipeline.JSONSerializable{Val: map[string]interface{}(nil), Null: false}, defaultLogger).
		Return(int64(3), pipeline.FinalResult{
			Values: []interface{}{decimal.NewFromInt(fetchedValue)},
			Errors: []error{nil},
		}, nil)
	tm.contractSubmitter.
		On("Submit", big.NewInt(4), big.NewInt(fetchedValue)).
		Return(nil).
		Once()
	tm.orm.
		On("UpdateFluxMonitorRoundStats",
			contractAddress,
			uint32(4),
			mock.AnythingOfType("int64"), //int64(3),
		).
		Return(nil).
		Once().
		Run(func(mock.Arguments) { close(chSafeToAssert) })

	fm.OnConnect()
	fm.Start()

	var logBroadcasts []*logmocks.Broadcast

	for i := 1; i <= 4; i++ {
		logBroadcast := new(logmocks.Broadcast)
		logBroadcast.On("DecodedLog").Return(&flux_aggregator_wrapper.FluxAggregatorNewRound{RoundId: big.NewInt(int64(i)), StartedAt: big.NewInt(0)})
		logBroadcast.On("WasAlreadyConsumed").Return(false, nil)
		logBroadcast.On("MarkConsumed").Return(nil)
		logBroadcasts = append(logBroadcasts, logBroadcast)
	}

	fm.HandleLog(logBroadcasts[0]) // Get the checker to start processing a log so we can freeze it
	<-chSafeToFillQueue
	fm.HandleLog(logBroadcasts[1]) // This log is evicted from the priority queue
	fm.HandleLog(logBroadcasts[2])
	fm.HandleLog(logBroadcasts[3])

	close(chBlock)
	<-chSafeToAssert
	fm.Close()
}

func TestFluxMonitor_TriggerIdleTimeThreshold(t *testing.T) {
	testCases := []struct {
		name              string
		idleTimerDisabled bool
		idleDuration      time.Duration
		expectedToSubmit  bool
	}{
		{"no idleDuration", true, 0, false},
		{"idleDuration > 0", false, 2 * time.Second, true},
	}

	store, nodeAddr := setupStoreWithKey(t)
	oracles := []common.Address{nodeAddr, cltest.NewAddress()}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var (
				orm = fluxmonitorv2.NewORM(store.DB)
			)

			fm, tm := setup(t,
				disablePollTicker(true),
				disableIdleTimer(tc.idleTimerDisabled),
				setIdleTimerPeriod(tc.idleDuration),
				withORM(orm),
			)

			tm.keyStore.On("Accounts").Return([]accounts.Account{{Address: nodeAddr}}).Once()

			const fetchedAnswer = 100
			answerBigInt := big.NewInt(fetchedAnswer * int64(math.Pow10(int(precision))))

			tm.fluxAggregator.On("GetOracles", nilOpts).Return(oracles, nil)
			tm.logBroadcaster.On("Register", mock.Anything, mock.Anything).Return(true, func() {})

			idleDurationOccured := make(chan struct{}, 3)

			tm.fluxAggregator.On("LatestRoundData", nilOpts).Return(freshContractRoundDataResponse()).Once()
			if tc.expectedToSubmit {
				// performInitialPoll()
				roundState1 := flux_aggregator_wrapper.OracleRoundState{RoundId: 1, EligibleToSubmit: false, LatestSubmission: answerBigInt, StartedAt: now()}
				tm.fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(0)).Return(roundState1, nil).Once()
				// idleDuration 1
				roundState2 := flux_aggregator_wrapper.OracleRoundState{RoundId: 1, EligibleToSubmit: false, LatestSubmission: answerBigInt, StartedAt: now()}
				tm.fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(0)).Return(roundState2, nil).Once().Run(func(args mock.Arguments) {
					idleDurationOccured <- struct{}{}
				})
			}

			fm.OnConnect()
			fm.Start()
			require.Len(t, idleDurationOccured, 0, "no Job Runs created")

			if tc.expectedToSubmit {
				require.Eventually(t, func() bool { return len(idleDurationOccured) == 1 }, 3*time.Second, 10*time.Millisecond)

				chBlock := make(chan struct{})
				// NewRound resets the idle timer
				roundState2 := flux_aggregator_wrapper.OracleRoundState{RoundId: 2, EligibleToSubmit: false, LatestSubmission: answerBigInt, StartedAt: now()}
				tm.fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(2)).Return(roundState2, nil).Once().Run(func(args mock.Arguments) {
					close(chBlock)
				})

				decodedLog := flux_aggregator_wrapper.FluxAggregatorNewRound{RoundId: big.NewInt(2), StartedAt: big.NewInt(0)}
				tm.logBroadcast.On("DecodedLog").Return(&decodedLog)
				tm.logBroadcast.On("WasAlreadyConsumed").Return(false, nil).Once()
				tm.logBroadcast.On("MarkConsumed").Return(nil).Once()
				fm.HandleLog(tm.logBroadcast)

				gomega.NewGomegaWithT(t).Eventually(chBlock).Should(gomega.BeClosed())

				// idleDuration 2
				roundState3 := flux_aggregator_wrapper.OracleRoundState{RoundId: 3, EligibleToSubmit: false, LatestSubmission: answerBigInt, StartedAt: now()}
				tm.fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(0)).Return(roundState3, nil).Once().Run(func(args mock.Arguments) {
					idleDurationOccured <- struct{}{}
				})
				require.Eventually(t, func() bool { return len(idleDurationOccured) == 2 }, 3*time.Second, 10*time.Millisecond)
			}

			fm.Close()

			if !tc.expectedToSubmit {
				require.Len(t, idleDurationOccured, 0)
			}
		})
	}
}

func TestFluxMonitor_IdleTimerResetsOnNewRound(t *testing.T) {
	t.Parallel()

	_, nodeAddr := setupStoreWithKey(t)
	oracles := []common.Address{nodeAddr, cltest.NewAddress()}

	fm, tm := setup(t,
		disablePollTicker(true),
		setIdleTimerPeriod(2*time.Second),
	)

	tm.keyStore.On("Accounts").Return([]accounts.Account{{Address: nodeAddr}}).Once()

	const fetchedAnswer = 100
	answerBigInt := big.NewInt(fetchedAnswer * int64(math.Pow10(int(precision))))

	tm.fluxAggregator.On("GetOracles", nilOpts).Return(oracles, nil)
	tm.logBroadcaster.On("Register", mock.Anything, mock.Anything).Return(true, func() {})

	tm.fluxAggregator.On("LatestRoundData", nilOpts).Return(freshContractRoundDataResponse()).Once()

	idleDurationOccured := make(chan struct{}, 3)

	// Initial Poll
	roundState1 := flux_aggregator_wrapper.OracleRoundState{RoundId: 1, EligibleToSubmit: false, LatestSubmission: answerBigInt, StartedAt: now()}
	tm.fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(0)).Return(roundState1, nil).Once()
	tm.orm.
		On("FindOrCreateFluxMonitorRoundStats", contractAddress, uint32(1)).
		Return(fluxmonitorv2.FluxMonitorRoundStatsV2{
			Aggregator:     contractAddress,
			RoundID:        1,
			NumSubmissions: 0,
		}, nil).Once()

	fm.OnConnect()
	fm.Start()
	t.Cleanup(func() { fm.Close() })
	require.Len(t, idleDurationOccured, 0, "no Job Runs created")

	// idleDuration 1 triggers using the same round id as the initial poll. This resets the idle timer
	roundState1Responded := flux_aggregator_wrapper.OracleRoundState{RoundId: 1, EligibleToSubmit: false, LatestSubmission: answerBigInt, StartedAt: now() + 1}
	tm.fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(0)).Return(roundState1Responded, nil).Once().Run(func(args mock.Arguments) {
		idleDurationOccured <- struct{}{}
	})
	// Finds an existing run created by the initial poll
	tm.orm.
		On("FindOrCreateFluxMonitorRoundStats", contractAddress, uint32(1)).
		Return(fluxmonitorv2.FluxMonitorRoundStatsV2{
			PipelineRunID:  corenull.NewInt64(int64(1), true),
			Aggregator:     contractAddress,
			RoundID:        1,
			NumSubmissions: 1,
		}, nil).Once()
	finishedAt := time.Now()
	tm.pipelineORM.On("FindRun", int64(1)).Return(pipeline.Run{
		FinishedAt: &finishedAt,
	}, nil)

	require.Eventually(t, func() bool { return len(idleDurationOccured) == 1 }, 3*time.Second, 10*time.Millisecond)

	// idleDuration 2 triggers a new round. Started at is 0
	roundState2 := flux_aggregator_wrapper.OracleRoundState{RoundId: 2, EligibleToSubmit: false, LatestSubmission: answerBigInt, StartedAt: 0}
	tm.fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(0)).Return(roundState2, nil).Once().Run(func(args mock.Arguments) {
		idleDurationOccured <- struct{}{}
	})
	tm.orm.
		On("FindOrCreateFluxMonitorRoundStats", contractAddress, uint32(2)).
		Return(fluxmonitorv2.FluxMonitorRoundStatsV2{
			Aggregator:     contractAddress,
			RoundID:        2,
			NumSubmissions: 0,
		}, nil).Once()

	require.Eventually(t, func() bool { return len(idleDurationOccured) == 2 }, 3*time.Second, 10*time.Millisecond)

	// idleDuration 3 triggers from the previous new round
	roundState3 := flux_aggregator_wrapper.OracleRoundState{RoundId: 3, EligibleToSubmit: false, LatestSubmission: answerBigInt, StartedAt: now()}
	tm.fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(0)).Return(roundState3, nil).Once().Run(func(args mock.Arguments) {
		idleDurationOccured <- struct{}{}
	})
	tm.orm.
		On("FindOrCreateFluxMonitorRoundStats", contractAddress, uint32(3)).
		Return(fluxmonitorv2.FluxMonitorRoundStatsV2{
			Aggregator:     contractAddress,
			RoundID:        3,
			NumSubmissions: 0,
		}, nil).Once()

	require.Eventually(t, func() bool { return len(idleDurationOccured) == 3 }, 3*time.Second, 10*time.Millisecond)
}

func TestFluxMonitor_RoundTimeoutCausesPoll_timesOutAtZero(t *testing.T) {
	t.Parallel()

	store, nodeAddr := setupStoreWithKey(t)

	var (
		oracles = []common.Address{nodeAddr, cltest.NewAddress()}
		orm     = fluxmonitorv2.NewORM(store.DB)
	)

	fm, tm := setup(t,
		disablePollTicker(true),
		disableIdleTimer(true),
		withORM(orm),
	)

	tm.keyStore.
		On("Accounts").
		Return([]accounts.Account{{Address: nodeAddr}}).
		Twice() // Once called from the test, once during start

	ch := make(chan struct{})

	const fetchedAnswer = 100
	answerBigInt := big.NewInt(fetchedAnswer * int64(math.Pow10(int(precision))))
	tm.logBroadcaster.On("Register", mock.Anything, mock.Anything).Return(true, func() {})

	tm.fluxAggregator.On("LatestRoundData", nilOpts).Return(makeRoundDataForRoundID(1), nil).Once()
	roundState0 := flux_aggregator_wrapper.OracleRoundState{RoundId: 1, EligibleToSubmit: false, LatestSubmission: answerBigInt, StartedAt: now()}
	tm.fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(1)).Return(roundState0, nil).Once() // initialRoundState()
	tm.fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(0)).Return(flux_aggregator_wrapper.OracleRoundState{
		RoundId:          1,
		EligibleToSubmit: false,
		LatestSubmission: answerBigInt,
		StartedAt:        0,
		Timeout:          0,
	}, nil).
		Run(func(mock.Arguments) { close(ch) }).
		Once()

	tm.fluxAggregator.On("GetOracles", nilOpts).Return(oracles, nil)

	fm.SetOracleAddress()
	fm.ExportedRoundState()
	fm.Start()
	fm.OnConnect()

	gomega.NewGomegaWithT(t).Eventually(ch).Should(gomega.BeClosed())

	fm.Close()
}

func TestFluxMonitor_UsesPreviousRoundStateOnStartup_RoundTimeout(t *testing.T) {
	store, nodeAddr := setupStoreWithKey(t)
	oracles := []common.Address{nodeAddr, cltest.NewAddress()}

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
			t.Parallel()

			var (
				orm = fluxmonitorv2.NewORM(store.DB)
			)

			fm, tm := setup(t,
				disablePollTicker(true),
				disableIdleTimer(true),
				withORM(orm),
			)

			tm.keyStore.On("Accounts").Return([]accounts.Account{{Address: nodeAddr}}).Once()

			tm.logBroadcaster.On("Register", mock.Anything, mock.Anything).Return(true, func() {})

			tm.fluxAggregator.On("GetOracles", nilOpts).Return(oracles, nil)

			tm.fluxAggregator.On("LatestRoundData", nilOpts).Return(makeRoundDataForRoundID(1), nil).Once()
			tm.fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(1)).Return(flux_aggregator_wrapper.OracleRoundState{
				RoundId:          1,
				EligibleToSubmit: false,
				StartedAt:        now(),
				Timeout:          test.timeout,
			}, nil).Once()

			// 2nd roundstate call means round timer triggered
			chRoundState := make(chan struct{})
			tm.fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(0)).Return(flux_aggregator_wrapper.OracleRoundState{
				RoundId:          1,
				EligibleToSubmit: false,
			}, nil).
				Run(func(mock.Arguments) { close(chRoundState) }).
				Maybe()

			fm.Start()
			fm.OnConnect()

			if test.expectedToSubmit {
				gomega.NewGomegaWithT(t).Eventually(chRoundState).Should(gomega.BeClosed())
			} else {
				gomega.NewGomegaWithT(t).Consistently(chRoundState).ShouldNot(gomega.BeClosed())
			}

			fm.Close()
		})
	}
}

func TestFluxMonitor_UsesPreviousRoundStateOnStartup_IdleTimer(t *testing.T) {
	store, nodeAddr := setupStoreWithKey(t)
	oracles := []common.Address{nodeAddr, cltest.NewAddress()}

	almostExpired := time.Now().
		Add(idleTimerPeriod * -1).
		Add(2 * time.Second).
		Unix()

	testCases := []struct {
		name             string
		startedAt        uint64
		expectedToSubmit bool
	}{
		{"active round exists - idleTimer about to expired", uint64(almostExpired), true},
		{"active round exists - idleTimer will not expire", 100, false},
		{"no active round", 0, false},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var (
				orm = fluxmonitorv2.NewORM(store.DB)
			)

			fm, tm := setup(t,
				disablePollTicker(true),
				withORM(orm),
			)

			tm.keyStore.On("Accounts").Return([]accounts.Account{{Address: nodeAddr}}).Once()

			tm.logBroadcaster.On("Register", mock.Anything, mock.Anything).Return(true, func() {})

			tm.fluxAggregator.On("GetOracles", nilOpts).Return(oracles, nil)
			tm.fluxAggregator.On("LatestRoundData", nilOpts).Return(makeRoundDataForRoundID(1), nil).Once()
			// first roundstate in setInitialTickers()
			tm.fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(1)).Return(flux_aggregator_wrapper.OracleRoundState{
				RoundId:          1,
				EligibleToSubmit: false,
				StartedAt:        tc.startedAt,
				Timeout:          10000, // round won't time out
			}, nil).Once()

			// 2nd roundstate in performInitialPoll()
			roundState := flux_aggregator_wrapper.OracleRoundState{RoundId: 1, EligibleToSubmit: false}
			tm.fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(0)).Return(roundState, nil).Once()

			// 3rd roundState call means idleTimer triggered
			chRoundState := make(chan struct{})
			tm.fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(0)).Return(roundState, nil).
				Run(func(mock.Arguments) { close(chRoundState) }).
				Maybe()

			fm.Start()
			fm.OnConnect()

			if tc.expectedToSubmit {
				gomega.NewGomegaWithT(t).Eventually(chRoundState).Should(gomega.BeClosed())
			} else {
				gomega.NewGomegaWithT(t).Consistently(chRoundState).ShouldNot(gomega.BeClosed())
			}

			fm.Close()
		})
	}
}

func TestFluxMonitor_RoundTimeoutCausesPoll_timesOutNotZero(t *testing.T) {
	t.Parallel()

	store, nodeAddr := setupStoreWithKey(t)
	oracles := []common.Address{nodeAddr, cltest.NewAddress()}

	var (
		orm = fluxmonitorv2.NewORM(store.DB)
	)

	fm, tm := setup(t,
		disablePollTicker(true),
		disableIdleTimer(true),
		withORM(orm),
	)

	tm.keyStore.On("Accounts").Return([]accounts.Account{{Address: nodeAddr}}).Once()

	const fetchedAnswer = 100
	answerBigInt := big.NewInt(fetchedAnswer * int64(math.Pow10(int(precision))))

	chRoundState1 := make(chan struct{})
	chRoundState2 := make(chan struct{})

	tm.logBroadcaster.On("Register", mock.Anything, mock.Anything).Return(true, func() {})

	tm.fluxAggregator.On("GetOracles", nilOpts).Return(oracles, nil)
	tm.fluxAggregator.On("LatestRoundData", nilOpts).Return(makeRoundDataForRoundID(1), nil).Once()
	tm.fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(1)).Return(flux_aggregator_wrapper.OracleRoundState{
		RoundId:          1,
		EligibleToSubmit: false,
		LatestSubmission: answerBigInt,
		StartedAt:        now(),
		Timeout:          uint64(1000000),
	}, nil).Once()

	startedAt := uint64(time.Now().Unix())
	timeout := uint64(3)
	tm.fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(0)).Return(flux_aggregator_wrapper.OracleRoundState{
		RoundId:          1,
		EligibleToSubmit: false,
		LatestSubmission: answerBigInt,
		StartedAt:        startedAt,
		Timeout:          timeout,
	}, nil).Once().
		Run(func(mock.Arguments) { close(chRoundState1) }).
		Once()
	tm.fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(0)).Return(flux_aggregator_wrapper.OracleRoundState{
		RoundId:          1,
		EligibleToSubmit: false,
		LatestSubmission: answerBigInt,
		StartedAt:        startedAt,
		Timeout:          timeout,
	}, nil).Once().
		Run(func(mock.Arguments) { close(chRoundState2) }).
		Once()

	fm.Start()
	fm.OnConnect()

	tm.logBroadcast.On("WasAlreadyConsumed").Return(false, nil)
	tm.logBroadcast.On("DecodedLog").Return(&flux_aggregator_wrapper.FluxAggregatorNewRound{RoundId: big.NewInt(0), StartedAt: big.NewInt(time.Now().UTC().Unix())})
	tm.logBroadcast.On("MarkConsumed").Return(nil)
	fm.HandleLog(tm.logBroadcast)

	gomega.NewGomegaWithT(t).Eventually(chRoundState1).Should(gomega.BeClosed())
	gomega.NewGomegaWithT(t).Eventually(chRoundState2).Should(gomega.BeClosed())

	time.Sleep(time.Duration(2*timeout) * time.Second)
	fm.Close()
}

func TestFluxMonitor_HandlesNilLogs(t *testing.T) {
	t.Parallel()

	fm, _ := setup(t)

	logBroadcast := new(logmocks.Broadcast)
	var logNewRound *flux_aggregator_wrapper.FluxAggregatorNewRound
	var logAnswerUpdated *flux_aggregator_wrapper.FluxAggregatorAnswerUpdated
	var randomType interface{}

	logBroadcast.On("DecodedLog").Return(logNewRound).Once()
	assert.NotPanics(t, func() {
		fm.HandleLog(logBroadcast)
	})

	logBroadcast.On("DecodedLog").Return(logAnswerUpdated).Once()
	assert.NotPanics(t, func() {
		fm.HandleLog(logBroadcast)
	})

	logBroadcast.On("DecodedLog").Return(randomType).Once()
	assert.NotPanics(t, func() {
		fm.HandleLog(logBroadcast)
	})

	logBroadcast.AssertExpectations(t)
}

func TestFluxMonitor_ConsumeLogBroadcast(t *testing.T) {
	t.Parallel()

	fm, tm := setup(t)

	tm.fluxAggregator.
		On("OracleRoundState", nilOpts, mock.Anything, mock.Anything).
		Return(flux_aggregator_wrapper.OracleRoundState{RoundId: 123}, nil)

	tm.logBroadcast.On("WasAlreadyConsumed").Return(false, nil).Once()
	tm.logBroadcast.On("DecodedLog").Return(&flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{})
	tm.logBroadcast.On("MarkConsumed").Return(nil).Once()

	fm.ExportedBacklog().Add(fluxmonitorv2.PriorityNewRoundLog, tm.logBroadcast)
	fm.ExportedProcessLogs()
}

func TestFluxMonitor_ConsumeLogBroadcast_Error(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		consumed bool
		err      error
	}{
		{"already consumed", true, nil},
		{"error determining already consumed", false, errors.New("err")},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fm, tm := setup(t)

			tm.logBroadcast.On("WasAlreadyConsumed").Return(tc.consumed, tc.err).Once()

			fm.ExportedBacklog().Add(fluxmonitorv2.PriorityNewRoundLog, tm.logBroadcast)
			fm.ExportedProcessLogs()
		})
	}
}

func TestFluxMonitor_DoesNotDoubleSubmit(t *testing.T) {
	t.Run("when NewRound log arrives, then poll ticker fires", func(t *testing.T) {
		store, nodeAddr := setupStoreWithKey(t)
		oracles := []common.Address{nodeAddr, cltest.NewAddress()}

		fm, tm := setup(t,
			disableIdleTimer(true),
			disablePollTicker(true),
		)

		var (
			paymentAmount  = store.Config.MinimumContractPayment().ToInt()
			availableFunds = big.NewInt(1).Mul(paymentAmount, big.NewInt(1000))
		)

		const (
			roundID = 3
			answer  = 100
		)

		tm.keyStore.On("Accounts").Return([]accounts.Account{{Address: nodeAddr}}).Once()

		// Mocks initiated by the New Round log
		tm.orm.On("MostRecentFluxMonitorRoundID", contractAddress).Return(uint32(roundID), nil).Once()
		tm.orm.
			On("FindOrCreateFluxMonitorRoundStats", contractAddress, uint32(roundID)).
			Return(fluxmonitorv2.FluxMonitorRoundStatsV2{
				Aggregator: contractAddress,
				RoundID:    roundID,
			}, nil).Once()
		tm.pipelineRunner.
			On("ExecuteAndInsertNewRun", context.Background(), pipelineSpec, mock.Anything, defaultLogger).
			Return(int64(1), pipeline.FinalResult{
				Values: []interface{}{decimal.NewFromInt(answer)},
				Errors: []error{nil},
			}, nil).Once()
		tm.contractSubmitter.On("Submit", big.NewInt(roundID), big.NewInt(answer)).Return(nil).Once()
		tm.orm.
			On("UpdateFluxMonitorRoundStats",
				contractAddress,
				uint32(roundID),
				int64(1),
			).
			Return(nil)

		tm.fluxAggregator.On("GetOracles", nilOpts).Return(oracles, nil)
		fm.SetOracleAddress()
		fm.OnConnect()

		tm.fluxAggregator.On("LatestRoundData", nilOpts).Return(flux_aggregator_wrapper.LatestRoundData{
			Answer:    big.NewInt(10),
			UpdatedAt: big.NewInt(100),
		}, nil)
		// Fire off the NewRound log, which the node should respond to
		tm.fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(roundID)).
			Return(flux_aggregator_wrapper.OracleRoundState{
				RoundId:          roundID,
				LatestSubmission: big.NewInt(answer),
				EligibleToSubmit: true,
				AvailableFunds:   availableFunds,
				PaymentAmount:    paymentAmount,
				OracleCount:      1,
			}, nil).
			Once()

		fm.ExportedRespondToNewRoundLog(&flux_aggregator_wrapper.FluxAggregatorNewRound{
			RoundId:   big.NewInt(roundID),
			StartedAt: big.NewInt(0),
		})

		// Mocks initiated by polling
		// Now force the node to try to poll and ensure it does not respond this time
		tm.fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(0)).
			Return(flux_aggregator_wrapper.OracleRoundState{
				RoundId:          roundID,
				LatestSubmission: big.NewInt(answer),
				EligibleToSubmit: true,
				AvailableFunds:   availableFunds,
				PaymentAmount:    paymentAmount,
				OracleCount:      1,
			}, nil).
			Once()
		tm.orm.
			On("FindOrCreateFluxMonitorRoundStats", contractAddress, uint32(roundID)).
			Return(fluxmonitorv2.FluxMonitorRoundStatsV2{
				PipelineRunID:  corenull.NewInt64(int64(1), true),
				Aggregator:     contractAddress,
				RoundID:        roundID,
				NumSubmissions: 1,
			}, nil).Once()
		now := time.Now()
		tm.pipelineORM.On("FindRun", int64(1)).Return(pipeline.Run{
			FinishedAt: &now,
		}, nil)

		fm.ExportedPollIfEligible(0, 0)
	})

	t.Run("when poll ticker fires, then NewRound log arrives", func(t *testing.T) {
		store, nodeAddr := setupStoreWithKey(t)
		oracles := []common.Address{nodeAddr, cltest.NewAddress()}
		fm, tm := setup(t,
			disableIdleTimer(true),
			disablePollTicker(true),
		)

		var (
			paymentAmount  = store.Config.MinimumContractPayment().ToInt()
			availableFunds = big.NewInt(1).Mul(paymentAmount, big.NewInt(1000))
		)

		const (
			roundID = 3
			answer  = 100
		)
		tm.keyStore.On("Accounts").Return([]accounts.Account{{Address: nodeAddr}}).Once()

		fm.OnConnect()

		// First, force the node to try to poll, which should result in a submission
		tm.fluxAggregator.On("LatestRoundData", nilOpts).Return(flux_aggregator_wrapper.LatestRoundData{
			Answer:    big.NewInt(10),
			UpdatedAt: big.NewInt(100),
		}, nil)
		tm.fluxAggregator.On("OracleRoundState", nilOpts, nodeAddr, uint32(0)).
			Return(flux_aggregator_wrapper.OracleRoundState{
				RoundId:          roundID,
				LatestSubmission: big.NewInt(answer),
				EligibleToSubmit: true,
				AvailableFunds:   availableFunds,
				PaymentAmount:    paymentAmount,
				OracleCount:      1,
			}, nil).
			Once()
		tm.orm.
			On("FindOrCreateFluxMonitorRoundStats", contractAddress, uint32(roundID)).
			Return(fluxmonitorv2.FluxMonitorRoundStatsV2{
				Aggregator: contractAddress,
				RoundID:    roundID,
			}, nil).Once()
		tm.pipelineRunner.
			On("ExecuteAndInsertNewRun", context.Background(), pipelineSpec, mock.Anything, defaultLogger).
			Return(int64(1), pipeline.FinalResult{
				Values: []interface{}{decimal.NewFromInt(answer)},
				Errors: []error{nil},
			}, nil).Once()
		tm.contractSubmitter.On("Submit", big.NewInt(roundID), big.NewInt(answer)).Return(nil).Once()
		tm.orm.
			On("UpdateFluxMonitorRoundStats",
				contractAddress,
				uint32(roundID),
				int64(1),
			).
			Return(nil).
			Once()

		tm.fluxAggregator.On("GetOracles", nilOpts).Return(oracles, nil)
		fm.SetOracleAddress()
		fm.ExportedPollIfEligible(0, 0)

		// Now fire off the NewRound log and ensure it does not respond this time
		tm.orm.On("MostRecentFluxMonitorRoundID", contractAddress).Return(uint32(roundID), nil)
		tm.orm.
			On("FindOrCreateFluxMonitorRoundStats", contractAddress, uint32(roundID)).
			Return(fluxmonitorv2.FluxMonitorRoundStatsV2{
				PipelineRunID:  corenull.NewInt64(int64(1), true),
				Aggregator:     contractAddress,
				RoundID:        roundID,
				NumSubmissions: 1,
			}, nil).Once()
		tm.pipelineORM.On("FindRun", int64(1)).Return(pipeline.Run{}, nil)

		fm.ExportedRespondToNewRoundLog(&flux_aggregator_wrapper.FluxAggregatorNewRound{
			RoundId:   big.NewInt(roundID),
			StartedAt: big.NewInt(0),
		})
	})
}
