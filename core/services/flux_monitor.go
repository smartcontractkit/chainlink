package services

import (
	"chainlink/core/eth"
	"chainlink/core/logger"
	"chainlink/core/store"
	"chainlink/core/store/models"
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
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
}

type concreteFluxMonitor struct {
	store          *store.Store
	runManager     RunManager
	checkers       map[string][]DeviationChecker
	checkerFactory DeviationCheckerFactory
	adds           chan addEntry
	removes        chan *models.ID
	ctx            context.Context
	cancel         context.CancelFunc
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

// Connect adds all persisted jobs and starts deviation checkers for each
// flux monitor initiator.
func (fm *concreteFluxMonitor) Connect(*models.Head) error {
	fm.ctx, fm.cancel = context.WithCancel(context.Background())
	fm.adds = make(chan addEntry)
	fm.removes = make(chan *models.ID)

	go fm.actionConsumer(fm.ctx, fm.adds, fm.removes) // start single goroutine consumer

	// enqueue addJob actions
	rchan := make(chan error, 1)
	count := 0
	err := fm.store.Jobs(func(j *models.JobSpec) bool { // improve scoping of sql query
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

// actionConsumer is run on the single goroutine to coordinate the
// collection of DeviationCheckers in a thread safe fashion.
// Deliberately without shared variables besides channels and a context, all
// thread safe.
func (fm *concreteFluxMonitor) actionConsumer(ctx context.Context, adds chan addEntry, removes chan *models.ID) {
	fm.checkers = map[string][]DeviationChecker{}
	for {
		select {
		case <-ctx.Done():
			return
		case entry := <-adds:
			entry.rchan <- fm.produceJobAction(entry.job)
		case jobID := <-removes:
			for _, checker := range fm.checkers[jobID.String()] {
				checker.Stop()
			}
			delete(fm.checkers, jobID.String())
		}
	}
}

// Disconnect cleans up running deviation checkers.
func (fm *concreteFluxMonitor) Disconnect() {
	fm.cancel()
}

// OnNewHead is a noop.
func (fm *concreteFluxMonitor) OnNewHead(*models.Head) {}

// AddJob created a DeviationChecker for any job initiators of type
// InitiatorFluxMonitor.
func (fm *concreteFluxMonitor) AddJob(job models.JobSpec) error {
	if !job.IsFluxMonitorInitiated() {
		return nil
	}

	// non-blocking send is ignored if actionConsumer isn't consuming,
	// such as when disconnected.
	rchan := make(chan error)
	select {
	case fm.adds <- addEntry{&job, rchan}:
		return <-rchan
	default:
		return fmt.Errorf("unable to add job %s to flux monitor, flux monitor disconnected", job.ID.String())
	}
}

func (fm *concreteFluxMonitor) produceJobAction(job *models.JobSpec) error {
	if _, ok := fm.checkers[job.ID.String()]; ok {
		return fmt.Errorf("job %s has already been added to flux monitor", job.ID.String())
	}
	validCheckers := []DeviationChecker{}
	for _, initr := range job.InitiatorsFor(models.InitiatorFluxMonitor) {
		checker, err := fm.checkerFactory.New(fm.ctx, initr, fm.runManager)
		if err != nil {
			return err
		}
		err = checker.Initialize(fm.store.TxManager)
		if err != nil {
			return err
		}
		validCheckers = append(validCheckers, checker)
	}
	for _, checker := range validCheckers {
		go checker.Start()
	}
	fm.checkers[job.ID.String()] = validCheckers
	return nil
}

// RemoveJob stops and removes the checker for all Flux Monitor initiators belonging
// to the passed job ID.
func (fm *concreteFluxMonitor) RemoveJob(ID *models.ID) {
	// non-blocking send is ignored if actionConsumer isn't consuming,
	// such as when disconnected.
	select {
	case fm.removes <- ID:
	default:
	}
}

// DeviationCheckerFactory holds the New method needed to create a new instance
// of a DeviationChecker.
type DeviationCheckerFactory interface {
	New(context.Context, models.Initiator, RunManager) (DeviationChecker, error)
}

type pollingDeviationCheckerFactory struct{}

func (f pollingDeviationCheckerFactory) New(parentCtx context.Context, initr models.Initiator, runManager RunManager) (DeviationChecker, error) {
	fetcher, err := newMedianFetcherFromURLs(
		defaultHTTPTimeout,
		initr.InitiatorParams.RequestData.String(),
		initr.InitiatorParams.Feeds...)

	if err != nil {
		return nil, err
	}

	return NewPollingDeviationChecker(parentCtx, initr, runManager, fetcher, 1*time.Second)
}

// DeviationChecker encapsulate methods needed to initialize and check prices
// for deviations, or swings.
type DeviationChecker interface {
	Initialize(eth.Client) error
	Start()
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
	fetcher      Fetcher
	ctx          context.Context
	cancel       context.CancelFunc
	delay        time.Duration
}

// defaultHTTPTimeout is the timeout used by the price adapter fetcher for outgoing HTTP requests.
const defaultHTTPTimeout = 5 * time.Second

// NewPollingDeviationChecker returns a new instance of PollingDeviationChecker.
func NewPollingDeviationChecker(
	parentCtx context.Context,
	initr models.Initiator,
	runManager RunManager,
	fetcher Fetcher,
	delay time.Duration,
) (*PollingDeviationChecker, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	return &PollingDeviationChecker{
		initr:        initr,
		address:      initr.InitiatorParams.Address,
		requestData:  initr.InitiatorParams.RequestData,
		threshold:    float64(initr.InitiatorParams.Threshold),
		precision:    initr.InitiatorParams.Precision,
		runManager:   runManager,
		currentPrice: decimal.NewFromInt(0),
		ctx:          ctx,
		cancel:       cancel,
		fetcher:      fetcher,
		delay:        delay,
	}, nil
}

// Initialize retrieves the price that's on-chain, with which we must check
// the deviation from.
func (p *PollingDeviationChecker) Initialize(client eth.Client) error {
	price, err := client.GetAggregatorPrice(p.address, p.precision)
	p.currentPrice = price
	return err
}

// Start begins a loop polling the price adapters set in the InitiatorFluxMonitor.
func (p *PollingDeviationChecker) Start() {
	for {
		logger.ErrorIf(p.Poll())

		select {
		case <-p.ctx.Done():
			return
		case <-time.After(p.delay):
		}
	}
}

// Stop stops this instance from polling, cleaning up resources.
func (p *PollingDeviationChecker) Stop() {
	p.cancel()
}

// CurrentPrice returns the price used to check deviations against.
func (p *PollingDeviationChecker) CurrentPrice() decimal.Decimal {
	return p.currentPrice
}

// Poll walks through the steps to check for a deviation, early exiting if deviation
// is not met, or triggering a new job run if deviation is met.
func (p *PollingDeviationChecker) Poll() error {
	nextPrice, err := p.fetchPrices()
	if err != nil {
		return err
	}

	deviated, err := p.checkIfOutsideDeviation(nextPrice)
	if err != nil {
		return err
	}

	if !deviated {
		return nil // early exit since deviation criteria not met.
	}

	err = p.createJobRun(nextPrice)
	if err != nil {
		return err
	}

	p.currentPrice = nextPrice
	return nil
}

func (p *PollingDeviationChecker) fetchPrices() (decimal.Decimal, error) {
	median, err := p.fetcher.Fetch()
	return decimal.NewFromFloat(median), errors.Wrap(err, "unable to fetch median price")
}

func (p *PollingDeviationChecker) checkIfOutsideDeviation(nextPrice decimal.Decimal) (bool, error) {
	curPrice := p.currentPrice
	diff := curPrice.Sub(nextPrice).Abs()
	perc := diff.Div(curPrice).Mul(decimal.NewFromInt(100))
	if perc.LessThan(decimal.NewFromFloat(p.threshold)) {
		logger.Debug(fmt.Sprintf("difference is %v%%, deviation threshold of %f%% not met", perc, p.threshold))
		return false, nil
	}
	logger.Infow(
		fmt.Sprintf("deviation of %v%% for threshold %f%% with previous price %v next price %v", perc, p.threshold, curPrice, nextPrice),
		"threshold", p.threshold,
		"deviation", perc,
	)
	return true, nil
}

func (p *PollingDeviationChecker) createJobRun(nextPrice decimal.Decimal) error {
	runData, err := models.JSON{}.Add("result", fmt.Sprintf("%v", nextPrice))
	if err != nil {
		return errors.Wrap(err, "unable to start chainlink run")
	}
	_, err = p.runManager.Create(p.initr.JobSpecID, &p.initr, &runData, nil, models.NewRunRequest())
	return err
}
