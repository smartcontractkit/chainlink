package services

import (
	"chainlink/core/eth"
	"chainlink/core/logger"
	"chainlink/core/store"
	"chainlink/core/store/models"
	"chainlink/core/utils"
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"go.uber.org/multierr"
)

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
	store               *store.Store
	runManager          RunManager
	checkerFactory      DeviationCheckerFactory
	adds                chan addEntry
	removes             chan *models.ID
	connect, disconnect chan struct{}
	ctx                 context.Context
	cancel              context.CancelFunc
}

type addEntry struct {
	job   *models.JobSpec
	rchan chan error
}

// NewFluxMonitor creates a service that manages a collection of DeviationCheckers,
// one per initiator of type InitiatorFluxMonitor for added jobs.
func NewFluxMonitor(store *store.Store, runManager RunManager) FluxMonitor {
	return &concreteFluxMonitor{
		store:          store,
		runManager:     runManager,
		checkerFactory: pollingDeviationCheckerFactory{},
	}
}

func (fm *concreteFluxMonitor) Start() error {
	fm.ctx, fm.cancel = context.WithCancel(context.Background())
	fm.adds = make(chan addEntry)
	fm.removes = make(chan *models.ID)
	fm.connect = make(chan struct{})
	fm.disconnect = make(chan struct{})

	go fm.actionConsumer(fm.ctx) // start single goroutine consumer

	rchan := make(chan error, 1)
	count := 0
	err := fm.store.Jobs(func(j *models.JobSpec) bool { // add persisted jobs
		fm.adds <- addEntry{j, rchan}
		count++
		return true
	}, models.InitiatorFluxMonitor)

	// Block until jobs have been added, returning errors if any.
	var merr error
	for i := 0; i < count; i++ {
		merr = multierr.Combine(merr, <-rchan)
	}
	return multierr.Append(err, merr)
}

// Connect initializes all DeviationCheckers and starts their listening.
func (fm *concreteFluxMonitor) Connect(*models.Head) error {
	fm.connect <- struct{}{}
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
			entry.rchan <- fm.addAction(connectionCtx, connected, entry.job, jobMap)
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
	rchan := make(chan error)
	fm.adds <- addEntry{&job, rchan}
	return <-rchan
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

func (fm *concreteFluxMonitor) addAction(ctx context.Context, connected bool, job *models.JobSpec, jobMap map[string][]DeviationChecker) error {
	if _, ok := jobMap[job.ID.String()]; ok {
		return fmt.Errorf("job %s has already been added to flux monitor", job.ID.String())
	}
	validCheckers := []DeviationChecker{}
	for _, initr := range job.InitiatorsFor(models.InitiatorFluxMonitor) {
		checker, err := fm.checkerFactory.New(initr, fm.runManager)
		if err != nil {
			return errors.Wrap(err, "factory unable to create checker")
		}
		if connected {
			err := connectSingleChecker(ctx, checker, fm.store.TxManager)
			if err != nil {
				return errors.Wrap(err, "unable to connect checker")
			}
		}
		validCheckers = append(validCheckers, checker)
	}

	if len(validCheckers) > 0 {
		jobMap[job.ID.String()] = validCheckers
	}
	return nil
}

func connectSingleChecker(ctx context.Context, checker DeviationChecker, client eth.Client) error {
	err := checker.Initialize(client)
	if err != nil {
		return errors.Wrap(err, "unable to initialize flux monitor checker")
	}
	go checker.Start(ctx)
	return nil
}

// RemoveJob stops and removes the checker for all Flux Monitor initiators belonging
// to the passed job ID.
func (fm *concreteFluxMonitor) RemoveJob(ID *models.ID) {
	fm.removes <- ID
}

// DeviationCheckerFactory holds the New method needed to create a new instance
// of a DeviationChecker.
type DeviationCheckerFactory interface {
	New(models.Initiator, RunManager) (DeviationChecker, error)
}

type pollingDeviationCheckerFactory struct{}

func (f pollingDeviationCheckerFactory) New(initr models.Initiator, runManager RunManager) (DeviationChecker, error) {
	fetcher, err := newMedianFetcherFromURLs(
		defaultHTTPTimeout,
		initr.InitiatorParams.RequestData.String(),
		initr.InitiatorParams.Feeds...)

	if err != nil {
		return nil, err
	}

	return NewPollingDeviationChecker(initr, runManager, fetcher, 1*time.Minute)
}

// DeviationChecker encapsulate methods needed to initialize and check prices
// for deviations, or swings.
type DeviationChecker interface {
	Initialize(eth.Client) error
	Start(context.Context)
	Stop()
}

// PollingDeviationChecker polls external price adapters via HTTP to check for price swings.
type PollingDeviationChecker struct {
	initr        models.Initiator
	address      common.Address
	requestData  models.JSON
	threshold    float64
	precision    int32
	runManager   RunManager
	currentPrice decimal.Decimal
	currentRound *big.Int
	fetcher      Fetcher
	delay        time.Duration
	cancel       context.CancelFunc
}

// defaultHTTPTimeout is the timeout used by the price adapter fetcher for outgoing HTTP requests.
const defaultHTTPTimeout = 5 * time.Second

// NewPollingDeviationChecker returns a new instance of PollingDeviationChecker.
func NewPollingDeviationChecker(
	initr models.Initiator,
	runManager RunManager,
	fetcher Fetcher,
	delay time.Duration,
) (*PollingDeviationChecker, error) {
	return &PollingDeviationChecker{
		initr:        initr,
		address:      initr.InitiatorParams.Address,
		requestData:  initr.InitiatorParams.RequestData,
		threshold:    float64(initr.InitiatorParams.Threshold),
		precision:    initr.InitiatorParams.Precision,
		runManager:   runManager,
		currentPrice: decimal.NewFromInt(0),
		currentRound: big.NewInt(0),
		fetcher:      fetcher,
		delay:        delay,
	}, nil
}

// Initialize retrieves the price that's on-chain, with which we must check
// the deviation from.
func (p *PollingDeviationChecker) Initialize(client eth.Client) error {
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

// Start begins a loop polling the price adapters set in the InitiatorFluxMonitor.
func (p *PollingDeviationChecker) Start(ctx context.Context) {
	pollingCtx, cancel := context.WithCancel(ctx)
	p.cancel = cancel
	for {
		logger.ErrorIf(p.Poll())

		select {
		case <-pollingCtx.Done():
			return
		case <-time.After(p.delay):
		}
	}
}

// Stop stops this instance from polling, cleaning up resources.
func (p *PollingDeviationChecker) Stop() {
	if p.cancel != nil {
		p.cancel()
	}
}

// CurrentPrice returns the price used to check deviations against.
func (p *PollingDeviationChecker) CurrentPrice() decimal.Decimal {
	return p.currentPrice
}

// CurrentRound returns the latest round.
func (p *PollingDeviationChecker) CurrentRound() *big.Int {
	return new(big.Int).Set(p.currentRound)
}

// Poll walks through the steps to check for a deviation, early exiting if deviation
// is not met, or triggering a new job run if deviation is met.
func (p *PollingDeviationChecker) Poll() error {
	nextPrice, err := p.fetchPrices()
	if err != nil {
		return err
	}

	if !OutsideDeviation(p.currentPrice, nextPrice, p.threshold) {
		return nil // early exit since deviation criteria not met.
	}

	nextRound := new(big.Int).Add(p.currentRound, big.NewInt(1))
	err = p.createJobRun(nextPrice, nextRound)
	if err != nil {
		return err
	}

	p.currentPrice = nextPrice
	p.currentRound = nextRound
	return nil
}

func (p *PollingDeviationChecker) fetchPrices() (decimal.Decimal, error) {
	median, err := p.fetcher.Fetch()
	return median, errors.Wrap(err, "unable to fetch median price")
}

var dec0 = decimal.NewFromInt(0)

// OutsideDeviation checks whether the next price is outside the threshold.
func OutsideDeviation(curPrice, nextPrice decimal.Decimal, threshold float64) bool {
	if curPrice.Equal(dec0) {
		logger.Infow("current price is 0, deviation automatically met")
		return true
	}

	diff := curPrice.Sub(nextPrice).Abs()
	percentage := diff.Div(curPrice).Mul(decimal.NewFromInt(100))
	if percentage.LessThan(decimal.NewFromFloat(threshold)) {
		logger.Debug(fmt.Sprintf("difference is %v%%, deviation threshold of %f%% not met", percentage, threshold))
		return false
	}
	logger.Infow(
		fmt.Sprintf("deviation of %v%% for threshold %f%% with previous price %v next price %v", percentage, threshold, curPrice, nextPrice),
		"threshold", threshold,
		"deviation", percentage,
	)
	return true
}

func (p *PollingDeviationChecker) createJobRun(nextPrice decimal.Decimal, nextRound *big.Int) error {
	aggregatorContract, err := eth.GetV5Contract(eth.PrepaidAggregatorName)
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
		nextPrice.String(),
		p.address.Hex(),
		hexutil.Encode(methodID),
		hexutil.Encode(nextRoundData))

	runData, err := models.ParseJSON([]byte(payload))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("unable to start chainlink run with payload %s", payload))
	}
	_, err = p.runManager.Create(p.initr.JobSpecID, &p.initr, &runData, nil, models.NewRunRequest())
	return err
}
