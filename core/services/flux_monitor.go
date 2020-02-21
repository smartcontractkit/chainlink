package services

import (
	"chainlink/core/eth"
	"chainlink/core/logger"
	"chainlink/core/store"
	"chainlink/core/store/models"
	"chainlink/core/store/orm"
	"chainlink/core/utils"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"net/url"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"go.uber.org/multierr"
)

//go:generate mockery -name FluxMonitor -output ../internal/mocks/ -case=underscore

// defaultHTTPTimeout is the timeout used by the price adapter fetcher for outgoing HTTP requests.
const defaultHTTPTimeout = 5 * time.Second

// MinimumPollingInterval is the smallest possible polling interval the Flux
// Monitor supports.
const MinimumPollingInterval = models.Duration(defaultHTTPTimeout)

// FluxMonitor is the interface encapsulating all functionality
// needed to listen to price deviations and new round requests.
type FluxMonitor interface {
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

// NewFluxMonitor creates a service that manages a collection of DeviationCheckers,
// one per initiator of type InitiatorFluxMonitor for added jobs.
func NewFluxMonitor(store *store.Store, runManager RunManager) FluxMonitor {
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

//go:generate mockery -name DeviationCheckerFactory -output ../internal/mocks/ -case=underscore

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

//go:generate mockery -name DeviationChecker -output ../internal/mocks/ -case=underscore

// DeviationChecker encapsulate methods needed to initialize and check prices
// for price deviations.
type DeviationChecker interface {
	Start(context.Context, eth.Client) error
	Stop()
}

// PollingDeviationChecker polls external price adapters via HTTP to check for price swings.
type PollingDeviationChecker struct {
	store         *store.Store
	initr         models.Initiator
	address       common.Address
	requestData   models.JSON
	idleThreshold time.Duration
	threshold     float64
	precision     int32
	runManager    RunManager
	currentPrice  decimal.Decimal
	currentRound  *big.Int
	fetcher       Fetcher
	delay         time.Duration
	cancel        context.CancelFunc
	newRounds     chan eth.Log

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
		store:         store,
		initr:         initr,
		address:       initr.InitiatorParams.Address,
		requestData:   initr.InitiatorParams.RequestData,
		idleThreshold: initr.InitiatorParams.IdleThreshold.Duration(),
		threshold:     float64(initr.InitiatorParams.Threshold),
		precision:     initr.InitiatorParams.Precision,
		runManager:    runManager,
		currentPrice:  decimal.NewFromInt(0),
		currentRound:  big.NewInt(0),
		fetcher:       fetcher,
		delay:         delay,
		newRounds:     make(chan eth.Log),

		waitOnStop: make(chan struct{}),
	}, nil
}

// Start begins the CSP consumer in a single goroutine to
// poll the price adapters and listen to NewRound events.
func (p *PollingDeviationChecker) Start(ctx context.Context, client eth.Client) error {
	logger.Debugw("Starting checker for job",
		"job", p.initr.JobSpecID.String(),
		"initr", p.initr.ID)
	err := p.fetchAggregatorData(client)
	if err != nil {
		return err
	}

	roundSubscription, err := p.subscribeToNewRounds(client)
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

func (p *PollingDeviationChecker) consume(ctx context.Context, roundSubscription eth.Subscription, client eth.Client) {
	defer roundSubscription.Unsubscribe()

	idleThreshold := p.idleThreshold
	if idleThreshold == 0 {
		idleThreshold = math.MaxInt64
	}

	idleThresholdTimer := time.NewTimer(idleThreshold)
	defer stopTimer(idleThresholdTimer)

	for {
		jobRunTriggered := false

		select {
		case <-ctx.Done():
			close(p.waitOnStop)
			return
		case err := <-roundSubscription.Err():
			logger.Error(errors.Wrap(err, "checker lost subscription to NewRound log events"))
		case log := <-p.newRounds:
			err := p.respondToNewRound(log)
			logger.ErrorIf(err, "checker unable to respond to new round")
		case <-time.After(p.delay):
			jobRunTriggered = p.pollIFRoundOpen(client)
		case <-idleThresholdTimer.C:
			ok, err := p.poll(0)
			logger.ErrorIf(err, "checker unable to poll")
			jobRunTriggered = ok
		}

		if jobRunTriggered {
			// Reset expects stopped or expired timer.
			stopTimer(idleThresholdTimer)
			idleThresholdTimer.Reset(idleThreshold)
		}
	}
}

func (p *PollingDeviationChecker) pollIFRoundOpen(client eth.Client) bool {
	open, err := p.isRoundOpen(client)
	logger.ErrorIf(err, "Unable to determine if round is open:")
	if !open {
		logger.Info("Round is currently not open to new submissions - polling paused")
		return false
	}
	ok, err := p.poll(p.threshold)
	logger.ErrorIf(err, "checker unable to poll")
	return ok
}

func (p *PollingDeviationChecker) isRoundOpen(client eth.Client) (bool, error) {
	latestRound, err := client.GetAggregatorRound(p.address)
	if err != nil {
		return false, err
	}
	nodeAddress := p.store.KeyStore.Accounts()[0].Address
	_, lastRoundAnswered, err := client.GetLatestSubmission(p.address, nodeAddress)
	if err != nil {
		return false, err
	}
	return lastRoundAnswered.Cmp(latestRound) <= 0, nil
}

// Stop stops this instance from polling, cleaning up resources.
func (p *PollingDeviationChecker) Stop() {
	if p.cancel != nil {
		p.cancel()
		<-p.waitOnStop
	}
}

// fetchAggregatorData retrieves the price that's on-chain, with which we check
// the deviation against.
func (p *PollingDeviationChecker) fetchAggregatorData(client eth.Client) error {
	price, err := client.GetAggregatorPrice(p.address, p.precision)
	if err != nil {
		return err
	}
	p.currentPrice = price

	round, err := client.GetAggregatorRound(p.address)
	if err != nil {
		return err
	}
	p.currentRound = round
	return nil
}

func (p *PollingDeviationChecker) subscribeToNewRounds(client eth.Client) (eth.Subscription, error) {
	filterQuery, err := models.FilterQueryFactory(p.initr, nil)
	if err != nil {
		return nil, err
	}

	subscription, err := client.SubscribeToLogs(p.newRounds, filterQuery)
	if err != nil {
		return nil, err
	}

	logger.Infow(
		"Flux Monitor Initiator subscribing to new rounds",
		"address", p.initr.Address.Hex())
	return subscription, nil
}

// respondToNewRound takes the round broadcasted in the log event, and responds
// on-chain with an updated price.
// Only invoked by the CSP consumer on the single goroutine for thread safety.
func (p *PollingDeviationChecker) respondToNewRound(log eth.Log) error {
	requestedRound, err := models.ParseNewRoundLog(log)
	if err != nil {
		return err
	}

	jobSpecID := p.initr.JobSpecID.String()
	promSetBigInt(promFMSeenRound.WithLabelValues(jobSpecID), requestedRound)

	// skip if requested is not greater than current.
	if requestedRound.Cmp(p.currentRound) < 1 {
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
	aggregatorContract, err := eth.GetV6Contract(eth.PrepaidAggregatorName)
	if err != nil {
		return err
	}
	methodID, err := aggregatorContract.GetMethodID("updateAnswer")
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
