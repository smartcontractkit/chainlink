package fluxmonitorv2

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flags_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitor/promfm"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/tevino/abool"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const hibernationPollPeriod = 24 * time.Hour

// FluxMonitor polls external price adapters via HTTP to check for price swings.
type FluxMonitor struct {
	jobID             int32
	contractAddress   common.Address
	oracleAddress     common.Address
	pipelineRun       PipelineRun
	orm               ORM
	jobORM            job.ORM
	pipelineORM       pipeline.ORM
	keyStore          KeyStoreInterface
	paymentChecker    *PaymentChecker
	contractSubmitter ContractSubmitter
	deviationChecker  *DeviationChecker
	submissionChecker *SubmissionChecker
	flags             Flags
	fluxAggregator    flux_aggregator_wrapper.FluxAggregatorInterface
	logBroadcaster    log.Broadcaster

	logger    *logger.Logger
	precision int32

	isHibernating bool
	connected     *abool.AtomicBool
	backlog       *utils.BoundedPriorityQueue
	chProcessLogs chan struct{}

	pollTicker       *PollTicker
	hibernationTimer utils.ResettableTimer
	idleTimer        *IdleTimer
	roundTimer       utils.ResettableTimer

	readyForLogs func()
	chStop       chan struct{}
	waitOnStop   chan struct{}
}

// NewFluxMonitor returns a new instance of PollingDeviationChecker.
func NewFluxMonitor(
	jobID int32,
	pipelineRun PipelineRun,
	orm ORM,
	jobORM job.ORM,
	pipelineORM pipeline.ORM,
	keyStore KeyStoreInterface,
	pollTicker *PollTicker,
	idleTimer *IdleTimer,
	paymentChecker *PaymentChecker,
	contractAddress common.Address,
	contractSubmitter ContractSubmitter,
	deviationChecker *DeviationChecker,
	submissionChecker *SubmissionChecker,
	flags Flags,
	fluxAggregator flux_aggregator_wrapper.FluxAggregatorInterface,
	logBroadcaster log.Broadcaster,
	precision int32,
	readyForLogs func(),
	fmLogger *logger.Logger,
) (*FluxMonitor, error) {
	fm := &FluxMonitor{
		jobID:             jobID,
		pipelineRun:       pipelineRun,
		orm:               orm,
		jobORM:            jobORM,
		pipelineORM:       pipelineORM,
		keyStore:          keyStore,
		pollTicker:        pollTicker,
		idleTimer:         idleTimer,
		paymentChecker:    paymentChecker,
		contractAddress:   contractAddress,
		contractSubmitter: contractSubmitter,
		deviationChecker:  deviationChecker,
		submissionChecker: submissionChecker,
		flags:             flags,

		readyForLogs:   readyForLogs,
		logBroadcaster: logBroadcaster,
		fluxAggregator: fluxAggregator,
		precision:      precision,
		logger:         fmLogger,

		hibernationTimer: utils.NewResettableTimer(),
		roundTimer:       utils.NewResettableTimer(),
		isHibernating:    false,
		connected:        abool.New(),
		backlog: utils.NewBoundedPriorityQueue(map[uint]uint{
			// We want reconnecting nodes to be able to submit to a round
			// that hasn't hit maxAnswers yet, as well as the newest round.
			PriorityNewRoundLog:      2,
			PriorityAnswerUpdatedLog: 1,
			PriorityFlagChangedLog:   2,
		}),
		chProcessLogs: make(chan struct{}, 1),
		chStop:        make(chan struct{}),
		waitOnStop:    make(chan struct{}),
	}

	return fm, nil
}

// NewFromJobSpec constructs an instance of FluxMonitor with sane defaults and
// validation.
func NewFromJobSpec(
	jobSpec job.Job,
	orm ORM,
	jobORM job.ORM,
	pipelineORM pipeline.ORM,
	keyStore KeyStoreInterface,
	ethClient eth.Client,
	logBroadcaster log.Broadcaster,
	pipelineRunner pipeline.Runner,
	cfg Config,
) (*FluxMonitor, error) {
	fmSpec := jobSpec.FluxMonitorSpec

	if !validatePollTimer(fmSpec.PollTimerDisabled, cfg.MinimumPollingInterval(), fmSpec.PollTimerPeriod) {
		return nil, fmt.Errorf(
			"pollTimerPeriod (%s), must be equal or greater than %s",
			fmSpec.PollTimerPeriod,
			cfg.MinimumPollingInterval(),
		)
	}

	// Set up the flux aggregator
	logBroadcaster.AddDependents(1)
	fluxAggregator, err := flux_aggregator_wrapper.NewFluxAggregator(
		fmSpec.ContractAddress.Address(),
		ethClient,
	)
	if err != nil {
		return nil, err
	}

	contractSubmitter := NewFluxAggregatorContractSubmitter(
		fluxAggregator,
		orm,
		keyStore,
		cfg.EthGasLimit,
		cfg.MaxUnconfirmedTransactions,
	)

	flags, err := NewFlags(cfg.FlagsContractAddress, ethClient)
	logger.ErrorIf(
		err,
		fmt.Sprintf(
			"unable to create Flags contract instance, check address: %s",
			cfg.FlagsContractAddress,
		),
	)

	paymentChecker := &PaymentChecker{
		MinContractPayment: cfg.MinContractPayment,
		MinJobPayment:      fmSpec.MinPayment,
	}

	pipelineRun := PipelineRun{
		runner: pipelineRunner,
		spec:   *jobSpec.PipelineSpec,
		jobID:  jobSpec.ID,
		logger: *logger.Default,
	}

	min, err := fluxAggregator.MinSubmissionValue(nil)
	if err != nil {
		return nil, err
	}

	max, err := fluxAggregator.MaxSubmissionValue(nil)
	if err != nil {
		return nil, err
	}

	fmLogger := logger.CreateLogger(
		logger.Default.With(
			"jobID", jobSpec.ID,
			"contract", fmSpec.ContractAddress.Hex(),
		),
	)

	return NewFluxMonitor(
		jobSpec.ID,
		pipelineRun,
		orm,
		jobORM,
		pipelineORM,
		keyStore,
		NewPollTicker(fmSpec.PollTimerPeriod, fmSpec.PollTimerDisabled),
		NewIdleTimer(fmSpec.IdleTimerPeriod, fmSpec.IdleTimerDisabled),
		paymentChecker,
		fmSpec.ContractAddress.Address(),
		contractSubmitter,
		NewDeviationChecker(
			float64(fmSpec.Threshold),
			float64(fmSpec.AbsoluteThreshold),
		),
		NewSubmissionChecker(min, max, fmSpec.Precision),
		*flags,
		fluxAggregator,
		logBroadcaster,
		fmSpec.Precision,
		func() { logBroadcaster.DependentReady() },
		fmLogger,
	)
}

const (
	PriorityFlagChangedLog   uint = 0
	PriorityNewRoundLog      uint = 1
	PriorityAnswerUpdatedLog uint = 2
)

// Start implements the job.Service interface. It begins the CSP consumer in a
// single goroutine to poll the price adapters and listen to NewRound events.
func (fm *FluxMonitor) Start() error {
	fm.logger.Debug("Starting Flux Monitor for job")

	go fm.consume()

	return nil
}

func (fm *FluxMonitor) setIsHibernatingStatus() {
	if !fm.flags.ContractExists() {
		fm.isHibernating = false

		return
	}

	isFlagLowered, err := fm.flags.IsLowered(fm.contractAddress)
	if err != nil {
		fm.logger.Errorf("unable to set hibernation status: %v", err)

		fm.isHibernating = false
	} else {
		fm.isHibernating = !isFlagLowered
	}
}

// Close implements the job.Service interface. It stops this instance from
// polling, cleaning up resources.
func (fm *FluxMonitor) Close() error {
	fm.pollTicker.Stop()
	fm.hibernationTimer.Stop()
	fm.idleTimer.Stop()
	fm.roundTimer.Stop()
	close(fm.chStop)
	<-fm.waitOnStop

	return nil
}

// OnConnect sets the poller as connected
func (fm *FluxMonitor) OnConnect() {
	fm.logger.Debugw("Flux Monitor connected to Ethereum node")

	fm.connected.Set()
}

// OnDisconnect sets the poller as disconnected
func (fm *FluxMonitor) OnDisconnect() {
	fm.logger.Debugw("Flux Monitor disconnected from Ethereum node")

	fm.connected.UnSet()
}

// JobID implements the listener.Listener interface.
//
// Since we don't have a v1 ID, we return a new v1 job id to satisfy the
// interface. This should not cause a problem as the log broadcaster will check
// if it is a v2 job before attempting to use this job id
func (fm *FluxMonitor) JobID() models.JobID {
	return models.NewJobID()
}

// JobIDV2 implements the listener.Listener interface.
//
// Returns the v2 job id
func (fm *FluxMonitor) JobIDV2() int32 { return fm.jobID }

// IsV2Job implements the listener.Listener interface.
//
// Returns true as this is a v2 job
func (fm *FluxMonitor) IsV2Job() bool { return true }

// HandleLog processes the contract logs
func (fm *FluxMonitor) HandleLog(broadcast log.Broadcast, err error) {
	if err != nil {
		fm.logger.Errorf("got error from LogBroadcaster: %v", err)
		return
	}

	log := broadcast.DecodedLog()
	if log == nil || reflect.ValueOf(log).IsNil() {
		fm.logger.Error("HandleLog: ignoring nil value")
		return
	}

	switch log := log.(type) {
	case *flux_aggregator_wrapper.FluxAggregatorNewRound:
		fm.backlog.Add(PriorityNewRoundLog, broadcast)

	case *flux_aggregator_wrapper.FluxAggregatorAnswerUpdated:
		fm.backlog.Add(PriorityAnswerUpdatedLog, broadcast)

	case *flags_wrapper.FlagsFlagRaised:
		if log.Subject == utils.ZeroAddress || log.Subject == fm.contractAddress {
			fm.backlog.Add(PriorityFlagChangedLog, broadcast)
		}

	case *flags_wrapper.FlagsFlagLowered:
		if log.Subject == utils.ZeroAddress || log.Subject == fm.contractAddress {
			fm.backlog.Add(PriorityFlagChangedLog, broadcast)
		}

	default:
		fm.logger.Warnf("unexpected log type %T", log)
		return
	}

	select {
	case fm.chProcessLogs <- struct{}{}:
	default:
	}
}

func (fm *FluxMonitor) consume() {
	defer close(fm.waitOnStop)

	if err := fm.SetOracleAddress(); err != nil {
		fm.logger.Warnw(
			"unable to set oracle address, this flux monitor job may not work correctly",
			"err",
			err,
		)
	}

	// Subscribe to contract logs
	isConnected := fm.logBroadcaster.Register(fm.fluxAggregator, fm)
	defer fm.logBroadcaster.Unregister(fm.fluxAggregator, fm)

	if fm.flags.ContractExists() {
		flagsConnected := fm.logBroadcaster.Register(fm.flags.Contract(), fm)
		isConnected = isConnected && flagsConnected
		defer fm.logBroadcaster.Unregister(fm.flags.Contract(), fm)
	}

	if isConnected {
		fm.connected.Set()
	} else {
		fm.connected.UnSet()
	}

	fm.readyForLogs()
	fm.setIsHibernatingStatus()
	fm.setInitialTickers()
	fm.performInitialPoll()

	for {
		select {
		case <-fm.chStop:
			return

		case <-fm.chProcessLogs:
			fm.processLogs()

		case <-fm.pollTicker.Ticks():
			fm.logger.Debugw("Poll ticker fired", fm.loggerFieldsForTick()...)
			fm.pollIfEligible(fm.deviationChecker)

		case <-fm.idleTimer.Ticks():
			fm.logger.Debugw("Idle ticker fired", fm.loggerFieldsForTick()...)
			fm.pollIfEligible(NewZeroDeviationChecker())

		case <-fm.roundTimer.Ticks():
			fm.logger.Debugw("Round timeout ticker fired", fm.loggerFieldsForTick()...)
			fm.pollIfEligible(fm.deviationChecker)

		case <-fm.hibernationTimer.Ticks():
			fm.logger.Debugw("Hibernation timout ticker fired", fm.loggerFieldsForTick()...)
			fm.pollIfEligible(NewZeroDeviationChecker())
		}
	}
}

// SetOracleAddress sets the oracle address which matches the node's keys.
// If none match, it uses the first available key
func (fm *FluxMonitor) SetOracleAddress() error {
	oracleAddrs, err := fm.fluxAggregator.GetOracles(nil)
	if err != nil {
		return errors.Wrap(err, "failed to get list of oracles from FluxAggregator contract")
	}
	accounts := fm.keyStore.Accounts()
	for _, acct := range accounts {
		for _, oracleAddr := range oracleAddrs {
			if acct.Address == oracleAddr {
				fm.oracleAddress = oracleAddr
				return nil
			}
		}
	}
	if len(accounts) > 0 {
		addr := accounts[0].Address
		fm.logger.Warnw("None of the node's keys matched any oracle addresses, using first available key. This flux monitor job may not work correctly", "address", addr)
		fm.oracleAddress = addr
	} else {
		fm.logger.Error("No keys found. This flux monitor job may not work correctly")
	}
	return errors.New("none of the node's keys matched any oracle addresses")

}

// performInitialPoll performs the initial poll if required
func (fm *FluxMonitor) performInitialPoll() {
	if fm.shouldPerformInitialPoll() {
		fm.pollIfEligible(fm.deviationChecker)
	}
}

func (fm *FluxMonitor) shouldPerformInitialPoll() bool {
	return !(fm.pollTicker.IsDisabled() && fm.idleTimer.IsDisabled() || fm.isHibernating)
}

// hibernate restarts the PollingDeviationChecker in hibernation mode
func (fm *FluxMonitor) hibernate() {
	fm.logger.Info("entering hibernation mode")
	fm.isHibernating = true
	fm.resetTickers(flux_aggregator_wrapper.OracleRoundState{})
}

// reactivate restarts the PollingDeviationChecker without hibernation mode
func (fm *FluxMonitor) reactivate() {
	fm.logger.Info("exiting hibernation mode, reactivating contract")
	fm.isHibernating = false
	fm.setInitialTickers()
	fm.pollIfEligible(NewZeroDeviationChecker())
}

func (fm *FluxMonitor) processLogs() {
	for !fm.backlog.Empty() {
		maybeBroadcast := fm.backlog.Take()
		broadcast, ok := maybeBroadcast.(log.Broadcast)
		if !ok {
			fm.logger.Errorf("Failed to convert backlog into LogBroadcast.  Type is %T", maybeBroadcast)
		}

		// If the log is a duplicate of one we've seen before, ignore it (this
		// happens because of the LogBroadcaster's backfilling behavior).
		consumed, err := broadcast.WasAlreadyConsumed()
		if err != nil {
			fm.logger.Errorf("Error determining if log was already consumed: %v", err)
			continue
		} else if consumed {
			fm.logger.Debug("Log was already consumed by Flux Monitor, skipping")
			continue
		}

		switch log := broadcast.DecodedLog().(type) {
		case *flux_aggregator_wrapper.FluxAggregatorNewRound:
			fm.respondToNewRoundLog(*log)
			err = broadcast.MarkConsumed()

		case *flux_aggregator_wrapper.FluxAggregatorAnswerUpdated:
			fm.respondToAnswerUpdatedLog(*log)
			err = broadcast.MarkConsumed()

		case *flags_wrapper.FlagsFlagRaised:
			// check the contract before hibernating, because one flag could be lowered
			// while the other flag remains raised
			var isFlagLowered bool
			isFlagLowered, err = fm.flags.IsLowered(fm.contractAddress)
			fm.logger.ErrorIf(err, "Error determining if flag is still raised")
			if !isFlagLowered {
				fm.hibernate()
			}
			err = broadcast.MarkConsumed()

		case *flags_wrapper.FlagsFlagLowered:
			fm.reactivate()
			err = broadcast.MarkConsumed()

		default:
			fm.logger.Errorf("unknown log %v of type %T", log, log)
		}

		fm.logger.ErrorIf(err, "Error marking log as consumed")
	}
}

// The AnswerUpdated log tells us that round has successfully closed with a new
// answer.  We update our view of the oracleRoundState in case this log was
// generated by a chain reorg.
func (fm *FluxMonitor) respondToAnswerUpdatedLog(log flux_aggregator_wrapper.FluxAggregatorAnswerUpdated) {
	fm.logger.Debugw("AnswerUpdated log", fm.loggerFieldsForAnswerUpdated(log)...)

	roundState, err := fm.roundState(0)
	if err != nil {
		logger.Errorw(fmt.Sprintf("could not fetch oracleRoundState: %v", err), fm.loggerFieldsForAnswerUpdated(log)...)
		return
	}
	fm.resetTickers(roundState)
}

// The NewRound log tells us that an oracle has initiated a new round.  This tells us that we
// need to poll and submit an answer to the contract regardless of the deviation.
func (fm *FluxMonitor) respondToNewRoundLog(log flux_aggregator_wrapper.FluxAggregatorNewRound) {
	fm.logger.Debugw("NewRound log", fm.loggerFieldsForNewRound(log)...)

	promfm.SetBigInt(promfm.SeenRound.WithLabelValues(fmt.Sprintf("%d", fm.jobID)), log.RoundId)

	//
	// NewRound answer submission logic:
	//   - Any log that reaches this point, regardless of chain reorgs or log backfilling, is one that we have
	//         not seen before.  Therefore, we should consider acting upon it.
	//   - We always take the round ID from the log, rather than the round ID suggested by `.RoundState`.  The
	//         reason is that if two NewRound logs come in in rapid succession, and we submit a tx for the first,
	//         the `.ReportableRoundID` field in the roundState() response for the 2nd log will not reflect the
	//         fact that we've submitted for the first round (assuming it hasn't been mined yet).
	//   - In the event of a reorg that pushes our previous submissions back into the mempool, we can rely on the
	//         TxManager to ensure they end up being mined into blocks, but this may cause them to revert if they
	//         are mined in an order that violates certain conditions in the FluxAggregator (restartDelay, etc.).
	//         Therefore, the cleanest solution at present is to resubmit for the reorged rounds.  The drawback
	//         of this approach is that one or the other submission tx for a given round will revert, costing the
	//         node operator some gas.  The benefit is that those submissions are guaranteed to be made, ensuring
	//         that we have high data availability (and also ensuring that node operators get paid).
	//   - There are a few straightforward cases where we don't want to submit:
	//         - When we're not eligible
	//         - When the aggregator is underfunded
	//         - When we were the initiator of the round (i.e. we've received our own NewRound log)
	//   - There are a few more nuanced cases as well:
	//         - When our node polls at the same time as another node, and both attempt to start a round.  In that
	//               case, it's possible that the other node will start the round, and our node will see the NewRound
	//               log and try to submit again.
	//         - When the poll ticker fires very soon after we've responded to a NewRound log.
	//
	//         To handle these more nuanced cases, we record round IDs and whether we've submitted for those rounds
	//         in the DB.  If we see we've already submitted for a given round, we simply bail out.
	//
	//         However, in the case of a chain reorganization, we might see logs with round IDs that we've already
	//         seen.  As mentioned above, we want to re-respond to these rounds to ensure high data availability.
	//         Therefore, if a log arrives with a round ID that is < the most recent that we submitted to, we delete
	//         all of the round IDs in the DB back to (and including) the incoming round ID.  This essentially
	//         rewinds the system back to a state wherein those reorg'ed rounds never occurred, allowing it to move
	//         forward normally.
	//
	//         There is one small exception: if the reorg is fairly shallow, and only un-starts a single round, we
	//         do not need to resubmit, because the TxManager will ensure that our existing submission gets back
	//         into the chain.  There is a very small risk that one of the nodes in the quorum (namely, whichever
	//         one started the previous round) will have its existing submission mined first, thereby violating
	//         the restartDelay, but as this risk is isolated to a single node, the round will not time out and
	//         go stale.  We consider this acceptable.
	//

	logRoundID := uint32(log.RoundId.Uint64())

	// We always want to reset the idle timer upon receiving a NewRound log, so we do it before any `return` statements.
	fm.resetIdleTimer(log.StartedAt.Uint64())

	mostRecentRoundID, err := fm.orm.MostRecentFluxMonitorRoundID(fm.contractAddress)
	if err != nil && err != gorm.ErrRecordNotFound {
		fm.logger.Errorw(fmt.Sprintf("error fetching Flux Monitor most recent round ID from DB: %v", err), fm.loggerFieldsForNewRound(log)...)
		return
	}

	if logRoundID < mostRecentRoundID {
		err = fm.orm.DeleteFluxMonitorRoundsBackThrough(fm.contractAddress, logRoundID)
		if err != nil {
			fm.logger.Errorw(fmt.Sprintf("error deleting reorged Flux Monitor rounds from DB: %v", err), fm.loggerFieldsForNewRound(log)...)
			return
		}
	}

	roundStats, jobRunStatus, err := fm.statsAndStatusForRound(logRoundID)
	if err != nil {
		fm.logger.Errorw(fmt.Sprintf("error determining round stats / run status for round: %v", err), fm.loggerFieldsForNewRound(log)...)
		return
	}

	if roundStats.NumSubmissions > 0 {
		// This indicates either that:
		//     - We tried to start a round at the same time as another node, and their transaction was mined first, or
		//     - The chain experienced a shallow reorg that unstarted the current round.
		// If our previous attempt is still pending, return early and don't re-submit
		// If our previous attempt is already over (completed or errored), we should retry
		if !jobRunStatus.Finished() {
			fm.logger.Debugw("Ignoring new round request: started round simultaneously with another node", fm.loggerFieldsForNewRound(log)...)
			return
		}
	}

	// Ignore rounds we started
	if fm.oracleAddress == log.StartedBy {
		fm.logger.Infow("Ignoring new round request: we started this round", fm.loggerFieldsForNewRound(log)...)
		return
	}

	// Ignore rounds we're not eligible for, or for which we won't be paid
	roundState, err := fm.roundState(logRoundID)
	if err != nil {
		fm.logger.Errorw(fmt.Sprintf("Ignoring new round request: error fetching eligibility from contract: %v", err), fm.loggerFieldsForNewRound(log)...)
		return
	}
	fm.resetTickers(roundState)
	err = fm.checkEligibilityAndAggregatorFunding(roundState)
	if err != nil {
		fm.logger.Infow(fmt.Sprintf("Ignoring new round request: %v", err), fm.loggerFieldsForNewRound(log)...)
		return
	}

	logger.Infow("Responding to new round request", fm.loggerFieldsForNewRound(log)...)

	// Call the v2 pipeline to execute a new job run
	runID, answer, err := fm.pipelineRun.Execute()
	if err != nil {
		fm.logger.Errorw(fmt.Sprintf("unable to fetch median price: %v", err), fm.loggerFieldsForNewRound(log)...)

		return
	}

	if !fm.isValidSubmission(logger.Default.SugaredLogger, *answer) {
		return
	}

	if roundState.PaymentAmount == nil {
		fm.logger.Error("roundState.PaymentAmount shouldn't be nil")
	}

	err = fm.submitTransaction(runID, *answer, roundState.RoundId)
	if err != nil {
		fm.logger.Errorw(fmt.Sprintf("unable to create job run: %v", err), fm.loggerFieldsForNewRound(log)...)

		return
	}
}

var (
	// ErrNotEligible defines when the round is not eligible for submission
	ErrNotEligible = errors.New("not eligible to submit")
	// ErrUnderfunded defines when the aggregator does not have sufficient funds
	ErrUnderfunded = errors.New("aggregator is underfunded")
	// ErrPaymentTooLow defines when the round payment is too low
	ErrPaymentTooLow = errors.New("round payment amount < minimum contract payment")
)

func (fm *FluxMonitor) checkEligibilityAndAggregatorFunding(roundState flux_aggregator_wrapper.OracleRoundState) error {
	if !roundState.EligibleToSubmit {
		return ErrNotEligible
	} else if !fm.paymentChecker.SufficientFunds(
		roundState.AvailableFunds,
		roundState.PaymentAmount,
		roundState.OracleCount,
	) {
		return ErrUnderfunded
	} else if !fm.paymentChecker.SufficientPayment(roundState.PaymentAmount) {
		return ErrPaymentTooLow
	}
	return nil
}

func (fm *FluxMonitor) pollIfEligible(deviationChecker *DeviationChecker) {
	l := fm.logger.With(
		"threshold", deviationChecker.Thresholds.Rel,
		"absoluteThreshold", deviationChecker.Thresholds.Abs,
	)

	if !fm.connected.IsSet() {
		l.Warnw("not connected to Ethereum node, skipping poll")

		return
	}

	//
	// Poll ticker submission logic:
	//   - We avoid saving on-chain state wherever possible.  Therefore, we do not know which round we should be
	//         submitting for when the pollTicker fires.
	//   - We pass 0 into `roundState()`, and the FluxAggregator returns a suggested roundID for us to
	//         submit to, as well as our eligibility to submit to that round.
	//   - If the poll ticker fires very soon after we've responded to a NewRound log, and our tx has not been
	//         mined, we risk double-submitting for a round.  To detect this, we check the DB to see whether
	//         we've responded to this round already, and bail out if so.
	//

	// Ask the FluxAggregator which round we should be submitting to, and what the state of that round is.
	roundState, err := fm.roundState(0)
	if err != nil {
		l.Errorw("unable to determine eligibility to submit from FluxAggregator contract", "err", err)
		fm.jobORM.RecordError(
			context.Background(),
			fm.jobID,
			"Unable to call roundState method on provided contract. Check contract address.",
		)

		return
	}

	fm.resetTickers(roundState)
	l = l.With("reportableRound", roundState.RoundId)

	roundStats, jobRunStatus, err := fm.statsAndStatusForRound(roundState.RoundId)
	if err != nil {
		l.Errorw("error determining round stats / run status for round", "err", err)

		return
	}

	// If we've already successfully submitted to this round (ie through a NewRound log)
	// and the associated JobRun hasn't errored, skip polling
	if roundStats.NumSubmissions > 0 && !jobRunStatus.Errored() {
		l.Infow("skipping poll: round already answered, tx unconfirmed")

		return
	}

	// Don't submit if we're not eligible, or won't get paid
	err = fm.checkEligibilityAndAggregatorFunding(roundState)
	if err != nil {
		l.Infow(fmt.Sprintf("skipping poll: %v", err))

		return
	}

	// Call the v2 pipeline to execute a new pipeline run
	// Note: we expect the FM pipeline to scale the fetched answer by the same
	// amount as precision
	runID, answer, err := fm.pipelineRun.Execute()
	if err != nil {
		l.Errorw("can't fetch answer", "err", err)
		fm.jobORM.RecordError(context.Background(), fm.jobID, "Error polling")

		return
	}

	if !fm.isValidSubmission(l, *answer) {
		return
	}

	jobID := fmt.Sprintf("%d", fm.jobID)
	latestAnswer := decimal.NewFromBigInt(roundState.LatestSubmission, -fm.precision)
	promfm.SetDecimal(promfm.SeenValue.WithLabelValues(jobID), *answer)

	l = l.With(
		"latestAnswer", latestAnswer,
		"answer", answer,
	)

	if roundState.RoundId > 1 && !deviationChecker.OutsideDeviation(latestAnswer, *answer) {
		l.Debugw("deviation < threshold, not submitting")

		return
	}

	if roundState.RoundId > 1 {
		l.Infow("deviation > threshold, starting new round")
	} else {
		l.Infow("starting first round")
	}

	// --> Create an ETH transaction by calling a 'submit' method on the contract

	if roundState.PaymentAmount == nil {
		l.Error("roundState.PaymentAmount shouldn't be nil")
	}

	err = fm.submitTransaction(runID, *answer, roundState.RoundId)
	if err != nil {
		l.Errorw("can't create job run", "err", err)

		return
	}

	promfm.SetDecimal(promfm.ReportedValue.WithLabelValues(jobID), *answer)
	promfm.SetUint32(promfm.ReportedRound.WithLabelValues(jobID), roundState.RoundId)
}

// If the answer is outside the allowable range, log an error and don't submit.
// to avoid an onchain reversion.
func (fm *FluxMonitor) isValidSubmission(l *zap.SugaredLogger, answer decimal.Decimal) bool {
	if fm.submissionChecker.IsValid(answer) {
		return true
	}

	l.Errorw("answer is outside acceptable range",
		"min", fm.submissionChecker.Min,
		"max", fm.submissionChecker.Max,
		"answer", answer,
	)
	fm.jobORM.RecordError(context.Background(), fm.jobID, "Answer is outside acceptable range")

	return false
}

func (fm *FluxMonitor) roundState(roundID uint32) (flux_aggregator_wrapper.OracleRoundState, error) {
	return fm.fluxAggregator.OracleRoundState(nil, fm.oracleAddress, roundID)
}

// initialRoundState fetches the round information that the fluxmonitor should use when starting
// new jobs. Choosing the correct round on startup is key to setting timers correctly.
func (fm *FluxMonitor) initialRoundState() flux_aggregator_wrapper.OracleRoundState {
	defaultRoundState := flux_aggregator_wrapper.OracleRoundState{
		StartedAt: uint64(time.Now().Unix()),
	}
	latestRoundData, err := fm.fluxAggregator.LatestRoundData(nil)
	if err != nil {
		fm.logger.Warnf(
			"unable to retrieve latestRoundData for FluxAggregator contract - defaulting "+
				"to current time for tickers: %v",
			err,
		)
		return defaultRoundState
	}
	roundID := uint32(latestRoundData.RoundId.Uint64())
	latestRoundState, err := fm.fluxAggregator.OracleRoundState(nil, fm.oracleAddress, roundID)
	if err != nil {
		fm.logger.Warnf(
			"unable to call roundState for latest round, round: %d, err: %v",
			latestRoundData.RoundId,
			err,
		)
		return defaultRoundState
	}
	return latestRoundState
}

func (fm *FluxMonitor) resetTickers(roundState flux_aggregator_wrapper.OracleRoundState) {
	fm.resetPollTicker()
	fm.resetHibernationTimer()
	fm.resetIdleTimer(roundState.StartedAt)
	fm.resetRoundTimer(roundStateTimesOutAt(roundState))
}

func (fm *FluxMonitor) setInitialTickers() {
	fm.resetTickers(fm.initialRoundState())
}

func (fm *FluxMonitor) resetPollTicker() {
	if fm.pollTicker.IsEnabled() && !fm.isHibernating {
		fm.pollTicker.Resume()
	} else {
		fm.pollTicker.Pause()
	}
}

func (fm *FluxMonitor) resetHibernationTimer() {
	if !fm.isHibernating {
		fm.hibernationTimer.Stop()
	} else {
		fm.hibernationTimer.Reset(hibernationPollPeriod)
	}
}

func (fm *FluxMonitor) resetRoundTimer(roundTimesOutAt uint64) {
	if fm.isHibernating {
		fm.roundTimer.Stop()
		return
	}

	loggerFields := fm.loggerFields("timesOutAt", roundTimesOutAt)

	if roundTimesOutAt == 0 {
		fm.roundTimer.Stop()
		fm.logger.Debugw("disabling roundTimer, no active round", loggerFields...)

	} else {
		timesOutAt := time.Unix(int64(roundTimesOutAt), 0)
		timeUntilTimeout := time.Until(timesOutAt)

		if timeUntilTimeout <= 0 {
			fm.roundTimer.Stop()
			fm.logger.Debugw("roundTimer has run down; disabling", loggerFields...)
		} else {
			fm.roundTimer.Reset(timeUntilTimeout)
			loggerFields = append(loggerFields, "value", roundTimesOutAt)
			fm.logger.Debugw("updating roundState.TimesOutAt", loggerFields...)
		}
	}
}

func (fm *FluxMonitor) resetIdleTimer(roundStartedAtUTC uint64) {
	if fm.isHibernating || fm.idleTimer.IsDisabled() {
		fm.idleTimer.Stop()
		return
	} else if roundStartedAtUTC == 0 {
		// There is no active round, so keep using the idleTimer we already have
		return
	}

	startedAt := time.Unix(int64(roundStartedAtUTC), 0)
	idleDeadline := startedAt.Add(fm.idleTimer.Period())
	timeUntilIdleDeadline := time.Until(idleDeadline)
	loggerFields := fm.loggerFields(
		"startedAt", roundStartedAtUTC,
		"timeUntilIdleDeadline", timeUntilIdleDeadline,
	)

	if timeUntilIdleDeadline <= 0 {
		fm.logger.Debugw("not resetting idleTimer, negative duration", loggerFields...)
		return
	}
	fm.idleTimer.Reset(timeUntilIdleDeadline)
	fm.logger.Debugw("resetting idleTimer", loggerFields...)
}

func (fm *FluxMonitor) submitTransaction(
	runID int64,
	answer decimal.Decimal,
	roundID uint32,
) error {
	// Submit the Eth Tx
	err := fm.contractSubmitter.Submit(
		new(big.Int).SetInt64(int64(roundID)),
		answer.BigInt(),
	)
	if err != nil {
		return err
	}

	// Update the flux monitor round stats
	err = fm.orm.UpdateFluxMonitorRoundStats(
		fm.contractAddress,
		roundID,
		runID,
	)
	if err != nil {
		fm.logger.Errorw(
			fmt.Sprintf("error updating FM round submission count: %v", err),
			"roundID", roundID,
		)

		return err
	}

	return nil
}

func (fm *FluxMonitor) loggerFields(added ...interface{}) []interface{} {
	return append(added, []interface{}{
		"pollFrequency", fm.pollTicker.Interval,
		"idleDuration", fm.idleTimer.Period,
	}...)
}

func (fm *FluxMonitor) loggerFieldsForNewRound(log flux_aggregator_wrapper.FluxAggregatorNewRound) []interface{} {
	return []interface{}{
		"round", log.RoundId,
		"startedBy", log.StartedBy.Hex(),
		"startedAt", log.StartedAt.String(),
	}
}

func (fm *FluxMonitor) loggerFieldsForAnswerUpdated(log flux_aggregator_wrapper.FluxAggregatorAnswerUpdated) []interface{} {
	return []interface{}{
		"round", log.RoundId,
		"answer", log.Current.String(),
		"timestamp", log.UpdatedAt.String(),
	}
}

func (fm *FluxMonitor) loggerFieldsForTick() []interface{} {
	return []interface{}{
		"pollPeriod", fm.pollTicker.Interval,
		"idleDuration", fm.idleTimer.Period,
	}
}

func (fm *FluxMonitor) statsAndStatusForRound(roundID uint32) (
	FluxMonitorRoundStatsV2,
	pipeline.RunStatus,
	error,
) {
	roundStats, err := fm.orm.FindOrCreateFluxMonitorRoundStats(fm.contractAddress, roundID)
	if err != nil {
		return FluxMonitorRoundStatsV2{}, pipeline.RunStatusUnknown, err
	}

	// JobRun will not exist if this is the first time responding to this round
	var run pipeline.Run
	if roundStats.PipelineRunID.Valid {
		run, err = fm.pipelineORM.FindRun(roundStats.PipelineRunID.Int64)
		if err != nil {
			return FluxMonitorRoundStatsV2{}, pipeline.RunStatusUnknown, err
		}
	}

	return roundStats, run.Status(), nil
}

func roundStateTimesOutAt(rs flux_aggregator_wrapper.OracleRoundState) uint64 {
	return rs.StartedAt + rs.Timeout
}
