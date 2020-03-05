package fluxmonitor

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"sync"
	"time"

	"chainlink/core/eth"
	"chainlink/core/eth/contracts"
	"chainlink/core/logger"
	"chainlink/core/store"
	"chainlink/core/store/models"
	"chainlink/core/store/orm"
	"chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

//go:generate mockery -name Service -output ../../internal/mocks/ -case=underscore
//go:generate mockery -name DeviationCheckerFactory -output ../../internal/mocks/ -case=underscore
//go:generate mockery -name DeviationChecker -output ../../internal/mocks/ -case=underscore

// defaultHTTPTimeout is the timeout used by the price adapter fetcher for outgoing HTTP requests.
const defaultHTTPTimeout = 5 * time.Second

// MinimumPollingInterval is the smallest possible polling interval the Flux
// Monitor supports.
const MinimumPollingInterval = models.Duration(defaultHTTPTimeout)

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
	checkerFactory DeviationCheckerFactory
	adds           chan addEntry
	removes        chan models.ID
	connect        chan *models.Head
	disconnect     chan struct{}
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
	return &concreteFluxMonitor{
		store:      store,
		runManager: runManager,
		checkerFactory: pollingDeviationCheckerFactory{
			store:          store,
			logBroadcaster: eth.NewLogBroadcaster(store.TxManager),
		},
		adds:       make(chan addEntry),
		removes:    make(chan models.ID),
		connect:    make(chan *models.Head),
		disconnect: make(chan struct{}),
		chStop:     make(chan struct{}),
		chDone:     make(chan struct{}),
	}
}

func (fm *concreteFluxMonitor) Start() error {
	go fm.processAddRemoveJobRequests()

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
	close(fm.chStop)
	<-fm.chDone
}

// actionConsumer is the CSP consumer. It's run on a single goroutine to
// coordinate the collection of DeviationCheckers in a thread-safe fashion.
func (fm *concreteFluxMonitor) processAddRemoveJobRequests() {
	defer close(fm.chDone)

	jobMap := map[models.ID][]DeviationChecker{}

	for {
		select {
		case entry := <-fm.adds:
			if _, ok := jobMap[entry.jobID]; ok {
				logger.Errorf("job %s has already been added to flux monitor", entry.jobID)
				return
			}
			for _, checker := range entry.checkers {
				checker.Start()
			}
			jobMap[entry.jobID] = entry.checkers

		case jobID := <-fm.removes:
			for _, checker := range jobMap[jobID] {
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
		checker, err := fm.checkerFactory.New(initr, fm.runManager, fm.store.ORM)
		if err != nil {
			return errors.Wrap(err, "factory unable to create checker")
		}
		validCheckers = append(validCheckers, checker)
	}
	if len(validCheckers) == 0 {
		return nil
	}

	fm.adds <- addEntry{*job.ID, validCheckers}
	return nil
}

// RemoveJob stops and removes the checker for all Flux Monitor initiators belonging
// to the passed job ID.
func (fm *concreteFluxMonitor) RemoveJob(id *models.ID) {
	if id == nil {
		logger.Warn("nil job ID passed to FluxMonitor#RemoveJob")
		return
	}
	fm.removes <- *id
}

// DeviationCheckerFactory holds the New method needed to create a new instance
// of a DeviationChecker.
type DeviationCheckerFactory interface {
	New(models.Initiator, RunManager, *orm.ORM) (DeviationChecker, error)
}

type pollingDeviationCheckerFactory struct {
	store          *store.Store
	logBroadcaster eth.LogBroadcaster
}

func (f pollingDeviationCheckerFactory) New(initr models.Initiator, runManager RunManager, orm *orm.ORM) (DeviationChecker, error) {
	if initr.InitiatorParams.PollingInterval < MinimumPollingInterval {
		return nil, fmt.Errorf("pollingInterval must be equal or greater than %s", MinimumPollingInterval)
	}

	urls, err := ExtractFeedURLs(initr.InitiatorParams.Feeds, orm)
	if err != nil {
		return nil, err
	}

	fetcher, err := newMedianFetcherFromURLs(
		defaultHTTPTimeout,
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
	pollDelay     time.Duration
	idleThreshold time.Duration

	connected         utils.AtomicBool
	chMaybeLogs       chan maybeLog
	reportableRoundID *big.Int

	chStop     chan struct{}
	waitOnStop chan struct{}
}

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
		store:          store,
		fluxAggregator: fluxAggregator,
		initr:          initr,
		requestData:    initr.InitiatorParams.RequestData,
		idleThreshold:  initr.InitiatorParams.IdleThreshold.Duration(),
		threshold:      float64(initr.InitiatorParams.Threshold),
		precision:      initr.InitiatorParams.Precision,
		runManager:     runManager,
		fetcher:        fetcher,
		pollDelay:      pollDelay,
		chMaybeLogs:    make(chan maybeLog),
		chStop:         make(chan struct{}),
		waitOnStop:     make(chan struct{}),
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
	p.connected.Set(true)
}

func (p *PollingDeviationChecker) OnDisconnect() {
	logger.Debugw("PollingDeviationChecker disconnected from Ethereum node",
		"address", p.initr.InitiatorParams.Address.Hex(),
	)
	p.connected.Set(false)
}

func (p *PollingDeviationChecker) HandleLog(log interface{}, err error) {
	select {
	case p.chMaybeLogs <- maybeLog{log, err}:
	case <-p.chStop:
	}
}

func (p *PollingDeviationChecker) consume() {
	defer close(p.waitOnStop)

	unsubscribeLogs := p.fluxAggregator.SubscribeToLogs(p)
	defer unsubscribeLogs()

	p.pollIfEligible(p.threshold)

	var idleTimeout <-chan time.Time
	if p.idleThreshold > 0 {
		idleTimeout = time.After(p.idleThreshold)
	}

	pollTimeout := time.NewTicker(p.pollDelay)
	defer pollTimeout.Stop()

Loop:
	for {
		select {
		case <-p.chStop:
			return
		case maybeLog := <-p.chMaybeLogs:
			if maybeLog.Err != nil {
				logger.Errorf("error received from log broadcaster: %v", maybeLog.Err)
				continue Loop
			}
			roundFinished, err := p.respondToLog(maybeLog.Log)
			logger.ErrorIf(err, fmt.Sprintf("checker unable to respond to %T log:", maybeLog.Log))

			// The idleThreshold resets after a finished round
			if roundFinished && p.idleThreshold > 0 {
				idleTimeout = time.After(p.idleThreshold)
			}

		case <-pollTimeout.C:
			p.pollIfEligible(p.threshold)
		case <-idleTimeout:
			p.pollIfEligible(0)
		}
	}
}

func (p *PollingDeviationChecker) respondToLog(log interface{}) (roundFinished bool, _ error) {
	switch log := log.(type) {
	case *contracts.LogNewRound:
		logger.Debugw("NewRound log",
			"round", log.RoundId,
			"startedBy", log.StartedBy.Hex(),
			"startedAt", log.StartedAt.String(),
			"contract", p.initr.InitiatorParams.Address.Hex(),
			"job", p.initr.JobSpecID,
		)
		return false, p.respondToNewRoundLog(log)
	case *contracts.LogAnswerUpdated:
		logger.Debugw("AnswerUpdated log",
			"round", log.RoundId,
			"current", log.Current.String(),
			"timestamp", log.Timestamp.String(),
			"contract", p.initr.InitiatorParams.Address.Hex(),
			"job", p.initr.JobSpecID,
		)
		return p.respondToAnswerUpdatedLog(log)
	default:
		return false, nil
	}
}

// The AnswerUpdated log tell us when a round has closed with an answer, either by timing out
// or because enough submissions have been received.  We use this to reset the idleThreshold timer.
//
// Only invoked by the CSP consumer on the single goroutine for thread safety.
func (p *PollingDeviationChecker) respondToAnswerUpdatedLog(log *contracts.LogAnswerUpdated) (roundFinished bool, _ error) {
	if p.reportableRoundID != nil && log.RoundId.Cmp(p.reportableRoundID) < 0 {
		return false, nil
	}
	p.reportableRoundID = log.RoundId
	return true, nil
}

// The NewRound log tells us that an oracle has initiated a new round.  This tells us that we
// need to poll and submit an answer to the contract regardless of the deviation.
//
// Only invoked by the CSP consumer on the single goroutine for thread safety.
func (p *PollingDeviationChecker) respondToNewRoundLog(log *contracts.LogNewRound) error {

	logKeysAndValues := []interface{}{
		"newRound", log.RoundId,
		"reportableRoundID", p.reportableRoundID,
		"address", log.Address.Hex(),
		"jobID", p.initr.JobSpecID,
	}

	// Ignore old rounds.
	if p.reportableRoundID != nil && log.RoundId.Cmp(p.reportableRoundID) <= 0 {
		logger.Infow("Ignoring new round request: new <= current", logKeysAndValues...)
		return nil
	}

	jobSpecID := p.initr.JobSpecID.String()
	promSetBigInt(promFMSeenRound.WithLabelValues(jobSpecID), log.RoundId)

	// It's possible for RoundState() to return a higher round ID than the one in the NewRound log
	// (for example, if a large set of logs are delayed and arrive all at once).  We trust the value
	// from RoundState() over the one in the log, and record it as the current ReportableRoundID.
	roundState, err := p.fluxAggregator.RoundState()
	if err != nil {
		logger.Infow(fmt.Sprintf("Ignoring new round request: error fetching eligibility from contract: %v", err), logKeysAndValues...)
		return err
	}
	p.reportableRoundID = roundState.ReportableRoundID

	if !roundState.EligibleToSubmit {
		logger.Infow("Ignoring new round request: not eligible to submit", logKeysAndValues...)
		return nil
	}

	logger.Infow("Responding to new round request: new > current", logKeysAndValues...)

	polledAnswer, err := p.fetcher.Fetch()
	if err != nil {
		return errors.Wrap(err, "unable to fetch median price")
	}

	return p.createJobRun(polledAnswer, roundState.ReportableRoundID)
}

// poll walks through the steps to check for a deviation, early exiting if deviation
// is not met, or triggering a new job run if deviation is met.
// Only invoked by the CSP consumer on the single goroutine for thread safety.
func (p *PollingDeviationChecker) pollIfEligible(threshold float64) {
	if p.connected.Get() == false {
		logger.Errorf("not connected to Ethereum node, skipping poll")
		return
	}

	roundState, err := p.fluxAggregator.RoundState()
	if err != nil {
		logger.Errorf("unable to determine eligibility to submit from FluxAggregator contract: %v", err)
		return
	}

	// It's pointless to listen to logs from before the current reporting round
	p.reportableRoundID = roundState.ReportableRoundID

	if !roundState.EligibleToSubmit {
		logger.Info("not eligible to submit, skipping poll")
		return
	}

	polledAnswer, err := p.fetcher.Fetch()
	if err != nil {
		logger.Errorf("can't fetch answer: %v", err)
		return
	}

	jobSpecID := p.initr.JobSpecID.String()
	promSetDecimal(promFMSeenValue.WithLabelValues(jobSpecID), polledAnswer)

	latestAnswer := decimal.NewFromBigInt(roundState.LatestAnswer, -p.precision)
	if !OutsideDeviation(latestAnswer, polledAnswer, threshold) {
		logger.Debugw("deviation < threshold, not submitting",
			"latestAnswer", latestAnswer,
			"polledAnswer", polledAnswer,
			"threshold", threshold,
		)
		return
	}

	logger.Infow("deviation > threshold, starting new round",
		"reportableRound", roundState.ReportableRoundID,
		"address", p.initr.Address.Hex(),
		"jobID", p.initr.JobSpecID,
	)
	err = p.createJobRun(polledAnswer, roundState.ReportableRoundID)
	if err != nil {
		logger.Errorf("can't create job run: %v", err)
		return
	}

	promSetDecimal(promFMReportedValue.WithLabelValues(jobSpecID), polledAnswer)
	promSetBigInt(promFMReportedRound.WithLabelValues(jobSpecID), roundState.ReportableRoundID)
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
	payload := fmt.Sprintf(`{
			"result": "%s",
			"address": "%s",
			"functionSelector": "%s",
			"dataPrefix": "%s"
	}`,
		polledAnswer,
		p.initr.InitiatorParams.Address.Hex(),
		hexutil.Encode(methodID),
		hexutil.Encode(nextRoundData))

	runData, err := models.ParseJSON([]byte(payload))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("unable to start chainlink run with payload %s", payload))
	}
	runRequest := models.NewRunRequest(runData)

	_, err = p.runManager.Create(p.initr.JobSpecID, &p.initr, nil, runRequest)
	return err
}

var dec0 = decimal.NewFromInt(0)

// OutsideDeviation checks whether the next price is outside the threshold.
func OutsideDeviation(curAnswer, nextAnswer decimal.Decimal, threshold float64) bool {
	if curAnswer.Equal(dec0) {
		logger.Infow("Current price is 0, deviation automatically met", "answer", dec0)
		return true
	}

	diff := curAnswer.Sub(nextAnswer).Abs()
	percentage := diff.Div(curAnswer).Mul(decimal.NewFromInt(100))
	if percentage.LessThan(decimal.NewFromFloat(threshold)) {
		logger.Debugw(
			"Deviation threshold not met",
			"difference", percentage,
			"threshold", threshold,
			"currentAnswer", curAnswer,
			"nextAnswer", nextAnswer)
		return false
	}
	logger.Infow(
		"Deviation threshold met",
		"difference", percentage,
		"threshold", threshold,
		"currentAnswer", curAnswer,
		"nextAnswer", nextAnswer,
	)
	return true
}
