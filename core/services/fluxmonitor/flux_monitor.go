package fluxmonitor

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"reflect"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/eth/contracts"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/tevino/abool"
)

//go:generate mockery -name Service -output ../../internal/mocks/ -case=underscore
//go:generate mockery -name DeviationCheckerFactory -output ../../internal/mocks/ -case=underscore
//go:generate mockery -name DeviationChecker -output ../../internal/mocks/ -case=underscore

type RunManager interface {
	Create(
		jobSpecID *models.ID,
		initiator *models.Initiator,
		creationHeight *big.Int,
		runRequest *models.RunRequest,
	) (*models.JobRun, error)
}

// Service is the interface encapsulating all functionality
// needed to listen to price deviations and new round requests.
type Service interface {
	AddJob(models.JobSpec) error
	RemoveJob(*models.ID)
	Start() error
	Stop()
}

type concreteFluxMonitor struct {
	store          *store.Store
	runManager     RunManager
	logBroadcaster eth.LogBroadcaster
	checkerFactory DeviationCheckerFactory
	chAdd          chan addEntry
	chRemove       chan models.ID
	chConnect      chan *models.Head
	chDisconnect   chan struct{}
	chStop         chan struct{}
	chDone         chan struct{}
	disabled       bool
	started        bool
}

type addEntry struct {
	jobID    models.ID
	checkers []DeviationChecker
}

// New creates a service that manages a collection of DeviationCheckers,
// one per initiator of type InitiatorFluxMonitor for added jobs.
func New(
	store *store.Store,
	runManager RunManager,
) Service {
	if store.Config.EthereumDisabled() {
		return &concreteFluxMonitor{disabled: true}
	}

	logBroadcaster := eth.NewLogBroadcaster(store)
	return &concreteFluxMonitor{
		store:          store,
		runManager:     runManager,
		logBroadcaster: logBroadcaster,
		checkerFactory: pollingDeviationCheckerFactory{
			store:          store,
			logBroadcaster: logBroadcaster,
		},
		chAdd:        make(chan addEntry),
		chRemove:     make(chan models.ID),
		chConnect:    make(chan *models.Head),
		chDisconnect: make(chan struct{}),
		chStop:       make(chan struct{}),
		chDone:       make(chan struct{}),
	}
}

func (fm *concreteFluxMonitor) Start() error {
	if fm.disabled {
		logger.Info("Flux monitor disabled: skipping start")
		return nil
	}

	go fm.serveInternalRequests()
	fm.started = true

	var wg sync.WaitGroup
	err := fm.store.Jobs(func(j *models.JobSpec) bool {
		if j == nil {
			err := errors.New("received nil job")
			logger.Error(err)
			return true
		}
		job := *j

		wg.Add(1)
		go func() {
			defer wg.Done()

			err := fm.AddJob(job)
			if err != nil {
				logger.Errorf("error adding FluxMonitor job: %v", err)
			}
		}()
		return true
	}, models.InitiatorFluxMonitor)

	wg.Wait()
	fm.logBroadcaster.Start()

	return err
}

// Disconnect cleans up running deviation checkers.
func (fm *concreteFluxMonitor) Stop() {
	if fm.disabled {
		logger.Info("Flux monitor disabled: cannot stop")
		return
	}

	fm.logBroadcaster.Stop()
	close(fm.chStop)
	if fm.started {
		fm.started = false
		<-fm.chDone
	}
}

// serveInternalRequests handles internal requests for state change via
// channels.  Inspired by the ideas of Communicating Sequential Processes, or
// CSP.
func (fm *concreteFluxMonitor) serveInternalRequests() {
	defer close(fm.chDone)

	jobMap := map[models.ID][]DeviationChecker{}

	for {
		select {
		case entry := <-fm.chAdd:
			if _, ok := jobMap[entry.jobID]; ok {
				logger.Errorf("job '%s' has already been added to flux monitor", entry.jobID.String())
				continue
			}
			for _, checker := range entry.checkers {
				checker.Start()
			}
			jobMap[entry.jobID] = entry.checkers

		case jobID := <-fm.chRemove:
			checkers, ok := jobMap[jobID]
			if !ok {
				logger.Debugf("job '%s' is missing from the flux monitor", jobID.String())
				continue
			}
			for _, checker := range checkers {
				checker.Stop()
			}
			delete(jobMap, jobID)

		case <-fm.chStop:
			for _, checkers := range jobMap {
				for _, checker := range checkers {
					checker.Stop()
				}
			}
			return
		}
	}
}

// AddJob created a DeviationChecker for any job initiators of type
// InitiatorFluxMonitor.
func (fm *concreteFluxMonitor) AddJob(job models.JobSpec) error {
	if job.ID == nil {
		err := errors.New("received job with nil ID")
		logger.Error(err)
		return err
	}

	var validCheckers []DeviationChecker
	for _, initr := range job.InitiatorsFor(models.InitiatorFluxMonitor) {
		logger.Debugw("Adding job to flux monitor",
			"job", job.ID.String(),
			"initr", initr.ID,
		)

		timeout := fm.store.Config.DefaultHTTPTimeout()
		checker, err := fm.checkerFactory.New(
			initr,
			job.MinPayment,
			fm.runManager,
			fm.store.ORM,
			timeout,
		)
		if err != nil {
			return errors.Wrap(err, "factory unable to create checker")
		}
		validCheckers = append(validCheckers, checker)
	}
	if len(validCheckers) == 0 {
		return nil
	}

	fm.chAdd <- addEntry{*job.ID, validCheckers}
	return nil
}

// RemoveJob stops and removes the checker for all Flux Monitor initiators belonging
// to the passed job ID.
func (fm *concreteFluxMonitor) RemoveJob(id *models.ID) {
	if id == nil {
		logger.Warn("nil job ID passed to FluxMonitor#RemoveJob")
		return
	}
	fm.chRemove <- *id
}

// DeviationCheckerFactory holds the New method needed to create a new instance
// of a DeviationChecker.
type DeviationCheckerFactory interface {
	New(models.Initiator, *assets.Link, RunManager, *orm.ORM, models.Duration) (DeviationChecker, error)
}

type pollingDeviationCheckerFactory struct {
	store          *store.Store
	logBroadcaster eth.LogBroadcaster
}

func (f pollingDeviationCheckerFactory) New(
	initr models.Initiator,
	minJobPayment *assets.Link,
	runManager RunManager,
	orm *orm.ORM,
	timeout models.Duration,
) (DeviationChecker, error) {
	minimumPollingInterval := models.Duration(f.store.Config.DefaultHTTPTimeout())

	if !initr.PollTimer.Disabled &&
		initr.PollTimer.Period.Shorter(minimumPollingInterval) {
		return nil, fmt.Errorf("pollTimer.period must be equal or greater than %s", minimumPollingInterval)
	}

	urls, err := ExtractFeedURLs(initr.Feeds, orm)
	if err != nil {
		return nil, err
	}

	fetcher, err := newMedianFetcherFromURLs(
		timeout,
		initr.RequestData.String(),
		urls)
	if err != nil {
		return nil, err
	}

	f.logBroadcaster.AddDependents(1)
	fluxAggregator, err := contracts.NewFluxAggregator(initr.Address, f.store.TxManager, f.logBroadcaster)
	if err != nil {
		return nil, err
	}

	return NewPollingDeviationChecker(
		f.store,
		fluxAggregator,
		initr,
		minJobPayment,
		runManager,
		fetcher,
		func() { f.logBroadcaster.DependentReady() },
	)
}

// ExtractFeedURLs extracts a list of url.URLs from the feeds parameter of the initiator params
func ExtractFeedURLs(feeds models.Feeds, orm *orm.ORM) ([]*url.URL, error) {
	var feedsData []interface{}
	var urls []*url.URL

	err := json.Unmarshal(feeds.Bytes(), &feedsData)
	if err != nil {
		return nil, err
	}

	for _, entry := range feedsData {
		var bridgeURL *url.URL
		var err error

		switch feed := entry.(type) {
		case string: // feed url - ex: "http://example.com"
			bridgeURL, err = url.ParseRequestURI(feed)
		case map[string]interface{}: // named feed - ex: {"bridge": "bridgeName"}
			bridgeName := feed["bridge"].(string)
			bridgeURL, err = GetBridgeURLFromName(bridgeName, orm) // XXX: currently an n query
		default:
			err = errors.New("unable to extract feed URLs from json")
		}

		if err != nil {
			return nil, err
		}
		urls = append(urls, bridgeURL)
	}

	return urls, nil
}

// GetBridgeURLFromName looks up a bridge in the DB by name, then extracts the url
func GetBridgeURLFromName(name string, orm *orm.ORM) (*url.URL, error) {
	task := models.TaskType(name)
	bridge, err := orm.FindBridge(task)
	if err != nil {
		return nil, err
	}
	bridgeURL := url.URL(bridge.URL)
	return &bridgeURL, nil
}

// DeviationChecker encapsulate methods needed to initialize and check prices
// for price deviations.
type DeviationChecker interface {
	Start()
	Stop()
}

// PollingDeviationChecker polls external price adapters via HTTP to check for price swings.
type PollingDeviationChecker struct {
	store          *store.Store
	fluxAggregator contracts.FluxAggregator
	runManager     RunManager
	fetcher        Fetcher

	initr         models.Initiator
	minJobPayment *assets.Link
	requestData   models.JSON
	precision     int32

	connected     *abool.AtomicBool
	backlog       *utils.BoundedPriorityQueue
	chProcessLogs chan struct{}
	pollTicker    <-chan time.Time
	idleTimer     <-chan time.Time
	roundTimer    <-chan time.Time

	readyForLogs func()
	chStop       chan struct{}
	waitOnStop   chan struct{}
}

// NewPollingDeviationChecker returns a new instance of PollingDeviationChecker.
func NewPollingDeviationChecker(
	store *store.Store,
	fluxAggregator contracts.FluxAggregator,
	initr models.Initiator,
	minJobPayment *assets.Link,
	runManager RunManager,
	fetcher Fetcher,
	readyForLogs func(),
) (*PollingDeviationChecker, error) {
	return &PollingDeviationChecker{
		readyForLogs:   readyForLogs,
		store:          store,
		fluxAggregator: fluxAggregator,
		initr:          initr,
		minJobPayment:  minJobPayment,
		requestData:    initr.RequestData,
		precision:      initr.Precision,
		runManager:     runManager,
		fetcher:        fetcher,
		pollTicker:     nil,
		idleTimer:      nil,
		roundTimer:     nil,
		connected:      abool.New(),
		backlog: utils.NewBoundedPriorityQueue(map[uint]uint{
			// We want reconnecting nodes to be able to submit to a round
			// that hasn't hit maxAnswers yet, as well as the newest round.
			PriorityNewRoundLog:      2,
			PriorityAnswerUpdatedLog: 1,
		}),
		chProcessLogs: make(chan struct{}, 1),
		chStop:        make(chan struct{}),
		waitOnStop:    make(chan struct{}),
	}, nil
}

const (
	PriorityNewRoundLog      uint = 0
	PriorityAnswerUpdatedLog uint = 1
)

// Start begins the CSP consumer in a single goroutine to
// poll the price adapters and listen to NewRound events.
func (p *PollingDeviationChecker) Start() {
	logger.Debugw("Starting checker for job",
		"job", p.initr.JobSpecID.String(),
		"initr", p.initr.ID,
	)

	go p.consume()
}

// Stop stops this instance from polling, cleaning up resources.
func (p *PollingDeviationChecker) Stop() {
	close(p.chStop)
	<-p.waitOnStop
}

func (p *PollingDeviationChecker) OnConnect() {
	logger.Debugw("PollingDeviationChecker connected to Ethereum node",
		"jobID", p.initr.JobSpecID.String(),
		"address", p.initr.Address.Hex(),
	)
	p.connected.Set()
}

func (p *PollingDeviationChecker) OnDisconnect() {
	logger.Debugw("PollingDeviationChecker disconnected from Ethereum node",
		"jobID", p.initr.JobSpecID.String(),
		"address", p.initr.Address.Hex(),
	)
	p.connected.UnSet()
}

func (p *PollingDeviationChecker) HandleLog(broadcast eth.LogBroadcast, err error) {
	if err != nil {
		logger.Errorf("got error from LogBroadcaster: %v", err)
		return
	}

	rawLog := broadcast.Log()
	if rawLog == nil || reflect.ValueOf(rawLog).IsNil() {
		logger.Error("HandleLog: ignoring nil value")
		return
	}

	switch log := rawLog.(type) {
	case *contracts.LogNewRound:
		p.backlog.Add(PriorityNewRoundLog, broadcast)

	case *contracts.LogAnswerUpdated:
		p.backlog.Add(PriorityAnswerUpdatedLog, broadcast)

	default:
		logger.Warnf("unexpected log type %T", log)
		return
	}

	select {
	case p.chProcessLogs <- struct{}{}:
	default:
	}
}

func (p *PollingDeviationChecker) consume() {
	defer close(p.waitOnStop)

	connected, unsubscribeLogs := p.fluxAggregator.SubscribeToLogs(p)
	defer unsubscribeLogs()

	if connected {
		p.connected.Set()
	} else {
		p.connected.UnSet()
	}

	p.readyForLogs()

	if !p.initr.PollTimer.Disabled {
		// Try to do an initial poll
		p.pollIfEligible(DeviationThresholds{
			Rel: float64(p.initr.Threshold),
			Abs: float64(p.initr.AbsoluteThreshold),
		})

		ticker := time.NewTicker(p.initr.PollTimer.Period.Duration())
		defer ticker.Stop()
		p.pollTicker = ticker.C
	}
	if !p.initr.IdleTimer.Disabled {
		p.idleTimer = time.After(p.initr.IdleTimer.Duration.Duration())
	}

	for {
		select {
		case <-p.chStop:
			return

		case <-p.chProcessLogs:
			p.processLogs()

		case <-p.pollTicker:
			logger.Debugw("Poll ticker fired",
				"pollPeriod", p.initr.PollTimer.Period,
				"idleDuration", p.initr.IdleTimer.Duration,
				"contract", p.initr.Address.Hex(),
			)
			p.pollIfEligible(DeviationThresholds{
				Rel: float64(p.initr.Threshold),
				Abs: float64(p.initr.AbsoluteThreshold),
			})

		case <-p.idleTimer:
			logger.Debugw("Idle ticker fired",
				"pollPeriod", p.initr.PollTimer.Period,
				"idleDuration", p.initr.IdleTimer.Duration,
				"contract", p.initr.Address.Hex(),
			)
			p.pollIfEligible(DeviationThresholds{Rel: 0, Abs: 0})

		case <-p.roundTimer:
			logger.Debugw("Round timeout ticker fired",
				"pollPeriod", p.initr.PollTimer.Period,
				"idleDuration", p.initr.IdleTimer.Duration,
				"contract", p.initr.Address.Hex(),
			)
			p.pollIfEligible(DeviationThresholds{
				Rel: float64(p.initr.Threshold),
				Abs: float64(p.initr.AbsoluteThreshold),
			})
		}
	}
}

func (p *PollingDeviationChecker) processLogs() {
	for !p.backlog.Empty() {
		broadcast := p.backlog.Take().(eth.LogBroadcast)

		// If the log is a duplicate of one we've seen before, ignore it (this
		// happens because of the LogBroadcaster's backfilling behavior).
		consumed, err := broadcast.WasAlreadyConsumed()
		if err != nil {
			logger.Errorf("Error determining if log was already consumed: %v", err)
			continue
		} else if consumed {
			continue
		}

		switch log := broadcast.Log().(type) {
		case *contracts.LogNewRound:
			p.respondToNewRoundLog(*log)

			err := broadcast.MarkConsumed()
			if err != nil {
				logger.Errorf("Error marking log as consumed: %v", err)
			}

		case *contracts.LogAnswerUpdated:
			p.respondToAnswerUpdatedLog(*log)

			err := broadcast.MarkConsumed()
			if err != nil {
				logger.Errorf("Error marking log as consumed: %v", err)
			}

		default:
			logger.Errorf("unknown log %v of type %T", log, log)
		}
	}
}

// The AnswerUpdated log tells us that round has successfully closed with a new
// answer.  We update our view of the oracleRoundState in case this log was
// generated by a chain reorg.
func (p *PollingDeviationChecker) respondToAnswerUpdatedLog(log contracts.LogAnswerUpdated) {
	logger.Debugw("AnswerUpdated log", p.loggerFieldsForAnswerUpdated(log)...)

	_, err := p.roundState(0)
	if err != nil {
		logger.Errorw(fmt.Sprintf("could not fetch oracleRoundState: %v", err), p.loggerFieldsForAnswerUpdated(log)...)
	}
}

// The NewRound log tells us that an oracle has initiated a new round.  This tells us that we
// need to poll and submit an answer to the contract regardless of the deviation.
func (p *PollingDeviationChecker) respondToNewRoundLog(log contracts.LogNewRound) {
	logger.Debugw("NewRound log", p.loggerFieldsForNewRound(log)...)

	promSetBigInt(promFMSeenRound.WithLabelValues(p.initr.JobSpecID.String()), log.RoundId)

	//
	// NewRound answer submission logic:
	//   - Any log that reaches this point, regardless of chain reorgs or log backfilling, is one that we have
	//         not seen before.  Therefore, we should act upon it.
	//   - We want to ensure that we respond to a previous, supersedable round in addition to the most recent
	//         round (sometimes logs come in quickly, all at once, as with a reorg or backfilling).  Therefore,
	//         we must use the roundID on the log, rather than the `.ReportableRoundID` value suggested by the
	//         FluxAggregator's `oracleRoundState` method.
	//   - We query the oracleRoundState of this specific roundID to ensure that it makes sense to submit.  We
	//         ignore its `reportableRoundID` suggestion, because we already have a specific roundID to consider.
	//   - In the event of a reorg that pushes our previous submissions back into the mempool, we can rely on the
	//         TxManager to ensure they end up being mined into blocks, but this may cause them to revert if they
	//         violate certain conditions in the FluxAggregator (restartDelay, etc.).  Therefore, the cleanest
	//         solution at present is to resubmit for the reorged rounds.  The drawback of this approach is that
	//         one or the other submission tx for a given round will revert, costing the node operator some gas.
	//         The benefit is that those submissions are guaranteed to be made, ensuring that we have high data
	//         availability (and also ensuring that node operators get paid).
	//

	// Ignore rounds we started
	acct, err := p.store.KeyStore.GetFirstAccount()
	if err != nil {
		logger.Errorw(fmt.Sprintf("error fetching account from keystore: %v", err), p.loggerFieldsForNewRound(log)...)
		return
	} else if log.StartedBy == acct.Address {
		logger.Infow("Ignoring new round request: we started this round", p.loggerFieldsForNewRound(log)...)
		return
	}

	// Ignore rounds we're not eligible for, or for which we won't be paid
	roundState, err := p.roundState(uint32(log.RoundId.Uint64()))
	if err != nil {
		logger.Errorw(fmt.Sprintf("Ignoring new round request: error fetching eligibility from contract: %v", err), p.loggerFieldsForNewRound(log)...)
		return
	}
	err = p.checkEligibilityAndAggregatorFunding(roundState)
	if err != nil {
		logger.Infow(fmt.Sprintf("Ignoring new round request: %v", err), p.loggerFieldsForNewRound(log)...)
		return
	}

	logger.Infow("Responding to new round request: new > current", p.loggerFieldsForNewRound(log)...)

	polledAnswer, err := p.fetcher.Fetch()
	if err != nil {
		logger.Errorw(fmt.Sprintf("unable to fetch median price: %v", err), p.loggerFieldsForNewRound(log)...)
		return
	}

	err = p.createJobRun(polledAnswer, uint32(log.RoundId.Uint64()))
	if err != nil {
		logger.Errorw(fmt.Sprintf("unable to create job run: %v", err), p.loggerFieldsForNewRound(log)...)
		return
	}
}

var (
	ErrNotEligible   = errors.New("not eligible to submit")
	ErrUnderfunded   = errors.New("aggregator is underfunded")
	ErrPaymentTooLow = errors.New("round payment amount < minimum contract payment")
)

func (p *PollingDeviationChecker) checkEligibilityAndAggregatorFunding(roundState contracts.FluxAggregatorRoundState) error {
	if !roundState.EligibleToSubmit {
		return ErrNotEligible
	} else if !p.sufficientFunds(roundState) {
		return ErrUnderfunded
	} else if !p.sufficientPayment(roundState.PaymentAmount) {
		return ErrPaymentTooLow
	}
	return nil
}

const MinFundedRounds int64 = 3

// sufficientFunds checks if the contract has sufficient funding to pay all the oracles on a
// conract for a minimum number of rounds, based on the payment amount in the contract
func (p *PollingDeviationChecker) sufficientFunds(state contracts.FluxAggregatorRoundState) bool {
	min := big.NewInt(int64(state.OracleCount))
	min = min.Mul(min, big.NewInt(MinFundedRounds))
	min = min.Mul(min, state.PaymentAmount)
	return state.AvailableFunds.Cmp(min) >= 0
}

// sufficientPayment checks if the available payment is enough to submit an answer. It compares
// the payment amount on chain with the min payment amount listed in the job spec / ENV var.
func (p *PollingDeviationChecker) sufficientPayment(payment *big.Int) bool {
	aboveOrEqMinGlobalPayment :=
		payment.Cmp(p.store.Config.MinimumContractPayment().ToInt()) >= 0
	aboveOrEqMinJobPayment := true
	if p.minJobPayment != nil {
		aboveOrEqMinJobPayment =
			payment.Cmp(p.minJobPayment.ToInt()) >= 0
	}
	return aboveOrEqMinGlobalPayment && aboveOrEqMinJobPayment
}

// DeviationThresholds carries parameters used by the threshold-trigger logic
type DeviationThresholds struct {
	Rel float64 // Relative change required, i.e. |new-old|/|new| >= Rel
	Abs float64 // Absolute change required, i.e. |new-old| >= Abs
}

func (p *PollingDeviationChecker) pollIfEligible(thresholds DeviationThresholds) (createdJobRun bool) {
	loggerFields := []interface{}{
		"jobID", p.initr.JobSpecID,
		"address", p.initr.InitiatorParams.Address,
		"threshold", thresholds.Rel,
		"absoluteThreshold", thresholds.Abs,
	}

	if !p.connected.IsSet() {
		logger.Warnw("not connected to Ethereum node, skipping poll", loggerFields...)
		return false
	}

	//
	// Poll ticker submission logic:
	//   - We avoid saving on-chain state wherever possible.  Therefore, we do not know which round we should be
	//         submitting for when the pollTicker fires.
	//   - We pass 0 into `roundState()`, and the FluxAggregator returns a suggested roundID for us to
	//         submit to, as well as our eligibility to submit to that round.
	//

	// Ask the FluxAggregator which round we should be submitting to, and what the state of that round is.
	roundState, err := p.roundState(0)
	if err != nil {
		logger.Errorw(fmt.Sprintf("unable to determine eligibility to submit from FluxAggregator contract: %v", err), loggerFields...)
		return false
	}
	loggerFields = append(loggerFields, "reportableRound", roundState.ReportableRoundID)

	// Don't submit if we're not eligible, or won't get paid
	err = p.checkEligibilityAndAggregatorFunding(roundState)
	if err != nil {
		logger.Infow(fmt.Sprintf("skipping poll: %v", err), loggerFields...)
		return false
	}

	polledAnswer, err := p.fetcher.Fetch()
	if err != nil {
		logger.Errorw(fmt.Sprintf("can't fetch answer: %v", err), loggerFields...)
		return false
	}

	jobSpecID := p.initr.JobSpecID.String()
	latestAnswer := decimal.NewFromBigInt(roundState.LatestAnswer, -p.precision)

	promSetDecimal(promFMSeenValue.WithLabelValues(jobSpecID), polledAnswer)
	loggerFields = append(loggerFields,
		"latestAnswer", latestAnswer,
		"polledAnswer", polledAnswer,
	)
	if roundState.ReportableRoundID > 1 && !OutsideDeviation(latestAnswer, polledAnswer, thresholds) {
		logger.Debugw("deviation < threshold, not submitting", loggerFields...)
		return false
	}

	if roundState.ReportableRoundID > 1 {
		logger.Infow("deviation > threshold, starting new round", loggerFields...)
	} else {
		logger.Infow("starting first round", loggerFields...)
	}

	err = p.createJobRun(polledAnswer, roundState.ReportableRoundID)
	if err != nil {
		logger.Errorw(fmt.Sprintf("can't create job run: %v", err), loggerFields...)
		return false
	}

	promSetDecimal(promFMReportedValue.WithLabelValues(jobSpecID), polledAnswer)
	promSetUint32(promFMReportedRound.WithLabelValues(jobSpecID), roundState.ReportableRoundID)
	return true
}

func (p *PollingDeviationChecker) roundState(roundID uint32) (contracts.FluxAggregatorRoundState, error) {
	acct, err := p.store.KeyStore.GetFirstAccount()
	if err != nil {
		return contracts.FluxAggregatorRoundState{}, err
	}
	roundState, err := p.fluxAggregator.RoundState(acct.Address, roundID)
	if err != nil {
		return contracts.FluxAggregatorRoundState{}, err
	}

	// Update our tickers to reflect the current on-chain round
	p.resetRoundTimeoutTicker(roundState)
	p.resetIdleTicker(roundState.StartedAt)

	return roundState, nil
}

func (p *PollingDeviationChecker) resetRoundTimeoutTicker(roundState contracts.FluxAggregatorRoundState) {
	loggerFields := p.loggerFields("timesOutAt", roundState.TimesOutAt())

	if roundState.TimesOutAt() == 0 {
		p.roundTimer = nil
		logger.Debugw("disabling roundTimer, no active round", loggerFields...)

	} else {
		timesOutAt := time.Unix(int64(roundState.TimesOutAt()), 0)
		timeUntilTimeout := time.Until(timesOutAt)

		if timeUntilTimeout <= 0 {
			p.roundTimer = nil
			logger.Debugw("roundTimer has run down; disabling", loggerFields...)
		} else {
			p.roundTimer = time.After(timeUntilTimeout)
			loggerFields = append(loggerFields, "value", roundState.TimesOutAt())
			logger.Debugw("updating roundState.TimesOutAt", loggerFields...)
		}
	}
}

func (p *PollingDeviationChecker) resetIdleTicker(roundStartedAtUTC uint64) {
	if p.initr.IdleTimer.Disabled {
		p.idleTimer = nil
		return
	} else if roundStartedAtUTC == 0 {
		// There is no active round, so keep using the idleTimer we already have
		return
	}

	startedAt := time.Unix(int64(roundStartedAtUTC), 0)
	idleDeadline := startedAt.Add(p.initr.IdleTimer.Duration.Duration())
	timeUntilIdleDeadline := time.Until(idleDeadline)
	loggerFields := p.loggerFields(
		"startedAt", roundStartedAtUTC,
		"timeUntilIdleDeadline", timeUntilIdleDeadline,
	)

	if timeUntilIdleDeadline <= 0 {
		logger.Debugw("not resetting idleTimer, negative duration", loggerFields...)
		return
	}
	p.idleTimer = time.After(timeUntilIdleDeadline)
	logger.Debugw("resetting idleTimer", loggerFields...)
}

// jobRunRequest is the request used to trigger a Job Run by the Flux Monitor.
type jobRunRequest struct {
	Result           decimal.Decimal `json:"result"`
	Address          string          `json:"address"`
	FunctionSelector string          `json:"functionSelector"`
	DataPrefix       string          `json:"dataPrefix"`
}

func (p *PollingDeviationChecker) createJobRun(polledAnswer decimal.Decimal, roundID uint32) error {
	methodID, err := p.fluxAggregator.GetMethodID("submit")
	if err != nil {
		return err
	}

	roundIDData := utils.EVMWordUint64(uint64(roundID))

	payload, err := json.Marshal(jobRunRequest{
		Result:           polledAnswer,
		Address:          p.initr.Address.Hex(),
		FunctionSelector: hexutil.Encode(methodID),
		DataPrefix:       hexutil.Encode(roundIDData),
	})
	if err != nil {
		return errors.Wrapf(err, "unable to encode Job Run request in JSON")
	}
	runData, err := models.ParseJSON(payload)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("unable to start chainlink run with payload %s", payload))
	}
	runRequest := models.NewRunRequest(runData)

	_, err = p.runManager.Create(p.initr.JobSpecID, &p.initr, nil, runRequest)
	if err != nil {
		return err
	}

	return nil
}

func (p *PollingDeviationChecker) loggerFields(added ...interface{}) []interface{} {
	return append(added, []interface{}{
		"pollFrequency", p.initr.PollTimer.Period,
		"idleDuration", p.initr.IdleTimer.Duration,
		"contract", p.initr.Address.Hex(),
		"jobID", p.initr.JobSpecID.String(),
	}...)
}

func (p *PollingDeviationChecker) loggerFieldsForNewRound(log contracts.LogNewRound) []interface{} {
	return []interface{}{
		"round", log.RoundId,
		"startedBy", log.StartedBy.Hex(),
		"startedAt", log.StartedAt.String(),
		"contract", log.Address.Hex(),
		"jobID", p.initr.JobSpecID,
	}
}

func (p *PollingDeviationChecker) loggerFieldsForAnswerUpdated(log contracts.LogAnswerUpdated) []interface{} {
	return []interface{}{
		"round", log.RoundId,
		"answer", log.Current.String(),
		"timestamp", log.Timestamp.String(),
		"contract", log.Address.Hex(),
		"job", p.initr.JobSpecID,
	}
}

func (p *PollingDeviationChecker) JobID() *models.ID {
	return p.initr.JobSpecID
}

// OutsideDeviation checks whether the next price is outside the threshold.
// If both thresholds are zero (default value), always returns true.
func OutsideDeviation(curAnswer, nextAnswer decimal.Decimal, thresholds DeviationThresholds) bool {
	loggerFields := []interface{}{
		"threshold", thresholds.Rel,
		"absoluteThreshold", thresholds.Abs,
		"currentAnswer", curAnswer,
		"nextAnswer", nextAnswer,
	}

	if thresholds.Rel == 0 && thresholds.Abs == 0 {
		logger.Debugw(
			"Deviation thresholds both zero; short-circuiting deviation checker to "+
				"true, regardless of feed values", loggerFields...)
		return true
	}
	diff := curAnswer.Sub(nextAnswer).Abs()
	loggerFields = append(loggerFields, "absoluteDeviation", diff)

	if !diff.GreaterThan(decimal.NewFromFloat(thresholds.Abs)) {
		logger.Debugw("Absolute deviation threshold not met", loggerFields...)
		return false
	}

	if curAnswer.IsZero() {
		if nextAnswer.IsZero() {
			logger.Debugw("Relative deviation is undefined; can't satisfy threshold", loggerFields...)
			return false
		}
		logger.Infow("Threshold met: relative deviation is ∞", loggerFields...)
		return true
	}

	// 100*|new-old|/|new|: Deviation (relative to curAnswer) as a percentage
	percentage := diff.Div(curAnswer.Abs()).Mul(decimal.NewFromInt(100))

	loggerFields = append(loggerFields, "percentage", percentage)

	if percentage.LessThan(decimal.NewFromFloat(thresholds.Rel)) {
		logger.Debugw("Relative deviation threshold not met", loggerFields...)
		return false
	}
	logger.Infow("Relative and absolute deviation thresholds both met", loggerFields...)
	return true
}

// MakeIdleTimer checks the log timestamp and calculates the idle time
// from that.
//
// This function makes the assumption that the local system time is
// relatively accurate (to within a second or so) and all participating nodes
// agree on that.
//
// If system time is not accurate (compared to the cluster) then you should
// expect poor behaviour here.
func MakeIdleTimer(log contracts.LogNewRound, idleThreshold models.Duration, clock utils.AfterNower) <-chan time.Time {
	timeNow := clock.Now()
	if log.StartedAt == nil {
		return defaultIdleTimer(idleThreshold, clock)
	}
	if !log.StartedAt.IsInt64() {
		logger.Errorf("Value for log.StartedAt %s would overflow int64, using default idle timer instead.", log.StartedAt.String())
		return defaultIdleTimer(idleThreshold, clock)
	}
	roundStarted := time.Unix(log.StartedAt.Int64(), 0)
	if roundStarted.After(timeNow) {
		logger.Warnf("Round started time of %s is later than current system time of %s, setting idle timer to %s from now. Most likely scenario is that this machine's clock is running slow. This is suboptimal! Please ensure your system clock is accurate.", roundStarted.String(), timeNow.String(), idleThreshold.Duration().String())
		return defaultIdleTimer(idleThreshold, clock)
	}
	// duration from now until idle threshold = log timestamp + idle threshold - current time
	durationUntilIdleThreshold := roundStarted.Add(idleThreshold.Duration()).Sub(timeNow)
	if durationUntilIdleThreshold < 0 {
		logger.Warnf("Idle threshold already passed, current time is %s and idle timer expired at %s (round started at %s with idle threshold of %s). It's possible you are processing an old round, or this machine has a fast clock. If this keeps happening, check your system clock and make sure it is accurate.", timeNow, roundStarted.Add(idleThreshold.Duration()).String(), roundStarted.String(), idleThreshold.Duration().String())
	}
	return clock.After(durationUntilIdleThreshold)
}

func defaultIdleTimer(idleThreshold models.Duration, clock utils.AfterNower) <-chan time.Time {
	return clock.After(idleThreshold.Duration())
}
