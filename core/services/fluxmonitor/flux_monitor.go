package fluxmonitor

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/eth/contracts"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jinzhu/gorm"
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
	logBroadcaster := eth.NewLogBroadcaster(store.TxManager, store.ORM)
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
	if fm.store.Config.EthereumDisabled() {
		return nil
	}

	fm.logBroadcaster.Start()

	go fm.serveInternalRequests()

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

	return err
}

// Disconnect cleans up running deviation checkers.
func (fm *concreteFluxMonitor) Stop() {
	fm.logBroadcaster.Stop()
	close(fm.chStop)
	<-fm.chDone
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
				logger.Errorf("job '%s' has already been added to flux monitor", entry.jobID)
				return
			}
			for _, checker := range entry.checkers {
				checker.Start()
			}
			jobMap[entry.jobID] = entry.checkers

		case jobID := <-fm.chRemove:
			checkers, ok := jobMap[jobID]
			if !ok {
				logger.Errorf("job '%s' is missing from the flux monitor", jobID)
				return
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
		checker, err := fm.checkerFactory.New(initr, fm.runManager, fm.store.ORM, timeout)
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
	New(models.Initiator, RunManager, *orm.ORM, time.Duration) (DeviationChecker, error)
}

type pollingDeviationCheckerFactory struct {
	store          *store.Store
	logBroadcaster eth.LogBroadcaster
}

func (f pollingDeviationCheckerFactory) New(
	initr models.Initiator,
	runManager RunManager,
	orm *orm.ORM,
	timeout time.Duration,
) (DeviationChecker, error) {
	minimumPollingInterval := models.Duration(f.store.Config.DefaultHTTPTimeout())

	if initr.InitiatorParams.PollingInterval < minimumPollingInterval {
		return nil, fmt.Errorf("pollingInterval must be equal or greater than %s", minimumPollingInterval)
	}

	urls, err := ExtractFeedURLs(initr.InitiatorParams.Feeds, orm)
	if err != nil {
		return nil, err
	}

	fetcher, err := newMedianFetcherFromURLs(
		timeout,
		initr.InitiatorParams.RequestData.String(),
		urls)
	if err != nil {
		return nil, err
	}

	fluxAggregator, err := contracts.NewFluxAggregator(initr.InitiatorParams.Address, f.store.TxManager, f.logBroadcaster)
	if err != nil {
		return nil, err
	}

	return NewPollingDeviationChecker(
		f.store,
		fluxAggregator,
		initr,
		runManager,
		fetcher,
		initr.InitiatorParams.PollingInterval.Duration(),
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
	requestData   models.JSON
	threshold     float64
	precision     int32
	idleThreshold time.Duration

	connected                  *abool.AtomicBool
	chMaybeLogs                chan maybeLog
	reportableRoundID          *big.Int
	mostRecentSubmittedRoundID uint64
	pollTicker                 *ResettableTicker
	idleTicker                 <-chan time.Time
	roundTimeoutTicker         <-chan time.Time

	chStop     chan struct{}
	waitOnStop chan struct{}
}

// maybeLog is just a tuple that allows us to send either an error or a log over the
// logs channel.  This is preferable to using two separate channels, as it ensures
// that we don't drop valid (but unprocessed) logs if we receive an error.
type maybeLog struct {
	Log interface{}
	Err error
}

// NewPollingDeviationChecker returns a new instance of PollingDeviationChecker.
func NewPollingDeviationChecker(
	store *store.Store,
	fluxAggregator contracts.FluxAggregator,
	initr models.Initiator,
	runManager RunManager,
	fetcher Fetcher,
	pollDelay time.Duration,
) (*PollingDeviationChecker, error) {
	return &PollingDeviationChecker{
		store:              store,
		fluxAggregator:     fluxAggregator,
		initr:              initr,
		requestData:        initr.InitiatorParams.RequestData,
		idleThreshold:      initr.InitiatorParams.IdleThreshold.Duration(),
		threshold:          float64(initr.InitiatorParams.Threshold),
		precision:          initr.InitiatorParams.Precision,
		runManager:         runManager,
		fetcher:            fetcher,
		pollTicker:         NewResettableTicker(pollDelay),
		idleTicker:         nil,
		roundTimeoutTicker: nil,
		connected:          abool.New(),
		chMaybeLogs:        make(chan maybeLog, 100),
		chStop:             make(chan struct{}),
		waitOnStop:         make(chan struct{}),
	}, nil
}

// Start begins the CSP consumer in a single goroutine to
// poll the price adapters and listen to NewRound events.
func (p *PollingDeviationChecker) Start() {
	logger.Debugw("Starting checker for job",
		"job", p.initr.JobSpecID.String(),
		"initr", p.initr.ID)

	go p.consume()
}

// Stop stops this instance from polling, cleaning up resources.
func (p *PollingDeviationChecker) Stop() {
	close(p.chStop)
	<-p.waitOnStop
}

func (p *PollingDeviationChecker) OnConnect() {
	logger.Debugw("PollingDeviationChecker connected to Ethereum node",
		"address", p.initr.InitiatorParams.Address.Hex(),
	)
	p.connected.Set()
}

func (p *PollingDeviationChecker) OnDisconnect() {
	logger.Debugw("PollingDeviationChecker disconnected from Ethereum node",
		"address", p.initr.InitiatorParams.Address.Hex(),
	)
	p.connected.UnSet()
}

type ResettableTicker struct {
	*time.Ticker
	d time.Duration
}

func NewResettableTicker(d time.Duration) *ResettableTicker {
	return &ResettableTicker{nil, d}
}

func (t *ResettableTicker) Tick() <-chan time.Time {
	if t.Ticker == nil {
		return nil
	}
	return t.Ticker.C
}

func (t *ResettableTicker) Stop() {
	if t.Ticker != nil {
		t.Ticker.Stop()
		t.Ticker = nil
	}
}

func (t *ResettableTicker) Reset() {
	t.Stop()
	t.Ticker = time.NewTicker(t.d)
}

func (p *PollingDeviationChecker) HandleLog(log interface{}, err error) {
	select {
	case p.chMaybeLogs <- maybeLog{log, err}:
	case <-p.chStop:
	}
}

func (p *PollingDeviationChecker) consume() {
	defer close(p.waitOnStop)

	p.determineMostRecentSubmittedRoundID()

	connected, unsubscribeLogs := p.fluxAggregator.SubscribeToLogs(p)
	defer unsubscribeLogs()

	if connected {
		p.connected.Set()
	} else {
		p.connected.UnSet()
	}

	// Try to do an initial poll
	p.pollIfEligible(p.threshold)
	p.pollTicker.Reset()
	defer p.pollTicker.Stop()

	if p.idleThreshold > 0 {
		p.idleTicker = time.After(p.idleThreshold)
	}

	for {
		select {
		case <-p.chStop:
			return

		case maybeLog := <-p.chMaybeLogs:
			if maybeLog.Err != nil {
				logger.Errorf("error received from log broadcaster: %v", maybeLog.Err)
				continue
			}
			p.respondToLog(maybeLog.Log)

		case <-p.pollTicker.Tick():
			logger.Debugw("Poll ticker fired",
				"pollDelay", p.pollTicker.d,
				"idleThreshold", p.idleThreshold,
				"mostRecentSubmittedRoundID", p.mostRecentSubmittedRoundID,
				"reportableRoundID", p.reportableRoundID,
				"contract", p.initr.InitiatorParams.Address.Hex(),
			)
			p.pollIfEligible(p.threshold)

		case <-p.idleTicker:
			logger.Debugw("Idle ticker fired",
				"pollDelay", p.pollTicker.d,
				"idleThreshold", p.idleThreshold,
				"mostRecentSubmittedRoundID", p.mostRecentSubmittedRoundID,
				"reportableRoundID", p.reportableRoundID,
				"contract", p.initr.InitiatorParams.Address.Hex(),
			)
			p.pollIfEligible(0)

		case <-p.roundTimeoutTicker:
			logger.Debugw("Round timeout ticker fired",
				"pollDelay", p.pollTicker.d,
				"idleThreshold", p.idleThreshold,
				"mostRecentSubmittedRoundID", p.mostRecentSubmittedRoundID,
				"reportableRoundID", p.reportableRoundID,
				"contract", p.initr.InitiatorParams.Address.Hex(),
			)
			p.pollIfEligible(p.threshold)
		}
	}
}

func (p *PollingDeviationChecker) determineMostRecentSubmittedRoundID() {
	myAccount, err := p.store.KeyStore.GetFirstAccount()
	if err != nil {
		logger.Error("error determining most recent submitted round ID: ", err)
		return
	}

	// Just to be particularly defensive against issues with the DB or TxManager, we
	// fetch the most recent 5 transactions we've submitted to this aggregator from our
	// Chainlink node address.  Take the highest round ID among them and store it so
	// that we avoid re-polling for a given round when our tx takes a while to confirm.
	txs, err := p.store.ORM.FindTxsBySenderAndRecipient(myAccount.Address, p.initr.InitiatorParams.Address, 0, 5)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		logger.Error("error determining most recent submitted round ID: ", err)
		return
	}

	// Parse the round IDs from the transaction data
	for _, tx := range txs {
		if len(tx.Data) != 68 {
			logger.Warnw("found Flux Monitor tx with bad data payload",
				"txID", tx.ID,
			)
			continue
		}

		roundIDBytes := tx.Data[4:36]
		roundID := big.NewInt(0).SetBytes(roundIDBytes).Uint64()
		if roundID > p.mostRecentSubmittedRoundID {
			p.mostRecentSubmittedRoundID = roundID
		}
	}
	logger.Infow(fmt.Sprintf("roundID of most recent submission is %v", p.mostRecentSubmittedRoundID),
		"jobID", p.initr.JobSpecID,
		"aggregator", p.initr.InitiatorParams.Address.Hex(),
	)
}

func (p *PollingDeviationChecker) respondToLog(log interface{}) {
	switch log := log.(type) {
	case *contracts.LogNewRound:
		logger.Debugw("NewRound log", p.loggerFieldsForNewRound(log)...)
		p.respondToNewRoundLog(log)

	case *contracts.LogAnswerUpdated:
		logger.Debugw("AnswerUpdated log", p.loggerFieldsForAnswerUpdated(log)...)
		p.respondToAnswerUpdatedLog(log)

	default:
	}
}

// The AnswerUpdated log tells us that round has successfully close with a new
// answer.  This tells us that we need to reset our poll ticker.
//
// Only invoked by the CSP consumer on the single goroutine for thread safety.
func (p *PollingDeviationChecker) respondToAnswerUpdatedLog(log *contracts.LogAnswerUpdated) {
	if p.reportableRoundID != nil && log.RoundId.Cmp(p.reportableRoundID) < 0 {
		// Ignore old rounds
		logger.Debugw("Ignoring stale AnswerUpdated log", p.loggerFieldsForAnswerUpdated(log)...)
		return
	}
	p.pollTicker.Reset()
}

// The NewRound log tells us that an oracle has initiated a new round.  This tells us that we
// need to poll and submit an answer to the contract regardless of the deviation.
//
// Only invoked by the CSP consumer on the single goroutine for thread safety.
func (p *PollingDeviationChecker) respondToNewRoundLog(log *contracts.LogNewRound) {
	// The idleThreshold resets when a new round starts
	if p.idleThreshold > 0 {
		p.idleTicker = time.After(p.idleThreshold)
	}

	jobSpecID := p.initr.JobSpecID.String()
	promSetBigInt(promFMSeenRound.WithLabelValues(jobSpecID), log.RoundId)

	// Ignore rounds we started
	acct, err := p.store.KeyStore.GetFirstAccount()
	if err != nil {
		logger.Errorw(fmt.Sprintf("error fetching account from keystore: %v", err), p.loggerFieldsForNewRound(log)...)
		return
	} else if log.StartedBy == acct.Address {
		logger.Infow("Ignoring new round request: we started this round", p.loggerFieldsForNewRound(log)...)
		return
	}

	// It's possible for RoundState() to return a higher round ID than the one in the NewRound log
	// (for example, if a large set of logs are delayed and arrive all at once).  We trust the value
	// from RoundState() over the one in the log, and record it as the current ReportableRoundID.
	roundState, err := p.roundState()
	if err != nil {
		logger.Errorw(fmt.Sprintf("Ignoring new round request: error fetching eligibility from contract: %v", err), p.loggerFieldsForNewRound(log)...)
		return
	}

	err = p.checkEligibilityAndAggregatorFunding(roundState)
	if errors.Cause(err) == ErrAlreadySubmitted {
		logger.Infow(fmt.Sprintf("Ignoring new round request: %v, possible chain reorg", err), p.loggerFieldsForNewRound(log)...)
		return
	} else if err != nil {
		logger.Infow(fmt.Sprintf("Ignoring new round request: %v", err), p.loggerFieldsForNewRound(log)...)
		return
	}

	// Ignore old rounds
	if log.RoundId.Cmp(p.reportableRoundID) < 0 {
		logger.Infow("Ignoring new round request: new < current", p.loggerFieldsForNewRound(log)...)
		return
	} else if log.RoundId.Uint64() <= p.mostRecentSubmittedRoundID {
		logger.Infow("Ignoring new round request: already submitted for this round", p.loggerFieldsForNewRound(log)...)
		return
	}

	logger.Infow("Responding to new round request: new > current", p.loggerFieldsForNewRound(log)...)

	polledAnswer, err := p.fetcher.Fetch()
	if err != nil {
		logger.Errorw(fmt.Sprintf("unable to fetch median price: %v", err), p.loggerFieldsForNewRound(log)...)
		return
	}

	p.createJobRun(polledAnswer, p.reportableRoundID)
}

var (
	ErrNotEligible      = errors.New("not eligible to submit")
	ErrUnderfunded      = errors.New("aggregator is underfunded")
	ErrPaymentTooLow    = errors.New("round payment amount < minimum contract payment")
	ErrAlreadySubmitted = errors.Errorf("already submitted for round")
)

func (p *PollingDeviationChecker) checkEligibilityAndAggregatorFunding(roundState contracts.FluxAggregatorRoundState) error {
	if !roundState.EligibleToSubmit {
		return ErrNotEligible
	} else if roundState.AvailableFunds.Cmp(roundState.PaymentAmount) < 0 {
		return ErrUnderfunded
	} else if roundState.PaymentAmount.Cmp(p.store.Config.MinimumContractPayment().ToInt()) < 0 {
		return ErrPaymentTooLow
	} else if p.mostRecentSubmittedRoundID >= uint64(roundState.ReportableRoundID) {
		return ErrAlreadySubmitted
	}
	return nil
}

func (p *PollingDeviationChecker) pollIfEligible(threshold float64) (createdJobRun bool) {
	loggerFields := []interface{}{
		"jobID", p.initr.JobSpecID,
		"address", p.initr.InitiatorParams.Address,
		"threshold", threshold,
	}

	if p.connected.IsSet() == false {
		logger.Warnw("not connected to Ethereum node, skipping poll", loggerFields...)
		return false
	}

	roundState, err := p.roundState()
	if err != nil {
		logger.Errorw(fmt.Sprintf("unable to determine eligibility to submit from FluxAggregator contract: %v", err), loggerFields...)
		return false
	}
	loggerFields = append(loggerFields, "reportableRound", roundState.ReportableRoundID)

	err = p.checkEligibilityAndAggregatorFunding(roundState)
	if errors.Cause(err) == ErrAlreadySubmitted {
		logger.Infow(fmt.Sprintf("skipping poll: %v, tx is pending", err), loggerFields...)
		return false
	} else if err != nil {
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
	if roundState.ReportableRoundID > 1 && !OutsideDeviation(latestAnswer, polledAnswer, threshold) {
		logger.Debugw("deviation < threshold, not submitting", loggerFields...)
		return false
	}

	if roundState.ReportableRoundID > 1 {
		logger.Infow("deviation > threshold, starting new round", loggerFields...)
	} else {
		logger.Infow("starting first round", loggerFields...)
	}

	err = p.createJobRun(polledAnswer, p.reportableRoundID)
	if err != nil {
		logger.Errorw(fmt.Sprintf("can't create job run: %v", err), loggerFields...)
		return false
	}

	promSetDecimal(promFMReportedValue.WithLabelValues(jobSpecID), polledAnswer)
	promSetBigInt(promFMReportedRound.WithLabelValues(jobSpecID), p.reportableRoundID)
	return true
}

func (p *PollingDeviationChecker) roundState() (contracts.FluxAggregatorRoundState, error) {
	acct, err := p.store.KeyStore.GetFirstAccount()
	if err != nil {
		return contracts.FluxAggregatorRoundState{}, err
	}
	roundState, err := p.fluxAggregator.RoundState(acct.Address)
	if err != nil {
		return contracts.FluxAggregatorRoundState{}, err
	}

	// It's pointless to listen to logs from before the current reporting round
	p.reportableRoundID = big.NewInt(int64(roundState.ReportableRoundID))

	// Update the roundTimeoutTicker using the .TimesOutAt field describing the current round
	if roundState.TimesOutAt == 0 {
		logger.Debugw("updating roundState.TimesOutAt",
			"value", roundState.TimesOutAt,
			"pollDelay", p.pollTicker.d,
			"idleThreshold", p.idleThreshold,
			"mostRecentSubmittedRoundID", p.mostRecentSubmittedRoundID,
			"reportableRoundID", p.reportableRoundID,
			"contract", p.initr.InitiatorParams.Address.Hex(),
		)
		p.roundTimeoutTicker = nil
	} else {
		timeUntilTimeout := time.Unix(int64(roundState.TimesOutAt), 0).Sub(time.Now())
		if timeUntilTimeout.Seconds() <= 0 {
			p.roundTimeoutTicker = nil
			logger.Debugw("NOT updating roundState.TimesOutAt, negative duration",
				"value", roundState.TimesOutAt,
				"pollDelay", p.pollTicker.d,
				"idleThreshold", p.idleThreshold,
				"mostRecentSubmittedRoundID", p.mostRecentSubmittedRoundID,
				"reportableRoundID", p.reportableRoundID,
				"contract", p.initr.InitiatorParams.Address.Hex(),
			)
		} else {
			p.roundTimeoutTicker = time.After(timeUntilTimeout)
			logger.Debugw("updating roundState.TimesOutAt",
				"value", roundState.TimesOutAt,
				"timeUntilTimeout", timeUntilTimeout,
				"pollDelay", p.pollTicker.d,
				"idleThreshold", p.idleThreshold,
				"mostRecentSubmittedRoundID", p.mostRecentSubmittedRoundID,
				"reportableRoundID", p.reportableRoundID,
				"contract", p.initr.InitiatorParams.Address.Hex(),
			)
		}
	}

	return roundState, nil
}

// jobRunRequest is the request used to trigger a Job Run by the Flux Monitor.
type jobRunRequest struct {
	Result           decimal.Decimal `json:"result"`
	Address          string          `json:"address"`
	FunctionSelector string          `json:"functionSelector"`
	DataPrefix       string          `json:"dataPrefix"`
}

func (p *PollingDeviationChecker) createJobRun(polledAnswer decimal.Decimal, nextRound *big.Int) error {
	methodID, err := p.fluxAggregator.GetMethodID("updateAnswer")
	if err != nil {
		return err
	}

	nextRoundData, err := utils.EVMWordBigInt(nextRound)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(jobRunRequest{
		Result:           polledAnswer,
		Address:          p.initr.InitiatorParams.Address.Hex(),
		FunctionSelector: hexutil.Encode(methodID),
		DataPrefix:       hexutil.Encode(nextRoundData),
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

	p.mostRecentSubmittedRoundID = nextRound.Uint64()

	return nil
}

func (p *PollingDeviationChecker) loggerFieldsForNewRound(log *contracts.LogNewRound) []interface{} {
	return []interface{}{
		"reportableRound", p.reportableRoundID,
		"round", log.RoundId,
		"startedBy", log.StartedBy.Hex(),
		"startedAt", log.StartedAt.String(),
		"contract", log.Address.Hex(),
		"jobID", p.initr.JobSpecID,
	}
}

func (p *PollingDeviationChecker) loggerFieldsForAnswerUpdated(log *contracts.LogAnswerUpdated) []interface{} {
	return []interface{}{
		"round", log.RoundId,
		"answer", log.Current.String(),
		"timestamp", log.Timestamp.String(),
		"contract", log.Address.Hex(),
		"job", p.initr.JobSpecID,
	}
}

// OutsideDeviation checks whether the next price is outside the threshold.
func OutsideDeviation(curAnswer, nextAnswer decimal.Decimal, threshold float64) bool {
	loggerFields := []interface{}{
		"threshold", threshold,
		"currentAnswer", curAnswer,
		"nextAnswer", nextAnswer,
	}

	if curAnswer.IsZero() {
		if nextAnswer.IsZero() {
			logger.Debugw("Deviation threshold not met", loggerFields...)
			return false
		}

		logger.Infow("Deviation threshold met", loggerFields...)
		return true
	}

	diff := curAnswer.Sub(nextAnswer).Abs()
	percentage := diff.Div(curAnswer.Abs()).Mul(decimal.NewFromInt(100))

	loggerFields = append(loggerFields, "percentage", percentage)

	if percentage.LessThan(decimal.NewFromFloat(threshold)) {
		logger.Debugw("Deviation threshold not met", loggerFields...)
		return false
	}
	logger.Infow("Deviation threshold met", loggerFields...)
	return true
}
