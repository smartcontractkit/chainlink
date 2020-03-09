package fluxmonitor

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"net/url"
	"time"

	"chainlink/core/eth"
	"chainlink/core/eth/contracts"
	"chainlink/core/logger"
	"chainlink/core/store"
	"chainlink/core/store/models"
	"chainlink/core/store/orm"
	"chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"go.uber.org/multierr"
)

//go:generate mockery -name FluxMonitor -output ../../internal/mocks/ -case=underscore
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
	store.HeadTrackable // (Dis)Connect methods handle initial boot and intermittent connectivity.
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
	removes        chan *models.ID
	connect        chan *models.Head
	disconnect     chan struct{}
	ctx            context.Context
	cancel         context.CancelFunc
}

type addEntry struct {
	jobID    string
	checkers []DeviationChecker
	errChan  chan error
}

// New creates a service that manages a collection of DeviationCheckers,
// one per initiator of type InitiatorFluxMonitor for added jobs.
func New(
	store *store.Store,
	runManager RunManager,
) Service {
	return &concreteFluxMonitor{
		store:          store,
		runManager:     runManager,
		checkerFactory: pollingDeviationCheckerFactory{store: store},
	}
}

func (fm *concreteFluxMonitor) Start() error {
	fm.ctx, fm.cancel = context.WithCancel(context.Background())
	fm.adds = make(chan addEntry)
	fm.removes = make(chan *models.ID)
	fm.connect = make(chan *models.Head)
	fm.disconnect = make(chan struct{})

	go fm.actionConsumer(fm.ctx)

	count := 0
	errChan := make(chan error)
	err := fm.store.Jobs(func(j *models.JobSpec) bool {
		go func(j *models.JobSpec) {
			errChan <- fm.AddJob(*j)
		}(j)
		count++
		return true
	}, models.InitiatorFluxMonitor)

	var merr error
	for i := 0; i < count; i++ {
		err := <-errChan
		merr = multierr.Combine(merr, err)
	}
	return multierr.Append(err, merr)
}

// Connect initializes all DeviationCheckers and starts their listening.
func (fm *concreteFluxMonitor) Connect(head *models.Head) error {
	fm.connect <- head
	return nil
}

// actionConsumer is the CSP consumer. It's run on a single goroutine to
// coordinate the collection of DeviationCheckers in a thread-safe fashion.
func (fm *concreteFluxMonitor) actionConsumer(ctx context.Context) {
	jobMap := map[string][]DeviationChecker{}

	// init w a noop cancel, so we never have to deal with nils
	connectionCtx, cancelConnection := context.WithCancel(ctx)
	var connected bool

	for {
		select {
		case <-ctx.Done():
			cancelConnection()
			return
		case <-fm.connect:
			// every connection, create a new ctx for canceling on disconnect.
			connectionCtx, cancelConnection = context.WithCancel(ctx)
			connectCheckers(connectionCtx, jobMap, fm.store.TxManager)
			connected = true
		case <-fm.disconnect:
			cancelConnection()
			connected = false
		case entry := <-fm.adds:
			entry.errChan <- fm.addAction(
				ctx,
				connected,
				jobMap,
				fm.store,
				entry.jobID,
				entry.checkers,
			)
		case jobID := <-fm.removes:
			for _, checker := range jobMap[jobID.String()] {
				checker.Stop()
			}
			delete(jobMap, jobID.String())
		}
	}
}

// Disconnect cleans up running deviation checkers.
func (fm *concreteFluxMonitor) Disconnect() {
	fm.disconnect <- struct{}{}
}

// Disconnect cleans up running deviation checkers.
func (fm *concreteFluxMonitor) Stop() {
	if fm.cancel != nil {
		fm.cancel()
	}
}

// OnNewHead is a noop.
func (fm *concreteFluxMonitor) OnNewHead(*models.Head) {}

// AddJob created a DeviationChecker for any job initiators of type
// InitiatorFluxMonitor.
func (fm *concreteFluxMonitor) AddJob(job models.JobSpec) error {
	validCheckers := []DeviationChecker{}
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

	errChan := make(chan error)
	fm.adds <- addEntry{job.ID.String(), validCheckers, errChan}
	return <-errChan
}

func connectCheckers(ctx context.Context, jobMap map[string][]DeviationChecker, client eth.Client) {
	for _, checkers := range jobMap {
		for _, checker := range checkers {
			// XXX: Add mechanism to asynchronously communicate when a job spec has
			// an ethereum interaction error.
			// https://www.pivotaltracker.com/story/show/170349568
			logger.ErrorIf(connectSingleChecker(ctx, checker, client))
		}
	}
}

func (fm *concreteFluxMonitor) addAction(
	ctx context.Context,
	connected bool,
	jobMap map[string][]DeviationChecker,
	store *store.Store,
	jobSpecID string,
	checkers []DeviationChecker,
) error {
	if _, ok := jobMap[jobSpecID]; ok {
		return fmt.Errorf(
			"job %s has already been added to flux monitor",
			jobSpecID,
		)
	}

	if connected {
		for _, checker := range checkers {
			err := connectSingleChecker(ctx, checker, fm.store.TxManager)
			if err != nil {
				return errors.Wrap(err, "unable to connect checker")
			}
		}
	}

	if len(checkers) > 0 {
		jobMap[jobSpecID] = checkers
	}
	return nil
}

func connectSingleChecker(ctx context.Context, checker DeviationChecker, client eth.Client) error {
	return checker.Start(ctx, client)
}

// RemoveJob stops and removes the checker for all Flux Monitor initiators belonging
// to the passed job ID.
func (fm *concreteFluxMonitor) RemoveJob(ID *models.ID) {
	fm.removes <- ID
}

// DeviationCheckerFactory holds the New method needed to create a new instance
// of a DeviationChecker.
type DeviationCheckerFactory interface {
	New(models.Initiator, RunManager, *orm.ORM) (DeviationChecker, error)
}

type pollingDeviationCheckerFactory struct {
	store *store.Store
}

func (f pollingDeviationCheckerFactory) New(initr models.Initiator, runManager RunManager, orm *orm.ORM) (DeviationChecker, error) {
	if initr.InitiatorParams.PollingInterval < MinimumPollingInterval {
		return nil, fmt.Errorf(
			"pollingInterval must be equal or greater than %s",
			MinimumPollingInterval,
		)
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

	return NewPollingDeviationChecker(
		f.store,
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
	Start(context.Context, eth.Client) error
	Stop()
}

// PollingDeviationChecker polls external price adapters via HTTP to check for price swings.
type PollingDeviationChecker struct {
	store          *store.Store
	fluxAggregator *contracts.FluxAggregator
	initr          models.Initiator
	address        common.Address
	requestData    models.JSON
	idleThreshold  time.Duration
	threshold      float64
	precision      int32
	runManager     RunManager
	currentPrice   decimal.Decimal
	currentRound   *big.Int
	fetcher        Fetcher
	delay          time.Duration
	cancel         context.CancelFunc
	chLogs         chan contracts.MaybeDecodedLog

	waitOnStop chan struct{}
}

// NewPollingDeviationChecker returns a new instance of PollingDeviationChecker.
func NewPollingDeviationChecker(
	store *store.Store,
	initr models.Initiator,
	runManager RunManager,
	fetcher Fetcher,
	delay time.Duration,
) (*PollingDeviationChecker, error) {
	return &PollingDeviationChecker{
		store:          store,
		fluxAggregator: fluxAggregator,
		initr:          initr,
		address:        initr.InitiatorParams.Address,
		requestData:    initr.InitiatorParams.RequestData,
		idleThreshold:  initr.InitiatorParams.IdleThreshold.Duration(),
		threshold:      float64(initr.InitiatorParams.Threshold),
		precision:      initr.InitiatorParams.Precision,
		runManager:     runManager,
		currentPrice:   decimal.NewFromInt(0),
		currentRound:   big.NewInt(0),
		fetcher:        fetcher,
		delay:          delay,
		newRounds:      make(chan eth.Log),

		waitOnStop: make(chan struct{}),
	}, nil
}

// Start begins the CSP consumer in a single goroutine to
// poll the price adapters and listen to NewRound events.
func (p *PollingDeviationChecker) Start(ctx context.Context, client eth.Client) error {
	logger.Debugw("Starting checker for job",
		"job", p.initr.JobSpecID.String(),
		"initr", p.initr.ID)

	fluxAggregator, err := contracts.NewFluxAggregator(ethClient, p.address)
	if err != nil {
		return nil, err
	}
	p.fluxAggregator = fluxAggregator

	err := p.fetchAggregatorData(client)
	if err != nil {
		return err
	}

	roundSubscription, err := p.fluxAggregator.SubscribeToLogs(nil)
	if err != nil {
		return err
	}

	_, err = p.poll(p.threshold)
	if err != nil {
		return err
	}

	ctx, p.cancel = context.WithCancel(ctx)
	go p.consume(ctx, roundSubscription, client)
	return nil
}

// Stop stops this instance from polling, cleaning up resources.
func (p *PollingDeviationChecker) Stop() {
	if p.cancel != nil {
		p.cancel()
		<-p.waitOnStop
	}
}

// stopTimer stops and clears the timer as suggested by the documentation.
func stopTimer(arg *time.Timer) {
	if !arg.Stop() && len(arg.C) > 0 {
		// Residual events are the timer's channel and need to be cleared.
		//
		// Refer to timer.Stop's documentation or
		// https://developpaper.com/detailed-explanation-of-the-trap-of-timer-in-golang/
		<-arg.C
	}
}

func (p *PollingDeviationChecker) consume(ctx context.Context, roundSubscription *contracts.LogSubscription, client eth.Client) {
	defer roundSubscription.Unsubscribe()

	idleThreshold := p.idleThreshold
	if idleThreshold == 0 {
		idleThreshold = math.MaxInt64
	}

	idleThresholdTimer := time.NewTimer(idleThreshold)
	defer stopTimer(idleThresholdTimer)

Loop:
	for {
		jobRunTriggered := false

		select {
		case <-ctx.Done():
			close(p.waitOnStop)
			return
		case maybeLog := <-roundSubscription.Logs():
			if maybeLog.Error != nil {
				logger.Error(errors.WithStack(maybeLog.Error))
				// @@TODO: other error handling?
				// @@TODO: detect broken subscription?
				continue Loop
			}
			_, _, err := p.respondToLog(maybeLog.Log)
			logger.ErrorIf(err, "checker unable to respond to new round")

		case <-time.After(p.delay):
			jobRunTriggered = p.pollIfEligible(client, p.threshold)
		case <-idleThresholdTimer.C:
			jobRunTriggered = p.pollIfEligible(client, 0)
		}

		if jobRunTriggered {
			// Reset expects stopped or expired timer.
			stopTimer(idleThresholdTimer)
			idleThresholdTimer.Reset(idleThreshold)
		}
	}
}

func (p *PollingDeviationChecker) respondToLog(log interface{}) (roundStarted bool, roundEnded bool, _ error) {
	switch log := log.(type) {
	case contracts.LogNewRound:
		return true, false, p.respondToNewRoundLog(log)
	case contracts.LogRoundDetailsUpdated:
		return false, false, p.respondToRoundDetailsUpdatedLog(log)
	case contracts.LogAnswerUpdated:
		return false, true, p.respondToAnswerUpdatedLog(log)
	default:
		panic("got unknown log")
	}
}

func (p *PollingDeviationChecker) respondToAnswerUpdatedLog(log contracts.LogAnswerUpdated) error {
	panic("@@TODO: unimplemented")
	return nil
}

func (p *PollingDeviationChecker) respondToRoundDetailsUpdatedLog(log contracts.LogRoundDetailsUpdated) error {
	p.roundTimeout = log.Timeout.Int64() * time.Second
	return nil
}

// respondToNewRoundLog takes the round broadcasted in the log event, and responds
// on-chain with an updated price.
// Only invoked by the CSP consumer on the single goroutine for thread safety.
func (p *PollingDeviationChecker) respondToNewRoundLog(log contracts.LogNewRound) error {
	requestedRound := log.RoundID

	jobSpecID := p.initr.JobSpecID.String()
	promSetBigInt(promFMSeenRound.WithLabelValues(jobSpecID), requestedRound)

	// skip if requested is not greater than current.
	if requestedRound.Cmp(p.currentRound) <= 0 {
		logger.Infow(
			fmt.Sprintf("Ignoring new round request: requested %s <= current %s", requestedRound, p.currentRound),
			"requestedRound", requestedRound,
			"currentRound", p.currentRound,
			"address", log.Address.Hex(),
			"jobID", p.initr.JobSpecID,
		)
		return nil
	}

	logger.Infow(
		fmt.Sprintf("Responding to new round request: requested %s > current %s", requestedRound, p.currentRound),
		"requestedRound", requestedRound,
		"currentRound", p.currentRound,
		"address", log.Address.Hex(),
		"jobID", p.initr.JobSpecID,
	)
	p.currentRound = requestedRound
	p.currentRoundStartedAt = log.StartedAt

	nextPrice, err := p.fetchPrices()
	if err != nil {
		return err
	}

	err = p.createJobRun(nextPrice, requestedRound)
	if err != nil {
		return err
	}

	p.currentPrice = nextPrice
	return nil
}

func (p *PollingDeviationChecker) pollIfEligible(client eth.Client, threshold float64) bool {
	open, err := p.isEligibleToPoll(client)
	logger.ErrorIf(err, "Unable to determine if round is open:")
	if !open {
		logger.Info("Round is currently not open to new submissions - polling paused")
		return false
	}
	ok, err := p.poll(threshold)
	logger.ErrorIf(err, "checker unable to poll")
	return ok
}

func (p *PollingDeviationChecker) isEligibleToPoll(client eth.Client) (bool, error) {
	roundIsOpen, err := p.isRoundOpen(client)
	if err != nil {
		return false, err
	} else if roundIsOpen {
		return true, nil
	}
	reportingRoundExpired, err := p.isRoundTimedOut(client)
	if err != nil {
		return false, err
	}
	return reportingRoundExpired, nil
}

func (p *PollingDeviationChecker) isRoundOpen(client eth.Client) (bool, error) {
	latestRound, err := p.fluxAggregator.LatestRound(p.address)
	if err != nil {
		return false, err
	}
	nodeAddress := p.store.KeyStore.Accounts()[0].Address
	_, lastRoundAnswered, err := p.fluxAggregator.LatestSubmission(p.address, nodeAddress)
	if err != nil {
		return false, err
	}
	roundIsOpen := lastRoundAnswered.Cmp(latestRound) <= 0
	return roundIsOpen, nil
}

func (p *PollingDeviationChecker) isRoundTimedOut(client eth.Client) (bool, error) {
	reportingRound, err := p.fluxAggregator.ReportingRound(p.address)
	if err != nil {
		return false, err
	}
	reportingRoundExpired, err := p.fluxAggregator.TimedOutStatus(p.address, reportingRound)
	if err != nil {
		return false, err
	}
	return reportingRoundExpired, nil
}

// fetchAggregatorData retrieves the price that's on-chain, with which we check
// the deviation against.
func (p *PollingDeviationChecker) fetchAggregatorData(client eth.Client) error {
	price, err := p.fluxAggregator.Price(p.address, p.precision)
	if err != nil {
		return err
	}
	p.currentPrice = price

	round, err := p.fluxAggregator.LatestRound(p.address)
	if err != nil {
		return err
	}
	p.currentRound = round
	return nil
}

// poll walks through the steps to check for a deviation, early exiting if deviation
// is not met, or triggering a new job run if deviation is met.
// Only invoked by the CSP consumer on the single goroutine for thread safety.
//
// True is returned when a Job Run was triggered.
func (p *PollingDeviationChecker) poll(threshold float64) (bool, error) {
	jobSpecID := p.initr.JobSpecID.String()

	nextPrice, err := p.fetchPrices()
	if err != nil {
		return false, err
	}
	promSetDecimal(promFMSeenValue.WithLabelValues(jobSpecID), nextPrice)

	if !OutsideDeviation(p.currentPrice, nextPrice, threshold) {
		return false, nil // early exit since deviation criteria not met.
	}

	nextRound := new(big.Int).Add(p.currentRound, big.NewInt(1)) // start new round
	logger.Infow("Detected change outside threshold, starting new round",
		"round", nextRound,
		"address", p.initr.Address.Hex(),
		"jobID", p.initr.JobSpecID,
	)
	err = p.createJobRun(nextPrice, nextRound)
	if err != nil {
		return false, err
	}

	p.currentPrice = nextPrice
	p.currentRound = nextRound

	promSetDecimal(promFMReportedValue.WithLabelValues(jobSpecID), p.currentPrice)
	promSetBigInt(promFMReportedRound.WithLabelValues(jobSpecID), p.currentRound)

	return true, nil
}

func (p *PollingDeviationChecker) fetchPrices() (decimal.Decimal, error) {
	median, err := p.fetcher.Fetch()
	return median, errors.Wrap(err, "unable to fetch median price")
}

func (p *PollingDeviationChecker) createJobRun(nextPrice decimal.Decimal, nextRound *big.Int) error {
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
		nextPrice,
		p.address.Hex(),
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
func OutsideDeviation(curPrice, nextPrice decimal.Decimal, threshold float64) bool {
	if curPrice.Equal(dec0) {
		logger.Infow("Current price is 0, deviation automatically met", "currentPrice", dec0)
		return true
	}

	diff := curPrice.Sub(nextPrice).Abs()
	percentage := diff.Div(curPrice).Mul(decimal.NewFromInt(100))
	if percentage.LessThan(decimal.NewFromFloat(threshold)) {
		logger.Debugw(
			"Deviation threshold not met",
			"difference", percentage,
			"threshold", threshold,
			"currentPrice", curPrice,
			"nextPrice", nextPrice)
		return false
	}
	logger.Infow(
		"Deviation threshold met",
		"difference", percentage,
		"threshold", threshold,
		"currentPrice", curPrice,
		"nextPrice", nextPrice,
	)
	return true
}
