package fluxmonitorv2

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	mrand "math/rand"
	"reflect"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/flags_wrapper"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/flux_aggregator_wrapper"
	evmclient "github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/log"
	evmutils "github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/recovery"
	"github.com/smartcontractkit/chainlink/v2/core/services/fluxmonitorv2/promfm"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// PollRequest defines a request to initiate a poll
type PollRequest struct {
	Type      PollRequestType
	Timestamp time.Time
}

// PollRequestType defines which method was used to request a poll
type PollRequestType int

const (
	PollRequestTypeUnknown PollRequestType = iota
	PollRequestTypeInitial
	PollRequestTypePoll
	PollRequestTypeIdle
	PollRequestTypeRound
	PollRequestTypeHibernation
	PollRequestTypeRetry
	PollRequestTypeAwaken
	PollRequestTypeDrumbeat
)

// DefaultHibernationPollPeriod defines the hibernation polling period
const DefaultHibernationPollPeriod = 24 * time.Hour

// FluxMonitor polls external price adapters via HTTP to check for price swings.
type FluxMonitor struct {
	services.Service
	eng    *services.Engine
	logger logger.SugaredLogger

	contractAddress   common.Address
	oracleAddress     common.Address
	jobSpec           job.Job
	spec              pipeline.Spec
	runner            pipeline.Runner
	ds                sqlutil.DataSource
	orm               ORM
	jobORM            job.ORM
	pipelineORM       pipeline.ORM
	keyStore          KeyStoreInterface
	pollManager       *PollManager
	paymentChecker    *PaymentChecker
	contractSubmitter ContractSubmitter
	deviationChecker  *DeviationChecker
	submissionChecker *SubmissionChecker
	flags             Flags
	fluxAggregator    flux_aggregator_wrapper.FluxAggregatorInterface
	logBroadcaster    log.Broadcaster
	chainID           *big.Int

	backlog       *utils.BoundedPriorityQueue[log.Broadcast]
	chProcessLogs chan struct{}
}

// NewFluxMonitor returns a new instance of PollingDeviationChecker.
func NewFluxMonitor(
	pipelineRunner pipeline.Runner,
	jobSpec job.Job,
	spec pipeline.Spec,
	ds sqlutil.DataSource,
	orm ORM,
	jobORM job.ORM,
	pipelineORM pipeline.ORM,
	keyStore KeyStoreInterface,
	pollManager *PollManager,
	paymentChecker *PaymentChecker,
	contractAddress common.Address,
	contractSubmitter ContractSubmitter,
	deviationChecker *DeviationChecker,
	submissionChecker *SubmissionChecker,
	flags Flags,
	fluxAggregator flux_aggregator_wrapper.FluxAggregatorInterface,
	logBroadcaster log.Broadcaster,
	lggr logger.Logger,
	chainID *big.Int,
) (*FluxMonitor, error) {
	fm := &FluxMonitor{
		ds:                ds,
		runner:            pipelineRunner,
		jobSpec:           jobSpec,
		spec:              spec,
		orm:               orm,
		jobORM:            jobORM,
		pipelineORM:       pipelineORM,
		keyStore:          keyStore,
		pollManager:       pollManager,
		paymentChecker:    paymentChecker,
		contractAddress:   contractAddress,
		contractSubmitter: contractSubmitter,
		deviationChecker:  deviationChecker,
		submissionChecker: submissionChecker,
		flags:             flags,
		logBroadcaster:    logBroadcaster,
		fluxAggregator:    fluxAggregator,
		chainID:           chainID,
		backlog: utils.NewBoundedPriorityQueue[log.Broadcast](map[uint]int{
			// We want reconnecting nodes to be able to submit to a round
			// that hasn't hit maxAnswers yet, as well as the newest round.
			PriorityNewRoundLog:      2,
			PriorityAnswerUpdatedLog: 1,
			PriorityFlagChangedLog:   2,
		}),
		chProcessLogs: make(chan struct{}, 1),
	}
	fm.Service, fm.eng = services.Config{
		Name:  "FluxMonitor",
		Start: fm.start,
		Close: fm.close,
	}.NewServiceEngine(lggr)
	fm.logger = logger.Sugared(fm.eng)

	return fm, nil
}

// NewFromJobSpec constructs an instance of FluxMonitor with sane defaults and
// validation.
func NewFromJobSpec(
	jobSpec job.Job,
	ds sqlutil.DataSource,
	orm ORM,
	jobORM job.ORM,
	pipelineORM pipeline.ORM,
	keyStore KeyStoreInterface,
	ethClient evmclient.Client,
	logBroadcaster log.Broadcaster,
	pipelineRunner pipeline.Runner,
	cfg Config,
	fcfg EvmFeeConfig,
	jcfg JobPipelineConfig,
	lggr logger.Logger,
) (*FluxMonitor, error) {
	fmSpec := jobSpec.FluxMonitorSpec
	chainId := ethClient.ConfiguredChainID()

	if !validatePollTimer(fmSpec.PollTimerDisabled, MinimumPollingInterval(jcfg), fmSpec.PollTimerPeriod) {
		return nil, fmt.Errorf(
			"PollTimerPeriod (%s), must be equal or greater than JobPipeline.HTTPRequest.DefaultTimeout (%s) ",
			fmSpec.PollTimerPeriod,
			MinimumPollingInterval(jcfg),
		)
	}

	// Set up the flux aggregator
	fluxAggregator, err := flux_aggregator_wrapper.NewFluxAggregator(
		fmSpec.ContractAddress.Address(),
		ethClient,
	)
	if err != nil {
		return nil, err
	}

	gasLimit := fcfg.LimitDefault()
	fmLimit := fcfg.LimitJobType().FM()
	if jobSpec.GasLimit.Valid {
		gasLimit = uint64(jobSpec.GasLimit.Uint32)
	} else if fmLimit != nil {
		gasLimit = uint64(*fmLimit)
	}

	contractSubmitter := NewFluxAggregatorContractSubmitter(
		fluxAggregator,
		orm,
		keyStore,
		gasLimit,
		jobSpec.ForwardingAllowed,
		chainId,
	)

	flags, err := NewFlags(cfg.FlagsContractAddress(), ethClient)
	logger.Sugared(lggr).ErrorIf(err,
		"Error creating Flags contract instance, check address: "+cfg.FlagsContractAddress(),
	)

	paymentChecker := &PaymentChecker{
		MinContractPayment: cfg.MinContractPayment(),
		MinJobPayment:      fmSpec.MinPayment,
	}

	min, err := fluxAggregator.MinSubmissionValue(nil)
	if err != nil {
		return nil, err
	}

	max, err := fluxAggregator.MaxSubmissionValue(nil)
	if err != nil {
		return nil, err
	}

	fmLogger := logger.With(lggr,
		"jobID", jobSpec.ID,
		"contract", fmSpec.ContractAddress.Hex(),
	)

	pollManager, err := NewPollManager(
		PollManagerConfig{
			PollTickerInterval:      fmSpec.PollTimerPeriod,
			PollTickerDisabled:      fmSpec.PollTimerDisabled,
			IdleTimerPeriod:         fmSpec.IdleTimerPeriod,
			IdleTimerDisabled:       fmSpec.IdleTimerDisabled,
			DrumbeatSchedule:        fmSpec.DrumbeatSchedule,
			DrumbeatEnabled:         fmSpec.DrumbeatEnabled,
			DrumbeatRandomDelay:     fmSpec.DrumbeatRandomDelay,
			HibernationPollPeriod:   DefaultHibernationPollPeriod, // Not currently configurable
			MinRetryBackoffDuration: 1 * time.Minute,
			MaxRetryBackoffDuration: 1 * time.Hour,
		},
		fmLogger,
	)
	if err != nil {
		return nil, err
	}

	return NewFluxMonitor(
		pipelineRunner,
		jobSpec,
		*jobSpec.PipelineSpec,
		ds,
		orm,
		jobORM,
		pipelineORM,
		keyStore,
		pollManager,
		paymentChecker,
		fmSpec.ContractAddress.Address(),
		contractSubmitter,
		NewDeviationChecker(
			float64(fmSpec.Threshold),
			float64(fmSpec.AbsoluteThreshold),
			fmLogger,
		),
		NewSubmissionChecker(min, max),
		flags,
		fluxAggregator,
		logBroadcaster,
		fmLogger,
		chainId,
	)
}

const (
	PriorityFlagChangedLog   uint = 0
	PriorityNewRoundLog      uint = 1
	PriorityAnswerUpdatedLog uint = 2
)

// Start implements the job.Service interface. It begins the CSP consumer in a
// single goroutine to poll the price adapters and listen to NewRound events.
func (fm *FluxMonitor) start(context.Context) error {
	fm.eng.Go(fm.consume)
	return nil
}

func (fm *FluxMonitor) IsHibernating() bool {
	if !fm.flags.ContractExists() {
		return false
	}

	isFlagLowered, err := fm.flags.IsLowered(fm.contractAddress)
	if err != nil {
		fm.logger.Errorf("unable to determine hibernation status: %v", err)

		return false
	}

	return !isFlagLowered
}

// close stops this instance from
// polling, cleaning up resources.
func (fm *FluxMonitor) close() error {
	fm.pollManager.Stop()

	return nil
}

// JobID implements the listener.Listener interface.
func (fm *FluxMonitor) JobID() int32 { return fm.spec.JobID }

// HandleLog processes the contract logs
func (fm *FluxMonitor) HandleLog(ctx context.Context, broadcast log.Broadcast) {
	log := broadcast.DecodedLog()
	if log == nil || reflect.ValueOf(log).IsNil() {
		fm.logger.Panic("HandleLog: failed to handle log of type nil")
	}

	switch log := log.(type) {
	case *flux_aggregator_wrapper.FluxAggregatorNewRound:
		fm.backlog.Add(PriorityNewRoundLog, broadcast)

	case *flux_aggregator_wrapper.FluxAggregatorAnswerUpdated:
		fm.backlog.Add(PriorityAnswerUpdatedLog, broadcast)

	case *flags_wrapper.FlagsFlagRaised:
		if log.Subject == evmutils.ZeroAddress || log.Subject == fm.contractAddress {
			fm.backlog.Add(PriorityFlagChangedLog, broadcast)
		}

	case *flags_wrapper.FlagsFlagLowered:
		if log.Subject == evmutils.ZeroAddress || log.Subject == fm.contractAddress {
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

func (fm *FluxMonitor) consume(ctx context.Context) {
	if err := fm.SetOracleAddress(ctx); err != nil {
		fm.logger.Warnw(
			"unable to set oracle address, this flux monitor job may not work correctly",
			"err", err,
		)
	}

	// Subscribe to contract logs
	unsubscribe := fm.logBroadcaster.Register(fm, log.ListenerOpts{
		Contract: fm.fluxAggregator.Address(),
		ParseLog: fm.fluxAggregator.ParseLog,
		LogsWithTopics: map[common.Hash][][]log.Topic{
			flux_aggregator_wrapper.FluxAggregatorNewRound{}.Topic():      nil,
			flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{}.Topic(): nil,
		},
		MinIncomingConfirmations: 0,
	})
	defer unsubscribe()

	if fm.flags.ContractExists() {
		unsubscribe := fm.logBroadcaster.Register(fm, log.ListenerOpts{
			Contract: fm.flags.Address(),
			ParseLog: fm.flags.ParseLog,
			LogsWithTopics: map[common.Hash][][]log.Topic{
				flags_wrapper.FlagsFlagLowered{}.Topic(): nil,
				flags_wrapper.FlagsFlagRaised{}.Topic():  nil,
			},
			MinIncomingConfirmations: 0,
		})
		defer unsubscribe()
	}

	fm.pollManager.Start(fm.IsHibernating(), fm.initialRoundState())

	tickLogger := fm.logger.With(
		"pollInterval", fm.pollManager.cfg.PollTickerInterval,
		"idlePeriod", fm.pollManager.cfg.IdleTimerPeriod,
	)

	for {
		select {
		case <-ctx.Done():
			return

		case <-fm.chProcessLogs:
			recovery.WrapRecover(fm.logger, func() { fm.processLogs(ctx) })

		case at := <-fm.pollManager.PollTickerTicks():
			tickLogger.Debugf("Poll ticker fired on %v", formatTime(at))
			recovery.WrapRecover(fm.logger, func() {
				fm.pollIfEligible(ctx, PollRequestTypePoll, fm.deviationChecker, nil)
			})

		case at := <-fm.pollManager.IdleTimerTicks():
			tickLogger.Debugf("Idle timer fired on %v", formatTime(at))
			recovery.WrapRecover(fm.logger, func() {
				fm.pollIfEligible(ctx, PollRequestTypeIdle, NewZeroDeviationChecker(fm.logger), nil)
			})

		case at := <-fm.pollManager.RoundTimerTicks():
			tickLogger.Debugf("Round timer fired on %v", formatTime(at))
			recovery.WrapRecover(fm.logger, func() {
				fm.pollIfEligible(ctx, PollRequestTypeRound, fm.deviationChecker, nil)
			})

		case at := <-fm.pollManager.HibernationTimerTicks():
			tickLogger.Debugf("Hibernation timer fired on %v", formatTime(at))
			recovery.WrapRecover(fm.logger, func() {
				fm.pollIfEligible(ctx, PollRequestTypeHibernation, NewZeroDeviationChecker(fm.logger), nil)
			})

		case at := <-fm.pollManager.RetryTickerTicks():
			tickLogger.Debugf("Retry ticker fired on %v", formatTime(at))
			recovery.WrapRecover(fm.logger, func() {
				fm.pollIfEligible(ctx, PollRequestTypeRetry, NewZeroDeviationChecker(fm.logger), nil)
			})

		case at := <-fm.pollManager.DrumbeatTicks():
			tickLogger.Debugf("Drumbeat ticker fired on %v", formatTime(at))
			recovery.WrapRecover(fm.logger, func() {
				fm.pollIfEligible(ctx, PollRequestTypeDrumbeat, NewZeroDeviationChecker(fm.logger), nil)
			})

		case request := <-fm.pollManager.Poll():
			switch request.Type {
			case PollRequestTypeUnknown:
				break
			default:
				recovery.WrapRecover(fm.logger, func() {
					fm.pollIfEligible(ctx, request.Type, fm.deviationChecker, nil)
				})
			}
		}
	}
}

func formatTime(at time.Time) string {
	ago := time.Since(at)
	return fmt.Sprintf("%v (%v ago)", at.UTC().Format(time.RFC3339), ago)
}

// SetOracleAddress sets the oracle address which matches the node's keys.
// If none match, it uses the first available key
func (fm *FluxMonitor) SetOracleAddress(ctx context.Context) error {
	oracleAddrs, err := fm.fluxAggregator.GetOracles(nil)
	if err != nil {
		fm.logger.Error("failed to get list of oracles from FluxAggregator contract")
		return errors.Wrap(err, "failed to get list of oracles from FluxAggregator contract")
	}
	keys, err := fm.keyStore.EnabledKeysForChain(ctx, fm.chainID)
	if err != nil {
		return errors.Wrap(err, "failed to load keys")
	}
	for _, k := range keys {
		for _, oracleAddr := range oracleAddrs {
			if k.Address == oracleAddr {
				fm.oracleAddress = oracleAddr
				return nil
			}
		}
	}

	log := fm.logger.With(
		"keys", keys,
		"oracleAddresses", oracleAddrs,
	)

	if len(keys) > 0 {
		addr := keys[0].Address
		log.Warnw("None of the node's keys matched any oracle addresses, using first available key. This flux monitor job may not work correctly",
			"address", addr.Hex(),
		)
		fm.oracleAddress = addr

		return nil
	}

	log.Error("No keys found. This flux monitor job may not work correctly")
	return errors.New("No keys found")
}

func (fm *FluxMonitor) processLogs(ctx context.Context) {
	for ctx.Err() == nil && !fm.backlog.Empty() {
		broadcast := fm.backlog.Take()
		fm.processBroadcast(ctx, broadcast)
	}
}

func (fm *FluxMonitor) processBroadcast(ctx context.Context, broadcast log.Broadcast) {
	// If the log is a duplicate of one we've seen before, ignore it (this
	// happens because of the LogBroadcaster's backfilling behavior).
	consumed, err := fm.logBroadcaster.WasAlreadyConsumed(ctx, broadcast)

	if err != nil {
		fm.logger.Errorf("Error determining if log was already consumed: %v", err)
		return
	} else if consumed {
		fm.logger.Debug("Log was already consumed by Flux Monitor, skipping")
		return
	}

	started := time.Now()
	decodedLog := broadcast.DecodedLog()
	switch log := decodedLog.(type) {
	case *flux_aggregator_wrapper.FluxAggregatorNewRound:
		fm.respondToNewRoundLog(ctx, *log, broadcast)
	case *flux_aggregator_wrapper.FluxAggregatorAnswerUpdated:
		fm.respondToAnswerUpdatedLog(*log)
		fm.markLogAsConsumed(ctx, broadcast, decodedLog, started)
	case *flags_wrapper.FlagsFlagRaised:
		fm.respondToFlagsRaisedLog()
		fm.markLogAsConsumed(ctx, broadcast, decodedLog, started)
	case *flags_wrapper.FlagsFlagLowered:
		// Only reactivate if it is hibernating
		if fm.pollManager.isHibernating.Load() {
			fm.pollManager.Awaken(fm.initialRoundState())
			fm.pollIfEligible(ctx, PollRequestTypeAwaken, NewZeroDeviationChecker(fm.logger), broadcast)
		}
	default:
		fm.logger.Errorf("unknown log %v of type %T", log, log)
	}
}

func (fm *FluxMonitor) markLogAsConsumed(ctx context.Context, broadcast log.Broadcast, decodedLog interface{}, started time.Time) {
	if err := fm.logBroadcaster.MarkConsumed(ctx, nil, broadcast); err != nil {
		fm.logger.Errorw("Failed to mark log as consumed",
			"err", err, "logType", fmt.Sprintf("%T", decodedLog), "log", broadcast.String(), "elapsed", time.Since(started))
	}
}

func (fm *FluxMonitor) respondToFlagsRaisedLog() {
	fm.logger.Debug("FlagsFlagRaised log")
	// check the contract before hibernating, because one flag could be lowered
	// while the other flag remains raised
	isFlagLowered, err := fm.flags.IsLowered(fm.contractAddress)
	fm.logger.ErrorIf(err, "Error determining if flag is still raised")
	if !isFlagLowered {
		fm.pollManager.Hibernate()
	}
}

// The AnswerUpdated log tells us that round has successfully closed with a new
// answer.  We update our view of the oracleRoundState in case this log was
// generated by a chain reorg.
func (fm *FluxMonitor) respondToAnswerUpdatedLog(log flux_aggregator_wrapper.FluxAggregatorAnswerUpdated) {
	answerUpdatedLogger := fm.logger.With(
		"round", log.RoundId,
		"answer", log.Current.String(),
		"timestamp", log.UpdatedAt.String(),
	)

	answerUpdatedLogger.Debug("AnswerUpdated log")

	roundState, err := fm.roundState(0)
	if err != nil {
		answerUpdatedLogger.Errorf("could not fetch oracleRoundState: %v", err)

		return
	}

	fm.pollManager.Reset(roundState)
}

// The NewRound log tells us that an oracle has initiated a new round.  This tells us that we
// need to poll and submit an answer to the contract regardless of the deviation.
func (fm *FluxMonitor) respondToNewRoundLog(ctx context.Context, log flux_aggregator_wrapper.FluxAggregatorNewRound, lb log.Broadcast) {
	started := time.Now()

	newRoundLogger := fm.logger.With(
		"round", log.RoundId,
		"startedBy", log.StartedBy.Hex(),
		"startedAt", log.StartedAt.String(),
		"startedAtUtc", time.Unix(log.StartedAt.Int64(), 0).UTC().Format(time.RFC3339),
	)
	var markConsumed = true
	defer func() {
		if markConsumed {
			if err := fm.logBroadcaster.MarkConsumed(ctx, nil, lb); err != nil {
				fm.logger.Errorw("Failed to mark log consumed", "err", err, "log", lb.String())
			}
		}
	}()

	newRoundLogger.Debug("NewRound log")
	promfm.SetBigInt(promfm.SeenRound.WithLabelValues(strconv.Itoa(int(fm.spec.JobID))), log.RoundId)

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
	fm.pollManager.ResetIdleTimer(log.StartedAt.Uint64())

	mostRecentRoundID, err := fm.orm.MostRecentFluxMonitorRoundID(ctx, fm.contractAddress)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		newRoundLogger.Errorf("error fetching Flux Monitor most recent round ID from DB: %v", err)
		return
	}

	roundStats, jobRunStatus, err := fm.statsAndStatusForRound(ctx, logRoundID, 1)
	if err != nil {
		newRoundLogger.Errorf("error determining round stats / run status for round: %v", err)
		return
	}

	if logRoundID < mostRecentRoundID && roundStats.NumNewRoundLogs > 0 {
		newRoundLogger.Debugf("Received an older round log (and number of previously received NewRound logs is: %v) - "+
			"a possible reorg, hence deleting round ids from %v to %v", roundStats.NumNewRoundLogs, logRoundID, mostRecentRoundID)
		err = fm.orm.DeleteFluxMonitorRoundsBackThrough(ctx, fm.contractAddress, logRoundID)
		if err != nil {
			newRoundLogger.Errorf("error deleting reorged Flux Monitor rounds from DB: %v", err)
			return
		}

		// as all newer stats were deleted, at this point a new round stats entry will be created
		roundStats, err = fm.orm.FindOrCreateFluxMonitorRoundStats(ctx, fm.contractAddress, logRoundID, 1)
		if err != nil {
			newRoundLogger.Errorf("error determining subsequent round stats for round: %v", err)
			return
		}
	}

	if roundStats.NumSubmissions > 0 {
		// This indicates either that:
		//     - We tried to start a round at the same time as another node, and their transaction was mined first, or
		//     - The chain experienced a shallow reorg that unstarted the current round.
		// If our previous attempt is still pending, return early and don't re-submit
		// If our previous attempt is already over (completed or errored), we should retry
		newRoundLogger.Debugf("There are already %v existing submissions to this round, while job run status is: %v", roundStats.NumSubmissions, jobRunStatus)
		if !jobRunStatus.Finished() {
			newRoundLogger.Debug("Ignoring new round request: started round simultaneously with another node")
			return
		}
	}

	// Ignore rounds we started
	if fm.oracleAddress == log.StartedBy {
		newRoundLogger.Info("Ignoring new round request: we started this round")
		return
	}

	// Ignore rounds we're not eligible for, or for which we won't be paid
	roundState, err := fm.roundState(logRoundID)
	if err != nil {
		newRoundLogger.Errorf("Ignoring new round request: error fetching eligibility from contract: %v", err)
		return
	}

	fm.pollManager.Reset(roundState)
	err = fm.checkEligibilityAndAggregatorFunding(roundState)
	if err != nil {
		newRoundLogger.Infof("Ignoring new round request: %v", err)
		return
	}

	newRoundLogger.Info("Responding to new round request")

	// Best effort to attach metadata.
	var metaDataForBridge map[string]interface{}
	lrd, err := fm.fluxAggregator.LatestRoundData(nil)
	if err != nil {
		newRoundLogger.Warnw("Couldn't read latest round data for request meta", "err", err)
	} else {
		metaDataForBridge, err = bridges.MarshalBridgeMetaData(lrd.Answer, lrd.UpdatedAt)
		if err != nil {
			newRoundLogger.Warnw("Error marshalling roundState for request meta", "err", err)
		}
	}

	vars := pipeline.NewVarsFrom(map[string]interface{}{
		"jobSpec": map[string]interface{}{
			"databaseID":    fm.jobSpec.ID,
			"externalJobID": fm.jobSpec.ExternalJobID,
			"name":          fm.jobSpec.Name.ValueOrZero(),
			"evmChainID":    fm.chainID.String(),
		},
		"jobRun": map[string]interface{}{
			"meta": metaDataForBridge,
		},
	})

	// Call the v2 pipeline to execute a new job run
	run, results, err := fm.runner.ExecuteRun(ctx, fm.spec, vars)
	if err != nil {
		newRoundLogger.Errorw(fmt.Sprintf("error executing new run for job ID %v name %v", fm.spec.JobID, fm.spec.JobName), "err", err)
		return
	}
	result, err := results.FinalResult().SingularResult()
	if err != nil || result.Error != nil {
		newRoundLogger.Errorw("can't fetch answer", "err", err, "result", result)
		fm.jobORM.TryRecordError(ctx, fm.spec.JobID, "Error polling")
		return
	}
	answer, err := utils.ToDecimal(result.Value)
	if err != nil {
		newRoundLogger.Errorw(fmt.Sprintf("error executing new run for job ID %v name %v", fm.spec.JobID, fm.spec.JobName), "err", err)
		return
	}

	if !fm.isValidSubmission(ctx, newRoundLogger, answer, started) {
		return
	}

	if roundState.PaymentAmount == nil {
		newRoundLogger.Error("roundState.PaymentAmount shouldn't be nil")
	}

	err = fm.Transact(ctx, func(tx sqlutil.DataSource) error {
		if err2 := fm.runner.InsertFinishedRun(ctx, tx, run, false); err2 != nil {
			return err2
		}
		if err2 := fm.queueTransactionForTxm(ctx, tx, run.ID, answer, roundState.RoundId, &log); err2 != nil {
			return err2
		}
		return fm.logBroadcaster.MarkConsumed(ctx, tx, lb)
	})
	// Either the tx failed and we want to reprocess the log, or it succeeded and already marked it consumed
	markConsumed = false
	if err != nil {
		newRoundLogger.Errorf("unable to create job run: %v", err)
		return
	}
}

func (fm *FluxMonitor) Transact(ctx context.Context, fn func(sqlutil.DataSource) error) error {
	return sqlutil.TransactDataSource(ctx, fm.ds, nil, fn)
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

func (fm *FluxMonitor) pollIfEligible(ctx context.Context, pollReq PollRequestType, deviationChecker *DeviationChecker, broadcast log.Broadcast) {
	started := time.Now()

	l := fm.logger.With(
		"threshold", deviationChecker.Thresholds.Rel,
		"absoluteThreshold", deviationChecker.Thresholds.Abs,
	)
	var markConsumed = true
	defer func() {
		if markConsumed && broadcast != nil {
			if err := fm.logBroadcaster.MarkConsumed(ctx, nil, broadcast); err != nil {
				l.Errorw("Failed to mark log consumed", "err", err, "log", broadcast.String())
			}
		}
	}()

	if pollReq != PollRequestTypeHibernation && fm.pollManager.isHibernating.Load() {
		l.Warnw("Skipping poll because a ticker fired while hibernating")
		return
	}

	if !fm.logBroadcaster.IsConnected() {
		l.Warnw("LogBroadcaster is not connected to Ethereum node, skipping poll")
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
		fm.jobORM.TryRecordError(ctx, fm.spec.JobID,
			"Unable to call roundState method on provided contract. Check contract address.",
		)

		return
	}

	l = l.With("reportableRound", roundState.RoundId)

	// Because drumbeat ticker may fire at the same time on multiple nodes, we wait a short random duration
	// after getting a recommended round id, to avoid starting multiple rounds in case of chains with instant tx confirmation
	if pollReq == PollRequestTypeDrumbeat && fm.pollManager.cfg.DrumbeatEnabled && fm.pollManager.cfg.DrumbeatRandomDelay > 0 {
		// #nosec
		delay := time.Duration(mrand.Int63n(int64(fm.pollManager.cfg.DrumbeatRandomDelay)))
		l.Infof("waiting %v (of max: %v) before continuing...", delay, fm.pollManager.cfg.DrumbeatRandomDelay)
		time.Sleep(delay)

		roundStateNew, err2 := fm.roundState(roundState.RoundId)
		if err2 != nil {
			l.Errorw("unable to determine eligibility to submit from FluxAggregator contract", "err", err2)
			fm.jobORM.TryRecordError(ctx, fm.spec.JobID,
				"Unable to call roundState method on provided contract. Check contract address.",
			)

			return
		}
		roundState = roundStateNew
	}

	fm.pollManager.Reset(roundState)
	// Retry if a idle timer fails
	defer func() {
		if pollReq == PollRequestTypeIdle {
			if err != nil {
				if fm.pollManager.StartRetryTicker() {
					min, max := fm.pollManager.retryTicker.Bounds()
					l.Debugw(fmt.Sprintf("started retry ticker (frequency between: %v - %v) because of error: '%v'", min, max, err.Error()))
				}
				return
			}
			fm.pollManager.StopRetryTicker()
		}
	}()

	roundStats, jobRunStatus, err := fm.statsAndStatusForRound(ctx, roundState.RoundId, 0)
	if err != nil {
		l.Errorw("error determining round stats / run status for round", "err", err)

		return
	}

	// If we've already successfully submitted to this round (ie through a NewRound log)
	// and the associated JobRun hasn't errored, skip polling
	if roundStats.NumSubmissions > 0 && !jobRunStatus.Errored() {
		l.Infow("skipping poll: round already answered, tx unconfirmed", "jobRunStatus", jobRunStatus)

		return
	}

	// Don't submit if we're not eligible, or won't get paid
	err = fm.checkEligibilityAndAggregatorFunding(roundState)
	if err != nil {
		l.Infof("skipping poll: %v", err)

		return
	}

	var metaDataForBridge map[string]interface{}
	lrd, err := fm.fluxAggregator.LatestRoundData(nil)
	if err != nil {
		l.Warnw("Couldn't read latest round data for request meta", "err", err)
	} else {
		metaDataForBridge, err = bridges.MarshalBridgeMetaData(lrd.Answer, lrd.UpdatedAt)
		if err != nil {
			l.Warnw("Error marshalling roundState for request meta", "err", err)
		}
	}

	// Call the v2 pipeline to execute a new pipeline run
	// Note: we expect the FM pipeline to scale the fetched answer by the same
	// amount as "decimals" in the FM contract.

	vars := pipeline.NewVarsFrom(map[string]interface{}{
		"jobSpec": map[string]interface{}{
			"databaseID":    fm.jobSpec.ID,
			"externalJobID": fm.jobSpec.ExternalJobID,
			"name":          fm.jobSpec.Name.ValueOrZero(),
			"evmChainID":    fm.chainID.String(),
		},
		"jobRun": map[string]interface{}{
			"meta": metaDataForBridge,
		},
	})

	run, results, err := fm.runner.ExecuteRun(ctx, fm.spec, vars)
	if err != nil {
		l.Errorw("can't fetch answer", "err", err)
		fm.jobORM.TryRecordError(ctx, fm.spec.JobID, "Error polling")
		return
	}
	result, err := results.FinalResult().SingularResult()
	if err != nil || result.Error != nil {
		l.Errorw("can't fetch answer", "err", err, "result", result)
		fm.jobORM.TryRecordError(ctx, fm.spec.JobID, "Error polling")
		return
	}
	answer, err := utils.ToDecimal(result.Value)
	if err != nil {
		l.Errorw(fmt.Sprintf("error executing new run for job ID %v name %v", fm.spec.JobID, fm.spec.JobName), "err", err)
		return
	}

	if !fm.isValidSubmission(ctx, l, answer, started) {
		return
	}

	jobID := strconv.Itoa(int(fm.spec.JobID))
	latestAnswer := decimal.NewFromBigInt(roundState.LatestSubmission, 0)
	promfm.SetDecimal(promfm.SeenValue.WithLabelValues(jobID), answer)

	l = l.With(
		"latestAnswer", latestAnswer,
		"answer", answer,
	)

	if roundState.RoundId > 1 && !deviationChecker.OutsideDeviation(latestAnswer, answer) {
		l.Debugw("deviation < threshold, not submitting")
		return
	}

	if roundState.RoundId > 1 {
		l.Infow("deviation > threshold, submitting")
	} else {
		l.Infow("starting first round")
	}

	if roundState.PaymentAmount == nil {
		l.Error("roundState.PaymentAmount shouldn't be nil")
	}

	err = fm.Transact(ctx, func(tx sqlutil.DataSource) error {
		if err2 := fm.runner.InsertFinishedRun(ctx, tx, run, true); err2 != nil {
			return err2
		}
		if err2 := fm.queueTransactionForTxm(ctx, tx, run.ID, answer, roundState.RoundId, nil); err2 != nil {
			return err2
		}
		if broadcast != nil {
			// In the case of a flag lowered, the pollEligible call is triggered by a log.
			return fm.logBroadcaster.MarkConsumed(ctx, tx, broadcast)
		}
		return nil
	})
	// Either the tx failed and we want to reprocess the log, or it succeeded and already marked it consumed
	markConsumed = false
	if err != nil {
		l.Errorw("can't create job run", "err", err)
		return
	}

	promfm.SetDecimal(promfm.ReportedValue.WithLabelValues(jobID), answer)
	promfm.SetUint32(promfm.ReportedRound.WithLabelValues(jobID), roundState.RoundId)
}

// If the answer is outside the allowable range, log an error and don't submit.
// to avoid an onchain reversion.
func (fm *FluxMonitor) isValidSubmission(ctx context.Context, l logger.Logger, answer decimal.Decimal, started time.Time) bool {
	if fm.submissionChecker.IsValid(answer) {
		return true
	}

	l.Errorw("answer is outside acceptable range",
		"min", fm.submissionChecker.Min,
		"max", fm.submissionChecker.Max,
		"answer", answer,
	)
	fm.jobORM.TryRecordError(ctx, fm.spec.JobID, "Answer is outside acceptable range")

	jobId := fm.spec.JobID
	jobName := fm.spec.JobName
	elapsed := time.Since(started)
	pipeline.PromPipelineTaskExecutionTime.WithLabelValues(strconv.Itoa(int(jobId)), jobName, "", job.FluxMonitor.String()).Set(float64(elapsed))
	pipeline.PromPipelineRunErrors.WithLabelValues(strconv.Itoa(int(jobId)), jobName).Inc()
	pipeline.PromPipelineRunTotalTimeToCompletion.WithLabelValues(strconv.Itoa(int(jobId)), jobName).Set(float64(elapsed))
	pipeline.PromPipelineTasksTotalFinished.WithLabelValues(strconv.Itoa(int(jobId)), jobName, "", job.FluxMonitor.String(), "", "error").Inc()
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

func (fm *FluxMonitor) queueTransactionForTxm(ctx context.Context, tx sqlutil.DataSource, runID int64, answer decimal.Decimal, roundID uint32, log *flux_aggregator_wrapper.FluxAggregatorNewRound) error {
	// Use pipeline run ID to generate globally unique key that can correlate this run to a Tx
	idempotencyKey := fmt.Sprintf("fluxmonitor-%d", runID)
	// Submit the Eth Tx
	err := fm.contractSubmitter.Submit(
		ctx,
		new(big.Int).SetInt64(int64(roundID)),
		answer.BigInt(),
		&idempotencyKey,
	)
	if err != nil {
		fm.logger.Errorw("failed to submit Tx to TXM", "err", err)
		return err
	}

	numLogs := uint(0)
	if log != nil {
		numLogs = 1
	}
	// Update the flux monitor round stats
	err = fm.orm.WithDataSource(tx).UpdateFluxMonitorRoundStats(
		ctx,
		fm.contractAddress,
		roundID,
		runID,
		numLogs,
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

func (fm *FluxMonitor) statsAndStatusForRound(ctx context.Context, roundID uint32, newRoundLogs uint) (FluxMonitorRoundStatsV2, pipeline.RunStatus, error) {
	roundStats, err := fm.orm.FindOrCreateFluxMonitorRoundStats(ctx, fm.contractAddress, roundID, newRoundLogs)
	if err != nil {
		return FluxMonitorRoundStatsV2{}, pipeline.RunStatusUnknown, err
	}

	// JobRun will not exist if this is the first time responding to this round
	var run pipeline.Run
	if roundStats.PipelineRunID.Valid {
		run, err = fm.pipelineORM.FindRun(ctx, roundStats.PipelineRunID.Int64)
		if err != nil {
			return FluxMonitorRoundStatsV2{}, pipeline.RunStatusUnknown, err
		}
	}

	return roundStats, run.Status(), nil
}
